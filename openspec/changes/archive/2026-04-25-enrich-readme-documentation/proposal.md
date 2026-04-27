## Why

The current README is a concise quickstart, but the project now has enough workflow and command surface that it needs to function as complete user-facing documentation. A fuller README will help future agents and humans understand what `sysi` is, how the design/build phases work, how `/system` is governed, and how Codex/Cursor/Claude integrations should be used.

## What Changes

- Expand `README.md` from a short overview into full project documentation.
- Document the purpose, mental model, lifecycle phases, `/system` structure, command reference, agent integrations, validation, OpenSpec relationship, Superpowers relationship, and common workflows.
- Add practical examples for initializing a project, designing with `sysi explore`/`sysi capture`, freezing, applying build changes, installing agent integrations, and validating the repository.
- Clarify v1 boundaries and limitations, including first-class Codex support and minimal Cursor/Claude instruction support.
- Keep the documentation aligned with the current Go CLI behavior and synced OpenSpec specs.

## Capabilities

### New Capabilities
- `project-documentation`: Defines the README as comprehensive project documentation for users, contributors, and agents.

### Modified Capabilities
- None.

## Impact

- Updates `README.md` only.
- Does not change CLI behavior, APIs, generated files, OpenSpec workflow, or agent skill behavior.
- Adds documentation expectations that future changes should preserve when command or workflow behavior changes.
