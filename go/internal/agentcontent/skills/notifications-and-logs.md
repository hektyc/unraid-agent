---
name: notifications-and-logs
description: Work with Unraid notifications and log sources — triage unread alerts, create custom notifications, and read Docker container logs. Use for alert review or when the user asks "what is my server trying to tell me".
---

# Notifications & Logs

## Notifications
- `notification_overview` — counts of unread and archived by severity (info/warning/alert).
- `notification_list` — browse with filters: type (UNREAD/ARCHIVE), importance (INFO/WARNING/ALERT), offset/limit.
- `notification_archive`, `notification_delete`, `notification_delete_archived` — housekeeping (delete needs ALLOW_NOTIFICATION_DELETE).
- `notification_create` — post your own notification (title, subject, description, importance). Useful to surface agent findings in the WebGUI bell. It is a mutation but non-destructive.

## Triage workflow
1. overview → 2. list unread alerts first (importance=ALERT) → 3. summarize for the user in plain language, grouped by subsystem → 4. propose actions (array issues → array-health-check skill; container issues → container-troubleshooting).

## Docker logs
- `docker_logs` (container_id, tail) — last N log lines for a container. The primary in-band log source; see container-troubleshooting skill for patterns.
