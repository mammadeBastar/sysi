## Why

The current workflow describes `openspec-apply` and Superpowers as required during build apply, but the generated `sys-apply` skill only says to "invoke or follow" the apply workflow. Also, `sys init` currently scaffolds sys state and `/system` only; it does not prepare the intended frontend/backend OpenSpec workspaces.

This change makes the orchestration contract stricter: apply must actually enter OpenSpec apply and Superpowers discipline, and initialization must create OpenSpec contexts only where implementation happens.

## What Changes

- Strengthen the Codex `sys-apply` skill so it explicitly invokes the project OpenSpec apply workflow for the named change before implementation.
- Require `sys-apply` to use Superpowers methods during implementation planning, TDD, debugging, and verification when those skills are available.
- Update `sys change apply <name>` from passive guidance toward a concrete OpenSpec apply handoff.
- Make `sys init` create `frontend/` and `backend/` directories when missing and run `openspec init` inside both directories.
- Keep the monorepo root and `/system` free of new OpenSpec initialization during `sys init`.
- Route build-phase `sys change` commands to the inferred `frontend/` or `backend/` OpenSpec workspace instead of the monorepo root.
- Add tests that prove initialization targets only frontend/backend and that apply behavior routes through OpenSpec apply semantics.

## Capabilities

### New Capabilities

- None.

### Modified Capabilities

- `project-lifecycle`: `sys init` must initialize frontend/backend OpenSpec workspaces without initializing root or `/system`.
- `build-workflow`: `sys change apply` must invoke or dispatch the OpenSpec apply workflow rather than only printing loose guidance.
- `agent-integration`: generated Codex `sys-apply` instructions must require `openspec-apply` and Superpowers discipline with stronger wording and stop conditions.

## Impact

- Affected code: `internal/sysapp/app.go`, `internal/sysapp/app_test.go`.
- Affected templates: `internal/sysapp/templates/agents/codex/sys-apply/SKILL.md`.
- Affected docs: `README.md`.
- External dependency surface: existing `openspec` CLI becomes required during `sys init` unless tests inject a fake executable.
- Existing sys projects: running `sys init` in an already initialized repository remains idempotent and must not reinitialize unrelated locations.
