# Naming Package

The naming package provides configurable default namespace support for container registry references in the model-distribution system.

## Problem

By default, the `github.com/google/go-containerregistry/pkg/name` package assumes Docker Hub (`index.docker.io`) as the default registry for image references that don't include an explicit registry hostname. For model distribution systems, it's often desirable to use a different default registry.

## Solution

This package provides wrapper functions around `name.ParseReference` and `name.ParseTag` that apply a configurable default namespace when no explicit registry is provided in the reference.

## Usage

### Client-specific Configuration

```go
import (
    "github.com/docker/model-distribution/registry"
    "github.com/docker/model-distribution/distribution"
)

// Configure registry client with custom default namespace
registryClient := registry.NewClient(
    registry.WithDefaultNamespace("registry.example.com"),
)

// Configure distribution client with custom default namespace
distClient, err := distribution.NewClient(
    distribution.WithStoreRootPath("/path/to/store"),
    distribution.WithDefaultRegistry("registry.example.com"),
)
```

### Global Configuration

```go
import "github.com/docker/model-distribution/internal/naming"

// Set global default namespace
naming.SetDefaultNamespace("registry.example.com")

// Use global convenience functions
ref, err := naming.ParseReference("mymodel:latest")
// Will resolve to registry.example.com/mymodel:latest

tag, err := naming.ParseTag("mymodel:latest") 
// Will resolve to registry.example.com/mymodel:latest
```

### Direct Usage

```go
import "github.com/docker/model-distribution/internal/naming"

// Create a namespace configuration
ns := &naming.DefaultNamespace{Registry: "registry.example.com"}

// Parse references with custom default
ref, err := ns.ParseReference("mymodel:latest")
// Will resolve to registry.example.com/mymodel:latest

// Explicit registries are preserved
ref, err := ns.ParseReference("other.registry.com/mymodel:latest")
// Will remain other.registry.com/mymodel:latest
```

## CLI Usage

The `mdltool` command-line interface supports a `--default-registry` flag:

```bash
# Use custom default registry
mdltool --default-registry registry.example.com pull mymodel:latest

# This will pull from registry.example.com/mymodel:latest instead of 
# the Docker Hub default index.docker.io/library/mymodel:latest
```

## Behavior

- **No explicit registry**: Default namespace is applied
  - Input: `mymodel:latest` → Output: `registry.example.com/mymodel:latest`
  - Input: `user/mymodel:latest` → Output: `registry.example.com/user/mymodel:latest`

- **Explicit registry**: Default namespace is ignored
  - Input: `other.registry.com/mymodel:latest` → Output: `other.registry.com/mymodel:latest`
  - Input: `localhost:5000/mymodel:latest` → Output: `localhost:5000/mymodel:latest`

- **No default configured**: Falls back to standard go-containerregistry behavior
  - Input: `mymodel:latest` → Output: `index.docker.io/library/mymodel:latest`

## Registry Detection

The package uses heuristics to detect whether a reference already contains an explicit registry:

- Contains a dot (`.`) before the first slash → Likely a registry hostname
- Contains a colon (`:`) followed by digits before the first slash → Likely a registry with port
- Contains a colon followed by non-digits → Likely a tag, not a registry port

Examples:
- `mymodel:latest` → No explicit registry (`:latest` is a tag)
- `registry.com/mymodel:latest` → Has explicit registry (contains dot)
- `localhost:5000/mymodel:latest` → Has explicit registry (port number)
- `user/mymodel:latest` → No explicit registry

## Integration Points

This functionality is integrated at the following levels:

1. **Registry Client**: `registry.Client` can be configured with a default namespace
2. **Distribution Client**: `distribution.Client` configures both registry client and store operations
3. **Store Operations**: Uses global namespace configuration for tag/reference parsing
4. **CLI Tool**: Accepts `--default-registry` flag to configure the default

## Backward Compatibility

This change is fully backward compatible:

- Existing code without default namespace configuration continues to work unchanged
- All explicit registry references continue to work as before
- The default behavior (Docker Hub) is preserved when no configuration is provided