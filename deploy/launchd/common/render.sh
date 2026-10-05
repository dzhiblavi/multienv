#!/usr/bin/env bash
set -euo pipefail

launchd_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
config_path="${CONFIG_PATH:?set CONFIG_PATH to the launchd YAML config}"

for command_name in yq jq; do
    if ! command -v "$command_name" >/dev/null 2>&1; then
        echo "$command_name is required" >&2
        exit 1
    fi
done

# launchd does not expand shell variables or ~ in paths.
config_json="$(
    yq -o=json '.' "$config_path" |
        jq -e --arg home "$HOME" '
            def expand_home:
                sub("^~(?=/|$)"; $home) |
                gsub("\\$\\{HOME\\}|\\$HOME"; $home);
            if (.label | type == "string" and test("^[A-Za-z0-9][A-Za-z0-9._-]*$"))
                and ([.program_path, .config_path, .stdout_path, .stderr_path] |
                    all(.[]; type == "string" and length > 0))
                and (.environment == null or (.environment | type == "object"))
                and ((.environment // {}) | all(.[]; type == "string"))
            then
                .program_path |= expand_home |
                .config_path |= expand_home |
                .stdout_path |= expand_home |
                .stderr_path |= expand_home |
                .environment = ((.environment // {}) |
                    with_entries(.value |= expand_home)) |
                if ([.program_path, .config_path, .stdout_path, .stderr_path] |
                    all(.[]; startswith("/")))
                then . else error("launchd paths must be absolute") end
            else error("invalid launchd config: label and all paths are required; " +
                "environment must contain strings")
            end
        '
)"

if [[ "${1:-}" == "--get" ]]; then
    [[ $# == 2 ]] || { echo "usage: $0 --get FIELD" >&2; exit 1; }
    case "$2" in
        label|program_path|config_path|stdout_path|stderr_path)
            jq -er --arg field "$2" '.[$field]' <<<"$config_json"
            ;;
        *) echo "unsupported field: $2" >&2; exit 1 ;;
    esac
elif [[ $# == 0 ]]; then
    printf '%s\n' "$config_json" |
        jinja2 --strict --format=json "$launchd_dir/plist.j2" -
else
    echo "usage: $0 [--get FIELD]" >&2
    exit 1
fi
