import type { CommitRef } from "./commit";

export type SessionType = "architect" | "ticket" | "freeform";

export type SessionCreatedBy = "desktop" | "architect_mcp";

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

export interface SessionIntent {
  id: string;
  architect_key: string;
  session_type: SessionType;
  context_id: string;
  prompt: string;
  workdir: string;
  instructions?: string;
  created_by?: SessionCreatedBy;
  created_at: string;
  updated_at: string;
  current_run?: SessionRun;
}

export interface SessionRun {
  id: string;
  session_intent_id: string;
  status: SessionRunStatus;
  agent_status?: string;
  profile_name: string;
  profile_snapshot?: AgentProfileSnapshot;
  workdir: string;
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
  session_intent_id: string;
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

export interface CreateSessionIntentRequest {
  session_type: SessionType;
  architect_key: string;
  ticket_id?: string;
  prompt?: string;
  workdir?: string;
  slug?: string;
}

export interface CreateSessionIntentParams {
  id: string;
  architect_key: string;
  session_type: SessionType;
  context_id: string;
  prompt: string;
  workdir: string;
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
  session_intent_id: string;
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
  session_intent_id: string;
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
  rejected: boolean;
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
  rejected: boolean;
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
}

export type TerminalPlacement = 'tab' | 'split';

export type CreateTerminalParams =
  | { placement: 'tab' }
  | { placement: 'split'; base_tab_id: string };

export interface SessionTab {
  type: string;
  id?: string;
  command?: string;
  status?: string;
  placement?: TerminalPlacement;
  base_tab_id?: string;
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
