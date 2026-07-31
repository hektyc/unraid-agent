import { GraphQLClient, type RequestDocument } from "graphql-request";
import { type UnraidConfig } from "./types.js";
import { createLogger } from "./logger.js";

const logger = createLogger({ name: "graphql-client" });

export function createGraphQLClient(config: UnraidConfig): GraphQLClient {
  return new GraphQLClient(config.apiUrl, {
    headers: {
      "Content-Type": "application/json",
      "X-API-Key": config.apiKey,
      Accept: "application/json",
    },
  });
}

export async function graphqlRequestWithRetry<T>(
  client: GraphQLClient,
  query: string | RequestDocument,
  variables?: Record<string, unknown>,
  maxRetries = 3,
): Promise<T> {
  let lastError: unknown;
  for (let attempt = 0; attempt < maxRetries; attempt++) {
    try {
      logger.debug({ query: typeof query === "string" ? query.slice(0, 200) : String(query), attempt }, "GraphQL request");
      const data = await client.request<T>(query, variables);
      return data;
    } catch (error) {
      lastError = error;
      if (attempt < maxRetries - 1) {
        const delay = Math.min(1000 * 2 ** attempt, 10000);
        logger.warn({ attempt: attempt + 1, delay }, "Retrying GraphQL request");
        await new Promise((resolve) => setTimeout(resolve, delay));
      }
    }
  }
  throw lastError;
}
