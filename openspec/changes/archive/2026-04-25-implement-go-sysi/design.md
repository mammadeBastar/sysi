## Context

This repository is a new sysi project with OpenSpec initialized and no application code yet. The product being introduced is a Go CLI named `sysi` that governs an agent-native monorepo workflow:

- During design phase, agents use `sysi explore` and `sysi capture` semantics to finalize decisions directly into `/system`; OpenSpec is intentionally not used for design.
- During build phase, OpenSpec owns implementation changes, `openspec-apply` remains the required apply path, and Superpowers discipline shapes the apply/debug/test/verify loop.
- `/system` is the ratified foundation: mostly mutable during design, controlled/frozen during build.
- Codex is the primary agent runtime for v1. Cursor and Claude Code receive minimal instruction-file support.

The initial repository is intentionally small, so v1 should avoid unnecessary dependencies and establish the protocol before optimizing the UI or adding deep integrations.

## Goals / Non-Goals

**Goals:**
- Implement a Go CLI named `sysi`.
- Scaffold `.sysi/` and the canonical `/system` tree.
- Track repo-local phase, freeze state, capture sessions, agent installs, and validation state.
- Provide `sysi status` as the primary terminal dashboard, including JSON output for agents.
- Provide design-phase commands that do not depend on OpenSpec.
- Provide build-phase commands that wrap or guide OpenSpec and Superpowers workflows.
- Install Codex-native skills with `sysi agent install codex`.
- Generate minimal Cursor and Claude Code instruction files with `sysi agent install cursor` and `sysi agent install claude`.

**Non-Goals:**
- Building frontend or backend application code.
- Enforcing OS-level filesystem sandboxing for agents.
- Creating deep Cursor or Claude Code runtime integrations.
- Implementing a full curses-style TUI framework in v1.
- Replacing OpenSpec, `openspec-apply`, or Superpowers.
- Auto-understanding arbitrary chat logs from the CLI alone.

## Decisions

### Use Go with standard-library-first implementation

The CLI will be a normal Go module with a `cmd/sysi` entrypoint and internal packages for commands, state, scaffolding, validation, agents, and OpenSpec execution. v1 should prefer the Go standard library for predictable bootstrapping in a new repo.

Alternatives considered:
- Cobra/Viper/Bubble Tea: better ergonomics and richer TUI, but increase dependency and setup weight before the protocol is stable.
- Shell scripts: faster to start, but weaker tests and harder cross-platform behavior.

### Store machine state under `.sysi/`

Repo-local machine state will live under `.sysi/`, with JSON files for v1:

```text
.sysi/
  state.json
  freeze.json
  allowlists.json
  captures/
  agents/
```

JSON keeps v1 dependency-free and easy to test. `/system` remains human-authored Markdown/YAML/SQL where appropriate. The state is operational metadata, not architectural truth.

Alternatives considered:
- YAML state: more human-friendly, but Go has no standard YAML parser.
- Global state: useful for a multi-repo registry, but wrong as the source of truth for a repo phase.

### Make `/system` canonical and small

`sysi init` will scaffold:

```text
system/
  architecture/
    system.md
    decisions/
  contracts/
    api.yaml
    events.asyncapi.yaml
    auth.md
  flows/
  modules/
    frontend.md
    backend.md
  data/
    schema.sql
    schema.md
    db/
      indexes.md
      triggers.md
      functions.md
  obs/
    metrics.md
    logging.md
    tracing.md
    alerts.md
    dashboards/
      grafana.md
```

There will be no generated `/system/views` in v1. Agent file access is expressed as allowlists and instructions.

Alternatives considered:
- Generated context views: convenient, but add drift and validation burden.
- A single system document: easier initially, but too coarse for targeted agent context and validation.

### Keep design phase independent from OpenSpec

`sysi explore` and Codex `sysi-explore` guide design conversations using the current `/system` state. `sysi capture` and Codex `sysi-capture` guide the agent to update `/system` and add a decision record when a decision is finalized. In design phase, this is direct and reversible through normal file history.

OpenSpec starts after the foundation exists, in build phase.

Alternatives considered:
- Using OpenSpec proposals during design: adds process overhead and misuses a change system before the base system exists.
- Capturing only at the end of the month: loses incremental memory and makes future agents restart from stale context.

### Treat build-phase foundation changes as explicit design changes

After `sysi design freeze`, controlled or frozen `/system` files require `sysi design-change` semantics. This command creates a heavier mutation path with rationale, impacted files, and build-change impact. Normal `sysi capture` is blocked in build phase.

Alternatives considered:
- Fully immutable `/system`: too rigid for real product evolution.
- Fully mutable `/system`: undermines the purpose of ratified foundation.

### Make Codex the first-class agent integration

`sysi agent install codex` will install project-local skills:

```text
.codex/skills/sysi-explore/SKILL.md
.codex/skills/sysi-capture/SKILL.md
.codex/skills/sysi-apply/SKILL.md
.codex/skills/sysi-design-change/SKILL.md
```

The skills will infer role from the current working directory and call/read CLI status where useful. Users should not need role-specific install commands.

Cursor and Claude Code support is intentionally minimal:

```text
.cursor/rules/sysi.mdc
CLAUDE.md
```

Alternatives considered:
- Equal first-class support for all runtimes: too much adapter work before the core protocol is proven.
- CLI-only integration: fails the requirement to work naturally from inside Codex.

### Implement `sysi status` as the dashboard command

`sysi status` will be the primary human dashboard. `sysi status --json` will be the machine-readable interface for skills and scripts. `sysi status --watch` will refresh the dashboard at an interval.

v1 can use ANSI terminal output and clear sections rather than a full TUI framework.

Alternatives considered:
- Separate `sysi dashboard` command: discoverable, but it splits the status concept. It can be added later as an alias if useful.

## Risks / Trade-offs

- CLI cannot read a Codex conversation by itself -> Codex `sysi-capture` skill will perform the actual extraction/editing, while the CLI supplies structure, status, templates, and validation.
- JSON state is less pleasant than YAML -> keep JSON hidden under `.sysi/` and reserve human-facing truth for `/system`.
- Minimal Cursor/Claude support may feel weaker than Codex -> document that v1 intentionally supports them through generated instructions only.
- Freeze enforcement based on recorded hashes can miss intent -> require `sysi validate` and clear status warnings; deeper git integration can come later.
- OpenSpec and Superpowers are external tools/instructions -> `sysi` should verify presence where possible and print exact next actions instead of pretending to replace them.

## Migration Plan

This is the first implementation. No migration is required.

Implementation can proceed in layers:

1. Create Go module, CLI entrypoint, command parser, and tests.
2. Add repo discovery and `.sysi` state.
3. Add `/system` scaffolding and validators.
4. Add status dashboard and JSON output.
5. Add design-phase command behavior.
6. Add build-phase OpenSpec command wrappers/checks.
7. Add agent installers and generated Codex skills.

Rollback is deleting generated artifacts from the working tree before adoption. Once users run `sysi init` in real projects, they should rely on git history for reversal.

## Open Questions

None for v1. The accepted boundary is full Go CLI plus first-class Codex support, with minimal Cursor and Claude Code instruction scaffolds.
