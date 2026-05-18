# code-cat

Go CLI tool (`ccat`) that automates Git workflows. Entry point: `cmd/ccat/main.go`.

All commands must work with both **GitHub** and **GitLab** (including self-hosted GitLab). Provider is auto-detected from the remote URL; no user configuration required.

## Command registration

Top-level subcommands (what appears under `ccat --help`) must be listed in the `topLevelCommands` slice in [`cmd/root.go`](cmd/root.go). Do not call `rootCmd.AddCommand` from individual command files—nested commands (e.g. `stash pop`) still attach via their parent in that command’s `init()`. [`cmd/root_test.go`](cmd/root_test.go) fails if the slice and the root command tree diverge.

## Testing

**After every bug fix or new feature**, spawn a sub-agent (Task tool) to write tests:

> "Write Go tests for the changes in [file(s)]. Use the standard `testing` package.
> Mirror patterns in the nearest existing `*_test.go` file. Run `go test -race -count=1 ./...`
> to verify all tests pass before finishing."

The sub-agent receives only the modified files — not the full conversation context.

Test command: `go test -race -count=1 ./...`

## Plugin maintenance

**When you add or modify a `ccat` command**, update `.claude-skill/commands/ccat.md` to reflect the change — new flags, new commands, or changed behavior. This keeps the `/ccat` slash command accurate for agents in repos that have the skill installed.
