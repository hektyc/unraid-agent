# Connecting to the Unraid MCP Server

This server supports multiple transport modes and works with every major AI client.

## Transport Modes

| Mode | Description |
|---|---|
| `stdio` | Subprocess pipe (default, most compatible) |
| `streamable-http` | HTTP endpoint at `http://host:port/mcp` |
| `sse` | Server-Sent Events (deprecated) |

## Claude Desktop

```json
{
  "mcpServers": {
    "unraid": {
      "command": "npx",
      "args": ["-y", "unraid-agent"],
      "env": {
        "UNRAID_API_URL": "https://your-unraid.local/graphql",
        "UNRAID_API_KEY": "your-api-key"
      }
    }
  }
}
```

## Claude Code

```bash
# Add marketplace
/plugin marketplace add hektyc/unraid-agent

# Install plugin
/plugin install unraid-agent@unraid-agent
```

## VS Code

```json
{
  "servers": {
    "unraid": {
      "type": "stdio",
      "command": "npx",
      "args": ["-y", "unraid-agent"],
      "env": {
        "UNRAID_API_URL": "https://your-unraid.local/graphql",
        "UNRAID_API_KEY": "your-api-key"
      }
    }
  }
}
```

## Codex CLI

```bash
export UNRAID_API_URL="https://your-unraid.local/graphql"
export UNRAID_API_KEY="your-api-key"
codex
```

## HTTP Transports (Kilo, OpenCode, Remote)

```bash
export TRANSPORT=streamable-http
npx unraid-agent start

# MCP endpoint: http://localhost:6970/mcp
# Health check: http://localhost:6970/health
```

## Docker

```bash
docker run -it --rm \
  -e UNRAID_API_URL=https://your-unraid.local/graphql \
  -e UNRAID_API_KEY=your-api-key \
  -p 6970:6970 \
  unraid-agent
```
