# Server

The server stores host snapshots in Redis, coordinates synchronization, forwards agent commands
and performs desktop actions on the master. Run one instance for the workspace.
See [requirements](requirements.md) and [architecture](architecture.md).

## Configuration

The default file is `~/.config/multienv/server.yaml`:

```yaml
server:
  address: 127.0.0.1
  port: 12401

tmux:
  hosts:
    - macbook
    - devbox

codex:
  hosts: []

storage:
  address: 127.0.0.1
  port: 6379
  db: 0
  credentials_path: ~/.config/multienv/redis-credentials.yaml
```

Use actual hostnames in the host lists; the server also uses these names for SSH.
It invokes `~/.bin/multienv-agent` on remote hosts.
`tmux.hosts` and `codex.hosts` select the machines collected for each domain.
Set `codex.hosts` to an empty list to disable Codex collection; add its hosts to enable it.
If omitted, both host lists default to the server's own hostname.

The HTTP endpoint defaults to `127.0.0.1:12123`; this example explicitly uses `12401`.
The API has no authentication. Bind it to loopback and use [SSH tunnels](forwarder.md) for other hosts.
Always specify the Redis port and credentials file. See [Redis connection](redis.md#connect-multienv).

## Start

```sh
multienv-server
multienv-server --config /path/to/server.yaml
```

Choose one invocation. Logs go to stdout. Stop with Ctrl-C, or manage the server with
[launchd on macOS](launchd.md) or your own systemd service.
Redis must be reachable when the server starts.

For diagnostics, run `multienv-agent doctor global` from a configured host.
See [agent diagnostics](agent.md#diagnostics).
