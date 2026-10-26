# TUI

```sh
multienv-tui                 # Start in tmux mode
multienv-tui codex           # Start in Codex mode
multienv-tui doctor          # Open the diagnostics view
multienv-tui --no-popup      # Use the current pane instead of a popup
multienv-tui --help
```

Use `--no-popup` for the current pane and `--doctor` for a diagnostic report.
It reads `~/.config/multienv/tui.yaml`; `--config` overrides only this TUI file.
The agent configuration is always read from `~/.config/multienv/agent.yaml`.

```yaml
host_aliases:
  mac: macbook
  dev: devbox

popup:
  enabled: true
  width: 90%
  height: 90%

log_file: /tmp/multienv-tui.log

keybindings:
  switch-view: [ctrl+b]
  doctor: [ctrl+/, ctrl+_, f1]

tmux:
  control_port: 12400
  remote_pane_preview: false
  keybindings:
    create: [ctrl+n]
    rename: [ctrl+r]
    move-window: [ctrl+o]
    delete: [ctrl+x]

codex:
  ignore_file: ~/.config/multienv/.multienv-tui-ignore
  keybindings:
    rename: [ctrl+r]
    archive: [ctrl+a]
    unarchive: [ctrl+u]
    delete: [ctrl+x]
```

The TUI uses Go regular expressions in the ignore file.
Key bindings map action names to lists of lowercase key names. Omitted actions use their defaults;
an empty list disables an action. Root bindings are shared; `tmux`, `codex` and `doctor` can override them.
The control port defaults to `7777` and must match the launcher and forwarded port.
`popup.width` and `popup.height` accept percentages (`90%`) or cell counts (`120`).
Remote pane previews are off by default. The ignore and log paths support `~` and environment variables.

Run the dashboard inside tmux. Type to search, use arrow keys or Ctrl-J/Ctrl-K to
move, and press Enter to open the current item. Tab marks sessions for batch actions;
Enter always opens just the current item. `?` toggles the help header. Ctrl-B switches
between tmux and the optional Codex view. Ctrl-/ or F1 opens doctor directly.
The doctor view shows a searchable, colored report from `multienv-tui --doctor`;
Ctrl-R reruns it, and Esc, Ctrl-/ or F1 returns to the previous view.
Esc closes the dashboard when the session list has focus.
The list and preview refresh every three seconds.
Search matches only the session name (or the title in Codex mode) and the displayed host.
Ctrl-V toggles searching all visible fields, including the tmux command and directory;
the prompt shows `[all]` while this mode is active. Press Ctrl-V again to return to name and host search.

If tmux consumes Ctrl-B as its prefix, use its send-prefix binding to pass that key
to the dashboard; with the default bindings, press Ctrl-B twice.

## Launcher

The workstation launcher owns the interactive terminal connection. It accepts
`multienv SESSION [HOSTNAME]`; the host defaults to the current machine.
It reads aliases and `tmux.control_port` from `~/.config/multienv/tui.yaml`.
For a remote host, the alias must resolve through the workstation's SSH configuration.

Start it outside tmux:

```sh
tmux new-session -d -s main
multienv main
```

Inside tmux, run `multienv-tui`. A remote selection asks the launcher to attach locally or over SSH.
Keep the [control-port tunnel](forwarder.md) active on remote hosts so they can return to the workstation.
Local switching also works when attached to tmux directly.

## Tmux sessions and windows

Sessions are ordered by their last attachment time across all hosts, newest first.
Each row shows the host, window count, command and directory of the active window's
active pane.

- A green `●` means at least one client is attached.
- A red `!` means the snapshot is older than 30 seconds and takes precedence over the dot.
- The right pane lists windows by index, one line each, with a dot on the active window.
- The active pane's contents appear below the window list when available. The
  preview follows the selected window when the window list has focus and refreshes
  every three seconds. Remote previews are disabled by default; set
  `tmux.remote_pane_preview: true` in `tui.yaml` to enable them.

| Key | Action |
| --- | --- |
| Enter in the session list | Switch to the session's active window |
| Ctrl-L | Focus the windows on the right; retain the left list, selection and query |
| j/k, arrows, Ctrl-J/Ctrl-K in the window list | Choose a window |
| Enter in the window list | Switch directly to the chosen window |
| Ctrl-H or Esc in the window list | Return focus to the sessions |
| Ctrl-V in the session list | Toggle name/host search and search across all visible fields |
| Ctrl-N in the session list | Create a session; choose its host and name |
| Ctrl-N in the window list | Choose a directory and create a window with an automatic name |
| Ctrl-O in the window list | Move the selected window to a new index |
| Ctrl-R in the session list | Rename the current session |
| Ctrl-R in the window list | Rename the selected window |
| Ctrl-X in the session list | Kill marked sessions, or the current session if none are marked |
| Ctrl-X in the window list | Kill the selected window; the last window is protected |

The directory picker starts with the selected window’s current directory and includes
directories from zoxide when it is installed on that host. Esc cancels creation.

## Optional Codex integration

Synchronization uses the shared host list from [server configuration](server.md#configuration).
The conversation list combines stored snapshots with a fresh local snapshot.

Install the [Codex dependencies](requirements.md#optional-tools) on those hosts. The integration uses the same
agent, SSH connections and tmux sessions as the main workspace. Press Ctrl-B in the
dashboard, or start with `multienv-tui codex`, to open the conversation view.

The list shows an activity indicator, state, host and conversation title. A green `●`
means the session has a mapped tmux pane. A red `!` means the snapshot is older than
30 seconds and takes precedence over the dot. `Working (bg)` identifies active work
without a mapped tmux pane. The preview includes the working directory, first and
latest messages, and snapshot age. Conversations are ordered by state:
ActionRequired, Idle, Working, Not Loaded, Unknown, Archived.

| Key | Action |
| --- | --- |
| Enter | Switch to a running session's pane or prepare and resume an inactive session |
| Ctrl-R | Rename the current conversation |
| Ctrl-A | Archive marked sessions, or the current session |
| Ctrl-U | Unarchive marked sessions, or the current session |
| Ctrl-X | Delete marked sessions, or the current session |
| Ctrl-G | Show or hide archived sessions |
| Ctrl-T | Show or hide sessions matching the ignore file |
| Ctrl-V | Toggle title/host search and search across all visible fields |

When a session needs resuming, choose a tmux session on its host and Multienv creates
a window for it. With `--no-popup`, resuming a local session replaces the dashboard
in the current pane. Archive and delete operations also close an associated tmux
pane when one is found.

Archived and ignored sessions are hidden initially. Their visibility is shown in the title;
the toggles are independent.

## Ignore rules and logs

Codex ignore rules live in `~/.config/multienv/.multienv-tui-ignore`. Each nonempty
line is a case-sensitive Go regular expression; lines beginning with optional
whitespace and `#` are comments. Rules are reread on refresh. For example:

```text
# Hide background work
^Working \(bg\)\t
# Hide titles beginning with scratch:
^[^\t]*\t[^\t]*\tscratch:
```

Rules match three tab-separated fields: state, host alias and title.

Override the ignore file or log path in `tui.yaml`:

```yaml
log_file: /tmp/multienv-tui.log

codex:
  ignore_file: ~/.config/multienv/custom-ignore
```

TUI errors are appended to the configured log file, `/tmp/multienv-tui.log` by default.
Launcher logs are written to `/tmp/multienv-launcher.log`.
For stale remote sessions, check server logs, SSH access from the master, and
`multienv-agent tmux list-local` on the affected host.

