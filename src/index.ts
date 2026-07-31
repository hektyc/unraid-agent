#!/usr/bin/env node

import { loadConfig } from "./core/config.ts";

async function main() {
  try {
    const config = loadConfig();

    switch (config.transport) {
      case "streamable-http": {
        const { runHttpServer } = await import("./server.ts");
        await runHttpServer();
        break;
      }
      case "sse": {
        const { runSseServer } = await import("./server.ts");
        await runSseServer();
        break;
      }
      default: {
        const { runStdioServer } = await import("./server.ts");
        await runStdioServer();
        break;
      }
    }
  } catch (error) {
    console.error("Failed to start Unraid MCP Server:", error instanceof Error ? error.message : error);
    process.exit(1);
  }
}

process.on("SIGINT", () => process.exit(0));
process.on("SIGTERM", () => process.exit(0));

main();
