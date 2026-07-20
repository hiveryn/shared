# @hiveryn/shared

Shared domain types for the Hiveryn ecosystem.

- Pure data types and simple exported error types only.
- No service interfaces, no runtime types, no internal daemon dependencies.
- Provides both Go package (`github.com/hiveryn/shared/domain`) and matching TypeScript definitions (`@hiveryn/shared/domain` via `domain/index.ts`).

## What belongs here

- Enums and string types: `SessionType`, `SessionCreatedBy`, `SessionRunStatus`, `SessionRunFailureReason`, `TicketStatus`
- Core entities: `Session`, `SessionRun`, `SessionEvent`, `SessionTab`, `Ticket`, `TicketSummary`, `TicketBoard`, `TicketConclusion`, etc.
- Param/result structs for operations: `Create*Params`, `Conclude*Params`, `Move*Params`, `Edit*Params`, etc.
- Supporting types: `CommitRef`, `ArchitectConclusion`, `ConclusionSummary`, `TerminalInfo`, `CreateTerminalParams`, `PreviewBrowserTabParams`, `BrowserTabInfo`, `AgentProfileSnapshot`
- Error types: `NotFoundError`, `ValidationError`, `ConflictError`, `InternalError` (and sentinel `ErrNotFound`)
- Wire types used across API/MCP: `CreateSessionRequest`, `AppendSessionEventParams`, etc.

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
