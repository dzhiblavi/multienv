#!/usr/bin/env bash

c_user_config_path="${HOME}/.config/multienv/tui.sh"

init_config() {
    # All these values can be overriden by the user

    c_ctrl_port=7777
    c_log_file_path="/tmp/multienv-tui.log"
    c_tui_ignore_file_path="${HOME}/.config/multienv/.multienv-tui-ignore"

    declare -gA c_host_alias_to_hostname=()
}

read_user_config_if_present() {
    [[ -f "${c_user_config_path}" ]] && source "${c_user_config_path}"
}

preprocess_config() {
    declare -gA c_hostname_to_host_alias=()
    declare -ga c_hostnames=()

    for host_alias in "${!c_host_alias_to_hostname[@]}"; do
        hostname="${c_host_alias_to_hostname[${host_alias}]}"
        c_hostname_to_host_alias["${hostname}"]="${host_alias}"
        c_hostnames+=("${hostname}")
    done
}

dump_config() {
    echo "configuration dump" >&2

    for var in ${!c_*}; do
        declare -p "$var" >&2
    done
}
