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

Installs to `/usr/local/bin` (will prompt for sudo if needed).

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
