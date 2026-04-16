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

### `ccat new-worktree [path]`

Create a git worktree in a separate directory on a new branch.

```
ccat new-worktree                              # prompts for path and branch
ccat new-worktree ../my-hotfix                 # prompts for branch name only
ccat new-worktree ../scratch --branch wip/exp  # fully specified
ccat new-worktree --base main                  # override base branch
```

## Install

**Go users:**
```sh
go install github.com/ma-cohen/code-cat@latest
```

**macOS / Linux (Homebrew):**
```sh
brew install ma-cohen/tap/ccat
```

**Windows (Scoop):**
```sh
scoop bucket add ma-cohen https://github.com/ma-cohen/scoop-bucket
scoop install ccat
```

**Direct download:** grab the binary for your platform from [GitHub Releases](https://github.com/ma-cohen/code-cat/releases) and put it on your PATH.

## Configuration

Place a `.code-cat.yml` in your repo root to override defaults:

```yaml
base_branch: main        # default: master
branch_prefix: "feat/"   # prepended to branch names in prompts
worktree_root: "../wt"   # default parent directory for new worktrees
```

See `.code-cat.yml.example` for all options.
