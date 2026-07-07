# Sysi V2 Native Core Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Remove the OpenSpec dependency from sysi and replace it with declared workspaces and a native change workflow (propose/apply/archive), per `docs/superpowers/specs/2026-07-07-sysi-v2-native-core-design.md`.

**Architecture:** Sysi stays a stdlib-only Go CLI. `.sysi/state.json` moves to schema version 2 with a `workspaces` field. Changes live as files in `<workspace>/changes/<name>/` (`proposal.md`, `design.md`, `tasks.md`, `meta.json`). New subsystems go in new focused files (`workspaces.go`, `changes.go`); removals happen in `app.go`. No migration from v1 — v1 projects keep the v1 binary.

**Tech Stack:** Go 1.22, standard library only. Tests via `go test` with `t.TempDir()`.

**Verification command (used throughout):** `GOCACHE=/tmp/sysi-go-cache go test ./...`
(The plain `go test ./...` form is fine when the default cache is writable.)

**Conventions for all tasks:**
- All tests live in package `sysiapp` and use the existing helpers `runApp`, `readFile`, `assertContainsAll` from `internal/sysiapp/app_test.go`.
- Every commit message follows the existing repo style: short imperative subject, no prefix (e.g. `status sees subfolders only, design-change creates artifact`). End each commit with the trailer `Co-Authored-By: Claude Fable 5 <noreply@anthropic.com>`.
- After implementing each task, run `gofmt -w` on touched files.

---

## File Structure Overview

| File | Role |
| --- | --- |
| `internal/sysiapp/app.go` | Existing. Modified: `Options`, `State`, dispatch, `init`, `explore`, `change`, status types/render, validation, allowlists, role inference. OpenSpec code deleted. |
| `internal/sysiapp/workspaces.go` | NEW. Workspace name validation, `--workspaces` flag parsing, workspace dirs, `sysi workspace list/add/remove`, `currentWorkspace`. |
| `internal/sysiapp/changes.go` | NEW. Change meta model, change file templates, propose/apply/archive implementations, change listing. |
| `internal/sysiapp/app_test.go` | Existing. OpenSpec-era tests deleted/rewritten. |
| `internal/sysiapp/workspaces_test.go` | NEW. Workspace command tests. |
| `internal/sysiapp/changes_test.go` | NEW. Change lifecycle tests. |
| `internal/sysiapp/templates/agents/**` | Rewritten for native workflow. |
| `README.md`, `openspec/specs/**` | Rewritten to describe v2. |

---

### Task 1: State v2, `sysi init --workspaces`, generalized roles and allowlists

This is the cutover task: init stops invoking OpenSpec entirely, workspaces become declared state, scaffolding/validation/roles/allowlists become workspace-driven. OpenSpec-era tests that test removed behavior are deleted here; later tasks add the native replacements.

**Files:**
- Create: `internal/sysiapp/workspaces.go`
- Modify: `internal/sysiapp/app.go`
- Modify: `internal/sysiapp/app_test.go`

- [ ] **Step 1: Write the failing tests**

In `internal/sysiapp/app_test.go`:

**Delete these tests entirely** (they test OpenSpec behavior that is being removed; native replacements arrive in Tasks 3–7):
- `TestInitScaffoldsProjectAndIsIdempotent` (rewritten below)
- `TestStatusAggregatesImplementationOpenSpecWorkspacesOnly` (replaced in Task 6)
- `TestInitUsesSysiOpenSpecEnvironment`
- `TestValidateReportsMissingImplementationOpenSpecWorkspace` (replaced in Task 7)
- `TestBuildWorkflowUsesFakeOpenSpecInBuildPhase` (replaced in Tasks 3–5)
- Helpers `runAppWithOpenSpec`, `writeFakeOpenSpec`, `shellQuote`

**Replace every remaining call** of the form `runAppWithOpenSpec(t, root, fakeOpenSpec, "init")` (and the accompanying `writeFakeOpenSpec` lines) with:

```go
if code, out, errOut := runApp(t, root, "init", "--workspaces", "frontend,backend"); code != 0 {
    t.Fatalf("init failed: code=%d stdout=%q stderr=%q", code, out, errOut)
}
```

**Add these new tests:**

```go
func TestInitRequiresWorkspacesFlag(t *testing.T) {
	root := t.TempDir()

	code, out, errOut := runApp(t, root, "init")
	if code == 0 {
		t.Fatalf("bare init should fail when not initialized: stdout=%q stderr=%q", out, errOut)
	}
	assertContainsAll(t, "bare init guidance", out+errOut, []string{
		"--workspaces",
		"sysi init --workspaces frontend,backend",
	})
	if _, err := os.Stat(filepath.Join(root, ".sysi")); err == nil {
		t.Fatalf("bare init must not create .sysi")
	}
}

func TestInitScaffoldsDeclaredWorkspacesAndIsIdempotent(t *testing.T) {
	root := t.TempDir()

	code, out, errOut := runApp(t, root, "init", "--workspaces", "api,web")
	if code != 0 {
		t.Fatalf("init failed: code=%d stdout=%q stderr=%q", code, out, errOut)
	}

	wantFiles := []string{
		".sysi/state.json",
		".sysi/freeze.json",
		".sysi/allowlists.json",
		"system/architecture/system.md",
		"system/contracts/api.yaml",
		"system/contracts/conventions.md",
		"system/contracts/errors.md",
		"system/modules/api.md",
		"system/modules/web.md",
		"system/security/model.md",
		"system/data/schema.sql",
		"system/obs/dashboards/grafana.md",
	}
	for _, rel := range wantFiles {
		if _, err := os.Stat(filepath.Join(root, rel)); err != nil {
			t.Fatalf("expected %s to exist: %v", rel, err)
		}
	}
	for _, ws := range []string{"api", "web"} {
		info, err := os.Stat(filepath.Join(root, ws, "changes"))
		if err != nil || !info.IsDir() {
			t.Fatalf("expected %s/changes directory: %v", ws, err)
		}
	}

	var state State
	if err := json.Unmarshal([]byte(readFile(t, filepath.Join(root, ".sysi", "state.json"))), &state); err != nil {
		t.Fatal(err)
	}
	if state.Version != 2 {
		t.Fatalf("state version = %d, want 2", state.Version)
	}
	if strings.Join(state.Workspaces, ",") != "api,web" {
		t.Fatalf("workspaces = %v, want [api web]", state.Workspaces)
	}

	// Idempotent: bare re-run reports already initialized, no flag needed.
	code, out, errOut = runApp(t, root, "init")
	if code != 0 {
		t.Fatalf("second init failed: code=%d stdout=%q stderr=%q", code, out, errOut)
	}
	if !strings.Contains(out, "already initialized") {
		t.Fatalf("second init should report already initialized, got %q", out)
	}
}

func TestInitRejectsInvalidWorkspaceNames(t *testing.T) {
	for _, invalid := range []string{"system", "Api", "a b", "-api", ""} {
		t.Run(invalid, func(t *testing.T) {
			root := t.TempDir()
			code, out, errOut := runApp(t, root, "init", "--workspaces", "good,"+invalid)
			if code == 0 {
				t.Fatalf("init should reject workspace name %q: stdout=%q stderr=%q", invalid, out, errOut)
			}
		})
	}
	root := t.TempDir()
	if code, out, errOut := runApp(t, root, "init", "--workspaces", "api,api"); code == 0 {
		t.Fatalf("init should reject duplicate workspace names: stdout=%q stderr=%q", out, errOut)
	}
}

func TestRoleInferenceUsesDeclaredWorkspaces(t *testing.T) {
	root := t.TempDir()
	if code, out, errOut := runApp(t, root, "init", "--workspaces", "api,web"); code != 0 {
		t.Fatalf("init failed: code=%d stdout=%q stderr=%q", code, out, errOut)
	}

	apiDir := filepath.Join(root, "api", "handlers")
	if err := os.MkdirAll(apiDir, 0o755); err != nil {
		t.Fatal(err)
	}

	cases := map[string]string{
		root:                              RoleDesign,
		apiDir:                            "api",
		filepath.Join(root, "system"):     RoleSystem,
		filepath.Join(root, "web"):        "web",
	}
	for dir, wantRole := range cases {
		code, out, errOut := runApp(t, dir, "status", "--json")
		if code != 0 {
			t.Fatalf("status in %s failed: code=%d stdout=%q stderr=%q", dir, code, out, errOut)
		}
		var status Status
		if err := json.Unmarshal([]byte(out), &status); err != nil {
			t.Fatal(err)
		}
		if status.Role != wantRole {
			t.Fatalf("role in %s = %q, want %q", dir, status.Role, wantRole)
		}
	}
}

func TestLoadStateRejectsV1State(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, ".sysi"), 0o755); err != nil {
		t.Fatal(err)
	}
	v1 := `{"version":1,"phase":"design","createdAt":"x","updatedAt":"x"}`
	if err := os.WriteFile(filepath.Join(root, ".sysi", "state.json"), []byte(v1), 0o644); err != nil {
		t.Fatal(err)
	}

	code, out, errOut := runApp(t, root, "status")
	if code == 0 {
		t.Fatalf("status should fail on v1 state: stdout=%q stderr=%q", out, errOut)
	}
	if !strings.Contains(out+errOut, "version") {
		t.Fatalf("v1 state error should mention version: stdout=%q stderr=%q", out, errOut)
	}
}
```

**Adjust existing tests:**
- `TestRootDiscoveryAndStatusJSON`: keep, init with `--workspaces frontend,backend`; the role assertion `status.Role != RoleFrontend` becomes `status.Role != "frontend"` (the `RoleFrontend` constant is deleted). The allowlist assertions stay (frontend/backend allowlists still contain `system/security/**`).
- `TestDesignCommandsMentionExpandedFoundationTargets`: init with `--workspaces frontend,backend`; assertions unchanged.
- All other init call sites: mechanical replacement as described above.

- [ ] **Step 2: Run tests to verify the new ones fail**

Run: `GOCACHE=/tmp/sysi-go-cache go test ./... 2>&1 | head -40`
Expected: compile errors (`State` has no `Workspaces`, `runAppWithOpenSpec` undefined where not yet cleaned, etc.) or failures of the new tests. That counts as failing.

- [ ] **Step 3: Implement**

Create `internal/sysiapp/workspaces.go`:

```go
package sysiapp

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var reservedWorkspaceNames = []string{"system", "docs", "openspec"}

func parseWorkspacesFlag(args []string) ([]string, bool, error) {
	var raw string
	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch {
		case arg == "--workspaces":
			if i+1 >= len(args) {
				return nil, true, errors.New("--workspaces requires a value, e.g. --workspaces frontend,backend")
			}
			raw = args[i+1]
			i++
		case strings.HasPrefix(arg, "--workspaces="):
			raw = strings.TrimPrefix(arg, "--workspaces=")
		}
	}
	if raw == "" {
		return nil, false, nil
	}
	var names []string
	for _, part := range strings.Split(raw, ",") {
		name := strings.TrimSpace(part)
		if err := validateWorkspaceName(name); err != nil {
			return nil, true, err
		}
		if containsString(names, name) {
			return nil, true, fmt.Errorf("duplicate workspace name %q", name)
		}
		names = append(names, name)
	}
	return names, true, nil
}

func validateWorkspaceName(name string) error {
	if name == "" {
		return errors.New("workspace name must not be empty")
	}
	if containsString(reservedWorkspaceNames, name) {
		return fmt.Errorf("workspace name %q is reserved", name)
	}
	first := rune(name[0])
	if first < 'a' || first > 'z' {
		return fmt.Errorf("workspace name %q must start with a lowercase letter", name)
	}
	for _, r := range name {
		isLower := r >= 'a' && r <= 'z'
		isDigit := r >= '0' && r <= '9'
		if !isLower && !isDigit && r != '-' {
			return fmt.Errorf("workspace name %q may only contain lowercase letters, digits, and hyphens", name)
		}
	}
	return nil
}

func ensureWorkspaceDirs(root string, workspaces []string) error {
	for _, ws := range workspaces {
		if err := os.MkdirAll(filepath.Join(root, ws, "changes"), 0o755); err != nil {
			return err
		}
	}
	return nil
}
```

In `internal/sysiapp/app.go`:

1. Constants: delete `RoleFrontend`, `RoleBackend`, `RoleChange`. Add `const stateVersion = 2`.

2. `State` gains `Workspaces []string \`json:"workspaces"\`` (after `UpdatedAt`).

3. `loadState`: after unmarshal and before the phase default, add:

```go
	if state.Version != stateVersion {
		return State{}, fmt.Errorf("state version %d is not supported; sysi v2 requires state version %d (v1 projects should keep using the v1 binary)", state.Version, stateVersion)
	}
```

4. `Run`: `case "init": err = a.init(args[1:])`.

5. Rewrite `init`:

```go
func (a *App) init(args []string) error {
	start, err := filepath.Abs(a.opts.Dir)
	if err != nil {
		return err
	}
	if root, ok := findRoot(start); ok {
		state, err := loadState(root)
		if err != nil {
			return err
		}
		if err := scaffoldSystem(root, state.Workspaces); err != nil {
			return err
		}
		if err := ensureAllowlists(root, state.Workspaces); err != nil {
			return err
		}
		if err := ensureWorkspaceDirs(root, state.Workspaces); err != nil {
			return err
		}
		fmt.Fprintf(a.opts.Stdout, "sysi already initialized at %s\n", root)
		return nil
	}

	workspaces, flagGiven, err := parseWorkspacesFlag(args)
	if err != nil {
		return err
	}
	if !flagGiven || len(workspaces) == 0 {
		fmt.Fprintln(a.opts.Stdout, `sysi init requires declared workspaces.

Usage:
  sysi init --workspaces <name>[,<name>...]

Examples:
  sysi init --workspaces frontend,backend
  sysi init --workspaces api,web,worker

Workspaces are the implementation directories where build changes live.`)
		return errors.New("missing --workspaces")
	}

	root := start
	now := time.Now().UTC().Format(time.RFC3339)
	state := State{
		Version:       stateVersion,
		Phase:         PhaseDesign,
		CreatedAt:     now,
		UpdatedAt:     now,
		Workspaces:    workspaces,
		AgentInstalls: map[string]string{},
	}

	if err := os.MkdirAll(filepath.Join(root, ".sysi", "captures"), 0o755); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Join(root, ".sysi", "agents"), 0o755); err != nil {
		return err
	}
	if err := saveJSON(filepath.Join(root, ".sysi", "state.json"), state); err != nil {
		return err
	}
	if err := saveJSON(filepath.Join(root, ".sysi", "freeze.json"), Freeze{Files: map[string]FreezeFile{}}); err != nil {
		return err
	}
	if err := saveJSON(filepath.Join(root, ".sysi", "allowlists.json"), defaultAllowlists(workspaces)); err != nil {
		return err
	}
	if err := scaffoldSystem(root, workspaces); err != nil {
		return err
	}
	if err := ensureWorkspaceDirs(root, workspaces); err != nil {
		return err
	}

	fmt.Fprintf(a.opts.Stdout, "initialized sysi repository at %s\n", root)
	fmt.Fprintf(a.opts.Stdout, "workspaces: %s\n", strings.Join(workspaces, ", "))
	fmt.Fprintln(a.opts.Stdout, "next: sysi status")
	return nil
}
```

6. `scaffoldSystem(root string, workspaces []string)`: remove the two hardcoded entries `system/modules/frontend.md` and `system/modules/backend.md` from the `files` map, and after the existing loop add:

```go
	for _, ws := range workspaces {
		content := fmt.Sprintf("# %s Modules\n\nDescribe %s components, responsibilities, and dependencies.\n", ws, ws)
		if err := writeFileIfMissing(filepath.Join(root, "system", "modules", ws+".md"), content); err != nil {
			return err
		}
	}
```

7. `requiredSystemFiles`: rename the var to `baseRequiredSystemFiles`, remove `system/modules/frontend.md` and `system/modules/backend.md` from it, and add:

```go
func requiredSystemFiles(workspaces []string) []string {
	files := append([]string(nil), baseRequiredSystemFiles...)
	for _, ws := range workspaces {
		files = append(files, "system/modules/"+ws+".md")
	}
	return files
}
```

8. `defaultAllowlists(workspaces []string)`:

```go
func defaultAllowlists(workspaces []string) map[string][]string {
	lists := map[string][]string{
		RoleDesign: {"system/**"},
		RoleSystem: {"system/**"},
	}
	for _, ws := range workspaces {
		lists[ws] = []string{
			"system/architecture/system.md",
			"system/contracts/**",
			"system/flows/**",
			"system/modules/" + ws + ".md",
			"system/data/**",
			"system/obs/**",
			"system/security/**",
		}
	}
	return lists
}
```

9. `ensureAllowlists(root string, workspaces []string)`: same merge logic, defaults come from `defaultAllowlists(workspaces)`.

10. `inferRole(root, dir string, workspaces []string) string`:

```go
func inferRole(root, dir string, workspaces []string) string {
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return RoleDesign
	}
	rel, err := filepath.Rel(root, absDir)
	if err != nil || rel == "." || strings.HasPrefix(rel, "..") {
		return RoleDesign
	}
	parts := strings.Split(filepath.ToSlash(rel), "/")
	if parts[0] == "system" {
		return RoleSystem
	}
	for _, ws := range workspaces {
		if parts[0] == ws {
			return ws
		}
	}
	return RoleDesign
}
```

10b. `allowlistForRole`: its fallback branch calls `defaultAllowlists()`, which now requires workspaces. Change the fallback to the role-independent minimum:

```go
func allowlistForRole(root, role string) []string {
	var lists map[string][]string
	if err := loadJSON(filepath.Join(root, ".sysi", "allowlists.json"), &lists); err != nil {
		lists = map[string][]string{RoleDesign: {"system/**"}, RoleSystem: {"system/**"}}
	}
	allowed := append([]string(nil), lists[role]...)
	if len(allowed) == 0 {
		allowed = append([]string(nil), lists[RoleDesign]...)
	}
	sort.Strings(allowed)
	return allowed
}
```

Update callers: `buildStatus` and `explore` pass `state.Workspaces`. The old `change` command's call to `requireImplementationOpenSpecDir` still compiles unchanged (that path is replaced in Task 3 and deleted in Task 8) — but `requireImplementationOpenSpecDir` itself references the deleted `RoleFrontend`/`RoleBackend` constants, so update it in place to keep the build green:

```go
func (a *App) requireImplementationOpenSpecDir(root string) (string, error) {
	return "", errors.New("build changes require an implementation workspace")
}
```

(The old OpenSpec change path is dead until Task 3 replaces it; no test covers it in between.)

11. `validateSystem`: replace the `requiredSystemFiles` loop source with `requiredSystemFiles(state.Workspaces)` and **delete** the `implementationOpenSpecTargets` warning loop (the `openspec/config.yaml` checks).

12. `ensureImplementationOpenSpec`: delete the call site in `init` (already gone in the rewrite above). Leave the function; Task 8 deletes it. If the compiler flags it as unused (it won't — it's a method), leave as is.

- [ ] **Step 4: Run tests to verify they pass**

Run: `GOCACHE=/tmp/sysi-go-cache go test ./...`
Expected: PASS (all packages).

- [ ] **Step 5: Commit**

```bash
git add -A
git commit -m "state v2 with declared workspaces, init drops openspec"
```

---

### Task 2: `sysi workspace list|add|remove`

**Files:**
- Modify: `internal/sysiapp/workspaces.go`
- Modify: `internal/sysiapp/app.go` (dispatch + help)
- Create: `internal/sysiapp/workspaces_test.go`

- [ ] **Step 1: Write the failing tests**

Create `internal/sysiapp/workspaces_test.go`:

```go
package sysiapp

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func initProject(t *testing.T, workspaces string) string {
	t.Helper()
	root := t.TempDir()
	if code, out, errOut := runApp(t, root, "init", "--workspaces", workspaces); code != 0 {
		t.Fatalf("init failed: code=%d stdout=%q stderr=%q", code, out, errOut)
	}
	return root
}

func TestWorkspaceListShowsDeclaredWorkspaces(t *testing.T) {
	root := initProject(t, "api,web")

	code, out, errOut := runApp(t, root, "workspace", "list")
	if code != 0 {
		t.Fatalf("workspace list failed: code=%d stdout=%q stderr=%q", code, out, errOut)
	}
	assertContainsAll(t, "workspace list", out, []string{"api", "web"})
}

func TestWorkspaceAddCreatesDirsAndModulesFile(t *testing.T) {
	root := initProject(t, "api")

	code, out, errOut := runApp(t, root, "workspace", "add", "worker")
	if code != 0 {
		t.Fatalf("workspace add failed: code=%d stdout=%q stderr=%q", code, out, errOut)
	}

	for _, rel := range []string{
		"worker/changes",
		"system/modules/worker.md",
	} {
		if _, err := os.Stat(filepath.Join(root, rel)); err != nil {
			t.Fatalf("expected %s to exist: %v", rel, err)
		}
	}

	code, out, _ = runApp(t, root, "workspace", "list")
	if code != 0 || !strings.Contains(out, "worker") {
		t.Fatalf("workspace list should include worker: %q", out)
	}

	// Duplicate add fails.
	if code, out, errOut := runApp(t, root, "workspace", "add", "worker"); code == 0 {
		t.Fatalf("duplicate workspace add should fail: stdout=%q stderr=%q", out, errOut)
	}
	// Invalid name fails.
	if code, out, errOut := runApp(t, root, "workspace", "add", "system"); code == 0 {
		t.Fatalf("reserved workspace add should fail: stdout=%q stderr=%q", out, errOut)
	}
}

func TestWorkspaceRemoveRefusesActiveChangesWithoutForce(t *testing.T) {
	root := initProject(t, "api,web")

	// Simulate an active change (native change scaffolding arrives in Task 3).
	changeDir := filepath.Join(root, "web", "changes", "add-login")
	if err := os.MkdirAll(changeDir, 0o755); err != nil {
		t.Fatal(err)
	}

	code, out, errOut := runApp(t, root, "workspace", "remove", "web")
	if code == 0 {
		t.Fatalf("remove with active changes should fail: stdout=%q stderr=%q", out, errOut)
	}
	if !strings.Contains(out+errOut, "add-login") {
		t.Fatalf("remove error should name the active change: stdout=%q stderr=%q", out, errOut)
	}

	code, out, errOut = runApp(t, root, "workspace", "remove", "web", "--force")
	if code != 0 {
		t.Fatalf("forced remove failed: code=%d stdout=%q stderr=%q", code, out, errOut)
	}
	if _, err := os.Stat(filepath.Join(root, "web")); err != nil {
		t.Fatalf("workspace directory should remain on disk after remove: %v", err)
	}

	code, out, _ = runApp(t, root, "workspace", "list")
	if code != 0 || strings.Contains(out, "web") {
		t.Fatalf("workspace list should no longer include web: %q", out)
	}

	// Removing an unknown workspace fails.
	if code, out, errOut := runApp(t, root, "workspace", "remove", "ghost"); code == 0 {
		t.Fatalf("remove of unknown workspace should fail: stdout=%q stderr=%q", out, errOut)
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `GOCACHE=/tmp/sysi-go-cache go test ./internal/sysiapp/ -run TestWorkspace -v`
Expected: FAIL with `unknown command "workspace"`.

- [ ] **Step 3: Implement**

In `internal/sysiapp/workspaces.go` add:

```go
func countActiveChangeNames(root, workspace string) []string {
	entries, err := os.ReadDir(filepath.Join(root, workspace, "changes"))
	if err != nil {
		return nil
	}
	var names []string
	for _, entry := range entries {
		if entry.IsDir() && entry.Name() != "archive" {
			names = append(names, entry.Name())
		}
	}
	return names
}

func (a *App) workspace(args []string) error {
	if len(args) == 0 {
		return errors.New("usage: sysi workspace list|add|remove <name> [--force]")
	}
	root, state, err := a.requireProject()
	if err != nil {
		return err
	}

	switch args[0] {
	case "list":
		for _, ws := range state.Workspaces {
			fmt.Fprintf(a.opts.Stdout, "%s: %d active change(s)\n", ws, len(countActiveChangeNames(root, ws)))
		}
		return nil
	case "add":
		if len(args) < 2 {
			return errors.New("usage: sysi workspace add <name>")
		}
		name := args[1]
		if err := validateWorkspaceName(name); err != nil {
			return err
		}
		if containsString(state.Workspaces, name) {
			return fmt.Errorf("workspace %q already declared", name)
		}
		state.Workspaces = append(state.Workspaces, name)
		state.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
		if err := ensureWorkspaceDirs(root, []string{name}); err != nil {
			return err
		}
		if err := scaffoldSystem(root, []string{name}); err != nil {
			return err
		}
		if err := ensureAllowlists(root, state.Workspaces); err != nil {
			return err
		}
		if err := saveState(root, state); err != nil {
			return err
		}
		fmt.Fprintf(a.opts.Stdout, "workspace added: %s\n", name)
		return nil
	case "remove":
		if len(args) < 2 {
			return errors.New("usage: sysi workspace remove <name> [--force]")
		}
		name := args[1]
		if !containsString(state.Workspaces, name) {
			return fmt.Errorf("workspace %q is not declared (declared: %s)", name, strings.Join(state.Workspaces, ", "))
		}
		active := countActiveChangeNames(root, name)
		if len(active) > 0 && !hasFlag(args, "--force") {
			return fmt.Errorf("workspace %q has active change(s): %s; use --force to remove anyway", name, strings.Join(active, ", "))
		}
		var kept []string
		for _, ws := range state.Workspaces {
			if ws != name {
				kept = append(kept, ws)
			}
		}
		state.Workspaces = kept
		state.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
		if err := saveState(root, state); err != nil {
			return err
		}
		fmt.Fprintf(a.opts.Stdout, "workspace removed from sysi state: %s (directory left on disk)\n", name)
		return nil
	default:
		return fmt.Errorf("unknown workspace command %q", args[0])
	}
}
```

Add `"time"` to the imports of `workspaces.go`.

In `app.go`:
- `Run` dispatch: add `case "workspace": err = a.workspace(args[1:])`.
- `printHelp`: replace the usage block with:

```go
	fmt.Fprintln(a.opts.Stdout, `sysi orchestrates agent-native system design and implementation.

Usage:
  sysi init --workspaces <name>[,<name>...]
  sysi status [--json|--watch]
  sysi validate
  sysi design start|freeze
  sysi explore [topic]
  sysi capture
  sysi design-change <name>
  sysi workspace list|add|remove <name> [--force]
  sysi change propose|apply|archive <name>
  sysi agent install codex|cursor|claude`)
```

Update `TestAgentInstallCommandNamesRemainStable`'s init call if not already done in Task 1 (it uses the shared init helper pattern).

- [ ] **Step 4: Run tests to verify they pass**

Run: `GOCACHE=/tmp/sysi-go-cache go test ./...`
Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add -A
git commit -m "workspace list, add, and remove commands"
```

---

### Task 3: Native change storage and `sysi change propose`

**Files:**
- Create: `internal/sysiapp/changes.go`
- Modify: `internal/sysiapp/app.go` (rewrite `change`)
- Create: `internal/sysiapp/changes_test.go`

- [ ] **Step 1: Write the failing tests**

Create `internal/sysiapp/changes_test.go`:

```go
package sysiapp

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func initBuildProject(t *testing.T, workspaces string) string {
	t.Helper()
	root := initProject(t, workspaces)
	if code, out, errOut := runApp(t, root, "design", "freeze"); code != 0 {
		t.Fatalf("design freeze failed: code=%d stdout=%q stderr=%q", code, out, errOut)
	}
	return root
}

func TestChangeProposeScaffoldsNativeChange(t *testing.T) {
	root := initBuildProject(t, "api,web")
	apiDir := filepath.Join(root, "api")

	code, out, errOut := runApp(t, apiDir, "change", "propose", "add-login")
	if code != 0 {
		t.Fatalf("change propose failed: code=%d stdout=%q stderr=%q", code, out, errOut)
	}

	base := filepath.Join(root, "api", "changes", "add-login")
	for _, rel := range []string{"proposal.md", "design.md", "tasks.md", "meta.json"} {
		if _, err := os.Stat(filepath.Join(base, rel)); err != nil {
			t.Fatalf("expected %s to exist: %v", rel, err)
		}
	}

	var meta ChangeMeta
	if err := json.Unmarshal([]byte(readFile(t, filepath.Join(base, "meta.json"))), &meta); err != nil {
		t.Fatal(err)
	}
	if meta.Name != "add-login" || meta.Workspace != "api" || meta.Status != ChangeStatusProposed {
		t.Fatalf("unexpected meta: %+v", meta)
	}

	proposal := readFile(t, filepath.Join(base, "proposal.md"))
	assertContainsAll(t, "proposal.md", proposal, []string{
		"# Change: add-login",
		"## Why",
		"## What Changes",
		"## Foundation Alignment",
		"sysi design-change",
		"## Out Of Scope",
	})
	design := readFile(t, filepath.Join(base, "design.md"))
	assertContainsAll(t, "design.md", design, []string{"## Decisions", "## Interfaces", "## Risks"})
	tasks := readFile(t, filepath.Join(base, "tasks.md"))
	assertContainsAll(t, "tasks.md", tasks, []string{"- [ ]", "/system"})
}

func TestChangeProposeGuardrails(t *testing.T) {
	root := initBuildProject(t, "api,web")
	apiDir := filepath.Join(root, "api")

	// Outside any workspace: error names declared workspaces.
	code, out, errOut := runApp(t, root, "change", "propose", "add-login")
	if code == 0 {
		t.Fatalf("propose at root should fail: stdout=%q stderr=%q", out, errOut)
	}
	assertContainsAll(t, "outside-workspace error", out+errOut, []string{"api", "web"})

	// Non-slug name fails.
	if code, out, errOut := runApp(t, apiDir, "change", "propose", "Add Login"); code == 0 {
		t.Fatalf("non-slug name should fail: stdout=%q stderr=%q", out, errOut)
	}

	// Duplicate fails.
	if code, _, _ := runApp(t, apiDir, "change", "propose", "add-login"); code != 0 {
		t.Fatal("first propose should succeed")
	}
	if code, out, errOut := runApp(t, apiDir, "change", "propose", "add-login"); code == 0 {
		t.Fatalf("duplicate propose should fail: stdout=%q stderr=%q", out, errOut)
	}

	// Name colliding with an archived change fails.
	archived := filepath.Join(root, "api", "changes", "archive", "2026-01-01-old-change")
	if err := os.MkdirAll(archived, 0o755); err != nil {
		t.Fatal(err)
	}
	if code, out, errOut := runApp(t, apiDir, "change", "propose", "old-change"); code == 0 {
		t.Fatalf("propose colliding with archive should fail: stdout=%q stderr=%q", out, errOut)
	}
}

func TestChangeProposeRequiresBuildPhase(t *testing.T) {
	root := initProject(t, "api")
	apiDir := filepath.Join(root, "api")

	code, out, errOut := runApp(t, apiDir, "change", "propose", "add-login")
	if code == 0 {
		t.Fatalf("propose in design phase should fail: stdout=%q stderr=%q", out, errOut)
	}
	if !strings.Contains(out+errOut, "build phase") {
		t.Fatalf("error should mention build phase: stdout=%q stderr=%q", out, errOut)
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `GOCACHE=/tmp/sysi-go-cache go test ./internal/sysiapp/ -run TestChange -v`
Expected: FAIL (compile error: `ChangeMeta` undefined).

- [ ] **Step 3: Implement**

Create `internal/sysiapp/changes.go`:

```go
package sysiapp

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

const (
	ChangeStatusProposed = "proposed"
	ChangeStatusApplying = "applying"
	ChangeStatusArchived = "archived"
)

type ChangeMeta struct {
	Name      string `json:"name"`
	Workspace string `json:"workspace"`
	Status    string `json:"status"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

func currentWorkspace(root, dir string, workspaces []string) (string, error) {
	role := inferRole(root, dir, workspaces)
	if containsString(workspaces, role) {
		return role, nil
	}
	return "", fmt.Errorf("change commands must run inside a declared workspace (declared: %s)", strings.Join(workspaces, ", "))
}

func changeDir(root, workspace, name string) string {
	return filepath.Join(root, workspace, "changes", name)
}

func loadChangeMeta(root, workspace, name string) (ChangeMeta, error) {
	var meta ChangeMeta
	err := loadJSON(filepath.Join(changeDir(root, workspace, name), "meta.json"), &meta)
	return meta, err
}

func saveChangeMeta(root, workspace, name string, meta ChangeMeta) error {
	return saveJSON(filepath.Join(changeDir(root, workspace, name), "meta.json"), meta)
}

// listChanges returns active (non-archived) changes for a workspace, sorted by name.
// Changes without a readable meta.json get status "unknown".
func listChanges(root, workspace string) []ChangeMeta {
	entries, err := os.ReadDir(filepath.Join(root, workspace, "changes"))
	if err != nil {
		return nil
	}
	var changes []ChangeMeta
	for _, entry := range entries {
		if !entry.IsDir() || entry.Name() == "archive" {
			continue
		}
		meta, err := loadChangeMeta(root, workspace, entry.Name())
		if err != nil {
			meta = ChangeMeta{Name: entry.Name(), Workspace: workspace, Status: "unknown"}
		}
		changes = append(changes, meta)
	}
	sort.Slice(changes, func(i, j int) bool { return changes[i].Name < changes[j].Name })
	return changes
}

func describeChanges(changes []ChangeMeta) string {
	if len(changes) == 0 {
		return "none"
	}
	var parts []string
	for _, change := range changes {
		parts = append(parts, fmt.Sprintf("%s (%s)", change.Name, change.Status))
	}
	return strings.Join(parts, ", ")
}

func archivedNameCollision(root, workspace, name string) bool {
	matches, err := filepath.Glob(filepath.Join(root, workspace, "changes", "archive", "*-"+name))
	return err == nil && len(matches) > 0
}

const changeProposalTemplate = `# Change: %[1]s

Status: proposed
Date: %[2]s
Workspace: %[3]s

## Why

Describe the problem or need this change addresses.

## What Changes

Describe the intended behavior change.

## Foundation Alignment

List the /system files this change relies on. If the change requires new or
different foundation truth, stop and use sysi design-change instead.

## Out Of Scope

List what this change deliberately does not do.
`

const changeDesignTemplate = `# Design: %[1]s

## Decisions

For each significant decision record the decision, the alternatives
considered, and why this one won.

## Interfaces

Describe new or changed interfaces this change introduces inside the
workspace: functions, endpoints, events, schemas.

## Risks

List the main risks and how the tasks mitigate them.
`

const changeTasksTemplate = `# Tasks: %[1]s

Work tasks in order. Check a task only after implementation and verification.

- [ ] 1. Read the /system files listed in proposal.md Foundation Alignment
- [ ] 2. Replace this template with the real task list for the change
`

func (a *App) changePropose(root, workspace, name string, now time.Time) error {
	if slugify(name) != name || name == "" {
		return fmt.Errorf("change name must be a lowercase slug (try %q)", slugify(name))
	}
	if name == "archive" {
		return errors.New(`change name "archive" is reserved`)
	}
	dir := changeDir(root, workspace, name)
	if exists(dir) {
		return fmt.Errorf("change %q already exists in %s", name, workspace)
	}
	if archivedNameCollision(root, workspace, name) {
		return fmt.Errorf("change name %q collides with an archived change in %s", name, workspace)
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	date := now.Format("2006-01-02")
	files := map[string]string{
		"proposal.md": fmt.Sprintf(changeProposalTemplate, name, date, workspace),
		"design.md":   fmt.Sprintf(changeDesignTemplate, name),
		"tasks.md":    fmt.Sprintf(changeTasksTemplate, name),
	}
	for rel, content := range files {
		if err := os.WriteFile(filepath.Join(dir, rel), []byte(content), 0o644); err != nil {
			return err
		}
	}
	stamp := now.Format(time.RFC3339)
	meta := ChangeMeta{Name: name, Workspace: workspace, Status: ChangeStatusProposed, CreatedAt: stamp, UpdatedAt: stamp}
	if err := saveChangeMeta(root, workspace, name, meta); err != nil {
		return err
	}
	rel, err := filepath.Rel(root, dir)
	if err != nil {
		rel = dir
	}
	fmt.Fprintf(a.opts.Stdout, "change proposed: %s\n", name)
	fmt.Fprintf(a.opts.Stdout, "location: %s\n", filepath.ToSlash(rel))
	fmt.Fprintln(a.opts.Stdout, "next: fill proposal.md, design.md, and tasks.md, then run sysi change apply "+name)
	return nil
}
```

In `app.go`, rewrite `change`:

```go
func (a *App) change(args []string) error {
	if len(args) < 2 {
		return errors.New("usage: sysi change propose|apply|archive <name>")
	}
	root, state, err := a.requireProject()
	if err != nil {
		return err
	}
	if state.Phase != PhaseBuild {
		return errors.New("build changes require build phase; run sysi design freeze first")
	}
	workspace, err := currentWorkspace(root, a.opts.Dir, state.Workspaces)
	if err != nil {
		return err
	}

	action, name := args[0], args[1]
	now := time.Now().UTC()
	switch action {
	case "propose":
		return a.changePropose(root, workspace, name, now)
	case "apply":
		return a.changeApply(root, workspace, name)
	case "archive":
		return a.changeArchive(root, workspace, name, now)
	default:
		return fmt.Errorf("unknown change command %q", action)
	}
}
```

For this task, add temporary stubs in `changes.go` so the build compiles (Tasks 4 and 5 replace them):

```go
func (a *App) changeApply(root, workspace, name string) error {
	return errors.New("not implemented")
}

func (a *App) changeArchive(root, workspace, name string, now time.Time) error {
	return errors.New("not implemented")
}
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `GOCACHE=/tmp/sysi-go-cache go test ./...`
Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add -A
git commit -m "native change storage and propose command"
```

---

### Task 4: `sysi change apply`

**Files:**
- Modify: `internal/sysiapp/changes.go` (replace the `changeApply` stub)
- Modify: `internal/sysiapp/changes_test.go`

- [ ] **Step 1: Write the failing tests**

Add to `internal/sysiapp/changes_test.go`:

```go
func TestChangeApplyMarksApplyingAndPrintsHandoff(t *testing.T) {
	root := initBuildProject(t, "api")
	apiDir := filepath.Join(root, "api")
	if code, _, _ := runApp(t, apiDir, "change", "propose", "add-login"); code != 0 {
		t.Fatal("propose failed")
	}

	code, out, errOut := runApp(t, apiDir, "change", "apply", "add-login")
	if code != 0 {
		t.Fatalf("change apply failed: code=%d stdout=%q stderr=%q", code, out, errOut)
	}
	assertContainsAll(t, "apply handoff", out, []string{
		"proposal.md",
		"design.md",
		"tasks.md",
		"Superpowers",
		"TDD",
		"design drift",
		"sysi design-change",
	})

	var meta ChangeMeta
	if err := json.Unmarshal([]byte(readFile(t, filepath.Join(root, "api", "changes", "add-login", "meta.json"))), &meta); err != nil {
		t.Fatal(err)
	}
	if meta.Status != ChangeStatusApplying {
		t.Fatalf("status = %q, want %q", meta.Status, ChangeStatusApplying)
	}

	// Re-apply is idempotent.
	if code, _, _ := runApp(t, apiDir, "change", "apply", "add-login"); code != 0 {
		t.Fatal("re-apply should succeed")
	}
}

func TestChangeApplyUnknownChangeListsAvailable(t *testing.T) {
	root := initBuildProject(t, "api")
	apiDir := filepath.Join(root, "api")
	if code, _, _ := runApp(t, apiDir, "change", "propose", "add-login"); code != 0 {
		t.Fatal("propose failed")
	}

	code, out, errOut := runApp(t, apiDir, "change", "apply", "ghost")
	if code == 0 {
		t.Fatalf("apply of unknown change should fail: stdout=%q stderr=%q", out, errOut)
	}
	assertContainsAll(t, "unknown change error", out+errOut, []string{"ghost", "add-login", "proposed"})
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `GOCACHE=/tmp/sysi-go-cache go test ./internal/sysiapp/ -run TestChangeApply -v`
Expected: FAIL (`not implemented`).

- [ ] **Step 3: Implement**

Replace the `changeApply` stub in `changes.go`:

```go
func (a *App) changeApply(root, workspace, name string) error {
	meta, err := loadChangeMeta(root, workspace, name)
	if err != nil {
		return fmt.Errorf("change %q not found in %s; available: %s", name, workspace, describeChanges(listChanges(root, workspace)))
	}
	if meta.Status == ChangeStatusArchived {
		return fmt.Errorf("change %q is archived", name)
	}
	meta.Status = ChangeStatusApplying
	meta.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	if err := saveChangeMeta(root, workspace, name, meta); err != nil {
		return err
	}

	rel := filepath.ToSlash(filepath.Join(workspace, "changes", name))
	fmt.Fprintf(a.opts.Stdout, "SYSI CHANGE APPLY: %s\n\n", name)
	fmt.Fprintf(a.opts.Stdout, "Workspace: %s\nStatus: %s\nLocation: %s\n\n", workspace, meta.Status, rel)
	fmt.Fprintln(a.opts.Stdout, "Read proposal.md, design.md, and tasks.md before editing implementation code.")
	fmt.Fprintln(a.opts.Stdout, "Work tasks in order with Superpowers discipline: planning, TDD, systematic debugging, verification.")
	fmt.Fprintln(a.opts.Stdout, "Check off each task in tasks.md only after implementation and verification.")
	fmt.Fprintln(a.opts.Stdout, "")
	fmt.Fprintln(a.opts.Stdout, "Stop and use sysi design-change if implementation reveals design drift from /system:")
	fmt.Fprintln(a.opts.Stdout, "new or changed endpoints, payload shapes, event contracts, auth/session/permission rules,")
	fmt.Fprintln(a.opts.Stdout, "shared error behavior, schema or data invariants, security invariants, or observability")
	fmt.Fprintln(a.opts.Stdout, "contracts that /system does not represent.")
	return nil
}
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `GOCACHE=/tmp/sysi-go-cache go test ./...`
Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add -A
git commit -m "change apply marks applying and prints discipline handoff"
```

---

### Task 5: `sysi change archive`

**Files:**
- Modify: `internal/sysiapp/changes.go` (replace the `changeArchive` stub)
- Modify: `internal/sysiapp/changes_test.go`

- [ ] **Step 1: Write the failing tests**

Add to `internal/sysiapp/changes_test.go`:

```go
func TestChangeArchiveMovesChangeAndWarnsOnUncheckedTasks(t *testing.T) {
	root := initBuildProject(t, "api")
	apiDir := filepath.Join(root, "api")
	if code, _, _ := runApp(t, apiDir, "change", "propose", "add-login"); code != 0 {
		t.Fatal("propose failed")
	}
	if code, _, _ := runApp(t, apiDir, "change", "apply", "add-login"); code != 0 {
		t.Fatal("apply failed")
	}

	// Archive with unchecked tasks warns but succeeds.
	code, out, errOut := runApp(t, apiDir, "change", "archive", "add-login")
	if code != 0 {
		t.Fatalf("archive failed: code=%d stdout=%q stderr=%q", code, out, errOut)
	}
	if !strings.Contains(out+errOut, "unchecked") {
		t.Fatalf("archive should warn about unchecked tasks: stdout=%q stderr=%q", out, errOut)
	}

	if _, err := os.Stat(filepath.Join(root, "api", "changes", "add-login")); err == nil {
		t.Fatal("active change dir should be gone after archive")
	}
	matches, err := filepath.Glob(filepath.Join(root, "api", "changes", "archive", "*-add-login"))
	if err != nil || len(matches) != 1 {
		t.Fatalf("expected one archived change dir, got %v (err=%v)", matches, err)
	}
	var meta ChangeMeta
	if err := json.Unmarshal([]byte(readFile(t, filepath.Join(matches[0], "meta.json"))), &meta); err != nil {
		t.Fatal(err)
	}
	if meta.Status != ChangeStatusArchived {
		t.Fatalf("archived status = %q, want %q", meta.Status, ChangeStatusArchived)
	}

	// Archiving an unknown change lists available ones.
	if code, out, errOut := runApp(t, apiDir, "change", "archive", "ghost"); code == 0 {
		t.Fatalf("archive of unknown change should fail: stdout=%q stderr=%q", out, errOut)
	}
}

func TestChangeArchiveNoWarningWhenTasksComplete(t *testing.T) {
	root := initBuildProject(t, "api")
	apiDir := filepath.Join(root, "api")
	if code, _, _ := runApp(t, apiDir, "change", "propose", "add-login"); code != 0 {
		t.Fatal("propose failed")
	}
	tasksPath := filepath.Join(root, "api", "changes", "add-login", "tasks.md")
	done := "# Tasks: add-login\n\n- [x] 1. Everything done\n"
	if err := os.WriteFile(tasksPath, []byte(done), 0o644); err != nil {
		t.Fatal(err)
	}

	code, out, errOut := runApp(t, apiDir, "change", "archive", "add-login")
	if code != 0 {
		t.Fatalf("archive failed: code=%d stdout=%q stderr=%q", code, out, errOut)
	}
	if strings.Contains(out+errOut, "unchecked") {
		t.Fatalf("archive should not warn when tasks are complete: stdout=%q stderr=%q", out, errOut)
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `GOCACHE=/tmp/sysi-go-cache go test ./internal/sysiapp/ -run TestChangeArchive -v`
Expected: FAIL (`not implemented`).

- [ ] **Step 3: Implement**

Replace the `changeArchive` stub in `changes.go`:

```go
func (a *App) changeArchive(root, workspace, name string, now time.Time) error {
	meta, err := loadChangeMeta(root, workspace, name)
	if err != nil {
		return fmt.Errorf("change %q not found in %s; available: %s", name, workspace, describeChanges(listChanges(root, workspace)))
	}

	tasksPath := filepath.Join(changeDir(root, workspace, name), "tasks.md")
	if data, err := os.ReadFile(tasksPath); err == nil {
		unchecked := strings.Count(string(data), "- [ ]")
		if unchecked > 0 {
			fmt.Fprintf(a.opts.Stdout, "warning: tasks.md has %d unchecked task(s)\n", unchecked)
		}
	}

	archiveDir := filepath.Join(root, workspace, "changes", "archive")
	if err := os.MkdirAll(archiveDir, 0o755); err != nil {
		return err
	}
	target := filepath.Join(archiveDir, now.Format("2006-01-02")+"-"+name)
	if exists(target) {
		return fmt.Errorf("archive target already exists: %s", target)
	}
	if err := os.Rename(changeDir(root, workspace, name), target); err != nil {
		return err
	}

	meta.Status = ChangeStatusArchived
	meta.UpdatedAt = now.Format(time.RFC3339)
	if err := saveJSON(filepath.Join(target, "meta.json"), meta); err != nil {
		return err
	}
	rel, err := filepath.Rel(root, target)
	if err != nil {
		rel = target
	}
	fmt.Fprintf(a.opts.Stdout, "change archived: %s -> %s\n", name, filepath.ToSlash(rel))
	return nil
}
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `GOCACHE=/tmp/sysi-go-cache go test ./...`
Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add -A
git commit -m "change archive moves changes and warns on unchecked tasks"
```

---

### Task 6: Status dashboard v2

**Files:**
- Modify: `internal/sysiapp/app.go` (status types, `buildStatus`, `renderStatus`; delete `OpenSpecStatus`, `OpenSpecWorkspaceStatus`, `openSpecStatus`, `openSpecWorkspaceStatus`)
- Modify: `internal/sysiapp/app_test.go`

- [ ] **Step 1: Write the failing tests**

Add to `internal/sysiapp/app_test.go`:

```go
func TestStatusShowsWorkspacesAndNativeChanges(t *testing.T) {
	root := initProject(t, "api,web")
	if code, _, _ := runApp(t, root, "design", "freeze"); code != 0 {
		t.Fatal("freeze failed")
	}
	apiDir := filepath.Join(root, "api")
	if code, _, _ := runApp(t, apiDir, "change", "propose", "add-login"); code != 0 {
		t.Fatal("propose failed")
	}

	code, out, errOut := runApp(t, root, "status", "--json")
	if code != 0 {
		t.Fatalf("status json failed: code=%d stdout=%q stderr=%q", code, out, errOut)
	}
	var status Status
	if err := json.Unmarshal([]byte(out), &status); err != nil {
		t.Fatalf("status output is not json: %v\n%s", err, out)
	}
	if len(status.Workspaces) != 2 {
		t.Fatalf("workspace count = %d, want 2", len(status.Workspaces))
	}
	byName := map[string]WorkspaceStatus{}
	for _, ws := range status.Workspaces {
		byName[ws.Name] = ws
	}
	api := byName["api"]
	if api.ActiveChanges != 1 || len(api.Changes) != 1 || api.Changes[0].Name != "add-login" || api.Changes[0].Status != ChangeStatusProposed {
		t.Fatalf("unexpected api workspace status: %+v", api)
	}
	if byName["web"].ActiveChanges != 0 {
		t.Fatalf("web should have no changes: %+v", byName["web"])
	}
	if strings.Contains(out, "openspec") {
		t.Fatalf("status JSON should not mention openspec:\n%s", out)
	}

	// Human dashboard shows workspaces and change statuses.
	code, out, errOut = runApp(t, root, "status")
	if code != 0 {
		t.Fatalf("status failed: code=%d stdout=%q stderr=%q", code, out, errOut)
	}
	assertContainsAll(t, "human status", out, []string{"Workspaces:", "api", "add-login", "proposed", "web"})
}
```

Also update `TestRoleInferenceUsesDeclaredWorkspaces` and `TestRootDiscoveryAndStatusJSON` if they reference `status.OpenSpec` (they should not after Task 1; verify).

- [ ] **Step 2: Run tests to verify they fail**

Run: `GOCACHE=/tmp/sysi-go-cache go test ./internal/sysiapp/ -run TestStatusShowsWorkspaces -v`
Expected: FAIL (compile error: `Status` has no field `Workspaces` / `WorkspaceStatus` undefined).

- [ ] **Step 3: Implement**

In `app.go`:

1. Delete types `OpenSpecStatus` and `OpenSpecWorkspaceStatus`; delete functions `openSpecStatus` and `openSpecWorkspaceStatus`.

2. Add:

```go
type ChangeSummary struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

type WorkspaceStatus struct {
	Name          string          `json:"name"`
	Present       bool            `json:"present"`
	ActiveChanges int             `json:"activeChanges"`
	Changes       []ChangeSummary `json:"changes"`
}
```

3. `Status`: replace `OpenSpec OpenSpecStatus \`json:"openspec"\`` with `Workspaces []WorkspaceStatus \`json:"workspaces"\``.

4. Add:

```go
func workspacesStatus(root string, workspaces []string) []WorkspaceStatus {
	var statuses []WorkspaceStatus
	for _, ws := range workspaces {
		status := WorkspaceStatus{
			Name:    ws,
			Present: exists(filepath.Join(root, ws)),
			Changes: []ChangeSummary{},
		}
		for _, change := range listChanges(root, ws) {
			status.Changes = append(status.Changes, ChangeSummary{Name: change.Name, Status: change.Status})
		}
		status.ActiveChanges = len(status.Changes)
		statuses = append(statuses, status)
	}
	return statuses
}
```

5. `buildStatus`: replace `OpenSpec: openSpecStatus(root)` with `Workspaces: workspacesStatus(root, state.Workspaces)`.

6. `renderStatus`: replace the two OpenSpec lines (`OpenSpec changes: ...` and the workspace loop) with:

```go
	fmt.Fprintln(a.opts.Stdout, "Workspaces:")
	for _, workspace := range status.Workspaces {
		fmt.Fprintf(a.opts.Stdout, "  - %s: present=%t changes=%d\n", workspace.Name, workspace.Present, workspace.ActiveChanges)
		for _, change := range workspace.Changes {
			fmt.Fprintf(a.opts.Stdout, "      %s: %s\n", change.Name, change.Status)
		}
	}
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `GOCACHE=/tmp/sysi-go-cache go test ./...`
Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add -A
git commit -m "status dashboard reports native workspace changes"
```

---

### Task 7: Validation v2

**Files:**
- Modify: `internal/sysiapp/app.go` (`validateSystem`)
- Modify: `internal/sysiapp/app_test.go`

- [ ] **Step 1: Write the failing tests**

Add to `internal/sysiapp/app_test.go`:

```go
func TestValidateReportsWorkspaceAndChangeProblems(t *testing.T) {
	root := initProject(t, "api,web")
	if code, _, _ := runApp(t, root, "design", "freeze"); code != 0 {
		t.Fatal("freeze failed")
	}
	apiDir := filepath.Join(root, "api")
	if code, _, _ := runApp(t, apiDir, "change", "propose", "add-login"); code != 0 {
		t.Fatal("propose failed")
	}

	// Break things: remove a workspace dir, corrupt a meta.json, add a bogus status.
	if err := os.RemoveAll(filepath.Join(root, "web")); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "api", "changes", "add-login", "meta.json"), []byte("{broken"), 0o644); err != nil {
		t.Fatal(err)
	}
	badDir := filepath.Join(root, "api", "changes", "bad-status")
	if err := os.MkdirAll(badDir, 0o755); err != nil {
		t.Fatal(err)
	}
	badMeta := `{"name":"bad-status","workspace":"api","status":"bogus","createdAt":"x","updatedAt":"x"}`
	if err := os.WriteFile(filepath.Join(badDir, "meta.json"), []byte(badMeta), 0o644); err != nil {
		t.Fatal(err)
	}

	code, out, errOut := runApp(t, root, "validate")
	if code == 0 {
		t.Fatalf("validate should fail: stdout=%q stderr=%q", out, errOut)
	}
	assertContainsAll(t, "validate warnings", out+errOut, []string{
		"missing workspace directory: web",
		"api/changes/add-login",
		"bad-status",
	})
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `GOCACHE=/tmp/sysi-go-cache go test ./internal/sysiapp/ -run TestValidateReportsWorkspace -v`
Expected: FAIL (no such warnings yet).

- [ ] **Step 3: Implement**

In `validateSystem` in `app.go`, after the required-files loop, add:

```go
	for _, ws := range state.Workspaces {
		if !exists(filepath.Join(root, ws)) {
			warnings = append(warnings, fmt.Sprintf("missing workspace directory: %s", ws))
			continue
		}
		entries, err := os.ReadDir(filepath.Join(root, ws, "changes"))
		if err != nil {
			continue
		}
		for _, entry := range entries {
			if !entry.IsDir() || entry.Name() == "archive" {
				continue
			}
			rel := filepath.ToSlash(filepath.Join(ws, "changes", entry.Name()))
			var meta ChangeMeta
			if err := loadJSON(filepath.Join(root, ws, "changes", entry.Name(), "meta.json"), &meta); err != nil {
				warnings = append(warnings, fmt.Sprintf("change %s has missing or invalid meta.json", rel))
				continue
			}
			if meta.Status != ChangeStatusProposed && meta.Status != ChangeStatusApplying {
				warnings = append(warnings, fmt.Sprintf("change %s has invalid status %q", rel, meta.Status))
			}
			if archivedNameCollision(root, ws, entry.Name()) {
				warnings = append(warnings, fmt.Sprintf("change %s collides with an archived change name", rel))
			}
		}
	}
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `GOCACHE=/tmp/sysi-go-cache go test ./...`
Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add -A
git commit -m "validation checks workspaces and native change metadata"
```

---

### Task 8: Delete dead OpenSpec code

**Files:**
- Modify: `internal/sysiapp/app.go`

- [ ] **Step 1: Delete the dead code**

Remove from `app.go`:
- `Options.OpenSpecPath` field
- `implementationOpenSpecTargets` var
- `requireImplementationOpenSpecDir` method
- `ensureImplementationOpenSpec` method
- `runOpenSpec` method (and the `SYSI_OPENSPEC` env lookup inside it)
- Now-unused imports: `"context"`, `"os/exec"` (verify with the compiler)

- [ ] **Step 2: Verify nothing references OpenSpec in Go code**

Run: `grep -ri openspec internal/ cmd/ --include='*.go'`
Expected: no output.

Run: `GOCACHE=/tmp/sysi-go-cache go test ./... && go vet ./...`
Expected: PASS, no vet errors.

- [ ] **Step 3: Commit**

```bash
git add -A
git commit -m "remove openspec integration code"
```

---

### Task 9: Rewrite the `sysi-apply` skill template

**Files:**
- Modify: `internal/sysiapp/templates/agents/codex/sysi-apply/SKILL.md` (full replacement)
- Modify: `internal/sysiapp/app_test.go` (`TestCodexInstructionPacksContainOperationalGuardrails` markers for `sysi-apply`)

- [ ] **Step 1: Update the test markers**

In `TestCodexInstructionPacksContainOperationalGuardrails`, replace the `"sysi-apply"` entry of `skillSpecific` with:

```go
		"sysi-apply": {
			"proposal.md",
			"design.md",
			"tasks.md",
			"sysi change apply",
			"Superpowers",
			"mandatory",
			"missing prerequisite",
			"frozen /system files",
			"new or changed HTTP endpoints",
			"user confirmation",
			"sysi design-change",
			"sysi-design-change",
			"does not agree",
			"system/security/**",
			"declared workspace",
		},
```

- [ ] **Step 2: Run test to verify it fails**

Run: `GOCACHE=/tmp/sysi-go-cache go test ./internal/sysiapp/ -run TestCodexInstructionPacks -v`
Expected: FAIL (markers like `proposal.md` and `declared workspace` missing).

- [ ] **Step 3: Replace the template**

Replace the full contents of `internal/sysiapp/templates/agents/codex/sysi-apply/SKILL.md` with:

```markdown
---
name: sysi-apply
description: Apply a sysi change in build phase using the native change workflow and Superpowers discipline.
---

## Purpose

Use this skill during build phase to implement a sysi change while preserving `/system` as the foundation truth. The change's own files are the work order: `proposal.md` says why, `design.md` says how, `tasks.md` says what remains. Superpowers governs the implementation/debug/test/verify loop.

## Initial Checks

1. Run or read `sysi status --json`.
2. Confirm the project is in build phase.
3. Confirm the current directory is inside a declared workspace and the named change exists under `<workspace>/changes/<name>/`.
4. Run `sysi change apply <name>` from the workspace to mark the change applying and print the handoff.
5. Confirm the relevant Superpowers workflows are available for implementation planning, TDD, systematic debugging, and verification.
6. Read the change's `proposal.md`, `design.md`, and `tasks.md` in full.
7. Read the relevant `/system` files allowed for the current role before editing implementation code, including `system/security/**` when security invariants affect the work.
8. Identify whether the requested implementation would introduce design drift from `/system` before changing behavior.

## Phase Rules

- Build phase is required for implementation.
- Design phase work should use `sysi-explore` and `sysi-capture` instead.
- The change's `tasks.md` owns task tracking during build.
- Running `sysi change apply <name>` is mandatory before implementation edits.
- Superpowers discipline is mandatory during implementation planning, test-driven development, systematic debugging, and verification.
- Frozen /system files are not implementation files.

## Role And File Access

- Role is the declared workspace containing the current working directory.
- Read the allowed `/system` files for that role before deciding how to implement.
- Treat `system/contracts/`, `system/flows/`, `system/modules/<workspace>.md`, `system/data/`, `system/obs/`, and `system/security/**` as build context when they affect the work.
- Keep implementation edits inside the current workspace.

## Workflow

1. Run `sysi change apply <name>` before editing implementation code.
2. Use Superpowers skills for implementation planning, TDD, debugging, and verification.
3. Treat a missing Superpowers workflow as a missing prerequisite and stop instead of implementing without it.
4. Work through `tasks.md` in order and check each task off only after implementation and verification.
5. Keep edits scoped to the change and the current task.
6. Compare implementation needs against `/system` truth before changing behavior.
7. Treat design drift as any implementation need that contradicts or extends foundation truth, including new or changed HTTP endpoints, request or response payload-shape changes, event contracts, auth/session/permission rules, shared error behavior, contract conventions, schema or data invariants, security invariants, metrics, logging, tracing, or alerting contracts.
8. If implementation reveals design drift, stop ordinary implementation work and explain the mismatch: what implementation needs, what `/system` currently says or omits, and which `/system` files likely own the truth.
9. Ask the user for explicit user confirmation before changing `/system` for detected drift.
10. If the user confirms the foundation change, run `sysi design-change <name>` and follow `sysi-design-change` before mutating controlled or frozen `/system` files.
11. If the user does not agree to the foundation change, do not continue implementation that contradicts `/system`; revise the change or implementation approach to fit current foundation truth.
12. When all tasks are checked and verified, run `sysi change archive <name>` from the workspace.

## Validation

- Run focused tests for the changed behavior.
- Run broader tests required by the change before completion.
- Re-read modified code and relevant `/system` files to check alignment.
- Confirm implementation respects contract conventions, error behavior, and security invariants when those files apply.
- Confirm detected design drift received user confirmation before any foundation mutation.
- Confirm agreed design drift went through `sysi design-change <name>` and `sysi-design-change` before controlled or frozen `/system` edits.
- Confirm no frozen /system files changed accidentally.
- Confirm `tasks.md` checkboxes accurately reflect completed work.

## Stop Conditions

- Stop if `sysi status --json` does not show build phase.
- Stop if the current directory is not inside a declared workspace.
- Stop if the named change is missing or archived.
- Stop if required Superpowers workflows are unavailable.
- Stop if the requested implementation contradicts `/system` truth.
- Stop if a foundation mutation is required and the user has not confirmed the drift.
- Stop if the user does not agree to a required foundation change and implementation would contradict `/system`.
- Stop if confirmed design drift has not gone through `sysi design-change <name>` and `sysi-design-change`.
- Stop if tests fail and systematic debugging has not isolated the cause.

## Do Not

- Do Not implement outside a sysi change during build phase.
- Do Not implement before running `sysi change apply <name>`.
- Do Not implement when a mandatory apply/debug/test/verify workflow is a missing prerequisite.
- Do Not mutate frozen /system files as part of normal apply.
- Do Not mutate `/system` for design drift without explicit user confirmation.
- Do Not continue implementation that contradicts `/system` when the user does not agree to the foundation change.
- Do Not treat new endpoints, payload shapes, auth rules, security invariants, data shapes, or observability contracts as ordinary implementation details when they are missing from `/system`.
- Do Not copy full Superpowers instructions into this skill; invoke or follow them.
- Do Not mark tasks complete without fresh verification.
- Do Not hide design drift by forcing code to fit an outdated proposal.
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `GOCACHE=/tmp/sysi-go-cache go test ./...`
Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add -A
git commit -m "rewrite sysi-apply skill for native change workflow"
```

---

### Task 10: Update remaining templates and CLI guidance text

**Files:**
- Modify: `internal/sysiapp/templates/agents/codex/sysi-explore/SKILL.md`
- Modify: `internal/sysiapp/templates/agents/codex/sysi-capture/SKILL.md`
- Modify: `internal/sysiapp/templates/agents/codex/sysi-design-change/SKILL.md`
- Modify: `internal/sysiapp/templates/agents/cursor/sysi.mdc`
- Modify: `internal/sysiapp/templates/agents/claude/CLAUDE.section.md`
- Modify: `internal/sysiapp/app.go` (`explore` output, `createDesignChangeArtifact`, `designChange` output)
- Modify: `internal/sysiapp/app_test.go` (marker updates)

- [ ] **Step 1: Update the test markers**

- In `TestCodexInstructionPacksContainOperationalGuardrails`, `"sysi-design-change"` markers: replace `"impacted OpenSpec changes"` with `"impacted workspace changes"`.
- In `TestCursorAndClaudeInstructionsContainWorkflowBoundaries`, `markers`: replace `"OpenSpec"` with `"sysi change"`.
- In `TestDesignChangeCreatesDatedDecisionArtifact`, replace marker `"## Impacted OpenSpec Changes"` with `"## Impacted Changes"`.
- Add to the end of `TestDesignCommandsDoNotRequireOpenSpec` (rename it to `TestDesignGuidanceDoesNotMentionOpenSpec`):

```go
	if strings.Contains(out, "OpenSpec") {
		t.Fatalf("design guidance should not mention OpenSpec: %q", out)
	}
```

(where `out` is the explore output; keep the existing assertions).

- [ ] **Step 2: Run tests to verify they fail**

Run: `GOCACHE=/tmp/sysi-go-cache go test ./internal/sysiapp/ -run 'TestCodexInstructionPacks|TestCursorAndClaude|TestDesignChangeCreates|TestDesignGuidance' -v`
Expected: FAIL.

- [ ] **Step 3: Implement the text changes**

**`app.go` `explore`**: replace the line

```go
	fmt.Fprintln(a.opts.Stdout, "During design phase, do not create OpenSpec changes.")
```

with

```go
	fmt.Fprintln(a.opts.Stdout, "During design phase, do not create build changes; capture decisions into /system.")
```

**`app.go` `createDesignChangeArtifact`**: in the artifact template, replace

```
## Impacted OpenSpec Changes

- frontend: TBD
- backend: TBD
```

with

```
## Impacted Changes

- List impacted workspace changes as workspace: change-name.
```

Also replace the sentence `Describe why normal OpenSpec apply work cannot continue without changing foundation truth.` with `Describe why normal build change work cannot continue without changing foundation truth.`

**`app.go` `designChange`**: replace the final `Fprintln` line's text with:

```go
	fmt.Fprintln(a.opts.Stdout, "Record rationale, affected /system files, impacted workspace changes, and migration notes before mutating frozen foundation files.")
```

**`sysi-explore/SKILL.md`** targeted replacements:
- Line 3 description: `... Does not create OpenSpec changes during design phase.` → `... Does not create build changes during design phase.`
- `- During design phase, do not create OpenSpec changes for design exploration.` → `- During design phase, do not create build changes for design exploration.`
- `- Do Not skip design work; avoid implementation until a build-phase OpenSpec change exists.` → `- Do Not skip design work; avoid implementation until a build-phase sysi change exists.`
- `- Do Not create OpenSpec changes during design phase.` → `- Do Not create build changes during design phase.`

**`sysi-capture/SKILL.md`**:
- `- Do not use OpenSpec for design-phase capture.` → `- Do not use build changes for design-phase capture.`
- `- Do Not create OpenSpec changes during design phase.` → `- Do Not create build changes during design phase.`

**`sysi-design-change/SKILL.md`** replacements (every OpenSpec mention):
- `4. Identify the current frontend or backend OpenSpec change, if any.` → `4. Identify the current workspace change(s), if any.`
- `...and every impacted OpenSpec change.` → `...and every impacted workspace change.`
- In the routing table row: `impacted OpenSpec changes` → `impacted workspace changes`
- `- If OpenSpec artifacts conflict with the new foundation truth, update or pause those artifacts before resuming implementation.` → `- If workspace change artifacts conflict with the new foundation truth, update or pause those changes before resuming implementation.`
- `5. List impacted OpenSpec changes and implementation tasks in the decision artifact.` → `5. List impacted workspace changes and implementation tasks in the decision artifact.`
- `- Confirm the decision artifact records rationale, affected files, impacted OpenSpec changes, ...` → `... impacted workspace changes, ...`
- `- Confirm impacted OpenSpec changes still describe the intended implementation.` → `- Confirm impacted workspace changes still describe the intended implementation.`
- `- Stop if impacted OpenSpec changes are unknown and the mutation would change implementation scope.` → `- Stop if impacted workspace changes are unknown and the mutation would change implementation scope.`
- `- Do Not leave OpenSpec artifacts inconsistent with changed /system truth.` → `- Do Not leave workspace changes inconsistent with changed /system truth.`

**`cursor/sysi.mdc`** and **`claude/CLAUDE.section.md`** (same edits in both):
- `- Respect phase boundaries: design phase captures decisions into /system; build phase implements OpenSpec changes.` → `- Respect phase boundaries: design phase captures decisions into /system; build phase implements sysi changes inside declared workspaces.`
- `- During design phase, use sysi explore and sysi capture semantics; do not create OpenSpec changes for design decisions.` → `- During design phase, use sysi explore and sysi capture semantics; do not create build changes for design decisions.`
- `- During build phase, implementation work must flow through OpenSpec and the local apply workflow.` → `- During build phase, implementation work must flow through sysi change propose|apply|archive from the owning workspace.`
- `- These rules do not replace OpenSpec, Superpowers, tests, or explicit user confirmation.` → `- These rules do not replace Superpowers, tests, or explicit user confirmation.`

- [ ] **Step 4: Run tests to verify they pass, and grep for leftovers**

Run: `GOCACHE=/tmp/sysi-go-cache go test ./...`
Expected: PASS.

Run: `grep -ri openspec internal/sysiapp/templates/`
Expected: no output.

- [ ] **Step 5: Commit**

```bash
git add -A
git commit -m "purge openspec from skill templates and cli guidance"
```

---

### Task 11: README rewrite and repo spec updates

**Files:**
- Modify: `README.md` (full rewrite of affected sections)
- Modify: `openspec/specs/project-lifecycle/spec.md`
- Modify: `openspec/specs/build-workflow/spec.md`
- Modify: `openspec/specs/status-dashboard/spec.md`
- Modify: `openspec/specs/design-workflow/spec.md`
- Modify: `openspec/specs/agent-integration/spec.md`
- Modify: `openspec/specs/system-foundation/spec.md`
- Modify: `openspec/specs/project-documentation/spec.md`

These files describe sysi itself. Rewrite them so every requirement describing OpenSpec integration, hardcoded frontend/backend, or the v1 shape describes the v2 native behavior instead. The design spec (`docs/superpowers/specs/2026-07-07-sysi-v2-native-core-design.md`) is the source of truth.

- [ ] **Step 1: Rewrite README.md**

Keep the existing tone and structure. Required content changes:

- The intro summary block becomes:

```text
/system      = ratified system truth
sysi changes = build-phase change protocol (native)
Superpowers  = apply-phase engineering discipline
sysi          = CLI that ties those pieces together
```

- **Mental Model**: replace "OpenSpec build changes" with "sysi workspace changes". The three-boundary table's OpenSpec row becomes: `sysi changes | Build-phase changes inside declared workspaces`.
- **Project Lifecycle**: design flow unchanged; build flow becomes `sysi change propose add-login` / `sysi change apply add-login` / `sysi change archive add-login` run from a workspace directory.
- **Quick Start**: `sysi init --workspaces frontend,backend` (note that init fails with guidance when the flag is missing); everything OpenSpec-related removed.
- **Repository Layout**: replace `frontend/openspec/` and `backend/openspec/` with `<workspace>/changes/` (+ `changes/archive/`); state.json documented as version 2 with `workspaces`.
- **The `/system` Foundation**: `system/modules/<workspace>.md` per declared workspace instead of hardcoded frontend.md/backend.md.
- **Design Phase**: role table becomes: repo root → `design`, `system/` → `system-maintainer`, inside declared workspace `<name>` → `<name>`.
- **Build Phase**: full rewrite describing native changes — the three files scaffolded by propose, statuses (`proposed`, `applying`, `archived`), apply handoff, archive warning on unchecked tasks. Remove all `SYSI_OPENSPEC` mentions.
- **New section: Workspaces** documenting `sysi workspace list|add|remove <name> [--force]`.
- **Agent Integrations**: unchanged structurally; reword skill descriptions to native changes.
- **Status And Validation**: JSON structure now has `workspaces` (with per-change statuses) instead of `openspec`; validation checks workspace dirs and change metadata.
- **Command Reference**: update `sysi init`, add `sysi workspace`, rewrite the three `sysi change` entries (no OpenSpec invocation).
- **Troubleshooting**: delete `openspec executable not found` and `OpenSpec PostHog Network Errors` sections; add `sysi init requires --workspaces` and `change commands must run inside a declared workspace` entries.
- **Contributor Notes**: test command stays `GOCACHE=/tmp/sysi-go-cache go test ./...`; drop `openspec validate --specs`.
- **V1 Boundaries** section becomes **V2 Boundaries**: remove "replacement behavior for OpenSpec" (it IS replaced now); keep no-sandboxing, no generated views, minimal Cursor/Claude, no curses UI; add "no multi-phase plan workflow yet (M2)" and "no gap analysis yet (M3)".

- [ ] **Step 2: Update the openspec spec files**

For each file, rewrite requirements/scenarios that mention OpenSpec or hardcoded frontend/backend:

- `project-lifecycle/spec.md`: init requires `--workspaces`; scaffolds declared workspaces with `changes/` dirs; no OpenSpec initialization; bare init on uninitialized repo prints usage and fails; re-init is idempotent; state is version 2 with workspaces; add workspace add/remove/list requirements.
- `build-workflow/spec.md`: propose/apply/archive are native; must run from a declared workspace; statuses proposed/applying/archived; archive moves to dated archive dir and warns on unchecked tasks; apply prints Superpowers handoff and drift stop conditions.
- `status-dashboard/spec.md`: dashboard and JSON report per-workspace native changes with statuses; no OpenSpec aggregation.
- `design-workflow/spec.md`: replace "does not call openspec" phrasing with "does not create build changes"; capture scenarios unchanged otherwise.
- `agent-integration/spec.md`: skills reference native change workflow; remove openspec-apply-change requirements.
- `system-foundation/spec.md`: modules files are scaffolded per declared workspace.
- `project-documentation/spec.md`: README documents v2 behavior (workspaces flag, native changes, no OpenSpec dependency).

- [ ] **Step 3: Verify**

Run: `grep -rn "openspec\|OpenSpec" README.md | grep -vi "native\|replaced\|no longer\|absorb"` — review remaining mentions; only historical/explanatory mentions are acceptable (ideally none).
Run: `GOCACHE=/tmp/sysi-go-cache go test ./...`
Expected: PASS.

- [ ] **Step 4: Commit**

```bash
git add -A
git commit -m "document v2 native workflow in readme and specs"
```

---

### Task 12: Final verification and smoke test

**Files:** none (verification only; fix anything found)

- [ ] **Step 1: Full checks**

```bash
GOCACHE=/tmp/sysi-go-cache go test ./...
go vet ./...
gofmt -l cmd internal
go build -o /tmp/sysi-smoke/sysi ./cmd/sysi
```

Expected: tests pass, no vet output, no gofmt output, build succeeds.

- [ ] **Step 2: End-to-end smoke test in a temp directory**

```bash
mkdir -p /tmp/sysi-smoke/demo && cd /tmp/sysi-smoke/demo
/tmp/sysi-smoke/sysi init --workspaces api,web
/tmp/sysi-smoke/sysi status
/tmp/sysi-smoke/sysi workspace add worker
/tmp/sysi-smoke/sysi workspace list
/tmp/sysi-smoke/sysi design freeze
cd api
/tmp/sysi-smoke/sysi change propose add-login
/tmp/sysi-smoke/sysi change apply add-login
/tmp/sysi-smoke/sysi change archive add-login
cd ..
/tmp/sysi-smoke/sysi status --json
/tmp/sysi-smoke/sysi validate
```

Expected: every command exits 0 (archive prints the unchecked-tasks warning); status JSON shows three workspaces and no `openspec` key; validate passes.

- [ ] **Step 3: Check the old binary artifact**

The repo root contains a stale compiled `sysi` binary (3.6MB, committed earlier). Rebuild it so it matches v2, or remove it and add `sysi` to `.gitignore` (preferred — compiled artifacts don't belong in git):

```bash
cd /Users/mmdbasi/other_projects/sysi
git rm --cached sysi
printf '%s\n' sysi >> .gitignore
```

- [ ] **Step 4: Commit any fixes**

```bash
git add -A
git commit -m "v2 verification fixes and drop committed binary"
```

---

## Spec Coverage Check

| Spec section | Tasks |
| --- | --- |
| Workspaces declared at init, stored in state v2 | Task 1 |
| Bare init prints usage, no guessing | Task 1 |
| Role inference generalizes | Task 1 |
| Generic allowlists | Task 1 |
| OpenSpec removal | Tasks 1, 6, 8 |
| Native change storage + templates | Task 3 |
| Workspace management commands | Task 2 |
| Change lifecycle propose/apply/archive | Tasks 3, 4, 5 |
| Status dashboard per-workspace changes | Task 6 |
| Validation v2 | Task 7 |
| Error messages teach | Tasks 2, 3, 4, 5 |
| Skill packs rewritten | Tasks 9, 10 |
| Cursor/Claude instructions | Task 10 |
| README + repo specs | Task 11 |
| Testing without external binary | All tasks; verified in Task 12 |
