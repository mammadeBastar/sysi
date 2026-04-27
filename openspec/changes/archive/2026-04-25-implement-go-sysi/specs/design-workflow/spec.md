## ADDED Requirements

### Requirement: Explore Without OpenSpec
The system SHALL support design-phase exploration without creating or requiring OpenSpec changes.

#### Scenario: Explore in design phase
- **WHEN** a user runs `sysi explore auth`
- **THEN** the system prints design-agent instructions based on current `/system` files and does not call `openspec`

### Requirement: Capture Finalized Decisions Into System
The system SHALL support finalized design decisions being captured directly into `/system` during design phase.

#### Scenario: Capture command in design phase
- **WHEN** a user runs `sysi capture`
- **THEN** the system prints capture rules, target file guidance, and decision-record requirements for the active agent

#### Scenario: Codex capture skill completes capture
- **WHEN** the Codex `sysi-capture` skill is invoked after a finalized decision
- **THEN** the agent updates the appropriate `/system` files and writes a decision record

### Requirement: Block Normal Capture In Build Phase
The system SHALL block normal design capture after the project enters build phase.

#### Scenario: Capture attempted during build phase
- **WHEN** the project phase is `build` and a user runs `sysi capture`
- **THEN** the system reports that `sysi design-change` is required instead

### Requirement: Infer Agent Role From Current Directory
The system SHALL infer agent role from the current working directory instead of requiring role-specific user commands.

#### Scenario: Frontend directory role
- **WHEN** a user invokes a sysi agent workflow from `frontend/`
- **THEN** the system treats the agent as a frontend agent and applies the frontend system-file allowlist
