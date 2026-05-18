package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestTopLevelCommandsSyncedWithRoot(t *testing.T) {
	got := rootCmd.Commands()
	want := make(map[*cobra.Command]struct{}, len(topLevelCommands))
	for _, c := range topLevelCommands {
		want[c] = struct{}{}
	}
	for _, c := range got {
		if _, ok := want[c]; ok {
			delete(want, c)
			continue
		}
		// After rootCmd.Execute(), cobra injects default help and completion commands.
		if c.Name() == "help" || c.Name() == "completion" {
			continue
		}
		t.Errorf("rootCmd has unknown command %q (not in topLevelCommands)", c.Name())
	}
	for c := range want {
		t.Errorf("topLevelCommands includes %q but it is not on rootCmd", c.Name())
	}
}

func TestHelpListsEveryTopLevelCommand(t *testing.T) {
	t.Cleanup(func() {
		rootCmd.SetArgs(nil)
	})
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"--help"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("help: %v", err)
	}
	out := buf.String()
	for _, c := range topLevelCommands {
		prefix := "\n  " + c.Name() + " "
		if !strings.Contains(out, prefix) {
			t.Errorf("help output missing subcommand line for %q", c.Name())
		}
	}
}
