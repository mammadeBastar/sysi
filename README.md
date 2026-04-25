# Sys Orchestrator

`sys` is a Go CLI for an agent-native monorepo workflow.

It gives agents and humans a durable system foundation before implementation starts. During design, decisions are captured directly into `/system`. During build, implementation changes flow through OpenSpec, while apply work is shaped by Superpowers-style planning, testing, debugging, and verification discipline.

The short version:

```text
/system      = ratified system truth
OpenSpec     = build-phase change protocol
Superpowers  = apply-phase engineering discipline
sys          = CLI that ties those pieces together
```

`sys` is intentionally pragmatic. It does not try to become a replacement for OpenSpec, a full documentation generator, or a hard filesystem sandbox. It gives the repository a clear lifecycle, a canonical system folder, validation, phase boundaries, and agent-native instructions.

## Table Of Contents

- [Mental Model](#mental-model)
- [Project Lifecycle](#project-lifecycle)
- [Install And Run](#install-and-run)
- [Quick Start](#quick-start)
- [Repository Layout](#repository-layout)
- [The `/system` Foundation](#the-system-foundation)
- [Design Phase](#design-phase)
- [Build Phase](#build-phase)
- [Agent Integrations](#agent-integrations)
- [Status And Validation](#status-and-validation)
- [Command Reference](#command-reference)
- [Troubleshooting](#troubleshooting)
- [Contributor Notes](#contributor-notes)
- [V1 Boundaries](#v1-boundaries)

## Mental Model

Most agent-driven projects lose context between conversations. Sys Orchestrator solves that by making the project foundation explicit and durable.

```text
Design conversations
        |
        v
Final decisions
        |
        v
/system
ratified current truth
        |
        v
OpenSpec build changes
        |
        v
Implementation with apply discipline
```

There are three important boundaries:

| Layer | Responsibility |
| --- | --- |
| `/system` | Current architecture, contracts, flows, modules, schema, observability, and decisions |
| OpenSpec | Build-phase changes to an already-understood system |
| Superpowers | Apply-phase method for planning, TDD, debugging, and verification |

Design phase is direct and decisive. Build phase is controlled and transactional.

## Project Lifecycle

Sys projects move through two primary phases.

### Phase 1: Design

Design phase is where the system foundation is created. Agents explore, discuss, and make hard decisions. When a decision is final, it is captured into `/system`.

OpenSpec is intentionally not used for design decisions in this phase. The base system does not exist yet, so there is no meaningful implementation change to propose.

Typical flow:

```bash
sys init
sys status
sys explore auth
sys capture
sys design freeze
```

### Phase 2: Build

Build phase starts after `sys design freeze`. At this point, `/system` becomes controlled. Normal implementation work should go through OpenSpec.

Typical flow:

```bash
sys change propose add-login
sys change apply add-login
sys change archive add-login
```

If a build task reveals that a foundational system decision must change, use:

```bash
sys design-change change-auth-boundary
```

Normal `sys capture` is blocked in build phase.

## Install And Run

This repository is a Go module.

Run from source:

```bash
go run ./cmd/sys --help
go run ./cmd/sys <command>
```

Build a local binary:

```bash
go build -o sys ./cmd/sys
./sys --help
```

If `sys` is installed on your `PATH`, use:

```bash
sys <command>
```

The CLI uses only the Go standard library in v1.

## Quick Start

Initialize a repository:

```bash
sys init
```

This creates:

```text
.sys-orchestrator/
system/
frontend/openspec/
backend/openspec/
```

Check project status:

```bash
sys status
```

Start design work:

```bash
sys explore auth
```

Capture finalized decisions:

```bash
sys capture
```

Freeze the design foundation:

```bash
sys design freeze
```

Install Codex integration:

```bash
sys agent install codex
```

Validate the system foundation:

```bash
sys validate
```

## Repository Layout

A sys-initialized repository contains:

```text
.
├── .sys-orchestrator/
│   ├── state.json
│   ├── freeze.json
│   ├── allowlists.json
│   ├── captures/
│   └── agents/
├── system/
│   ├── architecture/
│   ├── contracts/
│   ├── flows/
│   ├── modules/
│   ├── data/
│   └── obs/
├── frontend/
│   └── openspec/
└── backend/
    └── openspec/
```

`sys init` creates `.sys-orchestrator/`, `system/`, and the `frontend/` and `backend/` implementation directories when they are missing. It runs non-interactive OpenSpec initialization inside `frontend/` and `backend/` only. The monorepo root and `/system` are not initialized as OpenSpec workspaces by `sys init`.

### `.sys-orchestrator/`

`.sys-orchestrator/` stores repo-local operational state.

| File | Purpose |
| --- | --- |
| `state.json` | Tracks phase, version, timestamps, and installed agent integrations |
| `freeze.json` | Stores freeze baselines for controlled `/system` files |
| `allowlists.json` | Stores role-based system-file allowlists |
| `captures/` | Reserved for capture-related metadata |
| `agents/` | Reserved for agent integration metadata |

This is machine state. The architectural truth belongs in `/system`.

## The `/system` Foundation

`/system` is the canonical project foundation. It is designed for agents and humans to read before they build.

`sys init` scaffolds:

```text
system/
├── architecture/
│   ├── system.md
│   └── decisions/
├── contracts/
│   ├── api.yaml
│   ├── events.asyncapi.yaml
│   └── auth.md
├── flows/
├── modules/
│   ├── frontend.md
│   └── backend.md
├── data/
│   ├── schema.sql
│   ├── schema.md
│   └── db/
│       ├── indexes.md
│       ├── triggers.md
│       └── functions.md
└── obs/
    ├── metrics.md
    ├── logging.md
    ├── tracing.md
    ├── alerts.md
    └── dashboards/
        └── grafana.md
```

### Architecture

`system/architecture/system.md` describes the high-level system:

- services and applications
- responsibilities
- communication patterns
- technical decisions
- system-wide invariants

`system/architecture/decisions/` stores decision records. Decision records should explain what was decided, why, status, and affected files.

### Contracts

`system/contracts/` contains everything that crosses boundaries:

- `api.yaml` for HTTP API contracts
- `events.asyncapi.yaml` for event contracts
- `auth.md` for authentication and authorization boundaries

Frontend agents are expected to rely heavily on this folder.

### Flows

`system/flows/` is for user and system action flows. Each flow should describe the movement of data and service interactions for one user action or system behavior.

Examples:

```text
system/flows/login.flow.md
system/flows/create-project.flow.md
system/flows/export-report.flow.md
```

### Modules

`system/modules/frontend.md` describes frontend pages, components, responsibilities, and dependencies.

`system/modules/backend.md` describes backend services, modules, responsibilities, and dependencies.

### Data

`system/data/schema.sql` is the canonical Postgres schema file.

`system/data/schema.md` explains relationships, invariants, schema rationale, protobuf notes, and any data-shape context that is better expressed in prose.

Database-specific operational details live in:

```text
system/data/db/indexes.md
system/data/db/triggers.md
system/data/db/functions.md
```

### Observability

`system/obs/` records observability expectations:

- metrics exposed to `/metrics`
- logging strategy
- tracing strategy
- alert rules
- intended Grafana dashboard structure

## Design Phase

Design phase is for creating and refining the foundation.

### Explore

Use `sys explore` to print design-agent guidance based on the current project state.

```bash
sys explore
sys explore auth
sys explore "billing events"
```

The command reports:

- topic
- current phase
- inferred agent role
- allowed `/system` files for that role
- reminder not to create OpenSpec changes during design phase

Role is inferred from the current working directory:

| Directory | Inferred Role |
| --- | --- |
| repo root | `design` |
| `frontend/` | `frontend` |
| `backend/` | `backend` |
| `system/` | `system-maintainer` |
| `openspec/changes/` | `change` |

### Capture

Use `sys capture` after a decision is finalized.

```bash
sys capture
```

The CLI prints capture rules. The active agent should then:

1. Update the relevant `/system` files.
2. Add a decision record under `system/architecture/decisions/`.
3. Keep the decision focused and explicit.

A decision record should include:

- status
- decision
- rationale
- affected files

During build phase, `sys capture` fails and directs users to `sys design-change`.

### Freeze

Freeze the design foundation when the project is ready for build work:

```bash
sys design freeze
```

This sets the phase to `build` and records hash baselines for controlled files:

```text
system/architecture/system.md
system/contracts/api.yaml
system/contracts/events.asyncapi.yaml
system/contracts/auth.md
system/data/schema.sql
```

`system/architecture/system.md` is treated as frozen. The other listed files are controlled.

## Build Phase

Build phase uses OpenSpec for implementation changes. Run build change commands from the implementation workspace that owns the work, either `frontend/` or `backend/`. The sys CLI still discovers the monorepo root from those directories, but it runs OpenSpec in the inferred implementation workspace.

Before using build commands, freeze the design:

```bash
sys design freeze
```

### Propose A Change

```bash
cd frontend
sys change propose add-login
```

This requires build phase and invokes OpenSpec from the current implementation workspace:

```bash
openspec new change add-login
```

If the `openspec` executable is not on `PATH`, set:

```bash
SYS_OPENSPEC=/path/to/openspec sys change propose add-login
```

### Apply A Change

```bash
cd frontend
sys change apply add-login
```

This checks that `frontend/openspec/changes/add-login` exists, invokes from the frontend workspace:

```bash
openspec instructions apply --change add-login --json
```

and reports that implementation must continue through OpenSpec apply plus Superpowers discipline.

In Codex, use the generated `sys-apply` skill. It requires the local OpenSpec apply workflow, `openspec-apply-change`, before implementation edits and requires Superpowers methods for planning, TDD, debugging, and verification.

### Archive A Change

```bash
cd frontend
sys change archive add-login
```

This requires build phase and invokes OpenSpec archive from the current implementation workspace:

```bash
openspec archive add-login
```

### Foundation Changes During Build

If build work reveals that `/system` itself must change, use:

```bash
sys design-change change-auth-boundary
```

This prints the required design-change guidance:

- rationale
- affected `/system` files
- impacted OpenSpec changes
- migration notes

## Agent Integrations

Agent integration is installed once per project. Runtime role is inferred from the current working directory.

### Codex

Codex is the first-class v1 integration.

```bash
sys agent install codex
```

This creates:

```text
.codex/skills/sys-explore/SKILL.md
.codex/skills/sys-capture/SKILL.md
.codex/skills/sys-apply/SKILL.md
.codex/skills/sys-design-change/SKILL.md
```

These are full instruction packs, not placeholder files. Each skill includes:

- purpose and when to use it
- initial `sys status` checks
- phase rules for design and build work
- role and `/system` file-access guidance
- workflow steps
- validation expectations
- stop conditions
- explicit "do not" guardrails

The installed Codex skills are:

| Skill | Purpose |
| --- | --- |
| `sys-explore` | Explore design questions from `/system`, surface candidate decisions, and avoid OpenSpec during design phase |
| `sys-capture` | Write finalized design decisions into the right `/system` files and create decision records |
| `sys-apply` | Apply OpenSpec changes in build phase by invoking OpenSpec apply first, then using mandatory Superpowers implementation discipline |
| `sys-design-change` | Mutate controlled or frozen `/system` truth during build phase only after explicit confirmation |

Typical Codex usage:

```text
[$sys-explore]
Design auth and sessions.

[$sys-capture]
Capture the finalized auth decisions.

[$sys-apply]
Apply add-login.

[$sys-design-change]
Change the auth boundary during build phase.
```

### Cursor

Cursor support is intentionally minimal in v1.

```bash
sys agent install cursor
```

This writes:

```text
.cursor/rules/sys-orchestrator.mdc
```

The file contains explicit workflow boundaries, phase rules, `/system` authority, OpenSpec build expectations, design-change protection, and role inference guidance. It is intentionally minimal and is not a deep runtime integration.

### Claude Code

Claude Code support is intentionally minimal in v1.

```bash
sys agent install claude
```

This creates or updates a marked section in:

```text
CLAUDE.md
```

Existing unrelated `CLAUDE.md` content is preserved. The sys section is bounded by:

```text
<!-- SYS-ORCHESTRATOR:START -->
<!-- SYS-ORCHESTRATOR:END -->
```

The managed section mirrors the minimal Cursor boundaries: it tells Claude Code how to respect design/build phases, `/system` truth, OpenSpec build workflow, `sys design-change`, and inferred role access without claiming hard sandboxing or runtime enforcement.

## Status And Validation

### Human Status

```bash
sys status
```

The dashboard shows:

- root path
- phase
- inferred role
- system health
- freeze baselines
- OpenSpec change summary
- installed agent integrations
- validation warnings

### JSON Status

```bash
sys status --json
```

This is useful for agents and scripts. It includes:

- `root`
- `phase`
- `role`
- `validation`
- `freeze`
- `agents`
- `openspec`

### Watch Mode

```bash
sys status --watch
```

This refreshes the terminal dashboard until interrupted.

### Validate

```bash
sys validate
```

Validation checks required `/system` files and, in build phase, checks frozen/controlled file baselines.

If required files are missing, validation reports warnings such as:

```text
warning: missing required file: system/contracts/api.yaml
```

If a frozen file changes after `sys design freeze`, status and validation report that `sys design-change` is required.

## Command Reference

### `sys init`

Initializes a repo-local sys project.

```bash
sys init
```

Creates `.sys-orchestrator/`, scaffolds `/system`, creates `frontend/` and `backend/` when missing, initializes OpenSpec inside `frontend/` and `backend/`, records `design` phase, and prints the next command.

Running it again preserves existing state, reports that the project is already initialized, and ensures the frontend/backend OpenSpec workspaces still exist. Targets that already contain `openspec/config.yaml` are skipped.

### `sys status`

Prints the human dashboard.

```bash
sys status
```

Use JSON output for agents:

```bash
sys status --json
```

Use watch mode for repeated refresh:

```bash
sys status --watch
```

### `sys validate`

Validates the required `/system` files and freeze baselines.

```bash
sys validate
```

### `sys design start`

Sets or confirms design phase.

```bash
sys design start
```

### `sys design freeze`

Moves the project into build phase and records freeze baselines.

```bash
sys design freeze
```

### `sys explore [topic]`

Prints design-agent guidance.

```bash
sys explore
sys explore auth
```

This command does not invoke OpenSpec.

### `sys capture`

Prints capture guidance for finalized design decisions.

```bash
sys capture
```

In build phase this command fails and points users to `sys design-change`.

### `sys design-change <name>`

Prints controlled mutation guidance for foundational `/system` changes during build phase.

```bash
sys design-change change-auth-boundary
```

### `sys change propose <name>`

Requires build phase and must be run from `frontend/` or `backend/`. Invokes OpenSpec in that implementation workspace to create a change.

```bash
sys change propose add-login
```

### `sys change apply <name>`

Requires build phase and must be run from `frontend/` or `backend/`. Checks for the OpenSpec change in that implementation workspace, invokes the OpenSpec apply instruction workflow there, and prints the required OpenSpec apply plus Superpowers handoff.

```bash
sys change apply add-login
```

### `sys change archive <name>`

Requires build phase and must be run from `frontend/` or `backend/`. Invokes OpenSpec archive in that implementation workspace.

```bash
sys change archive add-login
```

### `sys agent install codex`

Installs Codex project-local skills.

```bash
sys agent install codex
```

### `sys agent install cursor`

Installs minimal Cursor rules.

```bash
sys agent install cursor
```

### `sys agent install claude`

Creates or updates the marked sys section in `CLAUDE.md`.

```bash
sys agent install claude
```

## Troubleshooting

### `sys project not initialized`

You are outside a directory tree containing:

```text
.sys-orchestrator/state.json
```

Run:

```bash
sys init
```

from the intended monorepo root.

### `build changes require build phase`

You ran a build-phase command before freezing design.

Run:

```bash
sys design freeze
```

then retry the build command.

### `normal capture is blocked in build phase`

`sys capture` is only for design phase.

During build phase, use:

```bash
sys design-change <name>
```

### `openspec executable not found`

`sys init` and build change commands that invoke OpenSpec require the `openspec` executable.

Either put `openspec` on `PATH`, or set:

```bash
SYS_OPENSPEC=/path/to/openspec
```

### OpenSpec PostHog Network Errors

In restricted-network environments, OpenSpec can print telemetry errors after successful commands if it cannot reach `edge.openspec.dev`.

Those errors are not necessarily sys or OpenSpec command failures. Check the command exit code and the main command output.

### Cursor And Claude Are Minimal

Codex receives project-local skills in v1. Cursor and Claude Code receive instruction files only.

This is intentional. The core protocol is stabilized first; richer adapters can be added later.

### No Generated `/system/views`

Sys does not generate `/system/views` in v1. Role-specific context is expressed through `.sys-orchestrator/allowlists.json` and agent instructions.

## Contributor Notes

Keep documentation aligned with behavior.

When a change affects commands, workflow, agent integration, phase behavior, or the `/system` scaffold:

1. Update the relevant OpenSpec specs.
2. Update `README.md`.
3. Add or update tests when behavior changes.
4. Run:

```bash
GOCACHE=/tmp/sys-orchestrator-go-cache go test ./...
openspec validate --specs
```

Use the `/tmp` Go cache form in sandboxed environments where the default Go cache is not writable.

Documentation-only changes should still verify OpenSpec status and validation.

## V1 Boundaries

V1 intentionally does not provide:

- hard OS-level filesystem sandboxing for agents
- generated `/system/views`
- deep Cursor runtime integration
- deep Claude Code runtime integration
- a full curses-style terminal UI
- replacement behavior for OpenSpec
- replacement behavior for Superpowers
- automatic chat-log extraction from the CLI alone

The intended v1 shape is smaller and stricter:

```text
Go CLI
repo-local state
canonical /system scaffold
design/build phase boundary
status and validation
Codex skills
minimal Cursor/Claude instructions
OpenSpec build handoff
Superpowers apply discipline
```
