package cmd

import (
	"fmt"
	"os"

	"github.com/ma-cohen/code-cat/internal/config"
	"github.com/spf13/cobra"
)

var version = "dev"

var rootCmd = &cobra.Command{
	Use:     "ccat",
	Short:   "code-cat — git workflow helpers",
	Version: version,
}

// topLevelCommands is the only place top-level subcommands are registered on the root.
// Keep this list in sync with new command files — nested commands use their parent's AddCommand.
// Tests enforce that rootCmd matches this slice so `ccat --help` cannot drift.
var topLevelCommands = []*cobra.Command{
	newTaskCmd,
	homeCmd,
	stashCmd,
	newWorktreeCmd,
	removeWorktreeCmd,
	prCmd,
	repoCmd,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(config.Load)
	for _, c := range topLevelCommands {
		rootCmd.AddCommand(c)
	}
}
