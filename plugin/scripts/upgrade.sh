#!/bin/bash
PLUGIN_NAME="unraid-mcp"
INSTALL_DIR="/usr/local/emhttp/plugins/${PLUGIN_NAME}"
CONFIG_DIR="/boot/config/plugins/${PLUGIN_NAME}"

# Backup config
if [ -d "$CONFIG_DIR" ]; then
    cp -f "$CONFIG_DIR/${PLUGIN_NAME}.cfg" "/tmp/${PLUGIN_NAME}-config-backup.cfg" 2>/dev/null || true
fi

# Reinstall
bash "$INSTALL_DIR/scripts/install.sh"

# Restore config if backup exists
if [ -f "/tmp/${PLUGIN_NAME}-config-backup.cfg" ]; then
    cp -f "/tmp/${PLUGIN_NAME}-config-backup.cfg" "$CONFIG_DIR/${PLUGIN_NAME}.cfg"
    rm -f "/tmp/unraid-mcp-config-backup.cfg"
fi

echo "Upgrade complete."
