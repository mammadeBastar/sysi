## 1. Init Lifecycle

- [x] 1.1 Add a helper that ensures `frontend/` and `backend/` directories exist during `sysi init`.
- [x] 1.2 Add a helper that runs non-interactive `openspec init <target> --tools none` from the sysi root for `frontend` and `backend`.
- [x] 1.3 Skip OpenSpec initialization for a target that already contains `openspec/config.yaml`.
- [x] 1.4 Preserve existing sysi state when `sysi init` runs in an already initialized repository while still ensuring frontend/backend OpenSpec workspaces.

## 2. Apply Dispatch

- [x] 2.1 Update `sysi change apply <name>` to verify build phase and change existence, then invoke `openspec instructions apply --change <name> --json`.
- [x] 2.2 Keep apply output clear that implementation must continue through OpenSpec apply plus Superpowers discipline.
- [x] 2.3 Ensure apply failures from the OpenSpec command are surfaced as command failures.
- [x] 2.4 Route `sysi change propose|apply|archive` through the inferred `frontend/` or `backend/` OpenSpec workspace.

## 3. Agent Instructions

- [x] 3.1 Harden `internal/sysiapp/templates/agents/codex/sysi-apply/SKILL.md` to require the local OpenSpec apply skill/workflow before implementation edits.
- [x] 3.2 Harden `sysi-apply` to require Superpowers methods for implementation planning, TDD, debugging, and verification.
- [x] 3.3 Add stop conditions for missing OpenSpec apply or Superpowers workflow.

## 4. Documentation

- [x] 4.1 Update `README.md` to document frontend/backend OpenSpec initialization during `sysi init`.
- [x] 4.2 Update `README.md` to document strict `sysi change apply` and `sysi-apply` behavior.

## 5. Tests And Validation

- [x] 5.1 Update init tests to use a fake OpenSpec executable and assert only `frontend` and `backend` receive OpenSpec initialization.
- [x] 5.2 Add or update apply tests to assert `sysi change apply` invokes the OpenSpec apply instruction command from the inferred implementation workspace.
- [x] 5.3 Update generated Codex skill guardrail tests for mandatory OpenSpec apply and Superpowers wording.
- [x] 5.4 Run `GOCACHE=/tmp/sysi-go-cache go test ./...`.
- [x] 5.5 Run `openspec validate enforce-apply-and-subproject-openspec-init --strict`.
- [x] 5.6 Run `openspec validate --specs`.
