import { type ToolDefinition, type ToolResult } from "../core/types.js";

export function getTools(): ToolDefinition[] {
  const make = (n: string, d: string, s: Record<string, unknown>, fn: (p: any, _ctx: any) => Promise<ToolResult>): ToolDefinition => ({
    name: n, description: d, inputSchema: s as any, execute: fn,
  });
  return [
    make("array_parity_status", "Get parity check progress. Read-only.", { type:"object", properties:{} },
      async () => ({ isSuccess:true, data: { status: "ok" } })),
    make("array_parity_history", "Get past parity results. Read-only.", { type:"object", properties:{} },
      async () => ({ isSuccess:true, data: { history: [] } })),
    make("array_assignable_disks", "Get disks not yet in array. Read-only.", { type:"object", properties:{} },
      async () => ({ isSuccess:true, data: { disks: [] } })),
    make("array_parity_start", "Start a parity check.", { type:"object", properties:{ correct:{type:"boolean"} } },
      async (p:any) => ({ isSuccess:true, data: p })),
    make("array_parity_pause", "Pause a parity check.", { type:"object", properties:{} },
      async () => ({ isSuccess:true, data: { status: "paused" } })),
    make("array_parity_resume", "Resume a parity check.", { type:"object", properties:{} },
      async () => ({ isSuccess:true, data: { status: "running" } })),
    make("array_parity_cancel", "Cancel a parity check.", { type:"object", properties:{} },
      async () => ({ isSuccess:true, data: { status: "cancelled" } })),
    make("array_start", "Start the Unraid array.", { type:"object", properties:{ confirm:{type:"boolean"} }, required:["confirm"] },
      async (p:any) => ({ isSuccess:true, data: p })),
    make("array_stop", "Stop the Unraid array.", { type:"object", properties:{ confirm:{type:"boolean"} }, required:["confirm"] },
      async (p:any) => ({ isSuccess:true, data: p })),
    make("array_add_disk", "Add a disk to the array.", { type:"object", properties:{ disk_id:{type:"string"}, slot:{type:"number"} }, required:["disk_id"] },
      async (p:any) => ({ isSuccess:true, data: p })),
    make("array_remove_disk", "Remove a disk from the array.", { type:"object", properties:{ disk_id:{type:"string"}, confirm:{type:"boolean"} }, required:["disk_id","confirm"] },
      async (p:any) => ({ isSuccess:true, data: p })),
    make("array_mount_disk", "Mount an array disk.", { type:"object", properties:{ disk_id:{type:"string"} }, required:["disk_id"] },
      async (p:any) => ({ isSuccess:true, data: p })),
    make("array_unmount_disk", "Unmount an array disk.", { type:"object", properties:{ disk_id:{type:"string"} }, required:["disk_id"] },
      async (p:any) => ({ isSuccess:true, data: p })),
    make("array_clear_disk_stats", "Clear disk I/O stats.", { type:"object", properties:{ disk_id:{type:"string"}, confirm:{type:"boolean"} }, required:["disk_id","confirm"] },
      async (p:any) => ({ isSuccess:true, data: p })),
  ];
}
