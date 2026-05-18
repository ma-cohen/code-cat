package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestPRCmdRegistration(t *testing.T) {
	if prCmd == nil {
		t.Fatal("prCmd is nil")
	}
	if prCmd.Use != "pr" {
		t.Errorf("prCmd.Use = %q, want %q", prCmd.Use, "pr")
	}
	if prCmd.Short == "" {
		t.Error("prCmd.Short should not be empty")
	}
	if prCmd.RunE == nil {
		t.Error("prCmd.RunE should not be nil")
	}
	if prCmd.Flags().HasFlags() {
		t.Error("prCmd should not define creation flags")
	}
}

func TestPRCmdIsSubcommandOfRoot(t *testing.T) {
	found := false
	for _, cmd := range rootCmd.Commands() {
		if cmd.Use == "pr" {
			found = true
			break
		}
	}
	if !found {
		t.Error("pr command not registered as subcommand of root")
	}
}

func TestRunPROpensGitHubPRView(t *testing.T) {
	makeTempGitRepo(t)
	runGit(t, "checkout", "-b", "feature/open-existing-pr")
	runGit(t, "remote", "add", "origin", "https://github.com/owner/repo.git")

	argsFile := writeFakeCLI(t, "gh")

	if err := runPR(nil, nil); err != nil {
		t.Fatalf("runPR: %v", err)
	}

	assertFakeCLIArgs(t, argsFile, []string{"pr", "view", "--web"})
}

func TestRunPROpensGitLabMRView(t *testing.T) {
	makeTempGitRepo(t)
	runGit(t, "checkout", "-b", "feature/open-existing-mr")
	runGit(t, "remote", "add", "origin", "https://gitlab.example.com/owner/repo.git")

	argsFile := writeFakeCLI(t, "glab")

	if err := runPR(nil, nil); err != nil {
		t.Fatalf("runPR: %v", err)
	}

	assertFakeCLIArgs(t, argsFile, []string{"mr", "view", "--web"})
}

func runGit(t *testing.T, args ...string) {
	t.Helper()

	cmd := exec.Command("git", args...)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
}

func writeFakeCLI(t *testing.T, name string) string {
	t.Helper()

	dir := t.TempDir()
	argsFile := filepath.Join(dir, name+".args")
	scriptPath := filepath.Join(dir, name)
	script := fmt.Sprintf("#!/bin/sh\nprintf '%%s\\n' \"$@\" > %q\n", argsFile)

	if err := os.WriteFile(scriptPath, []byte(script), 0755); err != nil {
		t.Fatalf("write fake %s: %v", name, err)
	}
	t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))

	return argsFile
}

func assertFakeCLIArgs(t *testing.T, argsFile string, want []string) {
	t.Helper()

	out, err := os.ReadFile(argsFile)
	if err != nil {
		t.Fatalf("read fake CLI args: %v", err)
	}

	got := strings.Fields(strings.TrimSpace(string(out)))
	if len(got) != len(want) {
		t.Fatalf("fake CLI args = %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("fake CLI arg %d = %q, want %q", i, got[i], want[i])
		}
	}
}
