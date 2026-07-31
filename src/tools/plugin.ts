import { type ToolDefinition, type ToolResult } from "../core/types.js";

export function getTools(): ToolDefinition[] {
  const make = (n: string, d: string, s: Record<string, unknown>, fn: (p: any, ctx: any) => Promise<ToolResult>): ToolDefinition => ({
    name: n, description: d, inputSchema: s as any, execute: fn,
  });
  return [
    make("plugin_list", "Get installed plugins. Read-only.", { type:"object", properties:{} },
      async (_p:any, ctx:any) => ({ isSuccess:true, data: await ctx.query("query { plugin { list { name version status } } }") })),
    make("plugin_installed_unraid", "Get raw .plg filenames. Read-only.", { type:"object", properties:{} },
      async (_p:any, ctx:any) => ({ isSuccess:true, data: await ctx.query("query { plugin { installed } }") })),
    make("plugin_install_operations", "Get async install operations. Read-only.", { type:"object", properties:{} },
      async (_p:any, ctx:any) => ({ isSuccess:true, data: await ctx.query("query { plugin { installOperations { id status pluginName } } }") })),
    make("plugin_install_operation", "Get one install operation status.", { type:"object", properties:{ operation_id:{type:"string"} }, required:["operation_id"] },
      async (p:any, ctx:any) => ({ isSuccess:true, data: await ctx.query("query InstallOp($id: ID!) { plugin { installOperation(id: $id) { id status progress } } }", p) })),
    make("plugin_add", "Install plugins by name.", { type:"object", properties:{ names:{type:"array"}, bundled:{type:"boolean"}, restart:{type:"boolean"} }, required:["names"] },
      async (p:any, ctx:any) => ({ isSuccess:true, data: await ctx.query("mutation AddPlugins($names: [String!]!) { plugin { add(names: $names) { status } } }", p) })),
    make("plugin_remove", "Uninstall plugins. Destructive.", { type:"object", properties:{ names:{type:"array"}, confirm:{type:"boolean"} }, required:["names","confirm"] },
      async (p:any, ctx:any) => ({ isSuccess:true, data: await ctx.query("mutation RemovePlugins($names: [String!]!) { plugin { remove(names: $names) { status } } }", p) })),
    make("plugin_install", "Async install a .plg URL. Destructive.", { type:"object", properties:{ url:{type:"string"}, confirm:{type:"boolean"}, plugin_name:{type:"string"}, forced:{type:"boolean"} }, required:["url","confirm"] },
      async (p:any, ctx:any) => ({ isSuccess:true, data: await ctx.query("mutation InstallPlugin($url: String!) { plugin { install(url: $url) { operationId } } }", p) })),
    make("plugin_install_language", "Install a language pack. Destructive.", { type:"object", properties:{ url:{type:"string"}, confirm:{type:"boolean"} }, required:["url","confirm"] },
      async (p:any, ctx:any) => ({ isSuccess:true, data: await ctx.query("mutation InstallLanguage($url: String!) { plugin { installLanguage(url: $url) { operationId } } }", p) })),
  ];
}
