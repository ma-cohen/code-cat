package cmd

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"unicode"

	"github.com/ma-cohen/code-cat/internal/config"
	"github.com/ma-cohen/code-cat/internal/git"
	"github.com/ma-cohen/code-cat/internal/prompt"
	"github.com/spf13/cobra"
)

var prCmd = &cobra.Command{
	Use:   "pr",
	Short: "Push the current branch and open a pull request via gh",
	Args:  cobra.NoArgs,
	RunE:  runPR,
}

func init() {
	rootCmd.AddCommand(prCmd)
	prCmd.Flags().String("base", "", "Base branch for the PR (overrides config)")
	prCmd.Flags().Bool("no-push", false, "Skip pushing the branch to origin")
	prCmd.Flags().Bool("draft", false, "Create the PR as a draft")
	prCmd.Flags().Bool("web", false, "Open the PR form in the browser instead of creating via CLI")
}

// branchToTitle converts a branch name into a human-readable PR title.
func branchToTitle(branch string) string {
	prefixes := []string{
		"feature/", "feat/",
		"fix/", "bugfix/", "bug/",
		"hotfix/",
		"chore/",
		"docs/",
		"refactor/",
		"test/",
		"ci/",
	}
	for _, p := range prefixes {
		branch = strings.TrimPrefix(branch, p)
	}
	branch = strings.ReplaceAll(branch, "-", " ")
	branch = strings.ReplaceAll(branch, "_", " ")
	if len(branch) == 0 {
		return branch
	}
	runes := []rune(branch)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

// runGH executes a gh command and returns trimmed stdout.
func runGH(args ...string) (string, error) {
	cmd := exec.Command("gh", args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		msg := strings.TrimSpace(stderr.String())
		if msg == "" {
			msg = err.Error()
		}
		return "", fmt.Errorf("gh %s: %s", strings.Join(args, " "), msg)
	}
	return strings.TrimSpace(stdout.String()), nil
}

func runPR(cmd *cobra.Command, args []string) error {
	if !git.IsInsideRepo() {
		return fmt.Errorf("not inside a git repository")
	}

	if _, err := exec.LookPath("gh"); err != nil {
		return fmt.Errorf("gh (GitHub CLI) is not installed or not in PATH\nInstall it from https://cli.github.com/ then run: gh auth login")
	}

	branch, err := git.CurrentBranch()
	if err != nil {
		return err
	}

	base, _ := cmd.Flags().GetString("base")
	if base == "" {
		base = config.C.BaseBranch
	}
	if base == "" {
		base, err = git.DefaultBranch()
		if err != nil {
			return err
		}
	}

	if branch == base {
		return fmt.Errorf("current branch %q is the base branch; check out a feature branch first", branch)
	}

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

	noPush, _ := cmd.Flags().GetBool("no-push")
	if !noPush {
		fmt.Printf("Pushing %s to origin...\n", branch)
		if _, err := git.Run("push", "-u", "origin", branch); err != nil {
			return err
		}
	}

	title, err := prompt.AskString("PR title", branchToTitle(branch))
	if err != nil {
		return err
	}

	body, err := prompt.AskOptionalString("PR body (leave empty to skip)", "")
	if err != nil {
		return err
	}

	draft, _ := cmd.Flags().GetBool("draft")
	web, _ := cmd.Flags().GetBool("web")

	ghArgs := []string{"pr", "create",
		"--base", base,
		"--title", title,
		"--body", body,
	}
	if draft {
		ghArgs = append(ghArgs, "--draft")
	}
	if web {
		ghArgs = append(ghArgs, "--web")
	}

	if web {
		fmt.Println("Opening PR form in your browser...")
		ghCmd := exec.Command("gh", ghArgs...)
		ghCmd.Stdout = os.Stdout
		ghCmd.Stderr = os.Stderr
		return ghCmd.Run()
	}

	url, err := runGH(ghArgs...)
	if err != nil {
		return err
	}
	fmt.Printf("Pull request created: %s\n", url)
	return nil
}
