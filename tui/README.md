# TUI session filters

In Codex mode, **Ctrl-T** shows or hides sessions matching the ignore file.
Matching sessions are hidden when the TUI starts. The prompt shows the current
setting, including when the help header is hidden. **Ctrl-G** independently
shows or hides archived sessions; showing ignored sessions does not show hidden
archives.

The default prompt is `(-archived,+filtered) ❯`: `+archived` means archived
sessions are shown, and `+filtered` means ignore rules are applied. Each toggle
updates its corresponding `+`/`-` indicator.

The default ignore file is `~/.config/multienv/.multienv-tui-ignore`.
Override its path in `~/.config/multienv/tui.sh`:

```bash
c_tui_ignore_file_path="${HOME}/.multienv-tui-ignore"
```

Each nonempty line is a case-sensitive regular expression, using jq's regex
syntax. Lines starting with `#` (optionally preceded by whitespace) are comments.
A session is ignored if any rule matches. Rules are read on each list refresh;
a missing file means no rules. Invalid expressions produce an error rather than
being silently skipped.

Rules match the complete row passed to fzf, with ANSI color codes removed. In
Codex mode its six fields are separated by actual tabs:

```text
status<TAB>host alias<TAB>title (or session ID)<TAB>hostname<TAB>session ID<TAB>base64 session JSON
```

Both `\t` in a regex and literal tabs match the separators. Only the first three
fields are displayed by fzf. Anchor a rule to the desired fields to avoid
matching elsewhere in the row. For example, to hide background work and titles
starting with `ailoops:`:

```text
# Working sessions without a tmux pane
^Working \(bg\)\t
# Titles starting with ailoops:, in any state
^[^\t]*\t[^\t]*\tailoops:
```

`Working (bg)` is a purple TUI label for `Working` sessions without
`tmux_pane_id`. It is also used in the preview and in the text matched by rules;
the session's underlying state remains `Working`. Regex parentheses must be
escaped when matching this label literally, as in the example above.

The preview shows `Synced` as a relative age (for example, `3 seconds ago` or
`2 minutes ago`). Values older than 10 seconds are red. The age refreshes every
two seconds, including when the selected session has not changed.

Tmux sessions are sorted by their last attachment time, newest first, across
all hosts. This order is preserved when filtering. Sessions without an
attachment time appear last. A green `●` before the name indicates one or more
attached clients. If the snapshot is older than 30 seconds, a red `!` replaces
the attachment indicator, including for detached sessions. Indicators refresh
every two seconds. Each row also shows the window count, command and directory
of the active pane in the active window. Home directories are abbreviated to
`~`, and long paths are shortened to their last 47 characters with a leading `…`.

The right pane shows the selected tmux session's windows, ordered by index,
with their names, commands and directories. A green `●` marks the active window.
Each window occupies one line; the command is omitted when it matches the window
name. Lines wider than the pane are truncated.
Press **Ctrl-L** to focus the windows on the right without changing the session
list, its selection or search query. Use **j/k**, arrow keys or **Ctrl-J/Ctrl-K**
to select a window, then **Enter** to switch directly to it. **Ctrl-H** or **Esc**
returns to session selection. The windows refresh every two seconds. **Enter**
in the session list still switches to the session's active window. Session
management shortcuts are disabled while the windows have focus.
