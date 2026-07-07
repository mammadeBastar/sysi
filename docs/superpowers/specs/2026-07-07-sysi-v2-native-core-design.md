# Sysi V2 Milestone 1: Native Core — Design

Date: 2026-07-07
Status: Approved for implementation

## Context And Motivation

Sysi v1 delegates build-phase change management to the external OpenSpec CLI and
hardcodes a `frontend/` + `backend/` monorepo layout. Real production use (an HA
project) showed OpenSpec was not earning its place, and public adoption requires
sysi to fit repos it has never seen.

Sysi v2 makes the tool self-contained. `/system` remains ratified truth; sysi
natively owns everything downstream of it.

### Roadmap (context, not M1 scope)

| Milestone | Content |
| --- | --- |
| M1 (this design) | Remove OpenSpec, declared workspaces, native change workflow |
| M2 | Multi-phase plans derived from `/system` (one file per phase) with plan skills |
| M3 | Gap analysis: plan coverage + change hygiene vs `/system` |
| M4 | Distribution: real module path, `go install`, release binaries |

### Explicit Non-Goals For M1

- No migration or compatibility shims for v1 projects. Existing v1 projects keep
  using the v1 binary. V2 is designed clean.
- No plans workflow (M2), no gap analysis (M3), no distribution work (M4).
- No implementation-reality scanning of code against `/system`.
- No deep Claude Code skill packs (possible later milestone).

## Design Principles

- CLI = structure, state, guardrails, validation. Skills = the intelligence.
- Predictable over clever: workspaces are declared, never guessed.
- Error messages do the teaching.
- Guardrails from v1 are preserved: design/build phase boundary, freeze
  baselines, drift forces `sysi design-change`, capture blocked in build phase.

## Section 1: Scope And Data Model

### Workspaces

- `sysi init --workspaces api,web,worker` declares workspaces. Each workspace
  directory is created if missing, with a `changes/` scaffold inside.
- Bare `sysi init` with no flag does not guess: it prints usage with examples
  and a suggested default (`--workspaces frontend,backend`).
- The workspace list is stored in `.sysi/state.json` (schema `version: 2`, new
  `workspaces` field).
- Role inference generalizes: cwd inside workspace `api` → role `api`;
  `system/` → `system-maintainer`; repo root → `design`.
- Allowlists use a generic per-workspace template instead of hardcoded
  frontend/backend entries.

### OpenSpec Removal

- `runOpenSpec`, `SYSI_OPENSPEC`, openspec-init during `sysi init`, and
  openspec workspace status reading are all deleted.
- No external binary is required for any sysi command.

### Native Change Storage

- Changes live in `<workspace>/changes/<name>/` with:
  - `proposal.md` — what and why
  - `design.md` — decisions with trade-offs
  - `tasks.md` — checkbox task list referencing `/system` files
  - `meta.json` — status and dates
- Archived changes move to `<workspace>/changes/archive/<date>-<name>/`.
- The layout deliberately mirrors OpenSpec's so prior users feel at home.

## Section 2: CLI Surface And Change Lifecycle

### Workspace Management Commands

```bash
sysi init --workspaces api,web        # declare at init; bare init prints usage
sysi workspace list                   # declared workspaces + active change counts
sysi workspace add worker             # create dir + changes/ scaffold, update state
sysi workspace remove worker          # remove from state; refuses if active changes exist (--force overrides)
```

### Change Lifecycle

Statuses in `meta.json`: `proposed` → `applying` → `archived`.

```bash
cd api
sysi change propose add-login    # scaffolds api/changes/add-login/, status=proposed
sysi change apply add-login      # status=applying; prints apply-discipline handoff
sysi change archive add-login    # moves to api/changes/archive/2026-07-07-add-login/, status=archived
```

- All change commands require build phase and must run from inside a declared
  workspace.
- `apply` prints the Superpowers discipline handoff (planning, TDD,
  verification) and the design-drift stop conditions.
- `archive` warns when `tasks.md` still has unchecked items.

### Status Dashboard

- The OpenSpec section is removed. Status shows, per declared workspace, active
  changes with their statuses.
- `--json` includes the same per-workspace change data so agents can enumerate
  changes without globbing.

### Scaffolded Templates

`proposal.md`, `design.md`, and `tasks.md` templates carry real structure
(what/why, decisions with trade-offs, checkbox task list referencing `/system`
files) — the same enriched-instruction philosophy as the Codex skill packs.

## Section 3: Skills, Validation, Errors, Testing, Docs

### Agent Skill Packs

- `sysi-apply` no longer invokes OpenSpec. It directs the agent to read the
  change's `proposal.md`/`design.md`/`tasks.md`, work through tasks with
  Superpowers discipline, check off tasks as completed, and honor the same
  design-drift stop conditions as v1.
- `sysi-explore`, `sysi-capture`, `sysi-design-change` lose OpenSpec references
  and speak in terms of declared workspaces instead of frontend/backend.
- Cursor and Claude instruction files get the same treatment.

### Validation

`sysi validate` keeps `/system` required-file and freeze-baseline checks, and
adds:

- every declared workspace directory exists
- every `changes/<name>/meta.json` parses
- statuses are legal values
- no active change name collides with an archived one

### Error Handling

- Change commands outside a declared workspace list the valid workspaces.
- `apply`/`archive` on a nonexistent change list the changes that do exist.
- `propose` with a duplicate name fails.
- `workspace remove` with active changes names them and requires `--force`.

### Testing

Removing the external binary removes all openspec stubbing. The whole lifecycle
is testable in a temp dir against the real code path:

- init with workspaces; bare init guidance
- workspace add / remove / remove-with-active-changes
- change lifecycle happy path: init → freeze → propose → apply → archive
- gating failures: wrong phase, outside workspace, duplicate name, missing change
- role inference across arbitrary workspace names
- status `--json` shape assertions
- archive warning on unchecked tasks

### Documentation

- README rewritten around v2: workspaces, native changes, no OpenSpec.
- The sysi repo's own OpenSpec specs (`build-workflow`, `project-lifecycle`,
  `agent-integration`, `status-dashboard`, `system-foundation` where touched)
  updated to match behavior, per contributor notes.
