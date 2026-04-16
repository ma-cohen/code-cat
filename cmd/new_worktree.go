package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/ma-cohen/code-cat/internal/config"
	"github.com/ma-cohen/code-cat/internal/git"
	"github.com/ma-cohen/code-cat/internal/prompt"
	"github.com/spf13/cobra"
)

var newWorktreeCmd = &cobra.Command{
	Use:   "new-worktree [path]",
	Short: "Create a new git worktree on a fresh branch",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runNewWorktree,
}

func init() {
	rootCmd.AddCommand(newWorktreeCmd)
	newWorktreeCmd.Flags().String("base", "", "Base branch to branch from (overrides config)")
	newWorktreeCmd.Flags().String("branch", "", "Name for the new branch in the worktree")
	newWorktreeCmd.Flags().Bool("no-fetch", false, "Skip git fetch")
}

func runNewWorktree(cmd *cobra.Command, args []string) error {
	if !git.IsInsideRepo() {
		return fmt.Errorf("not inside a git repository")
	}

	base, _ := cmd.Flags().GetString("base")
	if base == "" {
		base = config.C.BaseBranch
	}
	branchFlag, _ := cmd.Flags().GetString("branch")
	noFetch, _ := cmd.Flags().GetBool("no-fetch")

	// Fetch
	if !noFetch {
		fmt.Printf("Fetching origin...\n")
		if _, err := git.Run("fetch", "origin"); err != nil {
			return err
		}
	}

	// Determine branch name
	branchName := branchFlag
	if branchName == "" {
		var err error
		branchName, err = prompt.AskString("New branch name", config.C.BranchPrefix)
		if err != nil {
			return err
		}
	}

	// Determine worktree path
	var wtPath string
	if len(args) > 0 {
		wtPath = args[0]
	} else {
		defaultPath := filepath.Join(config.C.WorktreeRoot, filepath.Base(branchName))
		var err error
		wtPath, err = prompt.AskString("Worktree path", defaultPath)
		if err != nil {
			return err
		}
	}

	// Create worktree on new branch from origin/<base>
	if _, err := git.Run("worktree", "add", "-b", branchName, wtPath, "origin/"+base); err != nil {
		return err
	}

	absPath, _ := filepath.Abs(wtPath)
	fmt.Printf("Worktree created at: %s\n", absPath)
	fmt.Printf("Branch: %s\n", branchName)

	return nil
}
