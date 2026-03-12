# API Graduation Roadmap

This document defines the upgrade path for `dns.astradns.com` APIs.

## Current state

- Active version: `v1alpha1`
- Scope: MVP APIs for upstream pools, cache profiles, and external policy references

## Graduation plan

### Phase 1 - Harden `v1alpha1`

- Keep additive-only changes (new optional fields, new status fields, new printer columns)
- Avoid breaking renames/removals
- Continue validating schema compatibility in CI

### Phase 2 - Introduce `v1beta1`

- Create `api/v1beta1` types with conversion-safe field naming
- Add hub/spoke conversion implementation
- Register conversion webhook in operator deployment profile
- Publish migration notes from `v1alpha1` to `v1beta1`

### Phase 3 - Stabilize `v1`

- Freeze API semantics
- Remove alpha-only behavior flags
- Keep conversion support from `v1alpha1`/`v1beta1` for at least one minor cycle

## Compatibility policy

- `v1alpha1`: best-effort compatibility, no intentional breaking changes without migration path
- `v1beta1`: compatibility required for one minor release
- `v1`: backward compatibility required for supported fields

## Required operator work

- Enable conversion webhook serving
- Add end-to-end upgrade tests across versions
- Ensure Helm and Kustomize install paths include conversion resources
