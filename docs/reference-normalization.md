# Model Distribution Reference Normalization

This document describes the reference normalization functionality added to model-distribution, which allows overriding default domain and repository prefix settings at the model-distribution level.

## Overview

The normalization functionality provides a consistent way to handle model references across the model-distribution system. It allows "familiar" names to be automatically expanded to fully qualified references, and provides the ability to configure default values that apply system-wide.

## Key Features

- **Configurable defaults**: Override `defaultDomain`, `officialRepoPrefix`, and `defaultTag` at the model-distribution level
- **Per-call configuration**: Use custom configuration for specific operations without affecting global defaults
- **Familiar name conversion**: Convert between short "familiar" names and fully qualified references
- **Server-side implementation**: Implemented in the core/server-side to avoid requiring each client to implement this logic

## Usage Examples

### Basic Usage with Default Configuration

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/docker/model-distribution/distribution"
)

func main() {
    // Parse familiar names using default configuration
    ref, err := distribution.ParseReference("llama")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println(ref.String())         // Output: docker.io/library/llama:latest
    fmt.Println(distribution.FamiliarName(ref)) // Output: llama
}
```

### Global Configuration Override

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/docker/model-distribution/distribution"
)

func main() {
    // Override global defaults
    distribution.SetDefaultDomain("models.ai")
    distribution.SetOfficialRepoPrefix("official/")
    distribution.SetDefaultTag("v1.0")
    
    // Now all parsing uses the new defaults
    ref, err := distribution.ParseReference("llama")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println(ref.String())         // Output: models.ai/official/llama:v1.0
    fmt.Println(distribution.FamiliarName(ref)) // Output: llama
}
```

### Per-Call Configuration

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/docker/model-distribution/distribution"
)

func main() {
    // Use custom configuration for specific operations
    config := distribution.ReferenceConfig{
        DefaultDomain:      "hub.models",
        OfficialRepoPrefix: "public/",
        DefaultTag:         "stable",
    }
    
    ref, err := distribution.ParseReferenceWithConfig("llama", config)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println(ref.String())         // Output: hub.models/public/llama:stable
    fmt.Println(distribution.FamiliarNameWithConfig(ref, config)) // Output: llama
}
```

## Reference Normalization Rules

The normalization follows these rules:

1. **Familiar names** (single component without `/`):
   - `"llama"` → `"<defaultDomain>/<officialRepoPrefix>llama:<defaultTag>"`
   - `"llama:7b"` → `"<defaultDomain>/<officialRepoPrefix>llama:7b"`

2. **Domain detection**:
   - If first component contains `.` or `:`, it's treated as a domain
   - If first component is uppercase, it's treated as a domain
   - `localhost` is always treated as a domain

3. **Default domain handling**:
   - For the default domain, single-component paths get the official prefix
   - `"docker.io/llama"` → `"docker.io/library/llama:latest"`

4. **Familiar name generation**:
   - Official repositories on the default domain: `"docker.io/library/llama"` → `"llama"`
   - Other repositories on default domain: `"docker.io/myorg/model"` → `"myorg/model"`
   - Non-default domains: `"myregistry.com/models/llama"` → `"myregistry.com/models/llama"`

## Default Configuration

The default configuration matches Docker Hub conventions:

- **DefaultDomain**: `"docker.io"`
- **OfficialRepoPrefix**: `"library/"`  
- **DefaultTag**: `"latest"`

These can be overridden using the `SetDefaultDomain()`, `SetOfficialRepoPrefix()`, and `SetDefaultTag()` functions.

## API Reference

### Types

- `Reference`: Basic reference interface
- `Named`: Named reference with domain and path methods
- `ReferenceConfig`: Configuration structure for normalization rules

### Functions

- `ParseReference(s string) (Named, error)`: Parse using global defaults
- `ParseReferenceWithConfig(s string, config ReferenceConfig) (Named, error)`: Parse with custom config
- `FamiliarName(named Named) string`: Get familiar name using global defaults
- `FamiliarNameWithConfig(named Named, config ReferenceConfig) string`: Get familiar name with custom config
- `SetDefaultDomain(domain string)`: Override global default domain
- `SetOfficialRepoPrefix(prefix string)`: Override global official repo prefix  
- `SetDefaultTag(tag string)`: Override global default tag

This implementation provides the foundation for consistent reference handling across the model-distribution system while allowing flexibility for different deployment scenarios.