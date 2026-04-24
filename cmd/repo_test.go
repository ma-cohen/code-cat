package cmd

import "testing"

func TestRepoCmdRegistration(t *testing.T) {
	if repoCmd == nil {
		t.Fatal("repoCmd is nil")
	}
	if repoCmd.Use != "repo" {
		t.Errorf("repoCmd.Use = %q, want %q", repoCmd.Use, "repo")
	}
	if repoCmd.Short == "" {
		t.Error("repoCmd.Short should not be empty")
	}
	if repoCmd.RunE == nil {
		t.Error("repoCmd.RunE should not be nil")
	}
}

func TestRepoCmdIsSubcommandOfRoot(t *testing.T) {
	found := false
	for _, cmd := range rootCmd.Commands() {
		if cmd.Use == "repo" {
			found = true
			break
		}
	}
	if !found {
		t.Error("repo command not registered as subcommand of root")
	}
}
