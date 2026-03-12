# astradns-types

Shared Go types for the AstraDNS platform. This module defines the CRD API types, engine abstraction interfaces, configuration templates, and rendering contracts used by both the operator (control plane) and the agent (data plane).

This module is not meant to be deployed independently. It is imported as a Go dependency by `astradns-operator` and `astradns-agent`.

## Contents

| Package | Description |
|---|---|
| `api/v1alpha1` | CRD types: `DNSUpstreamPool`, `DNSCacheProfile`, `ExternalDNSPolicy` |
| `engine` | `Engine` interface and engine registry for pluggable DNS backends |
| `engineconfig` | `ConfigRenderer` interface, renderer registry, and config template types |

## API Group

All CRD types belong to the API group `dns.astradns.com/v1alpha1`.

## API Lifecycle

The graduation and compatibility strategy is documented in `docs/api-graduation-roadmap.md`.

## Usage

Import the module in your `go.mod`:

```
require github.com/astradns/astradns-types v0.0.0
```

Import paths:

```go
import (
    dnsv1alpha1 "github.com/astradns/astradns-types/api/v1alpha1"
    "github.com/astradns/astradns-types/engine"
    "github.com/astradns/astradns-types/engineconfig"
)
```

## Development

```sh
# Generate deepcopy methods for CRD types
make generate

# Run unit tests
make test

# Run static analysis
make vet
```

## Release

Tagging with `vX.Y.Z` triggers the release workflow (`.github/workflows/release.yml`) to run build/test/vet checks before publishing the module version.

## Contribution Policy

- Human and AI contributions: `CONTRIBUTING.md`
- OpenCode-specific guardrails: `OPENCODE_RULES.md`
- Repository-level AI constraints: `AGENTS.md`
