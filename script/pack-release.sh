#!/usr/bin/env bash
# 打出发版归档：前端嵌入 Go 单二进制（Linux / Windows amd64，CGO 关闭）。
# 用法（仓库根）：
#   bash script/pack-release.sh
# 环境变量：
#   SKIP_INSTALL=1    跳过 pnpm install（CI 已装好依赖）
#   SKIP_FRONTEND=1   跳过 vue-tsc / vite build（已有 embed dist）
#   SKIP_TEST=1       跳过 go test
#   NIUMA_SITE_VERSION  归档版本号，默认读 package.json
# CHANGELOG.md 必须含 ## [version] 段，否则失败。
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

VERSION="${NIUMA_SITE_VERSION:-}"
if [[ -z "$VERSION" ]]; then
  VERSION="$(node -p "require('./package.json').version")"
fi
VERSION="${VERSION#v}"

OUT="$ROOT/output"
LINUX_NAME="niuma-site-${VERSION}-linux-amd64"
WIN_NAME="niuma-site-${VERSION}-windows-amd64"
STAGE_LINUX="$OUT/stage-linux/${LINUX_NAME}"
STAGE_WIN="$OUT/stage-windows/${WIN_NAME}"
LINUX_TAR="$OUT/${LINUX_NAME}.tar.gz"
WIN_ZIP="$OUT/${WIN_NAME}.zip"

rm -rf "$OUT"
mkdir -p "$STAGE_LINUX/run" "$STAGE_LINUX/script" "$STAGE_LINUX/config"
mkdir -p "$STAGE_WIN/run" "$STAGE_WIN/script" "$STAGE_WIN/config"

echo "==> changelog ${VERSION}"
node "$ROOT/.github/scripts/changelog-notes.mjs" "$VERSION" "$OUT/release-notes.md" --require

if [[ "${SKIP_INSTALL:-}" != "1" ]]; then
  echo "==> pnpm install"
  if [[ -d "$ROOT/../niuma-ui" ]]; then
    pnpm --dir "$ROOT/../niuma-ui" install --frozen-lockfile
  fi
  pnpm install --frozen-lockfile
fi

if [[ "${SKIP_FRONTEND:-}" != "1" ]]; then
  echo "==> frontend build (embed -> server/internal/web/dist)"
  pnpm exec vue-tsc --noEmit
  pnpm exec vite build
fi

if [[ "${SKIP_TEST:-}" != "1" ]]; then
  echo "==> go test"
  (cd "$ROOT/server" && go test ./...)
fi

LDFLAGS="-s -w -X main.version=${VERSION}"
echo "==> go build linux-amd64 / windows-amd64 (version=${VERSION})"
(
  cd "$ROOT/server"
  CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="$LDFLAGS" \
    -o "$STAGE_LINUX/run/niuma-site" ./cmd/server
  CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -trimpath -ldflags="$LDFLAGS" \
    -o "$STAGE_WIN/run/niuma-site.exe" ./cmd/server
)

copy_runtime() {
  local dest="$1"
  printf '%s\n' "$VERSION" > "$dest/VERSION"
  cp "$ROOT/README.md" "$dest/README.md"
  cp "$ROOT/CHANGELOG.md" "$dest/CHANGELOG.md"
  cp "$ROOT/config/app.yaml.example" "$dest/config/app.yaml.example"
  cp "$ROOT/script/README.md" "$dest/script/README.md"
}

copy_runtime "$STAGE_LINUX"
cp "$ROOT/script/start.sh" "$STAGE_LINUX/script/start.sh"
cp "$ROOT/script/stop.sh" "$STAGE_LINUX/script/stop.sh"
cp "$ROOT/script/restart.sh" "$STAGE_LINUX/script/restart.sh"
chmod +x "$STAGE_LINUX/run/niuma-site" "$STAGE_LINUX/script/"*.sh

copy_runtime "$STAGE_WIN"
cp "$ROOT/script/start.ps1" "$STAGE_WIN/script/start.ps1"
cp "$ROOT/script/stop.ps1" "$STAGE_WIN/script/stop.ps1"
cp "$ROOT/script/restart.ps1" "$STAGE_WIN/script/restart.ps1"

zip_windows() {
  local src="$OUT/stage-windows"
  if command -v zip >/dev/null 2>&1; then
    (cd "$src" && zip -qr "$WIN_ZIP" "$WIN_NAME")
    return
  fi
  if [[ -x /c/Windows/System32/tar.exe ]]; then
    /c/Windows/System32/tar.exe -a -c -f "$WIN_ZIP" -C "$src" "$WIN_NAME"
    return
  fi
  echo "Need zip or Windows tar.exe to create the Windows archive." >&2
  exit 1
}

echo "==> archive"
tar -C "$OUT/stage-linux" -czf "$LINUX_TAR" "$LINUX_NAME"
zip_windows
node "$ROOT/.github/scripts/write-sha256.mjs" "$OUT"

echo "OK: $LINUX_TAR"
echo "OK: $WIN_ZIP"
echo "OK: $OUT/SHA256SUMS.txt"
echo "OK: $OUT/release-notes.md"
