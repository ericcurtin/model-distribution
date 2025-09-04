package distribution

import "github.com/docker/model-distribution/internal/reference"

// Reference represents a reference that can be either tagged or digested
type Reference interface {
	String() string
}

// Named represents a named reference
type Named interface {
	Reference
	Name() string
	Domain() string
	Path() string
}

// ReferenceConfig holds the normalization configuration for model references
type ReferenceConfig struct {
	// DefaultDomain is the default domain used for model references.
	DefaultDomain string

	// OfficialRepoPrefix is the namespace used for official models.
	OfficialRepoPrefix string

	// DefaultTag is the default tag if no tag is provided.
	DefaultTag string
}

// DefaultReferenceConfig provides the default normalization configuration
// for model distribution.
var DefaultReferenceConfig = ReferenceConfig{
	DefaultDomain:      "docker.io",
	OfficialRepoPrefix: "library/",
	DefaultTag:         "latest",
}

// SetDefaultDomain allows overriding the default domain at the model-distribution level.
// This affects all subsequent calls to ParseReference that don't specify a custom config.
func SetDefaultDomain(domain string) {
	reference.SetDefaultDomain(domain)
	DefaultReferenceConfig.DefaultDomain = domain
}

// SetOfficialRepoPrefix allows overriding the official repository prefix at the model-distribution level.
// This affects all subsequent calls to ParseReference that don't specify a custom config.
func SetOfficialRepoPrefix(prefix string) {
	reference.SetOfficialRepoPrefix(prefix)
	DefaultReferenceConfig.OfficialRepoPrefix = prefix
}

// SetDefaultTag allows overriding the default tag at the model-distribution level.
// This affects all subsequent calls to ParseReference that don't specify a custom config.
func SetDefaultTag(tag string) {
	reference.SetDefaultTag(tag)
	DefaultReferenceConfig.DefaultTag = tag
}

// ParseReference parses a string into a named reference, transforming a familiar name
// to a fully qualified reference using the current default configuration.
//
// Examples:
//   - "llama" -> "docker.io/library/llama:latest"
//   - "llama:7b" -> "docker.io/library/llama:7b"
//   - "myregistry.com/models/llama:7b" -> "myregistry.com/models/llama:7b"
//
// This function uses the defaults that can be overridden via SetDefaultDomain,
// SetOfficialRepoPrefix, and SetDefaultTag.
func ParseReference(s string) (Named, error) {
	ref, err := reference.ParseNormalizedNamed(s)
	if err != nil {
		return nil, err
	}
	return &namedWrapper{ref}, nil
}

// ParseReferenceWithConfig parses a string into a named reference using the provided
// configuration for normalization rules.
//
// This allows for per-call configuration without affecting the global defaults.
func ParseReferenceWithConfig(s string, config ReferenceConfig) (Named, error) {
	internalConfig := reference.Configuration{
		DefaultDomain:      config.DefaultDomain,
		OfficialRepoPrefix: config.OfficialRepoPrefix,
		DefaultTag:         config.DefaultTag,
	}
	
	ref, err := reference.ParseNormalizedNamedWithConfig(s, internalConfig)
	if err != nil {
		return nil, err
	}
	return &namedWrapper{ref}, nil
}

// FamiliarName returns a shortened version of the name familiar to users.
// Familiar names have the default domain and official repository prefix
// removed when appropriate.
//
// Examples:
//   - "docker.io/library/llama" -> "llama"
//   - "docker.io/someorg/model" -> "someorg/model"
//   - "myregistry.com/models/llama" -> "myregistry.com/models/llama"
func FamiliarName(named Named) string {
	if wrapper, ok := named.(*namedWrapper); ok {
		return reference.FamiliarName(wrapper.internal)
	}
	// Fallback for other implementations
	return named.String()
}

// FamiliarNameWithConfig returns a shortened version of the name using the provided
// configuration.
func FamiliarNameWithConfig(named Named, config ReferenceConfig) string {
	if wrapper, ok := named.(*namedWrapper); ok {
		internalConfig := reference.Configuration{
			DefaultDomain:      config.DefaultDomain,
			OfficialRepoPrefix: config.OfficialRepoPrefix,
			DefaultTag:         config.DefaultTag,
		}
		return reference.FamiliarNameWithConfig(wrapper.internal, internalConfig)
	}
	// Fallback for other implementations
	return named.String()
}

// namedWrapper wraps the internal reference.Named to expose it publicly
type namedWrapper struct {
	internal reference.Named
}

func (w *namedWrapper) String() string {
	return w.internal.String()
}

func (w *namedWrapper) Name() string {
	return w.internal.Name()
}

func (w *namedWrapper) Domain() string {
	return w.internal.Domain()
}

func (w *namedWrapper) Path() string {
	return w.internal.Path()
}