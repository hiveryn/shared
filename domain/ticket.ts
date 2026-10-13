import type { CommitRef } from "./commit";

export type TicketStatus = "backlog" | "progress" | "done";

export interface TicketWarning {
  code: string;
  message: string;
}

/** Outcome of inspecting a path reference where it lives; unverified means it could not be inspected. */
export type PathReferenceStatus = "present" | "missing" | "unverified";

/** Where a path reference was interpreted; unknown (primary-repo machine undetermined) is never local. */
export type PathReferenceLocation = "local" | "remote" | "unknown";

export interface TicketIdReference {
  value: string;
  type: "ticket";
  exists: boolean;
}

export interface PathReference {
  value: string;
  type: "path";
  /** Mirrors `status === "present"`. */
  exists: boolean;
  kind?: "file" | "directory";
  status: PathReferenceStatus;
  location: PathReferenceLocation;
  /** The machine name, when `location` is remote. */
  machine?: string;
  /** Why the path could not be inspected, when `status` is unverified. */
  error?: string;
}

export type TicketReference = TicketIdReference | PathReference;

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
