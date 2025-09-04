package registry

import (
	"testing"
)

func TestWithDefaultNamespace(t *testing.T) {
	// Test client without default namespace - should use standard behavior
	client1 := NewClient()
	if client1.defaultNamespace != nil {
		t.Error("expected nil default namespace")
	}
	
	// Test client with default namespace
	client2 := NewClient(WithDefaultNamespace("registry.example.com"))
	if client2.defaultNamespace == nil {
		t.Fatal("expected non-nil default namespace")
	}
	if client2.defaultNamespace.Registry != "registry.example.com" {
		t.Errorf("expected registry.example.com, got %s", client2.defaultNamespace.Registry)
	}
	
	// Test empty namespace (should not set)
	client3 := NewClient(WithDefaultNamespace(""))
	if client3.defaultNamespace != nil {
		t.Error("expected nil default namespace for empty string")
	}
}

func TestClientNamespaceBehavior(t *testing.T) {
	tests := []struct {
		name             string
		defaultRegistry  string
		input            string
		expectedRegistry string
	}{
		{
			name:             "no default - uses Docker Hub",
			defaultRegistry:  "",
			input:            "mymodel:latest",
			expectedRegistry: "index.docker.io",
		},
		{
			name:             "custom default applied",
			defaultRegistry:  "registry.example.com",
			input:            "mymodel:latest",
			expectedRegistry: "registry.example.com",
		},
		{
			name:             "explicit registry preserved",
			defaultRegistry:  "registry.example.com",
			input:            "other.registry.com/mymodel:latest",
			expectedRegistry: "other.registry.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var client *Client
			if tt.defaultRegistry == "" {
				client = NewClient()
			} else {
				client = NewClient(WithDefaultNamespace(tt.defaultRegistry))
			}

			// Test Model method parsing
			ref, err := client.defaultNamespace.ParseReference(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if ref.Context().Registry.RegistryStr() != tt.expectedRegistry {
				t.Errorf("expected registry %s, got %s", tt.expectedRegistry, ref.Context().Registry.RegistryStr())
			}
		})
	}
}