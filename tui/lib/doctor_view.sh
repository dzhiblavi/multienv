#!/usr/bin/env bash

doctor_report_format() {
    local line color
    while IFS= read -r line; do
        color=0
        case "$line" in
        '  OK '*) color=32 ;;
        '  WARN '*) color=33 ;;
        '  FAIL '*) color='1;31' ;;
        '  SKIPPED '*) color=90 ;;
        *' ['*']') color='1;36' ;;
        esac

        printf '\033[%sm%s\033[0m\n' "$color" "$line"
    done
}

doctor_report() {
    local status
    "$BASH" "${script_dir}/multienv-tui" --doctor 2>&1 | doctor_report_format
    status="${PIPESTATUS[0]}"

    if ((status == 0)); then
        printf '\n\033[1;32mDiagnostics complete: no failed checks\033[0m\n'
    else
        printf '\n\033[1;31mDiagnostics found problems (exit %s)\033[0m\n' "$status"
    fi
}

run_doctor_picker() {
    local query="${1:-}" args=()
    popup_enabled && args+=(--tmux=center,90%,90%)

    ui_fzf --ansi --layout=reverse --no-sort --wrap --no-multi --no-select-1 --no-expect \
        --preview= --preview-window=hidden \
        --expect=ctrl-b,ctrl-/,f1,esc --print-query --query="$query" --prompt='Doctor ❯ ' \
        --header="$(ui_header doctor)" --with-shell="$child_shell" \
        --bind='enter:ignore,double-click:ignore' \
        --bind="start,ctrl-r:reload($self_command --doctor-report)" \
        "${args[@]}" </dev/null
}
