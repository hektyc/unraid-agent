package customization

import (
	"github.com/hektyc/unraid-mcp-server/internal/config"
	"github.com/hektyc/unraid-mcp-server/internal/mcp"
)

func RegisterTools(s *mcp.Server, cfg *config.Config) {
	s.RegisterTool(&mcp.ToolDef{
		Name:        "customization_themes",
		Description: "Get available themes. Read-only.",
		Query:       `query { publicTheme { name showBannerImage showBannerGradient showHeaderDescription headerBackgroundColor headerPrimaryTextColor headerSecondaryTextColor } display { theme } }`,
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "customization_locales",
		Description: "Get available locales. Read-only.",
		Query:       `query { customization { availableLanguages { code name url } } }`,
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "customization_set_theme",
		Description: "Set UI theme (azure, black, gray, white).",
		Query:       `mutation($theme_id: ThemeName!) { customization { setTheme(theme: $theme_id) { name } } }`,
		Params: map[string]string{
			"theme_id": "string",
		},
	})
}
