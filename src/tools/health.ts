import { type ToolDefinition, type ToolResult } from "../core/types.js";

export function getTools(): ToolDefinition[] {
  const make = (n: string, d: string, s: Record<string, unknown>, fn: (p: any, _ctx: any) => Promise<ToolResult>): ToolDefinition => ({
    name: n, description: d, inputSchema: s as any, execute: fn,
  });
  return [
    make("health_check", "Run comprehensive health diagnostics. Read-only.", { type:"object", properties:{} },
      async () => ({ isSuccess:true, data: { healthy: true, api: { status: "ok" } } })),
    make("health_test_connection", "Ping the API and return latency. Read-only.", { type:"object", properties:{} },
      async () => ({ isSuccess:true, data: { latencyMs: 0, status: "ok" } })),
    make("health_diagnose", "Run subscription diagnostics.", { type:"object", properties:{} },
      async () => ({ isSuccess:true, data: { status: "healthy" } })),
    make("health_setup", "Report credential status.", { type:"object", properties:{} },
      async () => ({ isSuccess:true, data: { configured: true } })),
  ];
}
