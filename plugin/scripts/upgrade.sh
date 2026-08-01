#!/bin/bash
PLUGIN_NAME="unraid-agent"
INSTALL_DIR="/usr/local/emhttp/plugins/${PLUGIN_NAME}"
CONFIG_DIR="/boot/config/plugins/${PLUGIN_NAME}"

# Backup config
if [ -d "$CONFIG_DIR" ]; then
    cp -f "$CONFIG_DIR/config.cfg" "/tmp/${PLUGIN_NAME}-config-backup.cfg" 2>/dev/null || true
fi

# Reinstall
bash "$INSTALL_DIR/scripts/install.sh"

# Restore config if backup exists
if [ -f "/tmp/${PLUGIN_NAME}-config-backup.cfg" ]; then
    cp -f "/tmp/${PLUGIN_NAME}-config-backup.cfg" "$CONFIG_DIR/config.cfg"
    rm -f "/tmp/unraid-agent-config-backup.cfg"
fi

echo "Upgrade complete."
