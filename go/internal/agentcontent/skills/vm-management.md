---
name: vm-management
description: Virtual machine lifecycle on Unraid — list, start, stop, pause, resume, force stop, reboot, reset, plus per-VM permission overrides. Use for any request involving VMs.
---

# VM Management

## Read operations
- `vm_list` — all domains with id, name, state (running / shutoff / paused).

## State operations (need toggles)
- `vm_start`, `vm_stop` (graceful ACPI shutdown), `vm_pause`, `vm_resume`, `vm_reboot`, `vm_force_stop` (pull the plug — data loss risk like a hard power-off), `vm_reset` (hard reset).
- Toggles: ALLOW_VM_START/STOP/PAUSE/RESUME/REBOOT and the destructive ALLOW_VM_FORCE_STOP / ALLOW_VM_RESET, plus the ALLOW_VM_ACTIONS category.
- Per-VM overrides may allow or deny actions on a specific VM (Permissions → Virtual Machines).

## Good practice
- Always try graceful `vm_stop` first; escalate to `vm_force_stop` only if the guest is unresponsive and the user accepts the data-loss risk.
- `vm_reset` is for frozen guests; it is not a reboot — unsaved guest data is lost.
- VMs in paused state resume instantly with `vm_resume`.
