# AI Agent Guidelines -- astradns-types

This document governs AI-assisted contributions to the `astradns-types` repository. AI agents (Claude, Copilot, and others) must follow these guidelines alongside standard project conventions.

## Principles

AI contributions are held to the same quality standards as human contributions. There are no exceptions. Every change must be reviewable, testable, and justified.

## Rules

1. **All AI-generated code must pass existing tests and linters.** Run `make test` and `make vet` before proposing any change.
2. **Do not introduce new dependencies without explicit approval.** This module is imported by both the agent and the operator. Every dependency added here becomes a transitive dependency for the entire project.
3. **Do not modify API contracts without discussion.** The CRD types in `api/v1alpha1/` and the `engine.Engine` interface are consumed by multiple repositories. Any breaking change requires coordinated updates across `astradns-agent` and `astradns-operator`.
4. **Do not commit secrets, credentials, or PII.** No tokens, passwords, API keys, or personal data in code, comments, or test fixtures.
5. **Follow conventional commit format.** Use prefixes such as `feat:`, `fix:`, `refactor:`, `docs:`, `test:`, `chore:`.
6. **Respect import boundaries.** This module must not import from `astradns-agent` or `astradns-operator`. It is the shared foundation that flows in one direction only.

## Repo-Specific Context

This is the **shared contract module** for the AstraDNS project. It contains:

- **CRD types** (`api/v1alpha1/`) -- the Kubernetes custom resource definitions shared across the control plane and data plane.
- **Engine interface** (`engine/interface.go`) -- the abstraction that all DNS engine implementations must satisfy.
- **Engine registry** (`engine/registry.go`) -- the registration mechanism for engine implementations.
- **Config templates** (`engine/templates.go`) -- template utilities for engine configuration rendering.
- **Engine config types** (`engineconfig/`) -- additional configuration structures.

Changes here affect both `astradns-agent` and `astradns-operator`. Breaking changes (renamed fields, removed methods, altered semantics) require coordinated pull requests across all consuming repositories. When in doubt, prefer additive changes.

## Code Style

- **Language:** Go
- **Follow existing patterns** in the codebase. Do not introduce new structural conventions without discussion.
- **Structured logging:** Use `log/slog` from the standard library. Do not add external logging frameworks.
- **Naming:** Use clear, domain-aligned names. CRD types follow Kubernetes API conventions. Interface methods should be self-documenting.
- **Comments:** Exported types and functions require GoDoc comments. Follow the `// TypeName does X` convention.

## Testing Expectations

- All exported types must have corresponding tests validating serialization and deserialization where applicable.
- The `engine.Engine` interface changes must include updated or new tests in this module.
- DeepCopy methods are auto-generated (`zz_generated.deepcopy.go`). Do not edit generated files. Run `make generate` after modifying `*_types.go` files.
- Run `make test` to execute the full test suite. Run `make vet` for static analysis.

## Generated Files

Do not manually edit files matching these patterns:

- `zz_generated.deepcopy.go`
- Any file with a `// Code generated` header

Regenerate them with `make generate` after modifying source types.
