#!/usr/bin/env bash
set -euo pipefail

launchd_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
export CONFIG_PATH="${CONFIG_PATH:-${HOME}/.config/multienv/server.launchd.yaml}"

exec bash "$launchd_dir/../common/render.sh" "$@"
