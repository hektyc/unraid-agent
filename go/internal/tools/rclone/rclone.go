package rclone

import (
	"github.com/hektyc/unraid-mcp-server/internal/config"
	"github.com/hektyc/unraid-mcp-server/internal/mcp"
)

func RegisterTools(s *mcp.Server, cfg *config.Config) {
	s.RegisterTool(&mcp.ToolDef{
		Name:        "rclone_list_remotes",
		Description: "Get configured remotes. Read-only.",
		Query:       `query { rclone { remotes { name type parameters } } }`,
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "rclone_create_remote",
		Description: "Create a remote.",
		Query:       `mutation($name: String!, $config: JSON!) { rclone { createRCloneRemote(input: {name: $name, parameters: $config}) { name type } } }`,
		Params: map[string]string{
			"name":   "string",
			"config": "object",
		},
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "rclone_delete_remote",
		Description: "Delete a remote.",
		Query:       `mutation($name: String!) { rclone { deleteRCloneRemote(input: {name: $name}) } }`,
		Params: map[string]string{
			"name": "string",
		},
	})
}
