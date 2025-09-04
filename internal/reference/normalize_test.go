package reference

import (
	"testing"

	"github.com/opencontainers/go-digest"
)

func TestParseNormalizedNamed(t *testing.T) {
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
			name:     "with digest",
			input:    "llama@sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
			expected: "docker.io/library/llama@sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
			wantErr:  false,
		},
		{
			name:     "empty input",
			input:    "",
			expected: "",
			wantErr:  true,
		},
		{
			name:     "uppercase repo name",
			input:    "LLAMA",
			expected: "",
			wantErr:  true,
		},
		{
			name:     "invalid digest",
			input:    "llama@sha256:invalid",
			expected: "",
			wantErr:  true,
		},
		{
			name:     "ip address domain",
			input:    "192.168.1.1:5000/models/llama",
			expected: "192.168.1.1:5000/models/llama:latest",
			wantErr:  false,
		},
		{
			name:     "uppercase first segment treated as domain",
			input:    "MyRegistry/models/llama",
			expected: "MyRegistry/models/llama:latest",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseNormalizedNamed(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseNormalizedNamed() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && result.String() != tt.expected {
				t.Errorf("ParseNormalizedNamed() = %v, want %v", result.String(), tt.expected)
			}
		})
	}
}

func TestParseNormalizedNamedWithConfig(t *testing.T) {
	customConfig := Configuration{
		DefaultDomain:      "models.ai",
		OfficialRepoPrefix: "official/",
		DefaultTag:         "v1.0",
	}

	tests := []struct {
		name     string
		input    string
		config   Configuration
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
			config:   DefaultConfiguration,
			expected: "docker.io/library/llama:latest",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseNormalizedNamedWithConfig(tt.input, tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseNormalizedNamedWithConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && result.String() != tt.expected {
				t.Errorf("ParseNormalizedNamedWithConfig() = %v, want %v", result.String(), tt.expected)
			}
		})
	}
}

func TestSplitDomainWithConfig(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		config       Configuration
		expectedDom  string
		expectedRem  string
	}{
		{
			name:         "familiar name",
			input:        "llama",
			config:       DefaultConfiguration,
			expectedDom:  "docker.io",
			expectedRem:  "library/llama",
		},
		{
			name:         "full path",
			input:        "docker.io/library/llama",
			config:       DefaultConfiguration,
			expectedDom:  "docker.io",
			expectedRem:  "library/llama",
		},
		{
			name:         "localhost",
			input:        "localhost/models/llama",
			config:       DefaultConfiguration,
			expectedDom:  "localhost",
			expectedRem:  "models/llama",
		},
		{
			name:         "domain with port",
			input:        "registry.com:5000/models/llama",
			config:       DefaultConfiguration,
			expectedDom:  "registry.com:5000",
			expectedRem:  "models/llama",
		},
		{
			name:         "uppercase first segment",
			input:        "MyRegistry/models/llama",
			config:       DefaultConfiguration,
			expectedDom:  "MyRegistry",
			expectedRem:  "models/llama",
		},
		{
			name:         "default domain single segment",
			input:        "docker.io/llama",
			config:       DefaultConfiguration,
			expectedDom:  "docker.io",
			expectedRem:  "library/llama",
		},
		{
			name: "custom config",
			input: "modelname",
			config: Configuration{
				DefaultDomain:      "models.ai",
				OfficialRepoPrefix: "official/",
				DefaultTag:         "v1.0",
			},
			expectedDom: "models.ai",
			expectedRem: "official/modelname",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			domain, remoteName := splitDomainWithConfig(tt.input, tt.config)
			if domain != tt.expectedDom {
				t.Errorf("splitDomainWithConfig() domain = %v, want %v", domain, tt.expectedDom)
			}
			if remoteName != tt.expectedRem {
				t.Errorf("splitDomainWithConfig() remoteName = %v, want %v", remoteName, tt.expectedRem)
			}
		})
	}
}

func TestFamiliarName(t *testing.T) {
	tests := []struct {
		name     string
		ref      Named
		expected string
	}{
		{
			name: "official image",
			ref: namedReference{
				domain: "docker.io",
				path:   "library/llama",
				tag:    "latest",
			},
			expected: "llama",
		},
		{
			name: "custom namespace",
			ref: namedReference{
				domain: "docker.io",
				path:   "myorg/mymodel",
				tag:    "v1.0",
			},
			expected: "myorg/mymodel",
		},
		{
			name: "custom domain",
			ref: namedReference{
				domain: "myregistry.com",
				path:   "models/llama",
				tag:    "7b",
			},
			expected: "myregistry.com/models/llama",
		},
		{
			name: "localhost",
			ref: namedReference{
				domain: "localhost",
				path:   "models/llama",
				tag:    "latest",
			},
			expected: "localhost/models/llama",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FamiliarName(tt.ref)
			if result != tt.expected {
				t.Errorf("FamiliarName() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestFamiliarNameWithConfig(t *testing.T) {
	customConfig := Configuration{
		DefaultDomain:      "models.ai",
		OfficialRepoPrefix: "official/",
		DefaultTag:         "v1.0",
	}

	tests := []struct {
		name     string
		ref      Named
		config   Configuration
		expected string
	}{
		{
			name: "custom config official image",
			ref: namedReference{
				domain: "models.ai",
				path:   "official/llama",
				tag:    "v1.0",
			},
			config:   customConfig,
			expected: "llama",
		},
		{
			name: "custom config custom namespace",
			ref: namedReference{
				domain: "models.ai",
				path:   "myorg/mymodel",
				tag:    "v1.0",
			},
			config:   customConfig,
			expected: "myorg/mymodel",
		},
		{
			name: "default config",
			ref: namedReference{
				domain: "docker.io",
				path:   "library/llama",
				tag:    "latest",
			},
			config:   DefaultConfiguration,
			expected: "llama",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FamiliarNameWithConfig(tt.ref, tt.config)
			if result != tt.expected {
				t.Errorf("FamiliarNameWithConfig() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestIsNameOnly(t *testing.T) {
	tests := []struct {
		name     string
		ref      Named
		expected bool
	}{
		{
			name: "name only",
			ref: namedReference{
				domain: "docker.io",
				path:   "library/llama",
			},
			expected: true,
		},
		{
			name: "with tag",
			ref: namedReference{
				domain: "docker.io",
				path:   "library/llama",
				tag:    "latest",
			},
			expected: false,
		},
		{
			name: "with digest",
			ref: namedReference{
				domain: "docker.io",
				path:   "library/llama",
				digest: digest.Digest("sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"),
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsNameOnly(tt.ref)
			if result != tt.expected {
				t.Errorf("IsNameOnly() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestWithTag(t *testing.T) {
	baseRef := namedReference{
		domain: "docker.io",
		path:   "library/llama",
	}

	tests := []struct {
		name     string
		named    Named
		tag      string
		expected string
		wantErr  bool
	}{
		{
			name:     "valid tag",
			named:    baseRef,
			tag:      "v1.0",
			expected: "docker.io/library/llama:v1.0",
			wantErr:  false,
		},
		{
			name:     "empty tag",
			named:    baseRef,
			tag:      "",
			expected: "",
			wantErr:  true,
		},
		{
			name:     "tag with space",
			named:    baseRef,
			tag:      "v1 0",
			expected: "",
			wantErr:  true,
		},
		{
			name:     "tag with @",
			named:    baseRef,
			tag:      "v1@0",
			expected: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := WithTag(tt.named, tt.tag)
			if (err != nil) != tt.wantErr {
				t.Errorf("WithTag() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && result.String() != tt.expected {
				t.Errorf("WithTag() = %v, want %v", result.String(), tt.expected)
			}
		})
	}
}

func TestTagNameOnly(t *testing.T) {
	tests := []struct {
		name     string
		ref      Named
		expected string
	}{
		{
			name: "name only - should add default tag",
			ref: namedReference{
				domain: "docker.io",
				path:   "library/llama",
			},
			expected: "docker.io/library/llama:latest",
		},
		{
			name: "already has tag - should not change",
			ref: namedReference{
				domain: "docker.io",
				path:   "library/llama",
				tag:    "v1.0",
			},
			expected: "docker.io/library/llama:v1.0",
		},
		{
			name: "has digest - should not change",
			ref: namedReference{
				domain: "docker.io",
				path:   "library/llama",
				digest: digest.Digest("sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"),
			},
			expected: "docker.io/library/llama@sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := TagNameOnly(tt.ref)
			if result.String() != tt.expected {
				t.Errorf("TagNameOnly() = %v, want %v", result.String(), tt.expected)
			}
		})
	}
}

func TestTagNameOnlyWithConfig(t *testing.T) {
	customConfig := Configuration{
		DefaultDomain:      "models.ai",
		OfficialRepoPrefix: "official/",
		DefaultTag:         "v1.0",
	}

	tests := []struct {
		name     string
		ref      Named
		config   Configuration
		expected string
	}{
		{
			name: "name only with custom config",
			ref: namedReference{
				domain: "models.ai",
				path:   "official/llama",
			},
			config:   customConfig,
			expected: "models.ai/official/llama:v1.0",
		},
		{
			name: "already has tag with custom config",
			ref: namedReference{
				domain: "models.ai",
				path:   "official/llama",
				tag:    "v2.0",
			},
			config:   customConfig,
			expected: "models.ai/official/llama:v2.0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := TagNameOnlyWithConfig(tt.ref, tt.config)
			if result.String() != tt.expected {
				t.Errorf("TagNameOnlyWithConfig() = %v, want %v", result.String(), tt.expected)
			}
		})
	}
}

func TestConfigurationOverrides(t *testing.T) {
	// Test the global configuration override functions
	originalDefault := DefaultConfiguration.DefaultDomain
	originalPrefix := DefaultConfiguration.OfficialRepoPrefix
	originalTag := DefaultConfiguration.DefaultTag

	// Override defaults
	SetDefaultDomain("custom.ai")
	SetOfficialRepoPrefix("models/")
	SetDefaultTag("v2.0")

	// Test that the overrides work
	if DefaultConfiguration.DefaultDomain != "custom.ai" {
		t.Errorf("SetDefaultDomain() did not update DefaultConfiguration.DefaultDomain")
	}
	if DefaultConfiguration.OfficialRepoPrefix != "models/" {
		t.Errorf("SetOfficialRepoPrefix() did not update DefaultConfiguration.OfficialRepoPrefix")
	}
	if DefaultConfiguration.DefaultTag != "v2.0" {
		t.Errorf("SetDefaultTag() did not update DefaultConfiguration.DefaultTag")
	}

	// Test parsing with the new defaults
	result, err := ParseNormalizedNamed("llama")
	if err != nil {
		t.Fatalf("ParseNormalizedNamed() with overridden config failed: %v", err)
	}
	expected := "custom.ai/models/llama:v2.0"
	if result.String() != expected {
		t.Errorf("ParseNormalizedNamed() with overridden config = %v, want %v", result.String(), expected)
	}

	// Restore original values
	DefaultConfiguration.DefaultDomain = originalDefault
	DefaultConfiguration.OfficialRepoPrefix = originalPrefix
	DefaultConfiguration.DefaultTag = originalTag
}

func TestNamedReferenceImplementation(t *testing.T) {
	ref := namedReference{
		domain: "docker.io",
		path:   "library/llama",
		tag:    "7b",
		digest: digest.Digest("sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"),
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
	if ref.Tag() != "7b" {
		t.Errorf("Tag() = %v, want %v", ref.Tag(), "7b")
	}
	if ref.Digest() != digest.Digest("sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855") {
		t.Errorf("Digest() = %v, want %v", ref.Digest(), digest.Digest("sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"))
	}

	expected := "docker.io/library/llama:7b@sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
	if ref.String() != expected {
		t.Errorf("String() = %v, want %v", ref.String(), expected)
	}
}