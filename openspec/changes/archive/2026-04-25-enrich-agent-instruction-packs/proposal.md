## Why

Sysi is meant to be an agent-native workflow, but the current generated Codex, Cursor, and Claude instructions are too skeletal to reliably guide agents through design capture, build apply, foundation mutation, and role-scoped system access. The next version should make agent instruction packs explicit, comprehensive, and testable so the CLI does not only create files, but installs an actual agent operating procedure.

## What Changes

- Replace tiny embedded Codex skill strings with comprehensive instruction pack templates for `sysi-explore`, `sysi-capture`, `sysi-apply`, and `sysi-design-change`.
- Add shared agent guidance covering phase rules, role inference, `/system` authority, file allowlists, OpenSpec boundaries, Superpowers boundaries, validation expectations, and stop conditions.
- Keep `sysi agent install codex` as the one-command install path, but make the generated skills detailed enough to guide real agent work.
- Enrich Cursor and Claude Code generated instructions while keeping them minimal compared with Codex.
- Move generated instruction content out of ad hoc Go string functions into repo-native templates or template-like constants that are easier to inspect and maintain.
- Add tests that verify generated instruction packs contain required operational sections and guardrails, not merely that files exist.
- Update README documentation to describe the stronger agent instruction model.

## Capabilities

### New Capabilities
- None.

### Modified Capabilities
- `agent-integration`: Strengthen agent installation requirements so generated instructions are comprehensive, role-aware, phase-aware, and explicit about workflow guardrails.

## Impact

- Updates Go implementation for agent installation content and possibly template loading.
- Updates tests for Codex, Cursor, and Claude instruction generation.
- Updates `README.md` to document the richer agent instruction packs.
- Updates the `agent-integration` OpenSpec capability.
- Does not change the command surface: `sysi agent install codex|cursor|claude` remains the user API.
