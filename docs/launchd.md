# Server and forwarder with launchd

These examples run the server and one SSH forwarder as user LaunchAgents on macOS.
They start at login and restart after exiting. See Apple's
[launchd guide](https://developer.apple.com/library/archive/documentation/MacOSX/Conceptual/BPSystemStartup/)
and the local `man launchd.plist` and `man launchctl` references.

First [install and configure Multienv](install.md), start Redis and check that the server and forwarder
work when run in a terminal. Stop those manual instances before loading the jobs.
The interactive launcher remains in your terminal.

Replace `/Users/alice` and `/Users/alice/src/multienv` below with your home and checkout paths.
Use absolute paths in plists: `~` and shell variables are not expanded.
The examples assume Homebrew's Bash is `/opt/homebrew/bin/bash`; use your actual Bash 4+ path.

```sh
mkdir -p ~/Library/LaunchAgents ~/Library/Logs/multienv
```

## Server

Create `~/Library/LaunchAgents/local.multienv.server.plist`:

```xml
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>local.multienv.server</string>
    <key>ProgramArguments</key>
    <array>
        <string>/Users/alice/.bin/multienv-server</string>
        <string>--config</string>
        <string>/Users/alice/.config/multienv/server.yaml</string>
    </array>
    <key>EnvironmentVariables</key>
    <dict>
        <key>PATH</key>
        <string>/Users/alice/.bin:/opt/homebrew/bin:/usr/local/bin:/usr/bin:/bin:/usr/sbin:/sbin</string>
    </dict>
    <key>RunAtLoad</key>
    <true/>
    <key>KeepAlive</key>
    <true/>
    <key>StandardOutPath</key>
    <string>/Users/alice/Library/Logs/multienv/server.log</string>
    <key>StandardErrorPath</key>
    <string>/Users/alice/Library/Logs/multienv/server.err.log</string>
</dict>
</plist>
```

## Forwarder

Create `~/Library/LaunchAgents/local.multienv.forward.plist`:

```xml
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>local.multienv.forward</string>
    <key>ProgramArguments</key>
    <array>
        <string>/opt/homebrew/bin/bash</string>
        <string>/Users/alice/src/multienv/forward/bin/ssh_port_forwarder.sh</string>
        <string>--config</string>
        <string>/Users/alice/.config/multienv/forwarder.json</string>
    </array>
    <key>EnvironmentVariables</key>
    <dict>
        <key>PATH</key>
        <string>/Users/alice/.bin:/opt/homebrew/bin:/usr/local/bin:/usr/bin:/bin:/usr/sbin:/sbin</string>
    </dict>
    <key>RunAtLoad</key>
    <true/>
    <key>KeepAlive</key>
    <true/>
    <key>StandardOutPath</key>
    <string>/Users/alice/Library/Logs/multienv/forward.log</string>
    <key>StandardErrorPath</key>
    <string>/Users/alice/Library/Logs/multienv/forward.err.log</string>
</dict>
</plist>
```

Create another plist with a distinct label and configuration path for each additional remote host.
LaunchAgents do not source shell startup files. Include tool directories in `PATH`, and add actual values for
`SSH_AUTH_SOCK` or `CODEX_HOME` to `EnvironmentVariables` if your setup relies on them.
Also make these tools and settings available in noninteractive SSH sessions on managed hosts.

## Load, inspect and stop

From a logged-in macOS user session:

```sh
plutil -lint ~/Library/LaunchAgents/local.multienv.server.plist
plutil -lint ~/Library/LaunchAgents/local.multienv.forward.plist
launchctl enable "gui/$(id -u)/local.multienv.server"
launchctl enable "gui/$(id -u)/local.multienv.forward"
launchctl bootstrap "gui/$(id -u)" ~/Library/LaunchAgents/local.multienv.server.plist
launchctl bootstrap "gui/$(id -u)" ~/Library/LaunchAgents/local.multienv.forward.plist
```

Inspect status and logs:

```sh
launchctl print "gui/$(id -u)/local.multienv.server"
launchctl print "gui/$(id -u)/local.multienv.forward"
tail -f ~/Library/Logs/multienv/server.log ~/Library/Logs/multienv/forward.err.log
```

After editing a plist, unload it and bootstrap it again. To restart a process without reloading its plist:

```sh
launchctl kickstart -k "gui/$(id -u)/local.multienv.server"
```

To stop and disable the jobs:

```sh
launchctl bootout "gui/$(id -u)/local.multienv.forward"
launchctl bootout "gui/$(id -u)/local.multienv.server"
launchctl disable "gui/$(id -u)/local.multienv.forward"
launchctl disable "gui/$(id -u)/local.multienv.server"
```
