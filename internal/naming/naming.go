package naming

import (
	"strings"

	"github.com/google/go-containerregistry/pkg/name"
)

// DefaultNamespace holds the default registry namespace configuration
type DefaultNamespace struct {
	Registry string
}

// ParseReference parses a reference string, applying the default namespace if needed
func (dn *DefaultNamespace) ParseReference(reference string) (name.Reference, error) {
	// If no default namespace is configured, use standard parsing
	if dn == nil || dn.Registry == "" {
		return name.ParseReference(reference)
	}

	// If the reference already contains a registry (has a domain with dot or port), use as-is
	if hasExplicitRegistry(reference) {
		return name.ParseReference(reference)
	}

	// Apply default registry to the reference
	qualified := dn.Registry + "/" + reference
	return name.ParseReference(qualified)
}

// ParseTag parses a tag string, applying the default namespace if needed
func (dn *DefaultNamespace) ParseTag(tag string) (name.Tag, error) {
	// If no default namespace is configured, use standard parsing
	if dn == nil || dn.Registry == "" {
		return name.NewTag(tag)
	}

	// If the tag already contains a registry (has a domain with dot or port), use as-is
	if hasExplicitRegistry(tag) {
		return name.NewTag(tag)
	}

	// Apply default registry to the tag
	qualified := dn.Registry + "/" + tag
	return name.NewTag(qualified)
}

// hasExplicitRegistry checks if a reference already contains an explicit registry
// This is a simple heuristic: if it contains a dot before the first slash or
// a colon followed by a port number before the first slash, it's likely a registry hostname
func hasExplicitRegistry(reference string) bool {
	// Find the first slash
	slashIndex := strings.Index(reference, "/")
	
	// If no slash, check if it looks like a registry (contains dot, not just tag colon)
	if slashIndex == -1 {
		// If it contains a dot, it's likely a registry
		if strings.Contains(reference, ".") {
			return true
		}
		// If it contains a colon, check if it's followed by a numeric port
		colonIndex := strings.Index(reference, ":")
		if colonIndex != -1 {
			// Check if what comes after colon looks like a port number
			afterColon := reference[colonIndex+1:]
			// If it's all digits, it's a port; otherwise it's a tag
			for _, r := range afterColon {
				if r < '0' || r > '9' {
					return false // It's a tag, not a port
				}
			}
			return len(afterColon) > 0 // It's a port if non-empty and all digits
		}
		return false
	}
	
	// Check the part before the first slash
	beforeSlash := reference[:slashIndex]
	
	// If it contains a dot (domain), it's a registry
	if strings.Contains(beforeSlash, ".") {
		return true
	}
	
	// If it contains a colon, check if it's followed by a numeric port
	if strings.Contains(beforeSlash, ":") {
		colonIndex := strings.Index(beforeSlash, ":")
		afterColon := beforeSlash[colonIndex+1:]
		// Check if what comes after colon looks like a port number
		for _, r := range afterColon {
			if r < '0' || r > '9' {
				return false // Not a port
			}
		}
		return len(afterColon) > 0 // It's a port if non-empty and all digits
	}
	
	return false
}

// Global default namespace instance
var globalDefaultNamespace *DefaultNamespace

// SetDefaultNamespace sets the global default namespace
func SetDefaultNamespace(registry string) {
	globalDefaultNamespace = &DefaultNamespace{Registry: registry}
}

// GetDefaultNamespace returns the current global default namespace
func GetDefaultNamespace() *DefaultNamespace {
	return globalDefaultNamespace
}

// ParseReference is a convenience function that uses the global default namespace
func ParseReference(reference string) (name.Reference, error) {
	return globalDefaultNamespace.ParseReference(reference)
}

// ParseTag is a convenience function that uses the global default namespace
func ParseTag(tag string) (name.Tag, error) {
	return globalDefaultNamespace.ParseTag(tag)
}