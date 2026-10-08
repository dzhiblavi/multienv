# Multienv agent

`multienv-agent` provides three command trees:

```sh
multienv-agent tmux list-all
multienv-agent tmux sync-all
multienv-agent codex list-all
multienv-agent codex sync-all
printf '%s\n' '{"text":"Task finished","title":"Multienv"}' | multienv-agent master show-notification
printf '%s\n' '{"target":"https://example.com"}' | multienv-agent master open
```

The tmux and Codex trees retain their own session operations and synchronization commands.
The master tree executes commands on the server host. `show-notification` reads
a notification object from stdin: only `text` is required; `title`, `subtitle`,
and `sound_name` are optional. Successful notification commands produce no output.
`open` reads an object with a required `target` string from stdin and invokes
`/usr/bin/open` on the server host. The target can be a URL, file, or directory;
filesystem paths refer to the server host. Successful open commands produce no output.
Commands accepting parameters read JSON from stdin and return JSON where
applicable, as before.

All trees use `~/.config/multienv/agent.yaml`. Set `--config` (`-c`) to use
another file. A minimal configuration is:

```yaml
client:
  url: http://127.0.0.1:12123
  timeout: 2s
```

To migrate, copy the `client` section from the previous per-agent configuration
to `agent.yaml`. All command trees use the same server URL and timeout.

Install with `make -C deploy install-agent`. The server invokes
`~/.bin/multienv-agent tmux …` or `~/.bin/multienv-agent codex …` on target hosts.
Update the agent on every configured host together with the server and TUI.
HTTP endpoints `/tmux` and `/codex` and their request formats are unchanged.
Master commands use `/master` with the same `action`/`request` envelope.
