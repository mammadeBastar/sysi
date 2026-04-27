## Context

Sysi separates ratified foundation truth in `/system` from build-phase implementation work in OpenSpec. The current `sysi-apply` skill already tells agents to compare implementation needs against `/system` and to escalate when foundation truth must change, but the handoff is too implicit. In practice, a new endpoint, changed payload shape, new auth behavior, new error envelope, schema shift, or observability contract change can look like a "little implementation change" unless the apply instructions force the agent to classify it as design drift.

The change affects generated agent instructions more than CLI mechanics. The Go CLI already exposes `sysi design-change <name>` and the generated `sysi-design-change` skill already requires rationale, affected files, migration notes, validation, and explicit confirmation before controlled or frozen `/system` edits. The missing piece is the `sysi-apply` behavior at the moment drift is discovered.

## Goals / Non-Goals

**Goals:**

- Make `sysi-apply` classify implementation needs against `/system` before changing behavior.
- Make design drift examples concrete enough for agents to recognize common cases.
- Require the agent to double-check detected drift with the user before starting a foundation mutation.
- Require user-agreed drift to enter `sysi design-change <name>` and the generated `sysi-design-change` workflow before `/system` edits.
- Prevent implementation from continuing when the user does not agree to a required foundation change.

**Non-Goals:**

- Do not make the Go CLI automatically infer or apply `/system` mutations from source-code diffs.
- Do not bypass `sysi-design-change` confirmation rules.
- Do not make every minor implementation or UI polish change require a foundation update.
- Do not change OpenSpec archive behavior or implement a new sync engine.
- Do not add hard filesystem sandboxing.

## Decisions

### Treat Design Drift As A User-Confirmed Handoff

When `sysi-apply` detects that implementation would contradict or extend `/system`, it should stop ordinary apply work and present the mismatch to the user. The agent should name:

- what implementation needs
- what `/system` currently says or omits
- which `/system` files likely own the truth
- the proposed `sysi design-change <name>` path

Only after user agreement should the agent invoke `sysi design-change <name>` and follow `sysi-design-change` before editing controlled or frozen foundation files.

Alternative considered: have `sysi-apply` automatically run `sysi design-change` as soon as drift is detected. That is too eager because detection can be wrong, and foundation mutation intentionally requires human agreement.

### Make Drift Examples Concrete In Agent Instructions

The generated `sysi-apply` skill should explicitly call out examples:

- new or changed HTTP endpoints
- request or response payload-shape changes
- event contract changes
- auth/session/permission changes
- shared error or convention changes
- schema or data invariant changes
- security, metrics, logging, tracing, or alerting contract changes

These examples belong in the skill because they are runtime judgment aids for agents, not new CLI commands.

Alternative considered: document only in README. That helps humans, but the most important reader is the agent applying a change under pressure.

### Use Existing Design-Change Ceremony

The implementation should reuse the current `sysi-design-change` skill and `sysi design-change <name>` command instead of adding another command. That preserves one ceremony for controlled foundation mutation and keeps the behavior understandable:

```text
sysi-apply detects drift
        |
        v
agent asks user to confirm the foundation change
        |
        v
sysi design-change <name> + sysi-design-change
        |
        v
/system updated, validated, then apply resumes if still consistent
```

Alternative considered: add `sysi apply-drift` or `sysi change drift`. That fragments the model without adding real enforcement.

## Risks / Trade-offs

- [Risk] Agents over-classify small implementation details as design drift. Mitigation: document that ordinary refactors, UI polish, and implementation details behind unchanged contracts do not require design-change.
- [Risk] Agents under-classify API payload changes as implementation details. Mitigation: include concrete drift examples in `sysi-apply` and tests.
- [Risk] The user agrees to a foundation change but OpenSpec artifacts remain stale. Mitigation: `sysi-design-change` already requires impacted OpenSpec changes to stay consistent; reinforce that apply resumes only after artifacts and `/system` agree.
- [Risk] CLI behavior still cannot force every agent to comply. Mitigation: this is an instruction-pack enforcement improvement; hard runtime enforcement remains outside v1 boundaries.

## Migration Plan

1. Update `build-workflow` and `agent-integration` delta specs with the user-confirmed drift handoff requirement.
2. Update `sysi-apply` template wording in workflow, validation, stop conditions, and prohibited actions.
3. Update tests that assert generated instruction guardrails so the new drift-confirmation language cannot regress.
4. Update README build-phase guidance to explain that small implementation changes do not update `/system`, while contract/foundation drift requires user-agreed `sysi design-change`.
5. Run `GOCACHE=/tmp/sysi-go-cache go test ./...` and `openspec validate --specs`.

Rollback is a normal revert of the template, docs, specs, and tests if the stricter handoff proves too heavy.

## Open Questions

None. The intended behavior is deliberate confirmation, not automatic foundation mutation.
