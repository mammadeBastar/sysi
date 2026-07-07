## Purpose

Define the design-phase workflow where agents explore and capture finalized decisions without creating build changes.

## Requirements

### Requirement: Explore Without Build Changes
The system SHALL support design-phase exploration without creating or requiring build changes.

#### Scenario: Explore in design phase
- **WHEN** a user runs `sysi explore auth`
- **THEN** the system prints design-agent instructions based on current `/system` files, including relevant contracts and security foundation files, and does not create build changes

### Requirement: Capture Finalized Decisions Into System
The system SHALL support finalized design decisions being captured directly into `/system` during design phase.

#### Scenario: Capture command in design phase
- **WHEN** a user runs `sysi capture`
- **THEN** the system prints capture rules, target file guidance for architecture, contracts, flows, modules, data, observability, and security, and decision-record requirements for the active agent

#### Scenario: Codex capture skill completes capture
- **WHEN** the Codex `sysi-capture` skill is invoked after a finalized decision
- **THEN** the agent updates the appropriate `/system` files, including contract conventions, contract errors, or security model files when they own the decision, and writes a decision record

### Requirement: Block Normal Capture In Build Phase
The system SHALL block normal design capture after the project enters build phase.

#### Scenario: Capture attempted during build phase
- **WHEN** the project phase is `build` and a user runs `sysi capture`
- **THEN** the system reports that `sysi design-change` is required instead

### Requirement: Infer Agent Role From Current Directory
The system SHALL infer agent role from the current working directory instead of requiring role-specific user commands.

#### Scenario: Declared workspace directory role
- **WHEN** a user invokes a sysi agent workflow from inside a declared workspace `<name>`
- **THEN** the system treats the agent as the `<name>` workspace agent and applies that workspace's system-file allowlist

#### Scenario: Root and system directory roles
- **WHEN** a user invokes a sysi agent workflow from the repository root or from `system/`
- **THEN** the system infers the `design` role at the root and the `system-maintainer` role under `system/`
