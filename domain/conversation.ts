export type InteractionMode = "terminal" | "chat";
export interface ConversationCapabilities {
  history: boolean;
  send: boolean;
  native_permissions: boolean;
  execution_location: "local" | "remote";
  limitation?: string;
  /** "chat", or the provider route into a terminal-owned session; absent when send is unavailable. */
  input_route?: "chat" | "claude_async_rewake" | "codex_thread_queue" | "opencode_prompt_async";
}
export interface SubmissionRequest {
  request_id: string;
  text: string;
}
/**
 * User input is never queued by Hiveryn: delivering, then rejected, unknown or
 * accepted; accepted input continues to running and a turn outcome. "queued"
 * exists only for the chat owner's internal kickoff/notification entries.
 */
export type SubmissionStatus =
  | "queued"
  | "delivering"
  | "rejected"
  | "unknown"
  | "accepted"
  | "running"
  | "blocked"
  | "completed"
  | "failed"
  | "uncertain";

/** Stable codes of a send that did not establish acceptance. DELIVERY_UNKNOWN is ambiguous; all others are definitive. */
export type DeliveryCode =
  | "SESSION_BUSY"
  | "SESSION_NOT_RUNNING"
  | "SESSION_ENDING"
  | "SESSION_ARCHIVED"
  | "SEND_UNAVAILABLE"
  | "SEND_UNSUPPORTED"
  | "DELIVERY_REJECTED"
  | "DELIVERY_INTERRUPTED"
  | "QUEUED_BEFORE_CUTOVER"
  | "DELIVERY_UNKNOWN";

export interface Submission {
  id: string;
  session_id: string;
  request_id: string;
  sequence: number;
  origin: string;
  text: string;
  status: SubmissionStatus;
  native_turn_id?: string;
  /** Stable reason of a rejected or unknown delivery. */
  error_code?: DeliveryCode;
  error?: string;
  /** When provider acceptance was established, and the evidence that established it. */
  accepted_at?: string;
  acceptance?: string;
  created_at: string;
  updated_at: string;
}
export interface ConversationItem {
  id: string;
  submission_id: string;
  /**
   * Display position within the session, following the native record's order
   * (never timestamps). It can change: a late-discovered or out-of-order item
   * is moved between its native neighbours as a new revision. Order by
   * sequence, then id.
   */
  sequence: number;
  /** Feed cursor (discovery order); for incremental synchronization only. */
  revision: number;
  /** Turn identity, separate from order. */
  native_turn_id?: string;
  role: string;
  origin: string;
  /**
   * How a user-role item entered: typed in the TUI, an API submission, an
   * outbox notification, or the kickoff (the launch message Hiveryn sent,
   * recorded text including any injected Overview/workflows).
   */
  via?: "terminal" | "api" | "notification" | "kickoff";
  /** The provider's own unmodified provenance for the input (e.g. "human", "task_notification", "queue", "prompt"). */
  native_origin?: string;
  text?: string;
  tool?: string;
  status: string;
  at?: string;
}
export interface ConversationPage {
  items: ConversationItem[];
  snapshot_revision: number;
  next_revision: number;
  has_more: boolean;
}

export interface SubmissionPage {
  submissions: Submission[];
  next_sequence: number;
  has_more: boolean;
}
export interface ArchivedChat {
  id: string;
  sequence: number;
  architect_key: string;
  session_type: "architect" | "ticket";
  context_id: string;
  agent: string;
  model?: string;
  profile_name: string;
  created_at: string;
  ended_at: string;
  lifecycle: "concluded" | "discarded" | "interrupted";
  /** Mode of the ended session; absent on descriptors archived before terminal history retention (chat). */
  interaction_mode?: InteractionMode;
}
export interface ArchivedChatPage {
  chats: ArchivedChat[];
  next_sequence: number;
  has_more: boolean;
}
