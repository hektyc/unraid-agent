---
name: storage-and-shares
description: Understand storage layout and capacity — array capacity, per-disk usage and health, unassigned disks, and how shares consume space. Use for free-space, disk-inventory, or expansion questions.
---

# Storage & Shares

## Tools
- `system_array` — the workhorse: capacity (total/used/free for the array and per disk), every data/parity/cache disk with status, temperature, filesystem size/free.
- `array_assignable_disks` — disks present but not in the array: candidates for expansion (id, vendor, size, serial, interface, SMART status, temperature).

## Reading capacity
- Array usable capacity = sum of data disks (parity disks store no files). Free space trends matter more than a single snapshot — ask the user about growth rate before recommending expansion.
- A cache pool near full breaks the cache-first write path for shares — new writes may go straight to slower array disks or fail for cache-only shares.

## Expansion guidance
- Adding a disk: pick from `array_assignable_disks` by serial; new disks are wiped during add — confirm with the user (see array-management skill).
- Cache expansion follows the same pattern for pools via the WebGUI.

## Rules
- Storage layout questions are read-only by nature; only disk membership changes need permissions and confirmation.
