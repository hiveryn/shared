export type InteractionMode = "terminal" | "chat";
export interface ConversationCapabilities {
  history: boolean;
  send: boolean;
  native_permissions: boolean;
  execution_location: "local" | "remote";
  limitation?: string;
}
export interface SubmissionRequest {
  request_id: string;
  text: string;
}
export interface Submission {
  id: string;
  session_id: string;
  request_id: string;
  sequence: number;
  origin: string;
  text: string;
  status: "queued" | "started" | "completed" | "failed" | "blocked" | "uncertain";
  native_turn_id?: string;
  error?: string;
  created_at: string;
  updated_at: string;
}
export interface ConversationItem {
  id: string;
  submission_id: string;
  sequence: number;
  revision: number;
  native_turn_id?: string;
  role: string;
  origin: string;
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
}
export interface ArchivedChatPage {
  chats: ArchivedChat[];
  next_sequence: number;
  has_more: boolean;
}
