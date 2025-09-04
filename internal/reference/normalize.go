package reference

import (
	"fmt"
	"strings"

	"github.com/opencontainers/go-digest"
)

// Configuration holds the normalization configuration that can be overridden
// at the model-distribution level to avoid requiring each client to implement this.
type Configuration struct {
	// DefaultDomain is the default domain used for model references.
	// It is used to normalize "familiar" names to canonical names.
	DefaultDomain string

	// OfficialRepoPrefix is the namespace used for official models.
	// It is used to normalize "familiar" names to canonical names, for example,
	// to convert "llama" to "docker.io/library/llama:latest".
	OfficialRepoPrefix string

	// DefaultTag is the default tag if no tag is provided.
	DefaultTag string
}

// DefaultConfiguration provides the default normalization configuration
// for model distribution. These can be overridden as needed.
var DefaultConfiguration = Configuration{
	DefaultDomain:      "docker.io",
	OfficialRepoPrefix: "library/",
	DefaultTag:         "latest",
}

// SetDefaultDomain allows overriding the default domain at the model-distribution level
func SetDefaultDomain(domain string) {
	DefaultConfiguration.DefaultDomain = domain
}

// SetOfficialRepoPrefix allows overriding the official repository prefix at the model-distribution level
func SetOfficialRepoPrefix(prefix string) {
	DefaultConfiguration.OfficialRepoPrefix = prefix
}

// SetDefaultTag allows overriding the default tag at the model-distribution level
func SetDefaultTag(tag string) {
	DefaultConfiguration.DefaultTag = tag
}

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

// Tagged represents a tagged reference
type Tagged interface {
	Tag() string
}

// Digested represents a digested reference
type Digested interface {
	Digest() digest.Digest
}

// NamedTagged represents a reference that has both a name and a tag
type NamedTagged interface {
	Named
	Tagged
}

// namedReference implements Named interface
type namedReference struct {
	domain string
	path   string
	tag    string
	digest digest.Digest
}

func (r namedReference) String() string {
	result := r.domain + "/" + r.path
	if r.tag != "" {
		result += ":" + r.tag
	}
	if r.digest != "" {
		result += "@" + string(r.digest)
	}
	return result
}

func (r namedReference) Name() string {
	return r.domain + "/" + r.path
}

func (r namedReference) Domain() string {
	return r.domain
}

func (r namedReference) Path() string {
	return r.path
}

func (r namedReference) Tag() string {
	return r.tag
}

func (r namedReference) Digest() digest.Digest {
	return r.digest
}

// ParseNormalizedNamed parses a string into a named reference
// transforming a familiar name to a fully qualified reference.
// It uses the current DefaultConfiguration for normalization rules.
func ParseNormalizedNamed(s string) (Named, error) {
	return ParseNormalizedNamedWithConfig(s, DefaultConfiguration)
}

// ParseNormalizedNamedWithConfig parses a string into a named reference
// using the provided configuration for normalization rules.
func ParseNormalizedNamedWithConfig(s string, config Configuration) (Named, error) {
	if s == "" {
		return nil, fmt.Errorf("empty reference")
	}

	domain, remainder := splitDomainWithConfig(s, config)
	
	var remote string
	var tag string
	var digestStr string

	// Check for digest
	if digestSep := strings.IndexRune(remainder, '@'); digestSep > -1 {
		remote = remainder[:digestSep]
		digestStr = remainder[digestSep+1:]
		
		// Validate digest
		if _, err := digest.Parse(digestStr); err != nil {
			return nil, fmt.Errorf("invalid digest: %v", err)
		}
	} else {
		// Check for tag
		if tagSep := strings.IndexRune(remainder, ':'); tagSep > -1 {
			remote = remainder[:tagSep]
			tag = remainder[tagSep+1:]
		} else {
			remote = remainder
			tag = config.DefaultTag
		}
	}

	if remote == "" {
		return nil, fmt.Errorf("invalid reference format: missing repository name")
	}

	// Ensure remote name is lowercase
	if strings.ToLower(remote) != remote {
		return nil, fmt.Errorf("invalid reference format: repository name (%s) must be lowercase", remote)
	}

	var parsedDigest digest.Digest
	if digestStr != "" {
		var err error
		parsedDigest, err = digest.Parse(digestStr)
		if err != nil {
			return nil, fmt.Errorf("invalid digest: %v", err)
		}
	}

	return namedReference{
		domain: domain,
		path:   remote,
		tag:    tag,
		digest: parsedDigest,
	}, nil
}

// splitDomainWithConfig splits a repository name to domain and remote-name using the provided configuration.
// If no valid domain is found, the default domain from config is used.
func splitDomainWithConfig(name string, config Configuration) (domain, remoteName string) {
	maybeDomain, maybeRemoteName, ok := strings.Cut(name, "/")
	if !ok {
		// Fast-path for single element ("familiar" names), such as "llama"
		// or "llama:latest". Familiar names must be handled separately.
		//
		// Canonicalize them as "defaultDomain/officialRepoPrefix/name[:tag]"
		return config.DefaultDomain, config.OfficialRepoPrefix + name
	}

	switch {
	case maybeDomain == "localhost":
		// localhost is a reserved namespace and always considered a domain.
		domain, remoteName = maybeDomain, maybeRemoteName
	case strings.ContainsAny(maybeDomain, ".:"):
		// Likely a domain or IP-address:
		//
		// - contains a "." (e.g., "example.com" or "127.0.0.1")
		// - contains a ":" (e.g., "example:5000", "::1", or "[::1]:5000")
		domain, remoteName = maybeDomain, maybeRemoteName
	case strings.ToLower(maybeDomain) != maybeDomain:
		// Uppercase namespaces are not allowed, so if the first element
		// is not lowercase, we assume it to be a domain-name.
		domain, remoteName = maybeDomain, maybeRemoteName
	default:
		// None of the above: it's not a domain, so use the default, and
		// use the original name as the remote-name.
		domain, remoteName = config.DefaultDomain, name
	}

	if domain == config.DefaultDomain && !strings.ContainsRune(remoteName, '/') {
		// Canonicalize "familiar" names, but only on the default domain:
		//
		// "defaultDomain/model[:tag]" => "defaultDomain/officialRepoPrefix/model[:tag]"
		remoteName = config.OfficialRepoPrefix + remoteName
	}

	return domain, remoteName
}

// FamiliarName returns a shortened version of the name familiar
// to users. Familiar names have the default domain and official
// repository prefix removed when appropriate.
// For example, "docker.io/library/llama" will have the familiar
// name "llama" and "docker.io/someorg/model" will be "someorg/model".
func FamiliarName(named Named) string {
	return FamiliarNameWithConfig(named, DefaultConfiguration)
}

// FamiliarNameWithConfig returns a shortened version of the name using the provided configuration.
func FamiliarNameWithConfig(named Named, config Configuration) string {
	domain := named.Domain()
	path := named.Path()

	if domain == config.DefaultDomain {
		// Handle official repositories which have the pattern "officialRepoPrefix/<official repo name>"
		if strings.HasPrefix(path, config.OfficialRepoPrefix) {
			if remainder := strings.TrimPrefix(path, config.OfficialRepoPrefix); !strings.ContainsRune(remainder, '/') {
				return remainder
			}
		}
		// Return just the path for default domain
		return path
	}

	// For non-default domains, return domain/path
	return domain + "/" + path
}

// IsNameOnly returns true if the reference only contains a repository name
// without tag or digest.
func IsNameOnly(ref Named) bool {
	if tagged, ok := ref.(Tagged); ok && tagged.Tag() != "" {
		return false
	}
	if digested, ok := ref.(Digested); ok && digested.Digest() != "" {
		return false
	}
	return true
}

// WithTag creates a new NamedTagged reference with the given tag
func WithTag(named Named, tag string) (NamedTagged, error) {
	if tag == "" {
		return nil, fmt.Errorf("tag cannot be empty")
	}
	
	// Validate tag format (basic validation)
	if strings.ContainsAny(tag, " \t\n\r@") {
		return nil, fmt.Errorf("invalid tag format: %s", tag)
	}

	return namedReference{
		domain: named.Domain(),
		path:   named.Path(),
		tag:    tag,
	}, nil
}

// TagNameOnly adds the default tag to a reference if it only has
// a repo name.
func TagNameOnly(ref Named) Named {
	return TagNameOnlyWithConfig(ref, DefaultConfiguration)
}

// TagNameOnlyWithConfig adds the default tag from config to a reference if it only has a repo name.
func TagNameOnlyWithConfig(ref Named, config Configuration) Named {
	if IsNameOnly(ref) {
		namedTagged, err := WithTag(ref, config.DefaultTag)
		if err != nil {
			// Default tag must be valid, to create a NamedTagged
			// type with non-validated input the WithTag function
			// should be used instead
			panic(err)
		}
		return namedTagged
	}
	return ref
}