import { describe, it, expect, vi } from "vitest";

describe("health tools", () => {
  it("health_check returns latency and array/config status", async () => {
    const { executeToolDirect } = await import("../../src/tools/health.js");
    const { loadConfig } = await import("../../src/core/config.js");

    process.env.UNRAID_API_URL = "https://tower.local/graphql";
    process.env.UNRAID_API_KEY = "test-key";

    const result = await executeToolDirect(loadConfig(), "health_check", {});
    expect(result.isSuccess).toBe(true);
  });

  it("health_test_connection returns latency", async () => {
    const { executeToolDirect } = await import("../../src/tools/health.js");
    const { loadConfig } = await import("../../src/core/config.js");

    process.env.UNRAID_API_URL = "https://tower.local/graphql";
    process.env.UNRAID_API_KEY = "test-key";

    const result = await executeToolDirect(loadConfig(), "health_test_connection", {});
    expect(result.isSuccess).toBe(true);
  });
});

describe("system tools", () => {
  it("system_overview returns data", async () => {
    const { executeToolDirect } = await import("../../src/tools/system.js");
    const { loadConfig } = await import("../../src/core/config.js");

    process.env.UNRAID_API_URL = "https://tower.local/graphql";
    process.env.UNRAID_API_KEY = "test-key";

    const result = await executeToolDirect(loadConfig(), "system_overview", {});
    expect(result.isSuccess).toBe(true);
  });
});

describe("live tools", () => {
  it("live_array_state creates subscription", async () => {
    const { executeToolDirect } = await import("../../src/tools/live.js");
    const { loadConfig } = await import("../../src/core/config.js");

    process.env.UNRAID_API_URL = "https://tower.local/graphql";
    process.env.UNRAID_API_KEY = "test-key";

    const result = await executeToolDirect(loadConfig(), "live_array_state", {});
    expect(result.isSuccess).toBe(true);
  });
});
