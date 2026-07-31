# Unraid MCP Server

Private, safety-first MCP server for Unraid. Runs as a native Unraid plugin (`.plg`) with WebGUI settings, GraphQL-only delegation, and `nobody:nobody` daemon isolation.

## Architecture

```
Unraid WebUI (Settings > Plugins > Unraid MCP Server)
         |
    plugin.xml + Settings.page
         |
    scripts/*.sh (start/stop/install/remove)
         |
    /usr/local/emhttp/plugins/unraid-mcp/bin/unraid-mcp  (Go binary, nobody:nobody)
         |
    GraphQL only -> https://<unraid>:/graphql (X-API-Key)
```

- **Transport:** `stdio` (default) or `streamable-http` (opt-in)
- **Auth:** Bearer token for HTTP mode (manual entry, never auto-generated)
- **Safety:** `READ_ONLY=true` by default; per-action `ALLOW_*` toggles; explicit warnings when read-only is disabled
- **Memory:** `AGENT-SKILLS.md` + `AGENT-MEMORY.md` stored in plugin config dir

## Install

1. Build the plugin:
   ```bash
   ./build.sh
   ```
   This creates `unraid-mcp-0.0.1.plg`.

2. Upload via Unraid WebUI:
   - Go to **Settings > Plugins > Install Plugin**
   - Enter the raw URL of the `.plg` file:
     ```
     https://raw.githubusercontent.com/hektyc/unraid-mcp-server/dev/unraid-mcp-0.0.1.plg
     ```
   - Click **Install**

3. Configure:
   - Go to **Settings > Unraid MCP Server**
   - Enter your Unraid API URL and API Key
   - Keep **Read Only** enabled unless you need write access
   - Adjust action toggles as needed

4. Start:
   - Go to **Settings > Plugins**
   - Find **Unraid MCP Server** and click **Start**
   - Logs are at `/var/log/unraid-mcp.log`

## Client Integration

### Claude Desktop (stdio)
```json
{
  "mcpServers": {
    "unraid": {
      "command": "ssh",
      "args": ["root@tower.local", "/usr/local/emhttp/plugins/unraid-mcp/bin/unraid-mcp"]
    }
  }
}
```

### HTTP (streamable-http)
```json
{
  "mcpServers": {
    "unraid": {
      "url": "http://tower.local:6970/mcp",
      "headers": {
        "Authorization": "Bearer YOUR_TOKEN"
      }
    }
  }
}
```

## AGENT-SKILLS.md

Stored at `/boot/config/plugins/unraid-mcp/AGENT-SKILLS.md`. Edit to teach the agent how to carry out tasks on your system.

## AGENT-MEMORY.md

Stored at `/boot/config/plugins/unraid-mcp/AGENT-MEMORY.md`. Edit to add rules and remembered state.

## Build from Source

```bash
# Requires Go 1.24+
go build -o go/bin/unraid-mcp ./go/cmd/unraid-mcp
./build.sh
```

## Branching

- `dev` — active development
- `main` — stable releases

Version bump on merge to `main` creates a new `.plg` tagged release.

## Security

- MCP server runs as `nobody:nobody`
- All admin actions go through GraphQL only
- No shell execution
- Default `READ_ONLY=true`
- Sensitive values redacted in logs
