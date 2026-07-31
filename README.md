# Unraid MCP Server

![Version](https://img.shields.io/badge/version-0.0.1-blue)
![License](https://img.shields.io/badge/license-MIT-green)

Feature-rich, safety-enhanced MCP server for Unraid. Built with TypeScript, designed to work with every major AI client: Claude (Desktop/Code), Codex, OpenCode, Kilo, VS Code, Cursor, and more.

## Safety-First Design

| Safety Feature | Description |
|---|---|
| **Read-only mode** | Set `READ_ONLY=true` to block all write/mutation operations |
| **Granular toggles** | Individual `ALLOW_*` env vars for each destructive action |
| **Confirmation prompts** | Destructive tools require `confirm=true` in parameters |
| **Input validation** | All inputs validated through Zod schemas |
| **Sensitive data redaction** | API keys and tokens redacted from all logs |
| **TLS opt-in** | Skipping SSL verification requires explicit opt-in |
| **Bearer auth** | HTTP transport protected by static token unless behind a proxy |

## Feature Coverage

- **System**: 20+ tools for OS, CPU, memory, network, UPS, and time info
- **Array**: Parity checks, disk add/remove/mount, start/stop with safety guards
- **Docker**: Container lifecycle, networks, organizer folders, template management
- **VMs**: Start/stop/pause/resume/reboot with hard power-off and reset guards
- **Notifications**: Full CRUD with archive/unarchive
- **API Keys**: Role and permission management with preview
- **Plugins**: Install/uninstall with async operation tracking
- **Rclone**: Remote configuration and management
- **Settings**: System settings, UPS, SSH, time/NTP updates
- **Connect**: Unraid Connect remote access and device management
- **Customization**: Themes, locales, SSO
- **OIDC**: OpenID Connect provider management
- **Onboarding**: First-boot state management
- **Live Telemetry**: WebSocket subscriptions for CPU, memory, array, notifications, temperature, and more
- **Health**: Connection tests, diagnostics, and setup status

## Installation

```bash
npm install unraid-mcp-server
```

## Quick Start

```bash
# 1. Copy env file
cp .env.example .env

# 2. Configure your Unraid server
export UNRAID_API_URL=https://your-unraid-server.local/graphql
export UNRAID_API_KEY=your-api-key

# 3. Run in stdio mode (for Claude Desktop, Codex, etc.)
npm run dev

# Or run in HTTP mode
TRANSPORT=streamable-http npm start
```

## Transport Modes

| Transport | Description | Use Case |
|---|---|---|
| `stdio` | Subprocess pipe (default) | Claude Desktop, Claude Code, Codex, Cursor |
| `streamable-http` | HTTP endpoint on port | Remote access, proxies, Kilo, OpenCode |
| `sse` | Server-Sent Events (deprecated) | Legacy clients |

## Configuration

See [`.env.example`](.env.example) for all options. Key variables:

| Variable | Required | Default | Description |
|---|---|---|---|
| `UNRAID_API_URL` | Yes | — | GraphQL endpoint URL |
| `UNRAID_API_KEY` | Yes | — | Unraid API key |
| `TRANSPORT` | No | `stdio` | Transport mode |
| `UNRAID_MCP_HOST` | No | `127.0.0.1` | HTTP bind address |
| `UNRAID_MCP_PORT` | No | `6970` | HTTP port |
| `UNRAID_MCP_BEARER_TOKEN` | No | auto-generated | HTTP auth token |
| `UNRAID_MCP_DISABLE_HTTP_AUTH` | No | `false` | Disable bearer auth |
| `UNRAID_VERIFY_SSL` | No | `true` | TLS verification |
| `UNRAID_ALLOW_INSECURE_TLS` | No | `false` | Required to skip TLS |
| `UNRAID_MCP_LOG_LEVEL` | No | `info` | Log verbosity |
| `READ_ONLY` | No | `false` | Block all write operations |
| `ALLOW_*` | No | varies | Per-action destructive toggles |

## Client Integration

### Claude Desktop
Add to `claude_desktop_config.json`:
```json
{
  "mcpServers": {
    "unraid": {
      "command": "npx",
      "args": ["-y", "unraid-mcp-server"],
      "env": {
        "UNRAID_API_URL": "https://your-unraid.local/graphql",
        "UNRAID_API_KEY": "your-api-key"
      }
    }
  }
}
```

### Claude Code
```bash
# Local plugin
/plugin marketplace add hektyc/unraid-mcp-server
/plugin install unraid-mcp-server@unraid-mcp-server
```

### Codex
```bash
codex plugin marketplace add hektyc/unraid-mcp-server
```

### VS Code
Use the `mcp-server` configuration in `.vscode/mcp.json`:
```json
{
  "servers": {
    "unraid": {
      "type": "stdio",
      "command": "npx",
      "args": ["-y", "unraid-mcp-server"],
      "env": {
        "UNRAID_API_URL": "https://your-unraid.local/graphql",
        "UNRAID_API_KEY": "your-api-key"
      }
    }
  }
}
```

### Cursor
Add to `.cursor/mcp.json`:
```json
{
  "mcpServers": {
    "unraid": {
      "command": "npx",
      "args": ["-y", "unraid-mcp-server"],
      "env": {
        "UNRAID_API_URL": "https://your-unraid.local/graphql",
        "UNRAID_API_KEY": "your-api-key"
      }
    }
  }
}
```

### HTTP Proxies (Kilo, OpenCode, Remote)
```bash
export TRANSPORT=streamable-http
npx unraid-mcp-server start
# Connect via http://localhost:6970/mcp
```

## Security Best Practices

1. **Never expose HTTP transport directly** without auth or behind a reverse proxy
2. **Use `READ_ONLY=true`** when you only need monitoring
3. **Enable specific `ALLOW_*` toggles** instead of `ALLOW_DESTRUCTIVE=true`
4. **Rotate API keys regularly** via the `key` tools
5. **Use TLS** (`UNRAID_VERIFY_SSL=true`) for production deployments
6. **Set `UNRAID_MCP_BEARER_TOKEN`** to a strong random value for HTTP mode

## Development

```bash
# Install dependencies
npm install

# Run in dev mode with auto-reload
npm run dev

# Run tests
npm test

# Build for production
npm run build

# Lint and format
npm run lint
npm run format

# Type check
npm run typecheck
```

## Branching Strategy

- **dev**: Primary working branch. All features and fixes go here.
- **main**: Stable, tested releases. Auto-versioned on merge.

## Releasing

```bash
# Install the helper
chmod +x scripts/bump-version.sh

# Bump patch version (0.0.1 -> 0.0.2)
bash scripts/bump-version.sh patch

# Bump minor version
bash scripts/bump-version.sh minor

# CI automatically creates a GitHub Release on merge to main
```

## License

MIT
