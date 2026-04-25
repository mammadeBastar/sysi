## MODIFIED Requirements

### Requirement: Initialize Sys Repository
The system SHALL provide `sys init` to initialize a repo-local sys-orchestrated monorepo and prepare frontend/backend OpenSpec workspaces.

#### Scenario: Initialize empty repository
- **WHEN** a user runs `sys init` at the repository root
- **THEN** the system creates `.sys-orchestrator/`, scaffolds `/system`, creates `frontend/` and `backend/` when missing, initializes OpenSpec inside `frontend/` and `backend/`, records the phase as `design`, and prints the next available commands

#### Scenario: Initialize frontend and backend only
- **WHEN** a user runs `sys init`
- **THEN** the system initializes OpenSpec under `frontend/` and `backend/`
- **AND** the system does not initialize OpenSpec under the repository root or under `/system`

#### Scenario: Initialize already initialized repository
- **WHEN** a user runs `sys init` in a repository that already contains `.sys-orchestrator/state.json`
- **THEN** the system preserves existing state, ensures `frontend/` and `backend/` OpenSpec workspaces exist when missing, and reports that the repository is already initialized
