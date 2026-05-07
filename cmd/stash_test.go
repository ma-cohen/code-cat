package cmd

import (
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/ma-cohen/code-cat/internal/git"
)

func makeTempGitRepo(t *testing.T) string {
	t.Helper()

	dir := t.TempDir()
	orig, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(orig); err != nil {
			t.Errorf("restore cwd: %v", err)
		}
	})

	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}

	for _, args := range [][]string{
		{"init"},
		{"config", "user.email", "test@test.com"},
		{"config", "user.name", "Test"},
		{"config", "commit.gpgsign", "false"},
		{"commit", "--allow-empty", "-m", "init"},
	} {
		cmd := exec.Command("git", args...)
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("setup git %v: %v\n%s", args, err, out)
		}
	}

	return dir
}

func TestStashCmdRegistration(t *testing.T) {
	if stashCmd == nil {
		t.Fatal("stashCmd is nil")
	}
	if stashCmd.Use != "stash [name]" {
		t.Errorf("stashCmd.Use = %q, want %q", stashCmd.Use, "stash [name]")
	}
	if stashCmd.Short == "" {
		t.Error("stashCmd.Short should not be empty")
	}
	if stashCmd.RunE == nil {
		t.Error("stashCmd.RunE should not be nil")
	}
}

func TestStashCmdIsSubcommandOfRoot(t *testing.T) {
	found := false
	for _, cmd := range rootCmd.Commands() {
		if cmd.Use == "stash [name]" {
			found = true
			break
		}
	}
	if !found {
		t.Error("stash command not registered as subcommand of root")
	}
}

func TestStashPopCmdRegistration(t *testing.T) {
	if stashPopCmd == nil {
		t.Fatal("stashPopCmd is nil")
	}
	if stashPopCmd.Use != "pop" {
		t.Errorf("stashPopCmd.Use = %q, want %q", stashPopCmd.Use, "pop")
	}
	if stashPopCmd.Short == "" {
		t.Error("stashPopCmd.Short should not be empty")
	}
	if stashPopCmd.RunE == nil {
		t.Error("stashPopCmd.RunE should not be nil")
	}

	found := false
	for _, cmd := range stashCmd.Commands() {
		if cmd.Use == "pop" {
			found = true
			break
		}
	}
	if !found {
		t.Error("stash pop command not registered as subcommand of stash")
	}
}

func TestRunStashOutsideRepoReturnsError(t *testing.T) {
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

	err = runStash(nil, []string{"named stash"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "not inside a git repository") {
		t.Errorf("error = %q, want repository error", err.Error())
	}
}

func TestRunStashWithNameStashesChanges(t *testing.T) {
	makeTempGitRepo(t)

	if err := os.WriteFile("work.txt", []byte("work in progress"), 0644); err != nil {
		t.Fatal(err)
	}

	if err := runStash(nil, []string{"named stash"}); err != nil {
		t.Fatalf("runStash: %v", err)
	}

	dirty, err := git.HasUncommitted()
	if err != nil {
		t.Fatalf("HasUncommitted: %v", err)
	}
	if dirty {
		t.Error("expected clean repo after runStash")
	}

	stashes, err := git.ListStashes()
	if err != nil {
		t.Fatalf("ListStashes: %v", err)
	}
	if len(stashes) != 1 {
		t.Fatalf("len(ListStashes()) = %d, want 1: %#v", len(stashes), stashes)
	}
	if !strings.Contains(stashes[0].Message, "named stash") {
		t.Errorf("stash message = %q, want it to contain %q", stashes[0].Message, "named stash")
	}
}

func TestRunStashWithNoChangesDoesNothing(t *testing.T) {
	makeTempGitRepo(t)

	if err := runStash(nil, []string{"unused name"}); err != nil {
		t.Fatalf("runStash: %v", err)
	}

	stashes, err := git.ListStashes()
	if err != nil {
		t.Fatalf("ListStashes: %v", err)
	}
	if len(stashes) != 0 {
		t.Errorf("len(ListStashes()) = %d, want 0: %#v", len(stashes), stashes)
	}
}

func TestRunStashPopWithNoStashesDoesNothing(t *testing.T) {
	makeTempGitRepo(t)

	if err := runStashPop(nil, nil); err != nil {
		t.Fatalf("runStashPop: %v", err)
	}
}
