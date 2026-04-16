package cmd

import (
	"fmt"
	"os"

	"github.com/ma-cohen/code-cat/internal/config"
	"github.com/ma-cohen/code-cat/internal/git"
	"github.com/ma-cohen/code-cat/internal/prompt"
	"github.com/spf13/cobra"
)

var homeCmd = &cobra.Command{
	Use:   "home",
	Short: "Go to the latest base branch from origin",
	Args:  cobra.NoArgs,
	RunE:  runHome,
}

func init() {
	rootCmd.AddCommand(homeCmd)
	homeCmd.Flags().String("base", "", "Base branch to go to (overrides config)")
	homeCmd.Flags().Bool("no-fetch", false, "Skip git fetch")
}

func runHome(cmd *cobra.Command, args []string) error {
	if !git.IsInsideRepo() {
		return fmt.Errorf("not inside a git repository")
	}

	base, _ := cmd.Flags().GetString("base")
	if base == "" {
		base = config.C.BaseBranch
	}
	if base == "" {
		var err error
		base, err = git.DefaultBranch()
		if err != nil {
			return err
		}
	}
	noFetch, _ := cmd.Flags().GetBool("no-fetch")

	dirty, err := git.HasUncommitted()
	if err != nil {
		return err
	}
	if dirty {
		choice, err := prompt.AskSelect(
			"You have uncommitted changes. What would you like to do?",
			[]string{"Stash", "Discard", "Abort"},
		)
		if err != nil {
			return err
		}
		switch choice {
		case "Stash":
			fmt.Println("Stashing changes...")
			if _, err := git.Run("stash"); err != nil {
				return err
			}
		case "Discard":
			fmt.Println("Discarding changes...")
			if _, err := git.Run("reset", "--hard", "HEAD"); err != nil {
				return err
			}
			if _, err := git.Run("clean", "-fd"); err != nil {
				return err
			}
		case "Abort":
			fmt.Println("Aborted.")
			os.Exit(0)
		}
	}

	if !noFetch {
		fmt.Println("Fetching origin...")
		if _, err := git.Run("fetch", "origin"); err != nil {
			return err
		}
	}

	fmt.Printf("Switching to %s...\n", base)
	if _, err := git.Run("checkout", base); err != nil {
		return err
	}
	if _, err := git.Run("reset", "--hard", "origin/"+base); err != nil {
		return err
	}

	fmt.Printf("Up to date on %s.\n", base)
	return nil
}
