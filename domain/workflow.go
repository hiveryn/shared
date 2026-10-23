package domain

import "time"

// Workflow is a plain Markdown file selected explicitly for a worker run.
type Workflow struct {
	Name        string                `json:"name"`
	Path        string                `json:"path"`
	RelPath     string                `json:"rel_path"`
	Valid       bool                  `json:"valid"`
	ModifiedAt  *time.Time            `json:"modified_at"`
	Diagnostics []WorkspaceDiagnostic `json:"diagnostics"`
}

type WorkflowList struct {
	ArchitectKey string                `json:"architect_key"`
	Workflows    []Workflow            `json:"workflows"`
	Diagnostics  []WorkspaceDiagnostic `json:"diagnostics"`
}
