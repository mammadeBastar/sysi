## Purpose

Define `README.md` as the complete project documentation surface for Sysi users, contributors, and agents.

## Requirements

### Requirement: README Provides Complete Project Overview
The README SHALL explain what Sysi is, what problem it solves, and how `/system`, native sysi changes, and Superpowers cooperate.

#### Scenario: New reader opens README
- **WHEN** a new reader opens `README.md`
- **THEN** they can understand the project purpose, core mental model, and lifecycle without reading archived changes

### Requirement: README Documents Installation And Running From Source
The README SHALL document how to run the CLI from source, how to use the installed `sysi` command form, and that the CLI has no external dependencies.

#### Scenario: User wants to run the CLI locally
- **WHEN** a user reads the installation section
- **THEN** they find copy-pasteable examples for `go run ./cmd/sysi --help` and `go run ./cmd/sysi <command>`
- **AND** they learn that no external executable is required for any sysi command

### Requirement: README Documents Workspaces
The README SHALL document declared workspaces: the required `--workspaces` flag at first init, the `sysi workspace list|add|remove` commands, workspace naming rules, and reserved names.

#### Scenario: User wants to declare or manage workspaces
- **WHEN** a user reads the workspaces documentation
- **THEN** they learn that `sysi init` requires `--workspaces` on first init, that bare init prints guidance and fails, that workspace names are lowercase slugs with reserved names, and that `workspace remove` refuses with active changes unless `--force` is given and leaves the directory on disk

### Requirement: README Documents Design Phase Workflow
The README SHALL document design-phase commands and explain that build changes are not created during design.

#### Scenario: User wants to design a system
- **WHEN** a user reads the design phase section
- **THEN** they find the intended flow for `sysi explore`, `sysi capture`, decision records, and `sysi design freeze`

### Requirement: README Documents Build Phase Workflow
The README SHALL document the native build-phase change workflow and the Superpowers apply responsibilities.

#### Scenario: User wants to implement a feature
- **WHEN** a user reads the build phase section
- **THEN** they find how `sysi change propose`, `sysi change apply`, `sysi change archive`, and `sysi design-change` fit together
- **AND** they learn the change directory layout (`proposal.md`, `design.md`, `tasks.md`, `meta.json`), the status lifecycle, and that change commands must run from inside a declared workspace

### Requirement: README Documents System Foundation Structure
The README SHALL document the canonical `/system` tree and the role of key files.

#### Scenario: User wants to know where architecture truth lives
- **WHEN** a user reads the `/system` section
- **THEN** they can identify the purpose of architecture, contracts, contract conventions, contract errors, flows, per-workspace modules, data, database, security, observability, and decisions files

### Requirement: README Documents Agent Integrations
The README SHALL document Codex, Cursor, and Claude Code integration levels.

#### Scenario: User wants to install agent support
- **WHEN** a user reads the agent integration section
- **THEN** they can see that Codex installs project-local skills and Cursor/Claude receive minimal instruction scaffolds

### Requirement: README Documents Command Reference
The README SHALL include a command reference covering all v2 commands.

#### Scenario: User wants to find a command
- **WHEN** a user scans the command reference
- **THEN** they can find each v2 command, including the `sysi workspace` commands, its purpose, and a representative example

### Requirement: README Documents Validation And Troubleshooting
The README SHALL document status, validation, common warnings, and known limitations.

#### Scenario: User encounters unexpected status output
- **WHEN** a user reads troubleshooting documentation
- **THEN** they can understand validation warnings, the `--workspaces` init requirement, the declared-workspace requirement for change commands, build-phase capture blocking, and minimal adapter limitations
