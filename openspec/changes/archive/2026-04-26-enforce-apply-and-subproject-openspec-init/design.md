## Context

The current implementation establishes the right vocabulary but leaves two enforcement points too soft:

- `sys change apply <name>` verifies the change directory and prints a reminder instead of entering an OpenSpec apply handoff.
- `sys-apply` tells Codex to "invoke or follow" OpenSpec apply and to use Superpowers "when available", which allows agents to treat both as optional guidance.
- `sys init` creates sys state and `/system`, but not the frontend/backend OpenSpec workspaces that will own implementation changes.

The user wants `sys` to behave like an orchestrator: initialization prepares implementation workspaces, and apply explicitly composes OpenSpec apply plus Superpowers discipline.

## Goals / Non-Goals

**Goals:**

- Make `sys init` create `frontend/` and `backend/` directories when missing.
- Run non-interactive `openspec init` in both `frontend/` and `backend/`.
- Avoid creating or initializing OpenSpec in the root sys project or under `/system` as part of `sys init`.
- Make `sys change apply <name>` dispatch to the OpenSpec apply instruction path instead of only printing loose guidance.
- Harden generated Codex `sys-apply` wording so OpenSpec apply and Superpowers are mandatory for apply work.
- Add tests covering targeted OpenSpec initialization and apply dispatch.

**Non-Goals:**

- Do not implement OS-level filesystem sandboxing for frontend/backend agents.
- Do not install Superpowers automatically.
- Do not replace the OpenSpec apply skill or copy its full instructions into sys.
- Do not initialize OpenSpec for arbitrary monorepo packages beyond `frontend/` and `backend/`.
- Do not add role-specific agent install commands.

## Decisions

### Initialize OpenSpec Only In Implementation Subprojects

`sys init` should ensure `frontend/` and `backend/` exist, then run:

```bash
openspec init frontend --tools none
openspec init backend --tools none
```

The exact command should be issued through the existing `runOpenSpec` wrapper or a small variant that can run from the sys root. `--tools none` keeps initialization non-interactive and avoids installing agent-specific OpenSpec files into these subprojects by surprise.

Alternative considered: initialize OpenSpec at the monorepo root. That conflicts with the intended ownership model: `/system` is root-level foundation truth, while OpenSpec tracks implementation changes inside frontend/backend.

### Make Init Idempotent Across Existing Sys Repositories

If `.sys-orchestrator/state.json` already exists, `sys init` should preserve existing sys state but still ensure frontend/backend directories and OpenSpec initialization are present. If a target already contains `openspec/config.yaml`, skip reinitializing that target.

Alternative considered: keep the current early return. That would leave existing sys repositories unable to adopt the new frontend/backend OpenSpec layout with the normal `sys init` command.

### Dispatch Apply Through OpenSpec Apply Instructions

`sys change apply <name>` should continue requiring build phase and an existing change directory. After that, it should call the OpenSpec apply instruction path for the change, using:

```bash
openspec instructions apply --change <name> --json
```

This does not replace the agent skill that implements tasks. It proves that sys has entered the OpenSpec apply path and gives the agent the same schema-aware context that `openspec-apply-change` uses.

Alternative considered: only print `[$openspec-apply-change] <name>`. That is too advisory for the requested behavior and preserves the current weakness.

### Route Build Commands To Implementation Workspaces

Build-phase `sys change propose|apply|archive` should run OpenSpec from the inferred implementation workspace:

- commands invoked under `frontend/` use `frontend/`
- commands invoked under `backend/` use `backend/`
- commands invoked from root, `/system`, or other non-implementation paths fail with guidance to run from `frontend/` or `backend/`

This keeps the root and `/system` free of OpenSpec initialization while still letting the sys CLI find the repo state from nested implementation directories.

Alternative considered: keep running OpenSpec from the sys root. That contradicts the frontend/backend OpenSpec initialization model and would reintroduce a root OpenSpec dependency for normal build work.

### Make Codex Sys Apply A Mandatory Composition Layer

The `sys-apply` skill should say:

- The agent MUST invoke/read the local OpenSpec apply workflow for the named change before editing implementation code.
- In Codex, that workflow is the project OpenSpec apply skill, currently `openspec-apply-change`.
- Superpowers apply discipline is mandatory during implementation planning, test-driven work, debugging, and final verification.
- If required OpenSpec apply or Superpowers skills are unavailable in the runtime, the agent must stop and report the missing prerequisite instead of silently continuing.

Alternative considered: keep "when available" phrasing for portability. That is weaker than the intended Codex-first workflow.

## Risks / Trade-offs

- [Risk] `openspec init` may emit telemetry network errors in restricted environments even after successful initialization. Mitigation: tests should use a fake OpenSpec executable, and README troubleshooting should continue explaining telemetry noise.
- [Risk] Existing initialized repos may not expect `sys init` to run OpenSpec commands. Mitigation: skip targets that already contain `openspec/config.yaml` and keep state preservation behavior.
- [Risk] `sys change apply` cannot directly invoke a Codex skill from a Go process. Mitigation: invoke the OpenSpec apply instruction command in the CLI and make the Codex `sys-apply` skill explicitly invoke the OpenSpec apply skill before implementation.
- [Risk] Enforcing Superpowers in the Codex skill makes non-Codex runtimes less complete. Mitigation: keep Cursor/Claude minimal and document that Codex is the first-class integration for strict apply composition.

## Migration Plan

1. Update specs for lifecycle, build workflow, and agent integration.
2. Update `sys init` to ensure frontend/backend OpenSpec initialization.
3. Route `sys change` commands through the inferred frontend/backend OpenSpec workspace.
4. Update `sys change apply` to call OpenSpec apply instructions for the named change.
5. Harden `sys-apply` template wording and tests.
6. Update README command behavior and troubleshooting notes.

Rollback is straightforward: revert the init helper and apply dispatch changes, then restore prior template wording.

## Open Questions

- Should future versions allow custom implementation package names beyond `frontend` and `backend` through `.sys-orchestrator` config?
