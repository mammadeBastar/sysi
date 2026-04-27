## Why

`sysi-apply` currently tells agents to pause when implementation reveals design drift, but it does not spell out the required user-confirmation and `sysi design-change` handoff. This leaves small contract shifts, such as new endpoints or payload-shape changes, too easy to treat as ordinary implementation details during build phase.

## What Changes

- Require build-phase apply work to classify design drift against `/system` before continuing implementation.
- Require agents to double-check detected drift with the user before changing foundation truth.
- Require confirmed drift to use `sysi design-change <name>` and the generated `sysi-design-change` workflow before mutating `/system`.
- Clarify that if the user does not agree to the foundation change, `sysi-apply` must not continue with implementation that contradicts `/system`.
- Update generated Codex `sysi-apply` instructions so this behavior is explicit in workflow, validation, stop conditions, and prohibited actions.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `build-workflow`: Require a user-confirmed `sysi design-change` handoff when apply work detects design drift from `/system`.
- `agent-integration`: Require generated `sysi-apply` instructions to double-check detected drift with the user and use `sysi-design-change` after agreement.

## Impact

- `openspec/specs/build-workflow/spec.md`
- `openspec/specs/agent-integration/spec.md`
- `internal/sysiapp/templates/agents/codex/sysi-apply/SKILL.md`
- `internal/sysiapp/app_test.go`
- `README.md`
