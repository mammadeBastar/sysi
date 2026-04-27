---
name: sysi-capture
description: Capture finalized design decisions into /system during design phase.
---

## Purpose

Use this skill after the user has made a Finalized Decision and wants it written into `/system` so future agents inherit the same foundation.

## Initial Checks

1. Run or read `sysi status --json`.
2. Confirm the project is in design phase.
3. Run or read `sysi capture` for current capture guidance.
4. Infer role from the current working directory.
5. Read the target `/system` files before editing them.

## Phase Rules

- Design phase allows normal capture into `/system`.
- Build phase blocks normal capture. If status shows build phase, stop and use `sysi design-change`.
- Do not use OpenSpec for design-phase capture.
- Do not turn implementation discoveries into foundation mutations without a finalized user decision.

## Role And File Access

- Respect the inferred role and its allowed `/system` files.
- Frontend agents may update frontend modules, flows, security context, and contracts that affect frontend behavior.
- Backend agents may update backend modules, data, security context, observability, flows, and contracts that affect backend behavior.
- System-maintainer or design agents may update architecture and cross-cutting files when the user has finalized that foundation decision.
- If a needed file is outside the role allowlist, stop and ask before proceeding.

## Workflow

1. Identify the Finalized Decision in one sentence.
2. Choose the smallest set of target `/system` files that should own the truth.
3. Update the owning files first: architecture for invariants, contracts for boundaries, `system/contracts/conventions.md` for cross-cutting conventions, `system/contracts/errors.md` for error behavior, flows for user actions, modules for responsibilities, data for schemas, `system/security/model.md` for security posture, and observability for metrics/logs/traces/alerts.
4. Avoid duplicated truth by linking or summarizing instead of repeating full details across multiple files; explicitly avoid duplicated truth across `/system` files.
5. Add a Decision Record under `system/architecture/decisions/` when the decision affects architecture, contracts, data, observability, security, or cross-service behavior.
6. Keep prose direct and durable; future agents should be able to build from it without this conversation.

## Decision Record

Decision records should include:

- status
- decision
- rationale
- affected files
- consequences
- date

Use concise filenames under `system/architecture/decisions/`, such as `auth-session-boundary.md`.

## Validation

- Re-read every edited `/system` file.
- Confirm the edited files agree with each other.
- Confirm cross-boundary payloads are represented in `system/contracts/`.
- Confirm contract conventions are represented in `system/contracts/conventions.md` and error behavior is represented in `system/contracts/errors.md`.
- Confirm security posture changes are represented in `system/security/model.md`.
- Confirm data changes keep `system/data/schema.sql` and `system/data/schema.md` aligned.
- Confirm the decision record points to the files that now own the truth.

## Stop Conditions

- Stop if the decision is not finalized.
- Stop if the project is in build phase.
- Stop if required target files are outside allowed role access.
- Stop if two `/system` files would become conflicting sources of truth.

## Do Not

- Do Not capture brainstorming notes as final architecture.
- Do Not create OpenSpec changes during design phase.
- Do Not edit implementation files.
- Do Not duplicate full contracts, schemas, or flow details in unrelated files.
- Do Not bypass `sysi design-change` in build phase.
