import type { WorkspaceDiagnostic } from './workspace';

/** A plain Markdown file selected explicitly for a worker run. */
export interface Workflow {
 name: string;
 path: string;
 rel_path: string;
 valid: boolean;
 modified_at: string | null;
 diagnostics: WorkspaceDiagnostic[];
}

export interface WorkflowList {
 architect_key: string;
 workflows: Workflow[];
 diagnostics: WorkspaceDiagnostic[];
}
