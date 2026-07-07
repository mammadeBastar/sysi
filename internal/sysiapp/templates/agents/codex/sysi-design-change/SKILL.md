---
name: sysi-design-change
description: Control foundational /system mutations during build phase.
---

## Purpose

Use this skill when build-phase work needs to mutate controlled or frozen `/system` truth. It is the ceremony for changing the project foundation after design freeze.

## Initial Checks

1. Run or read `sysi status --json`.
2. Run `sysi design-change <name>` or read its output.
3. Open the created decision artifact under `system/architecture/decisions/<date>-<name>.md`.
4. Identify the current workspace change(s), if any.
5. Read the affected `/system` files before proposing edits.
6. Check validation warnings and freeze mutations before changing anything.

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

## Foundation Change Routing

Use this table before editing. The design-change artifact must name every owner file being changed, every file only being cross-linked, and every impacted workspace change.

| Path | Must Own | Must Not Contain | Cross-Link Instead When |
| --- | --- | --- | --- |
| `system/architecture/system.md` | System-wide architecture, major service/application responsibilities, communication patterns, deployment assumptions, technology decisions, and global invariants. | Endpoint schemas, event schemas, table DDL, flow steps, dashboard layouts, secret values, or implementation task lists. | The foundation change affects a boundary, schema, security invariant, operational signal, or module; summarize the architectural impact and link to the owner. |
| `system/architecture/decisions/*.md` | Rationale, explicit confirmation, affected files, impacted workspace changes, migration or compatibility notes, accepted decision, consequences, and date. | The live canonical copy of any contract, schema, security model, telemetry spec, or module definition. | Multiple owners change together; record why and point to each current owner. |
| `system/contracts/api.yaml` | HTTP contract changes: paths, methods, params, request/response payloads, status codes, and OpenAPI components. | Flow prose, DB DDL, module internals, shared error prose, auth policy prose, or migration plans. | HTTP behavior depends on conventions, errors, auth, data, flows, or security; link to the owner instead of duplicating. |
| `system/contracts/events.asyncapi.yaml` | Event contract changes: channels, operations, messages, producers, consumers, and AsyncAPI components. | HTTP routes, batch job internals, database schemas, operational runbooks, or replay procedures. | Event behavior depends on conventions, errors, data, observability, or flow semantics. |
| `system/contracts/auth.md` | Authentication, authorization, sessions, tokens, roles, permissions, and boundary rules visible through APIs/events. | Secret values, encryption internals, endpoint payload copies, or full threat models. | The change alters trust boundaries, sensitive data, or security invariants; link to `system/security/model.md`. |
| `system/contracts/conventions.md` | Shared API/event rules: pagination, filtering, idempotency, correlation IDs, timestamps, versioning, deprecation, rate-limit expression, compatibility, and retry semantics. | Endpoint-specific fields, one-off route/event behavior, data-model rationale, or flow narratives. | A concrete route/event needs the convention; keep the general rule here and reference it from the contract. |
| `system/contracts/errors.md` | Shared error envelope, error codes, retryability, validation shape, user-facing/internal boundary, and rate-limit error behavior. | Per-endpoint exhaustive error tables, stack traces, logging fields, UI copy, or incident procedures. | A route/event has specific errors; define shared shape here and reference it from `api.yaml` or `events.asyncapi.yaml`. |
| `system/flows/*.flow.md` | Changed behavior over time for one user or system action: actors, trigger, preconditions, data movement, service interactions, success path, failure branches, retries, and user-visible result. | Canonical API schemas, event schemas, table definitions, module inventories, or low-level code plans. | The flow depends on contracts, data, security, observability, or modules. |
| `system/modules/<workspace>.md` | The declared workspace's pages, routes, services, components, workers, state ownership, ownership boundaries, dependencies, coordination, and responsibility changes. | OpenAPI/AsyncAPI schema copies, database DDL, other workspaces' module internals, or dashboard/alert policy. | The module consumes or exposes contracts, implements a flow, owns data, emits telemetry, or enforces security. |
| `system/data/schema.sql` | Executable canonical Postgres schema changes: tables, columns, constraints, types, keys, foreign keys, and schema-level indexes. | Prose rationale, migration story, sample data, API payloads, or alternative schema drafts. | Schema evolution requires explanation, compatibility, or data lifecycle notes; put prose in `schema.md` and reference exact objects. |
| `system/data/schema.md` | Relationships, invariants, lifecycle, retention, schema rationale, protobuf notes, compatibility assumptions, and schema evolution explanation. | Duplicate DDL, endpoint schemas, event schemas, or implementation tasks. | Exact database shape belongs in `schema.sql`; refer to object names without copying DDL. |
| `system/data/db/indexes.md` | Index changes, query patterns, uniqueness rationale, cardinality assumptions, performance trade-offs, and rollback risks. | Table DDL ownership, API pagination contracts, dashboard layout, or broad data rationale. | The index supports a schema, flow, or module; link to the owner. |
| `system/data/db/triggers.md` | Trigger changes, timing, affected tables, invariants enforced, side effects, and failure behavior. | Application worker behavior, API validation, or non-database business flow. | The trigger protects a data invariant or affects a flow/module; link to the owner. |
| `system/data/db/functions.md` | Database function changes, inputs, outputs, volatility, permissions, invariants, side effects, and operational risk. | Application service logic, API schemas, or broad architecture decisions. | A function changes data invariants, module behavior, or contracts; link to the owner. |
| `system/security/model.md` | Trust boundary changes, attacker-controlled inputs, sensitive data handling, encryption expectations, secret handling, audit requirements, threat assumptions, and security invariants. | Secret values, endpoint payload copies, every permission rule, implementation checklists, or speculative vulnerabilities without accepted design impact. | Enforcement is in auth, contracts, modules, data, or flows; state the invariant here and link to enforcement owners. |
| `system/obs/metrics.md` | Metrics changes: names, labels, cardinality, counters, histograms, saturation signals, and why they exist. | Log fields, trace span structure, alert thresholds, dashboard layout, or code snippets. | A metric feeds an alert or dashboard; link to `alerts.md` or dashboard docs. |
| `system/obs/logging.md` | Logging changes: required fields, redaction, correlation IDs, retention, audit/debug purpose, and operator usage. | Metrics, traces, alert thresholds, raw sensitive examples, or incident history. | Logs support flows, security, or recovery; link to the owner. |
| `system/obs/tracing.md` | Trace boundary changes, span names, propagation, sampling, cross-service context, and latency attribution. | Metrics definitions, log retention, dashboard layout, or flow prose. | Traces connect modules, contracts, and flows; link to those owners. |
| `system/obs/alerts.md` | Alert changes: condition, severity, paging, escalation, runbook pointer, SLO/SLA impact, and recovery expectation. | Dashboard layout, raw metric implementation, incident narratives, or broad observability strategy. | An alert depends on metrics/logs/traces; link to the telemetry owner. |
| `system/obs/dashboards/grafana.md` | Dashboard changes: panels, variables, drilldowns, layout, and operator questions answered. | Alert policy, metric semantics, log schema, or screenshots as canonical truth. | A panel visualizes metrics, logs, traces, or alerts; link to the telemetry owner. |

Do not manually edit `.sysi/` operational state. Build-phase state, freeze baselines, allowlists, capture metadata, and agent-install metadata are machine-managed.

## Foundation Change Edge Cases

- If implementation wants a new endpoint, payload field, event, auth rule, error behavior, data shape, security invariant, or observability contract not in `/system`, treat that as design drift and use this skill before continuing.
- If the change is only a private refactor behind unchanged contracts, data, security, and observability, do not mutate `/system`.
- If one foundation change touches multiple owners, update all owners coherently after confirmation and record them in the decision artifact.
- If schema evolution is involved, record compatibility for old code, new code, old clients, new clients, existing data, migrations, rollbacks, and backfills.
- If a contract change is not backward compatible, record the versioning/deprecation plan or explicitly state that no existing consumers exist.
- If a security boundary changes, record the attacker-controlled inputs, authorization checks, sensitive data impact, audit requirement, and residual risk.
- If observability changes are required to operate the new truth, update the relevant `system/obs/` files in the same design-change.
- If workspace change artifacts conflict with the new foundation truth, update or pause those changes before resuming implementation.

## Workflow

1. State why normal apply cannot continue without changing foundation truth.
2. Use the routing table to identify exact owner files and cross-link-only files.
3. Update the decision artifact with rationale as the working record.
4. List affected `/system` files in the decision artifact.
5. List impacted workspace changes and implementation tasks in the decision artifact.
6. Describe migration or compatibility notes for already-built code, data, APIs, events, clients, security, and operations.
7. Ask for explicit user confirmation before editing controlled or frozen files.
8. Apply the smallest coherent `/system` mutation after confirmation.
9. Update the decision artifact with accepted decision and consequences.
10. Re-run validation before and after the mutation when possible.

## Validation

- Capture the before and after state of affected files in the conversation summary.
- Confirm the decision artifact records rationale, affected files, impacted workspace changes, compatibility notes, confirmation, decision, and consequences.
- Confirm every changed fact has exactly one owning file.
- Confirm duplicate-looking facts are replaced with cross-links.
- Confirm edited `/system` files do not conflict.
- Confirm security changes are captured in `system/security/model.md` when they affect trust boundaries, sensitive data, secrets, or security invariants.
- Confirm schema evolution keeps `system/data/schema.sql`, `system/data/schema.md`, and database detail files aligned.
- Confirm impacted workspace changes still describe the intended implementation.
- Confirm validation warnings are understood after the mutation.
- Confirm new foundation truth is precise enough for future agents.

## Stop Conditions

- Stop if the user has not given explicit user confirmation for controlled or frozen file edits.
- Stop if affected owner files cannot be identified from the routing table.
- Stop if impacted workspace changes are unknown and the mutation would change implementation scope.
- Stop if migration or compatibility notes are required but missing.
- Stop if the proposed change would make `/system` internally inconsistent.

## Do Not

- Do Not mutate foundation truth silently.
- Do Not use this skill for ordinary code changes.
- Do Not make broad rewrites when a smaller foundation correction is enough.
- Do Not edit implementation code before the foundation change is accepted.
- Do Not leave workspace changes inconsistent with changed `/system` truth.
- Do Not write secrets, credentials, tokens, private keys, or production sensitive values into `/system`.
