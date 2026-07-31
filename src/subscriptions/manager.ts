import { SubscriptionProtocol } from "./protocol.js";

export type { SubscriptionState } from "../core/types.js";

export interface SubscriptionHandle {
  id: string;
  name: string;
  state: string;
  disconnect: () => void;
}

export class SubscriptionManager {
  private subscriptions = new Map<string, SubscriptionHandle>();

  create(_options: any): SubscriptionHandle {
    const id = `${Date.now()}-${Math.random().toString(36).slice(2, 8)}`;
    const protocol = new SubscriptionProtocol({
      ..._options,
      onClose: () => {},
      onError: () => {},
      onEvent: () => {},
    });
    protocol.connect();

    const handle: SubscriptionHandle = {
      id,
      name: _options.query.slice(0, 50).replace(/\s+/g, " "),
      state: "active",
      disconnect: () => {
        protocol.disconnect();
        this.subscriptions.delete(id);
      },
    };

    this.subscriptions.set(id, handle);
    return handle;
  }

  remove(id: string): void {
    this.subscriptions.delete(id);
  }

  getActiveSubscriptions(): SubscriptionHandle[] {
    return Array.from(this.subscriptions.values()).filter((s) => s.state === "active");
  }

  disconnectAll(): void {
    for (const [, handle] of this.subscriptions) handle.disconnect();
    this.subscriptions.clear();
  }

  getSubscriptionCount(): number {
    return this.subscriptions.size;
  }
}

let globalManager: SubscriptionManager | null = null;

export function getSubscriptionManager(): SubscriptionManager {
  if (!globalManager) globalManager = new SubscriptionManager();
  return globalManager;
}

export function resetSubscriptionManager(): void {
  if (globalManager) {
    globalManager.disconnectAll();
    globalManager = null;
  }
}
