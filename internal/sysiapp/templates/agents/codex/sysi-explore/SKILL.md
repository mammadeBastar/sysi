---
name: sysi-explore
description: Explore system design using /system as the current project truth. Does not create OpenSpec changes during design phase.
---

## Purpose

Use this skill when the user wants to explore architecture, contracts, flows, modules, data, security, observability, or a feature idea before implementation. The goal is to help the user reach hard design decisions that can later be captured into `/system`.

## Initial Checks

1. Run or read `sysi status --json` when useful.
2. Identify the current phase, root, inferred role, validation warnings, and installed agent state.
3. Infer role from the current working directory: design, frontend, backend, system-maintainer, or change.
4. Read `.sysi/allowlists.json` when role access is unclear.
5. Read only the allowed `/system` files needed for the topic.

## Phase Rules

- Design phase is the normal phase for this skill.
- During design phase, do not create OpenSpec changes for design exploration.
- During build phase, use this skill only for understanding existing foundation truth; do not mutate controlled or frozen files.
- If exploration reveals that build-phase foundation truth must change, stop and direct the user to `sysi design-change`.

## Role And File Access

- Respect the role inferred from the current working directory.
- Frontend work should usually read `system/architecture/system.md`, `system/contracts/**`, `system/flows/**`, `system/modules/frontend.md`, and `system/security/**`.
- Backend work should usually read `system/architecture/system.md`, `system/contracts/**`, `system/flows/**`, `system/modules/backend.md`, `system/data/**`, `system/security/**`, and `system/obs/**`.
- Design and system-maintainer work may read broader `/system` context.
- Treat allowed /system files as the source of truth and avoid reading unrelated implementation files unless the user asks for codebase investigation.

## Workflow

1. Restate the topic and the current phase/role.
2. Read the relevant `/system` files and summarize only the facts that matter.
3. Explore architecture, contracts, flows, modules, data, security, and observability implications as relevant.
4. Ask focused questions only when a hard decision cannot be made from context.
5. Surface candidate decisions with trade-offs and the exact `/system` files they would affect.
6. Keep design choices precise enough that another agent can continue later.
7. When decisions are finalized, suggest `sysi-capture` and identify the likely target files.

## Validation

- Check that candidate decisions do not contradict existing `/system` truth.
- Check that every cross-boundary behavior belongs in `system/contracts/`.
- Check that contract conventions belong in `system/contracts/conventions.md` and error behavior belongs in `system/contracts/errors.md`.
- Check that trust boundaries, sensitive data rules, secret handling, and security invariants belong in `system/security/model.md`.
- Check that every user-visible action has or will need a matching `system/flows/*.flow.md` file.
- Check that data shape changes are reflected in both `system/data/schema.sql` and `system/data/schema.md` when relevant.

## Stop Conditions

- Stop if the user asks for implementation instead of design exploration.
- Stop if required context is outside the role allowlist and the user has not approved broader reading.
- Stop if the decision would mutate frozen or controlled build-phase truth; use `sysi design-change`.
- Stop if a decision is not finalized; continue exploring instead of capturing.

## Do Not

- Do Not implement code or edit application files.
- Do Not skip design work; avoid implementation until a build-phase OpenSpec change exists.
- Do Not create OpenSpec changes during design phase.
- Do Not capture tentative ideas as final truth.
- Do Not read broad codebase context when allowed `/system` files are sufficient.
- Do Not obscure uncertainty; label open questions directly.
