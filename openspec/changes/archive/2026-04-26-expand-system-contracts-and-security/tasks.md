## 1. CLI Foundation

- [x] 1.1 Add `system/contracts/conventions.md`, `system/contracts/errors.md`, and `system/security/model.md` to the system scaffold with concise starter content.
- [x] 1.2 Add `system/security/` to scaffolded directories.
- [x] 1.3 Add the new files to required `/system` validation.
- [x] 1.4 Add the new contract files and security model to freeze baseline handling as controlled foundation files.
- [x] 1.5 Update default role allowlists so frontend and backend workflows can read `system/security/**` alongside their existing `/system` context.
- [x] 1.6 Update `sysi explore` and `sysi capture` output so design guidance names contracts, contract conventions, contract errors, and security as possible target areas.

## 2. Agent Instruction Templates

- [x] 2.1 Update `sysi-explore` to include security in design exploration, role guidance, and validation checks.
- [x] 2.2 Update `sysi-capture` to route decisions about conventions, errors, and security posture to the new owning files.
- [x] 2.3 Update `sysi-apply` to include security context when frontend or backend implementation depends on security invariants.
- [x] 2.4 Update `sysi-design-change` to include `system/security/` when build-phase foundation changes affect security truth.
- [x] 2.5 Update Cursor and Claude instruction templates to identify contracts and security as `/system` foundation truth.

## 3. Documentation

- [x] 3.1 Update the README `/system` tree to include `contracts/conventions.md`, `contracts/errors.md`, and `security/model.md`.
- [x] 3.2 Update README file-purpose descriptions for contracts and security.
- [x] 3.3 Update freeze and validation documentation so the new controlled files are accurately described.

## 4. Tests And Validation

- [x] 4.1 Update scaffold tests to assert the new contract and security files are created.
- [x] 4.2 Update validation tests to cover missing `system/contracts/conventions.md`, `system/contracts/errors.md`, and `system/security/model.md`.
- [x] 4.3 Update freeze tests to verify build-phase mutations to the new controlled files require `sysi design-change`.
- [x] 4.4 Update agent-template tests to assert generated instructions mention security foundation and the new contract guidance.
- [x] 4.5 Run `gofmt`, `go test ./...`, and relevant OpenSpec status/validation commands.
