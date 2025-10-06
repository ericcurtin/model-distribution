package distribution

import (
	"testing"
)

func TestParseReference(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		wantErr  bool
	}{
		{
			name:     "familiar name",
			input:    "llama",
			expected: "docker.io/library/llama:latest",
			wantErr:  false,
		},
		{
			name:     "familiar name with tag",
			input:    "llama:7b",
			expected: "docker.io/library/llama:7b",
			wantErr:  false,
		},
		{
			name:     "full reference",
			input:    "docker.io/library/llama:7b",
			expected: "docker.io/library/llama:7b",
			wantErr:  false,
		},
		{
			name:     "custom domain",
			input:    "myregistry.com/models/llama:7b",
			expected: "myregistry.com/models/llama:7b",
			wantErr:  false,
		},
		{
			name:     "localhost",
			input:    "localhost/models/llama:7b",
			expected: "localhost/models/llama:7b",
			wantErr:  false,
		},
		{
			name:     "empty input",
			input:    "",
			expected: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseReference(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseReference() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && result.String() != tt.expected {
				t.Errorf("ParseReference() = %v, want %v", result.String(), tt.expected)
			}
		})
	}
}

func TestParseReferenceWithConfig(t *testing.T) {
	customConfig := ReferenceConfig{
		DefaultDomain:      "models.ai",
		OfficialRepoPrefix: "official/",
		DefaultTag:         "v1.0",
	}

	tests := []struct {
		name     string
		input    string
		config   ReferenceConfig
		expected string
		wantErr  bool
	}{
		{
			name:     "custom config familiar name",
			input:    "llama",
			config:   customConfig,
			expected: "models.ai/official/llama:v1.0",
			wantErr:  false,
		},
		{
			name:     "custom config with tag",
			input:    "llama:7b",
			config:   customConfig,
			expected: "models.ai/official/llama:7b",
			wantErr:  false,
		},
		{
			name:     "default config",
			input:    "llama",
			config:   DefaultReferenceConfig,
			expected: "docker.io/library/llama:latest",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseReferenceWithConfig(tt.input, tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseReferenceWithConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && result.String() != tt.expected {
				t.Errorf("ParseReferenceWithConfig() = %v, want %v", result.String(), tt.expected)
			}
		})
	}
}

func TestFamiliarName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "official image",
			input:    "docker.io/library/llama:latest",
			expected: "llama",
		},
		{
			name:     "custom namespace",
			input:    "docker.io/myorg/mymodel:v1.0",
			expected: "myorg/mymodel",
		},
		{
			name:     "custom domain",
			input:    "myregistry.com/models/llama:7b",
			expected: "myregistry.com/models/llama",
		},
		{
			name:     "localhost",
			input:    "localhost/models/llama:latest",
			expected: "localhost/models/llama",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ref, err := ParseReference(tt.input)
			if err != nil {
				t.Fatalf("ParseReference() failed: %v", err)
			}
			
			result := FamiliarName(ref)
			if result != tt.expected {
				t.Errorf("FamiliarName() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestFamiliarNameWithConfig(t *testing.T) {
	customConfig := ReferenceConfig{
		DefaultDomain:      "models.ai",
		OfficialRepoPrefix: "official/",
		DefaultTag:         "v1.0",
	}

	tests := []struct {
		name     string
		input    string
		config   ReferenceConfig
		expected string
	}{
		{
			name:     "custom config official image",
			input:    "models.ai/official/llama:v1.0",
			config:   customConfig,
			expected: "llama",
		},
		{
			name:     "custom config custom namespace",
			input:    "models.ai/myorg/mymodel:v1.0",
			config:   customConfig,
			expected: "myorg/mymodel",
		},
		{
			name:     "default config",
			input:    "docker.io/library/llama:latest",
			config:   DefaultReferenceConfig,
			expected: "llama",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ref, err := ParseReferenceWithConfig(tt.input, tt.config)
			if err != nil {
				t.Fatalf("ParseReferenceWithConfig() failed: %v", err)
			}
			
			result := FamiliarNameWithConfig(ref, tt.config)
			if result != tt.expected {
				t.Errorf("FamiliarNameWithConfig() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGlobalConfigurationOverrides(t *testing.T) {
	// Save original values
	originalDefault := DefaultReferenceConfig.DefaultDomain
	originalPrefix := DefaultReferenceConfig.OfficialRepoPrefix
	originalTag := DefaultReferenceConfig.DefaultTag

	// Override defaults
	SetDefaultDomain("custom.ai")
	SetOfficialRepoPrefix("models/")
	SetDefaultTag("v2.0")

	// Test that the overrides work
	if DefaultReferenceConfig.DefaultDomain != "custom.ai" {
		t.Errorf("SetDefaultDomain() did not update DefaultReferenceConfig.DefaultDomain")
	}
	if DefaultReferenceConfig.OfficialRepoPrefix != "models/" {
		t.Errorf("SetOfficialRepoPrefix() did not update DefaultReferenceConfig.OfficialRepoPrefix")
	}
	if DefaultReferenceConfig.DefaultTag != "v2.0" {
		t.Errorf("SetDefaultTag() did not update DefaultReferenceConfig.DefaultTag")
	}

	// Test parsing with the new defaults
	result, err := ParseReference("llama")
	if err != nil {
		t.Fatalf("ParseReference() with overridden config failed: %v", err)
	}
	expected := "custom.ai/models/llama:v2.0"
	if result.String() != expected {
		t.Errorf("ParseReference() with overridden config = %v, want %v", result.String(), expected)
	}

	// Test familiar name with the new defaults
	familiar := FamiliarName(result)
	if familiar != "llama" {
		t.Errorf("FamiliarName() with overridden config = %v, want %v", familiar, "llama")
	}

	// Restore original values
	SetDefaultDomain(originalDefault)
	SetOfficialRepoPrefix(originalPrefix)
	SetDefaultTag(originalTag)
}

func TestNamedInterface(t *testing.T) {
	ref, err := ParseReference("docker.io/library/llama:7b")
	if err != nil {
		t.Fatalf("ParseReference() failed: %v", err)
	}

	if ref.Name() != "docker.io/library/llama" {
		t.Errorf("Name() = %v, want %v", ref.Name(), "docker.io/library/llama")
	}
	if ref.Domain() != "docker.io" {
		t.Errorf("Domain() = %v, want %v", ref.Domain(), "docker.io")
	}
	if ref.Path() != "library/llama" {
		t.Errorf("Path() = %v, want %v", ref.Path(), "library/llama")
	}

	expected := "docker.io/library/llama:7b"
	if ref.String() != expected {
		t.Errorf("String() = %v, want %v", ref.String(), expected)
	}
}