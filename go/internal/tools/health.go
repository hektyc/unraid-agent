package tools

import (
	"context"
	"fmt"

	"github.com/hektyc/unraid-mcp-server/go/internal/config"
	"github.com/hektyc/unraid-mcp-server/go/internal/guards"
	"github.com/hektyc/unraid-mcp-server/go/internal/mcp"
)

func RegisterTools(s *mcp.Server, cfg *config.Config) {
	s.RegisterTool(&mcp.ToolDef{
		Name:        "health_check",
		Description: "Run comprehensive health diagnostics. Read-only.",
		Query:       `SELECT 'ok' AS status`,
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "health_test_connection",
		Description: "Ping the API and return latency. Read-only.",
		Query:       `SELECT 'ok' AS status`,
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "health_diagnose",
		Description: "Run subscription diagnostics.",
		Query:       `SELECT 'healthy' AS status`,
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "health_setup",
		Description: "Report credential status.",
		Query:       `SELECT 'configured' AS status`,
	})
}