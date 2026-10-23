#!/usr/bin/env bash

doctor_result() {
    printf '  %-7s %-12s %s: %s\n' "$1" tui "$2" "${3//$'\n'/ }"

    [[ "$1" != FAIL ]]
}

doctor_check_command() {
    local name="${1:?need command}" path output
    shift
    path="$(command -v "$name")" || {
        doctor_result FAIL "binary.$name" 'not found in PATH'
        return 1
    }

    if (($# == 0)); then
        doctor_result OK "binary.$name" "$path"
        return
    fi

    output="$("$name" "$@" 2>&1)" || {
        doctor_result FAIL "binary.$name" "$path: $output"
        return 1
    }

    doctor_result OK "binary.$name" "$path: ${output%%$'\n'*}"
}

doctor_check_bash() {
    if ((BASH_VERSINFO[0] < 4)); then
        doctor_result FAIL binary.bash "$BASH: version $BASH_VERSION; Bash 4+ must be first in PATH"
        return 1
    fi

    doctor_result OK binary.bash "$BASH: version $BASH_VERSION"
}

doctor_check_fzf() {
    local help output status=0
    doctor_check_command fzf --version || return

    help="$(fzf --help 2>&1)" || {
        doctor_result FAIL fzf.capabilities "$help"
        return 1
    }

    if [[ "$help" != *--tmux* || "$help" != *--id-nth* ]]; then
        doctor_result FAIL fzf.capabilities 'requires --tmux and --id-nth'
        return 1
    fi

    output="$(fzf --filter='' --bind='every(3):ignore' </dev/null 2>&1)" || status=$?
    if ((status > 1)); then
        doctor_result FAIL fzf.capabilities "$output"
        return 1
    fi

    doctor_result OK fzf.capabilities '--tmux, --id-nth and every(...) bindings available'
}

doctor_check_popup() {
    local version major minor
    version="$(tmux -V 2>&1)" || {
        doctor_result FAIL tmux.popup "$version"
        return 1
    }

    if [[ "$version" =~ ([0-9]+)\.([0-9]+) ]]; then
        major="${BASH_REMATCH[1]}"
        minor="${BASH_REMATCH[2]}"
        if ((major > 3 || (major == 3 && minor >= 3))); then
            doctor_result OK tmux.popup "$version"
            return
        fi
    fi

    doctor_result WARN tmux.popup "$version; popup mode requires tmux 3.3+; use --no-popup otherwise"
}

doctor_check_netcat() {
    local path help
    path="$(command -v nc)" || {
        doctor_result FAIL binary.nc 'not found in PATH'
        return 1
    }

    help="$(nc -h 2>&1)"
    if [[ "$help" != *-k* ]]; then
        doctor_result FAIL binary.nc "$path: requires OpenBSD-compatible nc with -k support"
        return 1
    fi

    doctor_result OK binary.nc "$path: -k available"
}

doctor_check_unix_tools() {
    local name missing=''
    for name in basename cat cksum cut date dirname hostname mkdir mktemp readlink rm sed sleep sort touch which; do
        command -v "$name" >/dev/null || missing+=" $name"
    done

    if [[ -n "$missing" ]]; then
        doctor_result FAIL unix.tools "missing:$missing"
        return 1
    fi

    doctor_result OK unix.tools 'standard Unix tools available'
}

doctor_check_config() {
    local path="${HOME}/.config/multienv/tui.sh" output
    if [[ ! -f "$path" ]]; then
        doctor_result OK config 'no tui.sh; using defaults'
        return
    fi

    output="$("$BASH" -n "$path" 2>&1)" || {
        doctor_result FAIL config "$output"
        return 1
    }

    doctor_result OK config "$path: shell syntax valid"
}

doctor_run_agent() {
    if ! command -v multienv-agent >/dev/null; then
        doctor_result FAIL binary.multienv-agent 'not found in PATH; agent diagnostics unavailable'
        return 1
    fi

    multienv-agent doctor global
}

tui_doctor() {
    case "${1:-}" in
    --help | -h)
        printf 'Usage: multienv-tui --doctor\n\n'
        printf 'Check local TUI dependencies and run multienv-agent doctor global.\n'
        printf 'Runs outside tmux; exits nonzero if either part fails.\n'
        return
        ;;
    '') ;;
    *)
        printf 'Unknown doctor option: %s\n' "$1" >&2
        return 2
        ;;
    esac

    local status=0
    printf '%s [tui]\n' "${HOSTNAME:-unknown}"
    doctor_check_bash || status=1
    doctor_check_command jq --version || status=1
    doctor_check_fzf || status=1
    doctor_check_command realpath . || status=1
    doctor_check_popup || status=1
    doctor_check_netcat || status=1
    doctor_check_unix_tools || status=1
    doctor_check_config || status=1

    printf '\n'
    doctor_run_agent || status=1

    return "$status"
}
