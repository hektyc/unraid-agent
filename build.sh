#!/bin/bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

VERSION="$(cat VERSION)"
PLUGIN_NAME="unraid-mcp"
PLG_FILE="${PLUGIN_NAME}.plg"
BUILD_DIR="/tmp/${PLUGIN_NAME}-build"
INSTALL_DIR="${BUILD_DIR}/usr/local/emhttp/plugins/${PLUGIN_NAME}"
CONFIG_DIR="${BUILD_DIR}/boot/config/plugins/${PLUGIN_NAME}"

echo "==> Building Go binary..."
cd "$SCRIPT_DIR/go"
go build -o bin/unraid-mcp cmd/unraid-mcp/main.go
cd "$SCRIPT_DIR"

echo "==> Creating build directory..."
rm -rf "$BUILD_DIR"
mkdir -p "$INSTALL_DIR"
mkdir -p "$CONFIG_DIR"

echo "==> Copying plugin files..."
cp -r plugin/* "$INSTALL_DIR/"
cp -r go "$INSTALL_DIR/"
cp -r memory "$INSTALL_DIR/"
mkdir -p "$INSTALL_DIR/bin"
cp go/bin/unraid-mcp "$INSTALL_DIR/bin/unraid-mcp"
chmod +x "$INSTALL_DIR/scripts/"*.sh

echo "==> Creating config.cfg..."
cat > "$CONFIG_DIR/config.cfg" << 'ENVEOF'
UNRAID_API_URL="http://localhost/graphql"
UNRAID_API_KEY=""
TRANSPORT="stdio"
UNRAID_MCP_HOST="127.0.0.1"
UNRAID_MCP_PORT="6970"
READ_ONLY="true"
ALLOW_ARRAY_STOP="false"
ALLOW_ARRAY_START="false"
ALLOW_CONTAINER_STOP="false"
ALLOW_CONTAINER_REMOVE="false"
ALLOW_VM_STOP="false"
ALLOW_PLUGIN_INSTALL="false"
ALLOW_DESTRUCTIVE="false"
UNRAID_MCP_LOG_LEVEL="info"
UNRAID_VERIFY_SSL="true"
UNRAID_ALLOW_INSECURE_TLS="false"
UNRAID_MCP_BEARER_TOKEN=""
UNRAID_MCP_DISABLE_HTTP_AUTH="false"
ENVEOF
chmod 600 "$CONFIG_DIR/config.cfg"

echo "==> Creating .plg package..."
cd "$BUILD_DIR"
tar -cf - . | xz -z > "$SCRIPT_DIR/$PLG_FILE"
cd "$SCRIPT_DIR"
rm -rf "$BUILD_DIR"

echo "==> Done: $PLG_FILE ($(du -h "$PLG_FILE" | cut -f1))"
