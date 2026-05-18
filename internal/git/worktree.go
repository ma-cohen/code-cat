package git

import (
	"fmt"
	"path/filepath"
	"strings"
)

// WorktreeEntry is one worktree from `git worktree list --porcelain`.
type WorktreeEntry struct {
	Path        string
	Head        string
	BranchRef   string
	BranchShort string
	Detached    bool
	Bare        bool
}

// ListWorktrees returns linked worktrees using porcelain output.
func ListWorktrees() ([]WorktreeEntry, error) {
	out, err := Run("worktree", "list", "--porcelain")
	if err != nil {
		return nil, err
	}
	return parseWorktreeListPorcelain(out), nil
}

// WorktreeTopLevel returns the absolute path to the top-level of the current worktree.
func WorktreeTopLevel() (string, error) {
	return Run("rev-parse", "--show-toplevel")
}

// PrimaryWorktreePath returns the absolute path to the primary worktree (always the first entry
// in `git worktree list --porcelain`).
func PrimaryWorktreePath() (string, error) {
	entries, err := ListWorktrees()
	if err != nil {
		return "", err
	}
	if len(entries) == 0 {
		return "", fmt.Errorf("no worktrees found")
	}
	abs, err := filepath.Abs(entries[0].Path)
	if err != nil {
		return "", err
	}
	return filepath.Clean(abs), nil
}

// RemoveWorktree runs `git worktree remove` for the given path.
func RemoveWorktree(path string, force bool) error {
	args := []string{"worktree", "remove"}
	if force {
		args = append(args, "--force")
	}
	args = append(args, path)
	_, err := Run(args...)
	return err
}

// PruneWorktrees runs `git worktree prune`.
func PruneWorktrees() error {
	_, err := Run("worktree", "prune")
	return err
}

func parseWorktreeListPorcelain(out string) []WorktreeEntry {
	if strings.TrimSpace(out) == "" {
		return nil
	}

	var entries []WorktreeEntry
	var cur *WorktreeEntry

	finish := func() {
		if cur != nil {
			entries = append(entries, *cur)
			cur = nil
		}
	}

	for _, line := range strings.Split(out, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "worktree ") {
			finish()
			p := strings.TrimSpace(strings.TrimPrefix(line, "worktree"))
			cur = &WorktreeEntry{Path: p}
			continue
		}
		if cur == nil {
			continue
		}
		switch {
		case strings.HasPrefix(line, "HEAD "):
			cur.Head = strings.TrimSpace(strings.TrimPrefix(line, "HEAD"))
		case strings.HasPrefix(line, "branch "):
			ref := strings.TrimSpace(strings.TrimPrefix(line, "branch"))
			cur.BranchRef = ref
			if strings.HasPrefix(ref, "refs/heads/") {
				cur.BranchShort = strings.TrimPrefix(ref, "refs/heads/")
			}
		case line == "detached":
			cur.Detached = true
		case line == "bare":
			cur.Bare = true
		}
	}
	finish()
	return entries
}

// WorktreeLabel builds a short human-readable label for prompts.
func WorktreeLabel(e WorktreeEntry) string {
	switch {
	case e.Bare:
		return fmt.Sprintf("%s  (bare)", e.Path)
	case e.Detached || e.BranchShort == "":
		return fmt.Sprintf("%s  (detached)", e.Path)
	default:
		return fmt.Sprintf("%s  (%s)", e.Path, e.BranchShort)
	}
}
