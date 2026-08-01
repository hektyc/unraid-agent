package plugin

import (
	"github.com/hektyc/unraid-mcp-server/internal/config"
	"github.com/hektyc/unraid-mcp-server/internal/mcp"
)

func RegisterTools(s *mcp.Server, cfg *config.Config) {
	s.RegisterTool(&mcp.ToolDef{
		Name:        "plugin_list",
		Description: "Get installed plugins. Read-only.",
		Query:       `query { plugins { name version hasApiModule hasCliModule } installedUnraidPlugins }`,
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "plugin_add",
		Description: "Install plugins by name.",
		Query:       `mutation($names: [String!]!) { addPlugin(input: {names: $names}) }`,
		Params: map[string]string{
			"names": "array",
		},
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "plugin_remove",
		Description: "Uninstall plugins.",
		Query:       `mutation($names: [String!]!) { removePlugin(input: {names: $names}) }`,
		Params: map[string]string{
			"names": "array",
		},
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "plugin_install",
		Description: "Async install a .plg URL.",
		Query:       `mutation($url: String!) { unraidPlugins { installPlugin(input: {url: $url}) { id url name status createdAt } } }`,
		Params: map[string]string{
			"url": "string",
		},
	})
}
