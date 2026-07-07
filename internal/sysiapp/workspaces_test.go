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

func TestWorkspaceRemoveIgnoresArchivedChanges(t *testing.T) {
	root := initProject(t, "api")

	// Only archived changes exist; the active-change guard must ignore archive.
	archivedDir := filepath.Join(root, "api", "changes", "archive", "old-change")
	if err := os.MkdirAll(archivedDir, 0o755); err != nil {
		t.Fatal(err)
	}

	code, out, errOut := runApp(t, root, "workspace", "remove", "api")
	if code != 0 {
		t.Fatalf("remove with only archived changes should succeed without --force: code=%d stdout=%q stderr=%q", code, out, errOut)
	}

	// Removing the last workspace leaves an empty list, not null.
	stateJSON := readFile(t, filepath.Join(root, ".sysi", "state.json"))
	if !strings.Contains(stateJSON, `"workspaces": []`) {
		t.Fatalf("state.json should contain \"workspaces\": [] after removing last workspace:\n%s", stateJSON)
	}
	if strings.Contains(stateJSON, `"workspaces": null`) {
		t.Fatalf("state.json must not contain null workspaces:\n%s", stateJSON)
	}

	code, out, errOut = runApp(t, root, "workspace", "list")
	if code != 0 {
		t.Fatalf("workspace list failed: code=%d stdout=%q stderr=%q", code, out, errOut)
	}
	if !strings.Contains(out, "no workspaces declared") {
		t.Fatalf("empty workspace list should print no workspaces declared: %q", out)
	}
}

func TestWorkspaceAddRemoveRejectFlagAsName(t *testing.T) {
	root := initProject(t, "api")

	code, out, errOut := runApp(t, root, "workspace", "add", "--force")
	if code == 0 {
		t.Fatalf("workspace add with flag as name should fail: stdout=%q stderr=%q", out, errOut)
	}
	if !strings.Contains(out+errOut, "usage: sysi workspace add <name>") {
		t.Fatalf("workspace add flag misuse should print usage: stdout=%q stderr=%q", out, errOut)
	}

	code, out, errOut = runApp(t, root, "workspace", "remove", "--force")
	if code == 0 {
		t.Fatalf("workspace remove with flag as name should fail: stdout=%q stderr=%q", out, errOut)
	}
	if !strings.Contains(out+errOut, "usage: sysi workspace remove <name> [--force]") {
		t.Fatalf("workspace remove flag misuse should print usage: stdout=%q stderr=%q", out, errOut)
	}
}
