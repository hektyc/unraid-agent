import { type ToolDefinition, type ToolResult } from "../core/types.js";

export function getTools(): ToolDefinition[] {
  const make = (n: string, d: string, s: Record<string, unknown>, fn: (p: any, ctx: any) => Promise<ToolResult>): ToolDefinition => ({
    name: n, description: d, inputSchema: s as any, execute: fn,
  });
  return [
    make("key_list", "Get all API keys. Read-only.", { type:"object", properties:{} },
      async (_p:any, ctx:any) => ({ isSuccess:true, data: await ctx.query("query { key { list { id name roles permissions } } }") })),
    make("key_get", "Get single API key details.", { type:"object", properties:{ key_id:{type:"string"} }, required:["key_id"] },
      async (p:any, ctx:any) => ({ isSuccess:true, data: await ctx.query("query KeyGet($id: ID!) { key { get(id: $id) { id name roles permissions } } }", p) })),
    make("key_possible_roles", "Get assignable roles. Read-only.", { type:"object", properties:{} },
      async (_p:any, ctx:any) => ({ isSuccess:true, data: await ctx.query("query { key { possibleRoles } }") })),
    make("key_possible_permissions", "Get grantable permissions. Read-only.", { type:"object", properties:{} },
      async (_p:any, ctx:any) => ({ isSuccess:true, data: await ctx.query("query { key { possiblePermissions } }") })),
    make("key_preview_permissions", "Get effective permissions.", { type:"object", properties:{ roles:{type:"array"}, permissions_input:{type:"array"} }, required:["roles"] },
      async (p:any, ctx:any) => ({ isSuccess:true, data: await ctx.query("query PreviewPermissions($roles: [String!], $perms: [PermissionInput!]) { key { previewPermissions(roles: $roles, permissions: $perms) { resource actions } } }", p) })),
    make("key_auth_actions", "Get available auth actions. Read-only.", { type:"object", properties:{} },
      async (_p:any, ctx:any) => ({ isSuccess:true, data: await ctx.query("query { key { authActions } }") })),
    make("key_creation_form_schema", "Get key creation schema. Read-only.", { type:"object", properties:{} },
      async (_p:any, ctx:any) => ({ isSuccess:true, data: await ctx.query("query { key { creationFormSchema } }") })),
    make("key_create", "Create an API key.", { type:"object", properties:{ name:{type:"string"}, roles:{type:"array"}, permissions:{type:"array"} }, required:["name"] },
      async (p:any, ctx:any) => ({ isSuccess:true, data: await ctx.query("mutation CreateKey($name: String!, $roles: [String!], $permissions: [PermissionInput!]) { key { create(name: $name, roles: $roles, permissions: $permissions) { id name token } } }", p) })),
    make("key_update", "Update an API key.", { type:"object", properties:{ key_id:{type:"string"}, name:{type:"string"}, roles:{type:"array"}, permissions:{type:"array"} }, required:["key_id"] },
      async (p:any, ctx:any) => ({ isSuccess:true, data: await ctx.query("mutation UpdateKey($id: ID!, $name: String, $roles: [String!], $permissions: [PermissionInput!]) { key { update(id: $id, name: $name, roles: $roles, permissions: $permissions) { id name } } }", p) })),
    make("key_delete", "Delete an API key. Destructive.", { type:"object", properties:{ key_id:{type:"string"}, confirm:{type:"boolean"} }, required:["key_id","confirm"] },
      async (p:any, ctx:any) => ({ isSuccess:true, data: await ctx.query("mutation DeleteKey($id: ID!) { key { delete(id: $id) { status } } }", p) })),
    make("key_add_role", "Add role to key.", { type:"object", properties:{ key_id:{type:"string"}, roles:{type:"array"} }, required:["key_id","roles"] },
      async (p:any, ctx:any) => ({ isSuccess:true, data: await ctx.query("mutation AddRole($id: ID!, $roles: [String!]!) { key { addRole(id: $id, roles: $roles) { id } } }", p) })),
    make("key_remove_role", "Remove role from key.", { type:"object", properties:{ key_id:{type:"string"}, roles:{type:"array"} }, required:["key_id","roles"] },
      async (p:any, ctx:any) => ({ isSuccess:true, data: await ctx.query("mutation RemoveRole($id: ID!, $roles: [String!]!) { key { removeRole(id: $id, roles: $roles) { id } } }", p) })),
  ];
}
