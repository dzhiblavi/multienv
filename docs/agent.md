# Agent

The agent runs on demand; there is no agent daemon. Its default configuration is
`~/.config/multienv/agent.yaml`:

```yaml
client:
  url: http://127.0.0.1:12401
  timeout: 3s
```

`client.url` must reach the server, including through an [SSH tunnel](forwarder.md) on remote hosts.
The defaults are `http://127.0.0.1:12123` and a `2s` timeout.
For remote integration operations, the agent asks the server to invoke `multienv-agent` from PATH.
Make it, tmux and optional Codex tools available in noninteractive SSH sessions;
set `CODEX_HOME` there too if needed.

## CLI

Use `--config PATH` or `-c PATH` to read a different agent configuration:

```sh
multienv-agent [--config PATH] tmux COMMAND
multienv-agent [--config PATH] codex COMMAND
multienv-agent [--config PATH] master COMMAND
multienv-agent [--config PATH] filesystem COMMAND
multienv-agent [--config PATH] host invoke
multienv-agent [--config PATH] doctor [local|global]
multienv-agent help tmux
```

Commands in the `host`, `tmux`, `codex`, `master` and `filesystem` trees that take parameters read one JSON
object from stdin. Session queries and synchronization commands need no input.
Commands that return a result write JSON to stdout; logs and errors go to stderr.
Successful commands without a result produce no stdout output.

## Host commands

Run a command on a list of hosts; `@master` selects the server host and `@all`
selects the shared host list from [server configuration](server.md#configuration):

```sh
printf '%s\n' '{"hostnames":["@master","devbox"],"command":"hostname"}' | multienv-agent host invoke
```

The request accepts `command`, an optional `args` array, optional `stdin` text and an optional
per-host `timeout` such as `10s` (default `30s`). Commands run directly; invoke a shell explicitly
when needed. The agent's `client.timeout` still limits the entire HTTP request.

Aliases expand to actual hostnames; duplicate targets run once. Results follow the expanded host order.
The JSON response contains `results`, one per host, with `hostname`, `stdout`,
`stderr`, `exit_code` and an optional `error` for execution failures. Hosts run in parallel;
one host's failure does not discard other results. The CLI prints the response and returns
a nonzero exit status if any host fails.

## Diagnostics

If something is not working, run `multienv-agent doctor global`.
Use `multienv-agent doctor local` to check only the current host.
Both scopes accept `--format json` for scripts and return a nonzero exit status on failed checks.
See each command's `--help` for details.

## Session queries and synchronization

Both `tmux` and `codex` provide:

| Command | Purpose |
| --- | --- |
| `list-local` | Read this host's sessions, keyed by session ID; no server request |
| `list-all` | Combine server snapshots with a fresh local snapshot, grouped by hostname and ID |
| `sync-local` | Publish this host's snapshot to the server |
| `sync-all` | Invoke `sync-local` on all configured hosts through `@all` |

```sh
multienv-agent tmux list-all | jq
multienv-agent codex list-local | jq
multienv-agent tmux sync-all
```

`list-all` falls back to local sessions if the server cannot be reached. `sync-all`
can complete while individual hosts fail; check the command's stderr and snapshot ages.
Snapshots are read and published through [JSON storage](storage.md), under `@tmux` and `@codex`.

## Tmux operations

| Command | JSON fields |
| --- | --- |
| `pane-preview` | `hostname`, `target` |
| `new-session` | `hostname`, `name` |
| `rename-session` | `hostname`, `session_id`, `new_name` |
| `rename-window` | `hostname`, `session_id`, `window_id`, `new_name` |
| `kill-session` | `hostname`, `session_id` |
| `kill-window` | `hostname`, `session_id`, `window_id` |
| `new-window` | `hostname`, `session_id`; optional `name`, `work_dir`, `command` (an argument array) |
| `move-window` | `hostname`, `session_id`, `window_id`, `new_index` |

```sh
printf '%s\n' '{"hostname":"devbox","name":"project"}' | multienv-agent tmux new-session
printf '%s\n' '{"hostname":"devbox","session_id":"$1","new_name":"backend"}' \
    | multienv-agent tmux rename-session
```

Use the IDs returned by the queries. `new-session` returns the session object;
`new-window` returns `{"tmux_pane_id":"%..."}` for the new window's initial pane.

## Codex operations, optional

Each of these takes `hostname` and `session_id`:

| Command | Purpose |
| --- | --- |
| `rename-session` | Rename the conversation; also requires `new_name` |
| `archive-session` | Archive the conversation and close its mapped tmux pane |
| `unarchive-session` | Make the conversation available again |
| `delete-session` | Delete the conversation and close its mapped tmux pane |
| `prepare-resume` | Unarchive if needed, check readiness and directory, and return the session |

```sh
printf '%s\n' '{"hostname":"devbox","session_id":"YOUR_CODEX_SESSION_ID"}' \
    | multienv-agent codex prepare-resume
```

`prepare-resume` prepares the session; the caller starts the Codex process.
Renaming updates local Codex data; an already open Codex TUI shows the new name after reopening.

## Filesystem operations

| Command | JSON fields |
| --- | --- |
| `directory-candidates` | `hostname` |

## Master operations

These always run on the server host, regardless of where you invoke the agent.

| Command | JSON fields |
| --- | --- |
| `show-notification` | Required `text`; optional `title`, `subtitle`, `sound_name` |
| `open` | Required `target`: a URL, file or directory on the master |

```sh
printf '%s\n' '{"text":"Build finished","title":"devbox"}' \
    | multienv-agent master show-notification
printf '%s\n' '{"target":"https://example.com/build/123"}' | multienv-agent master open
```

These commands use macOS desktop notifications and `/usr/bin/open`. File paths
refer to the master, so remote files must be accessible there to open them.

