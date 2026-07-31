package tools

import (
	"context"
	"fmt"

	"github.com/hektyc/unraid-mcp-server/internal/config"
	"github.com/hektyc/unraid-mcp-server/internal/guards"
	"github.com/hektyc/unraid-mcp-server/internal/mcp"
)

func RegisterTools(s *mcp.Server, cfg *config.Config) {
	s.RegisterTool(&mcp.ToolDef{
		Name:        "docker_list",
		Description: "Get all containers. Read-only.",
		Query:       `query { docker { containers { id names image state status autoStart } } }`,
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "docker_details",
		Description: "Get single container details.",
		Query:       `query DockerDetails($id: ID!) { docker { container(id: $id) { id names image state status ports mounts } } }`,
		Params: map[string]string{
			"container_id": "string", // required=true
		},
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "docker_logs",
		Description: "Not available via GraphQL API.",
		Query:       `SELECT 'Not available via GraphQL. Use docker logs ' || ? AS info`,
		Params: map[string]string{
			"container_id": "string", // required=true
		},
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "docker_ports",
		Description: "Get host port bindings. Read-only.",
		Query:       `query { docker { containers { ports { hostPort containerPort protocol } } } }`,
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "docker_start",
		Description: "Start a container.",
		Query:       `mutation StartContainer($id: ID!) { startDockerContainer(id: $id) { status } }`,
		Params: map[string]string{
			"container_id": "string", // required=true
		},
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "docker_stop",
		Description: "Stop a container.",
		Query:       `mutation StopContainer($id: ID!) { stopDockerContainer(id: $id) { status } }`,
		Params: map[string]string{
			"container_id": "string", // required=true
		},
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "docker_restart",
		Description: "Restart a container.",
		Query:       `mutation RestartContainer($id: ID!) { restartDockerContainer(id: $id) { status } }`,
		Params: map[string]string{
			"container_id": "string", // required=true
		},
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "docker_unpause",
		Description: "Unpause a container.",
		Query:       `mutation UnpauseContainer($id: ID!) { unpauseDockerContainer(id: $id) { status } }`,
		Params: map[string]string{
			"container_id": "string", // required=true
		},
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "docker_remove_container",
		Description: "Remove a container.",
		Query:       `mutation RemoveContainer($id: ID!, $wi: Boolean) { removeDockerContainer(id: $id, withImage: $wi) { status } }`,
		Params: map[string]string{
			"container_id": "string", // required=true
			"confirm": "boolean", // required=true
		},
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "docker_update_container",
		Description: "Apply pending image update.",
		Query:       `mutation UpdateContainer($id: ID!) { updateDockerContainer(id: $id) { status } }`,
		Params: map[string]string{
			"container_id": "string", // required=true
		},
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "docker_refresh_digests",
		Description: "Refresh image digests.",
		Query:       `mutation { refreshDockerDigests { status } }`,
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "docker_networks",
		Description: "Get all networks. Read-only.",
		Query:       `query { docker { networks { id name driver scope } } }`,
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "docker_network_details",
		Description: "Get single network details.",
		Query:       `query NetworkDetails($id: ID!) { docker { network(id: $id) { id name driver scope } } }`,
		Params: map[string]string{
			"network_id": "string", // required=true
		},
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "docker_create_folder",
		Description: "Create an organizer folder.",
		Query:       `mutation CreateFolder($name: String!) { createOrganizerFolder(name: $name) { id name } }`,
		Params: map[string]string{
			"name": "string", // required=true
		},
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "docker_delete_entries",
		Description: "Delete organizer entries.",
		Query:       `mutation DeleteEntries($ids: [ID!]!) { deleteOrganizerEntries(entryIds: $ids) { status } }`,
		Params: map[string]string{
			"entry_ids": "array", // required=true
			"confirm": "boolean", // required=true
		},
	})
}