---
name: network-and-access
description: Network layout and remote access — interfaces and their roles, WebGUI access URLs, and how users reach the server remotely (LAN, DNS, VPN/Tailscale, Unraid Connect). Use for connectivity, URL, or remote-access questions.
---

# Network & Access

## Tools
- `system_network` — all WebGUI access URLs (LAN IP, mDNS name, WAN/remote entries) with types.
- `system_network_interfaces` — every interface: physical NICs, bonding, VLANs, docker bridge, VPN tunnels (e.g. tailscale), with IPs, speed, duplex, operstate.
- `system_network_metrics` / `live_network_metrics` — throughput per interface.

## Concepts
- Unraid typically has: the LAN interface(s), optional VLANs (eth0.x), docker0 bridge, and VPN interfaces (Tailscale etc.).
- Access URLs: LAN IP is the direct path; the <name>.lan mDNS name works on most home networks; remote access may go through Unraid Connect dynamic DNS or a VPN.
- Containers with their own IPs appear on macvlan/custom networks — see docker_networks for those.

## Troubleshooting
- Interface down (operstate down) on the main NIC = total outage for local users — physical layer first.
- A WebGUI reachable by IP but not by name = mDNS/DNS issue on the client side, not the server.
