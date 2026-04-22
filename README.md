# ccat — code-cat

Git workflow helpers. Stop doing the same steps by hand.

## Commands

### `ccat new-task [branch-name]`

Fetch latest base branch and check out a new feature branch.

```
ccat new-task                     # prompts for branch name
ccat new-task feat/my-thing       # use provided name directly
ccat new-task --base main         # override base branch
ccat new-task --push              # push to origin after creation
ccat new-task --no-fetch          # skip git fetch (useful offline)
```

### `ccat home`

Go to the latest base branch from origin, handling uncommitted changes along the way.

```
ccat home                 # goes to configured base branch
ccat home --base main     # override base branch
ccat home --no-fetch      # skip git fetch
```

### `ccat new-worktree [path]`

Create a git worktree in a separate directory on a new branch.

```
ccat new-worktree                              # prompts for path and branch
ccat new-worktree ../my-hotfix                 # prompts for branch name only
ccat new-worktree ../scratch --branch wip/exp  # fully specified
ccat new-worktree --base main                  # override base branch
```

### `ccat pr`

Push the current branch to origin and open a pull request (GitHub) or merge request (GitLab). The provider is auto-detected from the remote URL — no configuration needed.

Requires the matching CLI to be installed and authenticated:
- **GitHub**: [gh](https://cli.github.com/) — `gh auth login`
- **GitLab**: [glab](https://gitlab.com/gitlab-org/cli) — `glab auth login`

```
ccat pr                          # push + prompts for title and optional body
ccat pr --draft                  # create as draft PR/MR
ccat pr --web                    # open the PR/MR form in the browser
ccat pr --no-push                # skip pushing (branch already on origin)
ccat pr --base develop           # override base branch
```

## Install

### Windows

```powershell
irm https://raw.githubusercontent.com/ma-cohen/code-cat/main/install.ps1 | iex
```

Installs to `~\.local\bin` and adds it to your PATH automatically. Restart your terminal and `ccat` will be available.

### macOS / Linux

```sh
curl -fsSL https://raw.githubusercontent.com/ma-cohen/code-cat/main/install.sh | sh
```

Installs to `~/.local/bin` by default (no sudo). If that directory is not already on your `PATH`, the script prints a one-line `export` to add — use a new terminal afterward (same idea as the Windows installer).

**System-wide install** (optional): set `INSTALL_DIR` before running, for example:

```sh
INSTALL_DIR=/usr/local/bin curl -fsSL https://raw.githubusercontent.com/ma-cohen/code-cat/main/install.sh | sh
```

That may prompt for `sudo` when the target directory is not writable by you.

**Troubleshooting:** If `INSTALL_DIR` exists but is not traversable for normal users (for example `/usr/local/bin` with mode `700` owned by root), the binary can be installed successfully yet still be unusable until directory permissions are fixed, or until you install to a user-local path such as the default `~/.local/bin`.

## Updating ccat

Run the same install command again — the script always fetches the latest release.

## Configuration

Place a `.code-cat.yml` in your repo root to override defaults:

```yaml
base_branch: main        # default: master
branch_prefix: "feat/"   # prepended to branch names in prompts
worktree_root: "../wt"   # default parent directory for new worktrees
```

See `.code-cat.yml.example` for all options.

## Claude Code Plugin

Install the ccat plugin to give Claude agents in any repo instant awareness of all ccat commands:

```sh
claude skills add ma-cohen/code-cat/.claude-skill
```

After installation, type `/ccat` in a Claude Code session to inject the full command reference into the agent's context.

### Global install (all projects)

To make `/ccat` available in every Claude Code session regardless of project, copy the command file to your user-level commands directory:

**macOS / Linux**
```sh
mkdir -p ~/.claude/commands
curl -fsSL https://raw.githubusercontent.com/ma-cohen/code-cat/main/.claude-skill/commands/ccat.md \
  -o ~/.claude/commands/ccat.md
```

**Windows (PowerShell)**
```powershell
New-Item -ItemType Directory -Force "$HOME\.claude\commands" | Out-Null
Invoke-WebRequest https://raw.githubusercontent.com/ma-cohen/code-cat/main/.claude-skill/commands/ccat.md `
  -OutFile "$HOME\.claude\commands\ccat.md"
```

The `/ccat` command will then be available in all your projects without any per-repo setup.
