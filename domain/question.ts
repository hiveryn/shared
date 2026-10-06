import type { SessionType } from "./session";

/**
 * Agent notifications and questions. An agent may `notify` the user (a phone
 * alert through ntfy, returning at once) or `askQuestion`: a phone alert plus a
 * pending question answered in the originating desktop session, by choosing a
 * suggested answer or typing free text. The agent's call waits at most
 * `QUESTION_TIMEOUT_MS`; the daemon enforces that expiry. A question is not an
 * Intent, and its recommended answer is never selected for the user.
 *
 * Lengths are counted in Unicode code points, as in question.go. See
 * question.go for the authoritative contract.
 */

/** Mirrors the bounds in question.go. */
export const MAX_NOTIFY_MESSAGE_LENGTH = 500;
export const MAX_QUESTION_LENGTH = 1000;
export const MIN_QUESTION_ANSWERS = 2;
export const MAX_QUESTION_ANSWERS = 6;
export const MAX_QUESTION_ANSWER_LENGTH = 200;
export const MAX_QUESTION_RESPONSE_LENGTH = 4000;
export const QUESTION_TIMEOUT_MS = 60 * 60 * 1000;

/** notify's exact reply once the ntfy server accepted the alert. */
export const NOTIFICATION_SENT_TEXT = "Notification sent to user";
/** askQuestion's exact reply when nobody answered within the hour. */
export const QUESTION_TIMEOUT_TEXT =
  "User didn't respond within 1 hour, stop here and wait for user to get back";

/**
 * `pending` until resolved exactly once: `answered` (the agent received the
 * answer), `expired` (nobody answered within the hour), `cancelled` (the agent
 * stopped waiting or its session ended; see `reason`) or `interrupted` (the
 * daemon stopped while it was pending).
 */
export type QuestionStatus = "pending" | "answered" | "expired" | "cancelled" | "interrupted";

/** Derived from the session record; `architect_key` is empty for an Action session. */
export interface QuestionOrigin {
  architect_key: string;
  session_id: string;
  session_type: SessionType;
  ticket_id?: string;
  action?: string;
  execution_id?: string;
}

/**
 * One question as rendered. It arrives on the session event stream as
 * `{ type: "question", status: "required", raw: { question_id, question,
 * answers, recommended_index, origin, created_at, expires_at } }` and ends with
 * `{ type: "question", status: "resolved", raw: { question_id, status, answer?,
 * reason? } }`, published for every resolution.
 */
export interface AgentQuestion {
  question_id: string;
  question: string;
  answers: string[];
  /** 0-based index into `answers`. */
  recommended_index: number;
  origin: QuestionOrigin;
  status: QuestionStatus;
  created_at: string;
  expires_at: string;
  answer?: string;
  reason?: string;
}

/** The desktop's answer: a suggested answer's text or free text, sent verbatim. */
export interface AnswerQuestionRequest {
  answer: string;
}

/** Mirrors `NormalizeQuestionResponse`: the trimmed answer, or the problem. */
export function normalizeQuestionResponse(
  answer: string,
): { ok: true; answer: string } | { ok: false; message: string } {
  const trimmed = answer.trim();
  if (trimmed === "") return { ok: false, message: "answer is required" };
  const length = [...trimmed].length;
  if (length > MAX_QUESTION_RESPONSE_LENGTH) {
    return {
      ok: false,
      message: `answer must be at most ${MAX_QUESTION_RESPONSE_LENGTH} characters, got ${length}`,
    };
  }
  return { ok: true, answer: trimmed };
}
