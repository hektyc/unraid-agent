package user

import (

	"github.com/hektyc/unraid-mcp-server/internal/config"
	"github.com/hektyc/unraid-mcp-server/internal/mcp"
)

func RegisterTools(s *mcp.Server, cfg *config.Config) {
	s.RegisterTool(&mcp.ToolDef{
		Name:        "user_me",
		Description: "Get authenticated user info. Read-only.",
		Query:       `query { me { name email avatar roles } }`,
	})
}