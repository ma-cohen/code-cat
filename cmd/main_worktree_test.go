package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/ma-cohen/code-cat/internal/git"
)

func TestMainWorktreeCmdRegistration(t *testing.T) {
	if mainWorktreeCmd == nil {
		t.Fatal("mainWorktreeCmd is nil")
	}
	if mainWorktreeCmd.Use != "main-worktree" {
		t.Errorf("mainWorktreeCmd.Use = %q, want %q", mainWorktreeCmd.Use, "main-worktree")
	}
	wantAliases := []string{"primary-worktree"}
	if !slices.Equal(mainWorktreeCmd.Aliases, wantAliases) {
		t.Errorf("Aliases = %v, want %v", mainWorktreeCmd.Aliases, wantAliases)
	}
	if mainWorktreeCmd.Short == "" {
		t.Error("mainWorktreeCmd.Short should not be empty")
	}
	if mainWorktreeCmd.RunE == nil {
		t.Error("mainWorktreeCmd.RunE should not be nil")
	}
}

func TestMainWorktreeCmdIsSubcommandOfRoot(t *testing.T) {
	found := false
	for _, cmd := range rootCmd.Commands() {
		if cmd.Use == "main-worktree" {
			found = true
			break
		}
	}
	if !found {
		t.Error("main-worktree command not registered as subcommand of root")
	}
}

func TestMainWorktreeCmdFlags(t *testing.T) {
	if mainWorktreeCmd.Flags().Lookup("print-path") == nil {
		t.Error("missing print-path flag")
	}
	if mainWorktreeCmd.Flags().Lookup("no-enter") == nil {
		t.Error("missing no-enter flag")
	}
}

func TestRunMainWorktreeOutsideRepoReturnsError(t *testing.T) {
	t.Cleanup(func() {
		rootCmd.SetArgs(nil)
		rootCmd.SetOut(os.Stdout)
		rootCmd.SetErr(os.Stderr)
	})

	orig, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(orig); err != nil {
			t.Errorf("restore cwd: %v", err)
		}
	})
	if err := os.Chdir(t.TempDir()); err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)
	rootCmd.SetArgs([]string{"main-worktree", "--no-enter"})

	err = rootCmd.Execute()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "not inside a git repository") {
		t.Errorf("error = %q, want repository error", err.Error())
	}
}

func TestRunMainWorktreePrintPathWhenInPrimary(t *testing.T) {
	t.Cleanup(func() {
		rootCmd.SetArgs(nil)
		rootCmd.SetOut(os.Stdout)
		rootCmd.SetErr(os.Stderr)
	})

	makeTempGitRepo(t)

	want, err := git.PrimaryWorktreePath()
	if err != nil {
		t.Fatal(err)
	}

	var outBuf, errBuf bytes.Buffer
	rootCmd.SetOut(&outBuf)
	rootCmd.SetErr(&errBuf)
	rootCmd.SetArgs([]string{"main-worktree", "--print-path", "--no-enter"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("Execute: %v", err)
	}

	got := strings.TrimSpace(outBuf.String())
	if filepath.Clean(got) != filepath.Clean(want) {
		t.Errorf("stdout = %q, want primary %q", got, want)
	}
	if !strings.Contains(errBuf.String(), "Already in primary worktree") {
		t.Errorf("stderr = %q, want message about already in primary", errBuf.String())
	}
}

func TestRunMainWorktreePrintPathWhenInLinkedWorktree(t *testing.T) {
	t.Cleanup(func() {
		rootCmd.SetArgs(nil)
		rootCmd.SetOut(os.Stdout)
		rootCmd.SetErr(os.Stderr)
	})

	makeTempGitRepo(t)

	want, err := git.PrimaryWorktreePath()
	if err != nil {
		t.Fatal(err)
	}

	wtDir := filepath.Join(t.TempDir(), "linked-wt")
	if _, err := git.Run("worktree", "add", "-b", "side-branch", wtDir); err != nil {
		t.Fatalf("worktree add: %v", err)
	}
	if err := os.Chdir(wtDir); err != nil {
		t.Fatal(err)
	}

	var outBuf, errBuf bytes.Buffer
	rootCmd.SetOut(&outBuf)
	rootCmd.SetErr(&errBuf)
	rootCmd.SetArgs([]string{"main-worktree", "--print-path", "--no-enter"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("Execute: %v", err)
	}

	got := strings.TrimSpace(outBuf.String())
	if filepath.Clean(got) != filepath.Clean(want) {
		t.Errorf("stdout = %q, want primary %q", got, want)
	}
	if !strings.Contains(errBuf.String(), "Primary worktree:") {
		t.Errorf("stderr = %q, want primary worktree hint", errBuf.String())
	}
}
