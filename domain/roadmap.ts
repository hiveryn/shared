import type { TicketOutcome, TicketStatus } from "./ticket";

export const ROADMAP_SCHEMA_VERSION = 1;

export type RoadmapItemKind = "goal" | "initiative" | "milestone";

export type RoadmapItemStatus = "planned" | "active" | "blocked" | "done";

export interface RoadmapItem {
  id: string;
  kind: RoadmapItemKind;
  title: string;
  status: RoadmapItemStatus;
  outcome: string;
  parent_id: string | null;
  order: number;
  success_criteria: string[];
  tickets: string[];
  depends_on: string[];
}

export interface Roadmap {
  schema_version: number;
  title: string;
  items: RoadmapItem[];
}

export interface RoadmapArchiveEntry {
  root_id: string;
  archived_at: string;
  summary?: string;
  original_parent_id: string | null;
  original_order: number;
  items: RoadmapItem[];
}

export interface RoadmapArchive {
  schema_version: number;
  entries: RoadmapArchiveEntry[];
}

export interface RoadmapArchiveEntrySummary {
  root_id: string;
  root_title: string;
  root_kind: RoadmapItemKind;
  archived_at: string;
  summary?: string;
  item_count: number;
}

export interface RoadmapTicketInfo {
  id: string;
  title: string;
  status: TicketStatus;
  repo?: string;
  additional_repos: string[];
  has_conclusion: boolean;
  conclusion_outcome?: TicketOutcome;
}

export type RoadmapOpType =
  | "create"
  | "update"
  | "move"
  | "link_ticket"
  | "unlink_ticket"
  | "archive"
  | "restore";

export interface RoadmapOp {
  type: RoadmapOpType;
  id?: string;
  kind?: RoadmapItemKind;
  title?: string;
  status?: RoadmapItemStatus;
  outcome?: string;
  success_criteria?: string[];
  depends_on?: string[];
  parent_id?: string;
  order?: number;
  ticket_id?: string;
  summary?: string;
}

export interface UpdateRoadmapParams {
  version: string;
  title?: string;
  ops: RoadmapOp[];
}

export interface RoadmapView {
  view: string;
  version: string;
  title?: string;
  items: RoadmapItem[];
  archive_entries?: RoadmapArchiveEntrySummary[];
  archive_entry?: RoadmapArchiveEntry;
  tickets: RoadmapTicketInfo[];
  warnings: string[];
}

export interface RoadmapOpResult {
  index: number;
  type: RoadmapOpType;
  item_id: string;
  detail: string;
}

export interface RoadmapUpdateResult {
  version: string;
  applied: RoadmapOpResult[];
  tickets: RoadmapTicketInfo[];
  warnings: string[];
}
