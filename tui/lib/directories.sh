#!/usr/bin/env bash

directory_options() {
    local response="${1:?need response}" current="${2:?need current directory}"
    jq -j --arg current "$current" '
        reduce ([$current] + .directories)[] as $directory ([];
            if index($directory) == null then . + [$directory] else . end) |
        .[] | ., "\u0000"
    ' <<<"$response"
}

choose_directory() {
    local host="${1:?need hostname}" current="${2:?need current directory}" response header
    header="Host: $host · Enter: select · Esc: cancel"

    if ! response="$(jq -n --arg host "$host" '{hostname: $host}' |
        multienv-agent filesystem directory-candidates 2>>"$log_file")" ||
        ! jq -e '.directories | type == "array" and all(.[]; type == "string")' \
            <<<"$response" >/dev/null 2>>"$log_file"; then
        show_error 'Directory lookup failed'
        response='{"directories":[]}'
        header+=$'\nDirectory lookup failed; showing current directory only'
    fi

    directory_options "$response" "$current" |
        fzf --no-popup --no-height --no-multi --no-expect --no-select-1 --no-print-query \
            --read0 --print0 --layout=reverse --scheme=path --tiebreak=index \
            --prompt='New window directory ❯ ' --header="$header" --preview-window=hidden
}
