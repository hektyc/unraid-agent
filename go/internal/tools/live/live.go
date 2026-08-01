package live

import (
	"github.com/hektyc/unraid-mcp-server/internal/config"
	"github.com/hektyc/unraid-mcp-server/internal/mcp"
)

func RegisterTools(s *mcp.Server, cfg *config.Config) {
	s.RegisterTool(&mcp.ToolDef{
		Name:        "live_cpu",
		Description: "Get live CPU usage. Read-only.",
		Query:       `query { metrics { cpu { percentTotal } } }`,
		Params: map[string]string{
			"collect_for": "number",
		},
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "live_memory",
		Description: "Get live memory usage. Read-only.",
		Query:       `query { metrics { memory { total used free available percentTotal swapTotal swapUsed } } }`,
		Params: map[string]string{
			"collect_for": "number",
		},
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "live_temperature",
		Description: "Get live temperature readings. Read-only.",
		Query:       `query { metrics { temperature { summary { average warningCount criticalCount hottest { name current { value } } } sensors { name location current { value } warning critical } } } }`,
		Params: map[string]string{
			"collect_for": "number",
		},
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "live_network_metrics",
		Description: "Get live network metrics. Read-only.",
		Query:       `query { metrics { network { name operstate bytesReceived bytesSent rxSec txSec utilizationPercent } } }`,
		Params: map[string]string{
			"collect_for": "number",
		},
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "live_array_state",
		Description: "Get live array state.",
		Query:       `query { array { state capacity { kilobytes { total used free } disks { total used free } } } }`,
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "live_parity_progress",
		Description: "Get live parity progress.",
		Query:       `query { array { parityCheckStatus { status progress speed errors correcting paused running duration } } }`,
		Params: map[string]string{
			"collect_for": "number",
		},
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "live_notifications_overview",
		Description: "Get live notifications overview.",
		Query:       `query { notifications { overview { unread { info warning alert total } archive { info warning alert total } } } }`,
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "live_docker_container_stats",
		Description: "Get live container stats.",
		Query:       `query { docker { containers { id names state status } } }`,
		Params: map[string]string{
			"collect_for": "number",
		},
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "live_collect",
		Description: "Collect events for N seconds.",
		Query:       `query { metrics { cpu { percentTotal } memory { total used free percentTotal } network { name rxSec txSec utilizationPercent } temperature { summary { average warningCount criticalCount } } } }`,
		Params: map[string]string{
			"subscription": "string",
			"collect_for":  "number",
		},
	})
}
