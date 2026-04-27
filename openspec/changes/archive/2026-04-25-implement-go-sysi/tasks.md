## 1. Go CLI Foundation

- [x] 1.1 Create the Go module, `cmd/sysi` entrypoint, and internal package layout.
- [x] 1.2 Implement a small command router using the Go standard library.
- [x] 1.3 Add common command output helpers for human text, errors, and JSON.
- [x] 1.4 Add baseline unit tests for command routing and error handling.

## 2. Repository State And Discovery

- [x] 2.1 Implement repository root discovery from nested directories.
- [x] 2.2 Implement `.sysi/state.json` read/write with default design phase.
- [x] 2.3 Implement `.sysi/freeze.json` and `.sysi/allowlists.json` models.
- [x] 2.4 Implement `sysi init` idempotently.
- [x] 2.5 Add tests for initialization, nested root discovery, and already-initialized repositories.

## 3. System Foundation

- [x] 3.1 Implement `/system` directory and file scaffolding.
- [x] 3.2 Scaffold Postgres-first `system/data/schema.sql` and explanatory `system/data/schema.md`.
- [x] 3.3 Scaffold `system/architecture/decisions/` and decision-record templates.
- [x] 3.4 Implement required-file validation for `/system`.
- [x] 3.5 Implement freeze baseline hashing for controlled and frozen `/system` files.
- [x] 3.6 Add tests for scaffolding, validation warnings, and freeze baseline detection.

## 4. Status Dashboard

- [x] 4.1 Implement `sysi status` human dashboard output.
- [x] 4.2 Implement `sysi status --json` with phase, root path, detected role, validation results, freeze state, and agent integration state.
- [x] 4.3 Implement `sysi status --watch` using repeated ANSI terminal refresh.
- [x] 4.4 Add uninitialized-project status handling.
- [x] 4.5 Add tests for status JSON and key dashboard sections.

## 5. Design Workflow

- [x] 5.1 Implement `sysi design start` to set or confirm design phase.
- [x] 5.2 Implement `sysi explore [topic]` to print design-agent instructions without invoking OpenSpec.
- [x] 5.3 Implement `sysi capture` design-phase guidance and build-phase blocking.
- [x] 5.4 Implement `sysi design freeze` to enter build phase and record freeze baselines.
- [x] 5.5 Implement role inference from current working directory.
- [x] 5.6 Add tests proving design commands do not call OpenSpec.

## 6. Build Workflow

- [x] 6.1 Implement OpenSpec executable detection and command execution wrapper.
- [x] 6.2 Implement `sysi change propose <name>` build-phase checks and OpenSpec handoff.
- [x] 6.3 Implement `sysi change apply <name>` checks and OpenSpec apply guidance.
- [x] 6.4 Implement `sysi change archive <name>` OpenSpec archive guidance and post-archive validation.
- [x] 6.5 Implement `sysi design-change <name>` guidance for controlled foundation mutation.
- [x] 6.6 Add tests using a fake OpenSpec executable.

## 7. Agent Integration

- [x] 7.1 Implement `sysi agent install codex`.
- [x] 7.2 Generate `.codex/skills/sysi-explore/SKILL.md`.
- [x] 7.3 Generate `.codex/skills/sysi-capture/SKILL.md`.
- [x] 7.4 Generate `.codex/skills/sysi-apply/SKILL.md`.
- [x] 7.5 Generate `.codex/skills/sysi-design-change/SKILL.md`.
- [x] 7.6 Implement `sysi agent install cursor` with minimal `.cursor/rules/sysi.mdc`.
- [x] 7.7 Implement `sysi agent install claude` with marked-section updates in `CLAUDE.md`.
- [x] 7.8 Add tests that agent installers preserve unrelated existing content.

## 8. Verification And Documentation

- [x] 8.1 Add CLI usage documentation for the v1 workflow.
- [x] 8.2 Run `go test ./...` and fix failures.
- [x] 8.3 Run OpenSpec status/validation for this change.
- [x] 8.4 Smoke test `sysi init`, `sysi status`, `sysi agent install codex`, and representative design/build commands in a temporary project.
