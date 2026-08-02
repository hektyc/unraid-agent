---
name: unraid-basics
description: Core Unraid concepts any agent must know — array, parity, cache pools, shares, flash drive, appdata, and Docker templates. Use before answering architecture or storage questions.
---

# Unraid Basics

## Storage model
- **Array**: the main storage. Data disks + up to 2 parity disks. Parity lets the array survive 1-2 disk failures and rebuild.
- **Parity is not a backup.** It protects against disk failure, not deletion, ransomware, or fire.
- **Cache pools**: fast SSD/NVMe pools holding new writes and appdata. The "mover" process relocates files from cache to array on a schedule.
- **Shares**: user-facing folders under /mnt/user/<share>. A share's settings decide whether files go to cache first, array only, or cache-only.
- Mount points: /mnt/user (shares), /mnt/diskN (array disks), /mnt/<pool> (cache pools), /boot (USB flash).

## Flash drive (/boot)
- Boots the OS and holds all configuration: /boot/config/plugins, /boot/config/domains (VMs), shares, docker templates.
- Everything under /boot/config persists across reboots. Plugin configs live in /boot/config/plugins/<name>/.

## Appdata
- /mnt/user/appdata/<container> holds each Docker container's persistent data (configs, databases).
- Mapping appdata into containers via bind mounts is the standard pattern.

## Docker on Unraid
- Containers are managed through templates (XML in /boot/config/plugins/dockerMan/templates-user/).
- Networks: bridge (default NAT), host (container shares host IP), macvlan/custom (container gets its own LAN IP), plus user-defined bridge networks.
- The icon shown in UIs usually comes from the template label net.unraid.docker.icon.

## Virtual Machines
- Defined by libvirt XML in /boot/config/domains/<name>/. States: running, shutoff, paused.
- Disk images live in /mnt/user/domains by default.
