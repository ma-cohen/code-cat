# ccat — Git Workflow CLI

`ccat` automates repetitive Git tasks: branching, worktrees, and repository navigation.

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

### `ccat stash [name]`

Stash current changes, including untracked files, with a message.

```
ccat stash                         # prompts for a stash name
ccat stash "wip auth cleanup"      # stash using the provided name
ccat stash pop                     # choose a stash interactively and pop it
```

Stash names are Git stash messages and are not guaranteed to be unique.

### `ccat new-worktree [path]`

Create a git worktree in a separate directory on a new branch.

```
ccat new-worktree                               # prompts for path and branch
ccat new-worktree ../my-hotfix                  # prompts for branch name only
ccat new-worktree ../scratch --branch wip/exp   # fully specified
ccat new-worktree --base main                   # override base branch
ccat new-worktree --no-fetch                    # skip git fetch
ccat new-worktree --print-path                  # stdout: only absolute path (use with cd "$(ccat ...)")
ccat new-worktree --no-enter                    # skip opening a shell in the new worktree (TTY only)
```

After creation, messages go to stderr. On an interactive TTY, `ccat` can open a subshell in the
new directory (decline with `--no-enter`). For a `cd` in the **current** shell, use `--print-path`.

### `ccat remove-worktree`

Remove linked worktrees interactively. Multi-select paths to remove, or choose **All removable
worktrees** to remove every linked worktree except the one you are in. You are always asked to
confirm before deleting. `--force` is passed to `git worktree remove` (for example when a
worktree has local changes).

```
ccat remove-worktree
ccat remove-worktree --force
```

### `ccat pr`

Open the pull request (GitHub) or merge request (GitLab) for the current branch in the browser. Provider is auto-detected from the remote URL.

Requires the matching CLI installed and authenticated:
- GitHub: `gh` — `gh auth login`
- GitLab: `glab` — `glab auth login`

```
ccat pr                     # opens the PR/MR for the current branch
```

If no PR/MR exists for the current branch, the provider CLI will display an error.

### `ccat repo`

Open the repository in the browser (GitHub or GitLab). Provider is auto-detected from the remote URL.

```
ccat repo                   # opens the repo homepage in your browser
```

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
- Use `ccat stash` to quickly save named work-in-progress changes, and `ccat stash pop` to restore one interactively.
- Use `ccat new-worktree` when you want to work on multiple branches simultaneously in separate directories.
- Use `ccat remove-worktree` to clean up extra linked worktrees (with an “all” option).
- Use `ccat pr` to open the existing PR/MR for the current branch.
