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

/**
 * The desktop-facing description of a pending, approvable tool call.
 *
 * `inputs`, when present, are fields the user completes as part of approving.
 * `unresolved_inputs` lists inputs whose defaults cannot satisfy the schema:
 * while it is non-empty the daemon never approves automatically and the intent
 * waits for the user.
 */
export interface Intent {
  intent_id: string;
  intent_type: IntentType;
  summary: string;
  payload?: Record<string, unknown>;
  inputs?: IntentInputField[];
  unresolved_inputs?: IntentInputIssue[];
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
 * One approval input. `default` has the field's value type and is what
 * automatic approval uses; the daemon validates it like user input.
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

/** The desktop's approve body. A missing value falls back to the default. */
export interface ApproveIntentRequest {
  inputs?: IntentInputValues;
}
