import { type UnraidConfig, type ToolResult } from "./types.js";
import { createLogger } from "./logger.js";
import { createGraphQLClient, graphqlRequestWithRetry } from "./client.js";

export interface ToolContext {
  config: UnraidConfig;
  logger: any;
  client: any;
  query<T>(q: string, v?: Record<string, unknown>): Promise<T>;
  run<T>(fn: () => Promise<ToolResult<T>>): Promise<ToolResult<T>>;
}

export function makeCtx(config: UnraidConfig): ToolContext {
  const logger = createLogger({ name: "tools", level: config.logLevel as any });
  const client = createGraphQLClient(config);

  async function query<T>(q: string, v?: Record<string, unknown>): Promise<T> {
    return graphqlRequestWithRetry<T>(client, q, v);
  }

  async function run<T>(fn: () => Promise<ToolResult<T>>): Promise<ToolResult<T>> {
    return fn();
  }

  return { config, logger, client, query, run };
}
