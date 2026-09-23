import type { CommitRef } from "./commit";
import type { TicketOutcome } from "./ticket";

export type SessionType = "architect" | "ticket" | "freeform";

export type SessionCreatedBy = "desktop";

export type SessionRunStatus = "running" | "completed" | "failed";

export type SessionRunFailureReason =
  | "launch_failed"
  | "process_exited"
  | "restore_failed"
  | "user_cancelled";

export interface AgentProfileSnapshot {
  agent: string;
  model?: string;
  yolo?: boolean;
  mode?: string;
  args: string[];
  env: Record<string, string>;
}

/**
 * The durable spawn record. `prompt` is the fixed daemon-rendered kickoff and
 * `instructions` the built-in role instructions, both resolved at create time.
 * `workflows` is the explicit workflow selection of a ticket session: canonical
 * absolute paths of files directly inside the architect workspace's
 * `workflows/` directory, exactly as the launch request named them. It may be
 * empty, is never inferred from repository membership, and is retained across
 * runs and resumes; the files are read live by the worker, never copied.
 */
export interface Session {
  id: string;
  architect_key: string;
  session_type: SessionType;
  context_id: string;
  prompt: string;
  workdir: string;
  additional_repos: string[];
  additional_workdirs: string[];
  workflows: string[];
  instructions?: string;
  created_by?: SessionCreatedBy;
  created_at: string;
  updated_at: string;
  current_run?: SessionRun;
}

export interface SessionRun {
  id: string;
  session_id: string;
  status: SessionRunStatus;
  agent_status?: string;
  profile_name: string;
  profile_snapshot?: AgentProfileSnapshot;
  workdir: string;
  additional_repos: string[];
  additional_workdirs: string[];
  native_id?: string;
  failure_reason?: SessionRunFailureReason;
  main_terminal_id?: string;
  started_at?: string;
  ended_at?: string;
  created_at: string;
  updated_at: string;
}

export interface SessionEvent {
  id: string;
  session_id: string;
  run_id?: string;
  seq: number;
  type: string;
  status?: string;
  tool?: string;
  message?: string;
  native_id?: string;
  primary_native_id?: string;
  native_session_role?: string;
  metadata?: Record<string, string>;
  raw?: Record<string, unknown>;
  at: string;
}

/**
 * The desktop's launch contract. `workflows` is accepted for ticket sessions
 * only: the explicit, user-confirmed selection of canonical workflow paths (as
 * listed by `GET /api/architects/{key}/workflows`). Omitted or empty is a valid
 * selection of nothing. Every entry is validated; an invalid one rejects the
 * request instead of being dropped or substituted.
 */
export interface CreateSessionRequest {
  session_type: SessionType;
  architect_key: string;
  ticket_id?: string;
  prompt?: string;
  workdir?: string;
  slug?: string;
  workflows?: string[];
}

export interface CreateSessionParams {
  id: string;
  architect_key: string;
  session_type: SessionType;
  context_id: string;
  prompt: string;
  workdir: string;
  workflows: string[];
  instructions: string;
  created_by: SessionCreatedBy;
}

export interface CreateSessionRunRequest {
  profile_name: string;
  cols?: number;
  rows?: number;
}

export interface CreateSessionRunParams {
  id: string;
  session_id: string;
  profile_name: string;
  profile_snapshot: AgentProfileSnapshot;
  workdir: string;
  native_id: string;
  started_at: string;
}

export interface CreateSessionRunResult {
  run: SessionRun;
  main_terminal_id?: string;
}

export interface AppendSessionEventParams {
  session_id: string;
  run_id: string;
  type: string;
  status: string;
  tool: string;
  message: string;
  native_id: string;
  primary_native_id: string;
  native_session_role: string;
  metadata: Record<string, string>;
  raw: Record<string, unknown>;
  at: string;
}

export interface ConcludeSessionParams {
  body: string;
  commits: CommitRef[];
  outcome: TicketOutcome;
  rejection_reason: string;
}

export interface ConcludeSessionResult {
  session_id: string;
  architect_key: string;
  ticket_id?: string;
}

export interface MoveTicketToDoneParams {
  body: string;
  commits: CommitRef[];
  outcome: TicketOutcome;
  rejection_reason: string;
}

export interface MoveTicketToDoneResult {
  ticket_id: string;
  architect_key: string;
}

export interface TerminalInfo {
  terminal_id: string;
  session_id: string;
  command: string;
  status: string;
  workdir_id?: string;
  workdir_title?: string;
  workdir_path?: string;
  workdir_display_path?: string;
}

export interface TerminalWorkdir { id: string; title: string; path: string; display_path: string; default: boolean; }

export type TerminalPlacement = 'tab' | 'split';

export type CreateTerminalParams =
  | { placement: 'tab'; workdir_id: string }
  | { placement: 'split'; base_tab_id: string; workdir_id: string };

export interface SessionTab {
  type: string;
  id?: string;
  command?: string;
  status?: string;
  placement?: TerminalPlacement;
  base_tab_id?: string;
  workdir_id?: string;
  workdir_title?: string;
  workdir_path?: string;
  workdir_display_path?: string;
}

export interface ArchitectConclusion {
  started_at: string;
  concluded_at: string;
  agent?: string;
  body: string;
}

export interface ConclusionSummary {
  id: string;
  concluded_at: string;
}
