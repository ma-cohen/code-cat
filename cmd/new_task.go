package cmd

import (
	"fmt"
	"os"

	"github.com/ma-cohen/code-cat/internal/config"
	"github.com/ma-cohen/code-cat/internal/git"
	"github.com/ma-cohen/code-cat/internal/prompt"
	"github.com/spf13/cobra"
)

var newTaskCmd = &cobra.Command{
	Use:   "new-task [branch-name]",
	Short: "Fetch latest base branch and check out a new feature branch",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runNewTask,
}

func init() {
	rootCmd.AddCommand(newTaskCmd)
	newTaskCmd.Flags().String("base", "", "Base branch to update from (overrides config)")
	newTaskCmd.Flags().Bool("no-fetch", false, "Skip git fetch")
	newTaskCmd.Flags().Bool("push", false, "Push new branch to origin after creation")
}

func runNewTask(cmd *cobra.Command, args []string) error {
	if !git.IsInsideRepo() {
		return fmt.Errorf("not inside a git repository")
	}

	base, _ := cmd.Flags().GetString("base")
	if base == "" {
		base = config.C.BaseBranch
	}
	noFetch, _ := cmd.Flags().GetBool("no-fetch")
	push, _ := cmd.Flags().GetBool("push")

	// Warn about uncommitted changes
	dirty, err := git.HasUncommitted()
	if err != nil {
		return err
	}
	if dirty {
		proceed, err := prompt.AskConfirm("You have uncommitted changes. Continue anyway?", false)
		if err != nil {
			return err
		}
		if !proceed {
			fmt.Println("Aborted.")
			os.Exit(0)
		}
	}

	// Fetch
	if !noFetch {
		fmt.Printf("Fetching origin...\n")
		if _, err := git.Run("fetch", "origin"); err != nil {
			return err
		}
	}

	// Update base branch
	fmt.Printf("Updating %s from origin/%s...\n", base, base)
	if _, err := git.Run("checkout", base); err != nil {
		return err
	}
	if _, err := git.Run("reset", "--hard", "origin/"+base); err != nil {
		return err
	}

	// Determine branch name
	var branchName string
	if len(args) > 0 {
		branchName = args[0]
	} else {
		branchName, err = prompt.AskString("New branch name", config.C.BranchPrefix)
		if err != nil {
			return err
		}
	}

	// Create and checkout new branch
	if _, err := git.Run("checkout", "-b", branchName); err != nil {
		return err
	}
	fmt.Printf("Switched to new branch '%s'\n", branchName)

	// Optionally push
	if push {
		fmt.Printf("Pushing %s to origin...\n", branchName)
		if _, err := git.Run("push", "-u", "origin", branchName); err != nil {
			return err
		}
	}

	return nil
}
