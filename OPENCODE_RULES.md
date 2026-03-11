# OpenCode Rules -- astradns-types

These rules define how AI agents may contribute to this repository.

1. Run `make test` and `make vet` before opening a PR.
2. Treat exported API changes as cross-repo changes.
3. Do not add dependencies without explicit maintainer approval.
4. Regenerate code with `make generate` when API type files change.
5. Do not manually edit generated files (`zz_generated.deepcopy.go`).
6. Never commit secrets, credentials, tokens, or personal data.
7. Follow `AGENTS.md` for repository-specific constraints.
