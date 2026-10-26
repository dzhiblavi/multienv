# Server

The server manages hosts, executes commands, provides JSON storage
and performs desktop actions on the master. Run one instance for the workspace.
See [requirements](requirements.md) and [architecture](architecture.md).

## Configuration

The default file is `~/.config/multienv/server.yaml`:

```yaml
server:
  address: 127.0.0.1
  port: 12401

hosts:
  - macbook
  - devbox

storage:
  address: 127.0.0.1
  port: 6379
  db: 0
  credentials_path: ~/.config/multienv/redis-credentials.yaml
```

Use actual hostnames in `hosts`; the server also uses these names for SSH.
`@master` can represent the server host in this list. If omitted, `hosts` defaults to the server's own hostname.
The `@all` invocation alias selects these hosts; duplicate targets are executed once.
The caller supplies the command, arguments and stdin. The server runs the command locally or over SSH
without interpreting integration-specific operations or data.
Generic command execution uses `POST /host` with action `invoke`; see [agent host commands](agent.md#host-commands).
JSON storage uses `POST /storage`; see [storage operations](storage.md).

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
The server checks Redis and runs `multienv-agent doctor local --format json` on itself and configured hosts.
This diagnostic command requires `multienv-agent` in PATH.
See [agent diagnostics](agent.md#diagnostics).
