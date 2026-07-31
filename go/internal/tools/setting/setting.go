package setting

import (

	"github.com/hektyc/unraid-mcp-server/internal/config"
	"github.com/hektyc/unraid-mcp-server/internal/mcp"
)

func RegisterTools(s *mcp.Server, cfg *config.Config) {
	s.RegisterTool(&mcp.ToolDef{
		Name:        "setting_list",
		Description: "Get all settings. Read-only.",
		Query:       `query { settings }`,
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "setting_get",
		Description: "Get single setting by name.",
		Query:       `query Setting($name: String!) { setting(name: $name) }`,
		Params: map[string]string{
			"name": "string", // required=true
		},
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "setting_update",
		Description: "Update a setting.",
		Query:       `mutation UpdateSetting($name: String!, $value: String!) { updateSetting(name: $name, value: $value) { status } }`,
		Params: map[string]string{
			"name": "string", // required=true
			"value": "string", // required=true
		},
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "setting_configure_ups",
		Description: "Configure UPS.",
		Query:       `mutation ConfigureUPS($driver: String!, $port: String!) { configureUPS(driver: $driver, port: $port) { status } }`,
		Params: map[string]string{
			"driver": "string", // required=true
			"port": "string", // required=true
		},
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "setting_update_ssh",
		Description: "Update SSH settings.",
		Query:       `mutation UpdateSSH($enabled: Boolean!) { updateSSH(enabled: $enabled) { status } }`,
		Params: map[string]string{
			"enabled": "boolean", // required=true
		},
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "setting_update_system_time",
		Description: "Update system time/NTP.",
		Query:       `mutation UpdateSystemTime($tz: String!) { updateSystemTime(timezone: $tz) { status } }`,
		Params: map[string]string{
			"timezone": "string", // required=true
		},
	})
}