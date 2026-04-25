package sysapp

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
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

func runAppWithOpenSpec(t *testing.T, dir, openspecPath string, args ...string) (int, string, string) {
	t.Helper()

	var stdout, stderr bytes.Buffer
	code := New(Options{
		Dir:          dir,
		Stdout:       &stdout,
		Stderr:       &stderr,
		OpenSpecPath: openspecPath,
	}).Run(args)

	return code, stdout.String(), stderr.String()
}

func TestInitScaffoldsProjectAndIsIdempotent(t *testing.T) {
	root := t.TempDir()
	fakeOpenSpec, logPath := writeFakeOpenSpec(t, root)

	code, out, errOut := runAppWithOpenSpec(t, root, fakeOpenSpec, "init")
	if code != 0 {
		t.Fatalf("init failed: code=%d stdout=%q stderr=%q", code, out, errOut)
	}

	wantFiles := []string{
		".sys-orchestrator/state.json",
		".sys-orchestrator/freeze.json",
		".sys-orchestrator/allowlists.json",
		"system/architecture/system.md",
		"system/contracts/api.yaml",
		"system/contracts/events.asyncapi.yaml",
		"system/contracts/auth.md",
		"system/modules/frontend.md",
		"system/modules/backend.md",
		"system/data/schema.sql",
		"system/data/schema.md",
		"system/data/db/indexes.md",
		"system/obs/dashboards/grafana.md",
		"frontend/openspec/config.yaml",
		"backend/openspec/config.yaml",
	}

	for _, rel := range wantFiles {
		if _, err := os.Stat(filepath.Join(root, rel)); err != nil {
			t.Fatalf("expected %s to exist: %v", rel, err)
		}
	}
	if _, err := os.Stat(filepath.Join(root, "openspec", "config.yaml")); err == nil {
		t.Fatalf("root openspec/config.yaml should not be initialized by sys init")
	}
	if _, err := os.Stat(filepath.Join(root, "system", "openspec", "config.yaml")); err == nil {
		t.Fatalf("system/openspec/config.yaml should not be initialized by sys init")
	}
	log := readFile(t, logPath)
	assertContainsAll(t, "openspec log", log, []string{
		"init frontend --tools none",
		"init backend --tools none",
	})
	if strings.Contains(log, "init .") || strings.Contains(log, "init system") {
		t.Fatalf("openspec init should only target frontend/backend, got:\n%s", log)
	}

	if err := os.Remove(filepath.Join(root, "backend", "openspec", "config.yaml")); err != nil {
		t.Fatal(err)
	}
	code, out, errOut = runAppWithOpenSpec(t, root, fakeOpenSpec, "init")
	if code != 0 {
		t.Fatalf("second init failed: code=%d stdout=%q stderr=%q", code, out, errOut)
	}
	if !strings.Contains(out, "already initialized") {
		t.Fatalf("second init should report already initialized, got %q", out)
	}
	log = readFile(t, logPath)
	if strings.Count(log, "init frontend --tools none") != 1 {
		t.Fatalf("frontend should not be reinitialized when config exists:\n%s", log)
	}
	if strings.Count(log, "init backend --tools none") != 2 {
		t.Fatalf("backend should be initialized again after config is missing:\n%s", log)
	}
}

func TestRootDiscoveryAndStatusJSON(t *testing.T) {
	root := t.TempDir()
	fakeOpenSpec, _ := writeFakeOpenSpec(t, root)
	if code, out, errOut := runAppWithOpenSpec(t, root, fakeOpenSpec, "init"); code != 0 {
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
	if status.Role != RoleFrontend {
		t.Fatalf("role = %q, want %q", status.Role, RoleFrontend)
	}
}

func TestValidateReportsMissingRequiredSystemFile(t *testing.T) {
	root := t.TempDir()
	fakeOpenSpec, _ := writeFakeOpenSpec(t, root)
	if code, out, errOut := runAppWithOpenSpec(t, root, fakeOpenSpec, "init"); code != 0 {
		t.Fatalf("init failed: code=%d stdout=%q stderr=%q", code, out, errOut)
	}
	if err := os.Remove(filepath.Join(root, "system", "contracts", "api.yaml")); err != nil {
		t.Fatal(err)
	}

	code, out, errOut := runApp(t, root, "validate")
	if code == 0 {
		t.Fatalf("validate should fail when required file is missing: stdout=%q stderr=%q", out, errOut)
	}
	if !strings.Contains(out+errOut, "system/contracts/api.yaml") {
		t.Fatalf("missing file warning not found in output: stdout=%q stderr=%q", out, errOut)
	}
}

func TestDesignFreezeRecordsBaselineAndCaptureBlocksInBuild(t *testing.T) {
	root := t.TempDir()
	fakeOpenSpec, _ := writeFakeOpenSpec(t, root)
	if code, out, errOut := runAppWithOpenSpec(t, root, fakeOpenSpec, "init"); code != 0 {
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

	archPath := filepath.Join(root, "system", "architecture", "system.md")
	if err := os.WriteFile(archPath, []byte("changed\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	code, out, errOut = runApp(t, root, "status", "--json")
	if code != 0 {
		t.Fatalf("status json failed: code=%d stdout=%q stderr=%q", code, out, errOut)
	}

	var status Status
	if err := json.Unmarshal([]byte(out), &status); err != nil {
		t.Fatal(err)
	}
	if len(status.Validation.Warnings) == 0 {
		t.Fatalf("expected freeze warning after architecture mutation: %#v", status)
	}
	if !strings.Contains(out, "design-change") {
		t.Fatalf("status should mention design-change after frozen file changes: %s", out)
	}
}

func TestDesignCommandsDoNotRequireOpenSpec(t *testing.T) {
	root := t.TempDir()
	fakeOpenSpec, _ := writeFakeOpenSpec(t, root)
	if code, out, errOut := runAppWithOpenSpec(t, root, fakeOpenSpec, "init"); code != 0 {
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
	fakeOpenSpec, _ := writeFakeOpenSpec(t, root)
	if code, out, errOut := runAppWithOpenSpec(t, root, fakeOpenSpec, "init"); code != 0 {
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
		".codex/skills/sys-explore/SKILL.md",
		".codex/skills/sys-capture/SKILL.md",
		".codex/skills/sys-apply/SKILL.md",
		".codex/skills/sys-design-change/SKILL.md",
		".cursor/rules/sys-orchestrator.mdc",
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
	if !strings.Contains(string(claude), "Keep this.") || !strings.Contains(string(claude), "SYS-ORCHESTRATOR") {
		t.Fatalf("CLAUDE.md did not preserve content and add marked section:\n%s", claude)
	}
}

func TestCodexInstructionPacksContainOperationalGuardrails(t *testing.T) {
	root := t.TempDir()
	fakeOpenSpec, _ := writeFakeOpenSpec(t, root)
	if code, out, errOut := runAppWithOpenSpec(t, root, fakeOpenSpec, "init"); code != 0 {
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
	for _, skill := range []string{"sys-explore", "sys-capture", "sys-apply", "sys-design-change"} {
		content := readFile(t, filepath.Join(root, ".codex", "skills", skill, "SKILL.md"))
		assertContainsAll(t, ".codex/skills/"+skill+"/SKILL.md", content, commonMarkers)
	}

	skillSpecific := map[string][]string{
		"sys-explore": {
			"allowed /system files",
			"candidate decisions",
			"sys-capture",
			"avoid implementation",
		},
		"sys-capture": {
			"Finalized Decision",
			"Decision Record",
			"avoid duplicated truth",
			"system/architecture/decisions",
		},
		"sys-apply": {
			"OpenSpec",
			"Superpowers",
			"openspec-apply-change",
			"mandatory",
			"missing prerequisite",
			"frozen /system files",
			"sys design-change",
		},
		"sys-design-change": {
			"explicit user confirmation",
			"migration or compatibility notes",
			"impacted OpenSpec changes",
			"before and after",
		},
	}
	for skill, markers := range skillSpecific {
		content := readFile(t, filepath.Join(root, ".codex", "skills", skill, "SKILL.md"))
		assertContainsAll(t, ".codex/skills/"+skill+"/SKILL.md", content, markers)
	}
}

func TestCursorAndClaudeInstructionsContainWorkflowBoundaries(t *testing.T) {
	root := t.TempDir()
	fakeOpenSpec, _ := writeFakeOpenSpec(t, root)
	if code, out, errOut := runAppWithOpenSpec(t, root, fakeOpenSpec, "init"); code != 0 {
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
		"sys design-change",
		"role",
		"minimal",
	}
	cursor := readFile(t, filepath.Join(root, ".cursor", "rules", "sys-orchestrator.mdc"))
	assertContainsAll(t, ".cursor/rules/sys-orchestrator.mdc", cursor, markers)

	claude := readFile(t, filepath.Join(root, "CLAUDE.md"))
	assertContainsAll(t, "CLAUDE.md", claude, markers)
}

func TestClaudeInstallReplacesOnlyManagedSection(t *testing.T) {
	root := t.TempDir()
	fakeOpenSpec, _ := writeFakeOpenSpec(t, root)
	if code, out, errOut := runAppWithOpenSpec(t, root, fakeOpenSpec, "init"); code != 0 {
		t.Fatalf("init failed: code=%d stdout=%q stderr=%q", code, out, errOut)
	}

	claudePath := filepath.Join(root, "CLAUDE.md")
	existing := "# Existing\n\nKeep before.\n\n<!-- SYS-ORCHESTRATOR:START -->\nold sys text\n<!-- SYS-ORCHESTRATOR:END -->\n\nKeep after.\n"
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
		"## Sys Orchestrator",
		"phase boundaries",
	})
	if strings.Contains(updated, "old sys text") {
		t.Fatalf("managed sys section was not replaced:\n%s", updated)
	}
	if strings.Count(updated, "<!-- SYS-ORCHESTRATOR:START -->") != 1 || strings.Count(updated, "<!-- SYS-ORCHESTRATOR:END -->") != 1 {
		t.Fatalf("managed sys section markers should appear exactly once:\n%s", updated)
	}
}

func TestAgentInstallCommandNamesRemainStable(t *testing.T) {
	root := t.TempDir()
	fakeOpenSpec, _ := writeFakeOpenSpec(t, root)
	if code, out, errOut := runAppWithOpenSpec(t, root, fakeOpenSpec, "init"); code != 0 {
		t.Fatalf("init failed: code=%d stdout=%q stderr=%q", code, out, errOut)
	}

	code, out, errOut := runApp(t, root, "help")
	if code != 0 {
		t.Fatalf("help failed: code=%d stdout=%q stderr=%q", code, out, errOut)
	}
	if !strings.Contains(out, "sys agent install codex|cursor|claude") {
		t.Fatalf("help output should keep stable agent install command names:\n%s", out)
	}
}

func TestBuildWorkflowUsesFakeOpenSpecInBuildPhase(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("shell-script fake executable is POSIX-only")
	}

	root := t.TempDir()
	fake, logPath := writeFakeOpenSpec(t, root)
	if code, out, errOut := runAppWithOpenSpec(t, root, fake, "init"); code != 0 {
		t.Fatalf("init failed: code=%d stdout=%q stderr=%q", code, out, errOut)
	}

	code, out, errOut := runApp(t, root, "change", "propose", "add-login")
	if code == 0 {
		t.Fatalf("change propose should fail before build phase: stdout=%q stderr=%q", out, errOut)
	}

	if code, out, errOut := runApp(t, root, "design", "freeze"); code != 0 {
		t.Fatalf("freeze failed: code=%d stdout=%q stderr=%q", code, out, errOut)
	}

	frontendDir := filepath.Join(root, "frontend")

	var stdout, stderr bytes.Buffer
	code = New(Options{Dir: frontendDir, Stdout: &stdout, Stderr: &stderr, OpenSpecPath: fake}).Run([]string{"change", "propose", "add-login"})
	if code != 0 {
		t.Fatalf("change propose failed: code=%d stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	}

	stdout.Reset()
	stderr.Reset()
	code = New(Options{Dir: frontendDir, Stdout: &stdout, Stderr: &stderr, OpenSpecPath: fake}).Run([]string{"change", "apply", "add-login"})
	if code != 0 {
		t.Fatalf("change apply failed: code=%d stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	}
	if !strings.Contains(stdout.String(), "OpenSpec apply") || !strings.Contains(stdout.String(), "Superpowers") {
		t.Fatalf("apply output should mention OpenSpec apply and Superpowers discipline: %q", stdout.String())
	}

	stdout.Reset()
	stderr.Reset()
	code = New(Options{Dir: frontendDir, Stdout: &stdout, Stderr: &stderr, OpenSpecPath: fake}).Run([]string{"change", "archive", "add-login"})
	if code != 0 {
		t.Fatalf("change archive failed: code=%d stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	}

	logBytes, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatal(err)
	}
	log := string(logBytes)
	if !strings.Contains(log, frontendDir+"|new change add-login") ||
		!strings.Contains(log, frontendDir+"|instructions apply --change add-login --json") ||
		!strings.Contains(log, frontendDir+"|archive add-login") {
		t.Fatalf("fake openspec did not receive expected calls:\n%s", log)
	}
	if _, err := os.Stat(filepath.Join(root, "openspec", "changes", "add-login")); err == nil {
		t.Fatalf("root OpenSpec change should not be created by frontend build command")
	}
}

func writeFakeOpenSpec(t *testing.T, root string) (string, string) {
	t.Helper()
	if runtime.GOOS == "windows" {
		t.Skip("shell-script fake executable is POSIX-only")
	}

	logPath := filepath.Join(root, "openspec.log")
	fake := filepath.Join(root, "fake-openspec")
	script := "#!/bin/sh\n" +
		"printf '%s|%s\\n' \"$PWD\" \"$*\" >> " + shellQuote(logPath) + "\n" +
		"if [ \"$1\" = \"init\" ]; then mkdir -p \"$2/openspec\"; echo 'project: fake' > \"$2/openspec/config.yaml\"; fi\n" +
		"if [ \"$1\" = \"new\" ]; then mkdir -p openspec/changes/add-login; fi\n" +
		"if [ \"$1\" = \"instructions\" ] && [ \"$2\" = \"apply\" ]; then echo '{\"state\":\"ready\"}'; fi\n" +
		"exit 0\n"
	if err := os.WriteFile(fake, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	return fake, logPath
}

func shellQuote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", "'\\''") + "'"
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
