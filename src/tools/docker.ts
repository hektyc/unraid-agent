import { type ToolDefinition, type ToolResult } from "../core/types.js";

export function getTools(): ToolDefinition[] {
  const make = (n: string, d: string, s: Record<string, unknown>, fn: (p: any, ctx: any) => Promise<ToolResult>): ToolDefinition => ({
    name: n, description: d, inputSchema: s as any, execute: fn,
  });
  return [
    make("docker_list", "Get all containers. Read-only.", { type:"object", properties:{} },
      async (_p:any, ctx:any) => ({ isSuccess:true, data: await ctx.query("query { docker { containers { id names image state status autoStart } } }") })),
    make("docker_details", "Get single container details.", { type:"object", properties:{ container_id:{type:"string"} }, required:["container_id"] },
      async (p:any, ctx:any) => ({ isSuccess:true, data: await ctx.query("query DockerDetails($id: ID!) { docker { container(id: $id) { id names image state status ports mounts } } }", p) })),
    make("docker_logs", "Not available via GraphQL API.", { type:"object", properties:{ container_id:{type:"string"} }, required:["container_id"] },
      async (p:any, _ctx:any) => ({ isSuccess:false, error: "Not available via GraphQL. Use docker logs " + p.container_id + " on host." })),
    make("docker_ports", "Get host port bindings. Read-only.", { type:"object", properties:{} },
      async (_p:any, ctx:any) => ({ isSuccess:true, data: await ctx.query("query { docker { containers { ports { hostPort containerPort protocol } } } }") })),
    make("docker_start", "Start a container.", { type:"object", properties:{ container_id:{type:"string"} }, required:["container_id"] },
      async (p:any, ctx:any) => ({ isSuccess:true, data: await ctx.query("mutation StartContainer($id: ID!) { startDockerContainer(id: $id) { status } }", p) })),
    make("docker_stop", "Stop a container.", { type:"object", properties:{ container_id:{type:"string"} }, required:["container_id"] },
      async (p:any, ctx:any) => ({ isSuccess:true, data: await ctx.query("mutation StopContainer($id: ID!) { stopDockerContainer(id: $id) { status } }", p) })),
    make("docker_restart", "Restart a container.", { type:"object", properties:{ container_id:{type:"string"} }, required:["container_id"] },
      async (p:any, ctx:any) => ({ isSuccess:true, data: await ctx.query("mutation RestartContainer($id: ID!) { restartDockerContainer(id: $id) { status } }", p) })),
    make("docker_unpause", "Unpause a container.", { type:"object", properties:{ container_id:{type:"string"} }, required:["container_id"] },
      async (p:any, ctx:any) => ({ isSuccess:true, data: await ctx.query("mutation UnpauseContainer($id: ID!) { unpauseDockerContainer(id: $id) { status } }", p) })),
    make("docker_remove_container", "Remove a container. Destructive.", { type:"object", properties:{ container_id:{type:"string"}, confirm:{type:"boolean"}, with_image:{type:"boolean"} }, required:["container_id","confirm"] },
      async (p:any, ctx:any) => ({ isSuccess:true, data: await ctx.query("mutation RemoveContainer($id: ID!, $wi: Boolean) { removeDockerContainer(id: $id, withImage: $wi) { status } }", p) })),
    make("docker_update_container", "Apply pending image update.", { type:"object", properties:{ container_id:{type:"string"} }, required:["container_id"] },
      async (p:any, ctx:any) => ({ isSuccess:true, data: await ctx.query("mutation UpdateContainer($id: ID!) { updateDockerContainer(id: $id) { status } }", p) })),
    make("docker_refresh_digests", "Refresh image digests.", { type:"object", properties:{} },
      async (_p:any, ctx:any) => ({ isSuccess:true, data: await ctx.query("mutation { refreshDockerDigests { status } }") })),
    make("docker_networks", "Get all networks. Read-only.", { type:"object", properties:{} },
      async (_p:any, ctx:any) => ({ isSuccess:true, data: await ctx.query("query { docker { networks { id name driver scope } } }") })),
    make("docker_network_details", "Get single network details.", { type:"object", properties:{ network_id:{type:"string"} }, required:["network_id"] },
      async (p:any, ctx:any) => ({ isSuccess:true, data: await ctx.query("query NetworkDetails($id: ID!) { docker { network(id: $id) { id name driver scope } } }", p) })),
    make("docker_create_folder", "Create an organizer folder.", { type:"object", properties:{ name:{type:"string"} }, required:["name"] },
      async (p:any, ctx:any) => ({ isSuccess:true, data: await ctx.query("mutation CreateFolder($name: String!) { createOrganizerFolder(name: $name) { id name } }", p) })),
    make("docker_delete_entries", "Delete organizer entries. Destructive.", { type:"object", properties:{ entry_ids:{type:"array"}, confirm:{type:"boolean"} }, required:["entry_ids","confirm"] },
      async (p:any, ctx:any) => ({ isSuccess:true, data: await ctx.query("mutation DeleteEntries($ids: [ID!]!) { deleteOrganizerEntries(entryIds: $ids) { status } }", p) })),
  ];
}
