package rclone

import (
	"context"
	"fmt"

	"github.com/hektyc/unraid-mcp-server/go/internal/config"
	"github.com/hektyc/unraid-mcp-server/go/internal/guards"
	"github.com/hektyc/unraid-mcp-server/go/internal/mcp"
)

func RegisterTools(s *mcp.Server, cfg *config.Config) {
	s.RegisterTool(&mcp.ToolDef{
		Name:        "rclone_list_remotes",
		Description: "Get configured remotes. Read-only.",
		Query:       `query { rclone { listRemotes { name type } } }`,
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "rclone_create_remote",
		Description: "Create a remote.",
		Query:       `mutation CreateRcloneRemote($name: String!, $config: JSON!) { rclone { createRemote(name: $name, config: $config) { name type } } }`,
		Params: map[string]string{
			"name": "string", // required=true
			"config": "object", // required=true
		},
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "rclone_delete_remote",
		Description: "Delete a remote.",
		Query:       `mutation DeleteRcloneRemote($name: String!) { rclone { deleteRemote(name: $name) { status } } }`,
		Params: map[string]string{
			"name": "string", // required=true
		},
	})
}