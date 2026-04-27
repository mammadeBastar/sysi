---
name: sysi-design-change
description: Control foundational /system mutations during build phase.
---

## Purpose

Use this skill when build-phase work needs to mutate controlled or frozen `/system` truth. It is the ceremony for changing the project foundation after design freeze.

## Initial Checks

1. Run or read `sysi status --json`.
2. Run `sysi design-change <name>` or read its output.
3. Identify the current OpenSpec change, if any.
4. Read the affected `/system` files before proposing edits.
5. Check validation warnings and freeze mutations before changing anything.

## Phase Rules

- Build phase requires controlled mutation for foundational changes.
- Design phase usually does not need this skill; use `sysi-capture` for normal design capture.
- Controlled files require rationale before mutation.
- Frozen files require explicit user confirmation before mutation.

## Role And File Access

- Infer role from the current working directory.
- Read only the `/system` files needed to evaluate the requested foundation change.
- If the change affects cross-boundary behavior, include `system/contracts/`.
- If the change affects data, include `system/data/schema.sql` and `system/data/schema.md`.
- If the change affects security truth, include `system/security/`.
- If the change affects operations, include `system/obs/`.

## Workflow

1. State why normal apply cannot continue without changing foundation truth.
2. List affected `/system` files.
3. List impacted OpenSpec changes and implementation tasks.
4. Describe migration or compatibility notes for already-built code, data, APIs, or operations.
5. Ask for explicit user confirmation before editing controlled or frozen files.
6. Apply the smallest coherent `/system` mutation after confirmation.
7. Update or create a decision record when the foundation decision has lasting architectural impact.
8. Re-run validation before and after the mutation when possible.

## Validation

- Capture the before and after state of affected files in the conversation summary.
- Confirm edited `/system` files do not conflict.
- Confirm security changes are captured in `system/security/model.md` when they affect trust boundaries, sensitive data, secrets, or security invariants.
- Confirm impacted OpenSpec changes still describe the intended implementation.
- Confirm validation warnings are understood after the mutation.
- Confirm new foundation truth is precise enough for future agents.

## Stop Conditions

- Stop if the user has not given explicit user confirmation for controlled or frozen file edits.
- Stop if affected files cannot be identified.
- Stop if impacted OpenSpec changes are unknown and the mutation would change implementation scope.
- Stop if migration or compatibility notes are required but missing.

## Do Not

- Do Not mutate foundation truth silently.
- Do Not use this skill for ordinary code changes.
- Do Not make broad rewrites when a smaller foundation correction is enough.
- Do Not edit implementation code before the foundation change is accepted.
- Do Not leave OpenSpec artifacts inconsistent with changed `/system` truth.
