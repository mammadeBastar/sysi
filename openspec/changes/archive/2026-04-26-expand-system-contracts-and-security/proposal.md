## Why

The current `/system` scaffold has the right shape, but its contract surface is missing durable defaults for API conventions and error behavior, and it has no first-class place for security posture. These gaps force agents to either overload `system/architecture/system.md` or make build-phase assumptions that should be settled during design.

## What Changes

- Add `system/contracts/errors.md` for error envelopes, codes, retryability, validation failures, and user-facing/internal error boundaries.
- Add `system/contracts/conventions.md` for cross-cutting API and event conventions such as pagination, filtering, idempotency, correlation IDs, timestamps, versioning, and deprecation rules.
- Add `system/security/model.md` as the minimal security foundation for trust boundaries, sensitive data, encryption expectations, secret handling, and security invariants.
- Update CLI scaffolding, validation, freeze baselines, allowlists, and capture/explore guidance so the new files are treated as canonical `/system` foundation.
- Update generated Codex, Cursor, and Claude instruction templates so agents know when to read or update the new contract and security files.
- Update README documentation to describe the expanded `/system` structure.

## Capabilities

### New Capabilities

- None.

### Modified Capabilities

- `system-foundation`: Expand the canonical `/system` scaffold, validation, controlled files, and role allowlists with contract conventions, contract errors, and security model files.
- `design-workflow`: Update design-phase CLI guidance so exploration and capture can target contract conventions, contract errors, and security model decisions.
- `agent-integration`: Update generated agent instruction packs to include the new contract and security foundation files in role guidance, workflow steps, validation, and guardrails.
- `project-documentation`: Update README requirements to document the expanded contracts and security structure.

## Impact

- Affected CLI code: `internal/sysiapp/app.go`.
- Affected tests: `internal/sysiapp/app_test.go`.
- Affected generated skills and agent instructions: `internal/sysiapp/templates/agents/**`.
- Affected docs: `README.md`.
- No external dependencies or runtime services are introduced.
