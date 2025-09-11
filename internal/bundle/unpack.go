package bundle

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/docker/model-distribution/types"
)

// Unpack creates and return a Bundle by unpacking files and config from model into dir.
func Unpack(dir string, model types.Model) (*Bundle, error) {
	bundle := &Bundle{
		dir: dir,
	}
	if err := unpackModelFiles(bundle, model); err != nil {
		return nil, fmt.Errorf("add model file(s) to runtime bundle: %w", err)
	}
	if err := unpackMultiModalProjector(bundle, model); err != nil {
		return nil, fmt.Errorf("add multi-model projector file to runtime bundle: %w", err)
	}
	if err := unpackRuntimeConfig(bundle, model); err != nil {
		return nil, fmt.Errorf("add config.json to runtime bundle: %w", err)
	}
	return bundle, nil
}

func unpackRuntimeConfig(bundle *Bundle, mdl types.Model) error {
	cfg, err := mdl.Config()
	if err != nil {
		return err
	}
	f, err := os.Create(filepath.Join(bundle.dir, "config.json"))
	if err != nil {
		return fmt.Errorf("create runtime config file: %w", err)
	}
	defer f.Close()
	if err := json.NewEncoder(f).Encode(cfg); err != nil {
		return fmt.Errorf("encode runtime config: %w", err)
	}
	bundle.runtimeConfig = cfg
	return nil
}

func unpackModelFiles(bundle *Bundle, mdl types.Model) error {
	// Try GGUF files first
	ggufPaths, err := mdl.GGUFPaths()
	if err == nil && len(ggufPaths) > 0 {
		return unpackGGUFs(bundle, ggufPaths)
	}
	
	// Try SafeTensors files
	safetensorsPaths, err := mdl.SafeTensorsPaths()
	if err == nil && len(safetensorsPaths) > 0 {
		return unpackSafeTensors(bundle, safetensorsPaths)
	}
	
	return fmt.Errorf("no supported model files found (GGUF or SafeTensors)")
}

func unpackGGUFs(bundle *Bundle, ggufPaths []string) error {
	if len(ggufPaths) == 1 {
		if err := unpackFile(filepath.Join(bundle.dir, "model.gguf"), ggufPaths[0]); err != nil {
			return err
		}
		bundle.modelFile = "model.gguf"
		return nil
	}

	for i := range ggufPaths {
		name := fmt.Sprintf("model-%05d-of-%05d.gguf", i+1, len(ggufPaths))
		if err := unpackFile(filepath.Join(bundle.dir, name), ggufPaths[i]); err != nil {
			return err
		}
		bundle.modelFile = name
	}

	return nil
}

func unpackSafeTensors(bundle *Bundle, safetensorsPaths []string) error {
	if len(safetensorsPaths) == 1 {
		if err := unpackFile(filepath.Join(bundle.dir, "model.safetensors"), safetensorsPaths[0]); err != nil {
			return err
		}
		bundle.modelFile = "model.safetensors"
		return nil
	}

	for i := range safetensorsPaths {
		name := fmt.Sprintf("model-%05d-of-%05d.safetensors", i+1, len(safetensorsPaths))
		if err := unpackFile(filepath.Join(bundle.dir, name), safetensorsPaths[i]); err != nil {
			return err
		}
		bundle.modelFile = name
	}

	return nil
}

func unpackMultiModalProjector(bundle *Bundle, mdl types.Model) error {
	path, err := mdl.MMPROJPath()
	if err != nil {
		return nil // no such file
	}
	if err = unpackFile(filepath.Join(bundle.dir, "model.mmproj"), path); err != nil {
		return err
	}
	bundle.mmprojPath = "model.mmproj"
	return nil
}

func unpackFile(bundlePath string, srcPath string) error {
	return os.Link(srcPath, bundlePath)
}
