---
name: permissions-model
description: The complete unraid-agent guardrail hierarchy — READ_ONLY, ALLOW_DESTRUCTIVE, per-entity overrides, category and specific toggles. Use to explain any "permission denied" error or to plan what a user must enable.
---

# Permissions Model

Every state-changing tool call passes this chain before reaching the Unraid API. Queries (reads) are always allowed.

## Resolution order
1. **unraid-agent self-removal** — always refused, no exceptions.
2. **READ_ONLY=true** (default) — blocks every mutation regardless of other settings.
3. **ALLOW_DESTRUCTIVE=true** — allows everything (explicit full-power mode).
4. **Per-entity overrides** (perms.json) — a container/VM/plugin set to Deny blocks that action on it; set to Allow permits it even when globals deny.
5. **Category toggles** — ALLOW_ARRAY_ACTIONS, ALLOW_DOCKER_ACTIONS, ALLOW_VM_ACTIONS cover their groups.
6. **Specific toggles** — one per action (ALLOW_CONTAINER_STOP, ALLOW_VM_FORCE_STOP, ALLOW_PLUGIN_REMOVE, ...).
7. Mutations with no toggle are allowed when not READ_ONLY (non-destructive only).

## Explaining denials to users
- "READ_ONLY mode is enabled" → Master Switch: Read Only = No, or keep it and accept read-only.
- "blocked by configuration — enable one of [...]" → the exact toggles to flip in Permissions; or set a per-entity override for that one container/VM.
- "explicitly denied ... per-entity settings" → the user set Deny on that entity's card; they must change it in the entity's modal.
- "Forbidden resource" (from the Unraid API, not this server) → the API key's role lacks that permission; adjust the key in Unraid → API Keys.

## API key layer
The plugin also authenticates to Unraid with an API key whose role independently gates mutations — both layers must permit an action.
