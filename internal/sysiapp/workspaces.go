package sysiapp

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
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
		if info, err := os.Stat(filepath.Join(root, name)); err == nil && !info.IsDir() {
			return fmt.Errorf("workspace %q conflicts with existing file %s", name, name)
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
