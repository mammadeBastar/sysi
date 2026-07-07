---
name: sysi-apply
description: Apply a sysi change in build phase using the native change workflow and Superpowers discipline.
---

## Purpose

Use this skill during build phase to implement a sysi change while preserving `/system` as the foundation truth. The change's own files are the work order: `proposal.md` says why, `design.md` says how, `tasks.md` says what remains. Superpowers governs the implementation/debug/test/verify loop.

## Initial Checks

1. Run or read `sysi status --json`.
2. Confirm the project is in build phase.
3. Confirm the current directory is inside a declared workspace and the named change exists under `<workspace>/changes/<name>/`.
4. Run `sysi change apply <name>` from the workspace to mark the change applying and print the handoff.
5. Confirm the relevant Superpowers workflows are available for implementation planning, TDD, systematic debugging, and verification.
6. Read the change's `proposal.md`, `design.md`, and `tasks.md` in full.
7. Read the relevant `/system` files allowed for the current role before editing implementation code, including `system/security/**` when security invariants affect the work.
8. Identify whether the requested implementation would introduce design drift from `/system` before changing behavior.

## Phase Rules

- Build phase is required for implementation.
- Design phase work should use `sysi-explore` and `sysi-capture` instead.
- The change's `tasks.md` owns task tracking during build.
- Running `sysi change apply <name>` is mandatory before implementation edits.
- Superpowers discipline is mandatory during implementation planning, test-driven development, systematic debugging, and verification.
- Frozen /system files are not implementation files.

## Role And File Access

- Role is the declared workspace containing the current working directory.
- Read the allowed `/system` files for that role before deciding how to implement.
- Treat `system/contracts/`, `system/flows/`, `system/modules/<workspace>.md`, `system/data/`, `system/obs/`, and `system/security/**` as build context when they affect the work.
- Keep implementation edits inside the current workspace.

## Workflow

1. Run `sysi change apply <name>` before editing implementation code.
2. Use Superpowers skills for implementation planning, TDD, debugging, and verification.
3. Treat a missing Superpowers workflow as a missing prerequisite and stop instead of implementing without it.
4. Work through `tasks.md` in order and check each task off only after implementation and verification.
5. Keep edits scoped to the change and the current task.
6. Compare implementation needs against `/system` truth before changing behavior.
7. Treat design drift as any implementation need that contradicts or extends foundation truth, including new or changed HTTP endpoints, request or response payload-shape changes, event contracts, auth/session/permission rules, shared error behavior, contract conventions, schema or data invariants, security invariants, metrics, logging, tracing, or alerting contracts.
8. If implementation reveals design drift, stop ordinary implementation work and explain the mismatch: what implementation needs, what `/system` currently says or omits, and which `/system` files likely own the truth.
9. Ask the user for explicit user confirmation before changing `/system` for detected drift.
10. If the user confirms the foundation change, run `sysi design-change <name>` and follow `sysi-design-change` before mutating controlled or frozen `/system` files.
11. If the user does not agree to the foundation change, do not continue implementation that contradicts `/system`; revise the change or implementation approach to fit current foundation truth.
12. When all tasks are checked and verified, run `sysi change archive <name>` from the workspace.

## Validation

- Run focused tests for the changed behavior.
- Run broader tests required by the change before completion.
- Re-read modified code and relevant `/system` files to check alignment.
- Confirm implementation respects contract conventions, error behavior, and security invariants when those files apply.
- Confirm detected design drift received user confirmation before any foundation mutation.
- Confirm agreed design drift went through `sysi design-change <name>` and `sysi-design-change` before controlled or frozen `/system` edits.
- Confirm no frozen /system files changed accidentally.
- Confirm `tasks.md` checkboxes accurately reflect completed work.

## Stop Conditions

- Stop if `sysi status --json` does not show build phase.
- Stop if the current directory is not inside a declared workspace.
- Stop if the named change is missing or archived.
- Stop if required Superpowers workflows are unavailable.
- Stop if the requested implementation contradicts `/system` truth.
- Stop if a foundation mutation is required and the user has not confirmed the drift.
- Stop if the user does not agree to a required foundation change and implementation would contradict `/system`.
- Stop if confirmed design drift has not gone through `sysi design-change <name>` and `sysi-design-change`.
- Stop if tests fail and systematic debugging has not isolated the cause.

## Do Not

- Do Not implement outside a sysi change during build phase.
- Do Not implement before running `sysi change apply <name>`.
- Do Not implement when a mandatory apply/debug/test/verify workflow is a missing prerequisite.
- Do Not mutate frozen /system files as part of normal apply.
- Do Not mutate `/system` for design drift without explicit user confirmation.
- Do Not continue implementation that contradicts `/system` when the user does not agree to the foundation change.
- Do Not treat new endpoints, payload shapes, auth rules, security invariants, data shapes, or observability contracts as ordinary implementation details when they are missing from `/system`.
- Do Not copy full Superpowers instructions into this skill; invoke or follow them.
- Do Not mark tasks complete without fresh verification.
- Do Not hide design drift by forcing code to fit an outdated proposal.
