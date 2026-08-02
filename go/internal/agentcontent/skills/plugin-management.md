---
name: plugin-management
description: List installed plugins, install new ones by name or URL, and remove plugins — including the absolute rule that unraid-agent can never be removed through this server. Use for plugin inventory, installs, or removals.
---

# Plugin Management

## Read operations
- `plugin_list` — everything installed. Note: the API's plugin registry only tracks API-module plugins; the authoritative installed set comes from the plugin files on the flash drive, which this tool reflects via installedUnraidPlugins.

## State operations (need toggles)
- `plugin_add` (names[]) — install by name from the plugin registry (ALLOW_PLUGIN_INSTALL).
- `plugin_install` (url) — async install from a .plg URL; track progress in the WebGUI Plugins page.
- `plugin_remove` (names[]) — uninstall (ALLOW_PLUGIN_REMOVE). Per-plugin overrides in Permissions → Plugins can allow or deny removal of individual plugins.

## Absolute rule
- **unraid-agent can never be removed through this MCP server** — any plugin_remove naming it is refused regardless of permissions. Removing the agent through the agent would kill the management channel. If a user insists, direct them to the WebGUI Plugins page.

## Good practice
- Install operations are async — after plugin_add/plugin_install, tell the user to watch the Plugins page for completion.
- Prefer removing one plugin per call so per-plugin deny verdicts are unambiguous.
