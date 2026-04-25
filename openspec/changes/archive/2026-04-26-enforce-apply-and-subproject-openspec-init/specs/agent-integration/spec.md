## MODIFIED Requirements

### Requirement: Codex Apply Skill Enforces Build Apply Boundaries
The system SHALL generate a `sys-apply` Codex skill that coordinates OpenSpec apply with Sys foundation rules and Superpowers discipline.

#### Scenario: Apply skill installed
- **WHEN** `sys-apply/SKILL.md` is generated
- **THEN** it requires build phase, mandatory OpenSpec apply invocation, mandatory Superpowers apply discipline, `/system` context review, frozen-file protection, design drift detection, and escalation to `sys-design-change` when foundational truth must change

#### Scenario: Apply skill lacks required external workflow
- **WHEN** the Codex `sys-apply` skill is invoked and the local OpenSpec apply or Superpowers workflow is unavailable
- **THEN** it instructs the agent to stop and report the missing prerequisite instead of implementing without the required apply discipline
