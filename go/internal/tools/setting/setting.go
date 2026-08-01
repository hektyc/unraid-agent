package setting

import (
	"github.com/hektyc/unraid-mcp-server/internal/config"
	"github.com/hektyc/unraid-mcp-server/internal/mcp"
)

func RegisterTools(s *mcp.Server, cfg *config.Config) {
	s.RegisterTool(&mcp.ToolDef{
		Name:        "setting_list",
		Description: "Get all settings. Read-only.",
		Query:       `query { settings { unified { values } } }`,
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "setting_get",
		Description: "Get single setting by name. Read-only.",
		Query:       `query { settings { unified { dataSchema uiSchema values } } }`,
		Params: map[string]string{
			"name": "string",
		},
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "setting_update",
		Description: "Update settings (pass a JSON object of setting keys to new values).",
		Query:       `mutation($settings: JSON!) { updateSettings(input: $settings) { __typename } }`,
		Params: map[string]string{
			"settings": "object",
		},
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "setting_update_ssh",
		Description: "Update SSH settings.",
		Query:       `mutation($enabled: Boolean!) { updateSshSettings(input: {enabled: $enabled}) { useSsh portssh } }`,
		Params: map[string]string{
			"enabled": "boolean",
		},
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "setting_update_system_time",
		Description: "Update system time/NTP.",
		Query:       `mutation($timezone: String!) { updateSystemTime(input: {timeZone: $timezone}) { currentTime timeZone useNtp } }`,
		Params: map[string]string{
			"timezone": "string",
		},
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "setting_configure_ups",
		Description: "Configure UPS.",
		Query:       `mutation($driver: String!, $port: String!) { configureUps(config: {upsType: $driver, device: $port}) }`,
		Params: map[string]string{
			"driver": "string",
			"port":   "string",
		},
	})
}
