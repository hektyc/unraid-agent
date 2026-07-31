# Unraid MCP Server - Agent Skills

## How to Check Array Status
```
Use tool: array_parity_status
Use tool: system_array
```

## How to Start Parity Check
```
Use tool: array_parity_start with params {"correct": false}
Note: Requires ALLOW_ARRAY_START=true in config
```

## How to Manage Docker Containers
```
List: docker_list
Start: docker_start with container_id
Stop: docker_stop with container_id (requires ALLOW_CONTAINER_STOP)
Logs: Not available via GraphQL. Use docker logs <name> via SSH
```

## How to Manage VMs
```
List: vm_list
Start: vm_start with vm_id
Stop: vm_stop with vm_id (requires ALLOW_VM_STOP)
Force stop: vm_force_stop with confirm=true (requires ALLOW_VM_FORCE_STOP)
```

## How to Read Logs
```
disk_log_files - lists available logs
disk_logs - reads log content with path and optional tail_lines
```

## How to Manage Plugins
```
plugin_list - shows installed plugins
plugin_add - installs by name (requires ALLOW_PLUGIN_INSTALL)
plugin_remove - uninstalls (requires ALLOW_PLUGIN_REMOVE and confirm=true)
```

## Safety Rules
1. Always check READ_ONLY status before destructive operations
2. Destructive actions require explicit ALLOW_* toggles
3. Stop/start array requires ALLOW_ARRAY_ACTIONS
4. Container/VM stop requires ALLOW_CONTAINER_STOP or ALLOW_VM_STOP
5. Plugin operations require ALLOW_PLUGIN_INSTALL or ALLOW_PLUGIN_REMOVE
