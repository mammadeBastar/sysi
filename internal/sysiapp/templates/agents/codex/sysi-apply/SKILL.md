---
name: sysi-apply
description: Apply OpenSpec changes using sysi workflow, openspec-apply, and Superpowers discipline.
---

## Purpose

Use this skill during build phase to implement an OpenSpec change while preserving `/system` as the foundation truth. This skill is a mandatory composition layer: OpenSpec apply starts the implementation path, and Superpowers governs the implementation/debug/test/verify loop.

## Initial Checks

1. Run or read `sysi status --json`.
2. Confirm the project is in build phase.
3. Confirm the named OpenSpec change exists.
4. Confirm the local OpenSpec apply workflow is available. In Codex, invoke or read `openspec-apply-change` for the named change before editing implementation code.
5. Confirm the relevant Superpowers workflows are available for implementation planning, TDD, systematic debugging, and verification.
6. Read the OpenSpec proposal, design, specs, and tasks for the change.
7. Read the relevant `/system` files allowed for the current role before editing implementation code.

## Phase Rules

- Build phase is required for implementation.
- Design phase work should use `sysi-explore` and `sysi-capture` instead.
- OpenSpec owns change planning and task tracking during build.
- OpenSpec apply is mandatory before implementation edits. In Codex, use `openspec-apply-change`.
- Superpowers discipline is mandatory during implementation planning, test-driven development, systematic debugging, and verification.
- Frozen /system files are not implementation files.

## Role And File Access

- Infer role from the current working directory.
- Read the allowed `/system` files for that role before deciding how to implement.
- Frontend agents should treat `system/contracts/`, `system/flows/`, `system/modules/frontend.md`, and `system/security/**` as their build context when security invariants affect the work.
- Backend agents should treat `system/contracts/`, `system/flows/`, `system/modules/backend.md`, `system/data/`, `system/security/**`, and `system/obs/` as their build context when relevant.
- Change agents may read OpenSpec change files and the `/system` files required by that change.

## Workflow

1. Invoke the local OpenSpec apply workflow for the named change before editing implementation code. In Codex, use `openspec-apply-change`.
2. Use Superpowers skills for implementation planning, TDD, debugging, and verification.
3. Treat a missing OpenSpec apply or Superpowers workflow as a missing prerequisite and stop instead of implementing without it.
4. Work through OpenSpec tasks in order and mark each task complete only after implementation and verification.
5. Keep edits scoped to the change and the current task.
6. Compare implementation needs against `/system` truth before changing behavior.
7. If implementation reveals design drift, pause and explain the mismatch.
8. Escalate to `sysi design-change` before mutating foundational truth.

## Validation

- Run focused tests for the changed behavior.
- Run broader tests required by the OpenSpec change before completion.
- Re-read modified code and relevant `/system` files to check alignment.
- Confirm implementation respects contract conventions, error behavior, and security invariants when those files apply.
- Confirm no frozen /system files changed accidentally.
- Confirm OpenSpec task checkboxes accurately reflect completed work.

## Stop Conditions

- Stop if `sysi status --json` does not show build phase.
- Stop if the OpenSpec change is missing or blocked.
- Stop if `openspec-apply-change` or the local OpenSpec apply workflow is unavailable.
- Stop if required Superpowers workflows are unavailable.
- Stop if the requested implementation contradicts `/system` truth.
- Stop if a foundation mutation is required; use `sysi design-change`.
- Stop if tests fail and systematic debugging has not isolated the cause.

## Do Not

- Do Not implement outside an OpenSpec change during build phase.
- Do Not implement before invoking the OpenSpec apply workflow.
- Do Not implement when a mandatory apply/debug/test/verify workflow is a missing prerequisite.
- Do Not mutate frozen /system files as part of normal apply.
- Do Not copy full OpenSpec or Superpowers instructions into this skill; invoke or follow them.
- Do Not mark tasks complete without fresh verification.
- Do Not hide design drift by forcing code to fit an outdated spec.
