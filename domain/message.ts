/**
 * Agent messages: durable, daemon-routed text delivered into a running agent
 * session through the provider's native route. The first kind is the
 * ticket-conclusion notice to a project architect; other kinds and recipients
 * share the record, router and status.
 *
 * States distinguish transport acceptance (`handed`), appearance in the
 * conversation (`appeared`) and an agent reply after it (`acknowledged`); none
 * means the agent acted on the message. A message is retained while no
 * recipient session runs and never launches one.
 */
/**
 * `question_escalated`: a ticket worker's question the user has not answered
 * after the configured delay, addressed to the project architect, which may
 * answer it with `answerQuestion`.
 */
export type AgentMessageKind = "ticket_concluded" | "action_ended" | "question_escalated" | "workspace_validation";

export type AgentMessageState =
  | "pending_commit"
  | "pending"
  | "handed"
  | "appeared"
  | "acknowledged"
  | "failed";

export interface AgentMessage {
  id: string;
  kind: AgentMessageKind;
  /** `architect:<key>` or `session:<id>`. */
  recipient: string;
  architect_key?: string;
  dedupe_key: string;
  body: string;
  payload?: Record<string, unknown>;
  state: AgentMessageState;
  /** Why an open message is not progressing right now. */
  condition?: string;
  target_session_id?: string;
  target_agent?: string;
  target_native_id?: string;
  target_config_root?: string;
  target_workdir?: string;
  mechanism?: string;
  provider_ref?: string;
  attempts: number;
  last_error?: string;
  next_attempt_at?: string;
  created_at: string;
  committed_at?: string;
  handed_at?: string;
  appeared_at?: string;
  acknowledged_at?: string;
  ended_at?: string;
  updated_at: string;
}

export interface AgentMessageList {
  messages: AgentMessage[];
}

export function agentMessageIsOpen(state: AgentMessageState): boolean {
  return state === "pending_commit" || state === "pending" || state === "handed";
}
