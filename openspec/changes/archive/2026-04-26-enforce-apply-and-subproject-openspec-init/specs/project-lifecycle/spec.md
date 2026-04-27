## MODIFIED Requirements

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
