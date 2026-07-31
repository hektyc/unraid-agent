#!/bin/bash
PLUGIN_NAME="unraid-mcp"
INSTALL_DIR="/usr/local/emhttp/plugins/${PLUGIN_NAME}"
CONFIG_DIR="/boot/config/plugins/${PLUGIN_NAME}"

$INSTALL_DIR/scripts/stop.sh >/dev/null 2>&1 || true
rm -rf "$INSTALL_DIR"
rm -rf "$CONFIG_DIR"
echo "Removal complete."
