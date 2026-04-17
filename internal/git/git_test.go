package git

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

// makeTempRepo creates a temp dir, initializes a git repo with an initial commit,
// and changes the working directory into it. CWD is restored after the test.
func makeTempRepo(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()

	orig, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Chdir(orig) })

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
		if _, err := Run(args...); err != nil {
			t.Fatalf("setup git %v: %v", args, err)
		}
	}
	return dir
}

// makeTempRepoWithRemote creates a bare remote and a clone, sets up a default
// branch, and changes CWD into the clone. CWD is restored after the test.
func makeTempRepoWithRemote(t *testing.T, branch string) string {
	t.Helper()

	orig, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Chdir(orig) })

	remoteDir := t.TempDir()
	if out, err := exec.Command("git", "init", "--bare", remoteDir).CombinedOutput(); err != nil {
		t.Fatalf("git init --bare: %s", out)
	}

	cloneDir := t.TempDir()
	if out, err := exec.Command("git", "clone", remoteDir, cloneDir).CombinedOutput(); err != nil {
		t.Fatalf("git clone: %s", out)
	}

	if err := os.Chdir(cloneDir); err != nil {
		t.Fatal(err)
	}

	for _, args := range [][]string{
		{"config", "user.email", "test@test.com"},
		{"config", "user.name", "Test"},
		{"config", "commit.gpgsign", "false"},
		{"commit", "--allow-empty", "-m", "init"},
		{"push", "origin", "HEAD:" + branch},
		{"remote", "set-head", "origin", branch},
	} {
		if _, err := Run(args...); err != nil {
			t.Fatalf("setup git %v: %v", args, err)
		}
	}
	return cloneDir
}

func TestRun(t *testing.T) {
	makeTempRepo(t)

	t.Run("valid command returns output", func(t *testing.T) {
		out, err := Run("rev-parse", "--git-dir")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(out) == 0 {
			t.Error("expected non-empty output")
		}
	})

	t.Run("invalid subcommand returns error", func(t *testing.T) {
		_, err := Run("not-a-real-subcommand")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "git not-a-real-subcommand") {
			t.Errorf("error message %q does not contain expected prefix", err.Error())
		}
	})

	t.Run("output is trimmed", func(t *testing.T) {
		out, err := Run("rev-parse", "--abbrev-ref", "HEAD")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if out != strings.TrimSpace(out) {
			t.Errorf("output %q has leading/trailing whitespace", out)
		}
	})
}

func TestCurrentBranch(t *testing.T) {
	t.Run("returns branch name after init", func(t *testing.T) {
		makeTempRepo(t)
		branch, err := CurrentBranch()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(branch) == 0 {
			t.Error("expected non-empty branch name")
		}
	})

	t.Run("returns new branch after checkout", func(t *testing.T) {
		makeTempRepo(t)
		if _, err := Run("checkout", "-b", "my-feature"); err != nil {
			t.Fatalf("checkout: %v", err)
		}
		branch, err := CurrentBranch()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if branch != "my-feature" {
			t.Errorf("got %q, want %q", branch, "my-feature")
		}
	})
}

func TestHasUncommitted(t *testing.T) {
	t.Run("clean repo returns false", func(t *testing.T) {
		makeTempRepo(t)
		dirty, err := HasUncommitted()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if dirty {
			t.Error("expected clean repo")
		}
	})

	t.Run("untracked file returns true", func(t *testing.T) {
		makeTempRepo(t)
		if err := os.WriteFile("foo.txt", []byte("x"), 0644); err != nil {
			t.Fatal(err)
		}
		dirty, err := HasUncommitted()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !dirty {
			t.Error("expected dirty repo")
		}
	})

	t.Run("staged file returns true", func(t *testing.T) {
		makeTempRepo(t)
		if err := os.WriteFile("foo.txt", []byte("x"), 0644); err != nil {
			t.Fatal(err)
		}
		if _, err := Run("add", "foo.txt"); err != nil {
			t.Fatalf("git add: %v", err)
		}
		dirty, err := HasUncommitted()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !dirty {
			t.Error("expected dirty repo after staging")
		}
	})

	t.Run("committed file returns false", func(t *testing.T) {
		makeTempRepo(t)
		if err := os.WriteFile("foo.txt", []byte("x"), 0644); err != nil {
			t.Fatal(err)
		}
		for _, args := range [][]string{
			{"add", "foo.txt"},
			{"commit", "-m", "add foo"},
		} {
			if _, err := Run(args...); err != nil {
				t.Fatalf("git %v: %v", args, err)
			}
		}
		dirty, err := HasUncommitted()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if dirty {
			t.Error("expected clean repo after commit")
		}
	})
}

func TestIsInsideRepo(t *testing.T) {
	t.Run("returns true inside repo", func(t *testing.T) {
		makeTempRepo(t)
		if !IsInsideRepo() {
			t.Error("expected true inside git repo")
		}
	})

	t.Run("returns false outside repo", func(t *testing.T) {
		dir := t.TempDir()
		orig, err := os.Getwd()
		if err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() { os.Chdir(orig) })
		if err := os.Chdir(dir); err != nil {
			t.Fatal(err)
		}
		if IsInsideRepo() {
			t.Error("expected false outside git repo")
		}
	})
}

func TestDefaultBranch(t *testing.T) {
	t.Run("detects main", func(t *testing.T) {
		makeTempRepoWithRemote(t, "main")
		branch, err := DefaultBranch()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if branch != "main" {
			t.Errorf("got %q, want %q", branch, "main")
		}
	})

	t.Run("detects master", func(t *testing.T) {
		makeTempRepoWithRemote(t, "master")
		branch, err := DefaultBranch()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if branch != "master" {
			t.Errorf("got %q, want %q", branch, "master")
		}
	})

	t.Run("error when no remote", func(t *testing.T) {
		makeTempRepo(t)
		_, err := DefaultBranch()
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "could not determine default branch") {
			t.Errorf("error %q does not contain expected message", err.Error())
		}
	})
}
