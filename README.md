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

Push the current branch to origin and open a pull request via the [GitHub CLI](https://cli.github.com/).

Requires `gh` to be installed and authenticated (`gh auth login`).

```
ccat pr                          # push + prompts for title and optional body
ccat pr --draft                  # create as draft PR
ccat pr --web                    # open the PR form in the browser
ccat pr --no-push                # skip pushing (branch already on origin)
ccat pr --base develop           # override base branch
```

## Install

### Windows

```powershell
go install github.com/ma-cohen/code-cat/cmd/ccat@latest
```

Then make sure Go's bin directory is on your PATH (one-time setup). Run this in PowerShell:

```powershell
[Environment]::SetEnvironmentVariable("PATH", $env:PATH + ";$env:USERPROFILE\go\bin", "User")
```

Restart your terminal and `ccat` will be available.

### macOS / Linux

```sh
go install github.com/ma-cohen/code-cat/cmd/ccat@latest
```

Add Go's bin directory to your PATH if not already set (add to `~/.bashrc` or `~/.zshrc`):

```sh
export PATH="$PATH:$(go env GOPATH)/bin"
```

### Direct download

Grab the binary for your platform from [GitHub Releases](https://github.com/ma-cohen/code-cat/releases), extract it, and place it somewhere on your PATH.

## Updating ccat

Run the same install command again — Go will fetch and install the latest version:

```sh
go install github.com/ma-cohen/code-cat/cmd/ccat@latest
```

## Configuration

Place a `.code-cat.yml` in your repo root to override defaults:

```yaml
base_branch: main        # default: master
branch_prefix: "feat/"   # prepended to branch names in prompts
worktree_root: "../wt"   # default parent directory for new worktrees
```

See `.code-cat.yml.example` for all options.
