package domain

import "time"

type TicketStatus string

const (
	TicketStatusBacklog  TicketStatus = "backlog"
	TicketStatusProgress TicketStatus = "progress"
	TicketStatusDone     TicketStatus = "done"
)

func (s TicketStatus) Valid() bool {
	switch s {
	case TicketStatusBacklog, TicketStatusProgress, TicketStatusDone:
		return true
	default:
		return false
	}
}

type TicketWarning struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type TicketSummary struct {
	ID     string       `json:"id"`
	Status TicketStatus `json:"status"`
	Title  string       `json:"title"`
	Repo   string       `json:"repo,omitempty"`
	// AdditionalRepos are extra repositories in scope beyond the primary Repo.
	// Order carries no meaning; keys must be unique and must not include Repo.
	AdditionalRepos []string        `json:"additional_repos"`
	Created         *time.Time      `json:"created,omitempty"`
	Updated         *time.Time      `json:"updated,omitempty"`
	References      []string        `json:"references"`
	HasConclusion   bool            `json:"has_conclusion"`
	Warnings        []TicketWarning `json:"warnings"`
}

// TicketOutcome distinguishes how a ticket conclusion ended: shipped commits,
// investigation/spike work that produced no commits, or an outright
// rejection.
type TicketOutcome string

const (
	TicketOutcomeCompleted   TicketOutcome = "completed"
	TicketOutcomeExploratory TicketOutcome = "exploratory"
	TicketOutcomeRejected    TicketOutcome = "rejected"
)

func (o TicketOutcome) Valid() bool {
	switch o {
	case TicketOutcomeCompleted, TicketOutcomeExploratory, TicketOutcomeRejected:
		return true
	default:
		return false
	}
}

type TicketConclusion struct {
	StartedAt       time.Time     `json:"started_at"`
	ConcludedAt     time.Time     `json:"concluded_at"`
	Agent           string        `json:"agent,omitempty"`
	Profile         string        `json:"profile,omitempty"`
	Outcome         TicketOutcome `json:"outcome"`
	RejectionReason string        `json:"rejection_reason,omitempty"`
	Commits         []CommitRef   `json:"commits"`
	Body            string        `json:"body"`
}

type Ticket struct {
	TicketSummary
	Body       string            `json:"body"`
	Conclusion *TicketConclusion `json:"conclusion"`
}

type TicketBoard struct {
	Backlog  []TicketSummary `json:"backlog"`
	Progress []TicketSummary `json:"progress"`
	Done     []TicketSummary `json:"done"`
}

type CreateTicketParams struct {
	Title           string
	Repo            string
	AdditionalRepos []string
	Body            string
	References      []string
	Now             time.Time
}

type EditTicketParams struct {
	OldString  string
	NewString  string
	ReplaceAll bool
	Now        time.Time
}

type UpdateTicketMetadataParams struct {
	Title           *string
	Repo            *string
	AdditionalRepos *[]string
	References      *[]string
	Now             time.Time
}

type MoveTicketParams struct {
	To TicketStatus
}
