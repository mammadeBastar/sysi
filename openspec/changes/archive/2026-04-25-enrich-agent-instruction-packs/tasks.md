## 1. Template Structure

- [x] 1.1 Create a maintainable instruction-pack layout for generated agent files.
- [x] 1.2 Move Codex skill content into template-backed assets or clearly separated template constants.
- [x] 1.3 Move Cursor and Claude instruction content into template-backed assets or clearly separated template constants.
- [x] 1.4 Keep `sysi agent install codex|cursor|claude` command behavior unchanged while sourcing the richer content.

## 2. Codex Instruction Packs

- [x] 2.1 Enrich `sysi-explore` with purpose, initial checks, phase rules, role and file access, exploration workflow, decision surfacing, validation, stop conditions, and prohibited actions.
- [x] 2.2 Enrich `sysi-capture` with finalized-decision criteria, target `/system` file selection, decision record guidance, duplicate-truth prevention, validation, stop conditions, and prohibited actions.
- [x] 2.3 Enrich `sysi-apply` with build-phase requirements, OpenSpec apply usage, Superpowers apply discipline, `/system` context review, frozen-file protection, design drift detection, validation, and `sysi design-change` escalation.
- [x] 2.4 Enrich `sysi-design-change` with rationale, affected `/system` files, impacted OpenSpec changes, migration or compatibility notes, before/after validation, explicit confirmation rules, stop conditions, and prohibited actions.

## 3. Cursor And Claude Instructions

- [x] 3.1 Enrich Cursor rules with explicit phase boundaries, `/system` authority, OpenSpec build workflow, design-change protection, role inference, and minimal integration limits.
- [x] 3.2 Enrich the Claude Code marked section with explicit phase boundaries, `/system` authority, OpenSpec build workflow, design-change protection, role inference, and minimal integration limits.
- [x] 3.3 Preserve existing Claude file content outside the managed sysi section when reinstalling Claude instructions.

## 4. Tests And Documentation

- [x] 4.1 Add tests that generated Codex skills contain required operational sections and guardrail phrases.
- [x] 4.2 Add tests that generated Cursor and Claude instructions contain required workflow boundaries.
- [x] 4.3 Add or update tests that command names and installation paths remain stable.
- [x] 4.4 Update README documentation to describe the richer agent instruction packs and their intended v1 boundaries.

## 5. Verification

- [x] 5.1 Run the Go test suite with a writable Go cache.
- [x] 5.2 Run OpenSpec validation for `enrich-agent-instruction-packs`.
- [x] 5.3 Confirm `openspec status --change "enrich-agent-instruction-packs"` reports the change as ready for implementation.
