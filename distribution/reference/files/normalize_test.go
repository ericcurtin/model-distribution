package files

import (
	"testing"
)

func TestNormalizeReference(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple model name with tag",
			input:    "smollm2:135M-Q4_K_M",
			expected: "ai/smollm2:135M-Q4_K_M",
		},
		{
			name:     "model name with latest tag",
			input:    "llama:latest",
			expected: "ai/llama:latest",
		},
		{
			name:     "model name without tag",
			input:    "smollm2",
			expected: "ai/smollm2",
		},
		{
			name:     "already has ai prefix",
			input:    "ai/smollm2:135M-Q4_K_M",
			expected: "ai/smollm2:135M-Q4_K_M",
		},
		{
			name:     "different namespace",
			input:    "huggingface/smollm2:135M-Q4_K_M",
			expected: "huggingface/smollm2:135M-Q4_K_M",
		},
		{
			name:     "full registry with namespace",
			input:    "registry.example.com/models/llama:v1.0",
			expected: "registry.example.com/models/llama:v1.0",
		},
		{
			name:     "docker hub library image",
			input:    "ubuntu:22.04",
			expected: "ai/ubuntu:22.04",
		},
		{
			name:     "digest reference",
			input:    "smollm2@sha256:abc123",
			expected: "ai/smollm2@sha256:abc123",
		},
		{
			name:     "digest reference with namespace",
			input:    "huggingface/smollm2@sha256:abc123",
			expected: "huggingface/smollm2@sha256:abc123",
		},
		{
			name:     "full registry with digest",
			input:    "registry.example.com/models/llama@sha256:def456",
			expected: "registry.example.com/models/llama@sha256:def456",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "localhost registry",
			input:    "localhost:5000/mymodel:v1",
			expected: "localhost:5000/mymodel:v1",
		},
		{
			name:     "registry with port",
			input:    "myregistry.com:8080/namespace/model:tag",
			expected: "myregistry.com:8080/namespace/model:tag",
		},
		{
			name:     "model with hyphen",
			input:    "my-model:1.0",
			expected: "ai/my-model:1.0",
		},
		{
			name:     "model with underscore",
			input:    "my_model:1.0",
			expected: "ai/my_model:1.0",
		},
		{
			name:     "model with version tag",
			input:    "model:v1.2.3",
			expected: "ai/model:v1.2.3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NormalizeReference(tt.input)
			if result != tt.expected {
				t.Errorf("NormalizeReference(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestNormalizeReferenceEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "malformed reference - too many colons",
			input:    "model:tag:extra",
			expected: "ai/model:tag:extra", // fallback behavior
		},
		{
			name:     "only colon",
			input:    ":",
			expected: "ai/:",
		},
		{
			name:     "only at sign",
			input:    "@",
			expected: "ai/@",
		},
		{
			name:     "just ai prefix",
			input:    "ai",
			expected: "ai/ai", // ai becomes ai/ai
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NormalizeReference(tt.input)
			if result != tt.expected {
				t.Errorf("NormalizeReference(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestOfficialRepoPrefixValue(t *testing.T) {
	// Test that the officialRepoPrefix constant has the expected value
	if officialRepoPrefix != "ai/" {
		t.Errorf("officialRepoPrefix = %q, want %q", officialRepoPrefix, "ai/")
	}
}