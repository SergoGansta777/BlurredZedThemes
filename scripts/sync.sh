#!/usr/bin/env bash
set -euo pipefail

if [[ $# -lt 1 ]]; then
  echo "usage: scripts/sync.sh <theme-name-without-extension>"
  echo "example: scripts/sync.sh evergarden-winter-hybrid"
  exit 1
fi

theme="$1"
palette="palettes/${theme}.json"
reference="themes/${theme}.json"

if [[ ! -f "$palette" ]]; then
  echo "missing palette: $palette"
  exit 1
fi
if [[ ! -f "$reference" ]]; then
  echo "missing reference theme: $reference"
  exit 1
fi

go run ./scripts/generate --palette "$palette" --compare "$reference" --write-alpha --write-overrides --prune-alpha-overrides --prune-overrides
go run ./scripts/generate --palette "$palette" --compare "$reference"
