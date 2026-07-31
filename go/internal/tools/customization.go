package tools

import (
	"context"
	"fmt"

	"github.com/hektyc/unraid-mcp-server/internal/config"
	"github.com/hektyc/unraid-mcp-server/internal/guards"
	"github.com/hektyc/unraid-mcp-server/internal/mcp"
)

func RegistercustomizationTools(s *mcp.Server, cfg *config.Config) {
	s.RegisterTool(&mcp.ToolDef{
		Name:        "customization_themes",
		Description: "Get available themes. Read-only.",
		Query:       `query { customization { themes { id name } } }`,
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "customization_locales",
		Description: "Get available locales. Read-only.",
		Query:       `query { customization { locales { id name } } }`,
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "customization_set_theme",
		Description: "Set UI theme.",
		Query:       `mutation SetTheme($themeId: String!) { setTheme(themeId: $themeId) { status } }`,
		Params: map[string]string{
			"theme_id": "string", // required=true
		},
	})
}