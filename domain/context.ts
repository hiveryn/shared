import type { ReferenceLocation } from "./ticket";

/** Outcome of a worker context read (readContext / listContext); see context.go. */
export type ContextReadStatus =
  | "ok"
  | "missing"
  | "permission_denied"
  | "not_a_file"
  | "not_a_directory"
  | "binary"
  | "too_large"
  | "changed"
  | "unreachable"
  | "outside_root";

export const MAX_CONTEXT_FILE_BYTES = 4 << 20;
export const DEFAULT_CONTEXT_LINES = 400;
export const MAX_CONTEXT_LINES = 2000;
export const MAX_CONTEXT_PAGE_BYTES = 256 << 10;
export const DEFAULT_CONTEXT_ENTRIES = 200;
export const MAX_CONTEXT_ENTRIES = 1000;

export interface ContextReadRequest {
  reference: string;
  /** 1-based first line; default 1. */
  offset?: number;
  limit?: number;
  expected_sha256?: string;
}

export interface ContextTarget {
  reference: string;
  owner: string;
  location: ReferenceLocation;
  machine?: string;
}

export interface ContextFile extends ContextTarget {
  status: ContextReadStatus;
  message?: string;
  sha256?: string;
  size?: number;
  modified_at?: string;
  total_lines?: number;
  start_line?: number;
  end_line?: number;
  content?: string;
  /** Next page's offset; absent when this page reaches the end. */
  next_offset?: number;
}

export interface ContextListRequest {
  reference: string;
  /** 0-based entry index; default 0. */
  offset?: number;
  limit?: number;
  expected_sha256?: string;
}

export type ContextEntryType = "file" | "directory" | "symlink" | "other";

export interface ContextEntry {
  name: string;
  reference: string;
  type: ContextEntryType;
  size?: number;
  modified_at?: string;
}

export interface ContextListing extends ContextTarget {
  status: ContextReadStatus;
  message?: string;
  sha256?: string;
  total_entries?: number;
  offset?: number;
  entries?: ContextEntry[];
  next_offset?: number;
}
