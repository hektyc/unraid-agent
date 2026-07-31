import { type ToolDefinition, type ToolResult } from "../core/types.js";

export function getSystemTools(): ToolDefinition[] {
  const make = (n: string, d: string, s: Record<string, unknown>, fn: (p: Record<string, unknown>, ctx: any) => Promise<ToolResult>): ToolDefinition => ({
    name: n, description: d, inputSchema: s as any, execute: fn,
  });

  return [
    make("system_overview", "Get OS, CPU, memory, versions, machine ID. Read-only.", { type: "object", properties: {} },
      async (_p: any, ctx: any) => ({ isSuccess: true, data: await ctx.query(`query { info { os { name version arch } cpu { name cores } memory { total free } machineId } }`) })),
    make("system_array", "Get array state, capacity, disk health. Read-only.", { type: "object", properties: {} },
      async (_p: any, ctx: any) => ({ isSuccess: true, data: await ctx.query(`query { array { state capacity disks { id name size type status } } }`) })),
    make("system_network", "Get access URLs, ports, LAN/WAN IPs. Read-only.", { type: "object", properties: {} },
      async (_p: any, ctx: any) => ({ isSuccess: true, data: await ctx.query(`query { network { accessUrls { type name ipv4 ipv6 } } }`) })),
    make("system_registration", "Get license info. Read-only.", { type: "object", properties: {} },
      async (_p: any, ctx: any) => ({ isSuccess: true, data: await ctx.query(`query { registration { type keyFile expiration } }`) })),
    make("system_variables", "Get full Unraid variables. Read-only.", { type: "object", properties: {} },
      async (_p: any, ctx: any) => ({ isSuccess: true, data: await ctx.query(`query { variables }`) })),
    make("system_metrics", "Get CPU and memory usage. Read-only.", { type: "object", properties: {} },
      async (_p: any, ctx: any) => ({ isSuccess: true, data: await ctx.query(`query { metrics { cpu memory { total used free usedPercent } } }`) })),
    make("system_network_metrics", "Get network throughput. Read-only.", { type: "object", properties: {} },
      async (_p: any, ctx: any) => ({ isSuccess: true, data: await ctx.query(`query { metrics { network { rxBytes txBytes rxSec txSec } } }`) })),
    make("system_services", "Get running services. Read-only.", { type: "object", properties: {} },
      async (_p: any, ctx: any) => ({ isSuccess: true, data: await ctx.query(`query { services { name online version } }`) })),
    make("system_display", "Get UI theme. Read-only.", { type: "object", properties: {} },
      async (_p: any, ctx: any) => ({ isSuccess: true, data: await ctx.query(`query { display { theme } }`) })),
    make("system_config", "Get config validity. Read-only.", { type: "object", properties: {} },
      async (_p: any, ctx: any) => ({ isSuccess: true, data: await ctx.query(`query { config { valid validations { type message } } }`) })),
    make("system_online", "Check server reachability. Returns boolean. Read-only.", { type: "object", properties: {} },
      async (_p: any, ctx: any) => ({ isSuccess: true, data: await ctx.query(`query { online }`) })),
    make("system_owner", "Get owner info. Read-only.", { type: "object", properties: {} },
      async (_p: any, ctx: any) => ({ isSuccess: true, data: await ctx.query(`query { owner { name avatar profileUrl } }`) })),
    make("system_server", "Get hostname, uptime, version. Read-only.", { type: "object", properties: {} },
      async (_p: any, ctx: any) => ({ isSuccess: true, data: await ctx.query(`query { server { hostname uptime version arrayState } }`) })),
    make("system_servers", "Get all registered servers. Read-only.", { type: "object", properties: {} },
      async (_p: any, ctx: any) => ({ isSuccess: true, data: await ctx.query(`query { servers { id name ip lan url gui } }`) })),
    make("system_flash", "Get flash drive info. Read-only.", { type: "object", properties: {} },
      async (_p: any, ctx: any) => ({ isSuccess: true, data: await ctx.query(`query { flash { device vendor product size } }`) })),
    make("system_ups_devices", "Get UPS devices. Read-only.", { type: "object", properties: {} },
      async (_p: any, ctx: any) => ({ isSuccess: true, data: await ctx.query(`query { ups { devices { name batteryCharge status } } }`) })),
    make("system_server_time", "Get server time and NTP. Read-only.", { type: "object", properties: {} },
      async (_p: any, ctx: any) => ({ isSuccess: true, data: await ctx.query(`query { serverTime { time timeZone ntp } }`) })),
  ];
}
