## Purpose

Define the canonical `/system` foundation, validation, decision records, and freeze behavior.

## Requirements

### Requirement: Scaffold Canonical System Tree
The system SHALL scaffold the canonical `/system` directory structure.

#### Scenario: System tree creation
- **WHEN** a user runs `sysi init --workspaces <names>`
- **THEN** the system creates architecture, contracts, flows, modules, data, database, security, observability, dashboard, and decision-record locations under `/system`

### Requirement: Scaffold Module Files Per Declared Workspace
The system SHALL scaffold one `system/modules/<workspace>.md` module file per declared workspace.

#### Scenario: Module files scaffolded at init
- **WHEN** a user runs `sysi init --workspaces frontend,backend`
- **THEN** the system creates `system/modules/frontend.md` and `system/modules/backend.md`

#### Scenario: Module file scaffolded on workspace add
- **WHEN** a user runs `sysi workspace add worker`
- **THEN** the system creates `system/modules/worker.md` when it is missing

### Requirement: Scaffold Contract Convention And Error Files
The system SHALL scaffold dedicated contract files for cross-cutting conventions and error behavior.

#### Scenario: Contract support files scaffolded
- **WHEN** a user runs `sysi init`
- **THEN** the system creates `system/contracts/conventions.md` for shared API and event conventions
- **AND** the system creates `system/contracts/errors.md` for shared error behavior

### Requirement: Scaffold Security Foundation
The system SHALL scaffold a minimal security foundation under `/system`.

#### Scenario: Security model scaffolded
- **WHEN** a user runs `sysi init`
- **THEN** the system creates `system/security/model.md` for trust boundaries, sensitive data rules, secret handling, and security invariants

### Requirement: Use Postgres SQL As Canonical Data Schema
The system SHALL treat `system/data/schema.sql` as the canonical database schema file.

#### Scenario: Schema files scaffolded
- **WHEN** `/system` is scaffolded
- **THEN** the system creates `system/data/schema.sql` for canonical database shape and `system/data/schema.md` for explanatory schema notes

### Requirement: Record Design Decisions
The system SHALL support append-only decision records under `system/architecture/decisions/`.

#### Scenario: Capture accepted decision
- **WHEN** a finalized design decision is captured
- **THEN** the system records a dated decision file containing status, decision, rationale, and affected system files

### Requirement: Validate Required System Files
The system SHALL validate the presence of required `/system` files, including the module file for every declared workspace.

#### Scenario: Missing contracts file
- **WHEN** `system/contracts/api.yaml` is missing
- **THEN** `sysi validate` and `sysi status` report a warning identifying the missing file

#### Scenario: Missing contract conventions file
- **WHEN** `system/contracts/conventions.md` is missing
- **THEN** `sysi validate` and `sysi status` report a warning identifying the missing file

#### Scenario: Missing security model file
- **WHEN** `system/security/model.md` is missing
- **THEN** `sysi validate` and `sysi status` report a warning identifying the missing file

#### Scenario: Missing workspace module file
- **WHEN** `system/modules/<workspace>.md` is missing for a declared workspace
- **THEN** `sysi validate` and `sysi status` report a warning identifying the missing file

### Requirement: Validate Workspace And Change Health
The system SHALL validate declared workspace directories and native change records.

#### Scenario: Missing workspace directory
- **WHEN** a declared workspace directory is missing or is not a directory
- **THEN** `sysi validate` and `sysi status` report a warning identifying the missing workspace directory

#### Scenario: Invalid change metadata
- **WHEN** an active change directory has a missing or unparseable `meta.json`
- **THEN** `sysi validate` and `sysi status` report a warning identifying the change

#### Scenario: Illegal change status
- **WHEN** an active change's `meta.json` records a status other than `proposed` or `applying`
- **THEN** `sysi validate` and `sysi status` report a warning identifying the invalid status

#### Scenario: Active change collides with archived name
- **WHEN** an active change name matches an archived change name in the same workspace
- **THEN** `sysi validate` and `sysi status` report a warning identifying the collision

### Requirement: Enforce Freeze Baselines
The system SHALL record and check freeze baselines for controlled and frozen `/system` files.

#### Scenario: Frozen architecture file changes during build phase
- **WHEN** the project is in build phase and `system/architecture/system.md` differs from its recorded freeze baseline
- **THEN** `sysi status` reports that a design-change workflow is required

#### Scenario: Controlled contract file changes during build phase
- **WHEN** the project is in build phase and `system/contracts/errors.md` differs from its recorded freeze baseline
- **THEN** `sysi status` reports that a design-change workflow is required

#### Scenario: Controlled security file changes during build phase
- **WHEN** the project is in build phase and `system/security/model.md` differs from its recorded freeze baseline
- **THEN** `sysi status` reports that a design-change workflow is required

### Requirement: Include Expanded Foundation In Role Allowlists
The system SHALL include contract support files and the security model in default role-based system-file allowlists for every declared workspace.

#### Scenario: Workspace role can read expanded foundation
- **WHEN** a user invokes a sysi agent workflow from inside a declared workspace `<name>`
- **THEN** the `<name>` role allowlist includes `system/architecture/system.md`, `system/contracts/**`, `system/flows/**`, `system/modules/<name>.md`, `system/data/**`, `system/obs/**`, and `system/security/**`
