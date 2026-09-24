import type { Session } from "./session";

/**
 * Actions are global, agent-operated procedures, independent of architects.
 * Each lives in its own Git repository at `HIVERYN_HOME/actions/<name>/`,
 * defined by `action.yaml` (`name`, `description`, `artifacts`) and
 * `KICKOFF.md` (containing `{{prompt}}` and `{{output_dir}}`). At most one
 * execution of an action runs at a time, globally. See action.go.
 */

/** Mirrors `ValidActionName` in action.go. */
export const ACTION_NAME_PATTERN = /^[a-z0-9][a-z0-9._-]{0,63}$/;

export interface ActionProblem {
  path: string;
  message: string;
}

export interface ActionDefinition {
  name: string;
  path: string;
  description: string;
  artifacts: string;
  valid: boolean;
  problems: ActionProblem[];
  running_execution_id?: string;
}

export interface ActionList {
  root: string;
  actions: ActionDefinition[];
}

/**
 * Manual launches start `running` and end `completed` or `failed`. An
 * architect's request starts `pending_approval` and moves to `denied`,
 * `failed` (could not start or was abandoned) or `running` once approved.
 * `completed` means the artifact package was delivered (it may still report
 * findings); `failed` means the action could not be executed or delivered.
 */
export type ActionRunStatus =
  | "pending_approval"
  | "denied"
  | "running"
  | "completed"
  | "failed";

export type ActionRunTrigger = "manual" | "architect";

export interface ActionRun {
  id: string;
  action: string;
  trigger: ActionRunTrigger;
  status: ActionRunStatus;
  prompt: string;
  profile_name: string;
  repo_path: string;
  output_dir: string;
  session_id?: string;
  summary?: string;
  error?: string;
  created_at: string;
  started_at?: string;
  ended_at?: string;
  /** Architect requests only: who may read the result, who asked, why it was denied. */
  architect_key?: string;
  requester_session_id?: string;
  reason?: string;
}

export interface LaunchActionRequest {
  prompt: string;
  profile_name: string;
  cols?: number;
  rows?: number;
}

export interface LaunchActionResult {
  run: ActionRun;
  session: Session;
  main_terminal_id?: string;
}

export type ActionConclusionOutcome = "completed" | "failed";

export const MAX_ACTION_SUMMARY_LENGTH = 2000;

export interface ConcludeActionRequest {
  outcome: ActionConclusionOutcome;
  summary: string;
}

export interface ActionConclusion {
  execution_id: string;
  status: ActionRunStatus;
  summary: string;
  ended_at?: string;
}

export interface ExecuteActionRequest {
  name: string;
  prompt: string;
}

/** The Actions an architect may request, in hiveryn.yaml `availableActions` order. */
export interface AvailableActionList {
  actions: ActionDefinition[];
}

/** `available: false` means nothing is known; `status` is then absent, not guessed. */
export interface ActionAgentActivity {
  available: boolean;
  status?: string;
}

/** An architect's view of one requested execution. See action.go. */
export interface ActionResult {
  execution_id: string;
  action: string;
  status: ActionRunStatus;
  prompt: string;
  profile_name?: string;
  requested_at: string;
  started_at?: string;
  ended_at?: string;
  elapsed_seconds?: number;
  reason?: string;
  error?: string;
  summary?: string;
  output_dir?: string;
  activity: ActionAgentActivity;
}

export const MAX_ACTION_WAIT_SECONDS = 30;

export interface ActionWaitResult {
  result: ActionResult;
  changed: boolean;
  timed_out: boolean;
}

export const ACTION_EVENT_TYPE = "action_changed";

/** An execution was requested, started or ended. No backlog: refetch executions on (re)connect. */
export interface ActionEvent {
  type: typeof ACTION_EVENT_TYPE;
  action: string;
  execution_id: string;
  status: ActionRunStatus;
  session_id?: string;
  at: string;
}
