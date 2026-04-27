## 1. Specification And Documentation

- [x] 1.1 Update main specs from the change deltas during archive or sync review.
- [x] 1.2 Update README build-phase guidance to describe user-confirmed design drift handoff.

## 2. Agent Instruction Changes

- [x] 2.1 Update `sysi-apply` template initial checks and workflow to classify concrete design drift examples against `/system`.
- [x] 2.2 Update `sysi-apply` template to require user confirmation before invoking `sysi design-change <name>` and `sysi-design-change`.
- [x] 2.3 Update `sysi-apply` validation, stop conditions, and prohibited actions so implementation cannot continue when required drift is not accepted by the user.

## 3. Tests And Verification

- [x] 3.1 Extend generated instruction-pack tests to assert user-confirmed drift and `sysi design-change` handoff wording.
- [x] 3.2 Run `GOCACHE=/tmp/sysi-go-cache go test ./...`.
- [x] 3.3 Run `openspec validate --specs`.
