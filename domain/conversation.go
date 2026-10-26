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
	Limitation        string `json:"limitation,omitempty"`
}
type SubmissionRequest struct {
	RequestID string `json:"request_id"`
	Text      string `json:"text"`
}
type Submission struct {
	ID           string    `json:"id"`
	SessionID    string    `json:"session_id"`
	RequestID    string    `json:"request_id"`
	Sequence     int64     `json:"sequence"`
	Origin       string    `json:"origin"`
	Text         string    `json:"text"`
	Status       string    `json:"status"`
	NativeTurnID string    `json:"native_turn_id,omitempty"`
	Error        string    `json:"error,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
type ConversationItem struct {
	ID           string    `json:"id"`
	SubmissionID string    `json:"submission_id"`
	Sequence     int64     `json:"sequence"`
	Revision     int64     `json:"revision"`
	NativeTurnID string    `json:"native_turn_id,omitempty"`
	Role         string    `json:"role"`
	Origin       string    `json:"origin"`
	Text         string    `json:"text,omitempty"`
	Tool         string    `json:"tool,omitempty"`
	Status       string    `json:"status"`
	At           time.Time `json:"at"`
}

// Items are ordered by revision. A page is bounded by SnapshotRevision;
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
