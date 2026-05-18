package cmd

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func prepareNestedRepoWithOrigin(t *testing.T, remoteBranch string) string {
	t.Helper()

	orig, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(orig); err != nil {
			t.Errorf("restore cwd: %v", err)
		}
	})

	bare := t.TempDir()
	if out, err := exec.Command("git", "init", "--bare", bare).CombinedOutput(); err != nil {
		t.Fatalf("git init --bare: %v\n%s", err, out)
	}
	cloneRoot := t.TempDir()
	if out, err := exec.Command("git", "clone", bare, cloneRoot).CombinedOutput(); err != nil {
		t.Fatalf("git clone: %v\n%s", err, out)
	}

	runGitInDir(t, cloneRoot, "config", "user.email", "t@t.com")
	runGitInDir(t, cloneRoot, "config", "user.name", "T")
	runGitInDir(t, cloneRoot, "config", "commit.gpgsign", "false")
	runGitInDir(t, cloneRoot, "commit", "--allow-empty", "-m", "init")
	runGitInDir(t, cloneRoot, "push", "origin", "HEAD:"+remoteBranch)
	runGitInDir(t, cloneRoot, "remote", "set-head", "origin", remoteBranch)

	nested := filepath.Join(cloneRoot, "deep", "nested")
	if err := os.MkdirAll(nested, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(nested); err != nil {
		t.Fatal(err)
	}
	return cloneRoot
}

func runGitInDir(t *testing.T, dir string, args ...string) {
	t.Helper()
	c := exec.Command("git", args...)
	c.Dir = dir
	if out, err := c.CombinedOutput(); err != nil {
		t.Fatalf("git %v in %s: %v\n%s", args, dir, err, out)
	}
}

func TestRunNewWorktreeFromNestedDir(t *testing.T) {
	t.Cleanup(func() {
		rootCmd.SetArgs(nil)
		rootCmd.SetOut(os.Stdout)
		rootCmd.SetErr(os.Stderr)
	})

	prepareNestedRepoWithOrigin(t, "main")

	wtDir := filepath.Join(t.TempDir(), "linked-wt")
	var stderr bytes.Buffer
	rootCmd.SetOut(io.Discard)
	rootCmd.SetErr(&stderr)
	rootCmd.SetArgs([]string{"new-worktree", wtDir, "--branch", "feat-side", "--no-fetch", "--no-enter"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("Execute: %v\nstderr:\n%s", err, stderr.String())
	}

	runGitInDir(t, wtDir, "rev-parse", "--is-inside-work-tree")
	out, err := exec.Command("git", "-C", wtDir, "rev-parse", "--abbrev-ref", "HEAD").CombinedOutput()
	if err != nil {
		t.Fatalf("HEAD: %v", err)
	}
	if strings.TrimSpace(string(out)) != "feat-side" {
		t.Fatalf("HEAD = %s", out)
	}
}

func TestNewWorktreeCmdRegistration(t *testing.T) {
	if newWorktreeCmd == nil {
		t.Fatal("newWorktreeCmd is nil")
	}
	if newWorktreeCmd.Use != "new-worktree [path]" {
		t.Errorf("Use = %q", newWorktreeCmd.Use)
	}
	if newWorktreeCmd.Short == "" {
		t.Error("Short should not be empty")
	}
	if newWorktreeCmd.Long == "" {
		t.Error("Long should not be empty")
	}
	if newWorktreeCmd.Flags().Lookup("base") == nil {
		t.Error("missing base flag")
	}
}
