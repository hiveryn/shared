package domain

import (
	"regexp"
	"time"
)

// Actions are global, agent-operated procedures, independent of architects.
// Each is registered in HIVERYN_HOME/actions.yaml, which maps its name to an
// existing Git checkout and an optional machine (machines.yaml key; omitted
// means local). Registration neither clones nor synchronizes the checkout.
// The checkout is defined by two files the daemon reads live, on the machine
// where it lives:
//
//   - action.yaml: the keys name (equal to the directory name), description
//     (what the action does and what the caller's prompt must contain) and
//     artifacts (the textual contract of the delivered package, including a
//     schema summary where useful), plus the optional suggestions: a list of
//     ready-made prompts offered in the manual launch form. No other key is
//     allowed.
//   - KICKOFF.md: the launch instructions, containing the placeholders
//     {{prompt}} and {{output_dir}} exactly once or more each; no other
//     {{...}} placeholder is allowed.
//
// A launch starts one agent session in the action's checkout, on its machine,
// with a fresh, daemon-created, empty working output directory outside it. At
// most one execution of an action runs at a time, globally. A completed
// execution names the folder it delivered; the daemon freezes that folder into
// a managed package on the execution machine and copies the package to each
// caller's machine on demand.

// actionNamePattern keeps an action name usable as a directory name, a path
// segment in API routes and an output-folder segment.
var actionNamePattern = regexp.MustCompile(`^[a-z0-9][a-z0-9._-]{0,63}$`)

// Bounds on action.yaml suggestions: at most MaxActionSuggestions entries,
// each nonblank, unique and at most MaxActionSuggestionLength characters.
const (
	MaxActionSuggestions      = 10
	MaxActionSuggestionLength = 1000
)

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
//
// Suggestions are the action.yaml suggested prompts, in file order: manual
// launch conveniences that only prefill the prompt. They are part of the
// library views (ActionList, a single definition) and are left out of an
// architect's AvailableActionList.
//
// Machine is where the checkout lives and the Action runs (empty means
// local); Path is a path on that machine. Variants lists the configured agent
// variants eligible to run it there.
type ActionDefinition struct {
	Name               string          `json:"name"`
	Path               string          `json:"path"`
	Machine            string          `json:"machine,omitempty"`
	Variants           []string        `json:"variants"`
	Description        string          `json:"description"`
	Artifacts          string          `json:"artifacts"`
	Suggestions        []string        `json:"suggestions,omitempty"`
	Valid              bool            `json:"valid"`
	Problems           []ActionProblem `json:"problems"`
	RunningExecutionID string          `json:"running_execution_id,omitempty"`
}

// ActionList is every Action registered in Registry (actions.yaml).
type ActionList struct {
	Registry string             `json:"registry"`
	Actions  []ActionDefinition `json:"actions"`
}

// ActionRunStatus is the lifecycle of one execution. Manual launches start in
// running (the user's launch is the approval) and end in completed or failed.
// An architect's request (executeAction) starts in pending_approval and moves
// to denied, to failed (it could not start, or its request was abandoned), or
// to running once it is approved (by the user, or when the approval window
// expires) and it launches — keeping the same
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
	// ActionRunTriggerWorker: a ticket worker requested it and the user approved.
	ActionRunTriggerWorker ActionRunTrigger = "worker"
)

// ActionRun is the durable record of one execution, addressed by its stable
// id. It outlives the agent session: SessionID names the session while it
// exists, OutputDir the delivered artifact folder. Summary is the agent's
// concluding summary; Error is the daemon's reason for an execution that
// failed without (or before) an agent conclusion. For an agent request
// (architect or worker), ArchitectKey is the project that may read the result,
// RequesterSessionID names the session that asked, RequesterTicketID the
// requesting worker's ticket, Reason keeps the user's denial reason, and
// ProfileName and StartedAt stay empty until the user approves and it starts.
//
// Machine is where the execution runs (empty means local); RepoPath and
// OutputDir — the working output directory created for the agent — are paths
// there. DeliveryDir is the folder the agent concluded with; Package is its
// frozen snapshot on the execution machine and Deliveries the verified copies
// made for callers on other machines.
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
	Machine     string           `json:"machine,omitempty"`
	DeliveryDir string           `json:"delivery_dir,omitempty"`
	SessionID   string           `json:"session_id,omitempty"`
	Summary     string           `json:"summary,omitempty"`
	Error       string           `json:"error,omitempty"`
	CreatedAt   time.Time        `json:"created_at"`
	StartedAt   *time.Time       `json:"started_at,omitempty"`
	EndedAt     *time.Time       `json:"ended_at,omitempty"`

	ArchitectKey       string `json:"architect_key,omitempty"`
	RequesterSessionID string `json:"requester_session_id,omitempty"`
	RequesterTicketID  string `json:"requester_ticket_id,omitempty"`
	Reason             string `json:"reason,omitempty"`

	Package    *ActionArtifactCopy  `json:"package,omitempty"`
	Deliveries []ActionArtifactCopy `json:"deliveries,omitempty"`

	// Attention is live, not part of the durable record: the agent's attention
	// while the execution runs, attached by the daemon when it serves a
	// running execution and absent otherwise.
	Attention *ActionAgentAttention `json:"attention,omitempty"`
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

// ConcludeActionRequest is the action agent's conclusion. DeliveryDir is the
// absolute path, on the execution machine, of the folder holding the
// delivered package: required for completed, refused for failed.
type ConcludeActionRequest struct {
	Outcome     ActionConclusionOutcome `json:"outcome"`
	Summary     string                  `json:"summary"`
	DeliveryDir string                  `json:"delivery_dir,omitempty"`
}

// ActionArtifactState is the state of a frozen package or of one caller copy.
type ActionArtifactState string

const (
	// ActionArtifactsNone: there is no package — the execution has not
	// completed, or it ended without one.
	ActionArtifactsNone ActionArtifactState = "none"
	// ActionArtifactsCopying: a copy is being made; it is not usable yet.
	ActionArtifactsCopying ActionArtifactState = "copying"
	// ActionArtifactsReady: the copy is complete and verified.
	ActionArtifactsReady ActionArtifactState = "ready"
	// ActionArtifactsFailed: the last copy attempt failed (Error says why);
	// reading the result again retries it.
	ActionArtifactsFailed ActionArtifactState = "failed"
	// ActionArtifactsExpired: the retention period ended; the files are gone
	// (or are being removed) while the execution record stays.
	ActionArtifactsExpired ActionArtifactState = "expired"
)

// ActionArtifactCopy is one Hiveryn-managed artifact directory: the frozen
// package on the execution machine, or a caller copy on another machine.
// Machine is where Path lives (empty means local). ExpiresAt counts from the
// execution's completion and is nil for a package that never expires (one
// delivered before frozen packages existed, kept at its original folder).
// CleanupError is the last failed attempt to remove it after expiry, retried
// at the next daemon start.
type ActionArtifactCopy struct {
	Machine      string              `json:"machine,omitempty"`
	Path         string              `json:"path"`
	State        ActionArtifactState `json:"state"`
	Error        string              `json:"error,omitempty"`
	CleanupError string              `json:"cleanup_error,omitempty"`
	Files        int                 `json:"files"`
	Bytes        int64               `json:"bytes"`
	ReadyAt      *time.Time          `json:"ready_at,omitempty"`
	ExpiresAt    *time.Time          `json:"expires_at,omitempty"`
}

// ActionArtifactLocation names a directory on a machine (empty means local).
type ActionArtifactLocation struct {
	Machine string `json:"machine,omitempty"`
	Path    string `json:"path"`
}

// ActionArtifacts is a caller's view of a completed execution's package, on
// the caller's own machine. State is ready once a usable copy exists there —
// the frozen package itself when the caller runs where the Action ran, a
// verified copy otherwise — and the result's OutputDir is then its path.
// copying means a copy is in progress (read the result again), failed that the
// last attempt failed (Error says why; reading the result again retries), and
// expired that the retention period ended. Delivery never changes the
// execution's status. Source is the frozen package on the execution machine.
type ActionArtifacts struct {
	State     ActionArtifactState     `json:"state"`
	Error     string                  `json:"error,omitempty"`
	Files     int                     `json:"files,omitempty"`
	Bytes     int64                   `json:"bytes,omitempty"`
	ExpiresAt *time.Time              `json:"expires_at,omitempty"`
	Source    *ActionArtifactLocation `json:"source,omitempty"`
}

// ActionConclusion is one earlier execution's concluding summary, as returned
// to an action agent reading its action's history.
type ActionConclusion struct {
	ExecutionID string          `json:"execution_id"`
	Status      ActionRunStatus `json:"status"`
	Summary     string          `json:"summary"`
	EndedAt     *time.Time      `json:"ended_at,omitempty"`
}

// ExecuteActionRequest is an architect's or worker's request to run one of its
// project's available Actions. The prompt must contain what the action's
// description asks for. Variant is the configured agent variant that runs it:
// required, never defaulted, and part of the request's identity, so asking for
// another variant is a different request.
type ExecuteActionRequest struct {
	Name    string `json:"name"`
	Prompt  string `json:"prompt"`
	Variant string `json:"variant"`
}

// ExecuteActionResponse is the resolved executeAction request. The call waits
// for the approval like createWorkTicket (wait-then-allow): Outcome is the
// intent verdict — approved or auto_approved (the Action was launched and now
// runs independently), denied_by_user (it never ran) or error (it could not
// start; Reason says why). Result is the execution record under the execution
// id, which is also the intent id; follow it with the result and wait tools.
// The call never waits for the Action to deliver.
type ExecuteActionResponse struct {
	Outcome IntentOutcome `json:"outcome"`
	Reason  string        `json:"reason,omitempty"`
	Result  ActionResult  `json:"result"`
}

// AvailableActionList is the Actions an architect may request: exactly the
// names in its hiveryn.yaml availableActions, in that order, each enriched
// from the library (without Suggestions, which are manual-launch only). A configured name with no directory in the library is
// listed invalid with a problem saying so, never dropped.
type AvailableActionList struct {
	Actions []ActionDefinition `json:"actions"`
}

// AddAvailableActionRequest is an architect's request to allow one more
// Action in its own project's hiveryn.yaml availableActions.
type AddAvailableActionRequest struct {
	Name string `json:"name"`
}

// AddAvailableActionResult reports whether the list changed (false when the
// name was already listed) and the resulting availableActions, in file order.
type AddAvailableActionResult struct {
	Changed          bool     `json:"changed"`
	AvailableActions []string `json:"available_actions"`
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

// ActionAttentionState says whether the Action agent is known to be waiting
// for the user at its terminal.
type ActionAttentionState string

const (
	// ActionAttentionInputRequired: an explicit provider signal — a hook
	// reporting a permission or question prompt, or a provider prompt
	// recognized on the agent's terminal screen — says the agent is waiting
	// for the user.
	ActionAttentionInputRequired ActionAttentionState = "input_required"
	// ActionAttentionNoneDetected: no such signal is present. This does not
	// prove the agent is not waiting: detection covers only known prompts.
	ActionAttentionNoneDetected ActionAttentionState = "none_detected"
	// ActionAttentionUnavailable: nothing can be inspected — the execution is
	// not running, or its agent terminal is not live.
	ActionAttentionUnavailable ActionAttentionState = "unavailable"
)

// ActionAttentionSource is where an input_required signal came from.
type ActionAttentionSource string

const (
	ActionAttentionSourceHook     ActionAttentionSource = "hook"
	ActionAttentionSourceTerminal ActionAttentionSource = "terminal"
)

// ActionAgentAttention reports whether the Action agent needs the user. It is
// separate from the execution status, which stays running, and from Activity:
// an idle or silent agent is never reported input_required. Reason is a
// stable code (a provider prompt such as folder_trust, or awaiting_input for a
// hook-reported prompt), Message says what the agent is asking and Since when
// it was first detected. Coverage states what the provider's detection can
// and cannot see, so none_detected is never read as "not waiting"; the
// execution's terminal, in the Actions window, is where the user checks and
// answers.
type ActionAgentAttention struct {
	State    ActionAttentionState  `json:"state"`
	Reason   string                `json:"reason,omitempty"`
	Message  string                `json:"message,omitempty"`
	Source   ActionAttentionSource `json:"source,omitempty"`
	Since    *time.Time            `json:"since,omitempty"`
	Coverage string                `json:"coverage,omitempty"`
}

// ActionResult is an architect's view of one requested execution, addressed by
// the execution id executeAction returned. The execution record is
// authoritative: pending_approval means nothing has started (StartedAt unset),
// running means the Action agent is working, and completed/failed/denied are
// final. ElapsedSeconds counts from StartedAt to EndedAt (or now) and is unset
// until it starts. Reason is the user's denial reason, Error the reason it
// failed without or before an agent conclusion, Summary the agent's
// conclusion. Machine is where the Action runs (empty means local). OutputDir
// is set only when Artifacts is ready: the package's path on the caller's
// machine.
type ActionResult struct {
	ExecutionID    string               `json:"execution_id"`
	Action         string               `json:"action"`
	Status         ActionRunStatus      `json:"status"`
	Prompt         string               `json:"prompt"`
	ProfileName    string               `json:"profile_name,omitempty"`
	Machine        string               `json:"machine,omitempty"`
	RequestedAt    time.Time            `json:"requested_at"`
	StartedAt      *time.Time           `json:"started_at,omitempty"`
	EndedAt        *time.Time           `json:"ended_at,omitempty"`
	ElapsedSeconds *int64               `json:"elapsed_seconds,omitempty"`
	Reason         string               `json:"reason,omitempty"`
	Error          string               `json:"error,omitempty"`
	Summary        string               `json:"summary,omitempty"`
	OutputDir      string               `json:"output_dir,omitempty"`
	Artifacts      ActionArtifacts      `json:"artifacts"`
	Activity       ActionAgentActivity  `json:"activity"`
	Attention      ActionAgentAttention `json:"attention"`
}

// MaxActionWaitSeconds caps one waitForActionResult call; a caller that needs
// longer waits again.
const MaxActionWaitSeconds = 30

// ActionWaitResult is one bounded wait: the result at return, and whether the
// wait ended on a change (Changed) or at its timeout (TimedOut). A change is
// a status change or an attention change — input_required appearing, clearing
// or changing reason — measured against the result when the wait began, so a
// repeated wait on an unchanged condition blocks until its timeout. An
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
// (pending_approval), started (running) or ended (denied, completed, failed),
// or, while running, that its agent's attention changed. The stream has no
// backlog, so a client reconciles by refetching executions on (re)connect.
type ActionEvent struct {
	Type        string          `json:"type"`
	Action      string          `json:"action"`
	ExecutionID string          `json:"execution_id"`
	Status      ActionRunStatus `json:"status"`
	SessionID   string          `json:"session_id,omitempty"`
	At          time.Time       `json:"at"`
}
