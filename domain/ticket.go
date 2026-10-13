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

type TicketReferenceType string
type PathReferenceKind string

// PathReferenceStatus is the outcome of inspecting a path reference where it
// lives: present, confirmed missing, or unverified (it could not be
// inspected; Error says why).
type PathReferenceStatus string

// PathReferenceLocation is where a path reference was interpreted: on the
// daemon's machine, on a remote machine (Machine names it), or unknown when
// the ticket's primary-repo machine could not be determined. Unknown is never
// local.
type PathReferenceLocation string

const (
	TicketReferenceTicket  TicketReferenceType = "ticket"
	TicketReferencePath    TicketReferenceType = "path"
	PathReferenceFile      PathReferenceKind   = "file"
	PathReferenceDirectory PathReferenceKind   = "directory"

	PathReferencePresent    PathReferenceStatus = "present"
	PathReferenceMissing    PathReferenceStatus = "missing"
	PathReferenceUnverified PathReferenceStatus = "unverified"

	PathReferenceLocal   PathReferenceLocation = "local"
	PathReferenceRemote  PathReferenceLocation = "remote"
	PathReferenceUnknown PathReferenceLocation = "unknown"
)

// TicketReference is the daemon-resolved view of an entry in References.
// Value remains the portable string serialized in ticket.md. Path references
// always carry Status and Location (Exists mirrors Status == present); a
// remote one names its Machine and an unverified one its Error. Ticket-ID
// references carry only Exists.
type TicketReference struct {
	Value    string                `json:"value"`
	Type     TicketReferenceType   `json:"type"`
	Exists   bool                  `json:"exists"`
	Kind     PathReferenceKind     `json:"kind,omitempty"`
	Status   PathReferenceStatus   `json:"status,omitempty"`
	Location PathReferenceLocation `json:"location,omitempty"`
	Machine  string                `json:"machine,omitempty"`
	Error    string                `json:"error,omitempty"`
}

type TicketSummary struct {
	ID     string       `json:"id"`
	Status TicketStatus `json:"status"`
	Title  string       `json:"title"`
	Repo   string       `json:"repo,omitempty"`
	// AdditionalRepos are extra repositories in scope beyond the primary Repo.
	// Order carries no meaning; keys must be unique and must not include Repo.
	AdditionalRepos    []string          `json:"additional_repos"`
	Created            *time.Time        `json:"created,omitempty"`
	Updated            *time.Time        `json:"updated,omitempty"`
	References         []string          `json:"references"`
	ResolvedReferences []TicketReference `json:"resolved_references"`
	HasConclusion      bool              `json:"has_conclusion"`
	Warnings           []TicketWarning   `json:"warnings"`
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
