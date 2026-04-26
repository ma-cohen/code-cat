package cmd

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
)

func TestResetCmdRegistration(t *testing.T) {
	if resetCmd == nil {
		t.Fatal("resetCmd is nil")
	}
	if resetCmd.Use != "reset" {
		t.Errorf("resetCmd.Use = %q, want %q", resetCmd.Use, "reset")
	}
	if resetCmd.Short == "" {
		t.Error("resetCmd.Short should not be empty")
	}
	if resetCmd.RunE == nil {
		t.Error("resetCmd.RunE should not be nil")
	}
}

func TestResetCmdIsSubcommandOfRoot(t *testing.T) {
	found := false
	for _, cmd := range rootCmd.Commands() {
		if cmd.Use == "reset" {
			found = true
			break
		}
	}
	if !found {
		t.Error("reset command not registered as subcommand of root")
	}
}

func TestResetCmdForceFlag(t *testing.T) {
	flag := resetCmd.Flags().Lookup("force")
	if flag == nil {
		t.Fatal("force flag is not registered")
	}
	if flag.DefValue != "false" {
		t.Errorf("force flag default = %q, want %q", flag.DefValue, "false")
	}

	got, err := resetCmd.Flags().GetBool("force")
	if err != nil {
		t.Fatalf("GetBool(force): %v", err)
	}
	if got {
		t.Error("force flag should default to false")
	}
}

func TestRunResetForceRemovesLocalWorkAndKeepsIgnoredFiles(t *testing.T) {
	remoteDir := t.TempDir()
	runTestGit(t, "", "init", "--bare", remoteDir)

	cloneDir := t.TempDir()
	runTestGit(t, "", "clone", remoteDir, cloneDir)

	orig, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Chdir(orig) })
	if err := os.Chdir(cloneDir); err != nil {
		t.Fatal(err)
	}

	runTestGit(t, cloneDir, "config", "user.email", "test@test.com")
	runTestGit(t, cloneDir, "config", "user.name", "Test")
	runTestGit(t, cloneDir, "config", "commit.gpgsign", "false")

	writeTestFile(t, cloneDir, ".gitignore", "ignored.txt\n")
	writeTestFile(t, cloneDir, "tracked.txt", "base\n")
	runTestGit(t, cloneDir, "add", ".gitignore", "tracked.txt")
	runTestGit(t, cloneDir, "commit", "-m", "init")
	runTestGit(t, cloneDir, "push", "-u", "origin", "HEAD:main")

	writeTestFile(t, cloneDir, "committed.txt", "local commit\n")
	runTestGit(t, cloneDir, "add", "committed.txt")
	runTestGit(t, cloneDir, "commit", "-m", "local commit")

	writeTestFile(t, cloneDir, "tracked.txt", "staged change\n")
	runTestGit(t, cloneDir, "add", "tracked.txt")
	writeTestFile(t, cloneDir, "tracked.txt", "unstaged change\n")
	writeTestFile(t, cloneDir, "untracked.txt", "untracked\n")
	writeTestFile(t, cloneDir, "ignored.txt", "ignored\n")

	reset := &cobra.Command{}
	reset.Flags().Bool("force", true, "")
	if err := runReset(reset, nil); err != nil {
		t.Fatalf("runReset: %v", err)
	}

	assertFileContent(t, cloneDir, "tracked.txt", "base\n")
	assertFileMissing(t, cloneDir, "committed.txt")
	assertFileMissing(t, cloneDir, "untracked.txt")
	assertFileContent(t, cloneDir, "ignored.txt", "ignored\n")

	if status := runTestGit(t, cloneDir, "status", "--porcelain"); status != "" {
		t.Errorf("status = %q, want clean working tree", status)
	}
}

func runTestGit(t *testing.T, dir string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", args...)
	if dir != "" {
		cmd.Dir = dir
	}
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
	return string(out)
}

func writeTestFile(t *testing.T, dir, name, content string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
}

func assertFileContent(t *testing.T, dir, name, want string) {
	t.Helper()
	got, err := os.ReadFile(filepath.Join(dir, name))
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != want {
		t.Errorf("%s = %q, want %q", name, got, want)
	}
}

func assertFileMissing(t *testing.T, dir, name string) {
	t.Helper()
	if _, err := os.Stat(filepath.Join(dir, name)); !os.IsNotExist(err) {
		t.Errorf("%s exists, want missing", name)
	}
}
