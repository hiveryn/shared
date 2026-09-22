import type { SessionType } from "./session";

/**
 * An Intent is a pending agent tool call awaiting the user's approval. The
 * agent's MCP call blocks in the daemon while the intent is pending; the
 * desktop renders it as a popup and approves or denies it. If nobody answers
 * within the wait window, the tool's hardcoded policy resolves it.
 *
 * Distinct from Session (the spawn record): a Session is "run an agent", an
 * Intent is "the agent wants to do something first".
 */
export type IntentType = "concludeSession" | "createWorkTicket";

/**
 * The agent-facing verdict. `denied_*` means the action did NOT run and must
 * not be retried; `error` means it failed and may be retried.
 */
export type IntentOutcome =
  | "approved"
  | "auto_approved"
  | "denied_by_user"
  | "auto_denied"
  | "error";

/** The tool's hardcoded behavior when the wait window expires. */
export type IntentPolicy = "auto-allow" | "wait-then-allow" | "wait-then-deny";

/**
 * Where an intent came from. Lets the UI render one popup for any tool across
 * any session. Derived from the Session record daemon-side.
 */
export interface IntentOrigin {
  architect_key: string;
  session_id: string;
  session_type: SessionType;
  ticket_id?: string;
}

/** The desktop-facing description of a pending, approvable tool call. */
export interface Intent {
  intent_id: string;
  intent_type: IntentType;
  summary: string;
  payload?: Record<string, unknown>;
  origin: IntentOrigin;
  wait_seconds: number;
  policy: IntentPolicy;
  created_at: string;
}
