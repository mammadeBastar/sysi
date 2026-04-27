## Why

The current workflow describes `openspec-apply` and Superpowers as required during build apply, but the generated `sysi-apply` skill only says to "invoke or follow" the apply workflow. Also, `sysi init` currently scaffolds sysi state and `/system` only; it does not prepare the intended frontend/backend OpenSpec workspaces.

This change makes the orchestration contract stricter: apply must actually enter OpenSpec apply and Superpowers discipline, and initialization must create OpenSpec contexts only where implementation happens.

## What Changes

- Strengthen the Codex `sysi-apply` skill so it explicitly invokes the project OpenSpec apply workflow for the named change before implementation.
- Require `sysi-apply` to use Superpowers methods during implementation planning, TDD, debugging, and verification when those skills are available.
- Update `sysi change apply <name>` from passive guidance toward a concrete OpenSpec apply handoff.
- Make `sysi init` create `frontend/` and `backend/` directories when missing and run `openspec init` inside both directories.
- Keep the monorepo root and `/system` free of new OpenSpec initialization during `sysi init`.
- Route build-phase `sysi change` commands to the inferred `frontend/` or `backend/` OpenSpec workspace instead of the monorepo root.
- Add tests that prove initialization targets only frontend/backend and that apply behavior routes through OpenSpec apply semantics.

## Capabilities

### New Capabilities

- None.

### Modified Capabilities

- `project-lifecycle`: `sysi init` must initialize frontend/backend OpenSpec workspaces without initializing root or `/system`.
- `build-workflow`: `sysi change apply` must invoke or dispatch the OpenSpec apply workflow rather than only printing loose guidance.
- `agent-integration`: generated Codex `sysi-apply` instructions must require `openspec-apply` and Superpowers discipline with stronger wording and stop conditions.

## Impact

- Affected code: `internal/sysiapp/app.go`, `internal/sysiapp/app_test.go`.
- Affected templates: `internal/sysiapp/templates/agents/codex/sysi-apply/SKILL.md`.
- Affected docs: `README.md`.
- External dependency surface: existing `openspec` CLI becomes required during `sysi init` unless tests inject a fake executable.
- Existing sysi projects: running `sysi init` in an already initialized repository remains idempotent and must not reinitialize unrelated locations.
