import type { Session } from "./session";

/**
 * Actions are global, agent-operated procedures, independent of architects.
 * Each lives in its own Git repository at `HIVERYN_HOME/actions/<name>/`,
 * defined by `action.yaml` (`name`, `description`, `artifacts`) and
 * `KICKOFF.md` (containing `{{prompt}}` and `{{output_dir}}`). At most one
 * execution of an action runs at a time, globally. See action.go.
 */

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
 * Manual launches start `running` and end `completed` or `failed`.
 * `pending_approval` and `denied` are reserved for approval-gated launches.
 * `completed` means the artifact package was delivered (it may still report
 * findings); `failed` means the action could not be executed or delivered.
 */
export type ActionRunStatus =
  | "pending_approval"
  | "denied"
  | "running"
  | "completed"
  | "failed";

export type ActionRunTrigger = "manual";

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

export const ACTION_EVENT_TYPE = "action_changed";

/** An execution started or ended. No backlog: refetch executions on (re)connect. */
export interface ActionEvent {
  type: typeof ACTION_EVENT_TYPE;
  action: string;
  execution_id: string;
  status: ActionRunStatus;
  session_id?: string;
  at: string;
}
