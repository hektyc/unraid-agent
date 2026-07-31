import { describe, it, expect, vi } from "vitest";
import { createGraphQLClient, graphqlRequest, graphqlRequestWithRetry } from "../../src/core/client.js";
import { createLogger } from "../../src/core/logger.js";

describe("client", () => {
  it("creates a GraphQL client with correct headers", () => {
    const client = createGraphQLClient({
      apiUrl: "https://tower.local/graphql",
      apiKey: "test-key",
      transport: "stdio",
      host: "127.0.0.1",
      port: 6970,
      maxResponseBytes: 40000,
      bearerToken: null,
      disableHttpAuth: false,
      trustProxy: false,
      verifySsl: true,
      allowInsecureTls: false,
      logLevel: "info",
      readOnly: false,
      safetyToggles: {} as any,
    });

    expect(client).toBeDefined();
  });

  it("graphqlRequestWithRetry handles eventual success", async () => {
    const logger = createLogger({ name: "test", level: "silent" });
    const mockClient = {
      request: vi.fn()
        .mockRejectedValueOnce(new Error("transient error"))
        .mockResolvedValueOnce({ data: "success" }),
    } as any;

    const result = await graphqlRequestWithRetry(mockClient, "query { test }", undefined, 3);
    expect(result).toEqual({ data: "success" });
    expect(mockClient.request).toHaveBeenCalledTimes(2);
  });

  it("graphqlRequestWithRetry throws after max retries", async () => {
    const mockClient = {
      request: vi.fn().mockRejectedValue(new Error("persistent error")),
    } as any;

    await expect(
      graphqlRequestWithRetry(mockClient, "query { test }", undefined, 2),
    ).rejects.toThrow("persistent error");
    expect(mockClient.request).toHaveBeenCalledTimes(2);
  });
});
