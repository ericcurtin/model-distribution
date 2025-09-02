# Default Official Repo Prefix Implementation

This document describes the implementation of the default "ai/" prefix functionality in the model-distribution backend.

## Overview

The model-distribution backend now automatically adds an "ai/" prefix to model references that don't specify a registry or namespace. This ensures that all clients using the model-distribution library will consistently use the official "ai/" namespace for simple model names.

## Implementation

### Files Added

1. **`distribution/reference/files/normalize.go`**
   - Contains the `NormalizeReference()` function that adds the "ai/" prefix
   - Handles both tag and digest references
   - Preserves existing namespaces and custom registries

2. **`distribution/reference/files/normalize_test.go`**
   - Comprehensive tests covering all normalization scenarios
   - Edge case testing for malformed references

3. **`distribution/normalize_integration_test.go`**
   - Integration tests ensuring the normalization works with the distribution client

### Files Modified

1. **`distribution/client.go`**
   - Import the normalization package
   - Apply normalization in `PullModel()` before calling registry
   - Apply normalization in `PushModel()` before calling registry

## Behavior Examples

| Input Reference | Normalized Reference | Notes |
|-----------------|---------------------|-------|
| `smollm2:135M-Q4_K_M` | `ai/smollm2:135M-Q4_K_M` | Simple model name gets ai/ prefix |
| `ai/smollm2:135M-Q4_K_M` | `ai/smollm2:135M-Q4_K_M` | Already prefixed, unchanged |
| `huggingface/model:tag` | `huggingface/model:tag` | Custom namespace, unchanged |
| `registry.com/model:tag` | `registry.com/model:tag` | Custom registry, unchanged |
| `model@sha256:abc123` | `ai/model@sha256:abc123` | Digest references supported |

## Benefits

1. **Consistency**: All clients automatically use the "ai/" namespace for official models
2. **Backwards Compatible**: Existing references with namespaces/registries are unchanged
3. **Transparent**: Works at the backend level, no client changes required
4. **Configurable**: The prefix is defined as a constant that can be easily changed

## Testing

The implementation includes:
- 20+ unit tests covering various reference formats
- Integration tests with the distribution client
- End-to-end CLI testing demonstrating the functionality
- All existing tests continue to pass

## Usage

No changes are required for clients. The normalization happens automatically in the backend:

```bash
# These commands now automatically use ai/ prefix:
model-distribution-tool pull smollm2:135M-Q4_K_M
model-distribution-tool push mymodel:latest

# Custom registries/namespaces work as before:
model-distribution-tool pull registry.example.com/models/llama:v1.0
model-distribution-tool pull huggingface/model:tag
```

The logs will show the normalized references being used internally while maintaining the expected user experience.