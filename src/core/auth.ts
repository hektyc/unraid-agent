import { type IncomingMessage, type ServerResponse } from "node:http";
import { randomBytes, timingSafeEqual } from "node:crypto";

export function createAuthMiddleware(options: { bearerToken: string | null; disabled: boolean; trustProxy: boolean }) {
  const { bearerToken, disabled } = options;

  return function authMiddleware(req: IncomingMessage, res: ServerResponse, next: () => void): void {
    const path = req.url || "";

    if (path === "/health" || path === "/ready" || path === "/.well-known/oauth-protected-resource") {
      return next();
    }

    if (disabled) {
      return next();
    }

    const auth = req.headers.authorization;
    const rawToken = bearerToken || "";
    if (!rawToken || !auth || !auth.startsWith("Bearer ")) {
      res.statusCode = 401;
      res.setHeader("Content-Type", "application/json");
      res.end(JSON.stringify({ error: "Missing or invalid authorization token" }));
      return;
    }

    const provided = auth.slice(7);
    let valid = false;
    try {
      const a = Buffer.from(provided);
      const b = Buffer.from(rawToken);
      if (a.length === b.length && timingSafeEqual(a, b) === 1) valid = true;
    } catch {
      valid = false;
    }

    if (!valid) {
      res.statusCode = 401;
      res.setHeader("Content-Type", "application/json");
      res.end(JSON.stringify({ error: "Missing or invalid authorization token" }));
      return;
    }

    return next();
  };
}

export function ensureBearerToken(token: string | null, _credentialsDir?: string): string {
  if (token && token.trim().length > 0) {
    return token.trim();
  }
  return randomBytes(32).toString("hex");
}
