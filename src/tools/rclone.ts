import { type ToolDefinition, type ToolResult } from "../core/types.js";

export function getTools(): ToolDefinition[] {
  const make = (n: string, d: string, s: Record<string, unknown>, fn: (p: any, ctx: any) => Promise<ToolResult>): ToolDefinition => ({
    name: n, description: d, inputSchema: s as any, execute: fn,
  });
  return [
    make("rclone_list_remotes", "Get configured remotes. Read-only.", { type:"object", properties:{} },
      async (_p:any, ctx:any) => ({ isSuccess:true, data: await ctx.query("query { rclone { listRemotes { name type } } }") })),
    make("rclone_config_form_schema", "Get remote config schema. Read-only.", { type:"object", properties:{} },
      async (_p:any, ctx:any) => ({ isSuccess:true, data: await ctx.query("query { rclone { configFormSchema } }") })),
    make("rclone_create_remote", "Create a remote.", { type:"object", properties:{ name:{type:"string"}, config:{type:"object"} }, required:["name","config"] },
      async (p:any, ctx:any) => ({ isSuccess:true, data: await ctx.query("mutation CreateRcloneRemote($name: String!, $config: JSON!) { rclone { createRemote(name: $name, config: $config) { name type } } }", p) })),
    make("rclone_delete_remote", "Delete a remote. Destructive.", { type:"object", properties:{ name:{type:"string"} }, required:["name"] },
      async (p:any, ctx:any) => ({ isSuccess:true, data: await ctx.query("mutation DeleteRcloneRemote($name: String!) { rclone { deleteRemote(name: $name) { status } } }", p) })),
  ];
}
