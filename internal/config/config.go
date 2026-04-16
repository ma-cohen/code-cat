package config

import (
	"github.com/spf13/viper"
)

// C holds the resolved configuration for the current invocation.
// Commands read from this after Load() has been called.
var C Config

type Config struct {
	BaseBranch   string `mapstructure:"base_branch"`
	BranchPrefix string `mapstructure:"branch_prefix"`
	WorktreeRoot string `mapstructure:"worktree_root"`
}

// Load reads .code-cat.yml from the current directory (repo-local config) and
// ~/.config/code-cat/config.yml (user-global config), merges them, and
// populates C. Repo config wins over user config wins over built-in defaults.
func Load() {
	viper.SetDefault("base_branch", "")
	viper.SetDefault("branch_prefix", "")
	viper.SetDefault("worktree_root", "..")

	// User-global config
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME/.config/code-cat")
	_ = viper.ReadInConfig()

	// Repo-local config overrides user config
	viper.SetConfigName(".code-cat")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	_ = viper.MergeInConfig()

	_ = viper.Unmarshal(&C)
}
