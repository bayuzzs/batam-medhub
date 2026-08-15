#!/usr/bin/env bash
# Top-level entrypoint for running the consolidated provider smoke verification suite.
set -Eeuo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
exec "$SCRIPT_DIR/scripts/smoke.sh" "$@"
