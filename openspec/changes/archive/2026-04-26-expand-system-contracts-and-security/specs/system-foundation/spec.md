## ADDED Requirements

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

### Requirement: Include Expanded Foundation In Role Allowlists
The system SHALL include contract support files and the security model in default role-based system-file allowlists.

#### Scenario: Frontend role can read security foundation
- **WHEN** a user invokes a sysi agent workflow from `frontend/`
- **THEN** the frontend role allowlist includes `system/contracts/**`, `system/flows/**`, `system/modules/frontend.md`, and `system/security/**`

#### Scenario: Backend role can read security foundation
- **WHEN** a user invokes a sysi agent workflow from `backend/`
- **THEN** the backend role allowlist includes `system/contracts/**`, `system/flows/**`, `system/modules/backend.md`, `system/data/**`, `system/obs/**`, and `system/security/**`

## MODIFIED Requirements

### Requirement: Scaffold Canonical System Tree
The system SHALL scaffold the canonical `/system` directory structure.

#### Scenario: System tree creation
- **WHEN** a user runs `sysi init`
- **THEN** the system creates architecture, contracts, flows, modules, data, database, security, observability, dashboard, and decision-record locations under `/system`

### Requirement: Validate Required System Files
The system SHALL validate the presence of required `/system` files.

#### Scenario: Missing contracts file
- **WHEN** `system/contracts/api.yaml` is missing
- **THEN** `sysi validate` and `sysi status` report a warning identifying the missing file

#### Scenario: Missing contract conventions file
- **WHEN** `system/contracts/conventions.md` is missing
- **THEN** `sysi validate` and `sysi status` report a warning identifying the missing file

#### Scenario: Missing security model file
- **WHEN** `system/security/model.md` is missing
- **THEN** `sysi validate` and `sysi status` report a warning identifying the missing file

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
