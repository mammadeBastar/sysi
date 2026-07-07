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
- Do not use build changes for design-phase capture.
- Do not turn implementation discoveries into foundation mutations without a finalized user decision.

## Role And File Access

- Respect the inferred role and its allowed `/system` files.
- Workspace agents read and update the /system files allowed for their declared workspace role, including `system/modules/<workspace>.md`.
- System-maintainer or design agents may update architecture and cross-cutting files when the user has finalized that foundation decision.
- If a needed file is outside the role allowlist, stop and ask before proceeding.

## System File Routing

Use the smallest owning file set. Do not duplicate full truth across files; explicitly avoid duplicated truth. Cross-link by relative path when another file owns the detail.

| Path | Must Own | Must Not Contain | Cross-Link Instead When |
| --- | --- | --- | --- |
| `system/architecture/system.md` | System shape, applications/services, major responsibilities, communication patterns, deployment assumptions, technology choices, and system-wide invariants. | Endpoint schemas, event message schemas, table DDL, UI flow steps, dashboard layouts, secret values, or detailed module internals. | A decision affects contracts, data, security, observability, or a specific flow/module; summarize the architectural implication and link to the owner. |
| `system/architecture/decisions/*.md` | Durable decision record: status, decision, rationale, affected files, consequences, and date. | The canonical copy of contracts, schemas, security policy, implementation tasks, or brainstorming history. | A decision spans multiple owners; list affected files and let those files own current truth. |
| `system/contracts/api.yaml` | HTTP paths, methods, parameters, request bodies, response bodies, status codes, and reusable OpenAPI components. | User journey prose, module internals, database DDL, security threat model, or shared error/convention prose. | Payloads need shared conventions, auth rules, errors, flows, or data rationale; link to `conventions.md`, `auth.md`, `errors.md`, flows, or data docs. |
| `system/contracts/events.asyncapi.yaml` | Event channels, operations, message schemas, producers, consumers, and reusable AsyncAPI components. | HTTP routes, batch job implementation plans, database schema, or operational runbooks. | Event semantics depend on shared conventions, errors, data, observability, or flows. |
| `system/contracts/auth.md` | Authentication, authorization, sessions, tokens, roles, permissions, and API/event boundary rules. | Secret values, encryption implementation details, endpoint payload schemas, or UI-only affordances. | The rule is a trust-boundary invariant or sensitive-data rule; link to `system/security/model.md`. The rule is attached to a route or event; link to the contract. |
| `system/contracts/conventions.md` | Cross-cutting API and event conventions: pagination, filtering, sorting, idempotency, correlation IDs, timestamps, versioning, deprecation, rate-limit expression, and compatibility rules. | Endpoint-specific fields, one-off event details, business process steps, or service internals. | A convention is used by a route/event; keep the generic rule here and reference it from the concrete contract. |
| `system/contracts/errors.md` | Error envelope, shared error codes, retryability, validation failure shape, user-facing versus internal error boundary, and rate-limit error semantics. | Every endpoint's full status-code table, stack traces, log formats, or UI copy. | A concrete route/event has specific error cases; define common behavior here and reference it from `api.yaml` or `events.asyncapi.yaml`. |
| `system/flows/*.flow.md` | One user or system action flow: actors, trigger, preconditions, data movement, service interactions, success path, failure branches, and externally visible outcome. | Canonical request/response schemas, event schemas, table definitions, component inventories, or low-level code steps. | The flow crosses an API/event boundary, uses data, relies on security, emits telemetry, or depends on modules. |
| `system/modules/<workspace>.md` | The declared workspace's pages, routes, services, components, workers, state ownership, responsibilities, dependencies, and boundary assumptions. | API schema copies, full flow narratives, other workspaces' module internals, table DDL, or operational alert rules. | A component or module consumes or exposes a contract, participates in a flow, owns data, emits telemetry, or enforces a security rule; reference the owning file. |
| `system/data/schema.sql` | Canonical Postgres DDL: tables, columns, constraints, types, primary keys, foreign keys, and schema-level indexes when they are part of the executable schema. | Prose rationale, API payload definitions, migration plans, sample data, or alternative schemas. | A human explanation, invariant, protobuf note, or relationship rationale is needed; put it in `schema.md`. |
| `system/data/schema.md` | Data relationships, invariants, lifecycle, retention, schema rationale, protobuf notes, and data-shape explanation that complements `schema.sql`. | Duplicate DDL, endpoint schemas, event schemas, or module implementation details. | The exact database shape belongs in `schema.sql`; refer to table and column names instead of restating the DDL. |
| `system/data/db/indexes.md` | Indexes, query patterns they support, cardinality assumptions, uniqueness rationale, and performance trade-offs. | Table ownership, full DDL copies, API pagination rules, or dashboard definitions. | The index depends on a table, flow, or module; link to `schema.sql`, `schema.md`, flows, or modules. |
| `system/data/db/triggers.md` | Database triggers, timing, affected tables, invariants enforced, side effects, and failure behavior. | Application worker logic, API validation, or business workflows outside the database. | The trigger exists to protect a data invariant; link to `schema.md` and the affected flow/module. |
| `system/data/db/functions.md` | Database functions, inputs, outputs, invariants, volatility, permissions, and operational risks. | Application service logic, API schemas, or broad architecture decisions. | Function behavior affects data invariants, modules, or contracts; link to the owner. |
| `system/security/model.md` | Trust boundaries, attacker-controlled inputs, sensitive data rules, encryption expectations, secret handling, audit expectations, threat assumptions, and security invariants. | Real secret values, endpoint payload schemas, every permission rule, vulnerability speculation without decisions, or implementation checklists. | A security invariant is enforced by auth, contracts, modules, data, or flows; state the invariant here and link to enforcement owners. |
| `system/obs/metrics.md` | Metrics names, labels, cardinality limits, success/failure counters, latency histograms, saturation signals, and why they exist. | Log field schemas, trace span trees, alert thresholds, dashboard layout, or implementation code. | A metric feeds an alert or dashboard; link to `alerts.md` or dashboard docs. |
| `system/obs/logging.md` | Logging strategy, required fields, redaction rules, correlation IDs, retention, and operator/debug use. | Metrics, trace spans, raw log examples with sensitive data, or alert policy. | Logs support a flow, security invariant, or incident response; reference the owner. |
| `system/obs/tracing.md` | Trace boundaries, span names, propagation, sampling, cross-service context, and latency attribution. | Metrics definitions, log retention, dashboard layout, or business flow prose. | Traces connect contracts, modules, and flows; link to those owners. |
| `system/obs/alerts.md` | Alert conditions, severity, paging rules, runbook pointers, SLO/SLA implications, and escalation expectations. | Dashboard layout, metric implementation details, or incident history. | An alert depends on a metric/log/trace; link to the telemetry owner. |
| `system/obs/dashboards/grafana.md` | Intended dashboard layout, panels, variables, drilldowns, and operator questions the dashboard answers. | Alert policy, metric semantics, log schema, or screenshots as source of truth. | A panel uses metrics, logs, traces, or alerts; link to the telemetry owner. |

Do not manually edit `.sysi/` operational state during capture. `.sysi/state.json`, `.sysi/freeze.json`, `.sysi/allowlists.json`, `.sysi/captures/`, and `.sysi/agents/` are machine-managed workflow state, not design truth.

## Routing Edge Cases

- If one decision changes multiple files, update each owning file once and add one Decision Record listing all affected files.
- If two files would own the same fact, choose the more canonical owner: contracts for boundaries, `schema.sql` for database shape, `security/model.md` for security invariants, observability files for telemetry expectations, flows for behavior over time, modules for component responsibility, and architecture for system-wide invariants.
- If a route/event payload includes a field stored in the database, define the payload in contracts and the persisted shape in data. Do not make either file a copy of the other.
- If a flow has errors, retries, idempotency, or auth behavior, summarize only the flow-specific branch and link to `contracts/errors.md`, `contracts/conventions.md`, or `contracts/auth.md`.
- If a security rule is both policy and enforcement, put the invariant in `security/model.md` and the enforcement location in contracts/modules/data/flows.
- If observability is required for a behavior, capture telemetry in `system/obs/` and reference it from the flow or module. Do not bury metrics or alerts inside flow prose.
- If the user gives implementation details that do not change foundation truth, do not capture them.
- If the decision conflicts with existing `/system` truth, stop, state the conflict, and ask the user to resolve it before editing.

## Workflow

1. Identify the Finalized Decision in one sentence.
2. Choose the smallest set of target `/system` files that should own the truth using the routing table.
3. Update owning files before dependents: architecture for invariants, contracts for boundaries, flows for behavior, modules for responsibility, data for persistence, security for trust and risk, and observability for operations.
4. Cross-link instead of copying full details across multiple files.
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
- Confirm each captured fact has exactly one owning file.
- Confirm every duplicate-looking fact is either removed or replaced with a cross-link.
- Confirm cross-boundary payloads are represented in `system/contracts/`.
- Confirm contract conventions are represented in `system/contracts/conventions.md` and error behavior is represented in `system/contracts/errors.md`.
- Confirm security posture changes are represented in `system/security/model.md`.
- Confirm data changes keep `system/data/schema.sql` and `system/data/schema.md` aligned without duplicating DDL.
- Confirm observability expectations are in the correct `system/obs/` file.
- Confirm the decision record points to the files that now own the truth.

## Stop Conditions

- Stop if the decision is not finalized.
- Stop if the project is in build phase.
- Stop if required target files are outside allowed role access.
- Stop if two `/system` files would become conflicting sources of truth.
- Stop if the correct owner file cannot be identified from the routing table and ask the user to choose the owner.

## Do Not

- Do Not capture brainstorming notes as final architecture.
- Do Not create build changes during design phase.
- Do Not edit implementation files.
- Do Not duplicate full contracts, schemas, or flow details in unrelated files.
- Do Not bypass `sysi design-change` in build phase.
- Do Not write secrets, credentials, tokens, private keys, or production sensitive values into `/system`.
