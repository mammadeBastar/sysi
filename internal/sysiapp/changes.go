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
	matches, err := filepath.Glob(filepath.Join(root, workspace, "changes", "archive", "????-??-??-"+name))
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
	if slugify(name) == "" {
		return errors.New("change name must contain letters or digits")
	}
	if slugify(name) != name {
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
	cleanup := func(err error) error {
		os.RemoveAll(dir)
		return err
	}
	date := now.Format("2006-01-02")
	files := []struct{ name, content string }{
		{"proposal.md", fmt.Sprintf(changeProposalTemplate, name, date, workspace)},
		{"design.md", fmt.Sprintf(changeDesignTemplate, name)},
		{"tasks.md", fmt.Sprintf(changeTasksTemplate, name)},
	}
	for _, file := range files {
		if err := os.WriteFile(filepath.Join(dir, file.name), []byte(file.content), 0o644); err != nil {
			return cleanup(err)
		}
	}
	stamp := now.Format(time.RFC3339)
	meta := ChangeMeta{Name: name, Workspace: workspace, Status: ChangeStatusProposed, CreatedAt: stamp, UpdatedAt: stamp}
	if err := saveChangeMeta(root, workspace, name, meta); err != nil {
		return cleanup(err)
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

func (a *App) changeApply(root, workspace, name string) error {
	return errors.New("not implemented")
}

func (a *App) changeArchive(root, workspace, name string, now time.Time) error {
	return errors.New("not implemented")
}
