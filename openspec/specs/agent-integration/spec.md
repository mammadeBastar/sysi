## Purpose

Define the agent-runtime integrations for Codex, Cursor, and Claude Code.

## Requirements

### Requirement: Install Codex Skills
The system SHALL install Codex-native sys skills with `sys agent install codex`.

#### Scenario: Install Codex integration
- **WHEN** a user runs `sys agent install codex`
- **THEN** the system creates project-local `sys-explore`, `sys-capture`, `sys-apply`, and `sys-design-change` skill directories under `.codex/skills/`

### Requirement: Avoid Role-Specific Install Commands
The system SHALL avoid requiring users to install separate design, frontend, or backend agent roles.

#### Scenario: Codex integration installed once
- **WHEN** a user has run `sys agent install codex`
- **THEN** the installed skills infer role from current working directory and project phase

### Requirement: Generate Minimal Cursor Instructions
The system SHALL generate minimal Cursor rules with `sys agent install cursor`.

#### Scenario: Install Cursor integration
- **WHEN** a user runs `sys agent install cursor`
- **THEN** the system writes `.cursor/rules/sys-orchestrator.mdc` with sys workflow rules and file-access guidance

### Requirement: Generate Minimal Claude Code Instructions
The system SHALL generate minimal Claude Code instructions with `sys agent install claude`.

#### Scenario: Install Claude Code integration
- **WHEN** a user runs `sys agent install claude`
- **THEN** the system creates or updates a clearly marked sys-orchestrator section in `CLAUDE.md`

### Requirement: Preserve Existing Agent Files
The system SHALL avoid overwriting unrelated user-authored agent instructions.

#### Scenario: Existing CLAUDE file
- **WHEN** `CLAUDE.md` already contains user-authored content
- **THEN** `sys agent install claude` updates only the marked sys-orchestrator section or appends one if missing

### Requirement: Install Comprehensive Codex Skills
The system SHALL generate comprehensive Codex skills for `sys-explore`, `sys-capture`, `sys-apply`, and `sys-design-change`.

#### Scenario: Codex skill content is operational
- **WHEN** a user runs `sys agent install codex`
- **THEN** each generated Codex skill contains purpose, phase rules, role and file access rules, workflow steps, validation expectations, stop conditions, and prohibited actions

### Requirement: Codex Explore Skill Guides Design Discovery
The system SHALL generate a `sys-explore` Codex skill that guides design-phase exploration without OpenSpec design changes.

#### Scenario: Explore skill installed
- **WHEN** `sys-explore/SKILL.md` is generated
- **THEN** it instructs agents to read status, infer role, read allowed `/system` files, explore architecture/contracts/flows/modules/data/observability as relevant, surface candidate decisions, avoid implementation, and suggest `sys-capture` only after decisions are finalized

### Requirement: Codex Capture Skill Guides System Mutation
The system SHALL generate a `sys-capture` Codex skill that guides finalized design decisions into `/system`.

#### Scenario: Capture skill installed
- **WHEN** `sys-capture/SKILL.md` is generated
- **THEN** it defines what counts as a finalized decision, how to select target `/system` files, how to avoid duplicated truth, how to write decision records, how to validate after capture, and when to stop instead of mutating files

### Requirement: Codex Apply Skill Enforces Build Apply Boundaries
The system SHALL generate a `sys-apply` Codex skill that coordinates OpenSpec apply with Sys foundation rules and Superpowers discipline.

#### Scenario: Apply skill installed
- **WHEN** `sys-apply/SKILL.md` is generated
- **THEN** it requires build phase, mandatory OpenSpec apply invocation, mandatory Superpowers apply discipline, `/system` context review, frozen-file protection, design drift detection, and escalation to `sys-design-change` when foundational truth must change

#### Scenario: Apply skill lacks required external workflow
- **WHEN** the Codex `sys-apply` skill is invoked and the local OpenSpec apply or Superpowers workflow is unavailable
- **THEN** it instructs the agent to stop and report the missing prerequisite instead of implementing without the required apply discipline

### Requirement: Codex Design Change Skill Defines Foundation Mutation Ceremony
The system SHALL generate a `sys-design-change` Codex skill that defines controlled foundation mutation during build phase.

#### Scenario: Design change skill installed
- **WHEN** `sys-design-change/SKILL.md` is generated
- **THEN** it requires rationale, affected `/system` files, impacted OpenSpec changes, migration or compatibility notes, validation before and after mutation, and explicit user confirmation before updating controlled or frozen files

### Requirement: Agent Instructions Use Maintainable Templates
The system SHALL keep generated agent instruction content in maintainable template-backed assets or clearly separated template constants.

#### Scenario: Developer reviews instruction content
- **WHEN** a developer inspects the repository
- **THEN** they can read the complete generated instruction content without reconstructing it from many small inline fragments

### Requirement: Cursor Instructions Are Explicit But Minimal
The system SHALL generate Cursor rules that remain minimal while explicitly covering phase boundaries, `/system` authority, OpenSpec build workflow, design-change protection, and role inference.

#### Scenario: Cursor rules installed
- **WHEN** a user runs `sys agent install cursor`
- **THEN** `.cursor/rules/sys-orchestrator.mdc` contains explicit workflow rules and safety boundaries, without claiming deep runtime integration

### Requirement: Claude Instructions Are Explicit But Minimal
The system SHALL generate a Claude Code section that remains minimal while explicitly covering phase boundaries, `/system` authority, OpenSpec build workflow, design-change protection, and role inference.

#### Scenario: Claude instructions installed
- **WHEN** a user runs `sys agent install claude`
- **THEN** the marked sys-orchestrator section in `CLAUDE.md` contains explicit workflow rules and safety boundaries, without claiming deep runtime integration

### Requirement: Instruction Pack Tests Check Required Guardrails
The system SHALL test generated agent instructions for required operational sections and guardrail phrases.

#### Scenario: Agent instructions regress to skeletal content
- **WHEN** generated skills omit required sections such as phase rules, role access, validation, stop conditions, or prohibited actions
- **THEN** the test suite fails
