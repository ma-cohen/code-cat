package cmd

import (
	"path/filepath"
	"reflect"
	"testing"

	"github.com/ma-cohen/code-cat/internal/git"
)

func TestResolveWorktreesToRemove(t *testing.T) {
	eligible := []git.WorktreeEntry{
		{Path: "/tmp/a", BranchShort: "a"},
		{Path: "/tmp/b", BranchShort: "b"},
	}
	labelA := git.WorktreeLabel(eligible[0])
	labelB := git.WorktreeLabel(eligible[1])
	allLabel := "All removable worktrees (2)"

	displayToPath := map[string]string{
		allLabel: "",
		labelA:   "/tmp/a",
		labelB:   "/tmp/b",
	}

	t.Run("all option expands", func(t *testing.T) {
		got := resolveWorktreesToRemove([]string{allLabel, labelA}, allLabel, eligible, displayToPath)
		want := []string{"/tmp/a", "/tmp/b"}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %#v want %#v", got, want)
		}
	})

	t.Run("individual paths", func(t *testing.T) {
		got := resolveWorktreesToRemove([]string{labelB}, allLabel, eligible, displayToPath)
		want := []string{"/tmp/b"}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %#v want %#v", got, want)
		}
	})

	t.Run("empty when nothing valid", func(t *testing.T) {
		got := resolveWorktreesToRemove([]string{}, allLabel, eligible, displayToPath)
		if len(got) != 0 {
			t.Errorf("got %#v", got)
		}
	})
}

func TestPathsEqual(t *testing.T) {
	if !pathsEqual("/foo/bar", filepath.Join("/foo", "bar")) {
		t.Fatal("expected equal")
	}
	if pathsEqual("/foo", "/bar") {
		t.Fatal("expected not equal")
	}
}
