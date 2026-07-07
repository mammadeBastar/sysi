## Purpose

Define the native build-phase change workflow that coordinates workspace changes and Superpowers apply discipline.

## Requirements

### Requirement: Manage Native Build Changes
The system SHALL provide build-phase commands `sysi change propose|apply|archive <name>` that operate natively on `<workspace>/changes/<name>/` directories without invoking any external tool.

#### Scenario: Propose build change
- **WHEN** a user runs `sysi change propose add-login` from inside the declared workspace `frontend/`
- **THEN** the system verifies the project is in build phase
- **AND** the system creates `frontend/changes/add-login/` containing templated `proposal.md`, `design.md`, and `tasks.md` plus a `meta.json` with status `proposed`

#### Scenario: Build change outside a declared workspace
- **WHEN** a user runs a `sysi change` command from the monorepo root, `system/`, or any directory not inside a declared workspace
- **THEN** the system reports that change commands must run inside a declared workspace and lists the declared workspaces

#### Scenario: Change name must be a lowercase slug
- **WHEN** a user runs a `sysi change` command with a name that is not a lowercase slug
- **THEN** the system rejects the name and suggests the slugified form

#### Scenario: Duplicate or colliding change name rejected
- **WHEN** a user runs `sysi change propose <name>` and `<name>` already exists as an active change or collides with an archived change name in the workspace
- **THEN** the system reports the conflict and does not create the change

### Requirement: Track Change Status
The system SHALL track each change's lifecycle status in its `meta.json` using the statuses `proposed`, `applying`, and `archived`.

#### Scenario: Status progresses through the lifecycle
- **WHEN** a change is proposed, then applied, then archived
- **THEN** its `meta.json` status is `proposed`, then `applying`, then `archived`, with updated timestamps at each transition

### Requirement: Apply Prints Superpowers Handoff
The system SHALL make `sysi change apply` mark the change as applying and print the apply-discipline handoff.

#### Scenario: Apply build change
- **WHEN** a user runs `sysi change apply add-login` from inside the owning workspace
- **THEN** the system verifies the project is in build phase and the change exists
- **AND** the system sets the change status to `applying`
- **AND** the system prints instructions to read `proposal.md`, `design.md`, and `tasks.md` before editing implementation code, to work tasks in order with Superpowers discipline, and to check off tasks only after implementation and verification
- **AND** the system prints the design-drift stop conditions directing the agent to `sysi design-change`

#### Scenario: Apply missing change
- **WHEN** a user runs `sysi change apply <name>` and the change does not exist in the workspace
- **THEN** the system reports the change was not found and lists the changes that do exist

#### Scenario: Apply archived change rejected
- **WHEN** a user runs `sysi change apply <name>` for a change whose status is `archived`
- **THEN** the system reports that the change is archived and does not modify it

### Requirement: Require Superpowers Discipline During Apply
The system SHALL make Superpowers apply, debugging, testing, and verification discipline part of the build apply workflow.

#### Scenario: Codex apply skill invoked
- **WHEN** the Codex `sysi-apply` skill is invoked for a sysi change
- **THEN** the agent runs `sysi change apply <name>` and reads the change's `proposal.md`, `design.md`, and `tasks.md` before editing implementation code
- **AND** the agent uses Superpowers methods for implementation planning, test-driven work, debugging, and verification, checking off tasks in `tasks.md` as they are completed and verified
- **AND** the agent stops instead of continuing if the required Superpowers workflow is unavailable

### Requirement: Support Explicit Design Changes During Build
The system SHALL provide `sysi design-change` for foundational `/system` mutations during build phase.

#### Scenario: Foundation change requested during build
- **WHEN** a user runs `sysi design-change change-auth-boundary`
- **THEN** the system verifies the project is in build phase
- **AND** the system creates a design-change artifact under `system/architecture/decisions/<date>-change-auth-boundary.md`
- **AND** the artifact records rationale, affected `/system` files, impacted workspace changes, migration or compatibility notes, explicit confirmation, decision, and consequences

### Requirement: Require User-Confirmed Design Drift Handoff During Apply
The system SHALL require apply work to stop and use a user-confirmed design-change handoff when implementation needs drift from `/system` foundation truth.

#### Scenario: Apply detects API contract drift
- **WHEN** an agent applying a build change discovers that implementation requires a new endpoint or changed request or response payload shape not represented in `/system`
- **THEN** the agent stops ordinary implementation work
- **AND** the agent identifies the affected `/system` files and the mismatch between implementation needs and foundation truth
- **AND** the agent asks the user to confirm whether the foundation should change
- **AND** the agent uses `sysi design-change <name>` and the generated design-change workflow before mutating `/system` if the user agrees

#### Scenario: User rejects required foundation change
- **WHEN** an agent applying a build change detects design drift and the user does not agree to update `/system`
- **THEN** the agent does not continue with implementation that contradicts `/system`
- **AND** the agent reports that the change or implementation approach must be revised to fit the existing foundation truth

### Requirement: Archive Completed Build Changes
The system SHALL support archiving completed changes into a dated archive directory inside the owning workspace.

#### Scenario: Archive change
- **WHEN** a user runs `sysi change archive add-login` from inside the owning workspace
- **THEN** the system moves the change to `<workspace>/changes/archive/<date>-add-login/` and sets its status to `archived`

#### Scenario: Archive warns on unchecked tasks
- **WHEN** a user archives a change whose `tasks.md` still contains unchecked tasks
- **THEN** the system warns with the number of unchecked tasks and still archives the change

#### Scenario: Archive warns on unexpected status
- **WHEN** a user archives a change whose status is not `proposed` or `applying`
- **THEN** the system warns that the change has an unexpected status

#### Scenario: Archive target already exists
- **WHEN** the dated archive target directory already exists
- **THEN** the system reports the conflict and does not move the change

#### Scenario: Archive meta update fails
- **WHEN** the archived `meta.json` cannot be updated after the move
- **THEN** the system moves the change back to its active location and reports the error
