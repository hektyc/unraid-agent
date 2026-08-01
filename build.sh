#!/bin/bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

VERSION="$(cat VERSION)"
PLUGIN_NAME="unraid-agent"
PACKAGE_FILE="unraid-agent.tar"
BUILD_DIR="/tmp/${PLUGIN_NAME}-build"
INSTALL_DIR="${BUILD_DIR}/usr/local/emhttp/plugins/${PLUGIN_NAME}"
CONFIG_DIR="${BUILD_DIR}/boot/config/plugins/${PLUGIN_NAME}"

echo "==> Building Go binary..."
cd "$SCRIPT_DIR/go"
go build -o bin/unraid-agent cmd/unraid-agent/main.go
cd "$SCRIPT_DIR"

echo "==> Creating build directory..."
rm -rf "$BUILD_DIR"
mkdir -p "$INSTALL_DIR/bin"
mkdir -p "$INSTALL_DIR/scripts"
mkdir -p "$INSTALL_DIR/images"
mkdir -p "$INSTALL_DIR/memory"
mkdir -p "$CONFIG_DIR"

echo "==> Assembling package files..."
cp go/bin/unraid-agent "$INSTALL_DIR/bin/unraid-agent"
cp plugin/scripts/install.sh "$INSTALL_DIR/scripts/install.sh"
cp plugin/scripts/start.sh "$INSTALL_DIR/scripts/start.sh"
cp plugin/scripts/stop.sh "$INSTALL_DIR/scripts/stop.sh"
cp plugin/scripts/status.sh "$INSTALL_DIR/scripts/status.sh"
cp plugin/scripts/upgrade.sh "$INSTALL_DIR/scripts/upgrade.sh"
cp plugin/scripts/apply.sh "$INSTALL_DIR/scripts/apply.sh"
cp plugin/scripts/remove.sh "$INSTALL_DIR/scripts/remove.sh"

# Page file goes in plugin ROOT (not webgui/ subdirectory)
cp plugin/unraid-agent.page "$INSTALL_DIR/unraid-agent.page"

# default.cfg goes in plugin root as a fallback default
# NOTE: config.cfg is intentionally NOT packaged — it is created on first
# install by the .plg copy-if-missing script and must survive updates.
cp plugin/webgui/config.cfg "$INSTALL_DIR/default.cfg"

# Icon images
cp plugin/images/unraid-agent-48.png "$INSTALL_DIR/images/unraid-agent-48.png"
cp plugin/images/unraid-agent-128.png "$INSTALL_DIR/images/unraid-agent-128.png"
cp plugin/images/unraid-agent.png "$INSTALL_DIR/images/unraid-agent.png"

cp memory/AGENT-SKILLS.md "$INSTALL_DIR/memory/AGENT-SKILLS.md"
cp memory/AGENT-MEMORY.md "$INSTALL_DIR/memory/AGENT-MEMORY.md"
cp plugin/unraid-agent.plg "$INSTALL_DIR/unraid-agent.plg"
cp plugin/unraid-agent.plg "$CONFIG_DIR/unraid-agent.plg"
chmod +x "$INSTALL_DIR/scripts/"*.sh
chmod 644 "$CONFIG_DIR/unraid-agent.plg" "$CONFIG_DIR/config.cfg"
chmod 644 "$INSTALL_DIR/images/"*.png

echo "==> Creating package tarball..."
cd "$BUILD_DIR"
tar --owner=0 --group=0 --no-same-owner -cf "$SCRIPT_DIR/$PACKAGE_FILE" .
cd "$SCRIPT_DIR"

rm -rf "$BUILD_DIR"

echo "==> Done: $PACKAGE_FILE ($(du -h "$PACKAGE_FILE" | cut -f1))"
