#!/usr/bin/env bash
# 构建官网：前端 → 嵌入目录 → Go 单二进制（Linux）
# 用法：bash script/build.sh
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

mkdir -p "$ROOT/run"

echo "==> pnpm install"
pnpm install

echo "==> frontend build (embed -> server/internal/web/dist)"
pnpm exec vue-tsc --noEmit
pnpm exec vite build

BIN="$ROOT/run/niuma-site"
echo "==> go build -> $BIN"
(
  cd "$ROOT/server"
  go build -o "$BIN" ./cmd/server
)

echo "OK: $BIN"
echo "Start: bash script/start.sh"
