## ADDED Requirements

### Requirement: Scaffold Canonical System Tree
The system SHALL scaffold the canonical `/system` directory structure.

#### Scenario: System tree creation
- **WHEN** a user runs `sysi init`
- **THEN** the system creates architecture, contracts, flows, modules, data, database, observability, dashboard, and decision-record locations under `/system`

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
The system SHALL validate the presence of required `/system` files.

#### Scenario: Missing contracts file
- **WHEN** `system/contracts/api.yaml` is missing
- **THEN** `sysi validate` and `sysi status` report a warning identifying the missing file

### Requirement: Enforce Freeze Baselines
The system SHALL record and check freeze baselines for controlled and frozen `/system` files.

#### Scenario: Frozen architecture file changes during build phase
- **WHEN** the project is in build phase and `system/architecture/system.md` differs from its recorded freeze baseline
- **THEN** `sysi status` reports that a design-change workflow is required
