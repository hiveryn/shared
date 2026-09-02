package domain

import "time"

type SessionType string

type SessionCreatedBy string

type SessionRunStatus string

type SessionRunFailureReason string

const (
	SessionTypeArchitect SessionType = "architect"
	SessionTypeTicket    SessionType = "ticket"
	SessionTypeFreeform  SessionType = "freeform"
)

const (
	SessionCreatedByDesktop      SessionCreatedBy = "desktop"
	SessionCreatedByArchitectMCP SessionCreatedBy = "architect_mcp"
)

const (
	SessionRunStatusRunning   SessionRunStatus = "running"
	SessionRunStatusCompleted SessionRunStatus = "completed"
	SessionRunStatusFailed    SessionRunStatus = "failed"
)

const (
	SessionRunFailureLaunchFailed  SessionRunFailureReason = "launch_failed"
	SessionRunFailureProcessExited SessionRunFailureReason = "process_exited"
	SessionRunFailureRestoreFailed SessionRunFailureReason = "restore_failed"
	SessionRunFailureUserCancelled SessionRunFailureReason = "user_cancelled"
)

type MCPServerSnapshot struct {
	Command           string            `json:"command,omitempty"`
	Args              []string          `json:"args,omitempty"`
	Env               map[string]string `json:"env,omitempty"`
	CWD               string            `json:"cwd,omitempty"`
	URL               string            `json:"url,omitempty"`
	BearerTokenEnvVar string            `json:"bearer_token_env_var,omitempty"`
}

type AgentProfileSnapshot struct {
	Agent string                       `json:"agent"`
	Model string                       `json:"model,omitempty"`
	Yolo  bool                         `json:"yolo,omitempty"`
	Mode  string                       `json:"mode,omitempty"`
	Args  []string                     `json:"args"`
	Env   map[string]string            `json:"env"`
	MCP   map[string]MCPServerSnapshot `json:"mcp,omitempty"`
}

type Session struct {
	ID                 string           `json:"id"`
	ArchitectKey       string           `json:"architect_key"`
	SessionType        SessionType      `json:"session_type"`
	ContextID          string           `json:"context_id"`
	Prompt             string           `json:"prompt"`
	Workdir            string           `json:"workdir"`
	AdditionalRepos    []string         `json:"additional_repos"`
	AdditionalWorkdirs []string         `json:"additional_workdirs"`
	Instructions       string           `json:"instructions,omitempty"`
	CreatedBy          SessionCreatedBy `json:"created_by,omitempty"`
	CreatedAt          time.Time        `json:"created_at"`
	UpdatedAt          time.Time        `json:"updated_at"`
	CurrentRun         *SessionRun      `json:"current_run,omitempty"`
}

type SessionRun struct {
	ID                 string                  `json:"id"`
	SessionID          string                  `json:"session_id"`
	Status             SessionRunStatus        `json:"status"`
	AgentStatus        string                  `json:"agent_status,omitempty"`
	ProfileName        string                  `json:"profile_name"`
	ProfileSnapshot    *AgentProfileSnapshot   `json:"profile_snapshot,omitempty"`
	Workdir            string                  `json:"workdir"`
	AdditionalRepos    []string                `json:"additional_repos"`
	AdditionalWorkdirs []string                `json:"additional_workdirs"`
	NativeID           string                  `json:"native_id,omitempty"`
	FailureReason      SessionRunFailureReason `json:"failure_reason,omitempty"`
	MainTerminalID     string                  `json:"main_terminal_id,omitempty"`
	StartedAt          *time.Time              `json:"started_at,omitempty"`
	EndedAt            *time.Time              `json:"ended_at,omitempty"`
	CreatedAt          time.Time               `json:"created_at"`
	UpdatedAt          time.Time               `json:"updated_at"`
}

type SessionEvent struct {
	ID                string            `json:"id"`
	SessionID         string            `json:"session_id"`
	RunID             string            `json:"run_id,omitempty"`
	Seq               int64             `json:"seq"`
	Type              string            `json:"type"`
	Status            string            `json:"status,omitempty"`
	Tool              string            `json:"tool,omitempty"`
	Message           string            `json:"message,omitempty"`
	NativeID          string            `json:"native_id,omitempty"`
	PrimaryNativeID   string            `json:"primary_native_id,omitempty"`
	NativeSessionRole string            `json:"native_session_role,omitempty"`
	Metadata          map[string]string `json:"metadata,omitempty"`
	Raw               map[string]any    `json:"raw,omitempty"`
	At                time.Time         `json:"at"`
}

type CreateSessionRequest struct {
	SessionType  SessionType `json:"session_type"`
	ArchitectKey string      `json:"architect_key"`
	TicketID     string      `json:"ticket_id,omitempty"`
	Prompt       string      `json:"prompt,omitempty"`
	Workdir      string      `json:"workdir,omitempty"`
	Slug         string      `json:"slug,omitempty"`
}

type CreateSessionParams struct {
	ID                 string
	ArchitectKey       string
	SessionType        SessionType
	ContextID          string
	Prompt             string
	Workdir            string
	AdditionalRepos    []string
	AdditionalWorkdirs []string
	Instructions       string
	CreatedBy          SessionCreatedBy
}

type CreateSessionRunRequest struct {
	ProfileName string `json:"profile_name"`
	Cols        uint16 `json:"cols,omitempty"`
	Rows        uint16 `json:"rows,omitempty"`
}

type CreateSessionRunParams struct {
	ID                 string
	SessionID          string
	ProfileName        string
	ProfileSnapshot    AgentProfileSnapshot
	Workdir            string
	AdditionalRepos    []string
	AdditionalWorkdirs []string
	NativeID           string
	StartedAt          time.Time
}

type CreateSessionRunResult struct {
	Run            SessionRun `json:"run"`
	MainTerminalID string     `json:"main_terminal_id,omitempty"`
}

type AppendSessionEventParams struct {
	SessionID         string
	RunID             string
	Type              string
	Status            string
	Tool              string
	Message           string
	NativeID          string
	PrimaryNativeID   string
	NativeSessionRole string
	Metadata          map[string]string
	Raw               map[string]any
	At                time.Time
}

// ConcludeSessionParams carries the structured conclusion input from the MCP
// conclude tools through to the daemon write path. The daemon renders the
// structured narrative/section fields into the canonical conclusion.md markdown
// body (see sessionruntime.render*ConclusionBody); Body is the rendered output
// carrier set by the daemon, not a wire input. The [fm] metadata fields
// (Commits/Outcome/RejectionReason, plus server-computed timestamps/agent)
// still flow into frontmatter exactly as before.
//
// Every presentational section is a Markdown string — the agent authors its own
// bullets/prose as text. Only Commits stays structured, because it is the one
// field persisted as data (YAML frontmatter) and read back per-commit, not just
// rendered. Markdown strings can never trip the MCP client's required-array drop
// bug, so no conclude section field is a required array.
type ConcludeSessionParams struct {
	Body            string
	Commits         []CommitRef
	Outcome         TicketOutcome
	RejectionReason string

	// Structured body fields (rendered into markdown sections). Which fields
	// are populated/required depends on the session type; see the per-type
	// render functions in the daemon.
	Summary         string // all types (required)
	Narrative       string // architect (required)
	Implementation  string // ticket (required unless rejected)
	Findings        string // freeform (required)
	Verification    string // ticket (optional)
	TicketsTouched  string // architect (optional, Markdown)
	Decisions       string // architect (optional, Markdown)
	ConfigChanges   string // architect (optional, Markdown)
	UserPriorities  string // architect (optional, Markdown)
	Deviations      string // ticket (optional, Markdown)
	FollowUps       string // ticket (optional, Markdown)
	Recommendations string // freeform (required, Markdown)
	OpenQuestions   string // all types; required for freeform, optional otherwise (Markdown)
	NextSteps       string // architect (required, Markdown)
}

type ConcludeSessionResult struct {
	SessionID    string `json:"session_id"`
	ArchitectKey string `json:"architect_key"`
	TicketID     string `json:"ticket_id,omitempty"`
}

type MoveTicketToDoneParams struct {
	Body            string
	Commits         []CommitRef
	Outcome         TicketOutcome
	RejectionReason string
}

type MoveTicketToDoneResult struct {
	TicketID     string `json:"ticket_id"`
	ArchitectKey string `json:"architect_key"`
}

type TerminalInfo struct {
	TerminalID         string `json:"terminal_id"`
	SessionID          string `json:"session_id"`
	Command            string `json:"command"`
	Status             string `json:"status"`
	WorkdirID          string `json:"workdir_id,omitempty"`
	WorkdirTitle       string `json:"workdir_title,omitempty"`
	WorkdirPath        string `json:"workdir_path,omitempty"`
	WorkdirDisplayPath string `json:"workdir_display_path,omitempty"`
}

type TerminalWorkdir struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Path        string `json:"path"`
	DisplayPath string `json:"display_path"`
	Default     bool   `json:"default"`
}

type TerminalPlacement string

const (
	TerminalPlacementTab   TerminalPlacement = "tab"
	TerminalPlacementSplit TerminalPlacement = "split"
)

type CreateTerminalParams struct {
	Placement TerminalPlacement `json:"placement"`
	BaseTabID string            `json:"base_tab_id,omitempty"`
	WorkdirID string            `json:"workdir_id"`
}

type SessionTab struct {
	Type               string            `json:"type"`
	ID                 string            `json:"id,omitempty"`
	Command            string            `json:"command,omitempty"`
	Status             string            `json:"status,omitempty"`
	Placement          TerminalPlacement `json:"placement,omitempty"`
	BaseTabID          string            `json:"base_tab_id,omitempty"`
	Target             string            `json:"target,omitempty"`
	WorkdirID          string            `json:"workdir_id,omitempty"`
	WorkdirTitle       string            `json:"workdir_title,omitempty"`
	WorkdirPath        string            `json:"workdir_path,omitempty"`
	WorkdirDisplayPath string            `json:"workdir_display_path,omitempty"`
}

type PreviewBrowserTabParams struct {
	Target string `json:"target"`
	TabID  string `json:"tab_id,omitempty"`
}

type BrowserTabInfo struct {
	TabID     string `json:"tab_id"`
	SessionID string `json:"session_id"`
	Target    string `json:"target"`
}

type ArchitectConclusion struct {
	StartedAt   time.Time `json:"started_at"`
	ConcludedAt time.Time `json:"concluded_at"`
	Agent       string    `json:"agent,omitempty"`
	Body        string    `json:"body"`
}

type ConclusionSummary struct {
	ID          string    `json:"id"`
	ConcludedAt time.Time `json:"concluded_at"`
}
