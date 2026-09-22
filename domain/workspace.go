package domain

import "time"

// ArtifactKind names one file-based architect workspace artifact. Tickets and
// conclusions are deliberately absent: they keep their own schemas and MCP
// lifecycle and are never described or validated by the workspace surface.
type ArtifactKind string

const (
	ArtifactProjectOverview ArtifactKind = "PROJECT_OVERVIEW"
	ArtifactProjectState    ArtifactKind = "PROJECT_STATE"
	ArtifactRoadmapCurrent  ArtifactKind = "ROADMAP_CURRENT"
	ArtifactRoadmapArchive  ArtifactKind = "ROADMAP_ARCHIVE"
	ArtifactArchitectSystem ArtifactKind = "ARCHITECT_SYSTEM"
	ArtifactWorkflow        ArtifactKind = "WORKFLOW"
	ArtifactHiverynYAML     ArtifactKind = "HIVERYN_YAML"
)

func (k ArtifactKind) Valid() bool {
	switch k {
	case ArtifactProjectOverview, ArtifactProjectState, ArtifactRoadmapCurrent,
		ArtifactRoadmapArchive, ArtifactArchitectSystem, ArtifactWorkflow, ArtifactHiverynYAML:
		return true
	default:
		return false
	}
}

// ArtifactKinds is the complete, ordered catalog. describeArtifact accepts
// exactly these values.
func ArtifactKinds() []ArtifactKind {
	return []ArtifactKind{
		ArtifactHiverynYAML,
		ArtifactProjectOverview,
		ArtifactProjectState,
		ArtifactRoadmapCurrent,
		ArtifactRoadmapArchive,
		ArtifactArchitectSystem,
		ArtifactWorkflow,
	}
}

// DiagnosticSeverity separates blocking structural problems from advisory
// ones. Only Error contributes to WorkspaceReport.Valid.
type DiagnosticSeverity string

const (
	DiagnosticError   DiagnosticSeverity = "error"
	DiagnosticWarning DiagnosticSeverity = "warning"
)

// Diagnostic codes. Structural only — nothing here encodes a judgment about
// whether a document's contents are factually current.
const (
	DiagMissingRequiredFile      = "MISSING_REQUIRED_FILE"
	DiagMissingRequiredDirectory = "MISSING_REQUIRED_DIRECTORY"
	DiagNotAFile                 = "NOT_A_FILE"
	DiagNotADirectory            = "NOT_A_DIRECTORY"
	DiagUnreadable               = "UNREADABLE"
	DiagInvalidUTF8              = "INVALID_UTF8"
	DiagIncompleteRead           = "INCOMPLETE_READ"
	DiagOutsideWorkspace         = "OUTSIDE_WORKSPACE"
	DiagEmptyBody                = "EMPTY_BODY"
	DiagMissingFrontmatter       = "MISSING_FRONTMATTER"
	DiagInvalidFrontmatter       = "INVALID_FRONTMATTER"
	DiagMissingField             = "MISSING_FIELD"
	DiagInvalidTimestamp         = "INVALID_TIMESTAMP"
	DiagUnexpectedField          = "UNEXPECTED_FIELD"
	DiagArchiveNameInvalid       = "ARCHIVE_NAME_INVALID"
	DiagArchiveDateMismatch      = "ARCHIVE_DATE_MISMATCH"
	DiagWorkflowSubdirectory     = "WORKFLOW_SUBDIRECTORY"
	DiagWorkflowInvalidAttach    = "WORKFLOW_INVALID_ATTACH"
	DiagWorkflowMissingRepos     = "WORKFLOW_MISSING_REPOS"
	DiagWorkflowUnknownRepo      = "WORKFLOW_UNKNOWN_REPO"
	DiagWorkflowDuplicateRepo    = "WORKFLOW_DUPLICATE_REPO"
	DiagWorkflowReposNotAllowed  = "WORKFLOW_REPOS_NOT_ALLOWED"
	DiagConfigInvalid            = "CONFIG_INVALID"
	DiagRepoPathMissing          = "REPO_PATH_MISSING"
	DiagRepoPathNotDirectory     = "REPO_PATH_NOT_DIRECTORY"
	DiagTicketScanFailed         = "TICKET_SCAN_FAILED"
)

// WorkspaceDiagnostic is one actionable finding. Path is workspace-relative
// (empty for workspace-wide findings). Line is the 1-indexed source line when
// the finding maps to one, and 0 otherwise — an absent field has no line.
type WorkspaceDiagnostic struct {
	Code     string             `json:"code"`
	Severity DiagnosticSeverity `json:"severity"`
	Path     string             `json:"path"`
	Line     int                `json:"line"`
	Message  string             `json:"message"`
}

// WorkspaceNodeType distinguishes the two node shapes in the report tree.
type WorkspaceNodeType string

const (
	WorkspaceNodeFile      WorkspaceNodeType = "file"
	WorkspaceNodeDirectory WorkspaceNodeType = "directory"
)

// WorkspaceNode is one entry of the workspace's expected shape: a root document,
// the config, or one of the two directories. Expected nodes are always present
// in the report even when the file is missing, so a broken workspace stays
// inspectable and repairable.
//
// DocumentUpdatedAt is the timestamp the document itself declares
// (lastUpdatedAt / archivedAt); ModifiedAt is the filesystem mtime. They are
// reported separately and neither certifies factual freshness.
type WorkspaceNode struct {
	Kind              ArtifactKind          `json:"kind"`
	Type              WorkspaceNodeType     `json:"type"`
	Path              string                `json:"path"`
	Required          bool                  `json:"required"`
	Exists            bool                  `json:"exists"`
	Valid             bool                  `json:"valid"`
	DocumentUpdatedAt *time.Time            `json:"document_updated_at"`
	ModifiedAt        *time.Time            `json:"modified_at"`
	Children          []WorkspaceEntry      `json:"children"`
	Diagnostics       []WorkspaceDiagnostic `json:"diagnostics"`
}

// WorkspaceEntry is one file discovered inside a directory node — a workflow or
// an archived roadmap.
//
// It has no children of its own, and that is a property of the workspace rather
// than a simplification: both directories are flat by rule, so a subdirectory
// inside them is reported as an error on the directory and never descended
// into. Keeping the distinction in the types also keeps the report
// non-recursive, which the MCP tool schema requires.
//
// It carries no Exists or Required: a discovered file exists by definition, and
// no individual workflow or archive is required.
type WorkspaceEntry struct {
	Kind              ArtifactKind          `json:"kind"`
	Path              string                `json:"path"`
	Valid             bool                  `json:"valid"`
	DocumentUpdatedAt *time.Time            `json:"document_updated_at"`
	ModifiedAt        *time.Time            `json:"modified_at"`
	Diagnostics       []WorkspaceDiagnostic `json:"diagnostics"`
}

// WorkspaceTicketTotals counts tickets per status from the current ticket
// system. The workspace check reports them; it does not validate ticket files.
type WorkspaceTicketTotals struct {
	Backlog  int `json:"backlog"`
	Progress int `json:"progress"`
	Done     int `json:"done"`
}

// WorkspaceTicketWarning is one warning raised by the ticket system for one
// backlog ticket. Advisory only — ticket warnings never make a workspace
// invalid.
type WorkspaceTicketWarning struct {
	TicketID string `json:"ticket_id"`
	Code     string `json:"code"`
	Message  string `json:"message"`
}

type WorkspaceTickets struct {
	Totals   WorkspaceTicketTotals    `json:"totals"`
	Warnings []WorkspaceTicketWarning `json:"warnings"`
}

// WorkspaceReport is the full result of a parameterless workspace check. It is
// read-only, deterministic and mutates nothing: node and diagnostic order is
// stable across runs for the same on-disk state.
//
// Valid is true only when no error-severity diagnostic appears anywhere in the
// report. It is a structural verdict, not an authorization and not a claim
// about the truth of any document's contents.
type WorkspaceReport struct {
	ArchitectKey  string                `json:"architect_key"`
	WorkspacePath string                `json:"workspace_path"`
	CheckedAt     time.Time             `json:"checked_at"`
	Valid         bool                  `json:"valid"`
	Nodes         []WorkspaceNode       `json:"nodes"`
	Tickets       WorkspaceTickets      `json:"tickets"`
	Diagnostics   []WorkspaceDiagnostic `json:"diagnostics"`
}

// ArtifactFieldType names the value shape a frontmatter/config field takes. It
// is what the validator actually enforces, rendered for describeArtifact.
type ArtifactFieldType string

const (
	ArtifactFieldString       ArtifactFieldType = "string"
	ArtifactFieldStringList   ArtifactFieldType = "string-list"
	ArtifactFieldStringMap    ArtifactFieldType = "string-map"
	ArtifactFieldEnum         ArtifactFieldType = "enum"
	ArtifactFieldTimestampUTC ArtifactFieldType = "rfc3339-utc-datetime"
	ArtifactFieldMapping      ArtifactFieldType = "mapping"
	ArtifactFieldMappingList  ArtifactFieldType = "mapping-list"
)

// ArtifactField describes one frontmatter or config field.
type ArtifactField struct {
	Name        string            `json:"name"`
	Required    bool              `json:"required"`
	Type        ArtifactFieldType `json:"type"`
	Enum        []string          `json:"enum,omitempty"`
	Description string            `json:"description"`
}

// ArtifactSchema is describeArtifact's result. It is rendered from the same
// definition the validators execute, so a rule described here is a rule
// enforced by the workspace check.
type ArtifactSchema struct {
	Kind          ArtifactKind    `json:"kind"`
	SchemaVersion int             `json:"schema_version"`
	Title         string          `json:"title"`
	Required      bool            `json:"required"`
	Format        string          `json:"format"`
	Location      string          `json:"location"`
	Naming        string          `json:"naming"`
	Fields        []ArtifactField `json:"fields"`
	Rules         []string        `json:"rules"`
	Example       string          `json:"example"`
}

// WorkerPreflight answers one question before a ticket session is created: is
// the architect workspace's required project context ready for a worker right
// now? It runs the project-context half of the worker launch validation — the
// same code the launch itself runs — and reports what would block it.
//
// It is deliberately not the workspace check's aggregate verdict: a worker is
// blocked only by hiveryn.yaml, PROJECT_OVERVIEW.md, PROJECT_STATE.md and an
// invalid-when-present ROADMAP_CURRENT.md. Architect-only artifacts, archived
// roadmaps and unselected invalid workflows never block one, so a workspace
// that the check calls invalid can still be launchable.
//
// Selected-workflow validity is not part of this answer: each workflow reports
// its own validity and diagnostics in the workflow listing, and the launch
// revalidates the whole selection against the live workspace regardless.
//
// Launchable is false exactly when Problems is non-empty; each problem is one
// actionable finding, in the same wording the launch error would carry.
type WorkerPreflight struct {
	ArchitectKey string    `json:"architect_key"`
	CheckedAt    time.Time `json:"checked_at"`
	Launchable   bool      `json:"launchable"`
	Problems     []string  `json:"problems"`
}
