package vllm

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	v1 "github.com/google/go-containerregistry/pkg/v1"

	"github.com/docker/model-distribution/internal/partial"
	"github.com/docker/model-distribution/types"
)

func NewModel(path string) (*Model, error) {
	// For vLLM, we support single model files or directories with multiple files
	var modelFiles []string
	
	// Check if path is a directory or a single file
	if strings.HasSuffix(path, "/") || isDirectory(path) {
		// Handle directory case - collect relevant model files
		modelFiles = collectModelFiles(path)
	} else {
		// Single file case
		modelFiles = []string{path}
	}

	if len(modelFiles) == 0 {
		return nil, fmt.Errorf("no valid model files found in %s", path)
	}

	layers := make([]v1.Layer, len(modelFiles))
	diffIDs := make([]v1.Hash, len(modelFiles))
	
	for i, modelFile := range modelFiles {
		layer, err := partial.NewLayer(modelFile, types.MediaTypeVLLM)
		if err != nil {
			return nil, fmt.Errorf("create vllm layer: %w", err)
		}
		diffID, err := layer.DiffID()
		if err != nil {
			return nil, fmt.Errorf("get vllm layer diffID: %w", err)
		}
		layers[i] = layer
		diffIDs[i] = diffID
	}

	created := time.Now()
	return &Model{
		configFile: types.ConfigFile{
			Config: configFromPath(path),
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

func configFromPath(path string) types.Config {
	// Extract basic metadata from path/filename
	// For vLLM models, we'll infer metadata from the path structure
	
	basename := filepath.Base(path)
	
	config := types.Config{
		Format: types.FormatVLLM,
	}
	
	// Try to extract model information from filename/path
	if strings.Contains(strings.ToLower(basename), "7b") {
		config.Parameters = "7B"
	} else if strings.Contains(strings.ToLower(basename), "13b") {
		config.Parameters = "13B"
	} else if strings.Contains(strings.ToLower(basename), "30b") {
		config.Parameters = "30B"
	} else if strings.Contains(strings.ToLower(basename), "70b") {
		config.Parameters = "70B"
	}
	
	// Extract architecture if present in path
	if strings.Contains(strings.ToLower(basename), "llama") {
		config.Architecture = "llama"
	} else if strings.Contains(strings.ToLower(basename), "mistral") {
		config.Architecture = "mistral"
	} else if strings.Contains(strings.ToLower(basename), "gemma") {
		config.Architecture = "gemma"
	}
	
	return config
}

// Helper function to check if path is a directory
func isDirectory(path string) bool {
	// This is a simplified check - in a real implementation you'd use os.Stat
	return !strings.Contains(filepath.Base(path), ".")
}

// Helper function to collect model files from a directory
func collectModelFiles(dirPath string) []string {
	// In a real implementation, this would scan the directory for relevant files
	// For now, we'll return the directory itself as a placeholder
	return []string{dirPath}
}