import { type ToolDefinition, type ToolResult } from "../core/types.js";

export function getTools(): ToolDefinition[] {
  const make = (n: string, d: string, s: Record<string, unknown>, fn: (p: any, ctx: any) => Promise<ToolResult>): ToolDefinition => ({
    name: n, description: d, inputSchema: s as any, execute: fn,
  });
  return [
    make("connect_status", "Get connection status. Read-only.", { type:"object", properties:{} },
      async (_p:any, ctx:any) => ({ isSuccess:true, data: await ctx.query("query { connect { status } }") })),
    make("connect_remote_access", "Get remote access info. Read-only.", { type:"object", properties:{} },
      async (_p:any, ctx:any) => ({ isSuccess:true, data: await ctx.query("query { connect { remoteAccess { enabled url } } }") })),
    make("connect_cloud", "Get cloud status. Read-only.", { type:"object", properties:{} },
      async (_p:any, ctx:any) => ({ isSuccess:true, data: await ctx.query("query { connect { cloud { status email } } }") })),
    make("connect_sign_in", "Sign in to Connect.", { type:"object", properties:{ email:{type:"string"}, password:{type:"string"} }, required:["email","password"] },
      async (p:any, ctx:any) => ({ isSuccess:true, data: await ctx.query("mutation SignIn($email: String!, $password: String!) { connect { signIn(email: $email, password: $password) { status } } }", p) })),
    make("connect_sign_out", "Sign out of Connect.", { type:"object", properties:{} },
      async (_p:any, ctx:any) => ({ isSuccess:true, data: await ctx.query("mutation { connect { signOut { status } } }") })),
    make("connect_devices", "Get paired devices. Read-only.", { type:"object", properties:{} },
      async (_p:any, ctx:any) => ({ isSuccess:true, data: await ctx.query("query { connect { devices { id name type status } } }") })),
    make("connect_pair_device", "Pair a device.", { type:"object", properties:{ name:{type:"string"} }, required:["name"] },
      async (p:any, ctx:any) => ({ isSuccess:true, data: await ctx.query("mutation PairDevice($name: String!) { connect { pairDevice(name: $name) { id code } } }", p) })),
    make("connect_remove_device", "Remove a device.", { type:"object", properties:{ device_id:{type:"string"} }, required:["device_id"] },
      async (p:any, ctx:any) => ({ isSuccess:true, data: await ctx.query("mutation RemoveDevice($id: ID!) { connect { removeDevice(id: $id) { status } } }", p) })),
  ];
}
