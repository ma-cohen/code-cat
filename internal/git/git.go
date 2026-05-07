package git

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// StashEntry describes one entry from git stash list.
type StashEntry struct {
	Ref     string
	Message string
}

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

// StashPush stashes tracked and untracked changes with a message.
func StashPush(message string) error {
	_, err := Run("stash", "push", "-u", "-m", message)
	return err
}

// ListStashes returns stash entries in Git's newest-first order.
func ListStashes() ([]StashEntry, error) {
	out, err := Run("stash", "list", "--format=%gd%x00%gs")
	if err != nil {
		return nil, err
	}
	return parseStashList(out), nil
}

func parseStashList(out string) []StashEntry {
	if out == "" {
		return nil
	}

	lines := strings.Split(out, "\n")
	stashes := make([]StashEntry, 0, len(lines))
	for _, line := range lines {
		ref, message, ok := strings.Cut(line, "\x00")
		if !ok || ref == "" {
			continue
		}
		stashes = append(stashes, StashEntry{
			Ref:     ref,
			Message: message,
		})
	}
	return stashes
}

// StashPop applies and removes the stash matching ref.
func StashPop(ref string) error {
	_, err := Run("stash", "pop", ref)
	return err
}

// IsInsideRepo returns true when the current directory is inside a git repository.
func IsInsideRepo() bool {
	_, err := Run("rev-parse", "--git-dir")
	return err == nil
}

// RemoteURL returns the fetch URL for the named remote.
func RemoteURL(remote string) (string, error) {
	return Run("remote", "get-url", remote)
}

// DefaultBranch returns the remote's default branch by reading origin/HEAD.
// Falls back to checking for "main" then "master" if the ref is not set.
func DefaultBranch() (string, error) {
	out, err := Run("symbolic-ref", "--short", "refs/remotes/origin/HEAD")
	if err == nil {
		// Returns e.g. "origin/main" — strip the remote prefix.
		if _, branch, ok := strings.Cut(out, "/"); ok {
			return branch, nil
		}
		return out, nil
	}
	// Remote HEAD not configured; probe for well-known branch names.
	for _, candidate := range []string{"main", "master"} {
		if _, err := Run("show-ref", "--verify", "--quiet", "refs/remotes/origin/"+candidate); err == nil {
			return candidate, nil
		}
	}
	return "", fmt.Errorf("could not determine default branch: set base_branch in .code-cat.yml or run: git remote set-head origin --auto")
}
