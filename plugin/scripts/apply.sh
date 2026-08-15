#!/bin/bash
PLUGIN_NAME="unraid-agent"
INSTALL_DIR="/usr/local/emhttp/plugins/${PLUGIN_NAME}"
CONFIG_DIR="/boot/config/plugins/${PLUGIN_NAME}"
PIDFILE="$CONFIG_DIR/server.pid"

source "$CONFIG_DIR/config.cfg"

set -a
source "$CONFIG_DIR/config.cfg"
set +a

if [ -z "${UNRAID_API_URL:-}" ]; then
    INI_FILE="/var/local/emhttp/var.ini"
    if [ -f "$INI_FILE" ]; then
        strip_quotes() { echo "$1" | sed 's/^"//;s/"$//;s/^'\''//;s/'\''$//'; }
        LOCAL_IP=$(strip_quotes "$(grep '^IPADDR=' "$INI_FILE" | head -1 | cut -d'=' -f2-)")
        USE_SSL=$(strip_quotes "$(grep '^USE_SSL=' "$INI_FILE" | head -1 | cut -d'=' -f2-)")
        [ -z "$USE_SSL" ] && USE_SSL=$(strip_quotes "$(grep '^USESSL=' "$INI_FILE" | head -1 | cut -d'=' -f2-)")
        PORT=$(strip_quotes "$(grep '^PORT=' "$INI_FILE" | head -1 | cut -d'=' -f2-)")
        PORTSSL=$(strip_quotes "$(grep '^PORTSSL=' "$INI_FILE" | head -1 | cut -d'=' -f2-)")
        [ -z "$PORTSSL" ] && PORTSSL=$(strip_quotes "$(grep '^PORT_SSL=' "$INI_FILE" | head -1 | cut -d'=' -f2-)")

        if [ "$USE_SSL" = "yes" ] || [ "$USE_SSL" = "true" ] || [ "$USE_SSL" = "1" ]; then
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

if [ -f "$PIDFILE" ]; then
    OLDPID=$(cat "$PIDFILE")
    if kill -0 "$OLDPID" 2>/dev/null; then
        kill "$OLDPID" 2>/dev/null
        sleep 1
        kill -9 "$OLDPID" 2>/dev/null
        rm -f "$PIDFILE"
    fi
fi

# stdio does not run as a background daemon — stop here after killing
# any previously running daemon (so switching to stdio stops the HTTP server).
if [ "${TRANSPORT:-streamable-http}" = "stdio" ]; then
    echo "stdio transport: daemon not started (client-spawned transport)."
    exit 0
fi

nohup "$INSTALL_DIR/bin/unraid-agent" -config "$CONFIG_DIR/config.cfg" >> /var/log/unraid-agent.log 2>&1 &
PID=$!
echo $PID > "$PIDFILE"
echo "Server restarted (PID $PID)"
