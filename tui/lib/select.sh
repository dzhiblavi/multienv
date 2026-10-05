#!/usr/bin/env bash

select_session() {
    if ! is_in_tmux; then
        die "should run within tmux session"
    fi

    local hostname="${1:?need hostname}"
    local session="${2:?need session}"

    if is_on_hostname "${hostname}"; then
        tmux switch-client -t "${session}"
    else
        send "${hostname}" "${session}"
        tmux detach
    fi
}

