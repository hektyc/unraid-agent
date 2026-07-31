# Unraid MCP Server

Private, safety-first MCP server for Unraid. Runs as a native Unraid plugin (`.plg`) with WebGUI settings, GraphQL-only delegation, and `nobody:nobody` daemon isolation.

## Quick Install

### Option A: GitHub Actions Build (Recommended)

1. Push code to the `dev` branch
2. Go to **Actions** tab in GitHub
3. Select **Build Unraid MCP Plugin** workflow
4. Click **Run workflow** → select `dev` branch → **Run**
5. When complete, download `unraid-mcp-plugin` artifact
6. In Unraid WebUI: **Settings > Plugins > Install Plugin**
7. Upload the downloaded `.plg` file

### Option B: Direct Download from GitHub Releases

1. Go to **Releases** page in GitHub
2. Download the latest `unraid-mcp.plg` asset
3. In Unraid WebUI: **Settings > Plugins > Install Plugin**
4. Paste the direct download URL or upload the file

### Option C: Build Locally

```bash
# Requires Go 1.24+
git clone https://github.com/hektyc/unraid-mcp-server.git
cd unraid-mcp-server
git checkout dev
chmod +x build.sh
./build.sh
# → creates: unraid-mcp.plg
```

Then upload the `.plg` via Unraid WebUI.

### Option D: Manual File Copy

If you cannot use the UI:

```bash
# Copy .plg to Unraid server
scp unraid-mcp.plg root@tower.local:/root/

# SSH into Unraid and install
ssh root@tower.local
cd /root
tar -xf unraid-mcp.plg -C /
```

## Configuration

After install, go to **Settings > Unraid MCP Server**:

| Setting | Default | Description |
|---|---|---|
| Unraid API URL | `http://localhost/graphql` | GraphQL endpoint |
| API Key | *(empty)* | Unraid API key |
| Transport | `stdio` | `stdio` or `streamable-http` |
| Bind Host | `127.0.0.1` | HTTP bind address |
| Port | `6970` | HTTP port |
| Read Only | `true` | **Keep enabled unless needed** |
| Allow Array Actions | `false` | Array start/stop |
| Allow Container Actions | `false` | Container stop/remove |
| Allow VM Actions | `false` | VM stop/force-stop |
| Allow Destructive | `false` | Shortcut for all destructive |
| Bearer Token | *(empty)* | Manual entry for HTTP auth |
| Verify SSL | `true` | TLS verification |

## Safety

- **Default is read-only.** All write/mutation operations are blocked unless explicitly enabled.
- Per-action toggles in the WebGUI control exactly which operations are allowed.
- Destructive actions show explicit warnings in tool descriptions.
- MCP server runs as `nobody:nobody` with GraphQL-only access.

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

## Files

| Path | Purpose |
|---|---|
| `/usr/local/emhttp/plugins/unraid-mcp/` | Plugin install dir |
| `/boot/config/plugins/unraid-mcp/config.cfg` | Persistent settings |
| `/boot/config/plugins/unraid-mcp/AGENT-SKILLS.md` | Agent skills (editable) |
| `/boot/config/plugins/unraid-mcp/AGENT-MEMORY.md` | Agent memory (editable) |
| `/var/log/unraid-mcp.log` | Plugin logs |
| `/usr/local/emhttp/plugins/unraid-mcp/bin/unraid-mcp` | Go binary |

## AGENT-SKILLS.md and AGENT-MEMORY.md

These files live in `/boot/config/plugins/unraid-mcp/` and can be edited by you or by the AI agent through the MCP tools. They provide:

- **AGENT-SKILLS.md** — How-to guides for carrying out tasks
- **AGENT-MEMORY.md** — Rules, constraints, and remembered state

## Branching

- `dev` — active development, CI builds `.plg` artifacts
- `main` — stable releases, auto-tagged `.plg` builds with GitHub Release assets

## Security

- `nobody:nobody` process isolation
- GraphQL-only delegation (no shell execution)
- Default `READ_ONLY=true`
- Sensitive values redacted in logs
- Manual bearer token entry (never auto-generated)
