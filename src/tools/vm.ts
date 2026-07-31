import { type ToolDefinition, type ToolResult } from "../core/types.js";

export function getTools(): ToolDefinition[] {
  const make = (n: string, d: string, s: Record<string, unknown>, fn: (p: any, ctx: any) => Promise<ToolResult>): ToolDefinition => ({
    name: n, description: d, inputSchema: s as any, execute: fn,
  });
  return [
    make("vm_list", "Get all VMs. Read-only.", { type:"object", properties:{} },
      async (_p:any, ctx:any) => ({ isSuccess:true, data: await ctx.query("query { vm { list { id name state uuid } } }") })),
    make("vm_details", "Get single VM details.", { type:"object", properties:{ vm_id:{type:"string"} }, required:["vm_id"] },
      async (p:any, ctx:any) => ({ isSuccess:true, data: await ctx.query("query VmDetails($id: ID!) { vm { dom(id: $id) { id name state uuid } } }", p) })),
    make("vm_start", "Start a VM.", { type:"object", properties:{ vm_id:{type:"string"} }, required:["vm_id"] },
      async (p:any, ctx:any) => ({ isSuccess:true, data: await ctx.query("mutation StartVm($id: ID!) { startVm(id: $id) { status } }", p) })),
    make("vm_stop", "Gracefully stop a VM.", { type:"object", properties:{ vm_id:{type:"string"} }, required:["vm_id"] },
      async (p:any, ctx:any) => ({ isSuccess:true, data: await ctx.query("mutation StopVm($id: ID!) { stopVm(id: $id) { status } }", p) })),
    make("vm_pause", "Pause a VM.", { type:"object", properties:{ vm_id:{type:"string"} }, required:["vm_id"] },
      async (p:any, ctx:any) => ({ isSuccess:true, data: await ctx.query("mutation PauseVm($id: ID!) { pauseVm(id: $id) { status } }", p) })),
    make("vm_resume", "Resume a VM.", { type:"object", properties:{ vm_id:{type:"string"} }, required:["vm_id"] },
      async (p:any, ctx:any) => ({ isSuccess:true, data: await ctx.query("mutation ResumeVm($id: ID!) { resumeVm(id: $id) { status } }", p) })),
    make("vm_reboot", "Reboot a VM.", { type:"object", properties:{ vm_id:{type:"string"} }, required:["vm_id"] },
      async (p:any, ctx:any) => ({ isSuccess:true, data: await ctx.query("mutation RebootVm($id: ID!) { rebootVm(id: $id) { status } }", p) })),
    make("vm_force_stop", "Hard power-off a VM.", { type:"object", properties:{ vm_id:{type:"string"}, confirm:{type:"boolean"} }, required:["vm_id","confirm"] },
      async (p:any, ctx:any) => ({ isSuccess:true, data: await ctx.query("mutation ForceStopVm($id: ID!) { forceStopVm(id: $id) { status } }", p) })),
    make("vm_reset", "Hard reset a VM.", { type:"object", properties:{ vm_id:{type:"string"}, confirm:{type:"boolean"} }, required:["vm_id","confirm"] },
      async (p:any, ctx:any) => ({ isSuccess:true, data: await ctx.query("mutation ResetVm($id: ID!) { resetVm(id: $id) { status } }", p) })),
  ];
}
