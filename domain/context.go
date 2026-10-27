package domain

import "time"

// Worker context reads: readContext / listContext. A session reads its own
// architect's workspace and that architect's configured repos, addressed with
// architect:<path> or <repoKey>:<path> (ParseReference). The scope comes from
// the calling session, never from the request. Reads return the current
// contents at the time of the call, not a launch snapshot; ExpectedSHA256
// detects a change between pages.
//
// Every outcome is a result with a Status, never an empty success: anything
// but ok carries a Message saying what happened and what to do.
type ContextReadStatus string

const (
	ContextOK               ContextReadStatus = "ok"
	ContextMissing          ContextReadStatus = "missing"
	ContextPermissionDenied ContextReadStatus = "permission_denied"
	// ContextNotFile: the target is a directory (or other non-regular file);
	// list it with listContext.
	ContextNotFile ContextReadStatus = "not_a_file"
	// ContextNotDirectory: the target is a file; read it with readContext.
	ContextNotDirectory ContextReadStatus = "not_a_directory"
	ContextBinary       ContextReadStatus = "binary"
	ContextTooLarge     ContextReadStatus = "too_large"
	// ContextChanged: the content no longer has ExpectedSHA256; restart
	// from the first page with the new hash.
	ContextChanged ContextReadStatus = "changed"
	// ContextUnreachable: the owning machine could not be reached. Existence
	// is unknown; it is never reported as missing.
	ContextUnreachable ContextReadStatus = "unreachable"
	// ContextOutsideRoot: the target resolves (through a symlink) outside its
	// root and is not readable as context.
	ContextOutsideRoot ContextReadStatus = "outside_root"
)

const (
	// MaxContextFileBytes is the largest file readContext reads.
	MaxContextFileBytes = 4 << 20
	// DefaultContextLines / MaxContextLines bound a readContext page.
	DefaultContextLines = 400
	MaxContextLines     = 2000
	// MaxContextPageBytes bounds a page's content; a page ends early, on a
	// line boundary, when it would exceed it (a single longer line is cut).
	MaxContextPageBytes = 256 << 10
	// DefaultContextEntries / MaxContextEntries bound a listContext page.
	DefaultContextEntries = 200
	MaxContextEntries     = 1000
)

// ContextReadRequest asks for lines [Offset, Offset+Limit) of a file, 1-based.
// Zero values mean line 1 and DefaultContextLines.
type ContextReadRequest struct {
	Reference      string `json:"reference"`
	Offset         int    `json:"offset,omitempty"`
	Limit          int    `json:"limit,omitempty"`
	ExpectedSHA256 string `json:"expected_sha256,omitempty"`
}

// ContextTarget is where a context reference was resolved.
type ContextTarget struct {
	Reference string            `json:"reference"`
	Owner     string            `json:"owner"`
	Location  ReferenceLocation `json:"location"`
	Machine   string            `json:"machine,omitempty"`
}

// ContextFile is one page of a text file. SHA256, Size and ModifiedAt
// identify the whole file's content; NextOffset is the next page's Offset, 0
// when this page reaches the end.
type ContextFile struct {
	ContextTarget
	Status     ContextReadStatus `json:"status"`
	Message    string            `json:"message,omitempty"`
	SHA256     string            `json:"sha256,omitempty"`
	Size       int64             `json:"size,omitempty"`
	ModifiedAt *time.Time        `json:"modified_at,omitempty"`
	TotalLines int               `json:"total_lines,omitempty"`
	StartLine  int               `json:"start_line,omitempty"`
	EndLine    int               `json:"end_line,omitempty"`
	Content    string            `json:"content,omitempty"`
	NextOffset int               `json:"next_offset,omitempty"`
}

// ContextListRequest asks for entries [Offset, Offset+Limit) of a directory,
// 0-based, in name order. A zero Limit means DefaultContextEntries.
type ContextListRequest struct {
	Reference      string `json:"reference"`
	Offset         int    `json:"offset,omitempty"`
	Limit          int    `json:"limit,omitempty"`
	ExpectedSHA256 string `json:"expected_sha256,omitempty"`
}

type ContextEntryType string

const (
	ContextEntryFile      ContextEntryType = "file"
	ContextEntryDirectory ContextEntryType = "directory"
	// ContextEntrySymlink is listed, not followed; reading it applies the
	// same containment check as any reference.
	ContextEntrySymlink ContextEntryType = "symlink"
	ContextEntryOther   ContextEntryType = "other"
)

// ContextEntry is one directory entry; Reference addresses it directly.
type ContextEntry struct {
	Name       string           `json:"name"`
	Reference  string           `json:"reference"`
	Type       ContextEntryType `json:"type"`
	Size       int64            `json:"size,omitempty"`
	ModifiedAt *time.Time       `json:"modified_at,omitempty"`
}

// ContextListing is one page of a directory's immediate entries (never
// recursive). SHA256 identifies the whole listing (names, types, sizes and
// modification times).
type ContextListing struct {
	ContextTarget
	Status       ContextReadStatus `json:"status"`
	Message      string            `json:"message,omitempty"`
	SHA256       string            `json:"sha256,omitempty"`
	TotalEntries int               `json:"total_entries,omitempty"`
	Offset       int               `json:"offset,omitempty"`
	Entries      []ContextEntry    `json:"entries,omitempty"`
	NextOffset   int               `json:"next_offset,omitempty"`
}
