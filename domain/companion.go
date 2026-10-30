package domain

import (
	"fmt"
	"strings"
	"time"
	"unicode/utf8"
)

// Companions.
//
// A companion is a persistent agent session with its own repository, started
// by the daemon with each new architect session for every entry of the
// architect's hiveryn.yaml `companions`. It belongs to that exact architect
// session (Session.OwnerSessionID) until it is stopped; configuration edits
// apply to later architect sessions only. The architect and its companions
// exchange durable asynchronous messages through the agent-message outbox,
// the companion asks the user (never the architect) to authorize external
// actions with requestAction, and a companion ends only through a cooperative
// shutdown: the architect's request_id, the companion's handoff and the
// daemon's confirmed termination. An architect session cannot conclude while
// any companion it owns is still not stopped.

// SessionTypeCompanion is a companion's session. Its ContextID is the
// configured companion name and OwnerSessionID the owning architect session.
const SessionTypeCompanion SessionType = "companion"

// SessionCreatedByCompanionStartup marks a companion session started by the
// daemon alongside its architect session.
const SessionCreatedByCompanionStartup SessionCreatedBy = "companion_startup"

// CompanionStatus is a companion's lifecycle within its architect session.
type CompanionStatus string

const (
	// CompanionStarting: the session is being created and launched.
	CompanionStarting CompanionStatus = "starting"
	// CompanionRunning: launched and accepting assignments.
	CompanionRunning CompanionStatus = "running"
	// CompanionShutdownRequested: the architect requested shutdown; the
	// instruction is queued and the companion accepts no new assignments.
	CompanionShutdownRequested CompanionStatus = "shutdown_requested"
	// CompanionStopping: the companion completed its handoff; the daemon is
	// terminating it.
	CompanionStopping CompanionStatus = "stopping"
	// CompanionStopped: termination was confirmed.
	CompanionStopped CompanionStatus = "stopped"
	// CompanionStopFailed: termination could not be confirmed (Error says
	// why). It stays visible and blocks the architect's conclusion; the
	// user's explicit terminate of the session is the remaining route.
	CompanionStopFailed CompanionStatus = "stop_failed"
	// CompanionStartFailed: the companion never became available (Error says
	// why); nothing of it remains running.
	CompanionStartFailed CompanionStatus = "start_failed"
	// CompanionFailed: a local companion's process ended and could not be
	// resumed; nothing of it remains running.
	CompanionFailed CompanionStatus = "failed"
)

// Settled reports whether nothing of the companion can still be running, so
// it no longer blocks its architect's conclusion.
func (s CompanionStatus) Settled() bool {
	switch s {
	case CompanionStopped, CompanionStartFailed, CompanionFailed:
		return true
	default:
		return false
	}
}

// Valid reports whether s is a known status.
func (s CompanionStatus) Valid() bool {
	switch s {
	case CompanionStarting, CompanionRunning, CompanionShutdownRequested, CompanionStopping,
		CompanionStopped, CompanionStopFailed, CompanionStartFailed, CompanionFailed:
		return true
	default:
		return false
	}
}

// Companion is one configured companion of one architect session: the
// configuration snapshot taken when that session started, its lifecycle and,
// once launched, its session. The record lives as long as the architect
// session.
//
// The trailing fields are hydrated from the live session when read
// (main terminal, SSH attachment state, agent status) and from the pending
// action requests; they are not stored.
type Companion struct {
	ArchitectKey       string          `json:"architect_key"`
	ArchitectSessionID string          `json:"architect_session_id"`
	Name               string          `json:"name"`
	Repo               string          `json:"repo"`
	Workdir            string          `json:"workdir"`
	Variant            string          `json:"variant"`
	Agent              string          `json:"agent,omitempty"`
	Machine            string          `json:"machine,omitempty"`
	SessionID          string          `json:"session_id,omitempty"`
	Status             CompanionStatus `json:"status"`
	Error              string          `json:"error,omitempty"`
	CreatedAt          time.Time       `json:"created_at"`
	UpdatedAt          time.Time       `json:"updated_at"`

	Shutdown       *CompanionShutdown `json:"shutdown,omitempty"`
	MainTerminalID string             `json:"main_terminal_id,omitempty"`
	Connection     string             `json:"connection,omitempty"`
	AgentStatus    string             `json:"agent_status,omitempty"`
	PendingActions []ActionRequest    `json:"pending_actions"`
}

// CompanionList is GET /api/architects/{key}/companions: the companions of the
// architect's running session, in name order (empty without one).
type CompanionList struct {
	ArchitectSessionID string      `json:"architect_session_id,omitempty"`
	Companions         []Companion `json:"companions"`
}

// CompanionShutdownStatus is the lifecycle of one shutdown request.
type CompanionShutdownStatus string

const (
	// ShutdownRequested: recorded and the instruction queued to the companion.
	ShutdownRequested CompanionShutdownStatus = "requested"
	// ShutdownCompleting: the companion's handoff is recorded; termination runs.
	ShutdownCompleting CompanionShutdownStatus = "completing"
	// ShutdownStopped: termination confirmed; the architect was notified.
	ShutdownStopped CompanionShutdownStatus = "stopped"
	// ShutdownStopFailed: termination could not be confirmed (Error).
	ShutdownStopFailed CompanionShutdownStatus = "stop_failed"
)

// CompanionShutdown is one architect-issued shutdown request, bound to the
// exact architect and companion sessions it was issued for.
type CompanionShutdown struct {
	RequestID          string                  `json:"request_id"`
	ArchitectKey       string                  `json:"architect_key"`
	ArchitectSessionID string                  `json:"architect_session_id"`
	CompanionSessionID string                  `json:"companion_session_id"`
	Companion          string                  `json:"companion"`
	Status             CompanionShutdownStatus `json:"status"`
	Handoff            string                  `json:"handoff,omitempty"`
	Error              string                  `json:"error,omitempty"`
	RequestedAt        time.Time               `json:"requested_at"`
	CompletedAt        *time.Time              `json:"completed_at,omitempty"`
	StoppedAt          *time.Time              `json:"stopped_at,omitempty"`
}

const (
	// MaxCompanionMessageLength bounds one architect/companion message.
	MaxCompanionMessageLength = 100000
	// MaxCompanionHandoffLength bounds a shutdown handoff.
	MaxCompanionHandoffLength = 100000
	// CompanionAccepted is the status of every accepted companion tool call.
	CompanionAccepted = "accepted"
)

// CompanionMessageRequest is sendMessageToArchitect (companion) or
// sendMessageToCompanion (architect; the companion name travels in the path).
type CompanionMessageRequest struct {
	Message string `json:"message"`
}

// CompanionReplyRequest is replyToArchitectMessage (companion) or
// replyToCompanionMessage (architect). MessageID is the message answered.
type CompanionReplyRequest struct {
	MessageID string `json:"message_id"`
	Message   string `json:"message"`
}

// CompanionMessageResult: the message is durably queued (MessageID is the
// queued message); InReplyTo is set for replies.
type CompanionMessageResult struct {
	MessageID string `json:"message_id"`
	InReplyTo string `json:"in_reply_to,omitempty"`
	Status    string `json:"status"`
}

// CompanionShutdownRequest is requestCompanionShutdown (architect).
type CompanionShutdownRequest struct {
	Companion string `json:"companion"`
	RequestID string `json:"request_id"`
}

// CompleteCompanionShutdownRequest is completeCompanionShutdown (companion).
type CompleteCompanionShutdownRequest struct {
	RequestID string `json:"request_id"`
	Handoff   string `json:"handoff"`
}

// CompanionShutdownResult acknowledges either shutdown call.
type CompanionShutdownResult struct {
	RequestID string `json:"request_id"`
	Status    string `json:"status"`
}

// NormalizeCompanionText trims a message or handoff and checks its bounds.
func NormalizeCompanionText(field, text string, max int) (string, error) {
	text = strings.TrimSpace(text)
	if text == "" {
		return "", &ValidationError{Field: field, Message: "is required"}
	}
	if n := utf8.RuneCountInString(text); n > max {
		return "", &ValidationError{Field: field, Message: fmt.Sprintf("must be at most %d characters, got %d", max, n)}
	}
	return text, nil
}

// Action requests (requestAction).
//
// A companion asks the user to authorize one concrete external action. The
// user's phone is alerted and the companion's call waits for a dialog
// decision — accept, deny, or reply with a message — for at most
// ActionRequestTimeout. Only an accepted outcome authorizes the proposal, and
// only the proposal as submitted; Hiveryn executes nothing. There is no
// automatic acceptance and no escalation to the architect. The request id is
// internal: it addresses the desktop's decision and never reaches the agent.
//
// The request travels on the companion session's event stream as
// `action_request` `required` (raw: action_id, title, description, origin,
// created_at, expires_at) and `resolved` (raw: action_id, status, message,
// reason, decided_by).

const (
	MaxActionTitleLength       = 200
	MaxActionDescriptionLength = 50000
	MaxActionReplyLength       = 4000
	ActionRequestTimeout       = time.Hour

	// SessionEventTypeActionRequest is the session event type of requests.
	SessionEventTypeActionRequest = "action_request"

	// ActionRequestTimeoutCode / ActionRequestTimeoutText: the user did not
	// decide within ActionRequestTimeout. Returned as a tool error.
	ActionRequestTimeoutCode = "ACTION_REQUEST_TIMEOUT"
	ActionRequestTimeoutText = "The user did not respond within one hour. No action is authorized. Call requestAction again with the proposal to submit another request and wait for a decision."
	// ActionRequestCancelledCode: the request was withdrawn before a decision
	// (shutdown requested, session ended, the call ended). A tool error.
	ActionRequestCancelledCode = "ACTION_REQUEST_CANCELLED"
	// ActionRequestInterruptedCode: the daemon stopped while it was pending.
	ActionRequestInterruptedCode = "ACTION_REQUEST_INTERRUPTED"
	// ActionRequestNotificationFailedCode: the phone alert failed, so no
	// request was raised.
	ActionRequestNotificationFailedCode = "ACTION_REQUEST_NOTIFICATION_FAILED"
)

// ActionRequestStatus is a request's lifecycle; it resolves exactly once.
type ActionRequestStatus string

const (
	ActionPending     ActionRequestStatus = "pending"
	ActionAccepted    ActionRequestStatus = "accepted"
	ActionDenied      ActionRequestStatus = "denied"
	ActionReplied     ActionRequestStatus = "replied"
	ActionExpired     ActionRequestStatus = "expired"
	ActionCancelled   ActionRequestStatus = "cancelled"
	ActionInterrupted ActionRequestStatus = "interrupted"
)

// ActionOrigin is derived from the requesting session.
type ActionOrigin struct {
	ArchitectKey       string      `json:"architect_key"`
	ArchitectSessionID string      `json:"architect_session_id,omitempty"`
	SessionID          string      `json:"session_id"`
	SessionType        SessionType `json:"session_type"`
	Companion          string      `json:"companion,omitempty"`
}

// ActionRequest is one request as shown in the decision dialog.
type ActionRequest struct {
	ID          string              `json:"id"`
	Title       string              `json:"title"`
	Description string              `json:"description"`
	Origin      ActionOrigin        `json:"origin"`
	Status      ActionRequestStatus `json:"status"`
	// Message is the user's reply (status replied).
	Message string `json:"message,omitempty"`
	// Reason explains a request that ended without a decision.
	Reason    string     `json:"reason,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	ExpiresAt time.Time  `json:"expires_at"`
	DecidedAt *time.Time `json:"decided_at,omitempty"`
}

// RequestActionRequest is the companion's requestAction call.
type RequestActionRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

// Normalize validates and trims the request.
func (r RequestActionRequest) Normalize() (RequestActionRequest, error) {
	title, err := NormalizeCompanionText("title", r.Title, MaxActionTitleLength)
	if err != nil {
		return RequestActionRequest{}, err
	}
	description, err := NormalizeCompanionText("description", r.Description, MaxActionDescriptionLength)
	if err != nil {
		return RequestActionRequest{}, err
	}
	return RequestActionRequest{Title: title, Description: description}, nil
}

// ActionDecision is the user's choice in the dialog.
type ActionDecision string

const (
	ActionDecisionAccept ActionDecision = "accept"
	ActionDecisionDeny   ActionDecision = "deny"
	ActionDecisionReply  ActionDecision = "reply"
)

// ActionDecisionRequest is POST /api/sessions/{id}/actions/{actionID}/decision.
// Message is required for reply and refused otherwise.
type ActionDecisionRequest struct {
	Decision ActionDecision `json:"decision"`
	Message  string         `json:"message,omitempty"`
}

// Normalize validates the decision.
func (r ActionDecisionRequest) Normalize() (ActionDecisionRequest, error) {
	switch r.Decision {
	case ActionDecisionAccept, ActionDecisionDeny:
		if strings.TrimSpace(r.Message) != "" {
			return ActionDecisionRequest{}, &ValidationError{Field: "message", Message: "is only accepted with decision reply"}
		}
		return ActionDecisionRequest{Decision: r.Decision}, nil
	case ActionDecisionReply:
		message, err := NormalizeCompanionText("message", r.Message, MaxActionReplyLength)
		if err != nil {
			return ActionDecisionRequest{}, err
		}
		return ActionDecisionRequest{Decision: r.Decision, Message: message}, nil
	default:
		return ActionDecisionRequest{}, &ValidationError{Field: "decision", Message: "must be accept, deny or reply"}
	}
}

// RequestActionResult is the daemon's answer to the companion's blocking
// call: an Outcome (accepted, denied, replied with Message) or, when the
// request ended without a decision, ErrorCode and Error, which the MCP tool
// returns as a tool error. Never both.
type RequestActionResult struct {
	Outcome   string `json:"outcome,omitempty"`
	Message   string `json:"message,omitempty"`
	ErrorCode string `json:"error_code,omitempty"`
	Error     string `json:"error,omitempty"`
}
