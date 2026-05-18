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

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(config.Load)
}
