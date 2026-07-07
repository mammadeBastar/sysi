---
name: sysi-explore
description: Explore system design using /system as the current project truth. Does not create build changes during design phase.
---

## Purpose

Use this skill as a principal-engineer design review partner when the user wants to explore architecture, contracts, flows, modules, data, security, observability, or a feature idea before implementation. The goal is to turn vague design space into hard, testable decisions that can later be captured into `/system`.

## Reference Materials

- DDIA Mental Model Reference: read `references/ddia-mental-model.md` when the topic touches data models, persistence, encoding, schema evolution, replication, partitioning, transactions, consistency, batch processing, stream processing, derived data, scaling, or reliability. Use it as a checklist of mental models, not as quoted book notes.
- Keep `SKILL.md` as the operating procedure and references as expandable design lenses, so future book mental models can be added under `references/` without bloating this file.

## Initial Checks

1. Run or read `sysi status --json` when useful.
2. Identify the current phase, root, inferred role, validation warnings, and installed agent state.
3. Infer role from the current working directory: design, frontend, backend, system-maintainer, or change.
4. Read `.sysi/allowlists.json` when role access is unclear.
5. Read only the allowed `/system` files needed for the topic.
6. Read the DDIA reference when the design depends on data, persistence, concurrency, distribution, evolution, or derived state.

## Phase Rules

- Design phase is the normal phase for this skill.
- During design phase, do not create build changes for design exploration.
- During build phase, use this skill only for understanding existing foundation truth; do not mutate controlled or frozen files.
- If exploration reveals that build-phase foundation truth must change, stop and direct the user to `sysi design-change`.

## Role And File Access

- Respect the role inferred from the current working directory.
- Frontend work should usually read `system/architecture/system.md`, `system/contracts/**`, `system/flows/**`, `system/modules/frontend.md`, and `system/security/**`.
- Backend work should usually read `system/architecture/system.md`, `system/contracts/**`, `system/flows/**`, `system/modules/backend.md`, `system/data/**`, `system/security/**`, and `system/obs/**`.
- Design and system-maintainer work may read broader `/system` context.
- Treat allowed `/system` files as the source of truth and avoid reading unrelated implementation files unless the user asks for codebase investigation. In generated guidance and summaries, call these the allowed /system files for the current role.

## Design Review Lens

Scale depth to risk. For a small local feature, use the lens lightly. For a foundation, contract, data, security, or operational change, review it like a principal engineer:

- Source of truth: identify the canonical owner file for each fact and reject duplicate truth.
- Invariants: state what must always remain true across requests, retries, restarts, deployments, migrations, and partial failures.
- Failure modes: ask how the design breaks under invalid input, timeouts, unavailable dependencies, stale clients, corrupt state, clock skew, quota exhaustion, and operator error.
- Concurrency: check races, lost updates, ordering assumptions, reentrancy, idempotency keys, optimistic or pessimistic locking, and background worker overlap.
- Retries and idempotency: define which operations can be retried, which must not be retried, how duplicate submissions are detected, and what clients see after ambiguous success.
- Schema evolution: check backward and forward compatibility, nullable and default fields, migrations, rollbacks, versioned payloads, and old clients.
- Observability: require metrics, logs, traces, alerts, and dashboards for the behavior's success path and failure path.
- Scaling limits: identify cardinality, hot keys, fan-out, pagination, queue depth, retention, throughput ceilings, and load-shedding behavior.
- Security boundaries: name trust boundaries, attacker-controlled inputs, authorization checks, sensitive data, secret handling, audit needs, and denial-of-service exposure.
- Migration paths: describe how existing data, clients, APIs, jobs, and operators move from old truth to new truth.
- Operational recovery: define restore, replay, backfill, rollback, repair, manual override, and incident-diagnosis paths.

## Workflow

1. Restate the topic and the current phase/role.
2. Read the relevant `/system` files and summarize only the facts that matter.
3. Identify the source of truth, invariants, trust boundaries, data ownership, and cross-boundary contracts before comparing options.
4. Explore architecture, contracts, flows, modules, data, security, and observability implications as relevant.
5. For each viable option, state trade-offs, failure modes, migration cost, operational burden, and which `/system` files would own the resulting truth.
6. Ask focused questions only when a hard decision cannot be made from context.
7. Surface candidate decisions with enough precision that another agent can capture or implement them later.
8. When decisions are finalized, suggest `sysi-capture` and identify the likely target files.

## Validation

- Check that candidate decisions do not contradict existing `/system` truth.
- Check that every cross-boundary behavior belongs in `system/contracts/`.
- Check that contract conventions belong in `system/contracts/conventions.md` and error behavior belongs in `system/contracts/errors.md`.
- Check that trust boundaries, sensitive data rules, secret handling, and security invariants belong in `system/security/model.md`.
- Check that every user-visible action has or will need a matching `system/flows/*.flow.md` file.
- Check that data shape changes are reflected in both `system/data/schema.sql` and `system/data/schema.md` when relevant.
- Check that stateful designs have explicit concurrency, retry, idempotency, schema evolution, observability, scaling, migration, and operational recovery answers.
- Check that derived data has a named source, freshness expectation, recomputation path, and repair strategy.

## Stop Conditions

- Stop if the user asks for implementation instead of design exploration.
- Stop if required context is outside the role allowlist and the user has not approved broader reading.
- Stop if the decision would mutate frozen or controlled build-phase truth; use `sysi design-change`.
- Stop if a decision is not finalized; continue exploring instead of capturing.
- Stop if the design cannot identify its source of truth, critical invariants, or trust boundaries.

## Do Not

- Do Not implement code or edit application files.
- Do Not skip design work; avoid implementation until a build-phase sysi change exists.
- Do Not create build changes during design phase.
- Do Not capture tentative ideas as final truth.
- Do Not read broad codebase context when allowed `/system` files are sufficient.
- Do Not hide uncertainty; label open questions directly.
- Do Not accept a design that relies on retries, caches, background jobs, or derived state without defining correctness and recovery behavior.
