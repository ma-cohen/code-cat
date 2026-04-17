package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
)

func writeYAML(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
}

// cdTemp changes CWD to a fresh temp dir and restores it after the test.
func cdTemp(t *testing.T) string {
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
	return dir
}

// resetViper clears all viper state between tests.
func resetViper(t *testing.T) {
	t.Helper()
	t.Cleanup(viper.Reset)
}

func TestLoad_Defaults(t *testing.T) {
	resetViper(t)
	cdTemp(t)
	t.Setenv("HOME", t.TempDir())

	Load()

	if C.BaseBranch != "" {
		t.Errorf("BaseBranch: got %q, want %q", C.BaseBranch, "")
	}
	if C.BranchPrefix != "" {
		t.Errorf("BranchPrefix: got %q, want %q", C.BranchPrefix, "")
	}
	if C.WorktreeRoot != ".." {
		t.Errorf("WorktreeRoot: got %q, want %q", C.WorktreeRoot, "..")
	}
}

func TestLoad_RepoLocalConfig(t *testing.T) {
	resetViper(t)
	dir := cdTemp(t)
	t.Setenv("HOME", t.TempDir())

	writeYAML(t, filepath.Join(dir, ".code-cat.yml"), "base_branch: develop\nbranch_prefix: \"feature/\"\n")

	Load()

	if C.BaseBranch != "develop" {
		t.Errorf("BaseBranch: got %q, want %q", C.BaseBranch, "develop")
	}
	if C.BranchPrefix != "feature/" {
		t.Errorf("BranchPrefix: got %q, want %q", C.BranchPrefix, "feature/")
	}
	if C.WorktreeRoot != ".." {
		t.Errorf("WorktreeRoot: got %q, want %q", C.WorktreeRoot, "..")
	}
}

func TestLoad_UserGlobalConfig(t *testing.T) {
	resetViper(t)
	cdTemp(t)
	home := t.TempDir()
	t.Setenv("HOME", home)

	writeYAML(t, filepath.Join(home, ".config", "code-cat", "config.yml"), "worktree_root: /tmp/wt\n")

	Load()

	if C.WorktreeRoot != "/tmp/wt" {
		t.Errorf("WorktreeRoot: got %q, want %q", C.WorktreeRoot, "/tmp/wt")
	}
}

func TestLoad_RepoOverridesUser(t *testing.T) {
	resetViper(t)
	dir := cdTemp(t)
	home := t.TempDir()
	t.Setenv("HOME", home)

	writeYAML(t, filepath.Join(home, ".config", "code-cat", "config.yml"), "base_branch: main\nworktree_root: /tmp/wt\n")
	writeYAML(t, filepath.Join(dir, ".code-cat.yml"), "base_branch: develop\n")

	Load()

	if C.BaseBranch != "develop" {
		t.Errorf("BaseBranch: got %q, want %q (repo config should win)", C.BaseBranch, "develop")
	}
	if C.WorktreeRoot != "/tmp/wt" {
		t.Errorf("WorktreeRoot: got %q, want %q (user value should be preserved)", C.WorktreeRoot, "/tmp/wt")
	}
}

func TestLoad_EmptyConfigFile(t *testing.T) {
	resetViper(t)
	dir := cdTemp(t)
	t.Setenv("HOME", t.TempDir())

	writeYAML(t, filepath.Join(dir, ".code-cat.yml"), "")

	Load()

	if C.WorktreeRoot != ".." {
		t.Errorf("WorktreeRoot: got %q, want %q", C.WorktreeRoot, "..")
	}
}
