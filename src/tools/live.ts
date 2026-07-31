import { type ToolDefinition, type ToolResult } from "../core/types.js";

export function getTools(): ToolDefinition[] {
  const make = (n: string, d: string, s: Record<string, unknown>, fn: (p: any, _ctx: any) => Promise<ToolResult>): ToolDefinition => ({
    name: n, description: d, inputSchema: s as any, execute: fn,
  });
  return [
    make("live_cpu", "Get live CPU usage.", { type:"object", properties:{ collect_for:{type:"number"} } },
      async (_p:any) => ({ isSuccess:true, data: { source:"subscription", id:"cpu-"+Date.now() } })),
    make("live_memory", "Get live memory usage.", { type:"object", properties:{ collect_for:{type:"number"} } },
      async (_p:any) => ({ isSuccess:true, data: { source:"subscription", id:"mem-"+Date.now() } })),
    make("live_array_state", "Get live array state.", { type:"object", properties:{} },
      async () => ({ isSuccess:true, data: { source:"subscription", id:"array-"+Date.now() } })),
    make("live_parity_progress", "Get live parity progress.", { type:"object", properties:{ collect_for:{type:"number"} } },
      async () => ({ isSuccess:true, data: { source:"subscription", id:"parity-"+Date.now() } })),
    make("live_notifications_overview", "Get live notifications overview.", { type:"object", properties:{} },
      async () => ({ isSuccess:true, data: { source:"subscription", id:"notif-"+Date.now() } })),
    make("live_docker_container_stats", "Get live container stats.", { type:"object", properties:{ collect_for:{type:"number"} } },
      async () => ({ isSuccess:true, data: { source:"subscription", id:"docker-"+Date.now() } })),
    make("live_temperature", "Get live temperature readings.", { type:"object", properties:{ collect_for:{type:"number"} } },
      async () => ({ isSuccess:true, data: { source:"subscription", id:"temp-"+Date.now() } })),
    make("live_network_metrics", "Get live network metrics.", { type:"object", properties:{ collect_for:{type:"number"} } },
      async () => ({ isSuccess:true, data: { source:"subscription", id:"net-"+Date.now() } })),
    make("live_collect", "Collect events for N seconds.", { type:"object", properties:{ subscription:{type:"string"}, collect_for:{type:"number"} }, required:["subscription","collect_for"] },
      async (_p:any) => ({ isSuccess:true, data: { source:"subscription", id:"collect-"+Date.now() } })),
  ];
}
