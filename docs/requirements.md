# Requirements

## Running components

Run one coordinating server and install the agent on every managed host, including the server host.
The examples call the workstation the **master**, with `macbook` and `devbox` as example hostnames.

| Component | Where | Startup guide |
| --- | --- | --- |
| Redis with JSON support | Reachable from the server | [Redis with Docker Compose](redis.md) |
| Multienv server | One host coordinating the workspace | [Server](server.md), [macOS launchd](launchd.md) |
| SSH server | Each remote host | [SSH access](forwarder.md#ssh-access) |
| SSH tunnels | Workstation to each remote host | [Forwarder](forwarder.md), [macOS launchd](launchd.md) |
| Launcher | Workstation terminal, outside tmux | [Launcher](tui.md#launcher) |
| tmux sessions | Hosts where work runs | [Installation and first launch](install.md) |

The launcher and its control-port tunnels enable switching between machines.
Local TUI switching also works when attached to tmux directly.
Agents run on demand; they are not background services. Open the TUI only when needed.

## Installed tools

| Host or feature | Requirements |
| --- | --- |
| Building binaries | Go 1.25.14+ and Make; see the [installer](../install/README.md) |
| Every managed host | tmux, `multienv-agent`, SSH client |
| Every TUI host | `multienv-tui`; tmux 3.3+ for popup support |
| Server host | SSH client and Redis access |
| Launcher workstation | Bash 4+ first in `PATH`, SSH client, jq, Mike Farah's `yq`, OpenBSD-compatible `nc` with `-k` |
| Forwarder host | Bash 4+, SSH client, jq |
| Remote hosts | SSH server allowing noninteractive login and TCP forwarding |

Built-in integrations request remote execution of `multienv-agent` through the server.
The noninteractive SSH environment must include it in PATH, along with tmux and any tools used by
enabled integrations. Global diagnostics also require the agent in PATH on the server host.

## Redis

Redis must support `JSON.GET`, `JSON.SET`, `JSON.DEL`, `JSON.TYPE`, `EVAL` and `EVALSHA`.
Transactions (`WATCH`, `UNWATCH`, `MULTI`, `EXEC`) are required for paths containing quotes or backslashes.
The account also needs `PING` and `COMMAND INFO` for connection and capability checks,
and read/write access to the JSON documents used by integrations.
The built-in integrations use `@tmux` and `@codex`.
Set the endpoint, database and credentials file in [server configuration](server.md#configuration).
See [Redis](redis.md) for an example instance and a JSON capability check.

## Optional tools

- Codex integration: `lsof`, a Codex CLI with `--no-daemon`, `resume`, `archive`, `unarchive` and `delete`,
  and compatible local session data. `CODEX_HOME` overrides the default `~/.codex` directory.
- Directory picker: zoxide on each host where its directory history should be available.
- Desktop notifications and opening links: macOS on the server host, `osascript` and `/usr/bin/open`.

If something is not working, run `multienv-agent doctor global` or `multienv-tui --doctor`.
See the [agent diagnostics](agent.md#diagnostics) and each command's `--help`.
