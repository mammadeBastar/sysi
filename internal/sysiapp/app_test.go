package sysiapp

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func runApp(t *testing.T, dir string, args ...string) (int, string, string) {
	t.Helper()

	var stdout, stderr bytes.Buffer
	code := New(Options{
		Dir:    dir,
		Stdout: &stdout,
		Stderr: &stderr,
	}).Run(args)

	return code, stdout.String(), stderr.String()
}

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

	// Uses the --workspaces=<value> equals form; other tests cover the space-separated form.
	code, out, errOut := runApp(t, root, "init", "--workspaces=api,web")
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
	for _, invalid := range []string{"system", "design", "system-maintainer", "Api", "a b", "-api", ""} {
		t.Run(invalid, func(t *testing.T) {
			root := t.TempDir()
			code, out, errOut := runApp(t, root, "init", "--workspaces", "good,"+invalid)
			if code == 0 {
				t.Fatalf("init should reject workspace name %q: stdout=%q stderr=%q", invalid, out, errOut)
			}
			if invalid != "" && !strings.Contains(out+errOut, "\""+invalid+"\"") {
				t.Fatalf("error should mention offending name %q: stdout=%q stderr=%q", invalid, out, errOut)
			}
			if _, err := os.Stat(filepath.Join(root, ".sysi")); err == nil {
				t.Fatalf("init must not create .sysi when a workspace name is invalid")
			}
		})
	}
	root := t.TempDir()
	code, out, errOut := runApp(t, root, "init", "--workspaces", "api,api")
	if code == 0 {
		t.Fatalf("init should reject duplicate workspace names: stdout=%q stderr=%q", out, errOut)
	}
	if _, err := os.Stat(filepath.Join(root, ".sysi")); err == nil {
		t.Fatalf("init must not create .sysi when workspace names are duplicated")
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
		root:                          RoleDesign,
		apiDir:                        "api",
		filepath.Join(root, "system"): RoleSystem,
		filepath.Join(root, "web"):    "web",
	}
	for dir, wantRole := range cases {
		code, out, errOut := runApp(t, dir, "status", "--json")
		if code != 0 {
			t.Fatalf("status in %s failed: code=%d stdout=%q stderr=%q", dir, code, out, errOut)
		}
		var status Status
		if err := json.Unmarshal([]byte(out), &status); err != nil {
			t.Fatalf("status output is not json: %v\n%s", err, out)
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

func TestLoadStateRejectsInvalidWorkspaceNames(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, ".sysi"), 0o755); err != nil {
		t.Fatal(err)
	}
	bad := `{"version":2,"phase":"design","createdAt":"x","updatedAt":"x","workspaces":["../../escape"]}`
	if err := os.WriteFile(filepath.Join(root, ".sysi", "state.json"), []byte(bad), 0o644); err != nil {
		t.Fatal(err)
	}

	code, out, errOut := runApp(t, root, "status")
	if code == 0 {
		t.Fatalf("status should fail on state with invalid workspace name: stdout=%q stderr=%q", out, errOut)
	}
	if !strings.Contains(out+errOut, "../../escape") && !strings.Contains(out+errOut, "invalid state") {
		t.Fatalf("error should mention the bad workspace name or invalid state: stdout=%q stderr=%q", out, errOut)
	}
}

func TestInitRefusesWorkspaceConflictingWithFile(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "backend"), []byte("not a directory\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	code, out, errOut := runApp(t, root, "init", "--workspaces", "backend")
	if code == 0 {
		t.Fatalf("init should fail when workspace conflicts with existing file: stdout=%q stderr=%q", out, errOut)
	}
	if !strings.Contains(out+errOut, "conflict") {
		t.Fatalf("error should mention the conflict: stdout=%q stderr=%q", out, errOut)
	}
	if _, err := os.Stat(filepath.Join(root, ".sysi")); err == nil {
		t.Fatalf("init must not create .sysi when a workspace conflicts with a file")
	}
}

func TestRootDiscoveryAndStatusJSON(t *testing.T) {
	root := t.TempDir()
	if code, out, errOut := runApp(t, root, "init", "--workspaces", "frontend,backend"); code != 0 {
		t.Fatalf("init failed: code=%d stdout=%q stderr=%q", code, out, errOut)
	}

	frontendDir := filepath.Join(root, "frontend", "app")
	if err := os.MkdirAll(frontendDir, 0o755); err != nil {
		t.Fatal(err)
	}

	code, out, errOut := runApp(t, frontendDir, "status", "--json")
	if code != 0 {
		t.Fatalf("status json failed: code=%d stdout=%q stderr=%q", code, out, errOut)
	}

	var status Status
	if err := json.Unmarshal([]byte(out), &status); err != nil {
		t.Fatalf("status output is not json: %v\n%s", err, out)
	}
	if status.Root != root {
		t.Fatalf("root = %q, want %q", status.Root, root)
	}
	if status.Phase != PhaseDesign {
		t.Fatalf("phase = %q, want %q", status.Phase, PhaseDesign)
	}
	if status.Role != "frontend" {
		t.Fatalf("role = %q, want %q", status.Role, "frontend")
	}

	var allowlists map[string][]string
	if err := json.Unmarshal([]byte(readFile(t, filepath.Join(root, ".sysi", "allowlists.json"))), &allowlists); err != nil {
		t.Fatal(err)
	}
	assertContainsAll(t, "frontend allowlist", strings.Join(allowlists["frontend"], "\n"), []string{"system/security/**"})
	assertContainsAll(t, "backend allowlist", strings.Join(allowlists["backend"], "\n"), []string{"system/security/**"})
}

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
	assertContainsAll(t, "json keys", out, []string{"\"workspaces\"", "\"activeChanges\"", "\"changes\""})
	if !strings.Contains(out, "\"changes\": []") {
		t.Fatalf("empty web workspace should emit \"changes\": []:\n%s", out)
	}

	// Human dashboard shows workspaces and change statuses.
	code, out, errOut = runApp(t, root, "status")
	if code != 0 {
		t.Fatalf("status failed: code=%d stdout=%q stderr=%q", code, out, errOut)
	}
	assertContainsAll(t, "human status", out, []string{"Workspaces:", "api", "add-login", "proposed", "web"})
}

func TestStatusWithZeroWorkspacesEmitsEmptyList(t *testing.T) {
	root := initProject(t, "api")
	if code, out, errOut := runApp(t, root, "workspace", "remove", "api", "--force"); code != 0 {
		t.Fatalf("workspace remove failed: code=%d stdout=%q stderr=%q", code, out, errOut)
	}

	code, out, errOut := runApp(t, root, "status", "--json")
	if code != 0 {
		t.Fatalf("status json failed: code=%d stdout=%q stderr=%q", code, out, errOut)
	}
	if strings.Contains(out, "\"workspaces\": null") {
		t.Fatalf("status JSON should not emit \"workspaces\": null:\n%s", out)
	}
	if !strings.Contains(out, "\"workspaces\": []") {
		t.Fatalf("status JSON should emit \"workspaces\": []:\n%s", out)
	}

	code, out, errOut = runApp(t, root, "status")
	if code != 0 {
		t.Fatalf("status failed: code=%d stdout=%q stderr=%q", code, out, errOut)
	}
	assertContainsAll(t, "human status with no workspaces", out, []string{"Workspaces:", "(none)"})
}

func TestValidateReportsMissingRequiredSystemFile(t *testing.T) {
	for _, rel := range []string{
		"system/contracts/api.yaml",
		"system/contracts/conventions.md",
		"system/contracts/errors.md",
		"system/security/model.md",
	} {
		t.Run(rel, func(t *testing.T) {
			root := t.TempDir()
			if code, out, errOut := runApp(t, root, "init", "--workspaces", "frontend,backend"); code != 0 {
				t.Fatalf("init failed: code=%d stdout=%q stderr=%q", code, out, errOut)
			}
			if err := os.Remove(filepath.Join(root, rel)); err != nil {
				t.Fatal(err)
			}

			code, out, errOut := runApp(t, root, "validate")
			if code == 0 {
				t.Fatalf("validate should fail when required file is missing: stdout=%q stderr=%q", out, errOut)
			}
			if !strings.Contains(out+errOut, rel) {
				t.Fatalf("missing file warning not found in output: stdout=%q stderr=%q", out, errOut)
			}
		})
	}
}

func TestValidateReportsWorkspaceAndChangeProblems(t *testing.T) {
	root := initProject(t, "api,web")
	if code, _, _ := runApp(t, root, "design", "freeze"); code != 0 {
		t.Fatal("freeze failed")
	}
	apiDir := filepath.Join(root, "api")
	if code, _, _ := runApp(t, apiDir, "change", "propose", "add-login"); code != 0 {
		t.Fatal("propose failed")
	}

	// Archive a change, then recreate an active change under the same name so
	// validation must flag the archived-name collision.
	if code, _, _ := runApp(t, apiDir, "change", "propose", "old-change"); code != 0 {
		t.Fatal("propose old-change failed")
	}
	if code, _, _ := runApp(t, apiDir, "change", "apply", "old-change"); code != 0 {
		t.Fatal("apply old-change failed")
	}
	if code, _, _ := runApp(t, apiDir, "change", "archive", "old-change"); code != 0 {
		t.Fatal("archive old-change failed")
	}
	collideDir := filepath.Join(root, "api", "changes", "old-change")
	if err := os.MkdirAll(collideDir, 0o755); err != nil {
		t.Fatal(err)
	}
	collideMeta := `{"name":"old-change","workspace":"api","status":"proposed","createdAt":"x","updatedAt":"x"}`
	if err := os.WriteFile(filepath.Join(collideDir, "meta.json"), []byte(collideMeta), 0o644); err != nil {
		t.Fatal(err)
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
		"has missing or invalid meta.json",
		"api/changes/bad-status",
		`has invalid status "bogus"`,
		"api/changes/old-change",
		"collides with an archived change name",
	})
}

func TestValidateFlagsWorkspacePathThatIsAFile(t *testing.T) {
	root := initProject(t, "api,web")
	if code, _, _ := runApp(t, root, "design", "freeze"); code != 0 {
		t.Fatal("freeze failed")
	}

	// Replace the web workspace directory with a plain file.
	if err := os.RemoveAll(filepath.Join(root, "web")); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "web"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}

	code, out, errOut := runApp(t, root, "validate")
	if code == 0 {
		t.Fatalf("validate should fail when workspace path is a file: stdout=%q stderr=%q", out, errOut)
	}
	assertContainsAll(t, "workspace-as-file warning", out+errOut, []string{
		"missing workspace directory: web",
	})
}

func TestStatusJSONEmitsEmptyWarningLists(t *testing.T) {
	root := initProject(t, "api")

	code, out, errOut := runApp(t, root, "status", "--json")
	if code != 0 {
		t.Fatalf("status json failed: code=%d stdout=%q stderr=%q", code, out, errOut)
	}
	assertContainsAll(t, "empty lists", out, []string{`"warnings": []`, `"mutations": []`})
}

func TestDesignFreezeRecordsBaselineAndCaptureBlocksInBuild(t *testing.T) {
	root := t.TempDir()
	if code, out, errOut := runApp(t, root, "init", "--workspaces", "frontend,backend"); code != 0 {
		t.Fatalf("init failed: code=%d stdout=%q stderr=%q", code, out, errOut)
	}

	code, out, errOut := runApp(t, root, "design", "freeze")
	if code != 0 {
		t.Fatalf("design freeze failed: code=%d stdout=%q stderr=%q", code, out, errOut)
	}

	code, out, errOut = runApp(t, root, "capture")
	if code == 0 {
		t.Fatalf("capture should fail during build phase: stdout=%q stderr=%q", out, errOut)
	}
	if !strings.Contains(out+errOut, "design-change") {
		t.Fatalf("capture output should mention design-change: stdout=%q stderr=%q", out, errOut)
	}

	for _, rel := range []string{
		"system/architecture/system.md",
		"system/contracts/errors.md",
		"system/security/model.md",
	} {
		if err := os.WriteFile(filepath.Join(root, rel), []byte("changed\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	code, out, errOut = runApp(t, root, "status", "--json")
	if code != 0 {
		t.Fatalf("status json failed: code=%d stdout=%q stderr=%q", code, out, errOut)
	}

	var status Status
	if err := json.Unmarshal([]byte(out), &status); err != nil {
		t.Fatal(err)
	}
	if len(status.Validation.Warnings) < 3 {
		t.Fatalf("expected freeze warnings after foundation mutations: %#v", status)
	}
	assertContainsAll(t, "freeze warning status", out, []string{
		"design-change",
		"system/architecture/system.md",
		"system/contracts/errors.md",
		"system/security/model.md",
	})
}

func TestDesignCommandsMentionExpandedFoundationTargets(t *testing.T) {
	root := t.TempDir()
	if code, out, errOut := runApp(t, root, "init", "--workspaces", "frontend,backend"); code != 0 {
		t.Fatalf("init failed: code=%d stdout=%q stderr=%q", code, out, errOut)
	}

	frontendDir := filepath.Join(root, "frontend")
	code, out, errOut := runApp(t, frontendDir, "explore", "security")
	if code != 0 {
		t.Fatalf("frontend explore failed: code=%d stdout=%q stderr=%q", code, out, errOut)
	}
	assertContainsAll(t, "frontend explore output", out, []string{
		"system/contracts/**",
		"system/security/**",
		"contract conventions",
		"contract errors",
	})

	code, out, errOut = runApp(t, root, "capture")
	if code != 0 {
		t.Fatalf("capture failed in design phase: code=%d stdout=%q stderr=%q", code, out, errOut)
	}
	assertContainsAll(t, "capture output", out, []string{
		"contracts",
		"conventions",
		"errors",
		"security",
	})
	if !strings.Contains(out, "decision record") {
		t.Fatalf("capture output should mention decision records: %s", out)
	}
}

func TestDesignChangeCreatesDatedDecisionArtifact(t *testing.T) {
	root := t.TempDir()
	if code, out, errOut := runApp(t, root, "init", "--workspaces", "frontend,backend"); code != 0 {
		t.Fatalf("init failed: code=%d stdout=%q stderr=%q", code, out, errOut)
	}
	if code, out, errOut := runApp(t, root, "design", "freeze"); code != 0 {
		t.Fatalf("design freeze failed: code=%d stdout=%q stderr=%q", code, out, errOut)
	}

	code, out, errOut := runApp(t, root, "design-change", "change auth boundary")
	if code != 0 {
		t.Fatalf("design-change failed: code=%d stdout=%q stderr=%q", code, out, errOut)
	}

	matches, err := filepath.Glob(filepath.Join(root, "system", "architecture", "decisions", "*-change-auth-boundary.md"))
	if err != nil {
		t.Fatal(err)
	}
	if len(matches) != 1 {
		t.Fatalf("expected one dated decision artifact, got %d: %v", len(matches), matches)
	}
	content := readFile(t, matches[0])
	assertContainsAll(t, "design-change artifact", content, []string{
		"# Design Change: change auth boundary",
		"Status: proposed",
		"## Rationale",
		"## Affected System Files",
		"## Impacted OpenSpec Changes",
		"## Migration Or Compatibility Notes",
	})
	if !strings.Contains(out, filepath.Base(matches[0])) {
		t.Fatalf("design-change output should mention artifact path: %q", out)
	}
}

func TestDesignCommandsDoNotRequireOpenSpec(t *testing.T) {
	root := t.TempDir()
	if code, out, errOut := runApp(t, root, "init", "--workspaces", "frontend,backend"); code != 0 {
		t.Fatalf("init failed: code=%d stdout=%q stderr=%q", code, out, errOut)
	}

	code, out, errOut := runApp(t, root, "explore", "auth")
	if code != 0 {
		t.Fatalf("explore failed: code=%d stdout=%q stderr=%q", code, out, errOut)
	}
	if !strings.Contains(out, "auth") || strings.Contains(out, "openspec new") {
		t.Fatalf("explore output did not look like design guidance: %q", out)
	}

	code, out, errOut = runApp(t, root, "capture")
	if code != 0 {
		t.Fatalf("capture failed in design phase: code=%d stdout=%q stderr=%q", code, out, errOut)
	}
	if !strings.Contains(out, "decision record") {
		t.Fatalf("capture output should mention decision records: %q", out)
	}
}

func TestAgentInstallersGenerateExpectedFilesAndPreserveClaudeContent(t *testing.T) {
	root := t.TempDir()
	if code, out, errOut := runApp(t, root, "init", "--workspaces", "frontend,backend"); code != 0 {
		t.Fatalf("init failed: code=%d stdout=%q stderr=%q", code, out, errOut)
	}

	claudePath := filepath.Join(root, "CLAUDE.md")
	if err := os.WriteFile(claudePath, []byte("# Existing\n\nKeep this.\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	for _, agent := range []string{"codex", "cursor", "claude"} {
		code, out, errOut := runApp(t, root, "agent", "install", agent)
		if code != 0 {
			t.Fatalf("agent install %s failed: code=%d stdout=%q stderr=%q", agent, code, out, errOut)
		}
	}

	wantFiles := []string{
		".codex/skills/sysi-explore/SKILL.md",
		".codex/skills/sysi-capture/SKILL.md",
		".codex/skills/sysi-apply/SKILL.md",
		".codex/skills/sysi-design-change/SKILL.md",
		".cursor/rules/sysi.mdc",
		"CLAUDE.md",
	}
	for _, rel := range wantFiles {
		if _, err := os.Stat(filepath.Join(root, rel)); err != nil {
			t.Fatalf("expected %s to exist: %v", rel, err)
		}
	}

	claude, err := os.ReadFile(claudePath)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(claude), "Keep this.") || !strings.Contains(string(claude), "SYSI") {
		t.Fatalf("CLAUDE.md did not preserve content and add marked section:\n%s", claude)
	}
}

func TestCodexInstructionPacksContainOperationalGuardrails(t *testing.T) {
	root := t.TempDir()
	if code, out, errOut := runApp(t, root, "init", "--workspaces", "frontend,backend"); code != 0 {
		t.Fatalf("init failed: code=%d stdout=%q stderr=%q", code, out, errOut)
	}
	if code, out, errOut := runApp(t, root, "agent", "install", "codex"); code != 0 {
		t.Fatalf("codex install failed: code=%d stdout=%q stderr=%q", code, out, errOut)
	}

	commonMarkers := []string{
		"## Purpose",
		"## Initial Checks",
		"## Phase Rules",
		"## Role And File Access",
		"## Workflow",
		"## Validation",
		"## Stop Conditions",
		"## Do Not",
	}
	for _, skill := range []string{"sysi-explore", "sysi-capture", "sysi-apply", "sysi-design-change"} {
		content := readFile(t, filepath.Join(root, ".codex", "skills", skill, "SKILL.md"))
		assertContainsAll(t, ".codex/skills/"+skill+"/SKILL.md", content, commonMarkers)
	}

	skillSpecific := map[string][]string{
		"sysi-explore": {
			"allowed /system files",
			"candidate decisions",
			"sysi-capture",
			"avoid implementation",
			"system/security/**",
			"principal-engineer design review",
			"## Design Review Lens",
			"source of truth",
			"invariants",
			"operational recovery",
			"DDIA Mental Model Reference",
			"references/ddia-mental-model.md",
		},
		"sysi-capture": {
			"Finalized Decision",
			"Decision Record",
			"avoid duplicated truth",
			"## System File Routing",
			"Must Own",
			"Must Not Contain",
			"Cross-Link Instead When",
			"system/architecture/decisions",
			"system/contracts/conventions.md",
			"system/contracts/errors.md",
			"system/security/model.md",
			"system/data/db/indexes.md",
		},
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
		"sysi-design-change": {
			"explicit user confirmation",
			"decision artifact",
			"migration or compatibility notes",
			"impacted OpenSpec changes",
			"## Foundation Change Routing",
			"Must Own",
			"Must Not Contain",
			"Cross-Link Instead When",
			"schema evolution",
			"before and after",
			"system/security/",
		},
	}
	for skill, markers := range skillSpecific {
		content := readFile(t, filepath.Join(root, ".codex", "skills", skill, "SKILL.md"))
		assertContainsAll(t, ".codex/skills/"+skill+"/SKILL.md", content, markers)
	}

	ddiaReference := readFile(t, filepath.Join(root, ".codex", "skills", "sysi-explore", "references", "ddia-mental-model.md"))
	assertContainsAll(t, "sysi-explore DDIA reference", ddiaReference, []string{
		"Data Models",
		"Storage And Retrieval",
		"Encoding And Evolution",
		"Replication",
		"Partitioning",
		"Transactions",
		"Consistency",
		"Batch And Stream Processing",
		"Derived Data",
	})
}

func TestCursorAndClaudeInstructionsContainWorkflowBoundaries(t *testing.T) {
	root := t.TempDir()
	if code, out, errOut := runApp(t, root, "init", "--workspaces", "frontend,backend"); code != 0 {
		t.Fatalf("init failed: code=%d stdout=%q stderr=%q", code, out, errOut)
	}
	for _, agent := range []string{"cursor", "claude"} {
		if code, out, errOut := runApp(t, root, "agent", "install", agent); code != 0 {
			t.Fatalf("%s install failed: code=%d stdout=%q stderr=%q", agent, code, out, errOut)
		}
	}

	markers := []string{
		"phase boundaries",
		"/system",
		"OpenSpec",
		"sysi design-change",
		"role",
		"minimal",
		"contracts",
		"security",
	}
	cursor := readFile(t, filepath.Join(root, ".cursor", "rules", "sysi.mdc"))
	assertContainsAll(t, ".cursor/rules/sysi.mdc", cursor, markers)

	claude := readFile(t, filepath.Join(root, "CLAUDE.md"))
	assertContainsAll(t, "CLAUDE.md", claude, markers)
}

func TestClaudeInstallReplacesOnlyManagedSection(t *testing.T) {
	root := t.TempDir()
	if code, out, errOut := runApp(t, root, "init", "--workspaces", "frontend,backend"); code != 0 {
		t.Fatalf("init failed: code=%d stdout=%q stderr=%q", code, out, errOut)
	}

	claudePath := filepath.Join(root, "CLAUDE.md")
	existing := "# Existing\n\nKeep before.\n\n<!-- SYSI:START -->\nold sysi text\n<!-- SYSI:END -->\n\nKeep after.\n"
	if err := os.WriteFile(claudePath, []byte(existing), 0o644); err != nil {
		t.Fatal(err)
	}
	if code, out, errOut := runApp(t, root, "agent", "install", "claude"); code != 0 {
		t.Fatalf("claude install failed: code=%d stdout=%q stderr=%q", code, out, errOut)
	}

	updated := readFile(t, claudePath)
	assertContainsAll(t, "CLAUDE.md", updated, []string{
		"Keep before.",
		"Keep after.",
		"## Sysi",
		"phase boundaries",
	})
	if strings.Contains(updated, "old sysi text") {
		t.Fatalf("managed sysi section was not replaced:\n%s", updated)
	}
	if strings.Count(updated, "<!-- SYSI:START -->") != 1 || strings.Count(updated, "<!-- SYSI:END -->") != 1 {
		t.Fatalf("managed sysi section markers should appear exactly once:\n%s", updated)
	}
}

func TestAgentInstallCommandNamesRemainStable(t *testing.T) {
	root := t.TempDir()
	if code, out, errOut := runApp(t, root, "init", "--workspaces", "frontend,backend"); code != 0 {
		t.Fatalf("init failed: code=%d stdout=%q stderr=%q", code, out, errOut)
	}

	code, out, errOut := runApp(t, root, "help")
	if code != 0 {
		t.Fatalf("help failed: code=%d stdout=%q stderr=%q", code, out, errOut)
	}
	if !strings.Contains(out, "sysi agent install codex|cursor|claude") {
		t.Fatalf("help output should keep stable agent install command names:\n%s", out)
	}
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

func assertContainsAll(t *testing.T, label, content string, markers []string) {
	t.Helper()
	for _, marker := range markers {
		if !strings.Contains(content, marker) {
			t.Fatalf("%s missing %q:\n%s", label, marker, content)
		}
	}
}
