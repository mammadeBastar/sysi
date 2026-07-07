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
