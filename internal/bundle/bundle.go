package bundle

import (
	"path/filepath"

	"github.com/docker/model-distribution/types"
)

// Bundle represents a runtime bundle containing a model and runtime config
type Bundle struct {
	dir           string
	mmprojPath    string
	ggufFile      string // path to GGUF file (first shard when model is split among files)
	vllmFile      string // path to vLLM model file
	runtimeConfig types.Config
}

// RootDir return the path to the bundle root directory
func (b *Bundle) RootDir() string {
	return b.dir
}

// GGUFPath return the path to model GGUF file. If the model is sharded this will be the path to the first shard,
// containing metadata headers.
func (b *Bundle) GGUFPath() string {
	if b.ggufFile == "" {
		return ""
	}
	return filepath.Join(b.dir, b.ggufFile)
}

// VLLMPath returns the path to vLLM model file or "" if none is present.
func (b *Bundle) VLLMPath() string {
	if b.vllmFile == "" {
		return ""
	}
	return filepath.Join(b.dir, b.vllmFile)
}

// MMPROJPath returns the path to a multi-modal projector file or "" if none is present.
func (b *Bundle) MMPROJPath() string {
	if b.mmprojPath == "" {
		return ""
	}
	return filepath.Join(b.dir, b.mmprojPath)
}

// RuntimeConfig returns config that should be respected by the backend at runtime.
func (b *Bundle) RuntimeConfig() types.Config {
	return b.runtimeConfig
}
