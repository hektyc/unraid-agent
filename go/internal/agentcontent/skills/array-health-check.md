---
name: array-health-check
description: Diagnose array and disk health — read array state, parity check status and history, identify disabled disks (DISK_DSBL), and triage notifications. Use whenever the user mentions errors, warnings, parity, or a disk problem.
---

# Array Health Check

## Workflow
1. `system_array` — overall state (STARTED/STOPPED), capacity, per-disk status, temps, filesystem usage.
2. `array_parity_status` — is a parity check running now? progress, speed, errors so far.
3. `array_parity_history` — past checks: dates, duration, error counts, whether they completed or were cancelled.
4. `notification_list` (importance=ALERT/WARNING) — current alerts the system has raised.

## Interpretation
- Disk status **DISK_DSBL**: the disk is disabled — writes to it stopped, array runs degraded using parity to simulate its contents. URGENT: tell the user immediately. Typical next steps in the WebGUI: review SMART/syslog, then replace or rebuild the disk. Do not attempt repairs via tools.
- **Parity errors > 0**: data mismatch between disks and parity. Many errors after an unclean shutdown can be benign-looking but need a correcting check (ALLOW_ARRAY_START required to start one — ask the user first).
- **CANCELLED checks**: a check was interrupted (shutdown or manual). Not an error by itself, but frequent cancellations plus rising error counts deserve attention.
- Hot disks: spinning disks >45°C or NVMe >70°C sustained deserve a cooling warning.

## Rules
- This skill is read-only diagnostics. Starting a correcting parity check is a state change — confirm with the user and check permissions first.
