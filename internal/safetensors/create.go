package safetensors

import (
	"fmt"
	"time"

	v1 "github.com/google/go-containerregistry/pkg/v1"

	"github.com/docker/model-distribution/internal/partial"
	"github.com/docker/model-distribution/types"
)

func NewModel(path string) (*Model, error) {
	// For now, handle single safetensors files
	// Future enhancement could support multi-file safetensors models
	layers := make([]v1.Layer, 1)
	diffIDs := make([]v1.Hash, 1)
	
	layer, err := partial.NewLayer(path, types.MediaTypeSafeTensors)
	if err != nil {
		return nil, fmt.Errorf("create safetensors layer: %w", err)
	}
	diffID, err := layer.DiffID()
	if err != nil {
		return nil, fmt.Errorf("get safetensors layer diffID: %w", err)
	}
	layers[0] = layer
	diffIDs[0] = diffID

	created := time.Now()
	return &Model{
		configFile: types.ConfigFile{
			Config: configFromFile(path),
			Descriptor: types.Descriptor{
				Created: &created,
			},
			RootFS: v1.RootFS{
				Type:    "rootfs",
				DiffIDs: diffIDs,
			},
		},
		layers: layers,
	}, nil
}

func configFromFile(path string) types.Config {
	// For SafeTensors files, we provide basic metadata
	// Future enhancement could parse SafeTensors headers for more detailed metadata
	return types.Config{
		Format: types.FormatSafeTensors,
		// Default values that could be enhanced with actual parsing
		Parameters:   "",
		Architecture: "",
		Quantization: "",
		Size:         "",
	}
}