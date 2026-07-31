import { type ToolDefinition, type ToolResult } from "../core/types.js";

export function getTools(): ToolDefinition[] {
  const make = (n: string, d: string, s: Record<string, unknown>, fn: (p: any, ctx: any) => Promise<ToolResult>): ToolDefinition => ({
    name: n, description: d, inputSchema: s as any, execute: fn,
  });
  return [
    make("disk_shares", "Get all user shares. Read-only.", { type:"object", properties:{} },
      async (_p:any, ctx:any) => ({ isSuccess:true, data: await ctx.query("query { shares { name size free used allocationMethod luks } }") })),
    make("disk_disks", "Get physical disk list. Read-only.", { type:"object", properties:{} },
      async (_p:any, ctx:any) => ({ isSuccess:true, data: await ctx.query("query { disks { id device name size type status } }") })),
    make("disk_details", "Get single disk details.", { type:"object", properties:{ disk_id:{type:"string"} }, required:["disk_id"] },
      async (p:any, ctx:any) => ({ isSuccess:true, data: await ctx.query("query DiskDetails($id: ID!) { disk(id: $id) { id serial size temperature type status } }", p) })),
    make("disk_log_files", "List available log files. Read-only.", { type:"object", properties:{} },
      async (_p:any, ctx:any) => ({ isSuccess:true, data: await ctx.query("query { logFiles { name path size modified } }") })),
    make("disk_logs", "Read log file content.", { type:"object", properties:{ log_path:{type:"string"}, tail_lines:{type:"number"} }, required:["log_path"] },
      async (p:any, ctx:any) => ({ isSuccess:true, data: await ctx.query("query Logs($path: String!, $tail: Int) { logs(path: $path, tail: $tail) { path content } }", p) })),
    make("disk_flash_backup", "Flash drive backup. Destructive.", { type:"object", properties:{ remote_name:{type:"string"}, source_path:{type:"string"}, destination_path:{type:"string"}, confirm:{type:"boolean"} }, required:["remote_name","source_path","destination_path","confirm"] },
      async (p:any, ctx:any) => ({ isSuccess:true, data: await ctx.query("mutation FlashBackup($r: String!, $s: String!, $d: String!) { flashBackup(remote: $r, sourcePath: $s, destinationPath: $d) { status jobId } }", p) })),
  ];
}
