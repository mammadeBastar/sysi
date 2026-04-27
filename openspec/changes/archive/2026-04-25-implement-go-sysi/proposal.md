## Why

The project needs a first working version of a Go-based orchestrator that turns the design conversation into durable `/system` truth and then coordinates implementation through OpenSpec and Superpowers. This should establish the project foundation before application code exists, while keeping the workflow usable from Codex and lightly portable to Cursor and Claude Code.

## What Changes

- Add a Go CLI named `sysi` with repo-local state, system scaffolding, phase management, status reporting, design capture support, build-phase change orchestration, validation, and agent integration installation.
- Introduce the canonical `/system` structure for architecture, contracts, flows, modules, Postgres-first schema, database details, observability, and decision records.
- Add design-phase commands that avoid OpenSpec and let agents explore, finalize, and capture decisions directly into `/system`.
- Add build-phase commands that use OpenSpec for implementation changes and require Superpowers discipline during apply workflows.
- Add Codex-native skill installation through `sysi agent install codex`.
- Add minimal Cursor and Claude Code support through generated instruction files only.
- Add an interactive `sysi status` dashboard plus JSON/text status output for agent consumption.

## Capabilities

### New Capabilities
- `project-lifecycle`: Initialize and manage a sysi-orchestrated monorepo with repo-local state, phase tracking, and conventional directories.
- `system-foundation`: Scaffold, maintain, freeze, and validate the canonical `/system` foundation files.
- `design-workflow`: Support design-phase exploration and capture without using OpenSpec for design decisions.
- `build-workflow`: Coordinate build-phase OpenSpec changes and Superpowers-influenced apply workflows.
- `agent-integration`: Install agent-native instructions, with full Codex skill support and minimal Cursor/Claude Code support.
- `status-dashboard`: Provide human and machine-readable project status, including an interactive terminal dashboard.

### Modified Capabilities
- None.

## Impact

- Adds a Go module and CLI entrypoint for `sysi`.
- Adds project scaffolding for `.sysi/`, `/system`, `.codex/skills/`, `.cursor/rules/`, and `CLAUDE.md`.
- Adds filesystem validation and status logic.
- Integrates with the existing `openspec` CLI by shelling out during build-phase commands.
- Does not add frontend/backend application code in v1.
