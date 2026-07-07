## Purpose

Define how `sysi` initializes and tracks a repo-local orchestrated monorepo lifecycle with declared workspaces.

## Requirements

### Requirement: Initialize Sysi Repository With Declared Workspaces
The system SHALL provide `sysi init --workspaces <name>[,<name>...]` to initialize a repo-local sysi-orchestrated monorepo with explicitly declared workspaces, without requiring any external tool.

#### Scenario: Initialize repository with workspaces
- **WHEN** a user runs `sysi init --workspaces frontend,backend` at the repository root
- **THEN** the system creates `.sysi/`, scaffolds `/system` including `system/modules/<workspace>.md` per declared workspace, creates each declared workspace directory with a `changes/` scaffold when missing, records the declared workspaces in state, records the phase as `design`, and prints the next available command

#### Scenario: Bare init prints guidance and fails
- **WHEN** a user runs `sysi init` without `--workspaces` in an uninitialized repository
- **THEN** the system prints usage guidance with examples, including a suggested default of `--workspaces frontend,backend`
- **AND** the system fails without writing any project state or scaffold

#### Scenario: Declared workspace conflicts with existing file
- **WHEN** a user runs `sysi init --workspaces <name>` and `<name>` already exists at the repository root as a regular file
- **THEN** the system reports the conflict and fails before writing any project state

#### Scenario: Invalid or duplicate workspace names rejected at init
- **WHEN** a user runs `sysi init --workspaces` with a name that is empty, reserved, not a lowercase slug, or duplicated in the list
- **THEN** the system reports the invalid name and fails before writing any project state

#### Scenario: Initialize already initialized repository
- **WHEN** a user runs `sysi init` in a repository that already contains `.sysi/state.json`
- **THEN** the system preserves existing state, re-ensures the `/system` scaffold, allowlists, and declared workspace `changes/` directories, and reports that the repository is already initialized

### Requirement: Manage Declared Workspaces
The system SHALL provide `sysi workspace list|add|remove` to manage the declared workspace list after initialization.

#### Scenario: List declared workspaces
- **WHEN** a user runs `sysi workspace list`
- **THEN** the system prints each declared workspace with its active change count

#### Scenario: Add a workspace
- **WHEN** a user runs `sysi workspace add worker`
- **THEN** the system validates the name as a lowercase slug that is not reserved and not already declared
- **AND** the system rejects the name if it collides with an existing regular file at the repository root
- **AND** the system creates `worker/changes/`, scaffolds `system/modules/worker.md`, updates the role allowlists, and records the workspace in state

#### Scenario: Remove a workspace with active changes
- **WHEN** a user runs `sysi workspace remove worker` and `worker/changes/` contains active changes
- **THEN** the system refuses, names the active changes, and requires `--force` to remove anyway

#### Scenario: Remove a workspace
- **WHEN** a user runs `sysi workspace remove worker` (or adds `--force` when active changes exist)
- **THEN** the system removes the workspace from state
- **AND** the system leaves the workspace directory and its contents on disk

### Requirement: Discover Repository Root
The system SHALL discover the sysi repository root from the current working directory.

#### Scenario: Command runs from nested workspace directory
- **WHEN** a user runs a `sysi` command from inside a declared workspace directory
- **THEN** the system locates the nearest ancestor containing `.sysi/state.json` and uses that path as the project root

### Requirement: Track Project Phase
The system SHALL track the current project phase as repo-local state.

#### Scenario: Default design phase
- **WHEN** a repository is initialized
- **THEN** the current phase is `design`

#### Scenario: Freeze enters build phase
- **WHEN** a user runs `sysi design freeze`
- **THEN** the system records the current phase as `build` and records freeze baselines for controlled system files

### Requirement: Use Version 2 State With Workspaces
The system SHALL store project state in `.sysi/state.json` with schema `version: 2` including the declared `workspaces` list, and SHALL reject state written by other schema versions.

#### Scenario: State stores declared workspaces
- **WHEN** a repository is initialized with `sysi init --workspaces frontend,backend`
- **THEN** `.sysi/state.json` records `version: 2` and `workspaces: ["frontend", "backend"]`

#### Scenario: V1 state rejected
- **WHEN** a `sysi` command loads a `.sysi/state.json` whose version is not `2`
- **THEN** the system fails with an error explaining that sysi v2 requires state version 2 and that v1 projects should keep using the v1 binary

#### Scenario: Invalid workspace names in state rejected
- **WHEN** a `sysi` command loads a state file whose `workspaces` list contains an invalid workspace name
- **THEN** the system fails with an error identifying the invalid state

### Requirement: Avoid Required Global State
The system SHALL operate without requiring global machine state.

#### Scenario: Repository moved to another machine
- **WHEN** the repository is copied with `.sysi/` and `/system`
- **THEN** `sysi status` reports project state without needing a prior global registration command
