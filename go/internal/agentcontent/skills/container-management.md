---
name: container-management
description: Day-to-day Docker container lifecycle — list, inspect, start, stop, restart, pause, update, and per-container permission overrides. Use for any request about apps, containers, or Docker services.
---

# Container Management

## Read operations (always allowed)
- `docker_list` — all containers with state/status/autostart.
- `docker_details` (container_id) — full inspect: mounts, labels, template path, registry URL, ports.
- `docker_logs` (container_id, tail) — recent log lines.
- `docker_networks`, `docker_ports` — networking layout and published ports.
- `docker_network_details` — per-network info.

## State operations (need toggles)
- `docker_start`, `docker_stop`, `docker_restart`, `docker_pause`, `docker_unpause`, `docker_update_container`, `docker_remove_container`.
- Each has its own toggle (ALLOW_CONTAINER_START/STOP/RESTART/PAUSE/UNPAUSE/REMOVE/UPDATE) plus the ALLOW_DOCKER_ACTIONS category.
- Per-container overrides may allow or deny a specific action on a specific container even when the global says otherwise — a "permission denied ... per-entity settings" error means the user set an override in Permissions → Containers.

## Good practice
- Before stopping a container, check what depends on it (databases, reverse proxies, tunnels are common linchpins).
- Prefer `docker_restart` over stop+start pairs.
- `docker_update_container` applies a pending image update; check for updates first in the WebGUI Docker page (update status tooling).
- Removing a container does not delete its appdata — tell the user their data stays in /mnt/user/appdata.
