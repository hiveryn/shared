import type { CommitRef } from "./commit";

export type TicketStatus = "backlog" | "progress" | "done";

export interface TicketWarning {
  code: string;
  message: string;
}

export interface TicketSummary {
  id: string;
  status: TicketStatus;
  title: string;
  repo?: string;
  created?: string;
  updated?: string;
  references: string[];
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
  references?: string[];
  now: string;
}

export interface MoveTicketParams {
  to: TicketStatus;
}
