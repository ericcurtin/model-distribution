package files

import (
	"strings"

	"github.com/google/go-containerregistry/pkg/name"
)

// officialRepoPrefix is the default namespace prefix for models that don't specify a registry
const officialRepoPrefix = "ai/"

// NormalizeReference adds the default officialRepoPrefix to references that don't have a registry/namespace specified.
// For example: "smollm2:135M-Q4_K_M" becomes "ai/smollm2:135M-Q4_K_M"
// References that already have a registry or namespace are left unchanged.
func NormalizeReference(reference string) string {
	if reference == "" {
		return reference
	}

	// Check if the reference already contains the official prefix by doing string matching first
	// This handles cases where we want to avoid double prefixing
	officialPrefixWithoutSlash := strings.TrimSuffix(officialRepoPrefix, "/")
	if strings.HasPrefix(reference, officialPrefixWithoutSlash+"/") {
		return reference
	}

	// Check if the reference looks like it has a custom registry (contains a dot or colon before any slash)
	// This is a fast check before parsing
	firstSlashIdx := strings.Index(reference, "/")
	if firstSlashIdx > 0 {
		beforeSlash := reference[:firstSlashIdx]
		if strings.Contains(beforeSlash, ".") || strings.Contains(beforeSlash, ":") {
			// This looks like a custom registry
			return reference
		}
		// This looks like it already has a namespace
		return reference
	}

	// For simple references like "model:tag" or "model", we need to add the prefix
	// Try to parse the reference to understand its structure
	ref, err := name.ParseReference(reference, name.WeakValidation)
	if err != nil {
		// If we can't parse it, try simple string manipulation as fallback
		return addPrefixFallback(reference)
	}

	// Check if this is a non-default registry
	if ref.Context().Registry.Name() != name.DefaultRegistry {
		return reference
	}

	// Get the repository part of the reference
	repository := ref.Context().RepositoryStr()

	// Docker Hub adds "library/" prefix automatically to simple names
	// We want to replace "library/" with our official prefix "ai/"
	if strings.HasPrefix(repository, "library/") {
		// Extract the original name without library prefix
		originalName := repository[8:] // len("library/") = 8

		// Construct the normalized reference with our prefix in short form
		normalizedRepo := officialRepoPrefix + originalName

		// Return in short form (without registry prefix for Docker Hub)
		if tagged, ok := ref.(name.Tag); ok {
			if tagged.TagStr() == "latest" && !strings.Contains(reference, ":") {
				// Original reference had no tag, keep it without :latest for simplicity
				return normalizedRepo
			}
			return normalizedRepo + ":" + tagged.TagStr()
		} else if digested, ok := ref.(name.Digest); ok {
			return normalizedRepo + "@" + digested.DigestStr()
		}

		// Fallback for references without tag/digest (should get :latest)
		return normalizedRepo
	}

	// For other cases, return as-is
	return reference
}

// addPrefixFallback is a simple fallback for cases where parsing fails
func addPrefixFallback(reference string) string {
	// Handle digest references
	if strings.Contains(reference, "@sha256:") {
		parts := strings.Split(reference, "@")
		if len(parts) == 2 {
			return officialRepoPrefix + parts[0] + "@" + parts[1]
		}
	}

	// Handle tag references
	if strings.Contains(reference, ":") {
		parts := strings.Split(reference, ":")
		if len(parts) >= 2 {
			return officialRepoPrefix + parts[0] + ":" + strings.Join(parts[1:], ":")
		}
	}

	// Simple name without tag/digest
	return officialRepoPrefix + reference
}
