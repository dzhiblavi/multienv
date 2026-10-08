# Multienv agent

`multienv-agent` provides two command trees:

```sh
multienv-agent tmux list-all
multienv-agent tmux sync-all
multienv-agent codex list-all
multienv-agent codex sync-all
```

Each tree retains its own session operations and synchronization commands.
Commands accepting parameters read JSON from stdin and return JSON where
applicable, as before.

Both trees use `~/.config/multienv/agent.yaml`. Set `--config` (`-c`) to use
another file. A minimal configuration is:

```yaml
client:
  url: http://127.0.0.1:12123
  timeout: 2s
```

To migrate, copy the `client` section from the previous per-agent configuration
to `agent.yaml`. Both command trees now use the same server URL and timeout.

Install with `make -C deploy install-agent`. The server invokes
`~/.bin/multienv-agent tmux …` or `~/.bin/multienv-agent codex …` on target hosts.
Update the agent on every configured host together with the server and TUI.
HTTP endpoints `/tmux` and `/codex` and their request formats are unchanged.
