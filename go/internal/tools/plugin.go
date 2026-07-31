package tools

import (
	"context"
	"fmt"

	"github.com/hektyc/unraid-mcp-server/go/internal/config"
	"github.com/hektyc/unraid-mcp-server/go/internal/guards"
	"github.com/hektyc/unraid-mcp-server/go/internal/mcp"
)

func RegisterTools(s *mcp.Server, cfg *config.Config) {
	s.RegisterTool(&mcp.ToolDef{
		Name:        "plugin_list",
		Description: "Get installed plugins. Read-only.",
		Query:       `query { plugin { list { name version status } } }`,
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "plugin_add",
		Description: "Install plugins by name.",
		Query:       `mutation AddPlugins($names: [String!]!) { plugin { add(names: $names) { status } } }`,
		Params: map[string]string{
			"names": "array", // required=true
		},
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "plugin_remove",
		Description: "Uninstall plugins.",
		Query:       `mutation RemovePlugins($names: [String!]!) { plugin { remove(names: $names) { status } } }`,
		Params: map[string]string{
			"names": "array", // required=true
			"confirm": "boolean", // required=true
		},
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "plugin_install",
		Description: "Async install a .plg URL.",
		Query:       `mutation InstallPlugin($url: String!) { plugin { install(url: $url) { operationId } } }`,
		Params: map[string]string{
			"url": "string", // required=true
			"confirm": "boolean", // required=true
		},
	})
}