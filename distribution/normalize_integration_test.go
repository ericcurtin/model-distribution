package distribution

import (
	"testing"

	"github.com/docker/model-distribution/distribution/reference/files"
)

func TestClientNormalizesReferences(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple model name gets ai prefix",
			input:    "smollm2:135M-Q4_K_M",
			expected: "ai/smollm2:135M-Q4_K_M",
		},
		{
			name:     "already prefixed model unchanged",
			input:    "ai/smollm2:135M-Q4_K_M",
			expected: "ai/smollm2:135M-Q4_K_M",
		},
		{
			name:     "custom registry unchanged",
			input:    "registry.example.com/models/llama:v1.0",
			expected: "registry.example.com/models/llama:v1.0",
		},
		{
			name:     "custom namespace unchanged",
			input:    "huggingface/smollm2:135M-Q4_K_M",
			expected: "huggingface/smollm2:135M-Q4_K_M",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := files.NormalizeReference(tt.input)
			if result != tt.expected {
				t.Errorf("NormalizeReference(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}