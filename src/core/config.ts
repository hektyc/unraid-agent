import "dotenv/config";
import dotenvExpand from "dotenv-expand";

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

export function loadConfig(): UnraidConfig {
  dotenvExpand({});

  const readOnly = envBool("READ_ONLY", false);

  const safetyToggles: Record<string, boolean> = {
    allowArrayStop: envBool("ALLOW_ARRAY_STOP", false),
    allowArrayStart: envBool("ALLOW_ARRAY_START", false),
    allowArrayAddDisk: envBool("ALLOW_ARRAY_ADD_DISK", false),
    allowArrayRemoveDisk: envBool("ALLOW_ARRAY_REMOVE_DISK", false),
    allowArrayClearStats: envBool("ALLOW_ARRAY_CLEAR_STATS", false),
    allowContainerStop: envBool("ALLOW_CONTAINER_STOP", false),
    allowContainerRemove: envBool("ALLOW_CONTAINER_REMOVE", false),
    allowContainerRestart: envBool("ALLOW_CONTAINER_RESTART", false),
    allowVmForceStop: envBool("ALLOW_VM_FORCE_STOP", false),
    allowVmReset: envBool("ALLOW_VM_RESET", false),
    allowVmStop: envBool("ALLOW_VM_STOP", false),
    allowPluginInstall: envBool("ALLOW_PLUGIN_INSTALL", false),
    allowPluginRemove: envBool("ALLOW_PLUGIN_REMOVE", false),
    allowSettingUpdates: envBool("ALLOW_SETTING_UPDATES", false),
    allowSshUpdate: envBool("ALLOW_SSH_UPDATE", false),
    allowTimeUpdate: envBool("ALLOW_TIME_UPDATE", false),
    allowNotificationDelete: envBool("ALLOW_NOTIFICATION_DELETE", false),
    allowApiKeyCreate: envBool("ALLOW_API_KEY_CREATE", false),
    allowApiKeyDelete: envBool("ALLOW_API_KEY_DELETE", false),
    allowFlashBackup: envBool("ALLOW_FLASH_BACKUP", false),
    allowRcloneOperations: envBool("ALLOW_RCLONE_OPERATIONS", false),
    allowConnectActions: envBool("ALLOW_CONNECT_ACTIONS", false),
    allowOnboardingActions: envBool("ALLOW_ONBOARDING_ACTIONS", false),
    allowDockerActions: envBool("ALLOW_DOCKER_ACTIONS", false),
    allowVmActions: envBool("ALLOW_VM_ACTIONS", false),
    allowArrayActions: envBool("ALLOW_ARRAY_ACTIONS", false),
    allowDestructive: envBool("ALLOW_DESTRUCTIVE", false),
  };

  const rawTransport = process.env.TRANSPORT || "stdio";
  const transport: string = ["stdio", "streamable-http", "sse"].includes(rawTransport) ? rawTransport : "stdio";

  const rawLogLevel = (process.env.UNRAID_MCP_LOG_LEVEL || "info").toLowerCase().trim();
  const logLevel: string = ["debug", "info", "warn", "error", "silent"].includes(rawLogLevel) ? rawLogLevel : "info";

  const verifySsl = envBool("UNRAID_VERIFY_SSL", true);
  const allowInsecureTls = envBool("UNRAID_ALLOW_INSECURE_TLS", false);

  if (!verifySsl && !allowInsecureTls) {
    throw new Error("UNRAID_VERIFY_SSL=false requires UNRAID_ALLOW_INSECURE_TLS=true");
  }

  if (envBool("UNRAID_MCP_DISABLE_HTTP_AUTH", false)) {
    const bindHost = process.env.UNRAID_MCP_HOST || "127.0.0.1";
    if (bindHost !== "127.0.0.1" && !envBool("UNRAID_MCP_TRUST_PROXY", false)) {
      throw new Error("UNRAID_MCP_DISABLE_HTTP_AUTH=true requires UNRAID_MCP_TRUST_PROXY=true");
    }
  }

  const apiUrl = process.env.UNRAID_API_URL;
  const apiKey = process.env.UNRAID_API_KEY;

  if (!apiUrl || !apiKey || apiUrl.trim() === "" || apiKey.trim() === "") {
    throw new Error("Missing required Unraid MCP config: UNRAID_API_URL, UNRAID_API_KEY");
  }

  return {
    apiUrl: apiUrl.trim(),
    apiKey: apiKey.trim(),
    transport,
    host: process.env.UNRAID_MCP_HOST || "127.0.0.1",
    port: envInt("UNRAID_MCP_PORT", 6970),
    maxResponseBytes: envInt("UNRAID_MCP_MAX_RESPONSE_BYTES", 40000),
    bearerToken: process.env.UNRAID_MCP_BEARER_TOKEN?.trim() || null,
    disableHttpAuth: envBool("UNRAID_MCP_DISABLE_HTTP_AUTH", false),
    trustProxy: envBool("UNRAID_MCP_TRUST_PROXY", false),
    verifySsl,
    allowInsecureTls,
    logLevel,
    readOnly,
    safetyToggles,
  };
}

function envBool(name: string, fallback = false): boolean {
  const raw = process.env[name];
  if (raw === null || raw === undefined || raw === "") return fallback;
  return !["0", "false", "no", "off"].includes(raw.toLowerCase().trim());
}

function envInt(name: string, fallback: number): number {
  const raw = process.env[name];
  if (raw === null || raw === undefined || raw === "") return fallback;
  const parsed = Number.parseInt(raw, 10);
  return Number.isFinite(parsed) ? parsed : fallback;
}
