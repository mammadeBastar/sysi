## Context

`README.md` currently explains the project in a compact quickstart form. Since the first Go implementation has been completed and archived, the README now needs to serve as the primary onboarding document for humans and agents. It should explain the system philosophy, the lifecycle, the `/system` foundation, OpenSpec/Superpowers boundaries, commands, integrations, and v1 limitations without requiring readers to inspect code or archived proposals first.

## Goals / Non-Goals

**Goals:**
- Make `README.md` complete enough to understand and use the v1 CLI.
- Preserve the README as a practical reference, not a marketing page.
- Document both human workflows and agent workflows.
- Include command examples that match the current implementation.
- Explain the relationship between `sysi`, `/system`, OpenSpec, Superpowers, Codex, Cursor, and Claude Code.
- Give contributors enough context to update documentation when behavior changes.

**Non-Goals:**
- Changing CLI behavior.
- Adding generated docs tooling.
- Creating separate docs files.
- Documenting future features as if they already exist.
- Replacing OpenSpec specs as the normative change-tracking mechanism.

## Decisions

### Keep documentation in a single README for v1

The README will be treated as the full documentation surface for now. This keeps the project easy to navigate while the tool is still small.

Alternatives considered:
- Split docs into `docs/`: useful later, but premature for v1.
- Keep README as a short landing page: insufficient for agent-native workflow setup.

### Structure the README around user tasks

The README should prioritize what a user or agent needs to do:

1. Understand the mental model.
2. Install/run the CLI.
3. Initialize a repo.
4. Use design phase.
5. Freeze and use build phase.
6. Install agent integrations.
7. Validate and troubleshoot.
8. Understand command reference and file layout.

This is more useful than organizing the document by internal package structure.

### Explicitly document boundaries and limitations

The README should state that Codex is first-class in v1, Cursor/Claude support is minimal, `/system/views` are intentionally absent, and build-phase implementation should flow through OpenSpec plus Superpowers discipline.

Alternatives considered:
- Avoiding limitations to keep docs shorter: makes the tool easier to misuse.
- Describing future adapters in detail: risks documenting behavior that does not exist.

### Keep command examples copy-pasteable

Examples should use real command names from the implementation and avoid unsupported aliases. The README should mention `go run ./cmd/sysi <command>` for source usage and `sysi <command>` for installed usage.

## Risks / Trade-offs

- [Risk] README becomes too long to scan -> Mitigation: use a table of contents, concise sections, and a command reference.
- [Risk] Documentation drifts from implementation -> Mitigation: base sections on current code and specs; include a contributor note to update README when workflows change.
- [Risk] Readers confuse design phase and build phase -> Mitigation: make phase boundaries explicit and repeat the OpenSpec/Superpowers responsibilities in the workflow sections.
- [Risk] Cursor/Claude users expect full parity with Codex -> Mitigation: document them as minimal instruction scaffolds in v1.

## Migration Plan

Replace the current README content with the enriched documentation in one edit. No runtime migration is required.

## Open Questions

None for v1.
