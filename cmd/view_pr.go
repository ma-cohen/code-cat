package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/ma-cohen/code-cat/internal/git"
	"github.com/ma-cohen/code-cat/internal/provider"
	"github.com/spf13/cobra"
)

var viewPRCmd = &cobra.Command{
	Use:   "view-pr",
	Short: "Open the pull request or merge request for the current branch in the browser",
	Args:  cobra.NoArgs,
	RunE:  runViewPR,
}

func init() {
	rootCmd.AddCommand(viewPRCmd)
}

func runViewPR(cmd *cobra.Command, args []string) error {
	if !git.IsInsideRepo() {
		return fmt.Errorf("not inside a git repository")
	}

	branch, err := git.CurrentBranch()
	if err != nil {
		return err
	}

	remoteURL, err := git.RemoteURL("origin")
	if err != nil {
		return fmt.Errorf("could not get remote URL: %w", err)
	}

	p := provider.Detect(remoteURL)

	if _, err := exec.LookPath(p.CLI); err != nil {
		installMsg := map[string]string{
			"gh":   "Install it from https://cli.github.com/ then run: gh auth login",
			"glab": "Install it from https://gitlab.com/gitlab-org/cli then run: glab auth login",
		}
		return fmt.Errorf("%s is not installed or not in PATH\n%s", p.CLI, installMsg[p.CLI])
	}

	fmt.Printf("Opening PR/MR for branch %s...\n", branch)
	c := exec.Command(p.CLI, p.ViewPRCmd...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}
