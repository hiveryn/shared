import type { CommitRef } from "./commit";

export type TicketStatus = "backlog" | "progress" | "done";

export interface TicketWarning {
  code: string;
  message: string;
}

/** Namespace a reference addresses (see reference.go ParseReference). */
export type ReferenceKind = "ticket" | "architect" | "repo";

/** Daemon verdict on a reference; unverified = could not be inspected, invalid = not a valid reference here. */
export type ReferenceStatus = "present" | "missing" | "unverified" | "invalid";

/** Where a reference's target lives; unknown is never local. */
export type ReferenceLocation = "local" | "remote" | "unknown";

/** The worker MCP tool that reads a reference. */
export type ReferenceAccess = "readTicket" | "readContext" | "listContext";

export type PathKind = "file" | "directory";

export const ARCHITECT_REFERENCE_PREFIX = "architect";

/**
 * Daemon-resolved view of one `references` entry. `value` is the exact string
 * in ticket.md; `exists` mirrors `status === "present"`. Architect and repo
 * references carry owner/target/location (machine when remote) and access.
 */
export interface TicketReference {
  value: string;
  /** Absent for an invalid reference. */
  kind?: ReferenceKind;
  status: ReferenceStatus;
  exists: boolean;
  /** `architect` or the repo key. */
  owner?: string;
  /** Path relative to the owner's root; `.` is the root. */
  target?: string;
  location?: ReferenceLocation;
  machine?: string;
  path_kind?: PathKind;
  access?: ReferenceAccess;
  /** Why the reference is invalid or unverified. */
  error?: string;
}

export interface TicketSummary {
  id: string;
  status: TicketStatus;
  title: string;
  repo?: string;
  additional_repos: string[];
  created?: string;
  updated?: string;
  references: string[];
  resolved_references: TicketReference[];
  has_conclusion: boolean;
  warnings: TicketWarning[];
}

export type TicketOutcome = "completed" | "exploratory" | "rejected";

export interface TicketConclusion {
  started_at: string;
  concluded_at: string;
  agent?: string;
  profile?: string;
  outcome: TicketOutcome;
  rejection_reason?: string;
  commits: CommitRef[];
  body: string;
}

export interface Ticket extends TicketSummary {
  body: string;
  conclusion: TicketConclusion | null;
}

export interface TicketBoard {
  backlog: TicketSummary[];
  progress: TicketSummary[];
  done: TicketSummary[];
}

export interface CreateTicketParams {
  title: string;
  repo: string;
  additionalRepos: string[];
  body: string;
  references: string[];
  now: string;
}

export interface EditTicketParams {
  oldString: string;
  newString: string;
  replaceAll: boolean;
  now: string;
}

export interface UpdateTicketMetadataParams {
  title?: string;
  repo?: string;
  additionalRepos?: string[];
  references?: string[];
  now: string;
}

export interface MoveTicketParams {
  to: TicketStatus;
}
