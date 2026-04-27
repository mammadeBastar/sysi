## Context

The current agent integrations are installed by `sysi agent install codex|cursor|claude`. The command surface is right, but the generated content is too thin:

- Codex skills are short embedded strings in `internal/sysiapp/app.go`.
- Cursor rules are a small bullet list.
- Claude Code receives a small marked section.

For an agent-native workflow, those instructions need to carry the operating procedure: phase boundaries, `/system` authority, role-scoped reading, capture rules, design-change rules, OpenSpec/Superpowers boundaries, validation expectations, and stop conditions.

## Goals / Non-Goals

**Goals:**
- Make Codex skills comprehensive enough to guide real agent work.
- Keep one-command installation: `sysi agent install codex`.
- Keep Cursor and Claude Code support minimal, but materially more explicit.
- Make instruction content easy to inspect and maintain.
- Test for required sections and guardrails in generated instruction files.
- Preserve existing behavior around role inference and phase detection.

**Non-Goals:**
- Adding new user-facing install commands.
- Building deep Cursor or Claude Code runtime integrations.
- Enforcing hard filesystem sandboxing.
- Replacing OpenSpec apply skills or Superpowers skills.
- Changing the core CLI phase model.

## Decisions

### Use Template-Backed Instruction Packs

Move large instruction content out of tiny ad hoc string functions and into template-backed assets.

Preferred structure:

```text
internal/sysiapp/templates/agents/
  codex/
    sysi-explore/SKILL.md
    sysi-capture/SKILL.md
    sysi-apply/SKILL.md
    sysi-design-change/SKILL.md
  cursor/
    sysi.mdc
  claude/
    CLAUDE.section.md
```

These can be embedded with Go `embed` or represented as clearly separated constants if embedding adds unnecessary complexity. The important property is maintainability: each instruction pack should read like a real artifact, not a small helper string.

Alternatives considered:
- Keep hardcoded tiny strings: simple, but does not meet agent-native needs.
- Generate instructions dynamically from many fragments: powerful, but premature for v1 and easier to make incoherent.

### Make Codex The Rich Instruction Surface

Codex should receive comprehensive skills because Codex skills are the strongest supported integration surface in v1.

Each Codex skill should include:

- purpose and when to use it
- required initial checks
- role and file access rules
- phase rules
- step-by-step workflow
- expected outputs/artifacts
- validation expectations
- stop conditions
- things the agent must not do

Cursor and Claude Code should remain lighter but should still include explicit workflow boundaries and safety rules.

### Do Not Duplicate OpenSpec And Superpowers Internals

`sysi-apply` should require OpenSpec apply and Superpowers discipline, but it should not paste full copies of those external skills. Instead it should state when to invoke/read them and what local sysi-specific constraints apply before and during apply.

This keeps Sysi as an orchestrator rather than a fork of those workflows.

### Test Instruction Semantics Through Required Markers

Tests should assert generated files contain required sections and phrases such as:

- `Phase Rules`
- `Role And File Access`
- `Stop Conditions`
- `Decision Record`
- `OpenSpec`
- `Superpowers`
- `sysi design-change`
- `Do Not`

This avoids brittle full-file snapshots while preventing accidental regression back to skeletal instructions.

## Risks / Trade-offs

- [Risk] Skills become too long and hard to follow -> Mitigation: structure with clear headings and operational checklists.
- [Risk] Instructions drift from CLI behavior -> Mitigation: tests check required command names and README documents the instruction-pack model.
- [Risk] Cursor/Claude support feels second-class -> Mitigation: document that this is intentional in v1 while still making their instructions explicit enough to avoid dangerous behavior.
- [Risk] Sysi duplicates OpenSpec/Superpowers guidance -> Mitigation: reference those workflows instead of copying them wholesale.

## Migration Plan

Existing users can rerun:

```bash
sysi agent install codex
sysi agent install cursor
sysi agent install claude
```

This overwrites generated sysi-owned Codex skill files and Cursor rule files, and updates only the marked sysi section in `CLAUDE.md`.

## Open Questions

None for this change. The intended boundary is comprehensive Codex skills and enriched minimal Cursor/Claude instructions without new commands.
