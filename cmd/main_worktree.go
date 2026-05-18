package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/ma-cohen/code-cat/internal/git"
	"github.com/ma-cohen/code-cat/internal/prompt"
	"github.com/spf13/cobra"
)

var mainWorktreeCmd = &cobra.Command{
	Use:     "main-worktree",
	Aliases: []string{"primary-worktree"},
	Short:   "Go to the repository's primary worktree (print path or open a shell)",
	Args:    cobra.NoArgs,
	RunE:    runMainWorktree,
}

func init() {
	mainWorktreeCmd.Flags().Bool("print-path", false, "Print only the primary worktree path on stdout (for cd \"$(ccat ...)\")")
	mainWorktreeCmd.Flags().Bool("no-enter", false, "Do not offer to open a shell in the primary worktree")
}

func runMainWorktree(cmd *cobra.Command, args []string) error {
	if !git.IsInsideRepo() {
		return fmt.Errorf("not inside a git repository")
	}

	printPath, _ := cmd.Flags().GetBool("print-path")
	noEnter, _ := cmd.Flags().GetBool("no-enter")

	cur, err := git.WorktreeTopLevel()
	if err != nil {
		return err
	}
	curAbs, err := filepath.Abs(cur)
	if err != nil {
		return err
	}
	curClean := filepath.Clean(curAbs)

	primary, err := git.PrimaryWorktreePath()
	if err != nil {
		return err
	}

	if curClean == primary {
		fmt.Fprintln(cmd.ErrOrStderr(), "Already in primary worktree.")
	} else {
		fmt.Fprintf(cmd.ErrOrStderr(), "Primary worktree: %s\n", primary)
		fmt.Fprintf(cmd.ErrOrStderr(), "Or from this shell: cd %q\n", primary)
	}

	if printPath {
		fmt.Fprintln(cmd.OutOrStdout(), primary)
	}

	if curClean == primary {
		return nil
	}

	if shouldOfferEnterShell(noEnter, printPath) {
		enter, err := prompt.AskConfirm("Open a shell in the primary worktree?", true)
		if err != nil {
			return err
		}
		if enter {
			return runShellInDir(primary)
		}
	}

	return nil
}
