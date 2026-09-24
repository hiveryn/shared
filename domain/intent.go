package domain

import "time"

// An Intent is a pending agent tool call awaiting the user's approval. The
// desktop renders it as a popup and approves or denies it. An intent is either
// blocking or deferred, by its tool's hardcoded policy:
//
//   - Blocking (auto-allow, wait-then-allow, wait-then-deny): the agent's MCP
//     call blocks in the daemon while the intent is pending; if nobody answers
//     within the wait window, the policy resolves it.
//   - Deferred (manual): the request returns at once with a stable intent id and
//     status pending_approval, and only the user resolves it — there is no timer.
//     Its pending or resolved outcome is a DeferredIntent, retrieved by that id.
//     Every intent that asks for approval inputs is deferred.
//
// This is deliberately distinct from Session (the spawn record): a Session is
// "run an agent", an Intent is "the agent wants to do something first".

type IntentType string

const (
	IntentTypeConcludeSession  IntentType = "concludeSession"
	IntentTypeCreateWorkTicket IntentType = "createWorkTicket"
	// IntentTypeExecuteAction is an architect's request to run an Action. It is
	// deferred: the intent id is also the Action execution id.
	IntentTypeExecuteAction IntentType = "executeAction"
)

func (t IntentType) Valid() bool {
	switch t {
	case IntentTypeConcludeSession, IntentTypeCreateWorkTicket, IntentTypeExecuteAction:
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
	// IntentPolicyManual: deferred. Ask, return immediately, and wait for the
	// user indefinitely; nothing resolves it automatically, and defaults never
	// stand in for an answer.
	IntentPolicyManual IntentPolicy = "manual"
)

func (p IntentPolicy) Valid() bool {
	switch p {
	case IntentPolicyAutoAllow, IntentPolicyWaitThenAllow, IntentPolicyWaitThenDeny, IntentPolicyManual:
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
//
// Inputs, when present, are fields the user completes as part of approving; the
// approved values reach the operation. An intent with inputs is always deferred
// (Policy manual), so it waits for the user however its defaults look.
// WaitSeconds is 0 for a deferred intent: there is no countdown.
type Intent struct {
	ID          string             `json:"intent_id"`
	Type        IntentType         `json:"intent_type"`
	Summary     string             `json:"summary"`
	Payload     map[string]any     `json:"payload,omitempty"`
	Inputs      []IntentInputField `json:"inputs,omitempty"`
	Origin      IntentOrigin       `json:"origin"`
	WaitSeconds int                `json:"wait_seconds"`
	Policy      IntentPolicy       `json:"policy"`
	CreatedAt   time.Time          `json:"created_at"`
}

// IntentInputType is the closed set of approval input kinds. The schema is
// deliberately small — ordinary form fields plus a required choice from
// supplied options — not a general form engine.
//
// Value types: text, textarea and choice carry a string; boolean carries a bool.
type IntentInputType string

const (
	// IntentInputText: a single-line string.
	IntentInputText IntentInputType = "text"
	// IntentInputTextarea: a multi-line string.
	IntentInputTextarea IntentInputType = "textarea"
	// IntentInputChoice: one Value from Options.
	IntentInputChoice IntentInputType = "choice"
	// IntentInputBoolean: true or false.
	IntentInputBoolean IntentInputType = "boolean"
)

func (t IntentInputType) Valid() bool {
	switch t {
	case IntentInputText, IntentInputTextarea, IntentInputChoice, IntentInputBoolean:
		return true
	default:
		return false
	}
}

// Bounds on a schema, so a pending approval stays a short form.
const (
	MaxIntentInputFields  = 16
	MaxIntentInputOptions = 200
)

// IntentInputOption is one allowed value of a choice input.
type IntentInputOption struct {
	Value       string `json:"value"`
	Label       string `json:"label,omitempty"` // display text; Value when empty
	Description string `json:"description,omitempty"`
}

// IntentInputField is one approval input. Default, when set, has the field's
// value type and only prefills the form: the daemon never approves with it or
// substitutes it for a missing submission, so a stale default costs the user a
// correction, never a wrong run.
// MaxLength (in characters) applies to text and textarea only; 0 is unbounded.
type IntentInputField struct {
	Name        string              `json:"name"`
	Label       string              `json:"label"`
	Description string              `json:"description,omitempty"`
	Type        IntentInputType     `json:"type"`
	Required    bool                `json:"required,omitempty"`
	Default     any                 `json:"default,omitempty"`
	Options     []IntentInputOption `json:"options,omitempty"`
	MaxLength   int                 `json:"max_length,omitempty"`
}

// IntentInputValues are submitted or resolved input values keyed by field name.
type IntentInputValues map[string]any

// IntentInputIssue is one problem with an input value, keyed by field name.
type IntentInputIssue struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ApproveIntentRequest is the desktop's approve body. Inputs may be omitted
// when the intent has none. Values are taken as submitted; defaults are not
// applied, so a required field must be sent.
type ApproveIntentRequest struct {
	Inputs IntentInputValues `json:"inputs,omitempty"`
}

// DeferredIntentStatus is the lifecycle of a deferred intent. Approval and
// execution are separate facts: running means the user approved and the
// operation started; only completed means it succeeded.
//
// The set deliberately matches the agreed Actions status model, so a future
// Action can keep the one caller-visible id from request to run completion.
type DeferredIntentStatus string

const (
	// DeferredIntentPendingApproval: waiting for the user. Nothing has run.
	DeferredIntentPendingApproval DeferredIntentStatus = "pending_approval"
	// DeferredIntentDenied: the user denied it; the operation never ran.
	// Reason carries the user's reason, if any.
	DeferredIntentDenied DeferredIntentStatus = "denied"
	// DeferredIntentRunning: approved; the operation is executing.
	DeferredIntentRunning DeferredIntentStatus = "running"
	// DeferredIntentCompleted: approved and the operation succeeded. Result
	// carries its result.
	DeferredIntentCompleted DeferredIntentStatus = "completed"
	// DeferredIntentFailed: it did not complete. Error says why. ApprovedAt
	// tells the two cases apart: unset means it was never approved and never
	// ran (for example the session ended or the daemon restarted first); set
	// means the operation started and failed or was interrupted, and may or may
	// not have taken effect.
	DeferredIntentFailed DeferredIntentStatus = "failed"
)

func (s DeferredIntentStatus) Valid() bool {
	switch s {
	case DeferredIntentPendingApproval, DeferredIntentDenied, DeferredIntentRunning,
		DeferredIntentCompleted, DeferredIntentFailed:
		return true
	default:
		return false
	}
}

// Terminal reports whether the status is final.
func (s DeferredIntentStatus) Terminal() bool {
	return s == DeferredIntentDenied || s == DeferredIntentCompleted || s == DeferredIntentFailed
}

// DeferredIntent is the durable, id-addressable outcome of a deferred intent:
// what the request returns immediately (status pending_approval) and what a
// lookup by ID returns at any later point, including after the session ended.
// It is scoped to its origin session. Inputs are the validated values the
// operation ran with, set once approved.
type DeferredIntent struct {
	ID         string               `json:"intent_id"`
	Type       IntentType           `json:"intent_type"`
	Summary    string               `json:"summary"`
	Payload    map[string]any       `json:"payload,omitempty"`
	Origin     IntentOrigin         `json:"origin"`
	Status     DeferredIntentStatus `json:"status"`
	Inputs     IntentInputValues    `json:"inputs,omitempty"`
	Result     any                  `json:"result,omitempty"`
	Reason     string               `json:"reason,omitempty"` // denial reason
	Error      string               `json:"error,omitempty"`  // failure detail
	CreatedAt  time.Time            `json:"created_at"`
	ApprovedAt *time.Time           `json:"approved_at,omitempty"`
	EndedAt    *time.Time           `json:"ended_at,omitempty"`
}
