---
name: array-management
description: Start/stop the array and add or remove disks — the most destructive operations available. Use when the user explicitly asks for array power actions or disk membership changes, and only after confirming permission toggles.
---

# Array Management

## Available operations
- `array_start` / `array_stop` — power the array up or down. Stopping unmounts all shares and stops dependent services; treat as a maintenance-window action.
- `array_add_disk` (disk_id, optional slot) — add an unassigned disk. New array disks are cleared (wiped) by Unraid before use — data on them is destroyed.
- `array_remove_disk` (disk_id) — remove a disk; the array may need to rebuild parity onto remaining disks afterward.

## Pre-flight checklist (always do this first)
1. `system_array` — confirm current state and which disk is involved.
2. `array_assignable_disks` — for add operations, list candidate disks (id, size, serial) and confirm the exact one with the user by serial number.
3. Check permission toggles: ALLOW_ARRAY_START / ALLOW_ARRAY_STOP / ALLOW_ARRAY_ADD_DISK / ALLOW_ARRAY_REMOVE_DISK or ALLOW_ARRAY_ACTIONS. If denied, report the exact toggle instead of retrying.
4. Confirm with the user in plain language: what will happen, that it may take a long time, and that new disks are wiped.

## Rules
- Never stop the array as a troubleshooting shortcut.
- Never add/remove disks without explicit user confirmation of the exact disk serial.
- Parity checks (`array_parity_start/pause/resume/cancel`) are safe operations but still report progress via `array_parity_status`.
