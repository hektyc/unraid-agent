import { describe, it, expect, vi } from "vitest";
import { loadConfig, isReadOnly, isActionAllowed, isDestructiveAction } from "../src/core/config.js";

vi.mock("dotenv/config", () => ({}));
vi.mock("dotenv-expand", () => ({ expand: () => {} }));

describe("config", () => {
  const originalEnv = process.env;

  beforeEach(() => {
    process.env = { ...originalEnv };
  });

  it("throws if required env vars are missing", () => {
    delete process.env.UNRAID_API_URL;
    delete process.env.UNRAID_API_KEY;
    expect(() => loadConfig()).toThrow("Missing required Unraid MCP config");
  });

  it("loads config from env vars", () => {
    process.env.UNRAID_API_URL = "https://tower.local/graphql";
    process.env.UNRAID_API_KEY = "test-key-123";
    process.env.TRANSPORT = "stdio";
    process.env.READ_ONLY = "true";

    const config = loadConfig();
    expect(config.apiUrl).toBe("https://tower.local/graphql");
    expect(config.apiKey).toBe("test-key-123");
    expect(config.transport).toBe("stdio");
    expect(config.readOnly).toBe(true);
  });

  it("defaults transport to stdio", () => {
    process.env.UNRAID_API_URL = "https://tower.local/graphql";
    process.env.UNRAID_API_KEY = "test-key";
    delete process.env.TRANSPORT;

    const config = loadConfig();
    expect(config.transport).toBe("stdio");
  });

  it("defaults port to 6970", () => {
    process.env.UNRAID_API_URL = "https://tower.local/graphql";
    process.env.UNRAID_API_KEY = "test-key";
    delete process.env.UNRAID_MCP_PORT;

    const config = loadConfig();
    expect(config.port).toBe(6970);
  });

  it("parses boolean env vars correctly", () => {
    process.env.UNRAID_API_URL = "https://tower.local/graphql";
    process.env.UNRAID_API_KEY = "test-key";

    process.env.ALLOW_ARRAY_STOP = "true";
    process.env.ALLOW_VM_FORCE_STOP = "false";
    process.env.ALLOW_DOCKER_ACTIONS = "1";

    const config = loadConfig();
    expect(config.safetyToggles.allowArrayStop).toBe(true);
    expect(config.safetyToggles.allowVmForceStop).toBe(false);
    expect(config.safetyToggles.allowDockerActions).toBe(true);
  });

  it("blocks non-loopback bind when auth disabled without trust proxy", () => {
    process.env.UNRAID_API_URL = "https://tower.local/graphql";
    process.env.UNRAID_API_KEY = "test-key";
    process.env.UNRAID_MCP_DISABLE_HTTP_AUTH = "true";
    process.env.UNRAID_MCP_HOST = "0.0.0.0";
    process.env.UNRAID_MCP_TRUST_PROXY = "false";

    expect(() => loadConfig()).toThrow("UNRAID_MCP_TRUST_PROXY=true");
  });

  it("blocks skip-ssl without explicit opt-in", () => {
    process.env.UNRAID_API_URL = "https://tower.local/graphql";
    process.env.UNRAID_API_KEY = "test-key";
    process.env.UNRAID_VERIFY_SSL = "false";
    delete process.env.UNRAID_ALLOW_INSECURE_TLS;

    expect(() => loadConfig()).toThrow("UNRAID_ALLOW_INSECURE_TLS=true");
  });
});

describe("guards", () => {
  it("returns read-only status correctly", () => {
    expect(isReadOnly({ readOnly: true } as any)).toBe(true);
    expect(isReadOnly({ readOnly: false } as any)).toBe(false);
  });

  it("detects destructive actions correctly", () => {
    expect(isDestructiveAction("array", "stop")).toBe(true);
    expect(isDestructiveAction("array", "parity_status")).toBe(false);
    expect(isDestructiveAction("docker", "remove_container")).toBe(true);
    expect(isDestructiveAction("vm", "start")).toBe(false);
  });

  it("allows safe actions when read-only is true", () => {
    const config = { readOnly: true, safetyToggles: {} } as any;
    expect(isActionAllowed(config, "system", "overview")).toBe(true);
    expect(isActionAllowed(config, "system", "metrics")).toBe(true);
  });

  it("blocks destructive actions when toggle is disabled", () => {
    const config = {
      readOnly: true,
      safetyToggles: {
        allowArrayStop: false,
        allowArrayActions: false,
        allowDestructive: false,
      },
    } as any;
    expect(isActionAllowed(config, "array", "stop")).toBe(false);
    expect(isActionAllowed(config, "array", "start")).toBe(false);
  });

  it("allows destructive actions when toggle is enabled", () => {
    const config = {
      readOnly: true,
      safetyToggles: {
        allowArrayStop: true,
        allowDestructive: false,
      },
    } as any;
    expect(isActionAllowed(config, "array", "stop")).toBe(true);
  });
});
