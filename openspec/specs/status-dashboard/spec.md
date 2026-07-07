## Purpose

Define the human and machine-readable status dashboard behavior for `sysi`.

## Requirements

### Requirement: Show Interactive Status Dashboard
The system SHALL provide `sysi status` as the primary terminal dashboard.

#### Scenario: Status in terminal
- **WHEN** a user runs `sysi status` in an initialized project
- **THEN** the system prints a structured dashboard showing root path, phase, inferred role, validation health, freeze baselines, a per-workspace summary of native changes with their statuses, agent installs, and warnings
- **AND** the workspace summary covers exactly the declared workspaces

### Requirement: Provide Machine-Readable Status
The system SHALL provide JSON status output for agents and scripts.

#### Scenario: JSON status requested
- **WHEN** a user runs `sysi status --json`
- **THEN** the system prints valid JSON containing `root`, `phase`, `role`, `validation`, `freeze`, `agents`, and a `workspaces` array
- **AND** each `workspaces` entry contains `name`, `present`, `activeChanges`, and a `changes` list of `{name, status}` objects

#### Scenario: Empty lists are never null
- **WHEN** a project has no declared workspaces, no active changes in a workspace, no validation warnings, or no freeze mutations
- **THEN** the corresponding JSON fields are empty arrays (`[]`), never `null`

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
