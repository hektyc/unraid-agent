#!/bin/bash
PLUGIN_NAME="unraid-mcp"
INSTALL_DIR="/usr/local/emhttp/plugins/${PLUGIN_NAME}"
CONFIG_DIR="/boot/config/plugins/${PLUGIN_NAME}"

mkdir -p "$INSTALL_DIR"
mkdir -p "$CONFIG_DIR"

cp -r "$ROOT/"* "$INSTALL_DIR/" 2>/dev/null || true
mkdir -p "$INSTALL_DIR/bin"

chmod +x "$INSTALL_DIR/scripts/"*.sh

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
echo "Install complete."
