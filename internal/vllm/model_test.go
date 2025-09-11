package vllm

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/docker/model-distribution/types"
)

func TestNewModel(t *testing.T) {
	// Create a temporary model file for testing
	tmpDir := t.TempDir()
	modelPath := filepath.Join(tmpDir, "test-model.safetensors")

	// Create a dummy model file
	if err := os.WriteFile(modelPath, []byte("dummy model content"), 0644); err != nil {
		t.Fatalf("Failed to create test model file: %v", err)
	}

	// Test creating a vLLM model
	model, err := NewModel(modelPath)
	if err != nil {
		t.Fatalf("Failed to create vLLM model: %v", err)
	}

	// Test that the model has the expected format
	config, err := model.Config()
	if err != nil {
		t.Fatalf("Failed to get model config: %v", err)
	}

	if config.Format != types.FormatVLLM {
		t.Errorf("Expected format %s, got %s", types.FormatVLLM, config.Format)
	}

	// Test that we can get layers
	layers, err := model.Layers()
	if err != nil {
		t.Fatalf("Failed to get model layers: %v", err)
	}

	if len(layers) != 1 {
		t.Errorf("Expected 1 layer, got %d", len(layers))
	}

	// Test VLLMPaths method
	vllmPaths, err := model.VLLMPaths()
	if err != nil {
		t.Fatalf("Failed to get vLLM paths: %v", err)
	}

	if len(vllmPaths) == 0 {
		t.Error("Expected at least one vLLM path")
	}

	// Test GGUFPaths method (should return empty for vLLM models)
	ggufPaths, err := model.GGUFPaths()
	if err != nil {
		t.Fatalf("Failed to get GGUF paths: %v", err)
	}

	if len(ggufPaths) != 0 {
		t.Errorf("Expected 0 GGUF paths for vLLM model, got %d", len(ggufPaths))
	}
}

func TestConfigFromPath(t *testing.T) {
	testCases := []struct {
		path                 string
		expectedParameters   string
		expectedArchitecture string
	}{
		{"/models/llama-7b.safetensors", "7B", "llama"},
		{"/models/mistral-13b.safetensors", "13B", "mistral"},
		{"/models/gemma-30b.safetensors", "30B", "gemma"},
		{"/models/some-70b.safetensors", "70B", ""},
		{"/models/unknown.safetensors", "", ""},
	}

	for _, tc := range testCases {
		t.Run(tc.path, func(t *testing.T) {
			config := configFromPath(tc.path)

			if config.Format != types.FormatVLLM {
				t.Errorf("Expected format %s, got %s", types.FormatVLLM, config.Format)
			}

			if config.Parameters != tc.expectedParameters {
				t.Errorf("Expected parameters %s, got %s", tc.expectedParameters, config.Parameters)
			}

			if config.Architecture != tc.expectedArchitecture {
				t.Errorf("Expected architecture %s, got %s", tc.expectedArchitecture, config.Architecture)
			}
		})
	}
}

func TestModelImplementsInterfaces(t *testing.T) {
	tmpDir := t.TempDir()
	modelPath := filepath.Join(tmpDir, "test-model.safetensors")

	if err := os.WriteFile(modelPath, []byte("dummy"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	model, err := NewModel(modelPath)
	if err != nil {
		t.Fatalf("Failed to create model: %v", err)
	}

	// Test that model implements types.ModelArtifact interface
	var _ types.ModelArtifact = model

	// Test basic interface methods
	id, err := model.ID()
	if err != nil {
		t.Errorf("Failed to get model ID: %v", err)
	}
	if id == "" {
		t.Error("Expected non-empty model ID")
	}

	descriptor, err := model.Descriptor()
	if err != nil {
		t.Errorf("Failed to get model descriptor: %v", err)
	}
	if descriptor.Created == nil {
		t.Error("Expected non-nil created time")
	}
	if time.Since(*descriptor.Created) > time.Minute {
		t.Error("Created time should be recent")
	}
}

func TestNewModelWithInvalidPath(t *testing.T) {
	// Test with non-existent path
	_, err := NewModel("/non/existent/path")
	if err == nil {
		t.Error("Expected error for non-existent path")
	}
}
