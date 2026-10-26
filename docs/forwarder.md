# SSH forwarder

The forwarder maintains SSH port forwards and reconnects after connection failures.
Run one instance on the workstation for each remote host.

## SSH access

The workstation must be able to SSH to each managed host without a password prompt.
For example, in `~/.ssh/config` on the workstation:

```sshconfig
Host dev devbox
    HostName devbox.example.net
    User alice
```

Here `devbox` matches the remote machine's actual hostname, while `dev` is a friendly alias.
Check the connection and agent installation:

```sh
ssh -T -o BatchMode=yes devbox 'multienv-agent tmux list-local'
```

The remote SSH server must allow TCP forwarding.

## Configuration

Create `~/.config/multienv/forwarder.json` on the workstation:

```json
{
  "ssh_host": "devbox",
  "connection": {
    "connect_timeout_seconds": 10,
    "connection_attempts": 1,
    "server_alive_interval_seconds": 10,
    "server_alive_count_max": 3,
    "retry_delay_seconds": 5
  },
  "forwards": [
    {"direction": "R", "spec": "127.0.0.1:12401:127.0.0.1:12401"},
    {"direction": "R", "spec": "127.0.0.1:12400:127.0.0.1:12400"}
  ]
}
```

The first reverse forward exposes the workstation's HTTP server on the remote host's loopback interface.
The second exposes the workstation launcher, allowing a remote TUI to switch back to another machine.
Match the ports to [server](server.md) and [TUI/launcher](tui.md) configuration.
On the remote host, `agent.yaml` can then use `http://127.0.0.1:12401` as its server URL.

`connection` is optional; the values above are the defaults. All its values must be positive integers.
`ssh_host` accepts an SSH hostname or alias. `direction` accepts `L` or `R`; `spec` uses SSH forwarding syntax.
Use a separate configuration file for each remote host.

## Start

From the repository root:

```sh
bash forward/bin/ssh_port_forwarder.sh --config ~/.config/multienv/forwarder.json
```

The script logs to stderr. Ctrl-C stops it and its SSH child process.
On macOS, see the [launchd example](launchd.md). You can also maintain equivalent SSH tunnels yourself.
