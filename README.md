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
- Intent approval contract: `Intent` (an approvable agent tool call; the agent's call blocks until the user answers or its tool's `IntentPolicy` resolves it when the wait window expires), `IntentOrigin` and the `IntentOutcome` catalog.
- Agent notification and question contract: `NotifyRequest` / `NotifyResult`, `AskQuestionRequest` (validated by `Normalize`: question ≤ `MaxQuestionLength`, `MinQuestionAnswers`–`MaxQuestionAnswers` distinct suggested answers ≤ `MaxQuestionAnswerLength`, and a required 0-based `recommended_index` into them), `AgentQuestion` / `QuestionOrigin` / `QuestionStatus` (`pending`, `answered`, `expired`, `cancelled`, `interrupted`; carried on the session stream as `question` `required`/`resolved` events), `AnswerQuestionRequest` (`NormalizeQuestionResponse`, ≤ `MaxQuestionResponseLength`), `AskQuestionResult`, `QuestionTimeout` (1h) and the exact agent replies `NotificationSentText` / `QuestionTimeoutText`. Answer provenance: `QuestionAnswerSource` (`user`, or `architect` when the owning architect answered an escalated ticket-worker question through `answerQuestion`) on `AgentQuestion.answered_by` / `AskQuestionResult.answered_by`, the worker-facing `QuestionArchitectAnswerPrefix`, and the architect's `ArchitectAnswerResult` / `ArchitectAnswerOutcome` (`answered`, `not_answerable` with its reason — a settled or unescalated question is a result, not an error). TypeScript mirrors the bounds, strings and `normalizeQuestionResponse`. A question is not an intent, and its recommendation is never selected for the user; an architect answer is never presented as the user's.
- Agent message contract: `AgentMessage` / `AgentMessageKind` (`ticket_concluded`, `question_escalated` — a ticket worker's question unanswered by the user after the daemon's configured delay, addressed to the project architect) / `AgentMessageState` (`pending_commit`, `pending`, `handed`, `appeared`, `acknowledged`, `failed`) / `AgentMessageList`, the daemon's durable outbox of text delivered into agent sessions through each provider's native route. `recipient` is `architect:<key>` (whichever session of that architect runs) or `session:<id>`; `condition` says why an open message is not progressing; `handed` is transport acceptance, `appeared` the provider's own record, `acknowledged` an agent reply after it — none means the agent acted on it. The architect stream reason `agent_message_changed` announces a change without naming the message.
- Session launch contract: `CreateSessionRequest.workflows` (explicit canonical workflow paths for a ticket session; empty is valid) and `Session.workflows` (the persisted selection, retained across runs and resumes). Architect-requested launches: `IntentTypeSpawnTicketWorker` (a blocking wait-then-allow intent), `SpawnTicketWorkerRequest` (ticket id, the required variant and workflow names — `NAME` or `NAME.md`, never paths; the list is the whole selection) / `SpawnTicketWorkerResponse` (intent outcome + `SpawnedTicketWorker`, the launched session identity, present only when the worker actually launched), and `SessionCreatedByArchitect` on sessions launched that way.
- Context reference contract: `ParseReference` (the one strict grammar — `<ticket-id>`, `architect:<relative path>` for the owning architect workspace on the daemon host, `<repoKey>:<relative path>` for a configured repo on its own machine; absolute, unqualified, URL, traversal and non-canonical paths are `ReferenceSyntaxError`s with remediation; `architect` is reserved and `ValidRepoKey` is the repo-key grammar configuration enforces) and `TicketReference` (`resolved_references`): every entry has a `ReferenceStatus` (`present`, `missing` — confirmed absent — `unverified` with its `error`, or `invalid` with why and how to fix it); architect/repo entries carry `owner`, `target`, `ReferenceLocation` (`local`, `remote` with its `machine`, `unknown` — never local), `path_kind` and `access` (the worker MCP tool that reads it: `readTicket`, `readContext`, `listContext`). Worker context reads: `ContextReadRequest`/`ContextFile` and `ContextListRequest`/`ContextListing`/`ContextEntry` with the closed `ContextReadStatus` catalog and the page/size bounds (`MaxContextFileBytes`, `MaxContextLines`, `MaxContextPageBytes`, `MaxContextEntries`); content identity is `sha256` + `size` + `modified_at`, and `expected_sha256` turns a change between pages into `changed`.
- Conversation API contract: `ConversationCapabilities`, `Submission`/`SubmissionRequest`, `ConversationItem`/`ConversationPage`. A `Submission` (receipt) is only returned by the send itself — there is no receipt list/read or archived-chat discovery contract (one connected client reconciles by repeating its request id). User input is delivered synchronously, never queued by Hiveryn: a `Submission` `status` is `delivering`, then `rejected` (a definitive `error_code`; never delivered later), `unknown` (`DELIVERY_UNKNOWN`: a hand-off was attempted and acceptance is not established; reconciled, never resent) or `accepted`, then `running`/`blocked`/`completed`/`failed`/`uncertain`; `queued` only names the chat owner's internal kickoff/notification entries. `accepted_at`/`acceptance` record provider acceptance separately from turn completion. A `ConversationItem` `via` of `kickoff` (with `origin: kickoff` on its turn) is the recorded launch message, shown as conversation content. `DeliveryError` (with the `Delivery*` codes) carries a send that did not establish acceptance and its durable receipt.
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
