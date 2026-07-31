#!/bin/bash
PLUGIN_NAME="unraid-mcp"
INSTALL_DIR="/usr/local/emhttp/plugins/${PLUGIN_NAME}"
CONFIG_DIR="/boot/config/plugins/${PLUGIN_NAME}"
PIDFILE="$CONFIG_DIR/server.pid"
LOGFILE="/var/log/unraid-mcp.log"

mkdir -p "$CONFIG_DIR"

if [ -f "$PIDFILE" ]; then
    OLDPID=$(cat "$PIDFILE")
    if kill -0 "$OLDPID" 2>/dev/null; then
        echo "Server already running (PID $OLDPID)"
        exit 0
    else
        rm -f "$PIDFILE"
    fi
fi

source "$CONFIG_DIR/config.cfg"

export UNRAID_API_URL UNRAID_API_KEY TRANSPORT UNRAID_MCP_HOST UNRAID_MCP_PORT
export READ_ONLY ALLOW_ARRAY_STOP ALLOW_ARRAY_START ALLOW_DOCKER_ACTIONS ALLOW_VM_ACTIONS
export ALLOW_DESTRUCTIVE UNRAID_MCP_LOG_LEVEL UNRAID_VERIFY_SSL UNRAID_ALLOW_INSECURE_TLS
export UNRAID_MCP_BEARER_TOKEN UNRAID_MCP_DISABLE_HTTP_AUTH

mkdir -p "$(dirname "$LOGFILE")"

nohup "$INSTALL_DIR/bin/unraid-mcp" >> "$LOGFILE" 2>&1 &
PID=$!
echo $PID > "$PIDFILE"
echo "Server started (PID $PID)"
