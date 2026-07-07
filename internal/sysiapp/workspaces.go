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
