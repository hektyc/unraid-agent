import { type ToolDefinition, type ToolResult } from "../core/types.js";

export function getTools(): ToolDefinition[] {
  const make = (n: string, d: string, s: Record<string, unknown>, fn: (p: any, ctx: any) => Promise<ToolResult>): ToolDefinition => ({
    name: n, description: d, inputSchema: s as any, execute: fn,
  });
  return [
    make("oidc_providers", "List providers. Read-only.", { type:"object", properties:{} },
      async (_p:any, ctx:any) => ({ isSuccess:true, data: await ctx.query("query { oidc { providers { id name enabled } } }") })),
    make("oidc_get", "Get provider details.", { type:"object", properties:{ provider_id:{type:"string"} }, required:["provider_id"] },
      async (p:any, ctx:any) => ({ isSuccess:true, data: await ctx.query("query OidcProvider($id: ID!) { oidc { provider(id: $id) { id name enabled clientId } } }", p) })),
    make("oidc_create", "Create provider.", { type:"object", properties:{ name:{type:"string"}, issuer:{type:"string"}, clientId:{type:"string"}, clientSecret:{type:"string"} }, required:["name","issuer","clientId"] },
      async (p:any, ctx:any) => ({ isSuccess:true, data: await ctx.query("mutation CreateOidc($name: String!, $issuer: String!, $clientId: String!, $clientSecret: String) { oidc { create(name: $name, issuer: $issuer, clientId: $clientId, clientSecret: $clientSecret) { id } } }", p) })),
    make("oidc_update", "Update provider.", { type:"object", properties:{ provider_id:{type:"string"}, name:{type:"string"}, issuer:{type:"string"}, clientId:{type:"string"}, enabled:{type:"boolean"} }, required:["provider_id"] },
      async (p:any, ctx:any) => ({ isSuccess:true, data: await ctx.query("mutation UpdateOidc($id: ID!, $name: String, $issuer: String, $clientId: String, $enabled: Boolean) { oidc { update(id: $id, name: $name, issuer: $issuer, clientId: $clientId, enabled: $enabled) { id } } }", p) })),
    make("oidc_delete", "Delete provider.", { type:"object", properties:{ provider_id:{type:"string"} }, required:["provider_id"] },
      async (p:any, ctx:any) => ({ isSuccess:true, data: await ctx.query("mutation DeleteOidc($id: ID!) { oidc { delete(id: $id) { status } } }", p) })),
  ];
}
