package builder_test

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/docker/model-distribution/builder"
	"github.com/docker/model-distribution/types"
)

func TestWithMultimodalProjector(t *testing.T) {
	// Create a builder from a GGUF file
	b, err := builder.FromGGUF(filepath.Join("..", "assets", "dummy.gguf"))
	if err != nil {
		t.Fatalf("Failed to create builder from GGUF: %v", err)
	}

	// Add multimodal projector
	b2, err := b.WithMultimodalProjector(filepath.Join("..", "assets", "dummy.mmproj"))
	if err != nil {
		t.Fatalf("Failed to add multimodal projector: %v", err)
	}

	// Build the model
	target := &fakeTarget{}
	if err := b2.Build(t.Context(), target, nil); err != nil {
		t.Fatalf("Failed to build model: %v", err)
	}

	// Verify the model has the expected layers
	manifest, err := target.artifact.Manifest()
	if err != nil {
		t.Fatalf("Failed to get manifest: %v", err)
	}

	// Should have 2 layers: GGUF + multimodal projector
	if len(manifest.Layers) != 2 {
		t.Fatalf("Expected 2 layers, got %d", len(manifest.Layers))
	}

	// Check that one layer has the multimodal projector media type
	foundMMProjLayer := false
	for _, layer := range manifest.Layers {
		if layer.MediaType == types.MediaTypeMultimodalProjector {
			foundMMProjLayer = true
			break
		}
	}

	if !foundMMProjLayer {
		t.Error("Expected to find a layer with multimodal projector media type")
	}

	// Note: We can't directly test MMPROJPath() on ModelArtifact interface
	// but we can verify the layer was added with correct media type above
}

func TestWithMultimodalProjectorInvalidPath(t *testing.T) {
	// Create a builder from a GGUF file
	b, err := builder.FromGGUF(filepath.Join("..", "assets", "dummy.gguf"))
	if err != nil {
		t.Fatalf("Failed to create builder from GGUF: %v", err)
	}

	// Try to add multimodal projector with invalid path
	_, err = b.WithMultimodalProjector("nonexistent/path/to/mmproj")
	if err == nil {
		t.Error("Expected error when adding multimodal projector with invalid path")
	}
}

func TestWithMultimodalProjectorChaining(t *testing.T) {
	// Create a builder from a GGUF file
	b, err := builder.FromGGUF(filepath.Join("..", "assets", "dummy.gguf"))
	if err != nil {
		t.Fatalf("Failed to create builder from GGUF: %v", err)
	}

	// Chain multiple operations: license + multimodal projector + context size
	b, err = b.WithLicense(filepath.Join("..", "assets", "license.txt"))
	if err != nil {
		t.Fatalf("Failed to add license: %v", err)
	}

	b, err = b.WithMultimodalProjector(filepath.Join("..", "assets", "dummy.mmproj"))
	if err != nil {
		t.Fatalf("Failed to add multimodal projector: %v", err)
	}

	b = b.WithContextSize(4096)

	// Build the model
	target := &fakeTarget{}
	if err := b.Build(t.Context(), target, nil); err != nil {
		t.Fatalf("Failed to build model: %v", err)
	}

	// Verify the final model has all expected layers and properties
	manifest, err := target.artifact.Manifest()
	if err != nil {
		t.Fatalf("Failed to get manifest: %v", err)
	}

	// Should have 3 layers: GGUF + license + multimodal projector
	if len(manifest.Layers) != 3 {
		t.Fatalf("Expected 3 layers, got %d", len(manifest.Layers))
	}

	// Check media types - using string comparison since we can't use types.MediaType directly
	expectedMediaTypes := map[string]bool{
		string(types.MediaTypeGGUF):                false,
		string(types.MediaTypeLicense):             false,
		string(types.MediaTypeMultimodalProjector): false,
	}

	for _, layer := range manifest.Layers {
		if _, exists := expectedMediaTypes[string(layer.MediaType)]; exists {
			expectedMediaTypes[string(layer.MediaType)] = true
		}
	}

	for mediaType, found := range expectedMediaTypes {
		if !found {
			t.Errorf("Expected to find layer with media type %s", mediaType)
		}
	}

	// Check context size
	config, err := target.artifact.Config()
	if err != nil {
		t.Fatalf("Failed to get config: %v", err)
	}

	if config.ContextSize == nil || *config.ContextSize != 4096 {
		t.Errorf("Expected context size 4096, got %v", config.ContextSize)
	}

	// Note: We can't directly test GGUFPath() and MMPROJPath() on ModelArtifact interface
	// but we can verify the layers were added with correct media types above
}

func TestFromVLLM(t *testing.T) {
	// Create a temporary vLLM model file for testing
	tmpDir := t.TempDir()
	modelPath := filepath.Join(tmpDir, "test-model.safetensors")
	
	// Create a dummy model file
	if err := os.WriteFile(modelPath, []byte("dummy vllm model content"), 0644); err != nil {
		t.Fatalf("Failed to create test model file: %v", err)
	}

	// Create a builder from a vLLM file
	b, err := builder.FromVLLM(modelPath)
	if err != nil {
		t.Fatalf("Failed to create builder from vLLM: %v", err)
	}

	// Build the model
	target := &fakeTarget{}
	if err := b.Build(t.Context(), target, nil); err != nil {
		t.Fatalf("Failed to build model: %v", err)
	}

	// Verify the model has the expected format
	config, err := target.artifact.Config()
	if err != nil {
		t.Fatalf("Failed to get config: %v", err)
	}

	if config.Format != types.FormatVLLM {
		t.Errorf("Expected format %s, got %s", types.FormatVLLM, config.Format)
	}

	// Verify the model has one vLLM layer
	manifest, err := target.artifact.Manifest()
	if err != nil {
		t.Fatalf("Failed to get manifest: %v", err)
	}

	if len(manifest.Layers) != 1 {
		t.Fatalf("Expected 1 layer, got %d", len(manifest.Layers))
	}

	if manifest.Layers[0].MediaType != types.MediaTypeVLLM {
		t.Errorf("Expected layer media type %s, got %s", types.MediaTypeVLLM, manifest.Layers[0].MediaType)
	}
}

var _ builder.Target = &fakeTarget{}

type fakeTarget struct {
	artifact types.ModelArtifact
}

func (ft *fakeTarget) Write(ctx context.Context, artifact types.ModelArtifact, writer io.Writer) error {
	ft.artifact = artifact
	return nil
}
