package tools

import (
	"context"
	"fmt"

	"github.com/hektyc/unraid-mcp-server/internal/config"
	"github.com/hektyc/unraid-mcp-server/internal/guards"
	"github.com/hektyc/unraid-mcp-server/internal/mcp"
)

func RegisterTools(s *mcp.Server, cfg *config.Config) {
	s.RegisterTool(&mcp.ToolDef{
		Name:        "onboarding_state",
		Description: "Get first-boot state. Read-only.",
		Query:       `query { onboarding { state } }`,
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "onboarding_complete",
		Description: "Complete onboarding.",
		Query:       `mutation { completeOnboarding { status } }`,
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "onboarding_reset",
		Description: "Reset onboarding.",
		Query:       `mutation { resetOnboarding { status } }`,
		Params: map[string]string{
			"confirm": "boolean", // required=true
		},
	})
}