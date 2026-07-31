package system

import (
	"context"
	"fmt"

	"github.com/hektyc/unraid-mcp-server/go/internal/config"
	"github.com/hektyc/unraid-mcp-server/go/internal/guards"
	"github.com/hektyc/unraid-mcp-server/go/internal/mcp"
)

func RegisterTools(s *mcp.Server, cfg *config.Config) {
	s.RegisterTool(&mcp.ToolDef{
		Name:        "system_overview",
		Description: "Get OS, CPU, memory, versions, machine ID. Read-only.",
		Query:       `query { info { os { name version arch } cpu { name cores } memory { total free } machineId } }`,
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "system_array",
		Description: "Get array state, capacity, disk health. Read-only.",
		Query:       `query { array { state capacity disks { id name size type status } } }`,
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "system_network",
		Description: "Get access URLs, ports, LAN/WAN IPs. Read-only.",
		Query:       `query { network { accessUrls { type name ipv4 ipv6 } } }`,
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "system_registration",
		Description: "Get license info. Read-only.",
		Query:       `query { registration { type keyFile expiration } }`,
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "system_variables",
		Description: "Get full Unraid variables. Read-only.",
		Query:       `query { variables }`,
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "system_metrics",
		Description: "Get CPU and memory usage. Read-only.",
		Query:       `query { metrics { cpu memory { total used free usedPercent } } }`,
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "system_network_metrics",
		Description: "Get network throughput. Read-only.",
		Query:       `query { metrics { network { rxBytes txBytes rxSec txSec } } }`,
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "system_services",
		Description: "Get running services. Read-only.",
		Query:       `query { services { name online version } }`,
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "system_display",
		Description: "Get UI theme. Read-only.",
		Query:       `query { display { theme } }`,
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "system_config",
		Description: "Get config validity. Read-only.",
		Query:       `query { config { valid validations { type message } } }`,
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "system_online",
		Description: "Check server reachability. Read-only.",
		Query:       `query { online }`,
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "system_owner",
		Description: "Get owner info. Read-only.",
		Query:       `query { owner { name avatar profileUrl } }`,
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "system_settings",
		Description: "Get unified settings. Read-only.",
		Query:       `query { settings }`,
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "system_server",
		Description: "Get hostname, uptime, version. Read-only.",
		Query:       `query { server { hostname uptime version arrayState } }`,
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "system_servers",
		Description: "Get all registered servers. Read-only.",
		Query:       `query { servers { id name ip lan url gui } }`,
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "system_flash",
		Description: "Get flash drive info. Read-only.",
		Query:       `query { flash { device vendor product size } }`,
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "system_ups_devices",
		Description: "Get UPS devices. Read-only.",
		Query:       `query { ups { devices { name batteryCharge status } } }`,
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "system_ups_config",
		Description: "Get UPS config. Read-only.",
		Query:       `query { ups { config { driver port } } }`,
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "system_server_time",
		Description: "Get server time and NTP. Read-only.",
		Query:       `query { serverTime { time timeZone ntp } }`,
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "system_timezones",
		Description: "Get available timezones. Read-only.",
		Query:       `query { timezones }`,
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "system_network_interfaces",
		Description: "Get network interfaces. Read-only.",
		Query:       `query { network { interfaces { name type ipv4 ipv6 mac } } }`,
	})
}