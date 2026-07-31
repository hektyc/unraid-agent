export type SubscriptionState = "idle" | "connecting" | "active" | "error" | "reconnecting";

export interface SubscriptionResource {
  name: string;
  action: string;
  subaction: string;
  state: SubscriptionState;
  lastError?: string;
}
