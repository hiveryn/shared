/**
 * The only type carried on the architect event stream
 * (GET /api/architects/{key}/events). It exists as a discriminator so clients
 * can ignore anything they synthesize locally onto the same channel.
 */
export const ARCHITECT_EVENT_TYPE = "workspace_changed";

/**
 * What actually changed. Ticket reasons drive the kanban board; session reasons
 * drive the set of session tabs in an architect window. A client that only knows
 * about the ticket reasons keeps working.
 *
 * `session_started` announces a session whose run is live (running run + main
 * terminal). It is the only signal that names a session id, and therefore the
 * only way a client can discover a session it did not create itself — an
 * architect MCP spawn, or a spawn from another window. `session_ended` announces
 * a session that is gone (concluded or discarded). `roadmap_updated` announces a
 * successful roadmap mutation (any architect MCP roadmap tool call, or a future
 * desktop-originated write).
 */
export type ArchitectEventReason =
  | "ticket_created"
  | "ticket_updated"
  | "ticket_moved"
  | "ticket_deleted"
  | "ticket_concluded"
  | "session_started"
  | "session_ended"
  | "roadmap_updated";

/**
 * Published per architect key and fanned out to every subscriber of that
 * architect's event stream. The stream has no backlog, so a client must
 * reconcile on (re)connect rather than treat delivery as guaranteed.
 */
export interface ArchitectEvent {
  type: typeof ARCHITECT_EVENT_TYPE;
  architect_key: string;
  /** Empty for reasons that are not about a ticket. */
  ticket_id: string;
  reason: ArchitectEventReason;
  /**
   * Set for the session reasons, and for ticket reasons caused by a session
   * (conclude, discard). Empty for plain ticket edits.
   */
  session_id: string;
  at: string;
}
