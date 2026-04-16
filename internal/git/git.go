package git

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// Run executes a git command in the current working directory and returns trimmed stdout.
func Run(args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("git %s: %s", strings.Join(args, " "), strings.TrimSpace(stderr.String()))
	}
	return strings.TrimSpace(stdout.String()), nil
}

// CurrentBranch returns the name of the currently checked-out branch.
func CurrentBranch() (string, error) {
	return Run("rev-parse", "--abbrev-ref", "HEAD")
}

// HasUncommitted returns true if there are uncommitted changes in the working tree.
func HasUncommitted() (bool, error) {
	out, err := Run("status", "--porcelain")
	if err != nil {
		return false, err
	}
	return len(out) > 0, nil
}

// IsInsideRepo returns true when the current directory is inside a git repository.
func IsInsideRepo() bool {
	_, err := Run("rev-parse", "--git-dir")
	return err == nil
}
