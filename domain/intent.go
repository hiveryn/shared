package domain

import "time"

// An Intent is a pending agent tool call awaiting the user's approval. The
// agent's MCP call blocks in the daemon while the intent is pending; the
// desktop renders it as a popup and approves or denies it. If nobody answers
// within the wait window, the tool's hardcoded policy resolves it.
//
// This is deliberately distinct from Session (the spawn record): a Session is
// "run an agent", an Intent is "the agent wants to do something first".

type IntentType string

const (
	IntentTypeConcludeSession  IntentType = "concludeSession"
	IntentTypeCreateWorkTicket IntentType = "createWorkTicket"
)

func (t IntentType) Valid() bool {
	switch t {
	case IntentTypeConcludeSession, IntentTypeCreateWorkTicket:
		return true
	default:
		return false
	}
}

// IntentOutcome is the agent-facing verdict. The catalog exists so an agent can
// tell "the user said no" apart from "it blew up" without parsing prose: the
// denied_* outcomes must read as do-not-retry, error as maybe-retry. Denials
// are returned as successful tool results carrying one of these values, never
// as tool errors — a tool error reads to an agent as "bad input, fix and retry",
// which is the retry-loop this catalog exists to prevent.
type IntentOutcome string

const (
	// IntentOutcomeApproved: the user approved; the action ran. Result present.
	IntentOutcomeApproved IntentOutcome = "approved"
	// IntentOutcomeAutoApproved: nobody answered in time and the tool's policy
	// is wait-then-allow; the action ran. Result present.
	IntentOutcomeAutoApproved IntentOutcome = "auto_approved"
	// IntentOutcomeDeniedByUser: the user denied; the action did NOT run.
	IntentOutcomeDeniedByUser IntentOutcome = "denied_by_user"
	// IntentOutcomeAutoDenied: nobody answered in time and the tool's policy is
	// wait-then-deny; the action did NOT run.
	IntentOutcomeAutoDenied IntentOutcome = "auto_denied"
	// IntentOutcomeError: the intent failed before completing. The action may or
	// may not have taken effect.
	IntentOutcomeError IntentOutcome = "error"
)

func (o IntentOutcome) Valid() bool {
	switch o {
	case IntentOutcomeApproved, IntentOutcomeAutoApproved,
		IntentOutcomeDeniedByUser, IntentOutcomeAutoDenied, IntentOutcomeError:
		return true
	default:
		return false
	}
}

// Approved reports whether the action actually ran and a result is present.
func (o IntentOutcome) Approved() bool {
	return o == IntentOutcomeApproved || o == IntentOutcomeAutoApproved
}

// Retryable reports whether an agent may reasonably retry the call. Only error
// qualifies: a denial is a decision, not a failure.
func (o IntentOutcome) Retryable() bool {
	return o == IntentOutcomeError
}

// IntentPolicy is the tool's hardcoded behavior when the wait window expires.
// It is carried on the wire so the desktop can label the countdown ("will
// auto-approve" vs "will auto-deny") without knowing the tool.
type IntentPolicy string

const (
	// IntentPolicyAutoAllow: never ask; run immediately.
	IntentPolicyAutoAllow IntentPolicy = "auto-allow"
	// IntentPolicyWaitThenAllow: ask, and run anyway if nobody answers.
	IntentPolicyWaitThenAllow IntentPolicy = "wait-then-allow"
	// IntentPolicyWaitThenDeny: ask, and deny if nobody answers.
	IntentPolicyWaitThenDeny IntentPolicy = "wait-then-deny"
)

func (p IntentPolicy) Valid() bool {
	switch p {
	case IntentPolicyAutoAllow, IntentPolicyWaitThenAllow, IntentPolicyWaitThenDeny:
		return true
	default:
		return false
	}
}

// IntentOrigin identifies where an intent came from. It is what lets the
// desktop render one popup for any tool across any session — the card can say
// "architect hiveryn / ticket X" without knowing what the tool does. Derived
// from the Session record daemon-side; never supplied by the agent.
type IntentOrigin struct {
	ArchitectKey string      `json:"architect_key"`
	SessionID    string      `json:"session_id"`
	SessionType  SessionType `json:"session_type"`
	TicketID     string      `json:"ticket_id,omitempty"` // Session.ContextID for ticket sessions
}

// Intent is the desktop-facing description of a pending, approvable tool call.
// Payload is tool-specific and rendered generically by the desktop.
type Intent struct {
	ID          string         `json:"intent_id"`
	Type        IntentType     `json:"intent_type"`
	Summary     string         `json:"summary"`
	Payload     map[string]any `json:"payload,omitempty"`
	Origin      IntentOrigin   `json:"origin"`
	WaitSeconds int            `json:"wait_seconds"`
	Policy      IntentPolicy   `json:"policy"`
	CreatedAt   time.Time      `json:"created_at"`
}
