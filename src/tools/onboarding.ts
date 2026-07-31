import { type ToolDefinition, type ToolResult } from "../core/types.js";

export function getTools(): ToolDefinition[] {
  const make = (n: string, d: string, s: Record<string, unknown>, fn: (p: any, ctx: any) => Promise<ToolResult>): ToolDefinition => ({
    name: n, description: d, inputSchema: s as any, execute: fn,
  });
  return [
    make("onboarding_state", "Get first-boot state. Read-only.", { type:"object", properties:{} },
      async (_p:any, ctx:any) => ({ isSuccess:true, data: await ctx.query("query { onboarding { state } }") })),
    make("onboarding_complete", "Complete onboarding.", { type:"object", properties:{} },
      async (_p:any, ctx:any) => ({ isSuccess:true, data: await ctx.query("mutation { completeOnboarding { status } }") })),
    make("onboarding_open", "Open onboarding.", { type:"object", properties:{} },
      async (_p:any, ctx:any) => ({ isSuccess:true, data: await ctx.query("mutation { openOnboarding { status } }") })),
    make("onboarding_close", "Close onboarding.", { type:"object", properties:{} },
      async (_p:any, ctx:any) => ({ isSuccess:true, data: await ctx.query("mutation { closeOnboarding { status } }") })),
    make("onboarding_resume", "Resume onboarding.", { type:"object", properties:{} },
      async (_p:any, ctx:any) => ({ isSuccess:true, data: await ctx.query("mutation { resumeOnboarding { status } }") })),
    make("onboarding_bypass", "Bypass onboarding.", { type:"object", properties:{} },
      async (_p:any, ctx:any) => ({ isSuccess:true, data: await ctx.query("mutation { bypassOnboarding { status } }") })),
    make("onboarding_reset", "Reset onboarding.", { type:"object", properties:{ confirm:{type:"boolean"} }, required:["confirm"] },
      async (p:any, ctx:any) => ({ isSuccess:true, data: await ctx.query("mutation { resetOnboarding { status } }", p) })),
    make("onboarding_create_internal_boot_pool", "Create boot pool.", { type:"object", properties:{} },
      async (_p:any, ctx:any) => ({ isSuccess:true, data: await ctx.query("mutation { createInternalBootPool { status } }") })),
  ];
}
