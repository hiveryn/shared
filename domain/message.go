package domain

import "time"

// Agent messages: durable, daemon-routed text delivered into a running agent
// session through the provider's native route (never by typing into its
// terminal). The first kind is the ticket-conclusion notice to a project
// architect; other kinds and recipients (any session) share the same record,
// router and status.
//
// A message is retained while no recipient session runs and is never a reason
// to launch one. States distinguish transport acceptance (handed), appearance
// in the conversation (appeared) and an agent reply after it (acknowledged);
// none of them means the agent acted on the message. Delivery is effectively
// once: a stable id, serialized attempts and a provider receipt check before
// any redelivery prevent avoidable duplicates, but provider crash windows
// documented per mechanism remain.

// AgentMessageKind names what a message announces.
type AgentMessageKind string

const AgentMessageWorkspaceValidation AgentMessageKind = "workspace_validation"

const (
	// AgentMessageTicketConcluded: a worker's approved, persisted ticket
	// conclusion. Payload: ticket_id, outcome, concluded_at.
	AgentMessageTicketConcluded AgentMessageKind = "ticket_concluded"
	// AgentMessageActionEnded: an agent-requested Action execution reached
	// its final status — the agent's approved conclusion, the user's cancel,
	// a launch or resume failure, or a daemon-restart interruption — and
	// getActionResult serves it. Addressed to the requesting agent. Payload:
	// execution_id, action, outcome (completed|failed), error (the daemon's
	// reason when the agent did not conclude), package (bool).
	AgentMessageActionEnded AgentMessageKind = "action_ended"
	// AgentMessageQuestionEscalated: a ticket worker's askQuestion has gone
	// unanswered by the user for the configured delay; the project architect
	// may answer it with answerQuestion. Addressed to the architect. Payload:
	// question_id, session_id, ticket_id, question, answers,
	// recommended_index, created_at, expires_at, elapsed_seconds. Exactly one
	// per question; withdrawn (failed) when the question resolves before it
	// is delivered.
	AgentMessageQuestionEscalated AgentMessageKind = "question_escalated"
)

// AgentMessageState is the delivery lifecycle.
type AgentMessageState string

const (
	// AgentMessagePendingCommit: the producer's side effect (for example the
	// conclusion files) is not yet durable; the message is never delivered in
	// this state and is reconciled against that side effect.
	AgentMessagePendingCommit AgentMessageState = "pending_commit"
	// AgentMessagePending: deliverable; waiting for a running recipient
	// session, its idle point or a retry. Condition says which.
	AgentMessagePending AgentMessageState = "pending"
	// AgentMessageHanded: the provider accepted it (waiter woken, queued,
	// HTTP accepted). Not yet seen in the conversation.
	AgentMessageHanded AgentMessageState = "handed"
	// AgentMessageAppeared: the provider's own record shows it in the
	// conversation.
	AgentMessageAppeared AgentMessageState = "appeared"
	// AgentMessageAcknowledged: the agent produced a reply after it.
	AgentMessageAcknowledged AgentMessageState = "acknowledged"
	// AgentMessageFailed: terminal; LastError says why (recipient gone,
	// producer never committed, superseded).
	AgentMessageFailed AgentMessageState = "failed"
)

// Open reports whether the message still needs delivery work.
func (s AgentMessageState) Open() bool {
	switch s {
	case AgentMessagePendingCommit, AgentMessagePending, AgentMessageHanded:
		return true
	default:
		return false
	}
}

// Recipient schemes: "architect:<key>" addresses whichever session of that
// architect is running (at most one); "session:<id>" addresses one session.
const (
	AgentMessageRecipientArchitect = "architect:"
	AgentMessageRecipientSession   = "session:"
)

// AgentMessage is the durable outbox record.
type AgentMessage struct {
	ID        string           `json:"id"`
	Kind      AgentMessageKind `json:"kind"`
	Recipient string           `json:"recipient"`
	// ArchitectKey scopes listing; empty for messages outside any project.
	ArchitectKey string `json:"architect_key,omitempty"`
	// DedupeKey is the producer's stable identity of the announcement; a
	// second message with the same key is the same message.
	DedupeKey string            `json:"dedupe_key"`
	Body      string            `json:"body"`
	Payload   map[string]any    `json:"payload,omitempty"`
	State     AgentMessageState `json:"state"`
	// Condition says why an open message is not progressing right now (no
	// running recipient, agent busy, no waiter armed, delivery failed...).
	Condition string `json:"condition,omitempty"`
	// Target is the session the latest attempt addressed and how.
	TargetSessionID string `json:"target_session_id,omitempty"`
	TargetAgent     string `json:"target_agent,omitempty"`
	TargetNativeID  string `json:"target_native_id,omitempty"`
	// TargetConfigRoot and TargetWorkdir locate the target's provider record
	// (transcript, rollout, database) after the session itself is gone.
	TargetConfigRoot string `json:"target_config_root,omitempty"`
	TargetWorkdir    string `json:"target_workdir,omitempty"`
	Mechanism        string `json:"mechanism,omitempty"`
	// ProviderRef is the provider-side identity of the handed message.
	ProviderRef string `json:"provider_ref,omitempty"`
	Attempts    int    `json:"attempts"`
	LastError   string `json:"last_error,omitempty"`
	// NextAttemptAt holds a failed attempt's retry time.
	NextAttemptAt  *time.Time `json:"next_attempt_at,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	CommittedAt    *time.Time `json:"committed_at,omitempty"`
	HandedAt       *time.Time `json:"handed_at,omitempty"`
	AppearedAt     *time.Time `json:"appeared_at,omitempty"`
	AcknowledgedAt *time.Time `json:"acknowledged_at,omitempty"`
	EndedAt        *time.Time `json:"ended_at,omitempty"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

// AgentMessageList is the architect-scoped listing.
type AgentMessageList struct {
	Messages []AgentMessage `json:"messages"`
}
