---
name: container-troubleshooting
description: Debug failing or unhealthy containers — read exit states and logs, spot port conflicts, and identify pending image updates. Use when a container is stopped, crashing, unhealthy, or unreachable.
---

# Container Troubleshooting

## Diagnostic workflow
1. `docker_list` — find the container; note state and status string.
2. `docker_details` — check mounts (does the appdata path exist?), labels (web UI URL), restart-related labels.
3. `docker_logs` with tail 50-100 — the last lines usually name the failure (missing env, permission denied on bind mount, port in use, crashed process).
4. `docker_ports` — verify the expected host port is actually published and not claimed by another container.

## Common patterns
- **Exited (1) or (137) right after start**: crash on boot — read logs; usually config or permissions on appdata.
- **Exited (0)**: clean exit; some containers are one-shot jobs by design.
- **"port already allocated"**: another container or service owns the host port — use `docker_ports` to find the conflict; fix by changing one container's port mapping (WebGUI edit), never by force.
- **Unhealthy status**: the container's healthcheck is failing even though the process runs — logs + the app's own URL tell more.
- **Updates pending**: containers show update status in the WebGUI Docker page; `docker_update_container` applies one at a time after user approval.

## Rules
- Diagnose read-only first; propose the fix; only act after user confirmation and permission checks.
