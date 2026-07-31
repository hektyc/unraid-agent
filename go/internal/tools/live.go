package live

import (
	"context"
	"fmt"

	"github.com/hektyc/unraid-mcp-server/go/internal/config"
	"github.com/hektyc/unraid-mcp-server/go/internal/guards"
	"github.com/hektyc/unraid-mcp-server/go/internal/mcp"
)

func RegisterTools(s *mcp.Server, cfg *config.Config) {
	s.RegisterTool(&mcp.ToolDef{
		Name:        "live_cpu",
		Description: "Get live CPU usage.",
		Query:       `SELECT 'subscription' AS source`,
		Params: map[string]string{
			"collect_for": "number", // required=false
		},
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "live_memory",
		Description: "Get live memory usage.",
		Query:       `SELECT 'subscription' AS source`,
		Params: map[string]string{
			"collect_for": "number", // required=false
		},
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "live_array_state",
		Description: "Get live array state.",
		Query:       `SELECT 'subscription' AS source`,
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "live_parity_progress",
		Description: "Get live parity progress.",
		Query:       `SELECT 'subscription' AS source`,
		Params: map[string]string{
			"collect_for": "number", // required=false
		},
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "live_notifications_overview",
		Description: "Get live notifications overview.",
		Query:       `SELECT 'subscription' AS source`,
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "live_docker_container_stats",
		Description: "Get live container stats.",
		Query:       `SELECT 'subscription' AS source`,
		Params: map[string]string{
			"collect_for": "number", // required=false
		},
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "live_temperature",
		Description: "Get live temperature readings.",
		Query:       `SELECT 'subscription' AS source`,
		Params: map[string]string{
			"collect_for": "number", // required=false
		},
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "live_network_metrics",
		Description: "Get live network metrics.",
		Query:       `SELECT 'subscription' AS source`,
		Params: map[string]string{
			"collect_for": "number", // required=false
		},
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "live_collect",
		Description: "Collect events for N seconds.",
		Query:       `SELECT 'subscription' AS source`,
		Params: map[string]string{
			"subscription": "string", // required=true
			"collect_for": "number", // required=true
		},
	})
}