import type { SessionType } from "./session";

/**
 * An Intent is a pending agent tool call awaiting the user's approval; the
 * desktop renders it as a popup and approves or denies it. It is either
 * blocking — the agent's call waits and, if nobody answers within the wait
 * window, the tool's policy resolves it — or deferred (policy `manual`): the
 * request returns at once with a stable id, only the user resolves it, and its
 * outcome is a `DeferredIntent` retrieved by that id. Every intent that asks for
 * approval inputs is deferred.
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

/**
 * The tool's hardcoded behavior when the wait window expires. `manual` is
 * deferred: there is no window, and nothing resolves it but the user.
 */
export type IntentPolicy = "auto-allow" | "wait-then-allow" | "wait-then-deny" | "manual";

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

/**
 * The desktop-facing description of a pending, approvable tool call.
 *
 * `inputs`, when present, are fields the user completes as part of approving;
 * such an intent is always deferred. `wait_seconds` is 0 for a deferred intent.
 */
export interface Intent {
  intent_id: string;
  intent_type: IntentType;
  summary: string;
  payload?: Record<string, unknown>;
  inputs?: IntentInputField[];
  origin: IntentOrigin;
  wait_seconds: number;
  policy: IntentPolicy;
  created_at: string;
}

/**
 * The closed set of approval input kinds. text, textarea and choice carry a
 * string; boolean carries a boolean.
 */
export type IntentInputType = "text" | "textarea" | "choice" | "boolean";

export const MAX_INTENT_INPUT_FIELDS = 16;
export const MAX_INTENT_INPUT_OPTIONS = 200;

/** One allowed value of a choice input. */
export interface IntentInputOption {
  value: string;
  /** Display text; `value` when absent. */
  label?: string;
  description?: string;
}

/**
 * One approval input. `default` has the field's value type and only prefills
 * the form: the daemon never approves with it or fills a missing value from it.
 * `max_length` (characters) applies to text and textarea only.
 */
export interface IntentInputField {
  name: string;
  label: string;
  description?: string;
  type: IntentInputType;
  required?: boolean;
  default?: string | boolean;
  options?: IntentInputOption[];
  max_length?: number;
}

export type IntentInputValue = string | boolean;

/** Submitted or resolved input values keyed by field name. */
export type IntentInputValues = Record<string, IntentInputValue>;

/** One problem with an input value, keyed by field name. */
export interface IntentInputIssue {
  field: string;
  message: string;
}

/** The desktop's approve body. Values are taken as submitted; no defaults apply. */
export interface ApproveIntentRequest {
  inputs?: IntentInputValues;
}

/**
 * Lifecycle of a deferred intent. Approval and execution are separate facts:
 * `running` means approved and started; only `completed` means it succeeded.
 * `failed` with no `approved_at` never ran; with `approved_at` it started and
 * may or may not have taken effect.
 */
export type DeferredIntentStatus =
  | "pending_approval"
  | "denied"
  | "running"
  | "completed"
  | "failed";

/**
 * The durable, id-addressable outcome of a deferred intent: returned
 * immediately on request (`pending_approval`) and by a lookup by id afterwards,
 * including after its session ended. Scoped to its origin session.
 */
export interface DeferredIntent {
  intent_id: string;
  intent_type: IntentType;
  summary: string;
  payload?: Record<string, unknown>;
  origin: IntentOrigin;
  status: DeferredIntentStatus;
  /** The validated values the operation ran with, once approved. */
  inputs?: IntentInputValues;
  result?: unknown;
  /** The user's denial reason. */
  reason?: string;
  /** Why it failed. */
  error?: string;
  created_at: string;
  approved_at?: string;
  ended_at?: string;
}
