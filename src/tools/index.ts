import { type ToolDefinition } from "../core/types.js";

const MODULES = [
  "./system.js", "./health.js", "./array.js", "./disk.js", "./docker.js",
  "./vm.js", "./notification.js", "./key.js", "./plugin.js", "./rclone.js",
  "./setting.js", "./connect.js", "./customization.js", "./oidc.js",
  "./onboarding.js", "./user.js", "./live.js", "./help.js",
];

export async function getAllTools(): Promise<ToolDefinition[]> {
  const tools: ToolDefinition[] = [];
  for (const path of MODULES) {
    const mod = await import(path);
    const fn = mod.getTools || mod.getSystemTools || mod.getHealthTools || mod.getArrayTools || mod.getDiskTools || mod.getDockerTools || mod.getVmTools || mod.getNotificationTools || mod.getKeyTools || mod.getPluginTools || mod.getRcloneTools || mod.getSettingTools || mod.getConnectTools || mod.getCustomizationTools || mod.getOidcTools || mod.getOnboardingTools || mod.getUserTools || mod.getLiveTools || mod.getHelpTools;
    if (fn) tools.push(...fn());
  }
  return tools;
}

export async function executeToolDirect(config: any, toolName: string, params: Record<string, unknown>): Promise<any> {
  const tools = await getAllTools();
  const tool = tools.find((t: any) => t.name === toolName);
  if (!tool) return { isSuccess: false, error: "Unknown tool: " + toolName };

  const { createGraphQLClient, graphqlRequestWithRetry } = await import("../core/client.js");
  const { createLogger } = await import("../core/logger.js");
  const client = createGraphQLClient(config);
  const logger = createLogger({ name: "tools", level: config.logLevel });
  const ctx = {
    config,
    logger,
    client,
    async query<T>(q: string, v?: Record<string, unknown>): Promise<T> {
      return graphqlRequestWithRetry<T>(client, q, v);
    },
    async run<T>(fn: () => Promise<any>): Promise<any> {
      return fn();
    },
  };

  try {
    return await tool.execute(params, ctx);
  } catch (error) {
    return { isSuccess: false, error: error instanceof Error ? error.message : String(error) };
  }
}
