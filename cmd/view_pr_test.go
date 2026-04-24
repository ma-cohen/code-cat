package cmd

import "testing"

func TestViewPRCmdRegistration(t *testing.T) {
	if viewPRCmd == nil {
		t.Fatal("viewPRCmd is nil")
	}
	if viewPRCmd.Use != "view-pr" {
		t.Errorf("viewPRCmd.Use = %q, want %q", viewPRCmd.Use, "view-pr")
	}
	if viewPRCmd.Short == "" {
		t.Error("viewPRCmd.Short should not be empty")
	}
	if viewPRCmd.RunE == nil {
		t.Error("viewPRCmd.RunE should not be nil")
	}
}

func TestViewPRCmdIsSubcommandOfRoot(t *testing.T) {
	found := false
	for _, cmd := range rootCmd.Commands() {
		if cmd.Use == "view-pr" {
			found = true
			break
		}
	}
	if !found {
		t.Error("view-pr command not registered as subcommand of root")
	}
}
