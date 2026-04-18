# ccat — Git Workflow CLI

`ccat` automates repetitive Git tasks: branching, worktrees, and pull requests.

## Commands

### `ccat new-task [branch-name]`

Fetch the latest base branch and check out a new feature branch.

```
ccat new-task                       # prompts for branch name
ccat new-task feat/my-thing         # use provided name directly
ccat new-task --base main           # override base branch
ccat new-task --push                # push new branch to origin after creation
ccat new-task --no-fetch            # skip git fetch (useful offline)
```

Warns if there are uncommitted changes before switching.

### `ccat home`

Return to the latest base branch from origin, handling uncommitted changes.

```
ccat home                   # goes to configured base branch
ccat home --base main       # override base branch
ccat home --no-fetch        # skip git fetch
```

When uncommitted changes exist, offers: **Stash**, **Discard**, or **Abort**.

### `ccat new-worktree [path]`

Create a git worktree in a separate directory on a new branch.

```
ccat new-worktree                               # prompts for path and branch
ccat new-worktree ../my-hotfix                  # prompts for branch name only
ccat new-worktree ../scratch --branch wip/exp   # fully specified
ccat new-worktree --base main                   # override base branch
ccat new-worktree --no-fetch                    # skip git fetch
```

Prints the absolute path and branch name when done.

### `ccat pr`

Push the current branch to origin and create a pull request via the GitHub CLI (`gh`).

Requires `gh` installed and authenticated (`gh auth login`).

```
ccat pr                     # push + prompts for title and optional body
ccat pr --draft             # create as draft PR
ccat pr --web               # open the PR form in the browser
ccat pr --no-push           # skip pushing (branch already on origin)
ccat pr --base develop      # override base branch
```

Auto-generates a PR title from the branch name (strips prefixes like `feat/`, `fix/`, `chore/`, etc.).

## Configuration

Place a `.code-cat.yml` in the repo root to override defaults:

```yaml
base_branch: main        # default: auto-detected from origin/HEAD, then main/master
branch_prefix: "feat/"   # prepended to branch names in interactive prompts
worktree_root: "../wt"   # default parent directory for new worktrees (default: ..)
```

Config precedence: repo-local `.code-cat.yml` > user-global `~/.config/code-cat/config.yml` > built-in defaults.

## Agent guidance

- Use `ccat new-task` to start any new piece of work — it keeps your base branch fresh.
- Use `ccat home` when you need to get back to the base branch cleanly.
- Use `ccat new-worktree` when you want to work on multiple branches simultaneously in separate directories.
- Use `ccat pr` as the final step when a branch is ready for review.
