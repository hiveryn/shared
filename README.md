# @hiveryn/shared

Shared domain types for the Hiveryn ecosystem.

- Pure data types and simple exported error types only.
- No service interfaces, no runtime types, no internal daemon dependencies.
- Provides both Go package (`github.com/hiveryn/shared/domain`) and matching TypeScript definitions (`@hiveryn/shared/domain` via `domain/index.ts`).

## What belongs here

- Enums and string types: `SessionType`, `SessionCreatedBy`, `SessionRunStatus`, `SessionRunFailureReason`, `TicketStatus`
- Core entities: `Session`, `SessionRun`, `SessionEvent`, `SessionTab`, `Ticket`, `TicketSummary`, `TicketBoard`, `TicketConclusion`, etc.
- Param/result structs for operations: `Create*Params`, `Conclude*Params`, `Move*Params`, `Edit*Params`, etc.
- Supporting types: `CommitRef`, `ArchitectConclusion`, `ConclusionSummary`, `TerminalInfo`, `CreateTerminalParams`, `AgentProfileSnapshot`
- Error types: `NotFoundError`, `ValidationError`, `ConflictError`, `InternalError` (and sentinel `ErrNotFound`)
- Wire types used across API/MCP: `CreateSessionRequest`, `AppendSessionEventParams`, etc.
- Stream contracts consumed by more than one repo: `ArchitectEvent` / `ArchitectEventReason` (the architect SSE stream — its `session_started` / `session_ended` reasons are the only place a session id is announced, so the desktop can discover sessions it did not create)
- Intent approval contract: `Intent` (an approvable agent tool call, blocking or deferred by its `IntentPolicy`; `manual` is deferred — no timer, only the user resolves it) and its optional approval-input schema — `IntentInputField` / `IntentInputType` (`text`, `textarea`, `choice`, `boolean`) / `IntentInputOption`, `IntentInputValues`, `IntentInputIssue` and `ApproveIntentRequest` (the desktop's `{"inputs": {...}}` approve body). Input-bearing intents are always deferred and defaults only prefill the form. `DeferredIntent` / `DeferredIntentStatus` (`pending_approval`, `denied`, `running`, `completed`, `failed`) is the durable outcome of a deferred intent, returned immediately with its stable id and retrievable by that id afterwards; approval (`approved_at`) and execution success (`completed`) are distinct.
- Actions contract: `ActionDefinition` / `ActionList` / `ActionProblem` (the global library under `HIVERYN_HOME/actions`, invalid definitions listed with their problems; optional action.yaml `suggestions` are manual-launch prompt suggestions bounded by `MaxActionSuggestions`/`MaxActionSuggestionLength` and never included in `AvailableActionList`), `ActionRun` / `ActionRunStatus` (`pending_approval`, `denied`, `running`, `completed`, `failed`) / `ActionRunTrigger` (the durable execution record, addressed by its stable execution id, carrying the output directory as structured data), `LaunchActionRequest` / `LaunchActionResult` (manual launch: prompt + variant), `ConcludeActionRequest` / `ActionConclusionOutcome` / `ActionConclusion` (the action agent's conclusion and its readback), `ActionEvent` (the `action_changed` stream) and `SessionTypeAction` (`action` sessions have no architect key; `context_id` is the execution id). Agent requests (architects and ticket workers): `IntentTypeExecuteAction` (a deferred intent whose id is the execution id), `ExecuteActionRequest`, `AvailableActionList` (the project's hiveryn.yaml `availableActions`, missing definitions listed invalid), `ActionResult` / `ActionAgentActivity` / `ActionWaitResult` (project-scoped result and bounded wait, capped at `MaxActionWaitSeconds`), `ActionRunTriggerArchitect` / `ActionRunTriggerWorker` (a worker request also records `requester_ticket_id`), `AddAvailableActionRequest` / `AddAvailableActionResult` (architect-only registration in its own hiveryn.yaml), and `ValidActionName`, shared by the library and the config.
- Session launch contract: `CreateSessionRequest.workflows` (explicit canonical workflow paths for a ticket session; empty is valid) and `Session.workflows` (the persisted selection, retained across runs and resumes).
- Architect workspace read contracts: `WorkspaceReport` / `WorkspaceNode` / `WorkspaceEntry` / `WorkspaceDiagnostic` (the structural check of the file-based workspace), `ArtifactKind` / `ArtifactSchema` / `ArtifactField` (the artifact schemas), `Workflow` / `WorkflowList` / `WorkflowAttach` (flat `workflows/*.md` discovery and repo applicability), and `WorkerPreflight` (whether the workspace's required project context can host a worker right now — the launch's own validation, narrower than `WorkspaceReport.valid`). The daemon produces them for both the architect MCP tools and the desktop, so they are shared rather than daemon-local.

## What does NOT belong here

- Service interfaces (`SessionService`, `TicketService`, `SessionRepository`, ...)
- Anything with runtime dependencies (terminals, MCP clients, stores, etc.)
- Internal-only helpers or daemon-specific envelopes (those stay in the using repo)

## Usage

### Go (daemon, future Go consumers)

In the consumer repo's `go.mod`, add a replace directive during early development:

```go
require github.com/hiveryn/shared v0.0.0

replace github.com/hiveryn/shared => ../shared
```

Import:

```go
import "github.com/hiveryn/shared/domain"

var _ = domain.SessionTypeArchitect
e := &domain.NotFoundError{Resource: "ticket", ID: "123"}
```

### TypeScript (desktop, tabplugin, future TS consumers)

The package is private and consumed via workspace or local path. In the consumer's `package.json`:

```json
{
  "dependencies": {
    "@hiveryn/shared": "workspace:*"
  }
}
```

Or with a path reference during early dev.

Import:

```ts
import {
  SessionType,
  TicketStatus,
  NotFoundError,
  type Session,
} from "@hiveryn/shared/domain";
```

For direct file import without package plumbing (early dev):

```ts
import type * as D from "../shared/domain";
```

## Versioning

Early development: breaking changes are expected and encouraged. Do not accumulate migration debt. If a cleaner model requires deleting all data, do it.

## Consumers (build-time)

- daemon (Go)
- tabplugin (TS)
- desktop (TS)
