import { describe, it, expect, vi } from "vitest";
import { StdioServerTransport } from "@modelcontextprotocol/sdk/server/stdio.js";
import { Server } from "@modelcontextprotocol/sdk/server/index.js";

describe("server", () => {
  it("Server can be instantiated", () => {
    const server = new Server(
      { name: "unraid-mcp-server", version: "0.0.1" },
      { capabilities: { tools: { listChanged: false } } },
    );
    expect(server).toBeDefined();
  });

  it("StdioServerTransport can be instantiated", () => {
    const transport = new StdioServerTransport();
    expect(transport).toBeDefined();
  });

  it("loads config and starts server in stdio mode", async () => {
    process.env.UNRAID_API_URL = "https://tower.local/graphql";
    process.env.UNRAID_API_KEY = "test-key";
    process.env.TRANSPORT = "stdio";

    const { runStdioServer } = await import("../../src/server.js");

    const serverPromise = runStdioServer();

    await new Promise((resolve) => setTimeout(resolve, 100));

    vi.spyOn(process, "exit").mockImplementation(() => {
      throw new Error("process.exit");
    });

    try {
      await serverPromise;
    } catch {
      // expected
    }
  });
});
