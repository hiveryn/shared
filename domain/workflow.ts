import type { WorkspaceDiagnostic } from "./workspace";

/**
 * How a workflow reaches a worker session. There are exactly two modes and no
 * third "always attached" category: manual workflows are only ever added by
 * hand, suggested workflows are preselected when their repos overlap the
 * session's writable scope. A repo match is a suggestion, never an obligation.
 */
export type WorkflowAttach = "manual" | "suggested";

/**
 * One discovered `workflows/*.md` file. Workflows are flat — immediate children
 * of `workflows/` only — have no dependencies on each other, and carry no
 * precedence: a semantic conflict between two workflow bodies is an agent/user
 * decision, not something this contract resolves.
 *
 * `path` is the canonical absolute path inside the architect workspace, which
 * is what a session records. Workflows are never copied or snapshotted.
 *
 * An invalid workflow is reported with its diagnostics rather than dropped, so
 * the problem is repairable; it is excluded from `suggested` regardless of repo
 * overlap.
 */
export interface Workflow {
  name: string;
  path: string;
  rel_path: string;
  attach: WorkflowAttach;
  repos: string[];
  valid: boolean;
  suggested: boolean;
  modified_at: string | null;
  diagnostics: WorkspaceDiagnostic[];
}

/**
 * The workflow read contract. `scope_repos` echoes the writable repo scope the
 * `suggested` flags were computed against; an empty scope yields no suggestions
 * at all rather than suggesting everything.
 */
export interface WorkflowList {
  architect_key: string;
  scope_repos: string[];
  workflows: Workflow[];
  diagnostics: WorkspaceDiagnostic[];
}
