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
		Name:        "help_full",
		Description: "Get comprehensive tool reference.",
		Query:       `SELECT 'See documentation' AS reference`,
	})
}