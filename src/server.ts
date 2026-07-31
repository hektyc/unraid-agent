import { Server } from "@modelcontextprotocol/sdk/server/index.ts";
import { StdioServerTransport } from "@modelcontextprotocol/sdk/server/stdio.ts";
import {
  CallToolRequestSchema,
  ListToolsRequestSchema,
  Tool,
} from "@modelcontextprotocol/sdk/types.ts";

interface UnraidConfig {
  apiUrl: string;
  apiKey: string;
  transport: string;
  host: string;
  port: number;
  maxResponseBytes: number;
  bearerToken: string | null;
  disableHttpAuth: boolean;
  trustProxy: boolean;
  verifySsl: boolean;
  allowInsecureTls: boolean;
  logLevel: string;
  readOnly: boolean;
  safetyToggles: Record<string, boolean>;
}

interface ToolResult {
  isSuccess: boolean;
  data?: unknown;
  error?: string;
}

type ToolDefinition = {
  name: string;
  description: string;
  inputSchema: Record<string, unknown>;
  execute: (params: Record<string, unknown>, context: any) => Promise<ToolResult>;
};

let config: UnraidConfig | null = null;

async function getAllTools(): Promise<ToolDefinition[]> {
  const modules = [
    "./system.ts", "./health.ts", "./array.ts", "./disk.ts", "./docker.ts",
    "./vm.ts", "./notification.ts", "./key.ts", "./plugin.ts", "./rclone.ts",
    "./setting.ts", "./connect.ts", "./customization.ts", "./oidc.ts",
    "./onboarding.ts", "./user.ts", "./live.ts", "./help.ts",
  ];
  const tools: ToolDefinition[] = [];
  for (const path of modules) {
    const mod = await import(path);
    if (mod.getTools) tools.push(...mod.getTools());
  }
  return tools;
}

async function buildMcpTools(): Promise<Tool[]> {
  return (await getAllTools()).map((t: any) => ({
    name: t.name,
    description: t.description,
    inputSchema: t.inputSchema as Tool["inputSchema"],
  }));
}

async function callToolHandler(name: string, args: Record<string, unknown>): Promise<{
  content: Array<{ type: string; text: string }>;
}> {
  if (!config) throw new Error("Config not initialized");

  const tools = await getAllTools();
  const tool = tools.find((t) => t.name === name);
  if (!tool) return { content: [{ type: "text", text: JSON.stringify({ isSuccess: false, error: "Unknown tool: " + name }) }] };

  const { createGraphQLClient, graphqlRequestWithRetry } = await import("./core/client.ts");
  const { createLogger } = await import("./core/logger.ts");
  const { expand } = await import("dotenv-expand");

  const client = createGraphQLClient(config);
  const logger = createLogger({ name: "tools", level: config.logLevel as any });
  const ctx = {
    config,
    logger,
    client,
    async query<T>(q: string, v?: Record<string, unknown>): Promise<T> {
      return graphqlRequestWithRetry<T>(client, q, v);
    },
    async run<T>(fn: () => Promise<ToolResult<T>>): Promise<ToolResult<T>> {
      return fn();
    },
  };

  try {
    const result = await tool.execute(args, ctx);
    return { content: [{ type: "text", text: JSON.stringify(result, null, 2) }] };
  } catch (error) {
    return { content: [{ type: "text", text: JSON.stringify({ isSuccess: false, error: error instanceof Error ? error.message : String(error) }, null, 2) }] };
  }
}

export async function runStdioServer(): Promise<void> {
  const server = new Server(
    { name: "unraid-mcp-server", version: "0.0.1" },
    { capabilities: { tools: { listChanged: false } } },
  );

  server.setRequestHandler(ListToolsRequestSchema, async () => ({
    tools: await buildMcpTools(),
  }));

  server.setRequestHandler(CallToolRequestSchema, async (request) => {
    const { name, arguments: args } = request.params;
    return callToolHandler(name, args || {});
  });

  const transport = new StdioServerTransport();
  await server.connect(transport);
}

export async function runHttpServer(): Promise<any> {
  const { createServer } = await import("node:http");
  const { WebSocketServer } = await import("ws");
  const { createAuthMiddleware, ensureBearerToken } = await import("./core/auth.ts");

  const server = new Server(
    { name: "unraid-mcp-server", version: "0.0.1" },
    { capabilities: { tools: { listChanged: false } } },
  );

  server.setRequestHandler(ListToolsRequestSchema, async () => ({
    tools: await buildMcpTools(),
  }));

  server.setRequestHandler(CallToolRequestSchema, async (request) => {
    const { name, arguments: args } = request.params;
    return callToolHandler(name, args || {});
  });

  const bearerToken = ensureBearerToken(config!.bearerToken, "");
  const authMiddleware = createAuthMiddleware({
    bearerToken,
    disabled: config!.disableHttpAuth,
    trustProxy: config!.trustProxy,
  });

  const httpServer = createServer(async (req, res) => {
    const url = new URL(req.url || "/", `http://${req.headers.host}`);

    if (url.pathname === "/health") {
      res.statusCode = 200;
      res.setHeader("Content-Type", "application/json");
      res.end(JSON.stringify({ status: "ok" }));
      return;
    }

    if (url.pathname === "/.well-known/oauth-protected-resource") {
      res.statusCode = 200;
      res.setHeader("Content-Type", "application/json");
      res.end(JSON.stringify({ issuer: "unraid-mcp", scopes_supported: ["openid"] }));
      return;
    }

    authMiddleware(req, res, async () => {
      if (url.pathname === "/mcp") {
        try {
          const { StreamableHTTPServerTransport } = await import("@modelcontextprotocol/sdk/server/streamableHttp.ts");
          const transport = new StreamableHTTPServerTransport({
            sessionIdGenerator: undefined,
            enableJsonResponse: true,
          });

          res.setHeader("Content-Type", "application/json, text/event-stream");
          res.setHeader("Cache-Control", "no-cache");
          res.setHeader("Connection", "keep-alive");

          await server.connect(transport);
        } catch (error) {
          if (!res.headersSent) {
            res.statusCode = 500;
            res.setHeader("Content-Type", "application/json");
            res.end(JSON.stringify({ error: error instanceof Error ? error.message : "Internal error" }));
          }
        }
      } else {
        res.statusCode = 404;
        res.setHeader("Content-Type", "application/json");
        res.end(JSON.stringify({ error: "Not found" }));
      }
    });
  });

  new WebSocketServer({ server: httpServer, path: "/ws" });

  return new Promise((resolve: (value: any) => void) => {
    httpServer.listen(config!.port, config!.host, () => {
      console.error(`Unraid MCP Server HTTP running on http://${config!.host}:${config!.port}`);
      console.error(`  MCP endpoint: http://${config!.host}:${config!.port}/mcp`);
      console.error(`  Health: http://${config!.host}:${config!.port}/health`);
      resolve(httpServer);
    });
  });
}

export async function runSseServer(): Promise<any> {
  console.error("WARNING: SSE transport is deprecated.");
  const { createServer } = await import("node:http");

  const httpServer = createServer((_req, res) => {
    res.statusCode = 200;
    res.setHeader("Content-Type", "text/event-stream");
    res.setHeader("Cache-Control", "no-cache");
    res.setHeader("Connection", "keep-alive");
    res.end("event: initialized\ndata: {}\n\n");
  });

  return new Promise((resolve: (value: any) => void) => {
    httpServer.listen(config!.port, config!.host, () => {
      console.error(`Unraid MCP Server SSE running on http://${config!.host}:${config!.port}`);
      resolve(httpServer);
    });
  });
}
