#!/usr/bin/env bash

lib_root="$(realpath $(dirname ${BASH_SOURCE}))"

log() {
    date="$(date "+%m/%d/%Y %I:%M:%S %p")"
    caller_module="$(basename "${BASH_SOURCE[1]:-$0}")"
    caller_func="${FUNCNAME[1]:-main}"

    echo "[$date] ($caller_module/$caller_func) $*" >>"${c_log_file_path}"
}

die() {
    log "DIE: " "$@"
    echo "DIE: " "$@" >&2
    exit 1
}

send() {
    local IFS=$'\t'

    if printf '%s\n' "$*" | nc -w 0 127.0.0.1 "${c_ctrl_port}"; then
        log "status=success args='$@'"
    else
        log "status=err args='$@'"
        return 1
    fi
}

is_in_tmux() {
    [[ -n "${TMUX:-}" ]]
}

is_on_remote() {
    [[ -n "${SSH_TTY:-}" ]]
}

is_on_hostname() {
    [[ "$(hostname)" == "${1:?need target hostname}" ]]
}

quote_shell_args() {
    local arg
    for arg in "$@"; do
        printf " %q" "$arg"
    done
}

run_remote() {
    local target_hostname="${1:?Need hostname}"
    local host_alias="${c_hostname_to_host_alias[$target_hostname]:-$target_hostname}"
    shift

    command="$1"
    shift

    log "host=$host_alias command='$command' args='$@'"
    ssh -n "$host_alias" "${command}" "$(quote_shell_args "$@")"
}

run_remote_interactive() {
    local target_hostname="${1:?Need hostname}"
    local host_alias="${c_hostname_to_host_alias[$target_hostname]:-$target_hostname}"
    shift

    command="$1"
    shift

    log "host=$host_alias command='$command' args='$@'"
    ssh -t "$host_alias" "${command}" "$(quote_shell_args "$@")"
}

_real_tmux=$(which tmux)

tmux() {
    local status
    log "call tmux: $@"

    if "${_real_tmux}" "$@"; then
        log "tmux call succeeded"
        return 0
    else
        status=$?
        log "tmux call failed with $status"
        return "$status"
    fi
}

_init_log() {
    mkdir -p "$(dirname ${c_log_file_path})" 2>/dev/null
    touch "${c_log_file_path}" 2>/dev/null || true
}

_cleanup() {
    log "cleanup called"
    [[ -z "$(jobs -p)" ]] || kill $(jobs -p)
}

trap _cleanup SIGINT SIGTERM EXIT

source "${lib_root}/config.sh"

init_config
read_user_config_if_present
preprocess_config
[[ -n $DEBUG ]] && dump_config

_init_log
