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
