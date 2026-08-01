# Unraid MCP Server - Agent Memory

## Server State
- MCP Server Version: 0.0.1
- Transport: stdio (default) or streamable-http
- Read-only mode: true by default
- Process runs as: nobody:nobody
- Config location: /boot/config/plugins/unraid-agent/config.cfg
- Binary location: /usr/local/emhttp/plugins/unraid-agent/bin/unraid-agent

## Security Rules
1. Always enforce READ_ONLY mode if enabled
2. Always check ALLOW_* toggles for destructive actions
3. Require confirm=true for disk_remove_disk, disk_clear_disk_stats, array_stop, vm_force_stop, vm_reset, notification_delete, key_delete, plugin_remove, plugin_install, rclone_delete_remote
4. Redact API keys and tokens from all logs
5. Never expose HTTP transport without auth unless explicitly disabled
6. Maintain bounded action surface: only GraphQL delegation, no shell execution

## GraphQL Endpoint
- URL: https://<unraid-host>/graphql
- Auth: X-API-Key header
- Subscription: WebSocket to wss://<unraid-host>/graphql

## Known Limitations
- Docker container logs not available via GraphQL
- Some advanced Unraid features require SSH fallback
- Plugin operations are async (track operationId)
