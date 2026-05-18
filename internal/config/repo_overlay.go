package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type dotCodeCatYAML struct {
	BaseBranch   string `yaml:"base_branch"`
	BranchPrefix string `yaml:"branch_prefix"`
	WorktreeRoot string `yaml:"worktree_root"`
}

// MergeDotCodeCatRepo overlays repo-root `.code-cat.yml` onto C for WorktreeRoot and BranchPrefix.
// BaseBranch is intentionally not applied here; callers such as new-worktree resolve the base from the remote default unless --base is set.
func MergeDotCodeCatRepo(repoRoot string) error {
	path := filepath.Join(repoRoot, ".code-cat.yml")
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("read %s: %w", path, err)
	}
	var overlay dotCodeCatYAML
	if err := yaml.Unmarshal(data, &overlay); err != nil {
		return fmt.Errorf("parse %s: %w", path, err)
	}
	if overlay.WorktreeRoot != "" {
		C.WorktreeRoot = overlay.WorktreeRoot
	}
	if overlay.BranchPrefix != "" {
		C.BranchPrefix = overlay.BranchPrefix
	}
	return nil
}
