package cmd

import (
	"fmt"

	"github.com/ma-cohen/code-cat/internal/git"
	"github.com/ma-cohen/code-cat/internal/prompt"
	"github.com/spf13/cobra"
)

var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Discard local changes and reset the current branch to upstream",
	Args:  cobra.NoArgs,
	RunE:  runReset,
}

func init() {
	rootCmd.AddCommand(resetCmd)
	resetCmd.Flags().Bool("force", false, "Skip confirmation prompt")
}

func runReset(cmd *cobra.Command, args []string) error {
	if !git.IsInsideRepo() {
		return fmt.Errorf("not inside a git repository")
	}

	upstream, err := git.UpstreamBranch()
	if err != nil {
		return err
	}

	force, _ := cmd.Flags().GetBool("force")
	if !force {
		confirmed, err := prompt.AskConfirm(
			fmt.Sprintf("Discard all local changes, untracked files, and unpushed commits by resetting to %s?", upstream),
			false,
		)
		if err != nil {
			return err
		}
		if !confirmed {
			fmt.Println("Aborted.")
			return nil
		}
	}

	fmt.Printf("Resetting current branch to %s...\n", upstream)
	if _, err := git.Run("reset", "--hard", upstream); err != nil {
		return err
	}
	if _, err := git.Run("clean", "-fd"); err != nil {
		return err
	}

	fmt.Println("Local changes removed.")
	return nil
}
