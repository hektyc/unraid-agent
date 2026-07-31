#!/bin/bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

VERSION="$(cat VERSION)"
PLUGIN_NAME="unraid-mcp"
PLG_FILE="${PLUGIN_NAME}.plg"
PACKAGE_FILE="${PLUGIN_NAME}-${VERSION}.tar.gz"
BUILD_DIR="/tmp/${PLUGIN_NAME}-build"
INSTALL_DIR="${BUILD_DIR}/usr/local/emhttp/plugins/${PLUGIN_NAME}"
CONFIG_DIR="${BUILD_DIR}/boot/config/plugins/${PLUGIN_NAME}"

echo "==> Building Go binary..."
cd "$SCRIPT_DIR/go"
go build -o bin/unraid-mcp cmd/unraid-mcp/main.go
cd "$SCRIPT_DIR"

echo "==> Creating build directory..."
rm -rf "$BUILD_DIR"
mkdir -p "$INSTALL_DIR" "$CONFIG_DIR"

echo "==> Assembling package files..."
mkdir -p "$INSTALL_DIR/bin"
mkdir -p "$INSTALL_DIR/scripts"
mkdir -p "$INSTALL_DIR/webgui"
mkdir -p "$INSTALL_DIR/memory"

cp go/bin/unraid-mcp "$INSTALL_DIR/bin/unraid-mcp"
cp -r plugin/scripts/*.sh "$INSTALL_DIR/scripts/"
cp plugin/webgui/Settings.page "$INSTALL_DIR/webgui/Settings.page"
cp -r memory/* "$INSTALL_DIR/memory/"
cp plugin/unraid-mcp.plg "$INSTALL_DIR/unraid-mcp.plg"
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
ALLOW_ARRAY_ADD_DISK="false"
ALLOW_ARRAY_REMOVE_DISK="false"
ALLOW_ARRAY_CLEAR_STATS="false"
ALLOW_CONTAINER_STOP="false"
ALLOW_CONTAINER_REMOVE="false"
ALLOW_CONTAINER_RESTART="false"
ALLOW_VM_FORCE_STOP="false"
ALLOW_VM_RESET="false"
ALLOW_VM_STOP="false"
ALLOW_PLUGIN_INSTALL="false"
ALLOW_PLUGIN_REMOVE="false"
ALLOW_SETTING_UPDATES="false"
ALLOW_SSH_UPDATE="false"
ALLOW_TIME_UPDATE="false"
ALLOW_NOTIFICATION_DELETE="false"
ALLOW_API_KEY_CREATE="false"
ALLOW_API_KEY_DELETE="false"
ALLOW_FLASH_BACKUP="false"
ALLOW_RCLONE_OPERATIONS="false"
ALLOW_CONNECT_ACTIONS="false"
ALLOW_ONBOARDING_ACTIONS="false"
ALLOW_DOCKER_ACTIONS="false"
ALLOW_VM_ACTIONS="false"
ALLOW_ARRAY_ACTIONS="false"
ALLOW_DESTRUCTIVE="false"
UNRAID_MCP_LOG_LEVEL="info"
UNRAID_VERIFY_SSL="true"
UNRAID_ALLOW_INSECURE_TLS="false"
UNRAID_MCP_BEARER_TOKEN=""
UNRAID_MCP_DISABLE_HTTP_AUTH="false"
ENVEOF
chmod 600 "$CONFIG_DIR/config.cfg"

echo "==> Creating package tarball..."
cd "$BUILD_DIR"
tar -czf "$SCRIPT_DIR/$PACKAGE_FILE" .
cd "$SCRIPT_DIR"

echo "==> Generating .plg manifest..."
cat > "$PLG_FILE" << XMLDOC
<?xml version="1.0" encoding="utf-8"?>
<plugin name="Unraid MCP Server" version="${VERSION}" url="https://github.com/hektyc/unraid-mcp-server">
  <description>MCP server for Unraid with safety controls and AI agent integration</description>
  <author>hektyc</author>
  <arch>amd64</arch>

  <minUnraidVersion>6.12.0</minUnraidVersion>
  <maxUnraidVersion>7.99.99</maxUnraidVersion>

  <package file="${PACKAGE_FILE}" url="https://github.com/hektyc/unraid-mcp-server/releases/download/v${VERSION}"/>

  <configdir>/boot/config/plugins/unraid-mcp</configdir>

  <scripts>
    <install type="shell">scripts/install.sh</install>
    <upgrade type="shell">scripts/upgrade.sh</upgrade>
    <remove type="shell">scripts/remove.sh</remove>
    <start type="shell">scripts/start.sh</start>
    <stop type="shell">scripts/stop.sh</stop>
    <status type="shell">scripts/status.sh</status>
  </scripts>

  <webgui>
    <page Settings.page />
  </webgui>

  <notes>
    <![CDATA[
    <div style="padding: 10px; color: #333;">
      <h3>Unraid MCP Server v${VERSION}</h3>
      <p><strong>Safety:</strong> Default read-only. Enable actions via Settings > Unraid MCP Server.</p>
      <p><strong>Transport:</strong> stdio (default) or streamable-http</p>
      <p><strong>Memory:</strong> AGENT-SKILLS.md and AGENT-MEMORY.md are stored in config dir</p>
      <p style="color: #c00;"><strong>Warning:</strong> Disabling read-only allows destructive operations. Enable only what you need.</p>
    </div>
    ]]>
  </notes>
</plugin>
XMLDOC

rm -rf "$BUILD_DIR"

echo "==> Done: $PLG_FILE ($(du -h "$PLG_FILE" | cut -f1)) + $PACKAGE_FILE ($(du -h "$PACKAGE_FILE" | cut -f1))"
