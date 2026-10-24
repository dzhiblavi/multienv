# Multienv

Multienv is a terminal workspace for moving between tmux sessions and windows across
multiple hosts. It helps you find a piece of work and return to it without remembering
which machine, SSH connection or window you left it in.

Open the dashboard inside tmux, search the combined session list, and press Enter to
switch. Names, commands and directories help you recognize each session; its windows
appear on the right. A workstation launcher handles connections between hosts, so
you can move between local and remote work from the same interface.

Optional Codex integration builds on that tmux workspace: find a conversation by
title or activity, jump to its running pane, or resume it in a tmux window. Desktop
commands also let you send notifications and open links on your workstation from
a remote shell.

## Getting started

Follow the [installation guide](docs/install.md) to install binaries and place configuration files.
Then start the launcher in your workstation terminal:

```sh
multienv main
```

Inside tmux, run `multienv-tui`. Type to search and press Enter to open a session.
Ctrl-B switches between tmux and Codex; Ctrl-/ or F1 opens diagnostics.
If something is not working, run `multienv-agent doctor global`.

## Documentation

| Guide | Contents |
| --- | --- |
| [Requirements](docs/requirements.md) | Running components and required tools |
| [Installation](docs/install.md) | Binary installation, configuration locations and first launch |
| [Agent](docs/agent.md) | Agent configuration, CLI operations and diagnostics |
| [Server](docs/server.md) | Host lists, storage and server startup |
| [TUI](docs/tui.md) | Launcher, configuration, key bindings and previews |
| [Architecture](docs/architecture.md) | Components and their interactions |
| [Forwarder](docs/forwarder.md) | SSH access, tunnels and reconnection settings |
| [Redis](docs/redis.md) | Example Redis instance using Docker Compose |
| [macOS launchd](docs/launchd.md) | Example user services for the server and forwarder |
