package connect

import (
	"github.com/hektyc/unraid-mcp-server/internal/config"
	"github.com/hektyc/unraid-mcp-server/internal/mcp"
)

func RegisterTools(s *mcp.Server, cfg *config.Config) {
	s.RegisterTool(&mcp.ToolDef{
		Name:        "connect_status",
		Description: "Get connection status. Read-only.",
		Query:       `query { connect { dynamicRemoteAccess settings { values } } }`,
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "connect_sign_in",
		Description: "Sign in to Connect.",
		Query:       `mutation($apiKey: String!) { connectSignIn(input: {apiKey: $apiKey}) }`,
		Params: map[string]string{
			"apiKey": "string",
		},
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "connect_sign_out",
		Description: "Sign out of Connect.",
		Query:       `mutation { connectSignOut }`,
	})
}
