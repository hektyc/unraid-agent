package health

import (
	"github.com/hektyc/unraid-mcp-server/internal/config"
	"github.com/hektyc/unraid-mcp-server/internal/mcp"
)

func RegisterTools(s *mcp.Server, cfg *config.Config) {
	s.RegisterTool(&mcp.ToolDef{
		Name:        "health_check",
		Description: "Run comprehensive health diagnostics. Read-only.",
		Query:       `query { online services { name online version } array { state } config { valid error } }`,
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "health_test_connection",
		Description: "Ping the API and return latency. Read-only.",
		Query:       `query { online }`,
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "health_diagnose",
		Description: "Run subscription diagnostics.",
		Query:       `query { registration { type state expiration updateExpiration } connect { dynamicRemoteAccess { enabledType runningType error } } }`,
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "health_setup",
		Description: "Report credential status.",
		Query:       `query { me { name description roles } }`,
	})
}
