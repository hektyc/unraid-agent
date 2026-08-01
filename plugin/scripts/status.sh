#!/bin/bash
PLUGIN_NAME="unraid-agent"
CONFIG_DIR="/boot/config/plugins/${PLUGIN_NAME}"
PIDFILE="$CONFIG_DIR/server.pid"

if [ -f "$PIDFILE" ]; then
    PID=$(cat "$PIDFILE")
    if kill -0 "$PID" 2>/dev/null; then
        echo "Running (PID $PID)"
        exit 0
    fi
fi
echo "Stopped"
exit 1
