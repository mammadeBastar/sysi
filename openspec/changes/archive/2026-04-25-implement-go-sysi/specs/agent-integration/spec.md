## ADDED Requirements

### Requirement: Install Codex Skills
The system SHALL install Codex-native sysi skills with `sysi agent install codex`.

#### Scenario: Install Codex integration
- **WHEN** a user runs `sysi agent install codex`
- **THEN** the system creates project-local `sysi-explore`, `sysi-capture`, `sysi-apply`, and `sysi-design-change` skill directories under `.codex/skills/`

### Requirement: Avoid Role-Specific Install Commands
The system SHALL avoid requiring users to install separate design, frontend, or backend agent roles.

#### Scenario: Codex integration installed once
- **WHEN** a user has run `sysi agent install codex`
- **THEN** the installed skills infer role from current working directory and project phase

### Requirement: Generate Minimal Cursor Instructions
The system SHALL generate minimal Cursor rules with `sysi agent install cursor`.

#### Scenario: Install Cursor integration
- **WHEN** a user runs `sysi agent install cursor`
- **THEN** the system writes `.cursor/rules/sysi.mdc` with sysi workflow rules and file-access guidance

### Requirement: Generate Minimal Claude Code Instructions
The system SHALL generate minimal Claude Code instructions with `sysi agent install claude`.

#### Scenario: Install Claude Code integration
- **WHEN** a user runs `sysi agent install claude`
- **THEN** the system creates or updates a clearly marked sysi section in `CLAUDE.md`

### Requirement: Preserve Existing Agent Files
The system SHALL avoid overwriting unrelated user-authored agent instructions.

#### Scenario: Existing CLAUDE file
- **WHEN** `CLAUDE.md` already contains user-authored content
- **THEN** `sysi agent install claude` updates only the marked sysi section or appends one if missing
