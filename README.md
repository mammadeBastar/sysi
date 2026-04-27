# Sysi

`sysi` is a Go CLI for an agent-native monorepo workflow.

It gives agents and humans a durable system foundation before implementation starts. During design, decisions are captured directly into `/system`. During build, implementation changes flow through OpenSpec, while apply work is shaped by Superpowers-style planning, testing, debugging, and verification discipline.

The short version:

```text
/system      = ratified system truth
OpenSpec     = build-phase change protocol
Superpowers  = apply-phase engineering discipline
sysi          = CLI that ties those pieces together
```

`sysi` is intentionally pragmatic. It does not try to become a replacement for OpenSpec, a full documentation generator, or a hard filesystem sandbox. It gives the repository a clear lifecycle, a canonical system folder, validation, phase boundaries, and agent-native instructions.

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
OpenSpec build changes
        |
        v
Implementation with apply discipline
```

There are three important boundaries:

| Layer | Responsibility |
| --- | --- |
| `/system` | Current architecture, contracts, flows, modules, schema, security, observability, and decisions |
| OpenSpec | Build-phase changes to an already-understood system |
| Superpowers | Apply-phase method for planning, TDD, debugging, and verification |

Design phase is direct and decisive. Build phase is controlled and transactional.

## Project Lifecycle

Sysi projects move through two primary phases.

### Phase 1: Design

Design phase is where the system foundation is created. Agents explore, discuss, and make hard decisions. When a decision is final, it is captured into `/system`.

OpenSpec is intentionally not used for design decisions in this phase. The base system does not exist yet, so there is no meaningful implementation change to propose.

Typical flow:

```bash
sysi init
sysi status
sysi explore auth
sysi capture
sysi design freeze
```

### Phase 2: Build

Build phase starts after `sysi design freeze`. At this point, `/system` becomes controlled. Normal implementation work should go through OpenSpec.

Typical flow:

```bash
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

The CLI uses only the Go standard library in v1.

## Quick Start

Initialize a repository:

```bash
sysi init
```

This creates:

```text
.sysi/
system/
frontend/openspec/
backend/openspec/
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
│   └── openspec/
└── backend/
    └── openspec/
```

`sysi init` creates `.sysi/`, `system/`, and the `frontend/` and `backend/` implementation directories when they are missing. It runs non-interactive OpenSpec initialization inside `frontend/` and `backend/` only. The monorepo root and `/system` are not initialized as OpenSpec workspaces by `sysi init`.

### `.sysi/`

`.sysi/` stores repo-local operational state.

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

`sysi init` scaffolds:

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

Build phase uses OpenSpec for implementation changes. Run build change commands from the implementation workspace that owns the work, either `frontend/` or `backend/`. The sysi CLI still discovers the monorepo root from those directories, but it runs OpenSpec in the inferred implementation workspace.

Before using build commands, freeze the design:

```bash
sysi design freeze
```

### Propose A Change

```bash
cd frontend
sysi change propose add-login
```

This requires build phase and invokes OpenSpec from the current implementation workspace:

```bash
openspec new change add-login
```

If the `openspec` executable is not on `PATH`, set:

```bash
SYSI_OPENSPEC=/path/to/openspec sysi change propose add-login
```

### Apply A Change

```bash
cd frontend
sysi change apply add-login
```

This checks that `frontend/openspec/changes/add-login` exists, invokes from the frontend workspace:

```bash
openspec instructions apply --change add-login --json
```

and reports that implementation must continue through OpenSpec apply plus Superpowers discipline.

In Codex, use the generated `sysi-apply` skill. It requires the local OpenSpec apply workflow, `openspec-apply-change`, before implementation edits and requires Superpowers methods for planning, TDD, debugging, and verification.

If apply work reveals design drift from `/system`, the agent must stop ordinary implementation work and double-check the mismatch with the user. Examples include new endpoints, changed request or response payload shapes, event contract changes, auth/session/permission changes, shared error behavior, data-shape changes, security invariants, or observability contracts that are not represented in `/system`.

If the user agrees that foundation truth should change, the agent must use `sysi design-change <name>` and the generated `sysi-design-change` workflow before mutating controlled or frozen `/system` files. If the user does not agree, implementation must not continue in a way that contradicts `/system`.

### Archive A Change

```bash
cd frontend
sysi change archive add-login
```

This requires build phase and invokes OpenSpec archive from the current implementation workspace:

```bash
openspec archive add-login
```

### Foundation Changes During Build

If build work reveals that `/system` itself must change, use:

```bash
sysi design-change change-auth-boundary
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
sysi agent install codex
```

This creates:

```text
.codex/skills/sysi-explore/SKILL.md
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
| `sysi-explore` | Explore design questions from `/system`, surface candidate decisions, and avoid OpenSpec during design phase |
| `sysi-capture` | Write finalized design decisions into the right `/system` files and create decision records |
| `sysi-apply` | Apply OpenSpec changes in build phase by invoking OpenSpec apply first, then using mandatory Superpowers implementation discipline |
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

Cursor support is intentionally minimal in v1.

```bash
sysi agent install cursor
```

This writes:

```text
.cursor/rules/sysi.mdc
```

The file contains explicit workflow boundaries, phase rules, `/system` authority, OpenSpec build expectations, design-change protection, and role inference guidance. It is intentionally minimal and is not a deep runtime integration.

### Claude Code

Claude Code support is intentionally minimal in v1.

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

The managed section mirrors the minimal Cursor boundaries: it tells Claude Code how to respect design/build phases, `/system` truth, OpenSpec build workflow, `sysi design-change`, and inferred role access without claiming hard sandboxing or runtime enforcement.

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
- OpenSpec change summary
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
- `openspec`

### Watch Mode

```bash
sysi status --watch
```

This refreshes the terminal dashboard until interrupted.

### Validate

```bash
sysi validate
```

Validation checks required `/system` files and, in build phase, checks frozen/controlled file baselines.

If required files are missing, validation reports warnings such as:

```text
warning: missing required file: system/contracts/api.yaml
warning: missing required file: system/contracts/conventions.md
warning: missing required file: system/security/model.md
```

If a frozen or controlled file changes after `sysi design freeze`, status and validation report that `sysi design-change` is required. Controlled files include API, event, auth, conventions, errors, security model, and canonical schema files.

## Command Reference

### `sysi init`

Initializes a repo-local sysi project.

```bash
sysi init
```

Creates `.sysi/`, scaffolds `/system` including contracts and security files, creates `frontend/` and `backend/` when missing, initializes OpenSpec inside `frontend/` and `backend/`, records `design` phase, and prints the next command.

Running it again preserves existing state, reports that the project is already initialized, and ensures the frontend/backend OpenSpec workspaces still exist. Targets that already contain `openspec/config.yaml` are skipped.

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

Validates the required `/system` files and freeze baselines.

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

This command does not invoke OpenSpec.

### `sysi capture`

Prints capture guidance for finalized design decisions.

```bash
sysi capture
```

In build phase this command fails and points users to `sysi design-change`.

### `sysi design-change <name>`

Prints controlled mutation guidance for foundational `/system` changes during build phase.

```bash
sysi design-change change-auth-boundary
```

### `sysi change propose <name>`

Requires build phase and must be run from `frontend/` or `backend/`. Invokes OpenSpec in that implementation workspace to create a change.

```bash
sysi change propose add-login
```

### `sysi change apply <name>`

Requires build phase and must be run from `frontend/` or `backend/`. Checks for the OpenSpec change in that implementation workspace, invokes the OpenSpec apply instruction workflow there, and prints the required OpenSpec apply plus Superpowers handoff.

```bash
sysi change apply add-login
```

### `sysi change archive <name>`

Requires build phase and must be run from `frontend/` or `backend/`. Invokes OpenSpec archive in that implementation workspace.

```bash
sysi change archive add-login
```

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
sysi init
```

from the intended monorepo root.

### `build changes require build phase`

You ran a build-phase command before freezing design.

Run:

```bash
sysi design freeze
```

then retry the build command.

### `normal capture is blocked in build phase`

`sysi capture` is only for design phase.

During build phase, use:

```bash
sysi design-change <name>
```

### `openspec executable not found`

`sysi init` and build change commands that invoke OpenSpec require the `openspec` executable.

Either put `openspec` on `PATH`, or set:

```bash
SYSI_OPENSPEC=/path/to/openspec
```

### OpenSpec PostHog Network Errors

In restricted-network environments, OpenSpec can print telemetry errors after successful commands if it cannot reach `edge.openspec.dev`.

Those errors are not necessarily sysi or OpenSpec command failures. Check the command exit code and the main command output.

### Cursor And Claude Are Minimal

Codex receives project-local skills in v1. Cursor and Claude Code receive instruction files only.

This is intentional. The core protocol is stabilized first; richer adapters can be added later.

### No Generated `/system/views`

Sysi does not generate `/system/views` in v1. Role-specific context is expressed through `.sysi/allowlists.json` and agent instructions.

## Contributor Notes

Keep documentation aligned with behavior.

When a change affects commands, workflow, agent integration, phase behavior, or the `/system` scaffold:

1. Update the relevant OpenSpec specs.
2. Update `README.md`.
3. Add or update tests when behavior changes.
4. Run:

```bash
GOCACHE=/tmp/sysi-go-cache go test ./...
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
