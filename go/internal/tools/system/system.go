package system

import (
	"github.com/hektyc/unraid-mcp-server/internal/config"
	"github.com/hektyc/unraid-mcp-server/internal/mcp"
)

func RegisterTools(s *mcp.Server, cfg *config.Config) {
	s.RegisterTool(&mcp.ToolDef{
		Name:        "system_overview",
		Description: "Get OS, CPU, memory, versions, machine ID. Read-only.",
		Query:       `query { info { os { platform distro release kernel arch hostname uptime } cpu { brand cores threads } machineId system { manufacturer model } } metrics { memory { total free used percentTotal } cpu { percentTotal } } }`,
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "system_array",
		Description: "Get array state, capacity, disk health. Read-only.",
		Query:       `query { array { state capacity { kilobytes { total used free } disks { total used free } } parities { id name device size status temp } disks { id name device size type status temp fsSize fsFree fsUsed isSpinning } caches { id name device size status fsType temp } } }`,
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "system_network",
		Description: "Get access URLs, ports, LAN/WAN IPs. Read-only.",
		Query:       `query { network { accessUrls { type name ipv4 ipv6 } } }`,
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "system_network_interfaces",
		Description: "Get network interfaces. Read-only.",
		Query:       `query { networkInterfaces { name description macAddress ipAddress netmask gateway useDhcp operstate speed duplex mtu type status } }`,
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "system_network_metrics",
		Description: "Get network throughput. Read-only.",
		Query:       `query { metrics { network { name operstate bytesReceived bytesSent rxSec txSec utilizationPercent } } }`,
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "system_registration",
		Description: "Get license info. Read-only.",
		Query:       `query { registration { type state expiration updateExpiration keyFile { location } } }`,
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "system_flash",
		Description: "Get flash drive info. Read-only.",
		Query:       `query { flash { guid vendor product } }`,
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "system_metrics",
		Description: "Get CPU and memory usage. Read-only.",
		Query:       `query { metrics { cpu { percentTotal } memory { total used free available percentTotal swapTotal swapUsed percentSwapTotal } } }`,
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "system_server",
		Description: "Get hostname, uptime, version. Read-only.",
		Query:       `query { server { name status lanip localurl } info { os { hostname uptime distro release } } vars { version } array { state } }`,
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "system_server_time",
		Description: "Get server time and NTP. Read-only.",
		Query:       `query { systemTime { currentTime timeZone useNtp ntpServers } }`,
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "system_services",
		Description: "Get running services. Read-only.",
		Query:       `query { services { name online version } }`,
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "system_settings",
		Description: "Get unified settings. Read-only.",
		Query:       `query { settings { unified { values } } }`,
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "system_variables",
		Description: "Get full Unraid variables. Read-only.",
		Query:       `query { vars { version name timeZone comment security workgroup domain sysModel useSsl port portssl localTld useSsh portssh useNtp startArray configValid mdState mdNumDisks deviceCount flashGuid flashProduct flashVendor regTy regState regTo sbState shareCount shareSmbEnabled shareNfsEnabled } }`,
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "system_timezones",
		Description: "Get available timezones. Read-only.",
		Query:       `query { timeZoneOptions { value label } }`,
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "system_online",
		Description: "Check server reachability. Read-only.",
		Query:       `query { online }`,
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "system_config",
		Description: "Get config validity. Read-only.",
		Query:       `query { config { valid error } }`,
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "system_display",
		Description: "Get UI theme. Read-only.",
		Query:       `query { display { theme locale unit } }`,
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "system_ups_config",
		Description: "Get UPS config. Read-only.",
		Query:       `query { upsConfiguration { service upsCable upsType device overrideUpsCapacity batteryLevel minutes timeout killUps nisIp netServer upsName modelName } }`,
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "system_ups_devices",
		Description: "Get UPS devices. Read-only.",
		Query:       `query { upsDevices { id name model status battery { chargeLevel estimatedRuntime health } power { inputVoltage outputVoltage loadPercentage nominalPower currentPower } } }`,
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "system_owner",
		Description: "Get owner info. Read-only.",
		Query:       `query { owner { username url avatar } }`,
	})
	s.RegisterTool(&mcp.ToolDef{
		Name:        "system_servers",
		Description: "Get all registered servers. Read-only.",
		Query:       `query { servers { id name comment status wanip lanip localurl remoteurl } }`,
	})
}
