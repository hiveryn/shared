package domain

import (
	"regexp"
	"time"
)

// Actions are global, agent-operated procedures, independent of architects.
// Each lives in its own Git repository at HIVERYN_HOME/actions/<name>/ and is
// defined by two files the daemon reads live:
//
//   - action.yaml: exactly the keys name (equal to the directory name),
//     description (what the action does and what the caller's prompt must
//     contain) and artifacts (the textual contract of the delivered package,
//     including a schema summary where useful).
//   - KICKOFF.md: the launch instructions, containing the placeholders
//     {{prompt}} and {{output_dir}} exactly once or more each; no other
//     {{...}} placeholder is allowed.
//
// A launch starts one agent session in the action's repository with a fresh,
// daemon-created, empty output directory outside it. At most one execution of
// an action runs at a time, globally.

// actionNamePattern keeps an action name usable as a directory name, a path
// segment in API routes and an output-folder segment.
var actionNamePattern = regexp.MustCompile(`^[a-z0-9][a-z0-9._-]{0,63}$`)

// ValidActionName reports whether name is a well-formed action name: lowercase
// letters, digits, '.', '_' or '-', starting with a letter or digit, at most 64
// characters. The actions library and hiveryn.yaml's availableActions share it.
func ValidActionName(name string) bool {
	return actionNamePattern.MatchString(name)
}

// ActionProblem is one reason an action definition cannot be launched. Path is
// the offending file (or the action directory itself).
type ActionProblem struct {
	Path    string `json:"path"`
	Message string `json:"message"`
}

// ActionDefinition is the inspected form of one action directory. An invalid
// definition is still listed, with its problems, so it can be repaired; it
// cannot be launched. RunningExecutionID names the execution currently
// occupying the action, if any.
type ActionDefinition struct {
	Name               string          `json:"name"`
	Path               string          `json:"path"`
	Description        string          `json:"description"`
	Artifacts          string          `json:"artifacts"`
	Valid              bool            `json:"valid"`
	Problems           []ActionProblem `json:"problems"`
	RunningExecutionID string          `json:"running_execution_id,omitempty"`
}

// ActionList is every directory found under the actions root.
type ActionList struct {
	Root    string             `json:"root"`
	Actions []ActionDefinition `json:"actions"`
}

// ActionRunStatus is the lifecycle of one execution. Manual launches start in
// running (the user's launch is the approval) and end in completed or failed.
// An architect's request (executeAction) starts in pending_approval and moves
// to denied, to failed (it could not start, or its request was abandoned), or
// to running once the user approves and it launches — keeping the same
// execution id from request to result.
type ActionRunStatus string

const (
	ActionRunPendingApproval ActionRunStatus = "pending_approval"
	ActionRunDenied          ActionRunStatus = "denied"
	ActionRunRunning         ActionRunStatus = "running"
	ActionRunCompleted       ActionRunStatus = "completed"
	ActionRunFailed          ActionRunStatus = "failed"
)

func (s ActionRunStatus) Terminal() bool {
	return s == ActionRunDenied || s == ActionRunCompleted || s == ActionRunFailed
}

// ActionRunTrigger records who started an execution.
type ActionRunTrigger string

const (
	// ActionRunTriggerManual: the user launched it from the Actions window.
	ActionRunTriggerManual ActionRunTrigger = "manual"
	// ActionRunTriggerArchitect: an architect requested it and the user approved.
	ActionRunTriggerArchitect ActionRunTrigger = "architect"
)

// ActionRun is the durable record of one execution, addressed by its stable
// id. It outlives the agent session: SessionID names the session while it
// exists, OutputDir the delivered artifact folder. Summary is the agent's
// concluding summary; Error is the daemon's reason for an execution that
// failed without (or before) an agent conclusion. For an architect request,
// ArchitectKey scopes who may read the result, RequesterSessionID names the
// architect session that asked, Reason keeps the user's denial reason, and
// ProfileName and StartedAt stay empty until the user approves and it starts.
//
// A completed run means the requested artifact package was delivered — the
// artifacts may still report findings such as failed checks. failed means the
// action could not be executed or delivered.
type ActionRun struct {
	ID          string           `json:"id"`
	Action      string           `json:"action"`
	Trigger     ActionRunTrigger `json:"trigger"`
	Status      ActionRunStatus  `json:"status"`
	Prompt      string           `json:"prompt"`
	ProfileName string           `json:"profile_name"`
	RepoPath    string           `json:"repo_path"`
	OutputDir   string           `json:"output_dir"`
	SessionID   string           `json:"session_id,omitempty"`
	Summary     string           `json:"summary,omitempty"`
	Error       string           `json:"error,omitempty"`
	CreatedAt   time.Time        `json:"created_at"`
	StartedAt   *time.Time       `json:"started_at,omitempty"`
	EndedAt     *time.Time       `json:"ended_at,omitempty"`

	ArchitectKey       string `json:"architect_key,omitempty"`
	RequesterSessionID string `json:"requester_session_id,omitempty"`
	Reason             string `json:"reason,omitempty"`
}

// LaunchActionRequest is the desktop's manual launch: the caller's prompt and
// the agent variant. Cols/Rows size the agent terminal.
type LaunchActionRequest struct {
	Prompt      string `json:"prompt"`
	ProfileName string `json:"profile_name"`
	Cols        uint16 `json:"cols,omitempty"`
	Rows        uint16 `json:"rows,omitempty"`
}

// LaunchActionResult is the running execution and its live agent session.
type LaunchActionResult struct {
	Run            ActionRun `json:"run"`
	Session        Session   `json:"session"`
	MainTerminalID string    `json:"main_terminal_id,omitempty"`
}

// ActionConclusionOutcome is the action agent's verdict on its own execution.
type ActionConclusionOutcome string

const (
	// ActionConclusionCompleted: the requested artifact package was delivered.
	ActionConclusionCompleted ActionConclusionOutcome = "completed"
	// ActionConclusionFailed: the action could not be executed or delivered.
	ActionConclusionFailed ActionConclusionOutcome = "failed"
)

// MaxActionSummaryLength bounds a conclusion summary; it is meant to be read
// at a glance and fed back to later executions.
const MaxActionSummaryLength = 2000

// ConcludeActionRequest is the action agent's conclusion.
type ConcludeActionRequest struct {
	Outcome ActionConclusionOutcome `json:"outcome"`
	Summary string                  `json:"summary"`
}

// ActionConclusion is one earlier execution's concluding summary, as returned
// to an action agent reading its action's history.
type ActionConclusion struct {
	ExecutionID string          `json:"execution_id"`
	Status      ActionRunStatus `json:"status"`
	Summary     string          `json:"summary"`
	EndedAt     *time.Time      `json:"ended_at,omitempty"`
}

// ExecuteActionRequest is an architect's request to run one of its available
// Actions. The prompt must contain what the action's description asks for.
type ExecuteActionRequest struct {
	Name   string `json:"name"`
	Prompt string `json:"prompt"`
}

// AvailableActionList is the Actions an architect may request: exactly the
// names in its hiveryn.yaml availableActions, in that order, each enriched
// from the library. A configured name with no directory in the library is
// listed invalid with a problem saying so, never dropped.
type AvailableActionList struct {
	Actions []ActionDefinition `json:"actions"`
}

// ActionAgentActivity is what is known about the Action agent's activity.
// Available is false when nothing is known — before the execution starts,
// after it ends, or while the agent has reported no status yet — and Status
// is then empty rather than guessed. Status is the agent's last reported
// state: active, idle, waiting or stopped.
type ActionAgentActivity struct {
	Available bool   `json:"available"`
	Status    string `json:"status,omitempty"`
}

// ActionResult is an architect's view of one requested execution, addressed by
// the execution id executeAction returned. The execution record is
// authoritative: pending_approval means nothing has started (StartedAt unset),
// running means the Action agent is working, and completed/failed/denied are
// final. ElapsedSeconds counts from StartedAt to EndedAt (or now) and is unset
// until it starts. Reason is the user's denial reason, Error the reason it
// failed without or before an agent conclusion, Summary the agent's
// conclusion; OutputDir is set once the execution started.
type ActionResult struct {
	ExecutionID    string              `json:"execution_id"`
	Action         string              `json:"action"`
	Status         ActionRunStatus     `json:"status"`
	Prompt         string              `json:"prompt"`
	ProfileName    string              `json:"profile_name,omitempty"`
	RequestedAt    time.Time           `json:"requested_at"`
	StartedAt      *time.Time          `json:"started_at,omitempty"`
	EndedAt        *time.Time          `json:"ended_at,omitempty"`
	ElapsedSeconds *int64              `json:"elapsed_seconds,omitempty"`
	Reason         string              `json:"reason,omitempty"`
	Error          string              `json:"error,omitempty"`
	Summary        string              `json:"summary,omitempty"`
	OutputDir      string              `json:"output_dir,omitempty"`
	Activity       ActionAgentActivity `json:"activity"`
}

// MaxActionWaitSeconds caps one waitForActionResult call; a caller that needs
// longer waits again.
const MaxActionWaitSeconds = 30

// ActionWaitResult is one bounded wait: the result at return, and whether the
// wait ended on a status change (Changed) or at its timeout (TimedOut). An
// already-final execution returns at once with neither set.
type ActionWaitResult struct {
	Result   ActionResult `json:"result"`
	Changed  bool         `json:"changed"`
	TimedOut bool         `json:"timed_out"`
}

// ActionEventType discriminates events on the actions stream
// (GET /api/actions/events).
const ActionEventType = "action_changed"

// ActionEvent announces that an execution changed status: it was requested
// (pending_approval), started (running) or ended (denied, completed, failed). The stream has no backlog, so a
// client reconciles by refetching executions on (re)connect.
type ActionEvent struct {
	Type        string          `json:"type"`
	Action      string          `json:"action"`
	ExecutionID string          `json:"execution_id"`
	Status      ActionRunStatus `json:"status"`
	SessionID   string          `json:"session_id,omitempty"`
	At          time.Time       `json:"at"`
}
