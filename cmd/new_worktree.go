package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/ma-cohen/code-cat/internal/config"
	"github.com/ma-cohen/code-cat/internal/git"
	"github.com/ma-cohen/code-cat/internal/prompt"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var newWorktreeCmd = &cobra.Command{
	Use:   "new-worktree [path]",
	Short: "Create a new worktree on a fresh branch from your remote default trunk",
	Long: `Creates a new git worktree with a new branch based on origin's default branch (origin/HEAD,
or otherwise main/master after fetch). Your current checked-out branch does not affect the base.

Repo-root .code-cat.yml values worktree_root and branch_prefix apply regardless of your current directory.

Pass --base to use a different remote-tracking branch instead of the detected default.`,
	Args: cobra.MaximumNArgs(1),
	RunE: runNewWorktree,
}

func init() {
	newWorktreeCmd.Flags().String("base", "", "Branch off origin/<value> instead of the remote default branch")
	newWorktreeCmd.Flags().String("branch", "", "Name for the new branch in the worktree")
	newWorktreeCmd.Flags().Bool("no-fetch", false, "Skip git fetch")
	newWorktreeCmd.Flags().Bool("print-path", false, "Print only the new worktree path on stdout (for cd \"$(ccat ...)\")")
	newWorktreeCmd.Flags().Bool("no-enter", false, "Do not offer to open a shell in the new worktree")
}

func runNewWorktree(cmd *cobra.Command, args []string) error {
	if !git.IsInsideRepo() {
		return fmt.Errorf("not inside a git repository")
	}

	repoRoot, err := git.WorktreeTopLevel()
	if err != nil {
		return err
	}
	if err := config.MergeDotCodeCatRepo(repoRoot); err != nil {
		return err
	}

	baseFlag, _ := cmd.Flags().GetString("base")
	branchFlag, _ := cmd.Flags().GetString("branch")
	noFetch, _ := cmd.Flags().GetBool("no-fetch")
	printPath, _ := cmd.Flags().GetBool("print-path")
	noEnter, _ := cmd.Flags().GetBool("no-enter")

	if !noFetch {
		fmt.Fprintf(cmd.ErrOrStderr(), "Fetching origin...\n")
		if _, err := git.RunIn(repoRoot, "fetch", "origin"); err != nil {
			return err
		}
	}

	base := baseFlag
	if base == "" {
		base, err = git.DefaultBranchIn(repoRoot)
		if err != nil {
			return err
		}
	}

	branchName := branchFlag
	if branchName == "" {
		branchName, err = prompt.AskString("New branch name", config.C.BranchPrefix)
		if err != nil {
			return err
		}
	}

	var wtPath string
	if len(args) > 0 {
		wtPath = args[0]
	} else {
		defaultPath := filepath.Join(repoRoot, config.C.WorktreeRoot, filepath.Base(branchName))
		defaultPath, err = filepath.Abs(defaultPath)
		if err != nil {
			return err
		}
		wtPath, err = prompt.AskString("Worktree path", defaultPath)
		if err != nil {
			return err
		}
	}

	if _, err := git.RunIn(repoRoot, "worktree", "add", "-b", branchName, wtPath, "origin/"+base); err != nil {
		return err
	}

	absPath, err := filepath.Abs(wtPath)
	if err != nil {
		return err
	}

	fmt.Fprintf(cmd.ErrOrStderr(), "Worktree created at: %s\n", absPath)
	fmt.Fprintf(cmd.ErrOrStderr(), "Branch: %s\n", branchName)
	fmt.Fprintf(cmd.ErrOrStderr(), "Or from this shell: cd %q\n", absPath)

	if printPath {
		fmt.Fprintln(cmd.OutOrStdout(), absPath)
	}

	if shouldOfferEnterShell(noEnter, printPath) {
		enter, err := prompt.AskConfirm("Open a shell in the new worktree?", true)
		if err != nil {
			return err
		}
		if enter {
			return runShellInDir(absPath)
		}
	}

	return nil
}

func shouldOfferEnterShell(noEnter, printPath bool) bool {
	if noEnter || printPath {
		return false
	}
	return term.IsTerminal(int(os.Stdin.Fd())) && term.IsTerminal(int(os.Stdout.Fd()))
}

func runShellInDir(dir string) error {
	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "/bin/sh"
	}
	c := exec.Command(shell)
	c.Dir = dir
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Env = os.Environ()
	return c.Run()
}
