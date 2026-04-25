package sysapp

import (
	"context"
	"crypto/sha256"
	"embed"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

const (
	PhaseDesign = "design"
	PhaseBuild  = "build"

	RoleDesign   = "design"
	RoleFrontend = "frontend"
	RoleBackend  = "backend"
	RoleSystem   = "system-maintainer"
	RoleChange   = "change"
)

type Options struct {
	Dir          string
	Stdout       io.Writer
	Stderr       io.Writer
	OpenSpecPath string
	WatchCount   int
}

type App struct {
	opts Options
}

type State struct {
	Version       int               `json:"version"`
	Phase         string            `json:"phase"`
	CreatedAt     string            `json:"createdAt"`
	UpdatedAt     string            `json:"updatedAt"`
	AgentInstalls map[string]string `json:"agentInstalls,omitempty"`
}

type Freeze struct {
	Files map[string]FreezeFile `json:"files"`
}

type FreezeFile struct {
	Level  string `json:"level"`
	SHA256 string `json:"sha256"`
}

type Validation struct {
	OK       bool     `json:"ok"`
	Warnings []string `json:"warnings"`
}

type FreezeStatus struct {
	Baselines int      `json:"baselines"`
	Mutations []string `json:"mutations"`
}

type AgentStatus struct {
	Codex  bool `json:"codex"`
	Cursor bool `json:"cursor"`
	Claude bool `json:"claude"`
}

type OpenSpecStatus struct {
	Present       bool `json:"present"`
	ActiveChanges int  `json:"activeChanges"`
}

type Status struct {
	Root       string         `json:"root"`
	Phase      string         `json:"phase"`
	Role       string         `json:"role"`
	Validation Validation     `json:"validation"`
	Freeze     FreezeStatus   `json:"freeze"`
	Agents     AgentStatus    `json:"agents"`
	OpenSpec   OpenSpecStatus `json:"openspec"`
}

var requiredSystemFiles = []string{
	"system/architecture/system.md",
	"system/contracts/api.yaml",
	"system/contracts/events.asyncapi.yaml",
	"system/contracts/auth.md",
	"system/modules/frontend.md",
	"system/modules/backend.md",
	"system/data/schema.sql",
	"system/data/schema.md",
	"system/data/db/indexes.md",
	"system/data/db/triggers.md",
	"system/data/db/functions.md",
	"system/obs/metrics.md",
	"system/obs/logging.md",
	"system/obs/tracing.md",
	"system/obs/alerts.md",
	"system/obs/dashboards/grafana.md",
}

var controlledSystemFiles = []string{
	"system/architecture/system.md",
	"system/contracts/api.yaml",
	"system/contracts/events.asyncapi.yaml",
	"system/contracts/auth.md",
	"system/data/schema.sql",
}

var implementationOpenSpecTargets = []string{
	"frontend",
	"backend",
}

//go:embed templates/agents/codex/*/SKILL.md templates/agents/cursor/sys-orchestrator.mdc templates/agents/claude/CLAUDE.section.md
var agentTemplates embed.FS

func New(opts Options) *App {
	if opts.Dir == "" {
		if wd, err := os.Getwd(); err == nil {
			opts.Dir = wd
		} else {
			opts.Dir = "."
		}
	}
	if opts.Stdout == nil {
		opts.Stdout = os.Stdout
	}
	if opts.Stderr == nil {
		opts.Stderr = os.Stderr
	}
	return &App{opts: opts}
}

func (a *App) Run(args []string) int {
	if len(args) == 0 || args[0] == "help" || args[0] == "--help" || args[0] == "-h" {
		a.printHelp()
		return 0
	}

	var err error
	switch args[0] {
	case "init":
		err = a.init()
	case "status":
		err = a.status(args[1:])
	case "validate":
		err = a.validateCmd()
	case "design":
		err = a.design(args[1:])
	case "explore":
		err = a.explore(args[1:])
	case "capture":
		err = a.capture()
	case "design-change":
		err = a.designChange(args[1:])
	case "change":
		err = a.change(args[1:])
	case "agent":
		err = a.agent(args[1:])
	default:
		err = fmt.Errorf("unknown command %q", args[0])
	}

	if err != nil {
		fmt.Fprintln(a.opts.Stderr, "error:", err)
		return 1
	}
	return 0
}

func (a *App) printHelp() {
	fmt.Fprintln(a.opts.Stdout, `sys orchestrates agent-native system design and implementation.

Usage:
  sys init
  sys status [--json|--watch]
  sys validate
  sys design start|freeze
  sys explore [topic]
  sys capture
  sys design-change <name>
  sys change propose|apply|archive <name>
  sys agent install codex|cursor|claude`)
}

func (a *App) init() error {
	start, err := filepath.Abs(a.opts.Dir)
	if err != nil {
		return err
	}
	if root, ok := findRoot(start); ok {
		if err := a.ensureImplementationOpenSpec(root); err != nil {
			return err
		}
		fmt.Fprintf(a.opts.Stdout, "sys already initialized at %s\n", root)
		return nil
	}

	root := start
	now := time.Now().UTC().Format(time.RFC3339)
	state := State{
		Version:       1,
		Phase:         PhaseDesign,
		CreatedAt:     now,
		UpdatedAt:     now,
		AgentInstalls: map[string]string{},
	}

	if err := os.MkdirAll(filepath.Join(root, ".sys-orchestrator", "captures"), 0o755); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Join(root, ".sys-orchestrator", "agents"), 0o755); err != nil {
		return err
	}
	if err := saveJSON(filepath.Join(root, ".sys-orchestrator", "state.json"), state); err != nil {
		return err
	}
	if err := saveJSON(filepath.Join(root, ".sys-orchestrator", "freeze.json"), Freeze{Files: map[string]FreezeFile{}}); err != nil {
		return err
	}
	if err := saveJSON(filepath.Join(root, ".sys-orchestrator", "allowlists.json"), defaultAllowlists()); err != nil {
		return err
	}
	if err := scaffoldSystem(root); err != nil {
		return err
	}
	if err := a.ensureImplementationOpenSpec(root); err != nil {
		return err
	}

	fmt.Fprintf(a.opts.Stdout, "initialized sys repository at %s\n", root)
	fmt.Fprintln(a.opts.Stdout, "next: sys status")
	return nil
}

func (a *App) status(args []string) error {
	jsonOut := hasFlag(args, "--json")
	watch := hasFlag(args, "--watch")

	root, state, err := a.requireProject()
	if err != nil {
		if jsonOut {
			return err
		}
		fmt.Fprintln(a.opts.Stdout, "sys is not initialized here.")
		fmt.Fprintln(a.opts.Stdout, "Run `sys init` at the monorepo root.")
		return err
	}

	render := func() error {
		status, err := a.buildStatus(root, state)
		if err != nil {
			return err
		}
		if jsonOut {
			enc := json.NewEncoder(a.opts.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(status)
		}
		a.renderStatus(status)
		return nil
	}

	if !watch {
		return render()
	}

	limit := a.opts.WatchCount
	for i := 0; limit == 0 || i < limit; i++ {
		if i > 0 {
			fmt.Fprint(a.opts.Stdout, "\033[H\033[2J")
		}
		if err := render(); err != nil {
			return err
		}
		time.Sleep(time.Second)
	}
	return nil
}

func (a *App) validateCmd() error {
	root, state, err := a.requireProject()
	if err != nil {
		return err
	}
	status, err := a.buildStatus(root, state)
	if err != nil {
		return err
	}
	if status.Validation.OK {
		fmt.Fprintln(a.opts.Stdout, "system validation passed")
		return nil
	}
	for _, warning := range status.Validation.Warnings {
		fmt.Fprintln(a.opts.Stdout, "warning:", warning)
	}
	return errors.New("system validation failed")
}

func (a *App) design(args []string) error {
	if len(args) == 0 {
		return errors.New("usage: sys design start|freeze")
	}

	root, state, err := a.requireProject()
	if err != nil {
		return err
	}

	switch args[0] {
	case "start":
		state.Phase = PhaseDesign
		state.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
		if err := saveState(root, state); err != nil {
			return err
		}
		fmt.Fprintln(a.opts.Stdout, "design phase active")
	case "freeze":
		freeze, err := computeFreeze(root)
		if err != nil {
			return err
		}
		state.Phase = PhaseBuild
		state.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
		if err := saveJSON(filepath.Join(root, ".sys-orchestrator", "freeze.json"), freeze); err != nil {
			return err
		}
		if err := saveState(root, state); err != nil {
			return err
		}
		fmt.Fprintln(a.opts.Stdout, "design frozen; build phase active")
	default:
		return fmt.Errorf("unknown design command %q", args[0])
	}
	return nil
}

func (a *App) explore(args []string) error {
	root, state, err := a.requireProject()
	if err != nil {
		return err
	}
	topic := "system"
	if len(args) > 0 {
		topic = strings.Join(args, " ")
	}
	role := inferRole(root, a.opts.Dir)
	fmt.Fprintf(a.opts.Stdout, "SYS EXPLORE\n\nTopic: %s\nPhase: %s\nRole: %s\n\n", topic, state.Phase, role)
	fmt.Fprintln(a.opts.Stdout, "Use current /system files as the project foundation.")
	fmt.Fprintln(a.opts.Stdout, "During design phase, do not create OpenSpec changes.")
	fmt.Fprintln(a.opts.Stdout, "When decisions are final, invoke sys capture or the Codex sys-capture skill.")
	fmt.Fprintln(a.opts.Stdout, "\nAllowed system files:")
	for _, allowed := range allowlistForRole(root, role) {
		fmt.Fprintf(a.opts.Stdout, "- %s\n", allowed)
	}
	return nil
}

func (a *App) capture() error {
	_, state, err := a.requireProject()
	if err != nil {
		return err
	}
	if state.Phase == PhaseBuild {
		return errors.New("normal capture is blocked in build phase; use sys design-change")
	}
	fmt.Fprintln(a.opts.Stdout, "SYS CAPTURE")
	fmt.Fprintln(a.opts.Stdout, "Capture only finalized decisions.")
	fmt.Fprintln(a.opts.Stdout, "Update the relevant /system files and add a decision record under system/architecture/decisions/.")
	fmt.Fprintln(a.opts.Stdout, "Each decision record should include status, decision, rationale, and affected files.")
	return nil
}

func (a *App) designChange(args []string) error {
	if len(args) == 0 {
		return errors.New("usage: sys design-change <name>")
	}
	root, state, err := a.requireProject()
	if err != nil {
		return err
	}
	fmt.Fprintf(a.opts.Stdout, "SYS DESIGN CHANGE: %s\n\n", args[0])
	fmt.Fprintf(a.opts.Stdout, "Phase: %s\n", state.Phase)
	fmt.Fprintf(a.opts.Stdout, "Root: %s\n\n", root)
	fmt.Fprintln(a.opts.Stdout, "Record rationale, affected /system files, impacted OpenSpec changes, and migration notes before mutating frozen foundation files.")
	return nil
}

func (a *App) change(args []string) error {
	if len(args) < 2 {
		return errors.New("usage: sys change propose|apply|archive <name>")
	}
	root, state, err := a.requireProject()
	if err != nil {
		return err
	}
	if state.Phase != PhaseBuild {
		return errors.New("build changes require build phase; run sys design freeze first")
	}
	openSpecDir, err := a.requireImplementationOpenSpecDir(root)
	if err != nil {
		return err
	}

	action, name := args[0], args[1]
	switch action {
	case "propose":
		if err := a.runOpenSpec(openSpecDir, "new", "change", name); err != nil {
			return err
		}
		fmt.Fprintf(a.opts.Stdout, "OpenSpec change proposed: %s\n", name)
	case "apply":
		if _, err := os.Stat(filepath.Join(openSpecDir, "openspec", "changes", name)); err != nil {
			return fmt.Errorf("OpenSpec change %q not found", name)
		}
		if err := a.runOpenSpec(openSpecDir, "instructions", "apply", "--change", name, "--json"); err != nil {
			return err
		}
		fmt.Fprintf(a.opts.Stdout, "OpenSpec apply instructions loaded for %s; continue implementation through OpenSpec apply and Superpowers discipline.\n", name)
	case "archive":
		if err := a.runOpenSpec(openSpecDir, "archive", name); err != nil {
			return err
		}
		fmt.Fprintf(a.opts.Stdout, "OpenSpec change archived: %s\n", name)
		_, _ = a.buildStatus(root, state)
	default:
		return fmt.Errorf("unknown change command %q", action)
	}
	return nil
}

func (a *App) agent(args []string) error {
	if len(args) != 2 || args[0] != "install" {
		return errors.New("usage: sys agent install codex|cursor|claude")
	}
	root, state, err := a.requireProject()
	if err != nil {
		return err
	}

	agent := args[1]
	switch agent {
	case "codex":
		err = installCodex(root)
	case "cursor":
		err = installCursor(root)
	case "claude":
		err = installClaude(root)
	default:
		return fmt.Errorf("unsupported agent %q", agent)
	}
	if err != nil {
		return err
	}
	if state.AgentInstalls == nil {
		state.AgentInstalls = map[string]string{}
	}
	state.AgentInstalls[agent] = time.Now().UTC().Format(time.RFC3339)
	state.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	if err := saveState(root, state); err != nil {
		return err
	}
	fmt.Fprintf(a.opts.Stdout, "installed %s integration\n", agent)
	return nil
}

func (a *App) requireProject() (string, State, error) {
	root, ok := findRoot(a.opts.Dir)
	if !ok {
		return "", State{}, errors.New("sys project not initialized")
	}
	state, err := loadState(root)
	return root, state, err
}

func (a *App) buildStatus(root string, state State) (Status, error) {
	validation, freezeStatus := validateSystem(root, state)
	return Status{
		Root:       root,
		Phase:      state.Phase,
		Role:       inferRole(root, a.opts.Dir),
		Validation: validation,
		Freeze:     freezeStatus,
		Agents:     agentStatus(root),
		OpenSpec:   openSpecStatus(root),
	}, nil
}

func (a *App) renderStatus(status Status) {
	health := "ok"
	if !status.Validation.OK {
		health = fmt.Sprintf("%d warning(s)", len(status.Validation.Warnings))
	}
	fmt.Fprintln(a.opts.Stdout, "SYS ORCHESTRATOR")
	fmt.Fprintf(a.opts.Stdout, "Root: %s\n", status.Root)
	fmt.Fprintf(a.opts.Stdout, "Phase: %s\n", status.Phase)
	fmt.Fprintf(a.opts.Stdout, "Role: %s\n", status.Role)
	fmt.Fprintf(a.opts.Stdout, "System health: %s\n", health)
	fmt.Fprintf(a.opts.Stdout, "Freeze baselines: %d\n", status.Freeze.Baselines)
	fmt.Fprintf(a.opts.Stdout, "OpenSpec changes: %d\n", status.OpenSpec.ActiveChanges)
	fmt.Fprintf(a.opts.Stdout, "Agents: codex=%t cursor=%t claude=%t\n", status.Agents.Codex, status.Agents.Cursor, status.Agents.Claude)
	if len(status.Validation.Warnings) > 0 {
		fmt.Fprintln(a.opts.Stdout, "\nWarnings:")
		for _, warning := range status.Validation.Warnings {
			fmt.Fprintf(a.opts.Stdout, "- %s\n", warning)
		}
	}
}

func findRoot(start string) (string, bool) {
	abs, err := filepath.Abs(start)
	if err != nil {
		return "", false
	}
	current := abs
	for {
		if _, err := os.Stat(filepath.Join(current, ".sys-orchestrator", "state.json")); err == nil {
			return current, true
		}
		parent := filepath.Dir(current)
		if parent == current {
			return "", false
		}
		current = parent
	}
}

func loadState(root string) (State, error) {
	var state State
	if err := loadJSON(filepath.Join(root, ".sys-orchestrator", "state.json"), &state); err != nil {
		return State{}, err
	}
	if state.Phase == "" {
		state.Phase = PhaseDesign
	}
	if state.AgentInstalls == nil {
		state.AgentInstalls = map[string]string{}
	}
	return state, nil
}

func saveState(root string, state State) error {
	return saveJSON(filepath.Join(root, ".sys-orchestrator", "state.json"), state)
}

func loadFreeze(root string) Freeze {
	var freeze Freeze
	if err := loadJSON(filepath.Join(root, ".sys-orchestrator", "freeze.json"), &freeze); err != nil {
		return Freeze{Files: map[string]FreezeFile{}}
	}
	if freeze.Files == nil {
		freeze.Files = map[string]FreezeFile{}
	}
	return freeze
}

func saveJSON(path string, value any) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return os.WriteFile(path, data, 0o644)
}

func loadJSON(path string, value any) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, value)
}

func scaffoldSystem(root string) error {
	dirs := []string{
		"system/architecture/decisions",
		"system/contracts",
		"system/flows",
		"system/modules",
		"system/data/db",
		"system/obs/dashboards",
	}
	for _, dir := range dirs {
		if err := os.MkdirAll(filepath.Join(root, dir), 0o755); err != nil {
			return err
		}
	}
	files := map[string]string{
		"system/architecture/system.md":         "# System Architecture\n\nDescribe services, applications, responsibilities, communication patterns, technical decisions, and system-wide invariants.\n",
		"system/contracts/api.yaml":             "openapi: 3.1.0\ninfo:\n  title: System API\n  version: 0.1.0\npaths: {}\n",
		"system/contracts/events.asyncapi.yaml": "asyncapi: 3.0.0\ninfo:\n  title: System Events\n  version: 0.1.0\nchannels: {}\noperations: {}\ncomponents:\n  messages: {}\n",
		"system/contracts/auth.md":              "# Auth Contract\n\nDescribe authentication, authorization, sessions, tokens, permissions, and boundary rules.\n",
		"system/modules/frontend.md":            "# Frontend Modules\n\nDescribe Next.js pages, components, responsibilities, and dependencies.\n",
		"system/modules/backend.md":             "# Backend Modules\n\nDescribe Go services, modules, responsibilities, and dependencies.\n",
		"system/data/schema.sql":                "-- Canonical Postgres schema.\n",
		"system/data/schema.md":                 "# Data Schema\n\nExplain database tables, relationships, invariants, protobuf files, and schema rationale. `schema.sql` is canonical for Postgres.\n",
		"system/data/db/indexes.md":             "# Database Indexes\n\nDocument database indexes and why they exist.\n",
		"system/data/db/triggers.md":            "# Database Triggers\n\nDocument database triggers and their invariants.\n",
		"system/data/db/functions.md":           "# Database Functions\n\nDocument database functions and their invariants.\n",
		"system/obs/metrics.md":                 "# Metrics\n\nDocument metrics exposed to `/metrics` and why they exist.\n",
		"system/obs/logging.md":                 "# Logging\n\nDocument logging strategy, required fields, and retention expectations.\n",
		"system/obs/tracing.md":                 "# Tracing\n\nDocument tracing strategy and span boundaries.\n",
		"system/obs/alerts.md":                  "# Alerts\n\nDocument alert rules and escalation expectations.\n",
		"system/obs/dashboards/grafana.md":      "# Grafana Dashboards\n\nDocument the ideal dashboard layout based on exposed metrics.\n",
	}
	for rel, content := range files {
		if err := writeFileIfMissing(filepath.Join(root, rel), content); err != nil {
			return err
		}
	}
	return nil
}

func writeFileIfMissing(path, content string) error {
	if _, err := os.Stat(path); err == nil {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(content), 0o644)
}

func (a *App) requireImplementationOpenSpecDir(root string) (string, error) {
	switch inferRole(root, a.opts.Dir) {
	case RoleFrontend:
		return filepath.Join(root, "frontend"), nil
	case RoleBackend:
		return filepath.Join(root, "backend"), nil
	default:
		return "", errors.New("build changes require an implementation workspace; run from frontend/ or backend/")
	}
}

func (a *App) ensureImplementationOpenSpec(root string) error {
	for _, target := range implementationOpenSpecTargets {
		if err := os.MkdirAll(filepath.Join(root, target), 0o755); err != nil {
			return err
		}
		if exists(filepath.Join(root, target, "openspec", "config.yaml")) {
			continue
		}
		if err := a.runOpenSpec(root, "init", target, "--tools", "none"); err != nil {
			return fmt.Errorf("initialize OpenSpec for %s: %w", target, err)
		}
	}
	return nil
}

func defaultAllowlists() map[string][]string {
	return map[string][]string{
		RoleDesign: {
			"system/**",
		},
		RoleFrontend: {
			"system/architecture/system.md",
			"system/contracts/**",
			"system/flows/**",
			"system/modules/frontend.md",
		},
		RoleBackend: {
			"system/architecture/system.md",
			"system/contracts/**",
			"system/flows/**",
			"system/modules/backend.md",
			"system/data/**",
			"system/obs/**",
		},
		RoleSystem: {
			"system/**",
		},
		RoleChange: {
			"openspec/changes/**",
			"system/**",
		},
	}
}

func allowlistForRole(root, role string) []string {
	var lists map[string][]string
	if err := loadJSON(filepath.Join(root, ".sys-orchestrator", "allowlists.json"), &lists); err != nil {
		lists = defaultAllowlists()
	}
	allowed := append([]string(nil), lists[role]...)
	if len(allowed) == 0 {
		allowed = append([]string(nil), lists[RoleDesign]...)
	}
	sort.Strings(allowed)
	return allowed
}

func computeFreeze(root string) (Freeze, error) {
	freeze := Freeze{Files: map[string]FreezeFile{}}
	for _, rel := range controlledSystemFiles {
		sum, err := hashFile(filepath.Join(root, rel))
		if err != nil {
			return Freeze{}, fmt.Errorf("cannot freeze %s: %w", rel, err)
		}
		level := "controlled"
		if rel == "system/architecture/system.md" {
			level = "frozen"
		}
		freeze.Files[rel] = FreezeFile{Level: level, SHA256: sum}
	}
	return freeze, nil
}

func validateSystem(root string, state State) (Validation, FreezeStatus) {
	var warnings []string
	for _, rel := range requiredSystemFiles {
		if _, err := os.Stat(filepath.Join(root, rel)); err != nil {
			warnings = append(warnings, fmt.Sprintf("missing required file: %s", rel))
		}
	}

	freeze := loadFreeze(root)
	var mutations []string
	if state.Phase == PhaseBuild {
		for rel, baseline := range freeze.Files {
			sum, err := hashFile(filepath.Join(root, rel))
			if err != nil {
				warnings = append(warnings, fmt.Sprintf("frozen file missing: %s requires sys design-change", rel))
				mutations = append(mutations, rel)
				continue
			}
			if sum != baseline.SHA256 {
				warnings = append(warnings, fmt.Sprintf("%s changed after freeze; use sys design-change", rel))
				mutations = append(mutations, rel)
			}
		}
	}

	return Validation{OK: len(warnings) == 0, Warnings: warnings}, FreezeStatus{Baselines: len(freeze.Files), Mutations: mutations}
}

func hashFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:]), nil
}

func inferRole(root, dir string) string {
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return RoleDesign
	}
	rel, err := filepath.Rel(root, absDir)
	if err != nil || rel == "." {
		return RoleDesign
	}
	parts := strings.Split(filepath.ToSlash(rel), "/")
	if len(parts) >= 2 && parts[0] == "openspec" && parts[1] == "changes" {
		return RoleChange
	}
	switch parts[0] {
	case "frontend":
		return RoleFrontend
	case "backend":
		return RoleBackend
	case "system":
		return RoleSystem
	default:
		return RoleDesign
	}
}

func agentStatus(root string) AgentStatus {
	return AgentStatus{
		Codex:  exists(filepath.Join(root, ".codex", "skills", "sys-explore", "SKILL.md")),
		Cursor: exists(filepath.Join(root, ".cursor", "rules", "sys-orchestrator.mdc")),
		Claude: exists(filepath.Join(root, "CLAUDE.md")),
	}
}

func openSpecStatus(root string) OpenSpecStatus {
	changesDir := filepath.Join(root, "openspec", "changes")
	entries, err := os.ReadDir(changesDir)
	if err != nil {
		return OpenSpecStatus{Present: exists(filepath.Join(root, "openspec"))}
	}
	count := 0
	for _, entry := range entries {
		if entry.IsDir() {
			count++
		}
	}
	return OpenSpecStatus{Present: true, ActiveChanges: count}
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func hasFlag(args []string, flag string) bool {
	for _, arg := range args {
		if arg == flag {
			return true
		}
	}
	return false
}

func (a *App) runOpenSpec(root string, args ...string) error {
	bin := a.opts.OpenSpecPath
	if bin == "" {
		bin = os.Getenv("SYS_OPENSPEC")
	}
	if bin == "" {
		found, err := exec.LookPath("openspec")
		if err != nil {
			return errors.New("openspec executable not found")
		}
		bin = found
	}
	cmd := exec.CommandContext(context.Background(), bin, args...)
	cmd.Dir = root
	output, err := cmd.CombinedOutput()
	if len(output) > 0 {
		fmt.Fprint(a.opts.Stdout, string(output))
	}
	if err != nil {
		return fmt.Errorf("openspec %s failed: %w", strings.Join(args, " "), err)
	}
	return nil
}

func installCodex(root string) error {
	skills := map[string]string{
		"sys-explore":       "codex/sys-explore/SKILL.md",
		"sys-capture":       "codex/sys-capture/SKILL.md",
		"sys-apply":         "codex/sys-apply/SKILL.md",
		"sys-design-change": "codex/sys-design-change/SKILL.md",
	}
	for name, templatePath := range skills {
		content, err := agentInstructionTemplate(templatePath)
		if err != nil {
			return err
		}
		path := filepath.Join(root, ".codex", "skills", name, "SKILL.md")
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			return err
		}
	}
	return nil
}

func installCursor(root string) error {
	content, err := agentInstructionTemplate("cursor/sys-orchestrator.mdc")
	if err != nil {
		return err
	}
	path := filepath.Join(root, ".cursor", "rules", "sys-orchestrator.mdc")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(content), 0o644)
}

func installClaude(root string) error {
	path := filepath.Join(root, "CLAUDE.md")
	var existing string
	if data, err := os.ReadFile(path); err == nil {
		existing = string(data)
	}

	section, err := agentInstructionTemplate("claude/CLAUDE.section.md")
	if err != nil {
		return err
	}
	start := "<!-- SYS-ORCHESTRATOR:START -->"
	end := "<!-- SYS-ORCHESTRATOR:END -->"

	var next string
	if strings.Contains(existing, start) && strings.Contains(existing, end) {
		before := existing[:strings.Index(existing, start)]
		after := existing[strings.Index(existing, end)+len(end):]
		next = strings.TrimRight(before, "\n") + "\n\n" + section + strings.TrimLeft(after, "\n")
	} else if strings.TrimSpace(existing) == "" {
		next = section
	} else {
		next = strings.TrimRight(existing, "\n") + "\n\n" + section
	}
	return os.WriteFile(path, []byte(next), 0o644)
}

func agentInstructionTemplate(rel string) (string, error) {
	data, err := agentTemplates.ReadFile("templates/agents/" + rel)
	if err != nil {
		return "", fmt.Errorf("agent template %s not found: %w", rel, err)
	}
	return string(data), nil
}
