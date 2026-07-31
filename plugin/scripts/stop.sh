#!/bin/bash
PLUGIN_NAME="unraid-mcp"
CONFIG_DIR="/boot/config/plugins/${PLUGIN_NAME}"
PIDFILE="$CONFIG_DIR/server.pid"

if [ -f "$PIDFILE" ]; then
    PID=$(cat "$PIDFILE")
    if kill -0 "$PID" 2>/dev/null; then
        kill "$PID" 2>/dev/null
        sleep 2
        if kill -0 "$PID" 2>/dev/null; then
            kill -9 "$PID" 2>/dev/null
        fi
        echo "Server stopped."
    else
        echo "Server not running."
    fi
    rm -f "$PIDFILE"
else
    echo "PID file not found. Server may not be running."
fi
