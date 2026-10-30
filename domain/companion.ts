import type { SessionType } from "./session";

/**
 * Companions: persistent agent sessions with their own repository, started by
 * the daemon with each new architect session for every `companions` entry of
 * its hiveryn.yaml and owned by that exact architect session until stopped.
 * See companion.go for the authoritative contract.
 */

/**
 * `starting` → `running` → `shutdown_requested` (the architect's request is
 * queued; no new assignments) → `stopping` (handoff recorded) → `stopped`
 * (termination confirmed) or `stop_failed` (not confirmed; blocks the
 * architect's conclusion). `start_failed`: never became available. `failed`:
 * a local companion's process ended and could not be resumed.
 */
export type CompanionStatus =
  | "starting"
  | "running"
  | "shutdown_requested"
  | "stopping"
  | "stopped"
  | "stop_failed"
  | "start_failed"
  | "failed";

/** Statuses in which nothing of the companion can still be running. */
export const SETTLED_COMPANION_STATUSES: readonly CompanionStatus[] = ["stopped", "start_failed", "failed"];

export type CompanionShutdownStatus = "requested" | "completing" | "stopped" | "stop_failed";

export interface CompanionShutdown {
  request_id: string;
  architect_key: string;
  architect_session_id: string;
  companion_session_id: string;
  companion: string;
  status: CompanionShutdownStatus;
  handoff?: string;
  error?: string;
  requested_at: string;
  completed_at?: string;
  stopped_at?: string;
}

/**
 * One configured companion of one architect session. `session_id` is its
 * session (terminals, events, decisions are addressed through it);
 * `main_terminal_id`, `connection` and `agent_status` are live values;
 * `pending_actions` are its undecided action requests.
 */
export interface Companion {
  architect_key: string;
  architect_session_id: string;
  name: string;
  repo: string;
  workdir: string;
  variant: string;
  agent?: string;
  machine?: string;
  session_id?: string;
  status: CompanionStatus;
  error?: string;
  created_at: string;
  updated_at: string;
  shutdown?: CompanionShutdown;
  main_terminal_id?: string;
  connection?: "connected" | "disconnected" | "missing" | string;
  agent_status?: string;
  pending_actions: ActionRequest[];
}

/** GET /api/architects/{key}/companions. */
export interface CompanionList {
  architect_session_id?: string;
  companions: Companion[];
}

export const MAX_COMPANION_MESSAGE_LENGTH = 100000;
export const MAX_COMPANION_HANDOFF_LENGTH = 100000;

/**
 * Action requests (requestAction). Carried on the companion session's event
 * stream as `{ type: "action_request", status: "required", raw: { action_id,
 * title, description, origin, created_at, expires_at } }` and `{ type:
 * "action_request", status: "resolved", raw: { action_id, status, message?,
 * reason? } }`. Decide with POST
 * /api/sessions/{session_id}/actions/{action_id}/decision.
 */
export const MAX_ACTION_TITLE_LENGTH = 200;
export const MAX_ACTION_DESCRIPTION_LENGTH = 50000;
export const MAX_ACTION_REPLY_LENGTH = 4000;
export const ACTION_REQUEST_TIMEOUT_MS = 60 * 60 * 1000;
export const SESSION_EVENT_TYPE_ACTION_REQUEST = "action_request";

export type ActionRequestStatus =
  | "pending"
  | "accepted"
  | "denied"
  | "replied"
  | "expired"
  | "cancelled"
  | "interrupted";

export interface ActionOrigin {
  architect_key: string;
  architect_session_id?: string;
  session_id: string;
  session_type: SessionType;
  companion?: string;
}

export interface ActionRequest {
  id: string;
  title: string;
  /** Markdown proposal. */
  description: string;
  origin: ActionOrigin;
  status: ActionRequestStatus;
  /** The user's reply when `status` is `replied`. */
  message?: string;
  reason?: string;
  created_at: string;
  expires_at: string;
  decided_at?: string;
}

export type ActionDecision = "accept" | "deny" | "reply";

/** `message` is required for `reply` and refused otherwise. */
export interface ActionDecisionRequest {
  decision: ActionDecision;
  message?: string;
}
