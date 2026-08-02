---
name: system-administration
description: Server-level administration — OS version and uptime, settings, time/NTP, registration, flash drive, services, and safe handling of system mutations. Use for OS-level questions or setting changes.
---

# System Administration

## Read operations
- `system_overview`, `system_server`, `system_server_time`, `system_services`, `system_settings`, `system_variables`, `system_registration`, `system_flash`, `system_config`, `system_display`, `system_online`, `system_owner`, `system_servers`, `system_timezones`, `user_me`.

## Key facts
- `system_variables` exposes Unraid's var.ini: version, name, ports, SSL mode, array start config, share counts, flash GUID/vendor.
- `system_services` shows unraid-api version and online state — version mismatches explain schema differences.
- `system_config` reports config validity; invalid config can explain odd behavior.

## Mutations (need toggles)
- `setting_update`, `setting_update_ssh` (ALLOW_SSH_UPDATE), `setting_update_system_time` (ALLOW_TIME_UPDATE), `customization_set_theme`, `setting_configure_ups`.
- These change the running OS — always confirm the exact change with the user first and expect some to require service restarts.

## Rules
- Prefer telling the user which WebGUI page owns a setting over changing it via tools when the change is complex.
