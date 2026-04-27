## Context

The current `/system` scaffold intentionally keeps design truth outside OpenSpec during design phase. It already has good top-level buckets for architecture, contracts, flows, modules, data, and observability, but the default contract set only covers HTTP APIs, events, and auth. That leaves common cross-boundary behavior such as error envelopes, pagination, idempotency, versioning, and correlation IDs without a clear home.

The scaffold also lacks a security foundation. `system/contracts/auth.md` covers authentication and authorization boundaries, but it is not the right place for broader trust boundaries, sensitive-data rules, encryption expectations, secret handling, or security invariants.

## Goals / Non-Goals

**Goals:**

- Add minimal contract files for cross-cutting API/event conventions and error behavior.
- Add a first-class `system/security/model.md` file without introducing a broader runtime folder.
- Treat the new files as canonical, validated `/system` foundation files.
- Include the new files in design freeze baselines so build-phase mutations are visible and directed through `sysi design-change`.
- Update role allowlists and generated agent instructions so agents know when to read or update contract and security truth.
- Update tests and README documentation for the expanded scaffold.

**Non-Goals:**

- Adding feature-specific contract files such as uploads, billing, notifications, search, or webhooks.
- Adding a `system/runtime/` folder.
- Adding security scanning, policy enforcement, secret management, or runtime authorization logic.
- Changing the OpenSpec build/apply workflow.
- Changing generated implementation workspaces.

## Decisions

### Keep Contracts Minimal But More Complete

Add exactly two default contract files:

```text
system/contracts/errors.md
system/contracts/conventions.md
```

`errors.md` owns error envelopes, error codes, retryability, validation error shape, and the distinction between internal and user-facing errors. `conventions.md` owns cross-cutting API/event conventions such as pagination, filtering, sorting, idempotency, correlation IDs, timestamp formats, versioning, deprecation, and rate-limit expression.

Alternatives considered:

- Put these concerns in `api.yaml`: OpenAPI can describe concrete endpoints, but it is a poor long-form home for rationale and conventions that apply across many paths.
- Add many feature-specific contract files by default: this would make the scaffold feel bureaucratic before the product shape is known.

### Add Security As Its Own Foundation Area

Add:

```text
system/security/model.md
```

The file owns trust boundaries, sensitive data classification, encryption expectations, secret handling, security invariants, and threat assumptions. It deliberately stays separate from `system/contracts/auth.md`; auth remains a boundary contract, while security is a cross-cutting model.

Alternatives considered:

- Expand `auth.md`: this would overload identity/permission boundaries with unrelated security posture.
- Add multiple security files by default: useful later, but too much for the clean minimal scaffold.

### Treat New Files As Required And Controlled

The CLI should scaffold the new files during `sysi init`, include them in `requiredSystemFiles`, and include them in freeze baselines. The new contract files and security model should be controlled after design freeze, like other cross-boundary foundation files.

This keeps the design/build boundary coherent: changes to error behavior, conventions, or security invariants after freeze should be visible as foundation drift and should use `sysi design-change`.

### Include Security In Role Guidance

Generated skills and CLI allowlists should expose `system/security/**` to frontend and backend agents. The security model is not a secret store; it is durable design truth that both sides need when building authentication screens, API callers, backend handlers, logging, and data-handling behavior.

The agent guidance should still keep ownership narrow:

- Contract behavior belongs in `system/contracts/`.
- Security invariants belong in `system/security/model.md`.
- Auth/session/permission boundaries stay in `system/contracts/auth.md`.
- Feature-specific files remain opt-in, not scaffolded by default.

## Risks / Trade-offs

- [Risk] The scaffold starts growing too many default files -> Mitigation: add only two contract files and one security file; keep feature-specific concerns opt-in.
- [Risk] Agents duplicate security details into auth or architecture files -> Mitigation: skill guidance should identify `system/security/model.md` as the owning file for security posture and use links/summaries elsewhere.
- [Risk] Freeze baselines become too strict for evolving design details -> Mitigation: these files represent cross-boundary foundation truth; build-phase changes should be explicit through `sysi design-change`.
- [Risk] Role allowlists expose security text too broadly -> Mitigation: `system/security/model.md` must document rules and assumptions, never real secrets or secret values.

## Migration Plan

1. Update scaffold creation to add `system/security/`, `system/contracts/errors.md`, `system/contracts/conventions.md`, and `system/security/model.md`.
2. Add the new files to required-file validation and freeze baseline handling.
3. Update role allowlists so relevant agents can read the security model.
4. Update generated Codex, Cursor, and Claude templates to mention the new files in exploration, capture, apply, and design-change guidance.
5. Update README structure documentation.
6. Update tests for scaffold creation, validation, freeze behavior, generated instruction content, and documentation markers.

Existing initialized projects can receive these files by rerunning the updated `sysi init`; the existing scaffold behavior writes missing files without overwriting existing files.

## Open Questions

None. Feature-specific contract files remain intentionally out of scope until a concrete system needs them.
