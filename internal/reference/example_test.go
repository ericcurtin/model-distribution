package reference

import (
	"fmt"
	"testing"
)

// TestDemonstrateFunctionality demonstrates the normalize functionality
func TestDemonstrateFunctionality(t *testing.T) {
	fmt.Println("\n=== Model Distribution Normalization Demo ===")
	
	// Example 1: Using default configuration
	fmt.Println("\n1. Default Configuration (docker.io, library/, latest):")
	examples := []string{
		"llama",
		"llama:7b",
		"docker.io/library/llama:7b",
		"myregistry.com/models/llama:7b",
		"localhost/models/llama",
	}

	for _, example := range examples {
		ref, err := ParseNormalizedNamed(example)
		if err != nil {
			t.Logf("Error parsing %s: %v", example, err)
			continue
		}
		fmt.Printf("  Input: %-30s -> Normalized: %-50s -> Familiar: %s\n", 
			example, ref.String(), FamiliarName(ref))
	}

	// Example 2: Using custom configuration via global overrides
	fmt.Println("\n2. Custom Configuration via Global Overrides (models.ai, official/, v1.0):")
	
	// Save original values
	originalDomain := DefaultConfiguration.DefaultDomain
	originalPrefix := DefaultConfiguration.OfficialRepoPrefix
	originalTag := DefaultConfiguration.DefaultTag
	
	// Override the global defaults
	SetDefaultDomain("models.ai")
	SetOfficialRepoPrefix("official/")
	SetDefaultTag("v1.0")

	for _, example := range examples {
		ref, err := ParseNormalizedNamed(example)
		if err != nil {
			t.Logf("Error parsing %s: %v", example, err)
			continue
		}
		fmt.Printf("  Input: %-30s -> Normalized: %-50s -> Familiar: %s\n", 
			example, ref.String(), FamiliarName(ref))
	}

	// Example 3: Using configuration passed to function
	fmt.Println("\n3. Per-Call Configuration (hub.models, public/, stable):")
	customConfig := Configuration{
		DefaultDomain:      "hub.models",
		OfficialRepoPrefix: "public/",
		DefaultTag:         "stable",
	}

	for _, example := range examples {
		ref, err := ParseNormalizedNamedWithConfig(example, customConfig)
		if err != nil {
			t.Logf("Error parsing %s: %v", example, err)
			continue
		}
		familiar := FamiliarNameWithConfig(ref, customConfig)
		fmt.Printf("  Input: %-30s -> Normalized: %-50s -> Familiar: %s\n", 
			example, ref.String(), familiar)
	}

	// Restore original values
	DefaultConfiguration.DefaultDomain = originalDomain
	DefaultConfiguration.OfficialRepoPrefix = originalPrefix
	DefaultConfiguration.DefaultTag = originalTag
	
	fmt.Println("\n=== Demo Complete ===")
}