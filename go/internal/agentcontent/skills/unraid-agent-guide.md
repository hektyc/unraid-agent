---
name: unraid-agent-guide
description: The complete operator's guide to unraid-agent itself — architecture, permissions, tool-domain token management, skills and memory systems, diagnostics, and client setup. Use when the user asks how to configure, secure, or get the most out of the unraid-agent plugin or its agent content features.
---

# unraid-agent Operator's Guide

## Architecture in one minute
- A Go daemon (the MCP server) runs on the Unraid server and proxies the Unraid GraphQL API.
- Transports: streamable-http daemon (recommended; host/port/bearer auth) or stdio (client-spawned, usually via SSH; gated by Enable stdio Transport).
- WebGUI tabs: unRAID Agent (connection + service control), Permissions (guardrails), MCP Access (auth), Advanced (logging/TLS/diagnostics), Agent Content (skills + memory).

## Permissions model (safety first)
Resolution order for every state-changing tool call: self-removal guard → READ_ONLY → ALLOW_DESTRUCTIVE → per-entity overrides → category toggles → specific toggles. Reads are always allowed. Per-container/VM/plugin overrides live in perms.json and beat globals in both directions. unraid-agent can never be removed through the MCP server, ever.

## Tool Domains and token management (Permissions → Tool Access)
- Every tool family (19 domains, ~120 tools) has a toggle. Disabled domains are **never registered** — their tools and schemas never reach the client, which shrinks the standing token payload of every session.
- Use it to match the toolset to the job: a read-only monitoring agent can run with everything except system/docker/array/logs disabled; a docker-only helper can drop vm, rclone, oidc, onboarding.
- Apply restarts the daemon, so changes take effect immediately. Disabled tools cannot be called — keep the domains your workflows actually use.
- Token economics: skills and memories cost zero standing tokens (they only cost when fetched). Tool schemas are the standing cost — domain toggles are the lever.

## Agent Skills
- Playbooks shipped in the daemon and synced to /boot config on start; plus your custom skills.
- **Topic discipline: one topic per skill.** A skill body only enters the model's context when fetched, so large knowledge should be split into focused skills (e.g. array-health-check, container-troubleshooting) rather than one giant document — the model loads only what the task needs. Keep bodies tight; link related skills by name.
- Editing a default skill saves a copy-on-write override (custom wins everywhere); deleting the override restores the shipped default.
- Export any skill as SKILL.md from the Agent Content tab to install it into an IDE's skills folder (Kilo Code, Claude Code, Cursor, etc.).
- Agent-side management: skills_list/skills_get (read), skills_create/update/delete (Allow Skills Write toggle).

## Agent Memory
- Per-client scopes: each connected agent (kilo, claude, etc.) gets isolated memory; the defaults scope holds the auto-generated server profile.
- Edit the profile freely — nothing overwrites it. Delete it and a fresh one regenerates from live server data on next start.
- Agent-side: memory_list/memory_get (read), memory_write/delete (Allow Memory Write toggle).

## Diagnostics
- Advanced → Diagnostics shows the last 20 endpoint log entries (WebGUI reads/saves/deletes/exports).
- Agents can tail it remotely with the read-only agent_endpoint_log tool (tail param, default 50, max 200).

## Client setup
- Use the Client Configuration box on the unRAID Agent tab — it generates paste-ready MCP config for your transport (HTTP remote URL or SSH-stdio command), including bearer header when configured and the correct SSH port.
- For stdio: enable Enable stdio Transport while needed, use the SSH snippet, then disable it again — the binary refuses stdio entirely when the gate is off.
