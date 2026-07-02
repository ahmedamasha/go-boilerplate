#!/usr/bin/env bash
# Prefer nvm Node over Homebrew node@22 (simdjson dylib crash on macOS after brew upgrade).
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
NVM_BASE="${NVM_DIR:-$HOME/.nvm}/versions/node"

use_nvm_node() {
  local ver="$1"
  local dir
  # .nvmrc may be "20" while nvm installs "v20.17.0" — match v20*
  for dir in "$NVM_BASE/v${ver}"*/bin "$NVM_BASE/v${ver}/bin"; do
    if [ -x "$dir/node" ]; then
      export PATH="$dir:$PATH"
      return 0
    fi
  done
  return 1
}

if [ -f "$ROOT/.nvmrc" ]; then
  NODE_VER="$(tr -d '[:space:]' < "$ROOT/.nvmrc")"
  use_nvm_node "$NODE_VER" || true
fi

if ! command -v node >/dev/null 2>&1 || ! node --version >/dev/null 2>&1; then
  echo "error: Node is unavailable (Homebrew node@22/simdjson crash?)." >&2
  echo "  Fix: nvm install ${NODE_VER:-20} && make web-dev" >&2
  exit 1
fi

cd "$ROOT/web"
exec npm "$@"
