## MODIFIED Requirements

### Requirement: Explore Without OpenSpec
The system SHALL support design-phase exploration without creating or requiring OpenSpec changes.

#### Scenario: Explore in design phase
- **WHEN** a user runs `sysi explore auth`
- **THEN** the system prints design-agent instructions based on current `/system` files, including relevant contracts and security foundation files, and does not call `openspec`

### Requirement: Capture Finalized Decisions Into System
The system SHALL support finalized design decisions being captured directly into `/system` during design phase.

#### Scenario: Capture command in design phase
- **WHEN** a user runs `sysi capture`
- **THEN** the system prints capture rules, target file guidance for architecture, contracts, flows, modules, data, observability, and security, and decision-record requirements for the active agent

#### Scenario: Codex capture skill completes capture
- **WHEN** the Codex `sysi-capture` skill is invoked after a finalized decision
- **THEN** the agent updates the appropriate `/system` files, including contract conventions, contract errors, or security model files when they own the decision, and writes a decision record
