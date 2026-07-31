export class GraphQLClientError extends Error {
  constructor(message: string) {
    super(message);
    this.name = "GraphQLClientError";
  }
}

export class ReadOnlyError extends Error {
  constructor(action: string, subaction: string) {
    super(`READ_ONLY mode is enabled. Action '${action}.${subaction}' is blocked.`);
    this.name = "ReadOnlyError";
  }
}

export class ActionNotAllowedError extends Error {
  constructor(action: string, subaction: string, toggle: string) {
    super(`Action '${action}.${subaction}' is blocked. Required toggle '${toggle}' is not enabled.`);
    this.name = "ActionNotAllowedError";
  }
}

export function formatGraphQLError(error: unknown): { isSuccess: boolean; error: string } {
  if (error instanceof Error) return { isSuccess: false, error: error.message };
  return { isSuccess: false, error: String(error) };
}

export function getDestructiveConfirmationText(action: string, subaction: string): string {
  return `[SAFETY] Action: ${action}.${subaction} requires explicit confirmation. Set the appropriate ALLOW_* toggle in your environment.`;
}
