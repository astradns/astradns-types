# Contributing to astradns-types

Thanks for contributing to the shared AstraDNS contract module.

## Pull Request Checklist

- Treat all exported API changes as cross-repository changes.
- Run `make test` and `make vet` locally.
- Regenerate generated code with `make generate` when type definitions change.
- Use conventional commits (`feat:`, `fix:`, `docs:`, `test:`, `refactor:`, `chore:`).
- Avoid introducing new dependencies unless there is no viable standard library alternative.

## AI/OpenCode Contributions

AI-assisted changes are welcome, but must follow repository guardrails in `AGENTS.md`.

Minimum requirements for AI-generated changes:

- No secrets, credentials, or personal data.
- Preserve module boundaries (`astradns-types` must not import `astradns-agent` or `astradns-operator`).
- Prefer additive changes to APIs to avoid breaking consumers.
- Do not manually edit generated files (`zz_generated.deepcopy.go`).
