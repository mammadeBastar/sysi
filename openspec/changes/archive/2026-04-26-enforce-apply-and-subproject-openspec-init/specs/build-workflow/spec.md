## MODIFIED Requirements

### Requirement: Coordinate OpenSpec Changes During Build
The system SHALL provide build-phase commands that coordinate with the existing OpenSpec CLI from the inferred implementation workspace.

#### Scenario: Propose build change
- **WHEN** a user runs `sys change propose add-login` from `frontend/`
- **THEN** the system verifies the project is in build phase and invokes or instructs the equivalent OpenSpec proposal workflow for `add-login`
- **AND** the system runs the OpenSpec command from the `frontend/` workspace

#### Scenario: Build change outside implementation workspace
- **WHEN** a user runs `sys change propose add-login` from the monorepo root or `/system`
- **THEN** the system reports that build changes must run from `frontend/` or `backend/`

### Requirement: Require OpenSpec Apply Path
The system SHALL require implementation changes to flow through OpenSpec apply semantics.

#### Scenario: Apply build change
- **WHEN** a user runs `sys change apply add-login` from `frontend/`
- **THEN** the system verifies the project is in build phase
- **AND** the system verifies the OpenSpec change exists in the inferred implementation workspace
- **AND** the system invokes the OpenSpec apply instruction workflow for `add-login`
- **AND** the system reports that implementation must continue through OpenSpec apply plus Superpowers discipline

### Requirement: Require Superpowers Discipline During Apply
The system SHALL make Superpowers apply, debugging, testing, and verification discipline part of the build apply workflow.

#### Scenario: Codex apply skill invoked
- **WHEN** the Codex `sys-apply` skill is invoked for an OpenSpec change
- **THEN** the agent invokes the local OpenSpec apply skill or workflow before editing implementation code
- **AND** the agent uses Superpowers methods for implementation planning, test-driven work, debugging, and verification
- **AND** the agent stops instead of continuing if the required OpenSpec apply or Superpowers workflow is unavailable
