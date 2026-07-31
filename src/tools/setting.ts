import { type ToolDefinition, type ToolResult } from "../core/types.js";

export function getTools(): ToolDefinition[] {
  const make = (n: string, d: string, s: Record<string, unknown>, fn: (p: any, ctx: any) => Promise<ToolResult>): ToolDefinition => ({
    name: n, description: d, inputSchema: s as any, execute: fn,
  });
  return [
    make("setting_list", "Get all settings. Read-only.", { type:"object", properties:{} },
      async (_p:any, ctx:any) => ({ isSuccess:true, data: await ctx.query("query { settings }") })),
    make("setting_get", "Get single setting by name.", { type:"object", properties:{ name:{type:"string"} }, required:["name"] },
      async (p:any, ctx:any) => ({ isSuccess:true, data: await ctx.query("query Setting($name: String!) { setting(name: $name) }", p) })),
    make("setting_update", "Update a setting.", { type:"object", properties:{ name:{type:"string"}, value:{type:"string"} }, required:["name","value"] },
      async (p:any, ctx:any) => ({ isSuccess:true, data: await ctx.query("mutation UpdateSetting($name: String!, $value: String!) { updateSetting(name: $name, value: $value) { status } }", p) })),
    make("setting_update_many", "Batch update settings.", { type:"object", properties:{ updates:{type:"array"} }, required:["updates"] },
      async (p:any, ctx:any) => ({ isSuccess:true, data: await ctx.query("mutation UpdateMany($u: [SettingInput!]!) { updateSettings(updates: $u) { status } }", p) })),
    make("setting_configure_ups", "Configure UPS.", { type:"object", properties:{ driver:{type:"string"}, port:{type:"string"} }, required:["driver","port"] },
      async (p:any, ctx:any) => ({ isSuccess:true, data: await ctx.query("mutation ConfigureUPS($driver: String!, $port: String!) { configureUPS(driver: $driver, port: $port) { status } }", p) })),
    make("setting_update_ssh", "Update SSH settings.", { type:"object", properties:{ enabled:{type:"boolean"}, port:{type:"number"} }, required:["enabled"] },
      async (p:any, ctx:any) => ({ isSuccess:true, data: await ctx.query("mutation UpdateSSH($enabled: Boolean!, $port: Int) { updateSSH(enabled: $enabled, port: $port) { status } }", p) })),
    make("setting_update_system_time", "Update system time/NTP.", { type:"object", properties:{ timezone:{type:"string"}, useNtp:{type:"boolean"}, ntpServers:{type:"array"} }, required:["timezone"] },
      async (p:any, ctx:any) => ({ isSuccess:true, data: await ctx.query("mutation UpdateSystemTime($tz: String!, $ntp: Boolean, $srv: [String!]) { updateSystemTime(timezone: $tz, useNtp: $ntp, ntpServers: $srv) { status } }", p) })),
  ];
}
