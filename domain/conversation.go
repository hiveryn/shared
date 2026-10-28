package domain

import "time"

type InteractionMode string

const (
	InteractionTerminal InteractionMode = "terminal"
	InteractionChat     InteractionMode = "chat"
)

type ConversationCapabilities struct {
	History           bool   `json:"history"`
	Send              bool   `json:"send"`
	NativePermissions bool   `json:"native_permissions"`
	ExecutionLocation string `json:"execution_location"`
	// Machine is the configured machine key of a remote session (reached over
	// SSH from the daemon host); empty for local sessions.
	Machine    string `json:"machine,omitempty"`
	Limitation string `json:"limitation,omitempty"`
	// InputRoute names how send reaches the native conversation: "chat" for the
	// daemon-owned structured transport, or the provider route into an existing
	// terminal-owned session ("claude_async_rewake", "codex_thread_queue",
	// "opencode_prompt_async"). Empty when send is unavailable.
	InputRoute string `json:"input_route,omitempty"`
}
type SubmissionRequest struct {
	RequestID string `json:"request_id"`
	Text      string `json:"text"`
}

// Submission is the durable receipt of one conversation input. User input is
// delivered synchronously by POST /messages: the receipt names the outcome of
// that attempt (Status, ErrorCode) and, separately, whether the provider
// accepted it (AcceptedAt/Acceptance). Turn completion is later and distinct.
type Submission struct {
	ID        string `json:"id"`
	SessionID string `json:"session_id"`
	RequestID string `json:"request_id"`
	Sequence  int64  `json:"sequence"`
	Origin    string `json:"origin"`
	Text      string `json:"text"`
	// Status is one of the SubmissionStatus constants.
	Status       string `json:"status"`
	NativeTurnID string `json:"native_turn_id,omitempty"`
	// ErrorCode is the stable reason for a rejected or unknown delivery (a
	// Delivery* / Session* / Send* code); Error is its readable detail, or the
	// turn's failure.
	ErrorCode string `json:"error_code,omitempty"`
	Error     string `json:"error,omitempty"`
	// AcceptedAt is when provider acceptance was established; Acceptance names
	// that evidence (for example "codex_thread_queue", "claude_native_record",
	// "opencode_prompt_async", "chat_turn_started"). Absent until accepted.
	AcceptedAt *time.Time `json:"accepted_at,omitempty"`
	Acceptance string     `json:"acceptance,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

// Submission statuses. User input is never queued by Hiveryn: it is
// delivering, then rejected, unknown or accepted; accepted input continues to
// running and a turn outcome. queued exists only for the chat owner's internal
// kickoff/notification entries.
const (
	SubmissionQueued     = "queued"
	SubmissionDelivering = "delivering"
	SubmissionRejected   = "rejected"
	SubmissionUnknown    = "unknown"
	SubmissionAccepted   = "accepted"
	SubmissionRunning    = "running"
	SubmissionBlocked    = "blocked"
	SubmissionCompleted  = "completed"
	SubmissionFailed     = "failed"
	SubmissionUncertain  = "uncertain"
)

// Stable delivery codes of a send that did not establish acceptance.
const (
	// Definitive: nothing reached the provider and nothing will later.
	DeliverySessionBusy         = "SESSION_BUSY"
	DeliverySessionNotRunning   = "SESSION_NOT_RUNNING"
	DeliverySessionEnding       = "SESSION_ENDING"
	DeliverySessionArchived     = "SESSION_ARCHIVED"
	DeliverySendUnavailable     = "SEND_UNAVAILABLE"
	DeliverySendUnsupported     = "SEND_UNSUPPORTED"
	DeliveryRejected            = "DELIVERY_REJECTED"
	DeliveryInterrupted         = "DELIVERY_INTERRUPTED"
	DeliveryQueuedBeforeCutover = "QUEUED_BEFORE_CUTOVER"
	// Ambiguous: a hand-off was attempted and its acceptance is not
	// established. Never resent; reconciled from native evidence.
	DeliveryUnknown = "DELIVERY_UNKNOWN"
)

// DeliveryError is a send that did not establish provider acceptance. Code is
// a Delivery* constant; Submission is the durable receipt when one exists
// (every case except an archived or missing session).
type DeliveryError struct {
	Code       string
	Message    string
	Submission *Submission
}

func (e *DeliveryError) Error() string { return e.Code + ": " + e.Message }

// Stable codes of a conversation history read that returned no page.
const (
	// HistoryUnsupported: this session's provider/location has no history
	// route (remote OpenCode, an agent without a reader). Not retryable.
	HistoryUnsupported = "HISTORY_UNSUPPORTED"
	// HistoryUnavailable: the route exists but not now (no native id yet).
	HistoryUnavailable = "HISTORY_UNAVAILABLE"
	// HistoryReadFailed: the native record could not be read (unreachable
	// machine, SSH failure, unreadable or ambiguous record). Retry later;
	// nothing was returned, never an empty or silently stale page.
	HistoryReadFailed = "HISTORY_READ_FAILED"
)

// ConversationError is a history read that returned no page; Code is a
// History* constant and Message keeps the original cause.
type ConversationError struct {
	Code    string
	Message string
}

func (e *ConversationError) Error() string { return e.Code + ": " + e.Message }

type ConversationItem struct {
	ID           string `json:"id"`
	SubmissionID string `json:"submission_id"`
	// Sequence is the display position within the session, following the
	// native record's order (never timestamps). It can change: an item
	// discovered late or stored out of record order is moved between its
	// native neighbours as a new revision. Order by Sequence, then ID.
	// Revision is the feed cursor (discovery order), for incremental
	// synchronization only; NativeTurnID/SubmissionID are turn identity.
	Sequence     int64  `json:"sequence"`
	Revision     int64  `json:"revision"`
	NativeTurnID string `json:"native_turn_id,omitempty"`
	Role         string `json:"role"`
	Origin       string `json:"origin"`
	// Via says how a user-role item entered the conversation: "terminal" (typed
	// in the session's own TUI), "api" (a conversation API submission),
	// "notification" (the daemon's agent-message outbox) or "kickoff" (the
	// launch message Hiveryn sent, as the provider recorded it, including any
	// injected Overview/workflows; never SYSTEM/developer instructions). Items
	// of the kickoff's turn carry origin "kickoff". NativeOrigin is the
	// provider's own unmodified provenance for that input (for example Claude
	// turnOrigin "human" or "task_notification", Codex turn trigger "user" or
	// "queue", OpenCode "prompt"). Both are omitted where unknown.
	Via          string    `json:"via,omitempty"`
	NativeOrigin string    `json:"native_origin,omitempty"`
	Text         string    `json:"text,omitempty"`
	Tool         string    `json:"tool,omitempty"`
	Status       string    `json:"status"`
	At           time.Time `json:"at,omitzero"`
}

// Items are in revision (feed) order, not display order. A page is bounded by SnapshotRevision;
// continue with that value, then use it as AfterRevision for the next poll.
type ConversationPage struct {
	Items            []ConversationItem `json:"items"`
	SnapshotRevision int64              `json:"snapshot_revision"`
	NextRevision     int64              `json:"next_revision"`
	HasMore          bool               `json:"has_more"`
}
type ConversationQuery struct {
	AfterRevision, SnapshotRevision int64
	Limit                           int
}

// Conversation stream (GET /api/sessions/{id}/stream, server-sent events).
// One connection per selected session carries three independently stored
// sources: conversation revisions (event "conversation", data a
// ConversationPage, SSE id = its next_revision — the reconnect cursor),
// question/intent controls ("controls", the complete durable control event
// history, as GET /controls) and session status ("status"). Controls and
// status are authoritative snapshots, sent on every connect and whenever they
// change: the latest one replaces the previous. Only conversation pages are
// incremental (upsert by id, highest revision wins). "reset" says the cursor
// is not part of this conversation; the client drops its items and the stream
// replays from revision 0.
const (
	StreamEventStatus       = "status"
	StreamEventControls     = "controls"
	StreamEventConversation = "conversation"
	StreamEventReset        = "reset"
)

// Stream states: live (the session exists; history and send follow its
// capabilities) or ended (concluded or discarded: the final history has been
// sent before this status and the stream closes; reconnecting is pointless).
const (
	StreamStateLive  = "live"
	StreamStateEnded = "ended"
)

// ConversationStreamStatus is the "status" event. Session is the live
// session record without its kickoff prompt, role instructions and variant
// environment/arguments. HistoryError is the latest native record read's
// failure (a History* code with the original cause); it clears when a later
// read succeeds, and already-sent history stays valid meanwhile.
type ConversationStreamStatus struct {
	State        string                   `json:"state"`
	Session      *Session                 `json:"session,omitempty"`
	HistoryError *ConversationStreamError `json:"history_error,omitempty"`
}

type ConversationStreamError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// ConversationStreamControls is the "controls" event: every durable
// question/intent event of the session in seq order, reduced by the client.
type ConversationStreamControls struct {
	Events []SessionEvent `json:"events"`
}

// ConversationStreamReset is the "reset" event.
type ConversationStreamReset struct {
	Reason string `json:"reason"`
}
