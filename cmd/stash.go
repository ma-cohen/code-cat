package cmd

import (
	"fmt"

	"github.com/ma-cohen/code-cat/internal/git"
	"github.com/ma-cohen/code-cat/internal/prompt"
	"github.com/spf13/cobra"
)

var stashCmd = &cobra.Command{
	Use:   "stash [name]",
	Short: "Stash changes with a name or pop a selected stash",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runStash,
}

var stashPopCmd = &cobra.Command{
	Use:   "pop",
	Short: "Interactively choose a stash to pop",
	Args:  cobra.NoArgs,
	RunE:  runStashPop,
}

func init() {
	rootCmd.AddCommand(stashCmd)
	stashCmd.AddCommand(stashPopCmd)
}

func runStash(cmd *cobra.Command, args []string) error {
	if !git.IsInsideRepo() {
		return fmt.Errorf("not inside a git repository")
	}

	dirty, err := git.HasUncommitted()
	if err != nil {
		return err
	}
	if !dirty {
		fmt.Println("No changes to stash.")
		return nil
	}

	name := ""
	if len(args) > 0 {
		name = args[0]
	} else {
		name, err = prompt.AskString("Stash name", "")
		if err != nil {
			return err
		}
	}

	fmt.Printf("Stashing changes as %q...\n", name)
	if err := git.StashPush(name); err != nil {
		return err
	}
	fmt.Println("Changes stashed.")
	return nil
}

func runStashPop(cmd *cobra.Command, args []string) error {
	if !git.IsInsideRepo() {
		return fmt.Errorf("not inside a git repository")
	}

	stashes, err := git.ListStashes()
	if err != nil {
		return err
	}
	if len(stashes) == 0 {
		fmt.Println("No stashes to pop.")
		return nil
	}

	labels := make([]string, len(stashes))
	refsByLabel := make(map[string]string, len(stashes))
	for i, stash := range stashes {
		label := fmt.Sprintf("%s: %s", stash.Ref, stash.Message)
		labels[i] = label
		refsByLabel[label] = stash.Ref
	}

	choice, err := prompt.AskSelect("Choose a stash to pop", labels)
	if err != nil {
		return err
	}

	ref := refsByLabel[choice]
	fmt.Printf("Popping %s...\n", ref)
	if err := git.StashPop(ref); err != nil {
		return err
	}
	fmt.Println("Stash popped.")
	return nil
}
