import { describe, it, expect, vi } from "vitest";
import { createLogger } from "../src/core/logger.js";

describe("logger", () => {
  it("creates a logger without crashing", () => {
    const logger = createLogger({ name: "test", level: "silent" });
    expect(logger).toBeDefined();
    expect(typeof logger.debug).toBe("function");
    expect(typeof logger.info).toBe("function");
    expect(typeof logger.warn).toBe("function");
    expect(typeof logger.error).toBe("function");
  });

  it("redacts sensitive strings", () => {
    const { redactSensitive } = require("../src/core/utils.js");
    expect(redactSensitive("api_key=eyJhbGciOiJIUzI1NiJ9")).toBe("[REDACTED]");
    expect(redactSensitive("token=github_pat_1234567890")).toBe("[REDACTED]");
    expect(redactSensitive("normal text")).toBe("normal text");
  });

  it("truncates strings", () => {
    const { truncateString } = require("../src/core/utils.js");
    expect(truncateString("hello", 10)).toBe("hello");
    expect(truncateString("hello world", 5)).toBe("hello[...truncated 6 chars]");
  });

  it("formats bytes", () => {
    const { formatBytes } = require("../src/core/utils.js");
    expect(formatBytes(0)).toBe("0 bytes");
    expect(formatBytes(1024)).toBe("1.00 KB");
    expect(formatBytes(1048576)).toBe("1.00 MB");
  });

  it("formats uptime", () => {
    const { formatUptime } = require("../src/core/utils.js");
    expect(formatUptime(0)).toBe("<1m");
    expect(formatUptime(3661)).toBe("1h 1m");
    expect(formatUptime(90061)).toBe("1d 1h 1m");
  });

  it("retries with backoff", async () => {
    const { retryWithBackoff } = require("../src/core/utils.js");
    let attempts = 0;
    await retryWithBackoff(async () => {
      attempts++;
      if (attempts < 3) throw new Error("fail");
      return "success";
    }, 5, 10);

    expect(attempts).toBe(3);
    expect(retryWithBackoff(async () => { throw new Error("always fail"); }, 2, 10)).rejects.toThrow();
  });
});
