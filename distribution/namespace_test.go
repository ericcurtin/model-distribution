package distribution

import (
	"os"
	"testing"

	"github.com/docker/model-distribution/internal/gguf"
)

func TestAddDefaultNamespace(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
		desc     string
	}{
		{
			input:    "gpt-oss",
			expected: "ai/gpt-oss",
			desc:     "simple model name gets ai/ prefix",
		},
		{
			input:    "ai/gpt-oss",
			expected: "ai/gpt-oss",
			desc:     "model with ai/ namespace stays unchanged",
		},
		{
			input:    "my-namespace/model",
			expected: "my-namespace/model",
			desc:     "model with custom namespace stays unchanged",
		},
		{
			input:    "some-repo:some-tag",
			expected: "some-repo:some-tag",
			desc:     "repo:tag format stays unchanged",
		},
		{
			input:    "registry.com/namespace/model:tag",
			expected: "registry.com/namespace/model:tag",
			desc:     "full registry reference stays unchanged",
		},
		{
			input:    "sha256:abcd1234567890",
			expected: "sha256:abcd1234567890",
			desc:     "SHA256 digest stays unchanged",
		},
		{
			input:    "llama2",
			expected: "ai/llama2",
			desc:     "another simple model name gets ai/ prefix",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			result := addDefaultNamespace(tc.input)
			if result != tc.expected {
				t.Errorf("addDefaultNamespace(%q) = %q, want %q", tc.input, result, tc.expected)
			}
		})
	}
}

func TestDefaultNamespaceIntegration(t *testing.T) {
	// Create temp directory for store
	tempDir, err := os.MkdirTemp("", "model-distribution-namespace-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create client
	client, err := NewClient(WithStoreRootPath(tempDir))
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Create a test model
	model, err := gguf.NewModel(testGGUFFile)
	if err != nil {
		t.Fatalf("Failed to create model: %v", err)
	}

	// Store model with the "ai/" namespace explicitly
	explicitTag := "ai/test-model"
	if err := client.store.Write(model, []string{explicitTag}, nil); err != nil {
		t.Fatalf("Failed to write model to store: %v", err)
	}

	t.Run("GetModel with default namespace", func(t *testing.T) {
		// Try to get the model using just the simple name (should add ai/ prefix)
		simpleRef := "test-model"
		retrievedModel, err := client.GetModel(simpleRef)
		if err != nil {
			t.Fatalf("Failed to get model with simple reference %q: %v", simpleRef, err)
		}

		// Verify it's the same model
		if len(retrievedModel.Tags()) == 0 {
			t.Fatal("Retrieved model has no tags")
		}

		found := false
		for _, tag := range retrievedModel.Tags() {
			if tag == explicitTag {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Retrieved model tags %v don't include expected tag %q", retrievedModel.Tags(), explicitTag)
		}
	})

	t.Run("IsModelInStore with default namespace", func(t *testing.T) {
		// Check if model exists using simple name
		simpleRef := "test-model"
		exists, err := client.IsModelInStore(simpleRef)
		if err != nil {
			t.Fatalf("Failed to check if model exists: %v", err)
		}
		if !exists {
			t.Errorf("Model should exist when referenced as %q", simpleRef)
		}

		// Check with explicit namespace
		exists, err = client.IsModelInStore(explicitTag)
		if err != nil {
			t.Fatalf("Failed to check if model exists: %v", err)
		}
		if !exists {
			t.Errorf("Model should exist when referenced as %q", explicitTag)
		}
	})

	t.Run("DeleteModel with default namespace", func(t *testing.T) {
		// Create another model for deletion test
		deleteTag := "ai/delete-model"
		if err := client.store.Write(model, []string{deleteTag}, nil); err != nil {
			t.Fatalf("Failed to write model to store: %v", err)
		}

		// Delete using simple name
		simpleRef := "delete-model"
		_, err := client.DeleteModel(simpleRef, false)
		if err != nil {
			t.Fatalf("Failed to delete model with simple reference %q: %v", simpleRef, err)
		}

		// Verify it's deleted
		exists, err := client.IsModelInStore(deleteTag)
		if err != nil {
			t.Fatalf("Failed to check if model exists after deletion: %v", err)
		}
		if exists {
			t.Errorf("Model should not exist after deletion")
		}
	})

	t.Run("No namespace added to existing qualified references", func(t *testing.T) {
		// Store model with custom namespace
		customTag := "custom/model:v1"
		if err := client.store.Write(model, []string{customTag}, nil); err != nil {
			t.Fatalf("Failed to write model to store: %v", err)
		}

		// Try to get with the same reference (should not add ai/ prefix)
		retrievedModel, err := client.GetModel(customTag)
		if err != nil {
			t.Fatalf("Failed to get model with qualified reference %q: %v", customTag, err)
		}

		// Verify it has the correct tag
		found := false
		for _, tag := range retrievedModel.Tags() {
			if tag == customTag {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Retrieved model tags %v don't include expected tag %q", retrievedModel.Tags(), customTag)
		}
	})
}