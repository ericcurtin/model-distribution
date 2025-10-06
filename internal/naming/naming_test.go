package naming

import (
	"testing"
)

func TestDefaultNamespace_ParseReference(t *testing.T) {
	tests := []struct {
		name               string
		defaultRegistry    string
		input              string
		expectedRegistry   string
		expectedRepository string
		expectedError      bool
	}{
		{
			name:               "no default namespace - standard behavior",
			defaultRegistry:    "",
			input:              "mymodel:latest",
			expectedRegistry:   "index.docker.io",
			expectedRepository: "library/mymodel",
			expectedError:      false,
		},
		{
			name:               "default registry applied to simple reference",
			defaultRegistry:    "registry.example.com",
			input:              "mymodel:latest",
			expectedRegistry:   "registry.example.com",
			expectedRepository: "mymodel",
			expectedError:      false,
		},
		{
			name:               "explicit registry preserved",
			defaultRegistry:    "registry.example.com",
			input:              "other.registry.com/user/mymodel:latest",
			expectedRegistry:   "other.registry.com",
			expectedRepository: "user/mymodel",
			expectedError:      false,
		},
		{
			name:               "localhost registry preserved",
			defaultRegistry:    "registry.example.com",
			input:              "localhost:5000/mymodel:latest",
			expectedRegistry:   "localhost:5000",
			expectedRepository: "mymodel",
			expectedError:      false,
		},
		{
			name:               "docker hub user namespace preserved with default",
			defaultRegistry:    "registry.example.com",
			input:              "user/mymodel:latest",
			expectedRegistry:   "registry.example.com",
			expectedRepository: "user/mymodel",
			expectedError:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dn := &DefaultNamespace{Registry: tt.defaultRegistry}
			ref, err := dn.ParseReference(tt.input)

			if tt.expectedError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if ref.Context().Registry.RegistryStr() != tt.expectedRegistry {
				t.Errorf("expected registry %s, got %s", tt.expectedRegistry, ref.Context().Registry.RegistryStr())
			}

			if ref.Context().RepositoryStr() != tt.expectedRepository {
				t.Errorf("expected repository %s, got %s", tt.expectedRepository, ref.Context().RepositoryStr())
			}
		})
	}
}

func TestDefaultNamespace_ParseTag(t *testing.T) {
	tests := []struct {
		name               string
		defaultRegistry    string
		input              string
		expectedRegistry   string
		expectedRepository string
		expectedError      bool
	}{
		{
			name:               "no default namespace - standard behavior",
			defaultRegistry:    "",
			input:              "mymodel:latest",
			expectedRegistry:   "index.docker.io",
			expectedRepository: "library/mymodel",
			expectedError:      false,
		},
		{
			name:               "default registry applied to simple tag",
			defaultRegistry:    "registry.example.com",
			input:              "mymodel:latest",
			expectedRegistry:   "registry.example.com",
			expectedRepository: "mymodel",
			expectedError:      false,
		},
		{
			name:               "explicit registry preserved in tag",
			defaultRegistry:    "registry.example.com",
			input:              "other.registry.com/user/mymodel:latest",
			expectedRegistry:   "other.registry.com",
			expectedRepository: "user/mymodel",
			expectedError:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dn := &DefaultNamespace{Registry: tt.defaultRegistry}
			tag, err := dn.ParseTag(tt.input)

			if tt.expectedError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if tag.Context().Registry.RegistryStr() != tt.expectedRegistry {
				t.Errorf("expected registry %s, got %s", tt.expectedRegistry, tag.Context().Registry.RegistryStr())
			}

			if tag.Context().RepositoryStr() != tt.expectedRepository {
				t.Errorf("expected repository %s, got %s", tt.expectedRepository, tag.Context().RepositoryStr())
			}
		})
	}
}

func TestHasExplicitRegistry(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"mymodel:latest", false},
		{"user/mymodel:latest", false},
		{"registry.example.com/mymodel:latest", true},
		{"localhost:5000/mymodel:latest", true},
		{"example.com/user/mymodel:latest", true},
		{"sub.domain.com/ns/mymodel:latest", true},
		{"localhost/mymodel:latest", false}, // localhost without port/dot after is not a registry
		{"model", false},
		{"registry.com", true},
		{"host:8080", true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := hasExplicitRegistry(tt.input)
			if result != tt.expected {
				t.Errorf("hasExplicitRegistry(%q) = %v, expected %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestGlobalDefaultNamespace(t *testing.T) {
	// Save original state
	original := globalDefaultNamespace
	defer func() {
		globalDefaultNamespace = original
	}()

	// Test setting and getting
	SetDefaultNamespace("test.registry.com")
	dn := GetDefaultNamespace()
	if dn == nil || dn.Registry != "test.registry.com" {
		t.Errorf("expected registry test.registry.com, got %v", dn)
	}

	// Test convenience functions
	ref, err := ParseReference("mymodel:latest")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ref.Context().Registry.RegistryStr() != "test.registry.com" {
		t.Errorf("expected registry test.registry.com, got %s", ref.Context().Registry.RegistryStr())
	}

	tag, err := ParseTag("mymodel:latest")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tag.Context().Registry.RegistryStr() != "test.registry.com" {
		t.Errorf("expected registry test.registry.com, got %s", tag.Context().Registry.RegistryStr())
	}
}

func TestNilDefaultNamespace(t *testing.T) {
	var dn *DefaultNamespace

	// Should fall back to standard behavior
	ref, err := dn.ParseReference("mymodel:latest")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ref.Context().Registry.RegistryStr() != "index.docker.io" {
		t.Errorf("expected Docker Hub default, got %s", ref.Context().Registry.RegistryStr())
	}

	tag, err := dn.ParseTag("mymodel:latest")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tag.Context().Registry.RegistryStr() != "index.docker.io" {
		t.Errorf("expected Docker Hub default, got %s", tag.Context().Registry.RegistryStr())
	}
}