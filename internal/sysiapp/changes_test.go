package sysiapp

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
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

	// An archived change whose name merely ends with the proposed name is not a collision.
	unrelated := filepath.Join(root, "api", "changes", "archive", "2026-01-01-foo-suffix")
	if err := os.MkdirAll(unrelated, 0o755); err != nil {
		t.Fatal(err)
	}
	if code, out, errOut := runApp(t, apiDir, "change", "propose", "suffix"); code != 0 {
		t.Fatalf("propose of suffix should not collide with 2026-01-01-foo-suffix: code=%d stdout=%q stderr=%q", code, out, errOut)
	}

	// Reserved name fails.
	if code, out, errOut := runApp(t, apiDir, "change", "propose", "archive"); code == 0 {
		t.Fatalf("propose of reserved name archive should fail: stdout=%q stderr=%q", out, errOut)
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
		"SYSI CHANGE APPLY",
		"Status: applying",
		"proposal.md",
		"design.md",
		"tasks.md",
		"Superpowers",
		"TDD",
		"design drift",
		"sysi design-change",
	})

	metaPath := filepath.Join(root, "api", "changes", "add-login", "meta.json")
	var meta ChangeMeta
	if err := json.Unmarshal([]byte(readFile(t, metaPath)), &meta); err != nil {
		t.Fatal(err)
	}
	if meta.Status != ChangeStatusApplying {
		t.Fatalf("status = %q, want %q", meta.Status, ChangeStatusApplying)
	}

	// Re-apply is idempotent.
	if code, _, _ := runApp(t, apiDir, "change", "apply", "add-login"); code != 0 {
		t.Fatal("re-apply should succeed")
	}
	meta = ChangeMeta{}
	if err := json.Unmarshal([]byte(readFile(t, metaPath)), &meta); err != nil {
		t.Fatal(err)
	}
	if meta.Status != ChangeStatusApplying {
		t.Fatalf("status after re-apply = %q, want %q", meta.Status, ChangeStatusApplying)
	}
}

func seedChangeStatus(t *testing.T, root, workspace, name, status string) {
	t.Helper()
	metaPath := filepath.Join(root, workspace, "changes", name, "meta.json")
	var meta ChangeMeta
	if err := json.Unmarshal([]byte(readFile(t, metaPath)), &meta); err != nil {
		t.Fatal(err)
	}
	meta.Status = status
	data, err := json.Marshal(meta)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(metaPath, data, 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestChangeApplyRejectsPathTraversalName(t *testing.T) {
	root := initBuildProject(t, "api,web")
	webDir := filepath.Join(root, "web")
	if code, _, _ := runApp(t, webDir, "change", "propose", "cross-target"); code != 0 {
		t.Fatal("propose in web failed")
	}
	webMetaPath := filepath.Join(root, "web", "changes", "cross-target", "meta.json")
	before := readFile(t, webMetaPath)

	apiDir := filepath.Join(root, "api")
	code, out, errOut := runApp(t, apiDir, "change", "apply", "../../web/changes/cross-target")
	if code == 0 {
		t.Fatalf("traversal name should fail: stdout=%q stderr=%q", out, errOut)
	}
	if !strings.Contains(out+errOut, "slug") {
		t.Fatalf("error should mention slug: stdout=%q stderr=%q", out, errOut)
	}
	if after := readFile(t, webMetaPath); after != before {
		t.Fatalf("traversal apply mutated other workspace's change meta:\nbefore=%q\nafter=%q", before, after)
	}
}

func TestChangeApplyCorruptMetaReportsReadError(t *testing.T) {
	root := initBuildProject(t, "api")
	apiDir := filepath.Join(root, "api")
	if code, _, _ := runApp(t, apiDir, "change", "propose", "add-login"); code != 0 {
		t.Fatal("propose failed")
	}
	metaPath := filepath.Join(root, "api", "changes", "add-login", "meta.json")
	if err := os.WriteFile(metaPath, []byte("{broken"), 0o644); err != nil {
		t.Fatal(err)
	}

	code, out, errOut := runApp(t, apiDir, "change", "apply", "add-login")
	if code == 0 {
		t.Fatalf("apply with corrupt meta should fail: stdout=%q stderr=%q", out, errOut)
	}
	combined := out + errOut
	if !strings.Contains(combined, "meta.json") {
		t.Fatalf("error should mention meta.json: %q", combined)
	}
	if strings.Contains(combined, "not found") {
		t.Fatalf("corrupt meta must not be reported as not found: %q", combined)
	}
}

func TestChangeApplyArchivedStatusErrors(t *testing.T) {
	root := initBuildProject(t, "api")
	apiDir := filepath.Join(root, "api")
	if code, _, _ := runApp(t, apiDir, "change", "propose", "add-login"); code != 0 {
		t.Fatal("propose failed")
	}
	seedChangeStatus(t, root, "api", "add-login", ChangeStatusArchived)

	code, out, errOut := runApp(t, apiDir, "change", "apply", "add-login")
	if code == 0 {
		t.Fatalf("apply of archived change should fail: stdout=%q stderr=%q", out, errOut)
	}
	if !strings.Contains(out+errOut, "archived") {
		t.Fatalf("error should mention archived: stdout=%q stderr=%q", out, errOut)
	}
}

func TestChangeApplyUnexpectedStatusErrors(t *testing.T) {
	root := initBuildProject(t, "api")
	apiDir := filepath.Join(root, "api")
	if code, _, _ := runApp(t, apiDir, "change", "propose", "add-login"); code != 0 {
		t.Fatal("propose failed")
	}
	seedChangeStatus(t, root, "api", "add-login", "bogus")

	code, out, errOut := runApp(t, apiDir, "change", "apply", "add-login")
	if code == 0 {
		t.Fatalf("apply with unexpected status should fail: stdout=%q stderr=%q", out, errOut)
	}
	if !strings.Contains(out+errOut, "bogus") {
		t.Fatalf("error should name the unexpected status: stdout=%q stderr=%q", out, errOut)
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
	code, out, errOut = runApp(t, apiDir, "change", "archive", "ghost")
	if code == 0 {
		t.Fatalf("archive of unknown change should fail: stdout=%q stderr=%q", out, errOut)
	}
	assertContainsAll(t, "unknown change error", out+errOut, []string{"not found", "available"})
}

func TestChangeArchiveWarnsOnUnexpectedStatus(t *testing.T) {
	root := initBuildProject(t, "api")
	apiDir := filepath.Join(root, "api")
	if code, _, _ := runApp(t, apiDir, "change", "propose", "add-login"); code != 0 {
		t.Fatal("propose failed")
	}
	seedChangeStatus(t, root, "api", "add-login", "bogus")

	code, out, errOut := runApp(t, apiDir, "change", "archive", "add-login")
	if code != 0 {
		t.Fatalf("archive with unexpected status should still succeed: code=%d stdout=%q stderr=%q", code, out, errOut)
	}
	if !strings.Contains(out+errOut, `unexpected status "bogus"`) {
		t.Fatalf("archive should warn about unexpected status: stdout=%q stderr=%q", out, errOut)
	}
}

func TestChangeArchiveFailsWhenTargetExists(t *testing.T) {
	root := initBuildProject(t, "api")
	apiDir := filepath.Join(root, "api")
	if code, _, _ := runApp(t, apiDir, "change", "propose", "add-login"); code != 0 {
		t.Fatal("propose failed")
	}
	today := time.Now().UTC().Format("2006-01-02")
	target := filepath.Join(root, "api", "changes", "archive", today+"-add-login")
	if err := os.MkdirAll(target, 0o755); err != nil {
		t.Fatal(err)
	}

	code, out, errOut := runApp(t, apiDir, "change", "archive", "add-login")
	if code == 0 {
		t.Fatalf("archive onto existing target should fail: stdout=%q stderr=%q", out, errOut)
	}
	if !strings.Contains(out+errOut, "archive target already exists") {
		t.Fatalf("error should mention existing archive target: stdout=%q stderr=%q", out, errOut)
	}
	if _, err := os.Stat(filepath.Join(root, "api", "changes", "add-login")); err != nil {
		t.Fatalf("active change dir should remain after failed archive: %v", err)
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
