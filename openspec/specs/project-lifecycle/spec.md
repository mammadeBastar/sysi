## Purpose

Define how `sysi` initializes and tracks a repo-local orchestrated monorepo lifecycle.

## Requirements

### Requirement: Initialize Sysi Repository
The system SHALL provide `sysi init` to initialize a repo-local sysi-orchestrated monorepo and prepare frontend/backend OpenSpec workspaces.

#### Scenario: Initialize empty repository
- **WHEN** a user runs `sysi init` at the repository root
- **THEN** the system creates `.sysi/`, scaffolds `/system`, creates `frontend/` and `backend/` when missing, initializes OpenSpec inside `frontend/` and `backend/`, records the phase as `design`, and prints the next available commands

#### Scenario: Initialize frontend and backend only
- **WHEN** a user runs `sysi init`
- **THEN** the system initializes OpenSpec under `frontend/` and `backend/`
- **AND** the system does not initialize OpenSpec under the repository root or under `/system`

#### Scenario: Initialize already initialized repository
- **WHEN** a user runs `sysi init` in a repository that already contains `.sysi/state.json`
- **THEN** the system preserves existing state, ensures `frontend/` and `backend/` OpenSpec workspaces exist when missing, and reports that the repository is already initialized

### Requirement: Discover Repository Root
The system SHALL discover the sysi repository root from the current working directory.

#### Scenario: Command runs from nested frontend directory
- **WHEN** a user runs a `sysi` command from `frontend/`
- **THEN** the system locates the nearest ancestor containing `.sysi/state.json` and uses that path as the project root

### Requirement: Track Project Phase
The system SHALL track the current project phase as repo-local state.

#### Scenario: Default design phase
- **WHEN** a repository is initialized
- **THEN** the current phase is `design`

#### Scenario: Freeze enters build phase
- **WHEN** a user runs `sysi design freeze`
- **THEN** the system records the current phase as `build` and records freeze baselines for controlled system files

### Requirement: Avoid Required Global State
The system SHALL operate without requiring global machine state.

#### Scenario: Repository moved to another machine
- **WHEN** the repository is copied with `.sysi/` and `/system`
- **THEN** `sysi status` reports project state without needing a prior global registration command
