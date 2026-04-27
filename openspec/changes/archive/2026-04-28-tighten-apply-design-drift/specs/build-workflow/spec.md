## ADDED Requirements

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
- **AND** the agent reports that the OpenSpec change or implementation approach must be revised to fit the existing foundation truth
