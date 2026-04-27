## Purpose

Define the agent-runtime integrations for Codex, Cursor, and Claude Code.

## Requirements

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

### Requirement: Install Comprehensive Codex Skills
The system SHALL generate comprehensive Codex skills for `sysi-explore`, `sysi-capture`, `sysi-apply`, and `sysi-design-change`.

#### Scenario: Codex skill content is operational
- **WHEN** a user runs `sysi agent install codex`
- **THEN** each generated Codex skill contains purpose, phase rules, role and file access rules, workflow steps, validation expectations, stop conditions, and prohibited actions

### Requirement: Codex Explore Skill Guides Design Discovery
The system SHALL generate a `sysi-explore` Codex skill that guides design-phase exploration without OpenSpec design changes.

#### Scenario: Explore skill installed
- **WHEN** `sysi-explore/SKILL.md` is generated
- **THEN** it instructs agents to read status, infer role, read allowed `/system` files, explore architecture/contracts/flows/modules/data/security/observability as relevant, surface candidate decisions, avoid implementation, and suggest `sysi-capture` only after decisions are finalized

### Requirement: Codex Capture Skill Guides System Mutation
The system SHALL generate a `sysi-capture` Codex skill that guides finalized design decisions into `/system`.

#### Scenario: Capture skill installed
- **WHEN** `sysi-capture/SKILL.md` is generated
- **THEN** it defines what counts as a finalized decision, how to select target `/system` files including contract conventions, contract errors, and security model files, how to avoid duplicated truth, how to write decision records, how to validate after capture, and when to stop instead of mutating files

### Requirement: Codex Apply Skill Enforces Build Apply Boundaries
The system SHALL generate a `sysi-apply` Codex skill that coordinates OpenSpec apply with Sysi foundation rules and Superpowers discipline.

#### Scenario: Apply skill installed
- **WHEN** `sysi-apply/SKILL.md` is generated
- **THEN** it requires build phase, mandatory OpenSpec apply invocation, mandatory Superpowers apply discipline, `/system` context review including contracts, flows, modules, data, observability, and security where relevant, frozen-file protection, design drift detection, and escalation to `sysi-design-change` when foundational truth must change

#### Scenario: Apply skill lacks required external workflow
- **WHEN** the Codex `sysi-apply` skill is invoked and the local OpenSpec apply or Superpowers workflow is unavailable
- **THEN** it instructs the agent to stop and report the missing prerequisite instead of implementing without the required apply discipline

### Requirement: Codex Design Change Skill Defines Foundation Mutation Ceremony
The system SHALL generate a `sysi-design-change` Codex skill that defines controlled foundation mutation during build phase.

#### Scenario: Design change skill installed
- **WHEN** `sysi-design-change/SKILL.md` is generated
- **THEN** it requires rationale, affected `/system` files including security files when security truth changes, impacted OpenSpec changes, migration or compatibility notes, validation before and after mutation, and explicit user confirmation before updating controlled or frozen files

### Requirement: Agent Instructions Use Maintainable Templates
The system SHALL keep generated agent instruction content in maintainable template-backed assets or clearly separated template constants.

#### Scenario: Developer reviews instruction content
- **WHEN** a developer inspects the repository
- **THEN** they can read the complete generated instruction content without reconstructing it from many small inline fragments

### Requirement: Cursor Instructions Are Explicit But Minimal
The system SHALL generate Cursor rules that remain minimal while explicitly covering phase boundaries, `/system` authority, OpenSpec build workflow, design-change protection, and role inference.

#### Scenario: Cursor rules installed
- **WHEN** a user runs `sysi agent install cursor`
- **THEN** `.cursor/rules/sysi.mdc` contains explicit workflow rules and safety boundaries, references contracts and security as foundation truth, and does not claim deep runtime integration

### Requirement: Claude Instructions Are Explicit But Minimal
The system SHALL generate a Claude Code section that remains minimal while explicitly covering phase boundaries, `/system` authority, OpenSpec build workflow, design-change protection, and role inference.

#### Scenario: Claude instructions installed
- **WHEN** a user runs `sysi agent install claude`
- **THEN** the marked sysi section in `CLAUDE.md` contains explicit workflow rules and safety boundaries, references contracts and security as foundation truth, and does not claim deep runtime integration

### Requirement: Instruction Pack Tests Check Required Guardrails
The system SHALL test generated agent instructions for required operational sections and guardrail phrases.

#### Scenario: Agent instructions regress to skeletal content
- **WHEN** generated skills omit required sections such as phase rules, role access, validation, stop conditions, prohibited actions, or references to the security foundation
- **THEN** the test suite fails
