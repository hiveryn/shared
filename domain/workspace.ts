/**
 * One file-based architect workspace artifact. Tickets and conclusions are
 * deliberately absent: they keep their own schemas and MCP lifecycle and are
 * never described or validated by the workspace surface.
 */
export type ArtifactKind =
  | "PROJECT_OVERVIEW"
  | "PROJECT_STATE"
  | "ROADMAP_CURRENT"
  | "ARCHITECT_SYSTEM"
  | "WORKFLOW"
  | "HIVERYN_YAML";

/** The complete, ordered catalog. describeArtifact accepts exactly these. */
export const ARTIFACT_KINDS: readonly ArtifactKind[] = [
  "HIVERYN_YAML",
  "PROJECT_OVERVIEW",
  "PROJECT_STATE",
  "ROADMAP_CURRENT",
  "ARCHITECT_SYSTEM",
  "WORKFLOW",
];

/**
 * Separates blocking structural problems from advisory ones. Only `error`
 * contributes to `WorkspaceReport.valid`.
 */
export type DiagnosticSeverity = "error" | "warning";

/**
 * One actionable finding. `path` is workspace-relative (empty for
 * workspace-wide findings). `line` is the 1-indexed source line when the
 * finding maps to one, and 0 otherwise — an absent field has no line.
 */
export interface WorkspaceDiagnostic {
  code: string;
  severity: DiagnosticSeverity;
  path: string;
  line: number;
  message: string;
}

export type WorkspaceNodeType = "file" | "directory";

/**
 * One entry of the workspace's expected shape: a root document, the config, or
 * the workflows directory. Expected nodes are always present in the report
 * even when the file is missing, so a broken workspace stays inspectable and
 * repairable.
 *
 * `document_updated_at` is the timestamp the document itself declares
 * (lastUpdatedAt); `modified_at` is the filesystem mtime. They are
 * reported separately and neither certifies factual freshness.
 */
export interface WorkspaceNode {
  kind: ArtifactKind;
  type: WorkspaceNodeType;
  path: string;
  required: boolean;
  exists: boolean;
  valid: boolean;
  document_updated_at: string | null;
  modified_at: string | null;
  children: WorkspaceEntry[];
  diagnostics: WorkspaceDiagnostic[];
}

/**
 * One workflow file discovered inside the workflows directory node.
 *
 * It has no children of its own, and that is a property of the workspace rather
 * than a simplification: workflows/ is flat by rule, so a subdirectory inside it
 * is reported as an error on the directory and never descended into.
 *
 * It carries no `exists` or `required`: a discovered file exists by definition,
 * and no individual workflow is required.
 */
export interface WorkspaceEntry {
  kind: ArtifactKind;
  path: string;
  valid: boolean;
  document_updated_at: string | null;
  modified_at: string | null;
  diagnostics: WorkspaceDiagnostic[];
}

export interface WorkspaceTicketTotals {
  backlog: number;
  progress: number;
  done: number;
}

/**
 * One warning raised by the ticket system for one backlog ticket. Advisory
 * only — ticket warnings never make a workspace invalid.
 */
export interface WorkspaceTicketWarning {
  ticket_id: string;
  code: string;
  message: string;
}

export interface WorkspaceTickets {
  totals: WorkspaceTicketTotals;
  warnings: WorkspaceTicketWarning[];
}

/**
 * The full result of a parameterless workspace check. Read-only, deterministic
 * and mutates nothing: node and diagnostic order is stable across runs for the
 * same on-disk state.
 *
 * `valid` is true only when no error-severity diagnostic appears anywhere in
 * the report. It is a structural verdict, not an authorization and not a claim
 * about the truth of any document's contents.
 */
export interface WorkspaceReport {
  architect_key: string;
  workspace_path: string;
  checked_at: string;
  valid: boolean;
  nodes: WorkspaceNode[];
  tickets: WorkspaceTickets;
  diagnostics: WorkspaceDiagnostic[];
}

/** The value shape the validator actually enforces for a field. */
export type ArtifactFieldType =
  | "string"
  | "string-list"
  | "string-map"
  | "enum"
  | "rfc3339-utc-datetime"
  | "mapping"
  | "mapping-list";

export interface ArtifactField {
  name: string;
  required: boolean;
  type: ArtifactFieldType;
  enum?: string[];
  description: string;
}

/**
 * describeArtifact's result. Rendered from the same definition the validators
 * execute, so a rule described here is a rule enforced by the workspace check.
 */
export interface ArtifactSchema {
  kind: ArtifactKind;
  schema_version: number;
  title: string;
  required: boolean;
  format: string;
  location: string;
  naming: string;
  fields: ArtifactField[];
  rules: string[];
  example: string;
}

/**
 * The answer to one question asked before a ticket session is created: is the
 * architect workspace's required project context ready for a worker right now?
 * It runs the project-context half of the worker launch validation — the same
 * code the launch itself runs — and reports what would block it.
 *
 * It is deliberately not `WorkspaceReport.valid`: a worker is blocked only by
 * `hiveryn.yaml`, `PROJECT_OVERVIEW.md`, `PROJECT_STATE.md` and an
 * invalid-when-present `ROADMAP_CURRENT.md`. Architect-only artifacts and
 * unselected invalid workflows never block one, so a workspace the
 * check calls invalid can still be launchable.
 *
 * Selected-workflow validity is not part of this answer: each workflow reports
 * its own validity and diagnostics in the workflow listing, and the launch
 * revalidates the whole selection against the live workspace regardless.
 *
 * `launchable` is false exactly when `problems` is non-empty; each problem is
 * one actionable finding, worded as the launch error would word it.
 */
export interface WorkerPreflight {
  architect_key: string;
  checked_at: string;
  launchable: boolean;
  problems: string[];
}
