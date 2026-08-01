#!/bin/bash
PLUGIN_NAME="unraid-agent"
CONFIG_DIR="/boot/config/plugins/${PLUGIN_NAME}"
PIDFILE="$CONFIG_DIR/server.pid"

if [ -f "$PIDFILE" ]; then
    PID=$(cat "$PIDFILE")
    if kill -0 "$PID" 2>/dev/null; then
        kill "$PID" 2>/dev/null
        # Wait up to 5 seconds for graceful shutdown
        for i in 1 2 3 4 5; do
            if ! kill -0 "$PID" 2>/dev/null; then
                break
            fi
            sleep 1
        done
        if kill -0 "$PID" 2>/dev/null; then
            kill -9 "$PID" 2>/dev/null
            sleep 1
        fi
        echo "Server stopped."
    else
        echo "Server not running."
    fi
    rm -f "$PIDFILE"
else
    echo "PID file not found. Server may not be running."
fi
