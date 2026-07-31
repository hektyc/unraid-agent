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

set -a
source "$CONFIG_DIR/config.cfg"
set +a

mkdir -p "$(dirname "$LOGFILE")"

nohup "$INSTALL_DIR/bin/unraid-mcp" >> "$LOGFILE" 2>&1 &
PID=$!
echo $PID > "$PIDFILE"
echo "Server started (PID $PID)"
