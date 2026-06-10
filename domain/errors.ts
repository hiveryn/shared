export type ErrorCode = "VALIDATION" | "CONFLICT" | "NOT_FOUND" | "INTERNAL";

export const ErrNotFound = new Error("resource not found");

export class ValidationError extends Error {
  readonly field: string;
  readonly message: string;
  constructor(field: string, message: string) {
    super(`${field} ${message}`);
    this.name = "ValidationError";
    this.field = field;
    this.message = message;
  }
}

export class ConflictError extends Error {
  readonly resource: string;
  readonly field: string;
  readonly message: string;
  constructor(resource: string, field: string, message: string) {
    super(`${resource} ${field}: ${message}`);
    this.name = "ConflictError";
    this.resource = resource;
    this.field = field;
    this.message = message;
  }
}

export class NotFoundError extends Error {
  readonly resource: string;
  readonly id: string;
  constructor(resource: string, id: string) {
    super(`${resource} ${id} not found`);
    this.name = "NotFoundError";
    this.resource = resource;
    this.id = id;
  }
}

export class InternalError extends Error {
  constructor(message: string) {
    super(message);
    this.name = "InternalError";
  }
}
