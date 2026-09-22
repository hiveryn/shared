package domain

import "time"

// ArchitectEventType is the only type carried on the architect event stream
// (GET /api/architects/{key}/events). It exists as a discriminator so clients
// can ignore anything they synthesize locally onto the same channel.
const ArchitectEventType = "workspace_changed"

// ArchitectEventReason says what actually changed. Ticket reasons drive the
// kanban board; session reasons drive the set of session tabs in an architect
// window. A client that only knows about the ticket reasons keeps working.
type ArchitectEventReason string

const (
	ArchitectEventTicketCreated   ArchitectEventReason = "ticket_created"
	ArchitectEventTicketUpdated   ArchitectEventReason = "ticket_updated"
	ArchitectEventTicketMoved     ArchitectEventReason = "ticket_moved"
	ArchitectEventTicketDeleted   ArchitectEventReason = "ticket_deleted"
	ArchitectEventTicketConcluded ArchitectEventReason = "ticket_concluded"
	// ArchitectEventSessionStarted announces a session whose run is live: it has
	// a running run and a main terminal. It is the only signal that names a
	// session id, and therefore the only way a client can discover a session it
	// did not create itself (a spawn from another window).
	ArchitectEventSessionStarted ArchitectEventReason = "session_started"
	// ArchitectEventSessionEnded announces a session that is gone (concluded or
	// discarded), so a client that missed the session-scoped lifecycle event can
	// still drop it.
	ArchitectEventSessionEnded ArchitectEventReason = "session_ended"
)

func (r ArchitectEventReason) Valid() bool {
	switch r {
	case ArchitectEventTicketCreated,
		ArchitectEventTicketUpdated,
		ArchitectEventTicketMoved,
		ArchitectEventTicketDeleted,
		ArchitectEventTicketConcluded,
		ArchitectEventSessionStarted,
		ArchitectEventSessionEnded:
		return true
	default:
		return false
	}
}

// ArchitectEvent is published per architect key and fans out to every subscriber
// of that architect's event stream. The stream has no backlog, so a client must
// reconcile on (re)connect rather than treat delivery as guaranteed.
type ArchitectEvent struct {
	Type         string               `json:"type"`
	ArchitectKey string               `json:"architect_key"`
	Reason       ArchitectEventReason `json:"reason"`
	// TicketID is empty for reasons that are not about a ticket.
	TicketID string `json:"ticket_id"`
	// SessionID is set for the session reasons, and for ticket reasons caused by
	// a session (conclude, discard). Empty for plain ticket edits.
	SessionID string    `json:"session_id"`
	At        time.Time `json:"at"`
}
