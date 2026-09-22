package domain

import "time"

// WorkflowAttach is how a workflow reaches a worker session. There are exactly
// two modes and no third "always attached" category: manual workflows are only
// ever added by hand, suggested workflows are preselected when their repos
// overlap the session's writable scope. A repo match is a suggestion, never an
// obligation.
type WorkflowAttach string

const (
	WorkflowAttachManual    WorkflowAttach = "manual"
	WorkflowAttachSuggested WorkflowAttach = "suggested"
)

func (a WorkflowAttach) Valid() bool {
	switch a {
	case WorkflowAttachManual, WorkflowAttachSuggested:
		return true
	default:
		return false
	}
}

// Workflow is one discovered `workflows/*.md` file. Workflows are flat —
// immediate children of workflows/ only — have no dependencies on each other,
// and carry no precedence: a semantic conflict between two workflow bodies is
// an agent/user decision, not something this contract resolves.
//
// Path is the canonical absolute path inside the architect workspace, which is
// what a session records. Workflows are never copied or snapshotted.
//
// An invalid workflow is reported with its diagnostics rather than dropped, so
// the problem is repairable; it is excluded from Suggested regardless of repo
// overlap.
type Workflow struct {
	Name        string                `json:"name"`
	Path        string                `json:"path"`
	RelPath     string                `json:"rel_path"`
	Attach      WorkflowAttach        `json:"attach"`
	Repos       []string              `json:"repos"`
	Valid       bool                  `json:"valid"`
	Suggested   bool                  `json:"suggested"`
	ModifiedAt  *time.Time            `json:"modified_at"`
	Diagnostics []WorkspaceDiagnostic `json:"diagnostics"`
}

// WorkflowList is the workflow read contract. ScopeRepos echoes the writable
// repo scope the Suggested flags were computed against; an empty scope yields
// no suggestions at all rather than suggesting everything.
type WorkflowList struct {
	ArchitectKey string                `json:"architect_key"`
	ScopeRepos   []string              `json:"scope_repos"`
	Workflows    []Workflow            `json:"workflows"`
	Diagnostics  []WorkspaceDiagnostic `json:"diagnostics"`
}
