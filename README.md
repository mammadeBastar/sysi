# Sysi

`sysi` is a Go CLI for an agent-native monorepo workflow.

It gives agents and humans a durable system foundation before implementation starts. During design, decisions are captured directly into `/system`. During build, implementation changes flow through sysi's native change workflow inside declared workspaces, while apply work is shaped by Superpowers-style planning, testing, debugging, and verification discipline.

The short version:

```text
/system      = ratified system truth
sysi changes = build-phase change protocol (native)
Superpowers  = apply-phase engineering discipline
sysi          = CLI that ties those pieces together
```

`sysi` is intentionally pragmatic. It does not try to become a full documentation generator or a hard filesystem sandbox. It gives the repository a clear lifecycle, a canonical system folder, declared workspaces, a native change protocol, validation, phase boundaries, and agent-native instructions.

## Table Of Contents

- [Mental Model](#mental-model)
- [Project Lifecycle](#project-lifecycle)
- [Install And Run](#install-and-run)
- [Quick Start](#quick-start)
- [Repository Layout](#repository-layout)
- [Workspaces](#workspaces)
- [The `/system` Foundation](#the-system-foundation)
- [Design Phase](#design-phase)
- [Build Phase](#build-phase)
- [Agent Integrations](#agent-integrations)
- [Status And Validation](#status-and-validation)
- [Command Reference](#command-reference)
- [Troubleshooting](#troubleshooting)
- [Contributor Notes](#contributor-notes)
- [V2 Boundaries](#v2-boundaries)

## Mental Model

Most agent-driven projects lose context between conversations. Sysi solves that by making the project foundation explicit and durable.

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
sysi workspace changes
        |
        v
Implementation with apply discipline
```

There are three important boundaries:

| Layer | Responsibility |
| --- | --- |
| `/system` | Current architecture, contracts, flows, modules, schema, security, observability, and decisions |
| `sysi changes` | Build-phase changes inside declared workspaces |
| Superpowers | Apply-phase method for planning, TDD, debugging, and verification |

Design phase is direct and decisive. Build phase is controlled and transactional.

## Project Lifecycle

Sysi projects move through two primary phases.

### Phase 1: Design

Design phase is where the system foundation is created. Agents explore, discuss, and make hard decisions. When a decision is final, it is captured into `/system`.

Build changes are intentionally not created in this phase. The base system does not exist yet, so there is no meaningful implementation change to propose.

Typical flow:

```bash
sysi init --workspaces frontend,backend
sysi status
sysi explore auth
sysi capture
sysi design freeze
```

### Phase 2: Build

Build phase starts after `sysi design freeze`. At this point, `/system` becomes controlled. Normal implementation work goes through sysi changes, run from the workspace that owns the work.

Typical flow:

```bash
cd frontend
sysi change propose add-login
sysi change apply add-login
sysi change archive add-login
```

If a build task reveals that a foundational system decision must change, use:

```bash
sysi design-change change-auth-boundary
```

Normal `sysi capture` is blocked in build phase.

## Install And Run

This repository is a Go module.

Run from source:

```bash
go run ./cmd/sysi --help
go run ./cmd/sysi <command>
```

Build a local binary:

```bash
go build -o sysi ./cmd/sysi
./sysi --help
```

If `sysi` is installed on your `PATH`, use:

```bash
sysi <command>
```

The CLI uses only the Go standard library and has no external dependencies. No other executable is required for any sysi command.

## Quick Start

Initialize a repository with its declared workspaces:

```bash
sysi init --workspaces frontend,backend
```

Workspaces are never guessed. A bare `sysi init` on an uninitialized repository prints usage guidance with examples and fails, so declare the implementation directories explicitly.

Initialization creates:

```text
.sysi/
system/
frontend/changes/
backend/changes/
```

Check project status:

```bash
sysi status
```

Start design work:

```bash
sysi explore auth
```

Capture finalized decisions:

```bash
sysi capture
```

Add another workspace later if the system grows:

```bash
sysi workspace add worker
```

Freeze the design foundation:

```bash
sysi design freeze
```

Install Codex integration:

```bash
sysi agent install codex
```

Validate the system foundation:

```bash
sysi validate
```

## Repository Layout

A sysi-initialized repository contains:

```text
.
├── .sysi/
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
│   ├── security/
│   └── obs/
├── frontend/
│   └── changes/
│       └── archive/
└── backend/
    └── changes/
        └── archive/
```

`sysi init --workspaces <names>` creates `.sysi/`, `system/`, and one directory per declared workspace with a `changes/` scaffold inside. Active build changes live under `<workspace>/changes/<name>/`; completed changes move to `<workspace>/changes/archive/<date>-<name>/` (the archive directory is created on first archive).

### `.sysi/`

`.sysi/` stores repo-local operational state.

| File | Purpose |
| --- | --- |
| `state.json` | Version 2 state: tracks phase, timestamps, declared `workspaces`, and installed agent integrations |
| `freeze.json` | Stores freeze baselines for controlled `/system` files |
| `allowlists.json` | Stores role-based system-file allowlists, one role per declared workspace plus design roles |
| `captures/` | Reserved for capture-related metadata |
| `agents/` | Reserved for agent integration metadata |

`state.json` uses schema `version: 2`. States written by sysi v1 are rejected with an explanatory error; v1 projects should keep using the v1 binary.

This is machine state. The architectural truth belongs in `/system`.

## Workspaces

Workspaces are the implementation directories where build changes live. They are declared explicitly, first at `sysi init --workspaces <names>` and later with the `sysi workspace` commands. Sysi never guesses a repository layout.

Workspace names are lowercase slugs: they start with a lowercase letter and may contain lowercase letters, digits, and hyphens. The names `system`, `docs`, `openspec`, `design`, and `system-maintainer` are reserved.

List declared workspaces and their active change counts:

```bash
sysi workspace list
```

Add a workspace:

```bash
sysi workspace add worker
```

This validates the name, rejects duplicates and names that collide with an existing file at the repo root, creates `worker/changes/`, scaffolds `system/modules/worker.md`, and updates the allowlists and state.

Remove a workspace:

```bash
sysi workspace remove worker
```

Removal refuses if the workspace has active changes; it names them and requires `--force` to proceed:

```bash
sysi workspace remove worker --force
```

Removing a workspace only removes it from sysi state. The directory and its contents are left on disk.

## The `/system` Foundation

`/system` is the canonical project foundation. It is designed for agents and humans to read before they build.

`sysi init` scaffolds (module files shown for `--workspaces frontend,backend`):

```text
system/
├── architecture/
│   ├── system.md
│   └── decisions/
├── contracts/
│   ├── api.yaml
│   ├── events.asyncapi.yaml
│   ├── auth.md
│   ├── conventions.md
│   └── errors.md
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
├── security/
│   └── model.md
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
- `conventions.md` for cross-cutting API and event conventions such as pagination, filtering, idempotency, correlation IDs, timestamps, versioning, deprecation, and rate-limit expression
- `errors.md` for error envelopes, error codes, retryability, validation failures, and user-facing/internal error boundaries

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

`system/modules/<workspace>.md` describes one module file per declared workspace: its components, responsibilities, and dependencies. Init scaffolds a module file for every declared workspace, and `sysi workspace add` scaffolds one for each workspace added later.

### Data

`system/data/schema.sql` is the canonical Postgres schema file.

`system/data/schema.md` explains relationships, invariants, schema rationale, protobuf notes, and any data-shape context that is better expressed in prose.

Database-specific operational details live in:

```text
system/data/db/indexes.md
system/data/db/triggers.md
system/data/db/functions.md
```

### Security

`system/security/model.md` records security posture:

- trust boundaries
- sensitive data rules
- encryption expectations
- secret handling
- security invariants
- threat assumptions

This file documents rules and assumptions. It must not contain real secret values.

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

Use `sysi explore` to print design-agent guidance based on the current project state.

```bash
sysi explore
sysi explore auth
sysi explore "billing events"
```

The command reports:

- topic
- current phase
- inferred agent role
- allowed `/system` files for that role
- reminder not to create build changes during design phase

Role is inferred from the current working directory:

| Directory | Inferred Role |
| --- | --- |
| repo root | `design` |
| `system/` | `system-maintainer` |
| inside declared workspace `<name>` | `<name>` |

### Capture

Use `sysi capture` after a decision is finalized.

```bash
sysi capture
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

During build phase, `sysi capture` fails and directs users to `sysi design-change`.

### Freeze

Freeze the design foundation when the project is ready for build work:

```bash
sysi design freeze
```

This sets the phase to `build` and records hash baselines for controlled files:

```text
system/architecture/system.md
system/contracts/api.yaml
system/contracts/events.asyncapi.yaml
system/contracts/auth.md
system/contracts/conventions.md
system/contracts/errors.md
system/security/model.md
system/data/schema.sql
```

`system/architecture/system.md` is treated as frozen. The other listed files are controlled.

## Build Phase

Build phase uses sysi's native change workflow. Changes are plain files inside the workspace that owns the work; no external tool is involved. Run change commands from inside a declared workspace directory — from anywhere else they fail and list the declared workspaces.

Before using build commands, freeze the design:

```bash
sysi design freeze
```

A change lives at `<workspace>/changes/<name>/` and contains:

| File | Purpose |
| --- | --- |
| `proposal.md` | Why the change exists, what changes, foundation alignment, out of scope |
| `design.md` | Decisions with alternatives, interfaces, risks |
| `tasks.md` | Checkbox task list worked in order during apply |
| `meta.json` | Machine state: name, workspace, status, timestamps |

Change statuses progress `proposed` → `applying` → `archived`.

Change names must be lowercase slugs (lowercase letters, digits, hyphens). A name that duplicates an active change or collides with an archived change in the same workspace is rejected, and the name `archive` is reserved.

### Propose A Change

```bash
cd frontend
sysi change propose add-login
```

This requires build phase and scaffolds `frontend/changes/add-login/` with templated `proposal.md`, `design.md`, and `tasks.md` plus a `meta.json` with status `proposed`. Fill the three documents before applying.

### Apply A Change

```bash
cd frontend
sysi change apply add-login
```

This marks the change `applying` and prints the apply handoff: read `proposal.md`, `design.md`, and `tasks.md` before editing implementation code, work tasks in order with Superpowers discipline (planning, TDD, systematic debugging, verification), and check off each task only after implementation and verification.

The handoff also prints the design-drift stop conditions. If apply work reveals design drift from `/system`, the agent must stop ordinary implementation work and double-check the mismatch with the user. Examples include new or changed endpoints, payload shapes, event contracts, auth/session/permission rules, shared error behavior, schema or data invariants, security invariants, or observability contracts that are not represented in `/system`.

If the user agrees that foundation truth should change, the agent must use `sysi design-change <name>` and the generated `sysi-design-change` workflow before mutating controlled or frozen `/system` files. If the user does not agree, implementation must not continue in a way that contradicts `/system`.

In Codex, use the generated `sysi-apply` skill. It requires running `sysi change apply <name>` before implementation edits and requires Superpowers methods for planning, TDD, debugging, and verification.

### Archive A Change

```bash
cd frontend
sysi change archive add-login
```

This moves the change to `frontend/changes/archive/<date>-add-login/` and sets its status to `archived`. Archiving warns when `tasks.md` still has unchecked tasks and when the change has an unexpected status; it fails if the dated archive target already exists. If updating the archived `meta.json` fails, the move is rolled back.

### Foundation Changes During Build

If build work reveals that `/system` itself must change, use:

```bash
sysi design-change change-auth-boundary
```

This prints the required design-change guidance and creates a dated artifact under `system/architecture/decisions/` recording:

- rationale
- affected `/system` files
- impacted workspace changes
- migration notes

## Agent Integrations

Agent integration is installed once per project. Runtime role is inferred from the current working directory.

### Codex

Codex is the first-class integration.

```bash
sysi agent install codex
```

This creates:

```text
.codex/skills/sysi-explore/SKILL.md
.codex/skills/sysi-explore/references/ddia-mental-model.md
.codex/skills/sysi-capture/SKILL.md
.codex/skills/sysi-apply/SKILL.md
.codex/skills/sysi-design-change/SKILL.md
```

These are full instruction packs, not placeholder files. Each skill includes:

- purpose and when to use it
- initial `sysi status` checks
- phase rules for design and build work
- role and `/system` file-access guidance
- workflow steps
- validation expectations
- stop conditions
- explicit "do not" guardrails

The installed Codex skills are:

| Skill | Purpose |
| --- | --- |
| `sysi-explore` | Explore design questions from `/system`, surface candidate decisions, and avoid creating build changes during design phase |
| `sysi-capture` | Write finalized design decisions into the right `/system` files and create decision records |
| `sysi-apply` | Apply a sysi change in build phase: run `sysi change apply`, read the change's `proposal.md`, `design.md`, and `tasks.md`, then work the tasks with mandatory Superpowers implementation discipline |
| `sysi-design-change` | Mutate controlled or frozen `/system` truth during build phase only after explicit confirmation |

Typical Codex usage:

```text
[$sysi-explore]
Design auth and sessions.

[$sysi-capture]
Capture the finalized auth decisions.

[$sysi-apply]
Apply add-login.

[$sysi-design-change]
Change the auth boundary during build phase.
```

### Cursor

Cursor support is intentionally minimal.

```bash
sysi agent install cursor
```

This writes:

```text
.cursor/rules/sysi.mdc
```

The file contains explicit workflow boundaries, phase rules, `/system` authority, the native change workflow expectation (`sysi change propose|apply|archive` from the owning workspace), design-change protection, and role inference guidance. It is intentionally minimal and is not a deep runtime integration.

### Claude Code

Claude Code support is intentionally minimal.

```bash
sysi agent install claude
```

This creates or updates a marked section in:

```text
CLAUDE.md
```

Existing unrelated `CLAUDE.md` content is preserved. The sysi section is bounded by:

```text
<!-- SYSI:START -->
<!-- SYSI:END -->
```

The managed section mirrors the minimal Cursor boundaries: it tells Claude Code how to respect design/build phases, `/system` truth, the native change workflow inside declared workspaces, `sysi design-change`, and inferred role access without claiming hard sandboxing or runtime enforcement.

## Status And Validation

### Human Status

```bash
sysi status
```

The dashboard shows:

- root path
- phase
- inferred role
- system health
- freeze baselines
- per-workspace change summary with each change's status
- installed agent integrations
- validation warnings

### JSON Status

```bash
sysi status --json
```

This is useful for agents and scripts. It includes:

- `root`
- `phase`
- `role`
- `validation`
- `freeze`
- `agents`
- `workspaces`

`workspaces` is an array with one entry per declared workspace: `name`, `present` (whether the directory exists), `activeChanges`, and `changes`, a list of `{name, status}` objects. Agents can enumerate active changes from it without globbing the filesystem. Empty lists are emitted as `[]`, never `null`.

### Watch Mode

```bash
sysi status --watch
```

This refreshes the terminal dashboard until interrupted.

### Validate

```bash
sysi validate
```

Validation checks:

- required `/system` files, including `system/modules/<workspace>.md` for every declared workspace
- every declared workspace directory exists and is a directory
- every active change has a parseable `meta.json`
- every active change has a legal status (`proposed` or `applying`)
- no active change name collides with an archived change name
- in build phase, frozen/controlled file baselines

If required files are missing, validation reports warnings such as:

```text
warning: missing required file: system/contracts/api.yaml
warning: missing required file: system/contracts/conventions.md
warning: missing required file: system/security/model.md
```

If a frozen or controlled file changes after `sysi design freeze`, status and validation report that `sysi design-change` is required. Controlled files include API, event, auth, conventions, errors, security model, and canonical schema files.

## Command Reference

### `sysi init --workspaces <names>`

Initializes a repo-local sysi project with declared workspaces.

```bash
sysi init --workspaces frontend,backend
```

The equals form is also accepted:

```bash
sysi init --workspaces=frontend,backend
```

Creates `.sysi/` with version 2 state, scaffolds `/system` including contracts, security, and one `system/modules/<workspace>.md` per declared workspace, creates each workspace directory with a `changes/` scaffold, records `design` phase, and prints the next command. Workspace names that conflict with an existing file are rejected before anything is written.

`--workspaces` is required on first init: a bare `sysi init` on an uninitialized repository prints usage guidance and fails without writing any state.

Running `sysi init` again in an initialized repository preserves existing state, re-ensures the `/system` scaffold, allowlists, and workspace `changes/` directories, and reports that the project is already initialized.

### `sysi status`

Prints the human dashboard.

```bash
sysi status
```

Use JSON output for agents:

```bash
sysi status --json
```

Use watch mode for repeated refresh:

```bash
sysi status --watch
```

### `sysi validate`

Validates the required `/system` files, workspace and change health, and freeze baselines.

```bash
sysi validate
```

### `sysi design start`

Sets or confirms design phase.

```bash
sysi design start
```

### `sysi design freeze`

Moves the project into build phase and records freeze baselines.

```bash
sysi design freeze
```

### `sysi explore [topic]`

Prints design-agent guidance.

```bash
sysi explore
sysi explore auth
```

This command does not create build changes.

### `sysi capture`

Prints capture guidance for finalized design decisions.

```bash
sysi capture
```

In build phase this command fails and points users to `sysi design-change`.

### `sysi design-change <name>`

Requires build phase and creates a controlled mutation artifact for foundational `/system` changes.

```bash
sysi design-change change-auth-boundary
```

The artifact is created under `system/architecture/decisions/<date>-change-auth-boundary.md` and is the working record for rationale, affected `/system` files, impacted workspace changes (its `## Impacted Changes` section), migration notes, confirmation, decision, and consequences.

### `sysi workspace list`

Lists declared workspaces with their active change counts.

```bash
sysi workspace list
```

### `sysi workspace add <name>`

Declares a new workspace.

```bash
sysi workspace add worker
```

Validates the name (lowercase slug, not reserved, not already declared, no file conflict), creates `<name>/changes/`, scaffolds `system/modules/<name>.md`, and updates allowlists and state.

### `sysi workspace remove <name> [--force]`

Removes a workspace from sysi state.

```bash
sysi workspace remove worker
```

Refuses if the workspace has active changes; `--force` overrides. The directory is left on disk.

### `sysi change propose <name>`

Requires build phase and must be run from inside a declared workspace. Scaffolds `<workspace>/changes/<name>/` with `proposal.md`, `design.md`, `tasks.md`, and `meta.json` at status `proposed`.

```bash
sysi change propose add-login
```

The name must be a lowercase slug and must not duplicate an active change or collide with an archived change in the workspace.

### `sysi change apply <name>`

Requires build phase and must be run from inside a declared workspace. Marks the change `applying` and prints the Superpowers apply handoff plus the design-drift stop conditions.

```bash
sysi change apply add-login
```

If the change does not exist, the error lists the changes that do. Archived changes cannot be applied.

### `sysi change archive <name>`

Requires build phase and must be run from inside a declared workspace. Moves the change to `<workspace>/changes/archive/<date>-<name>/` and sets its status to `archived`.

```bash
sysi change archive add-login
```

Warns on unchecked tasks in `tasks.md` and on unexpected statuses; fails if the archive target already exists.

### `sysi agent install codex`

Installs Codex project-local skills.

```bash
sysi agent install codex
```

### `sysi agent install cursor`

Installs minimal Cursor rules.

```bash
sysi agent install cursor
```

### `sysi agent install claude`

Creates or updates the marked sysi section in `CLAUDE.md`.

```bash
sysi agent install claude
```

## Troubleshooting

### `sysi project not initialized`

You are outside a directory tree containing:

```text
.sysi/state.json
```

Run:

```bash
sysi init --workspaces <names>
```

from the intended monorepo root.

### `sysi init requires declared workspaces`

You ran a bare `sysi init` on an uninitialized repository and it failed with `error: missing --workspaces`. Sysi never guesses workspaces; declare them explicitly:

```bash
sysi init --workspaces frontend,backend
```

Workspaces are the implementation directories where build changes live. They can be adjusted later with `sysi workspace add|remove`.

### `build changes require build phase`

You ran a build-phase command before freezing design.

Run:

```bash
sysi design freeze
```

then retry the build command.

### `change commands must run inside a declared workspace`

You ran `sysi change propose|apply|archive` from the repo root, `system/`, or another directory that is not inside a declared workspace. The error lists the declared workspaces.

Run the command from the workspace that owns the change:

```bash
cd frontend
sysi change propose add-login
```

If the directory should be a workspace but is not declared, declare it first with `sysi workspace add <name>`.

### `normal capture is blocked in build phase`

`sysi capture` is only for design phase.

During build phase, use:

```bash
sysi design-change <name>
```

### Cursor And Claude Are Minimal

Codex receives project-local skills. Cursor and Claude Code receive instruction files only.

This is intentional. The core protocol is stabilized first; richer adapters can be added later.

### No Generated `/system/views`

Sysi does not generate `/system/views`. Role-specific context is expressed through `.sysi/allowlists.json` and agent instructions.

## Contributor Notes

Keep documentation aligned with behavior.

When a change affects commands, workflow, agent integration, phase behavior, or the `/system` scaffold:

1. Update the relevant spec files under `openspec/specs/`.
2. Update `README.md`.
3. Add or update tests when behavior changes.
4. Run:

```bash
GOCACHE=/tmp/sysi-go-cache go test ./...
```

Use the `/tmp` Go cache form in sandboxed environments where the default Go cache is not writable.

Documentation-only changes should still run the test suite.

## V2 Boundaries

V2 intentionally does not provide:

- hard OS-level filesystem sandboxing for agents
- generated `/system/views`
- deep Cursor runtime integration
- deep Claude Code runtime integration
- a full curses-style terminal UI
- replacement behavior for Superpowers
- automatic chat-log extraction from the CLI alone
- a multi-phase plan workflow yet (planned M2)
- gap analysis of plan coverage and change hygiene yet (planned M3)

The intended v2 shape is smaller and stricter:

```text
Go CLI with no external dependencies
repo-local state
canonical /system scaffold
declared workspaces
design/build phase boundary
native build changes
status and validation
Codex skills
minimal Cursor/Claude instructions
Superpowers apply discipline
```
