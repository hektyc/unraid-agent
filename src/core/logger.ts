import pino from "pino";

export function createLogger(options: { level?: string; name?: string } = {}) {
  const level = options.level || "info";
  const name = options.name || "unraid-mcp";

  const logger = pino({
    level,
    name,
    base: { service: name },
    redact: {
      paths: ["*.api_key", "*.bearer_token", "*.password", "*.secret", "*.token", "*authorization"],
      remove: true,
    },
  });

  return {
    debug: (obj: Record<string, unknown>, msg?: string) => logger.debug(obj, msg),
    info: (obj: Record<string, unknown>, msg?: string) => logger.info(obj, msg),
    warn: (obj: Record<string, unknown>, msg?: string) => logger.warn(obj, msg),
    error: (obj: Record<string, unknown>, msg?: string) => logger.error(obj, msg),
    child: (name: string) => createLogger({ level, name }),
  };
}
