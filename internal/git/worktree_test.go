package git

import (
	"path/filepath"
	"reflect"
	"testing"
)

func TestParseWorktreeListPorcelain(t *testing.T) {
	t.Run("multiple worktrees", func(t *testing.T) {
		input := `worktree /repo/main
HEAD deadbeef
branch refs/heads/main

worktree /repo/other
HEAD cafebabe
detached

worktree /repo/third
HEAD abc12345
branch refs/heads/feature/x
`

		got := parseWorktreeListPorcelain(input)
		want := []WorktreeEntry{
			{Path: "/repo/main", Head: "deadbeef", BranchRef: "refs/heads/main", BranchShort: "main"},
			{Path: "/repo/other", Head: "cafebabe", Detached: true},
			{Path: "/repo/third", Head: "abc12345", BranchRef: "refs/heads/feature/x", BranchShort: "feature/x"},
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("parseWorktreeListPorcelain(...) = %#v\nwant %#v", got, want)
		}
	})

	t.Run("bare worktree", func(t *testing.T) {
		input := `worktree /srv/bare.git
HEAD deadbeef
bare
`
		got := parseWorktreeListPorcelain(input)
		want := []WorktreeEntry{
			{Path: "/srv/bare.git", Head: "deadbeef", Bare: true},
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %#v want %#v", got, want)
		}
	})

	t.Run("empty", func(t *testing.T) {
		if got := parseWorktreeListPorcelain(""); got != nil {
			t.Errorf("got %#v want nil", got)
		}
		if got := parseWorktreeListPorcelain("\n  \n"); got != nil {
			t.Errorf("got %#v want nil", got)
		}
	})
}

func TestWorktreeLabel(t *testing.T) {
	t.Run("branch", func(t *testing.T) {
		s := WorktreeLabel(WorktreeEntry{Path: "/p", BranchShort: "feat"})
		if s != "/p  (feat)" {
			t.Errorf("got %q", s)
		}
	})
	t.Run("detached flag", func(t *testing.T) {
		s := WorktreeLabel(WorktreeEntry{Path: "/p", Detached: true})
		if s != "/p  (detached)" {
			t.Errorf("got %q", s)
		}
	})
	t.Run("bare", func(t *testing.T) {
		s := WorktreeLabel(WorktreeEntry{Path: "/b", Bare: true})
		if s != "/b  (bare)" {
			t.Errorf("got %q", s)
		}
	})
}

func TestListWorktrees(t *testing.T) {
	makeTempRepo(t)

	mainTop, err := WorktreeTopLevel()
	if err != nil {
		t.Fatal(err)
	}

	wtDir := filepath.Join(t.TempDir(), "extra-wt")
	if _, err := Run("worktree", "add", "-b", "side-branch", wtDir); err != nil {
		t.Fatalf("worktree add: %v", err)
	}

	list, err := ListWorktrees()
	if err != nil {
		t.Fatalf("ListWorktrees: %v", err)
	}
	if len(list) != 2 {
		t.Fatalf("len = %d, want 2: %#v", len(list), list)
	}

	var sawMain, sawSide bool
	for _, e := range list {
		if filepath.Clean(e.Path) == filepath.Clean(mainTop) {
			sawMain = true
		}
		if e.BranchShort == "side-branch" {
			sawSide = true
		}
	}
	if !sawMain || !sawSide {
		t.Errorf("expected main and side-branch entries: %#v", list)
	}
}
