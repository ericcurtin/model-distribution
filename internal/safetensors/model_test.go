package safetensors

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/docker/model-distribution/types"
)

func TestNewModel(t *testing.T) {
	// Create a temporary safetensors file for testing
	tmpDir := t.TempDir()
	safetensorsPath := filepath.Join(tmpDir, "test.safetensors")
	
	// Create a dummy safetensors file
	content := []byte("dummy safetensors content for testing")
	if err := os.WriteFile(safetensorsPath, content, 0644); err != nil {
		t.Fatalf("Failed to create test safetensors file: %v", err)
	}

	// Test creating a model from the safetensors file
	model, err := NewModel(safetensorsPath)
	if err != nil {
		t.Fatalf("Failed to create model: %v", err)
	}

	// Verify the model has the expected format
	config, err := model.Config()
	if err != nil {
		t.Fatalf("Failed to get model config: %v", err)
	}

	if config.Format != types.FormatSafeTensors {
		t.Errorf("Expected format %v, got %v", types.FormatSafeTensors, config.Format)
	}

	// Verify we can get layers
	layers, err := model.Layers()
	if err != nil {
		t.Fatalf("Failed to get layers: %v", err)
	}

	if len(layers) != 1 {
		t.Errorf("Expected 1 layer, got %d", len(layers))
	}

	// Verify the layer has the correct media type
	if len(layers) > 0 {
		mediaType, err := layers[0].MediaType()
		if err != nil {
			t.Fatalf("Failed to get layer media type: %v", err)
		}

		if mediaType != types.MediaTypeSafeTensors {
			t.Errorf("Expected media type %v, got %v", types.MediaTypeSafeTensors, mediaType)
		}
	}
}