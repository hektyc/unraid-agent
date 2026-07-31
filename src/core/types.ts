export interface UnraidConfig {
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

export interface ToolResult<T = unknown> {
  isSuccess: boolean;
  data?: T;
  error?: string;
}

export interface ToolDefinition {
  name: string;
  description: string;
  inputSchema: Record<string, unknown>;
  execute: (params: Record<string, unknown>, context: any) => Promise<ToolResult>;
}

export type SubscriptionState = "idle" | "connecting" | "active" | "error" | "reconnecting";
