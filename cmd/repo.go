package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/ma-cohen/code-cat/internal/git"
	"github.com/ma-cohen/code-cat/internal/provider"
	"github.com/spf13/cobra"
)

var repoCmd = &cobra.Command{
	Use:   "repo",
	Short: "Open the repository in the browser (GitHub or GitLab)",
	Args:  cobra.NoArgs,
	RunE:  runRepo,
}

func init() {
	rootCmd.AddCommand(repoCmd)
}

func runRepo(cmd *cobra.Command, args []string) error {
	if !git.IsInsideRepo() {
		return fmt.Errorf("not inside a git repository")
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

	fmt.Println("Opening repository in browser...")
	c := exec.Command(p.CLI, p.BrowseRepoCmd...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}
