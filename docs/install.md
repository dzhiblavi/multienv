# Installation and setup

Check the [requirements](requirements.md) first. The examples use `macbook` for the workstation
and `devbox` for a remote host. Replace them with each machine's actual `hostname` output.
Example HTTP and launcher ports are `12401` and `12400`.

## Install binaries

Clone the repository on each host. From the repository root, run the appropriate target:

```sh
mkdir -p ~/.bin ~/.config/multienv
make -C install local    # Workstation: server, agent, launcher and TUI
make -C install remote   # Other hosts: agent and TUI
```

Add `~/.bin` to `PATH` and install matching versions across hosts.
Keep the checkout in place: the launcher is a symlink to its script.
See [installer commands](../install/README.md) for individual targets and `INSTALL_DIR`.

## Place configuration files

| File | Host | Contents |
| --- | --- | --- |
| `~/.config/multienv/server.yaml` | Workstation | [Server and Redis connection](server.md#configuration) |
| `~/.config/multienv/redis-credentials.yaml` | Workstation | [Redis username and password](redis.md#connect-multienv) |
| `~/.config/multienv/agent.yaml` | Every host | [Server URL and request timeout](agent.md) |
| `~/.config/multienv/tui.yaml` | TUI hosts and launcher workstation | [Aliases, launcher port and UI options](tui.md) |
| `~/.config/multienv/forwarder.json` | Workstation, one per remote host | [SSH tunnels](forwarder.md) |

The launcher also reads `tui.yaml`. Keep its aliases resolvable through the workstation's SSH configuration,
and match the server URL and launcher port to the forwarded ports.

## Start the workspace

1. Make [Redis](redis.md) available to the workstation.
2. Check passwordless [SSH access](forwarder.md#ssh-access) from the workstation to the managed hosts.
3. Start the [server](server.md#start) and one [forwarder](forwarder.md#start) per remote host.
   On macOS, [launchd](launchd.md) can manage these processes.
4. Start the launcher outside tmux:

```sh
tmux new-session -d -s main
multienv main
```

Inside tmux, run `multienv-tui`. See the [TUI guide](tui.md) for navigation and shortcuts.
Run `multienv-agent doctor global` if a host cannot be reached.

## Uninstall

Stop the launcher, server and forwarder, then remove the installed programs:

```sh
rm -f ~/.bin/multienv ~/.bin/multienv-tui ~/.bin/multienv-agent ~/.bin/multienv-server
```

Redis data, tmux sessions, Codex conversations and configuration files remain in place.
