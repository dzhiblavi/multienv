# Installation

Requires Go 1.25.14+ and Make. Run commands from the repository root.
See the [setup guide](../docs/install.md) for configuration and first launch.

```sh
mkdir -p ~/.bin
make -C install local    # Workstation: server, agent, launcher and TUI
make -C install remote   # Other hosts: agent and TUI
```

Choose the appropriate target for each host. Binaries are built into `~/.bin`;
`multienv` is a symlink to the launcher script, so keep the checkout in place.
Add `~/.bin` to `PATH`. Override the destination with `INSTALL_DIR=/path/to/bin`.

All targets below use `make -C install TARGET`:

| Target | Action |
| --- | --- |
| `local` | Install server, agent, launcher and TUI |
| `remote` | Install agent and TUI |
| `install-server` | Build `multienv-server` |
| `install-agent` | Build `multienv-agent` |
| `install-tui` | Build `multienv-tui` |
| `install-launcher` | Link `multienv` to the launcher script |
