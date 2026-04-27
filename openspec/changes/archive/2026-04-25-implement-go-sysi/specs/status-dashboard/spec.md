## ADDED Requirements

### Requirement: Show Interactive Status Dashboard
The system SHALL provide `sysi status` as the primary terminal dashboard.

#### Scenario: Status in terminal
- **WHEN** a user runs `sysi status` in an initialized project
- **THEN** the system prints a structured dashboard showing phase, freeze state, validation health, OpenSpec summary, agent installs, recent decisions, and warnings

### Requirement: Provide Machine-Readable Status
The system SHALL provide JSON status output for agents and scripts.

#### Scenario: JSON status requested
- **WHEN** a user runs `sysi status --json`
- **THEN** the system prints valid JSON containing phase, root path, detected role, validation results, freeze state, and agent integration state

### Requirement: Watch Status
The system SHALL support a watch mode for repeated status refreshes.

#### Scenario: Watch mode
- **WHEN** a user runs `sysi status --watch`
- **THEN** the system refreshes the dashboard until interrupted

### Requirement: Report Uninitialized Project
The system SHALL clearly report when commands run outside a sysi-initialized repository.

#### Scenario: Status outside project
- **WHEN** a user runs `sysi status` outside a directory tree containing `.sysi/state.json`
- **THEN** the system reports that the project is not initialized and suggests running `sysi init`
