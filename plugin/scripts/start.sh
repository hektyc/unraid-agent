#!/bin/bash
PLUGIN_NAME="unraid-agent"
INSTALL_DIR="/usr/local/emhttp/plugins/${PLUGIN_NAME}"
CONFIG_DIR="/boot/config/plugins/${PLUGIN_NAME}"
PIDFILE="$CONFIG_DIR/server.pid"
LOGFILE="/var/log/unraid-agent.log"

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

if [ -z "${UNRAID_API_URL:-}" ]; then
    INI_FILE="/var/local/emhttp/var.ini"
    if [ -f "$INI_FILE" ]; then
        LOCAL_IP=$(grep '^IPADDR=' "$INI_FILE" | head -1 | cut -d'=' -f2-)
        USE_SSL=$(grep '^USE_SSL=' "$INI_FILE" | head -1 | cut -d'=' -f2-)
        PORT=$(grep '^PORT=' "$INI_FILE" | head -1 | cut -d'=' -f2-)
        PORTSSL=$(grep '^PORTSSL=' "$INI_FILE" | head -1 | cut -d'=' -f2-)

        if [ "$USE_SSL" = "yes" ]; then
            PROTO="https"
            LOCAL_PORT="${PORTSSL:-443}"
        else
            PROTO="http"
            LOCAL_PORT="${PORT:-80}"
        fi

        if [ -n "$LOCAL_IP" ] && [ "$LOCAL_IP" != "0.0.0.0" ]; then
            if [ "$LOCAL_PORT" = "80" ] || [ "$LOCAL_PORT" = "443" ]; then
                UNRAID_API_URL="${PROTO}://${LOCAL_IP}/graphql"
            else
                UNRAID_API_URL="${PROTO}://${LOCAL_IP}:${LOCAL_PORT}/graphql"
            fi
        else
            UNRAID_API_URL="http://127.0.0.1/graphql"
        fi
    else
        UNRAID_API_URL="http://127.0.0.1/graphql"
    fi
fi
export UNRAID_API_URL

mkdir -p "$(dirname "$LOGFILE")"

nohup "$INSTALL_DIR/bin/unraid-agent" -config "$CONFIG_DIR/config.cfg" >> "$LOGFILE" 2>&1 &
PID=$!
echo $PID > "$PIDFILE"
echo "Server started (PID $PID)"
