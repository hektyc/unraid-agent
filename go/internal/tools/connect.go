package connect

import (
	"context"
	"fmt"

	"github.com/hektyc/unraid-mcp-server/go/internal/config"
	"github.com/hektyc/unraid-mcp-server/go/internal/guards"
	"github.com/hektyc/unraid-mcp-server/go/internal/mcp"
)

func RegisterTools(s *mcp.Server, cfg *config.Config) {
	s.RegisterTool(&mcp.ToolDef{
		Name:        "connect_status",
		Description: "Get connection status. Read-only.",
		Query:       `query { connect { status } }`,
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "connect_sign_in",
		Description: "Sign in to Connect.",
		Query:       `mutation SignIn($email: String!, $password: String!) { connect { signIn(email: $email, password: $password) { status } } }`,
		Params: map[string]string{
			"email": "string", // required=true
			"password": "string", // required=true
		},
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "connect_sign_out",
		Description: "Sign out of Connect.",
		Query:       `mutation { connect { signOut { status } } }`,
	})
}