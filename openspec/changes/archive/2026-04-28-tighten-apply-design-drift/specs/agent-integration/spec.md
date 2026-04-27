## MODIFIED Requirements

### Requirement: Codex Apply Skill Enforces Build Apply Boundaries
The system SHALL generate a `sysi-apply` Codex skill that coordinates OpenSpec apply with Sysi foundation rules and Superpowers discipline.

#### Scenario: Apply skill installed
- **WHEN** `sysi-apply/SKILL.md` is generated
- **THEN** it requires build phase, mandatory OpenSpec apply invocation, mandatory Superpowers apply discipline, `/system` context review including contracts, flows, modules, data, observability, and security where relevant, frozen-file protection, concrete design drift detection examples, user confirmation before foundation mutation, and escalation to `sysi-design-change` through `sysi design-change <name>` when foundational truth must change

#### Scenario: Apply skill lacks required external workflow
- **WHEN** the Codex `sysi-apply` skill is invoked and the local OpenSpec apply or Superpowers workflow is unavailable
- **THEN** it instructs the agent to stop and report the missing prerequisite instead of implementing without the required apply discipline

#### Scenario: Apply skill detects design drift
- **WHEN** the Codex `sysi-apply` skill detects that implementation requires a new or changed endpoint, payload shape, event, auth rule, error behavior, data shape, security invariant, or observability contract not represented in `/system`
- **THEN** it instructs the agent to stop ordinary implementation work
- **AND** it instructs the agent to double-check the drift with the user before changing `/system`
- **AND** it instructs the agent to invoke `sysi design-change <name>` and follow `sysi-design-change` before mutating controlled or frozen foundation files if the user agrees
- **AND** it instructs the agent not to continue implementation that contradicts `/system` if the user does not agree to the foundation change
