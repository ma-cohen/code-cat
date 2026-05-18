package cmd

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/ma-cohen/code-cat/internal/git"
	"github.com/ma-cohen/code-cat/internal/prompt"
	"github.com/spf13/cobra"
)

const removeAllWorktreesLabel = "All removable worktrees"

var removeWorktreeCmd = &cobra.Command{
	Use:   "remove-worktree",
	Short: "Remove linked git worktrees interactively",
	Args:  cobra.NoArgs,
	RunE:  runRemoveWorktree,
}

func init() {
	rootCmd.AddCommand(removeWorktreeCmd)
	removeWorktreeCmd.Flags().Bool("force", false, "Force remove even if there are uncommitted changes or untracked files")
}

func runRemoveWorktree(cmd *cobra.Command, args []string) error {
	if !git.IsInsideRepo() {
		return fmt.Errorf("not inside a git repository")
	}

	force, _ := cmd.Flags().GetBool("force")

	curTop, err := git.WorktreeTopLevel()
	if err != nil {
		return err
	}
	curTop, err = filepath.Abs(curTop)
	if err != nil {
		return err
	}

	entries, err := git.ListWorktrees()
	if err != nil {
		return err
	}

	var eligible []git.WorktreeEntry
	for _, e := range entries {
		if e.Bare {
			continue
		}
		p, err := filepath.Abs(e.Path)
		if err != nil {
			continue
		}
		if pathsEqual(p, curTop) {
			continue
		}
		eligible = append(eligible, e)
	}

	if len(eligible) == 0 {
		fmt.Println("No removable worktrees (only this checkout is linked).")
		return nil
	}

	allLabel := fmt.Sprintf("%s (%d)", removeAllWorktreesLabel, len(eligible))
	displayOptions := []string{allLabel}
	displayToPath := map[string]string{allLabel: ""}
	for _, e := range eligible {
		label := git.WorktreeLabel(e)
		displayOptions = append(displayOptions, label)
		p, _ := filepath.Abs(e.Path)
		displayToPath[label] = p
	}

	selected, err := prompt.AskMultiSelect("Choose worktrees to remove (space = toggle, enter = confirm)", displayOptions)
	if err != nil {
		return err
	}

	toRemove := resolveWorktreesToRemove(selected, allLabel, eligible, displayToPath)
	if len(toRemove) == 0 {
		fmt.Println("Nothing selected. Aborted.")
		return nil
	}

	sort.Strings(toRemove)

	ok, err := prompt.AskConfirm(fmt.Sprintf("Remove %d worktree(s)?", len(toRemove)), false)
	if err != nil {
		return err
	}
	if !ok {
		fmt.Println("Aborted.")
		return nil
	}

	var removed int
	for _, p := range toRemove {
		fmt.Printf("Removing %s\n", p)
		if err := git.RemoveWorktree(p, force); err != nil {
			return fmt.Errorf("remove %s: %w", p, err)
		}
		removed++
	}

	if removed > 0 {
		_ = git.PruneWorktrees()
	}
	return nil
}

func resolveWorktreesToRemove(
	selected []string,
	allLabel string,
	eligible []git.WorktreeEntry,
	displayToPath map[string]string,
) []string {
	hasAll := false
	for _, s := range selected {
		if s == allLabel {
			hasAll = true
			break
		}
	}

	seen := make(map[string]struct{})
	var out []string

	if hasAll {
		for _, e := range eligible {
			p, err := filepath.Abs(e.Path)
			if err != nil {
				continue
			}
			if _, ok := seen[p]; ok {
				continue
			}
			seen[p] = struct{}{}
			out = append(out, p)
		}
		return out
	}

	for _, s := range selected {
		p := displayToPath[s]
		if p == "" {
			continue
		}
		if _, ok := seen[p]; ok {
			continue
		}
		seen[p] = struct{}{}
		out = append(out, p)
	}
	return out
}

func pathsEqual(a, b string) bool {
	return strings.EqualFold(filepath.Clean(a), filepath.Clean(b))
}
