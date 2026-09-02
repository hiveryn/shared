package domain

import "time"

// RoadmapSchemaVersion is the current on-disk schema version for both
// roadmap/current.yaml and roadmap/archive.yaml.
const RoadmapSchemaVersion = 1

type RoadmapItemKind string

const (
	RoadmapItemGoal       RoadmapItemKind = "goal"
	RoadmapItemInitiative RoadmapItemKind = "initiative"
	RoadmapItemMilestone  RoadmapItemKind = "milestone"
)

func (k RoadmapItemKind) Valid() bool {
	switch k {
	case RoadmapItemGoal, RoadmapItemInitiative, RoadmapItemMilestone:
		return true
	default:
		return false
	}
}

type RoadmapItemStatus string

const (
	RoadmapStatusPlanned RoadmapItemStatus = "planned"
	RoadmapStatusActive  RoadmapItemStatus = "active"
	RoadmapStatusBlocked RoadmapItemStatus = "blocked"
	RoadmapStatusDone    RoadmapItemStatus = "done"
)

func (s RoadmapItemStatus) Valid() bool {
	switch s {
	case RoadmapStatusPlanned, RoadmapStatusActive, RoadmapStatusBlocked, RoadmapStatusDone:
		return true
	default:
		return false
	}
}

// RoadmapItem is one node in the flat roadmap item graph. Hierarchy is
// expressed through ParentID (nil = top-level), sibling sequence through
// Order. The same shape is persisted as YAML and served over the wire.
type RoadmapItem struct {
	ID              string            `json:"id" yaml:"id" jsonschema:"Stable human-readable kebab-case identifier, unique across current and archived items"`
	Kind            RoadmapItemKind   `json:"kind" yaml:"kind" jsonschema:"goal, initiative, or milestone"`
	Title           string            `json:"title" yaml:"title" jsonschema:"Non-empty item title"`
	Status          RoadmapItemStatus `json:"status" yaml:"status" jsonschema:"planned, active, blocked, or done"`
	Outcome         string            `json:"outcome" yaml:"outcome" jsonschema:"Plain-text intended result, not an implementation checklist"`
	ParentID        *string           `json:"parent_id" yaml:"parent_id" jsonschema:"ID of the parent item; null for top-level items"`
	Order           int               `json:"order" yaml:"order" jsonschema:"Integer used for deterministic sibling ordering"`
	SuccessCriteria []string          `json:"success_criteria" yaml:"success_criteria" jsonschema:"Ordered plain-text success criteria"`
	Tickets         []string          `json:"tickets" yaml:"tickets" jsonschema:"Ordered, deduplicated ticket IDs on the same architect board"`
	DependsOn       []string          `json:"depends_on" yaml:"depends_on" jsonschema:"Deduplicated roadmap item IDs this item depends on"`
}

// Roadmap is the full contents of roadmap/current.yaml.
type Roadmap struct {
	SchemaVersion int           `json:"schema_version" yaml:"schema_version"`
	Title         string        `json:"title" yaml:"title"`
	Items         []RoadmapItem `json:"items" yaml:"items"`
}

// RoadmapArchiveEntry is one archived subtree: the root item plus the
// complete flat snapshot of its descendants. The root's parent_id is null
// inside Items; OriginalParentID/OriginalOrder preserve where the subtree
// was attached so restore can put it back.
type RoadmapArchiveEntry struct {
	RootID           string        `json:"root_id" yaml:"root_id"`
	ArchivedAt       time.Time     `json:"archived_at" yaml:"archived_at"`
	Summary          string        `json:"summary,omitempty" yaml:"summary,omitempty"`
	OriginalParentID *string       `json:"original_parent_id" yaml:"original_parent_id"`
	OriginalOrder    int           `json:"original_order" yaml:"original_order"`
	Items            []RoadmapItem `json:"items" yaml:"items"`
}

// RoadmapArchive is the full contents of roadmap/archive.yaml.
type RoadmapArchive struct {
	SchemaVersion int                   `json:"schema_version" yaml:"schema_version"`
	Entries       []RoadmapArchiveEntry `json:"entries" yaml:"entries"`
}

// RoadmapArchiveEntrySummary is the compact listing shape for archive reads
// without a specific entry ID.
type RoadmapArchiveEntrySummary struct {
	RootID     string          `json:"root_id"`
	RootTitle  string          `json:"root_title"`
	RootKind   RoadmapItemKind `json:"root_kind"`
	ArchivedAt time.Time       `json:"archived_at"`
	Summary    string          `json:"summary,omitempty"`
	ItemCount  int             `json:"item_count"`
}

// RoadmapTicketInfo is compact, decision-relevant metadata for a ticket
// linked from a roadmap item, resolved against the same architect board.
type RoadmapTicketInfo struct {
	ID                string        `json:"id"`
	Title             string        `json:"title"`
	Status            TicketStatus  `json:"status"`
	Repo              string        `json:"repo,omitempty"`
	AdditionalRepos   []string      `json:"additional_repos"`
	HasConclusion     bool          `json:"has_conclusion"`
	ConclusionOutcome TicketOutcome `json:"conclusion_outcome,omitempty"`
}

type RoadmapOpType string

const (
	RoadmapOpCreate       RoadmapOpType = "create"
	RoadmapOpUpdate       RoadmapOpType = "update"
	RoadmapOpMove         RoadmapOpType = "move"
	RoadmapOpLinkTicket   RoadmapOpType = "link_ticket"
	RoadmapOpUnlinkTicket RoadmapOpType = "unlink_ticket"
	RoadmapOpArchive      RoadmapOpType = "archive"
	RoadmapOpRestore      RoadmapOpType = "restore"
)

func (t RoadmapOpType) Valid() bool {
	switch t {
	case RoadmapOpCreate, RoadmapOpUpdate, RoadmapOpMove,
		RoadmapOpLinkTicket, RoadmapOpUnlinkTicket,
		RoadmapOpArchive, RoadmapOpRestore:
		return true
	default:
		return false
	}
}

// RoadmapOp is one operation in a roadmap update batch (the daemon HTTP
// contract; agents use flat single-op MCP tools that each wrap one of these).
// Type discriminates which fields apply; unused fields are ignored per op.
// ParentID is a tri-state where the op has a prior location: for move, nil
// keeps the current parent and "" moves to top level; for restore, nil uses
// the original parent and "" restores to top level; for create, nil and ""
// both mean top level. A batch may not mix archive and restore ops (their
// crash-safety rename orders conflict).
type RoadmapOp struct {
	Type            RoadmapOpType      `json:"type"`
	ID              string             `json:"id,omitempty"`
	Kind            RoadmapItemKind    `json:"kind,omitempty"`
	Title           *string            `json:"title,omitempty"`
	Status          *RoadmapItemStatus `json:"status,omitempty"`
	Outcome         *string            `json:"outcome,omitempty"`
	SuccessCriteria *[]string          `json:"success_criteria,omitempty"`
	DependsOn       *[]string          `json:"depends_on,omitempty"`
	ParentID        *string            `json:"parent_id,omitempty"`
	Order           *int               `json:"order,omitempty"`
	TicketID        string             `json:"ticket_id,omitempty"`
	Summary         string             `json:"summary,omitempty"`
}

// UpdateRoadmapParams is the wire body for a roadmap update: a version-guarded
// atomic batch. Title, when present, replaces the roadmap-level title.
type UpdateRoadmapParams struct {
	Version string      `json:"version"`
	Title   *string     `json:"title,omitempty"`
	Ops     []RoadmapOp `json:"ops"`
}

// RoadmapView is the read result for both views. Items carries current-view
// items; ArchiveEntries carries compact summaries for archive listings;
// ArchiveEntry carries one complete stored subtree for focused archive reads.
type RoadmapView struct {
	View           string                       `json:"view"`
	Version        string                       `json:"version"`
	Title          string                       `json:"title,omitempty"`
	Items          []RoadmapItem                `json:"items"`
	ArchiveEntries []RoadmapArchiveEntrySummary `json:"archive_entries,omitempty"`
	ArchiveEntry   *RoadmapArchiveEntry         `json:"archive_entry,omitempty"`
	Tickets        []RoadmapTicketInfo          `json:"tickets"`
	Warnings       []string                     `json:"warnings"`
}

// RoadmapOpResult summarizes one applied operation.
type RoadmapOpResult struct {
	Index  int           `json:"index"`
	Type   RoadmapOpType `json:"type"`
	ItemID string        `json:"item_id"`
	Detail string        `json:"detail"`
}

// RoadmapUpdateResult is returned after a successful atomic batch.
type RoadmapUpdateResult struct {
	Version  string              `json:"version"`
	Applied  []RoadmapOpResult   `json:"applied"`
	Tickets  []RoadmapTicketInfo `json:"tickets"`
	Warnings []string            `json:"warnings"`
}
