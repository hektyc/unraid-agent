import { WebSocket, type RawData } from "ws";

export type SubscriptionEventCallback = (message: SubscriptionMessage) => void;

export type SubscriptionMessage =
  | { type: "connection_ack" }
  | { type: "connection_error"; payload: { message: string } }
  | { type: "data"; payload: { data: unknown } }
  | { type: "error"; payload: { message: string } }
  | { type: "complete" }
  | { type: "ping" }
  | { type: "pong" };

export type TransportProtocol = "graphql-transport-ws" | "graphql-ws";

export interface SubscriptionProtocolOptions {
  config: any;
  query: string;
  variables?: Record<string, unknown>;
  protocol: TransportProtocol;
  onEvent: SubscriptionEventCallback;
  onClose: (code: number, reason: string) => void;
  onError: (error: Error) => void;
}

export class SubscriptionProtocol {
  private ws: WebSocket | null = null;

  constructor(private readonly options: SubscriptionProtocolOptions) {}

  connect(): void {
    try {
      const url = this.buildWebSocketUrl();
      this.ws = new WebSocket(url, this.options.config.verifySsl ? undefined : { rejectUnauthorized: false });

      this.ws.on("open", () => {
        this.ws?.send(JSON.stringify({
          type: "connection_init",
          payload: { "X-API-Key": this.options.config.apiKey },
        }));
        this.ws?.send(JSON.stringify({
          id: this.generateId(),
          type: "start",
          payload: { query: this.options.query, variables: this.options.variables || {} },
        }));
      });

      this.ws.on("message", (data: RawData) => this.handleMessage(data.toString()));
      this.ws.on("close", (code: number, reason: Buffer) => {
        this.options.onClose(code, reason.toString());
      });
      this.ws.on("error", (err) => this.options.onError(new Error(`WebSocket error: ${err.message}`)));
    } catch (error) {
      this.options.onError(error instanceof Error ? error : new Error(String(error)));
    }
  }

  disconnect(): void {
    if (this.ws && this.ws.readyState === WebSocket.OPEN) {
      this.ws.close(1000, "Client disconnecting");
    }
    this.ws = null;
  }

  private buildWebSocketUrl(): string {
    const url = new URL(this.options.config.apiUrl);
    url.protocol = url.protocol === "https:" ? "wss:" : "ws:";
    url.pathname = url.pathname.replace(/\/graphql\/?$/, "") + "/graphql";
    return url.toString();
  }

  private handleMessage(raw: string): void {
    let message: unknown;
    try { message = JSON.parse(raw); } catch { return; }
    if (typeof message !== "object" || message === null || !("type" in message)) return;

    const msg = message as { type: string; payload?: unknown };
    switch (msg.type) {
      case "connection_ack":
        this.options.onEvent({ type: "connection_ack" });
        break;
      case "connection_error": {
        const payload = msg.payload as { message?: string } | undefined;
        this.options.onEvent({ type: "connection_error", payload: { message: payload?.message || "Connection error" } });
        break;
      }
      case "next": {
        const payload = msg.payload as { data?: unknown } | undefined;
        if (payload?.data !== undefined) {
          this.options.onEvent({ type: "data", payload: { data: payload.data } });
        }
        break;
      }
      case "error": {
        const payload = msg.payload as { message?: string } | undefined;
        this.options.onEvent({ type: "error", payload: { message: typeof payload?.message === "string" ? payload.message : "Subscription error" } });
        break;
      }
      case "complete":
        this.options.onEvent({ type: "complete" });
        break;
      default:
        break;
    }
  }

  private generateId(): string {
    return Math.random().toString(36).slice(2, 10);
  }
}
