## Purpose

Define `README.md` as the complete project documentation surface for Sysi users, contributors, and agents.

## Requirements

### Requirement: README Provides Complete Project Overview
The README SHALL explain what Sysi is, what problem it solves, and how `/system`, OpenSpec, and Superpowers cooperate.

#### Scenario: New reader opens README
- **WHEN** a new reader opens `README.md`
- **THEN** they can understand the project purpose, core mental model, and lifecycle without reading archived OpenSpec changes

### Requirement: README Documents Installation And Running From Source
The README SHALL document how to run the CLI from source and how to use the installed `sysi` command form.

#### Scenario: User wants to run the CLI locally
- **WHEN** a user reads the installation section
- **THEN** they find copy-pasteable examples for `go run ./cmd/sysi --help` and `go run ./cmd/sysi <command>`

### Requirement: README Documents Design Phase Workflow
The README SHALL document design-phase commands and explain that OpenSpec is not used for design decisions.

#### Scenario: User wants to design a system
- **WHEN** a user reads the design phase section
- **THEN** they find the intended flow for `sysi explore`, `sysi capture`, decision records, and `sysi design freeze`

### Requirement: README Documents Build Phase Workflow
The README SHALL document build-phase commands and explain the OpenSpec and Superpowers responsibilities.

#### Scenario: User wants to implement a feature
- **WHEN** a user reads the build phase section
- **THEN** they find how `sysi change propose`, `sysi change apply`, `sysi change archive`, and `sysi design-change` fit together

### Requirement: README Documents System Foundation Structure
The README SHALL document the canonical `/system` tree and the role of key files.

#### Scenario: User wants to know where architecture truth lives
- **WHEN** a user reads the `/system` section
- **THEN** they can identify the purpose of architecture, contracts, contract conventions, contract errors, flows, modules, data, database, security, observability, and decisions files

### Requirement: README Documents Agent Integrations
The README SHALL document Codex, Cursor, and Claude Code integration levels.

#### Scenario: User wants to install agent support
- **WHEN** a user reads the agent integration section
- **THEN** they can see that Codex installs project-local skills and Cursor/Claude receive minimal instruction scaffolds

### Requirement: README Documents Command Reference
The README SHALL include a command reference covering all v1 commands.

#### Scenario: User wants to find a command
- **WHEN** a user scans the command reference
- **THEN** they can find each v1 command, its purpose, and a representative example

### Requirement: README Documents Validation And Troubleshooting
The README SHALL document status, validation, common warnings, and known v1 limitations.

#### Scenario: User encounters unexpected status output
- **WHEN** a user reads troubleshooting documentation
- **THEN** they can understand validation warnings, build-phase capture blocking, OpenSpec telemetry errors, and minimal adapter limitations
