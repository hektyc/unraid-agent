---
name: performance-triage
description: Read live CPU, memory, temperature, and network metrics to assess server load and cooling. Use for "is my server healthy/fast", fan/temp questions, or before recommending load changes.
---

# Performance Triage

## Tools (all read-only)
- `system_metrics` — CPU% and memory totals in one call (best first stop).
- `live_cpu`, `live_memory` — current utilization snapshots.
- `live_temperature` — all sensors: CPU, motherboard, drives; summary has average and hottest sensor.
- `live_network_metrics` — per-interface throughput (rx/tx per second).
- `system_overview` — hardware context: CPU model/cores, RAM size, motherboard, uptime.

## Interpretation
- Sustained CPU >85% with high iowait suggests disk-bound work (parity check, mover, Plex scans) — check `array_parity_status` and container activity.
- Memory above ~90% used on Unraid is often ZFS/cache usage, not leakage — compare uptime and swap before alarming.
- Drive temps: HDDs healthy <45°C, warn 45-50, critical >50. NVMe: warn >70, critical >80.
- A hot parity disk or cache drive during parity checks is normal short-term.

## Rules
- Metrics are point-in-time snapshots; for trends, sample a few times before concluding.
