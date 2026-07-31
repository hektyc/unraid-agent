import { type ToolDefinition, type ToolResult } from "../core/types.js";

export function getTools(): ToolDefinition[] {
  const make = (n: string, d: string, s: Record<string, unknown>, fn: (p: any, ctx: any) => Promise<ToolResult>): ToolDefinition => ({
    name: n, description: d, inputSchema: s as any, execute: fn,
  });
  return [
    make("customization_themes", "Get available themes. Read-only.", { type:"object", properties:{} },
      async (_p:any, ctx:any) => ({ isSuccess:true, data: await ctx.query("query { customization { themes { id name author } } }") })),
    make("customization_locales", "Get available locales. Read-only.", { type:"object", properties:{} },
      async (_p:any, ctx:any) => ({ isSuccess:true, data: await ctx.query("query { customization { locales { id name language } } }") })),
    make("customization_set_theme", "Set UI theme.", { type:"object", properties:{ theme_id:{type:"string"} }, required:["theme_id"] },
      async (p:any, ctx:any) => ({ isSuccess:true, data: await ctx.query("mutation SetTheme($themeId: String!) { setTheme(themeId: $themeId) { status } }", p) })),
    make("customization_set_locale", "Set locale.", { type:"object", properties:{ locale_id:{type:"string"} }, required:["locale_id"] },
      async (p:any, ctx:any) => ({ isSuccess:true, data: await ctx.query("mutation SetLocale($localeId: String!) { setLocale(localeId: $localeId) { status } }", p) })),
    make("customization_sso", "Get SSO config. Read-only.", { type:"object", properties:{} },
      async (_p:any, ctx:any) => ({ isSuccess:true, data: await ctx.query("query { customization { sso { enabled provider } } }") })),
    make("customization_display_details", "Get display details. Read-only.", { type:"object", properties:{} },
      async (_p:any, ctx:any) => ({ isSuccess:true, data: await ctx.query("query { display { case theme temperatureDisplay thresholds locale } }") })),
  ];
}
