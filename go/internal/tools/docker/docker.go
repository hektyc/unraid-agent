package docker

import (
	"github.com/hektyc/unraid-mcp-server/internal/config"
	"github.com/hektyc/unraid-mcp-server/internal/mcp"
)

func RegisterTools(s *mcp.Server, cfg *config.Config) {
	s.RegisterTool(&mcp.ToolDef{
		Name:        "docker_list",
		Description: "Get all containers. Read-only.",
		Query:       `query { docker { containers { id names image state status autoStart } } }`,
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "docker_networks",
		Description: "Get all networks. Read-only.",
		Query:       `query { docker { networks { id name driver scope } } }`,
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "docker_ports",
		Description: "Get host port bindings. Read-only.",
		Query:       `query { docker { containers { names ports } } }`,
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "docker_details",
		Description: "Get single container details.",
		Query:       `query($container_id: PrefixedID!) { docker { container(id: $container_id) { id names image imageId command created state status ports mounts labels autoStart autoStartOrder templatePath projectUrl registryUrl sizeRootFs sizeRw } } }`,
		Params: map[string]string{
			"container_id": "string",
		},
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "docker_network_details",
		Description: "Get single network details (returns all networks; match by id).",
		Query:       `query { docker { networks { id name driver scope } } }`,
		Params: map[string]string{
			"network_id": "string",
		},
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "docker_create_folder",
		Description: "Create an organizer folder.",
		Query:       `mutation($name: String!) { createDockerFolder(name: $name) { id } }`,
		Params: map[string]string{
			"name": "string",
		},
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "docker_delete_entries",
		Description: "Delete organizer entries.",
		Query:       `mutation($entry_ids: [PrefixedID!]!) { deleteDockerEntries(entryIds: $entry_ids) { id } }`,
		Params: map[string]string{
			"entry_ids": "array",
		},
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "docker_refresh_digests",
		Description: "Refresh image digests.",
		Query:       `mutation { refreshDockerDigests }`,
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "docker_logs",
		Description: "Get container logs (tail N lines).",
		Query:       `query($container_id: PrefixedID!, $tail: Int) { docker { logs(id: $container_id, tail: $tail) { containerId lines cursor } } }`,
		Params: map[string]string{
			"container_id": "string",
			"tail":         "number",
		},
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "docker_start",
		Description: "Start a container.",
		Query:       `mutation($container_id: PrefixedID!) { docker { start(id: $container_id) { id names state status } } }`,
		Params: map[string]string{
			"container_id": "string",
		},
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "docker_stop",
		Description: "Stop a container.",
		Query:       `mutation($container_id: PrefixedID!) { docker { stop(id: $container_id) { id names state status } } }`,
		Params: map[string]string{
			"container_id": "string",
		},
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "docker_restart",
		Description: "Restart a container.",
		Query:       `mutation($container_id: PrefixedID!) { docker { restart(id: $container_id) { id names state status } } }`,
		Params: map[string]string{
			"container_id": "string",
		},
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "docker_pause",
		Description: "Pause a container.",
		Query:       `mutation($container_id: PrefixedID!) { docker { pause(id: $container_id) { id names state status } } }`,
		Params: map[string]string{
			"container_id": "string",
		},
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "docker_unpause",
		Description: "Unpause a container.",
		Query:       `mutation($container_id: PrefixedID!) { docker { unpause(id: $container_id) { id names state status } } }`,
		Params: map[string]string{
			"container_id": "string",
		},
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "docker_remove_container",
		Description: "Remove a container.",
		Query:       `mutation($container_id: PrefixedID!, $withImage: Boolean) { docker { removeContainer(id: $container_id, withImage: $withImage) } }`,
		Params: map[string]string{
			"container_id": "string",
		},
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "docker_update_container",
		Description: "Apply pending image update.",
		Query:       `mutation($container_id: PrefixedID!) { docker { updateContainer(id: $container_id) { id names state status } } }`,
		Params: map[string]string{
			"container_id": "string",
		},
	})
}
