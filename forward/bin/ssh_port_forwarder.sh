#!/usr/bin/env bash

set -u

log() {
    echo "$(date "+%m/%d/%Y %I:%M:%S %p") $@" >&2
}

usage() {
    echo "Usage: $0 --config CONFIG.json"
}

case "${1:-}" in
    -h|--help)
        usage
        exit 0
        ;;
    -c|--config)
        [[ $# == 2 ]] || { usage >&2; exit 2; }
        config=$2
        ;;
    *)
        [[ $# == 1 ]] || { usage >&2; exit 2; }
        config=$1
        ;;
esac

if ! command -v jq >/dev/null 2>&1; then
    log "jq is required"
    exit 1
fi

if ! jq -e '
    def positive_integer: type == "number" and . == floor and . > 0;
    (.ssh_host | type == "string" and length > 0 and (contains("\n") | not)) and
    (.forwards | type == "array" and length > 0) and
    (.connection // {} | type == "object") and
    all([
    .connection.connect_timeout_seconds,
    .connection.connection_attempts,
    .connection.server_alive_interval_seconds,
    .connection.server_alive_count_max,
    .connection.retry_delay_seconds
    ][]; . == null or positive_integer) and
    all(.forwards[];
    (.direction == "L" or .direction == "R") and
    (.spec | type == "string" and length > 0 and (contains("\n") | not))
)
' "$config" >/dev/null; then
    log "Invalid forwarder config: $config"
    exit 1
fi

ssh_host=$(jq -r '.ssh_host' "$config")
IFS=$'\t' read -r connect_timeout connection_attempts server_alive_interval \
    server_alive_count_max retry_delay < <(
    jq -r '[
    .connection.connect_timeout_seconds // 10,
    .connection.connection_attempts // 1,
    .connection.server_alive_interval_seconds // 10,
    .connection.server_alive_count_max // 3,
    .connection.retry_delay_seconds // 5
    ] | @tsv' "$config"
)
ssh_args=(
    -n -N -T
    -o BatchMode=yes
    -o "ConnectTimeout=$connect_timeout"
    -o "ConnectionAttempts=$connection_attempts"
    -o ExitOnForwardFailure=yes
    -o "ServerAliveInterval=$server_alive_interval"
    -o "ServerAliveCountMax=$server_alive_count_max"
    -o ControlMaster=no
    -o ControlPath=none
    -o ForwardAgent=no
)

while IFS= read -r direction && IFS= read -r spec; do
    ssh_args+=("-$direction" "$spec")
done < <(jq -r '.forwards[] | .direction, .spec' "$config")

child_pid=
stop() {
    trap - INT TERM
    if [[ -n $child_pid ]]; then
        kill "$child_pid" 2>/dev/null || true
        wait "$child_pid" 2>/dev/null || true
    fi
    exit 0
}
trap stop INT TERM

delay() {
    sleep "$retry_delay" &
    child_pid=$!
    wait "$child_pid" || true
    child_pid=
}

log "[info] script started"

while true; do
    log "[info] starting ssh"
    ssh "${ssh_args[@]}" -- "$ssh_host" &
    child_pid=$!

    log "[info] ssh started"
    wait "$child_pid"

    status=$?
    child_pid=

    log "[info] ssh exited with status $status; retrying in $retry_delay seconds"
    delay
done
