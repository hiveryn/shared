package domain

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

// Context reference grammar. A ticket's `references` entry, and the
// `reference` argument of the worker context tools, is exactly one of:
//
//	<ticket-id>              a ticket on the same board
//	architect:<relative>     a file or directory in the owning architect workspace
//	<repoKey>:<relative>     a file or directory in that configured repo, on its own machine
//
// <relative> is a canonical slash-separated path below the root: no leading
// "/", "~", backslash, control character, empty, "." or ".." segment and no
// trailing "/". The single segment "." names the root itself. Nothing else is
// a reference: absolute and unqualified paths are rejected rather than
// interpreted, so a reference always says whose filesystem it addresses.
const (
	// ArchitectReferencePrefix names the owning architect workspace. It is
	// reserved: no repository may use it as its key.
	ArchitectReferencePrefix = "architect"
	// ReferenceRootTarget is the target naming a root itself.
	ReferenceRootTarget = "."
)

var (
	ticketIDGrammar   = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}-\d{4}-[A-Za-z0-9][A-Za-z0-9._-]*$`)
	repoKeyGrammar    = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._-]*$`)
	ticketIDLikeStart = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}-\d{4}-`)
)

// ValidRepoKey reports whether key can be addressed as a reference prefix.
// Configuration enforces it for every repo key, and rejects the reserved
// ArchitectReferencePrefix separately.
func ValidRepoKey(key string) bool { return repoKeyGrammar.MatchString(key) }

// ReferenceKind is which namespace a parsed reference addresses.
type ReferenceKind string

const (
	ReferenceKindTicket    ReferenceKind = "ticket"
	ReferenceKindArchitect ReferenceKind = "architect"
	ReferenceKindRepo      ReferenceKind = "repo"
)

// ParsedReference is the syntax-only reading of a reference. Whether a repo
// key is configured, where it lives and whether the target exists are the
// daemon's resolution, not the parser's.
type ParsedReference struct {
	Kind ReferenceKind
	// TicketID is set for ticket references.
	TicketID string
	// Owner is ArchitectReferencePrefix or the repo key.
	Owner string
	// Target is the canonical relative path, or ReferenceRootTarget.
	Target string
}

// String renders the canonical form.
func (r ParsedReference) String() string {
	if r.Kind == ReferenceKindTicket {
		return r.TicketID
	}
	return r.Owner + ":" + r.Target
}

// ReferenceSyntaxError explains why a value is not a reference, with the
// remediation the author needs.
type ReferenceSyntaxError struct {
	Value  string
	Reason string
}

func (e *ReferenceSyntaxError) Error() string {
	return fmt.Sprintf("invalid reference %q: %s", e.Value, e.Reason)
}

const referenceForms = "use a same-board ticket id, architect:<path relative to the architect workspace> or <repo key>:<path relative to that repo>"

// ParseReference parses one reference strictly. Surrounding whitespace is
// not trimmed: a stored value must already be canonical.
func ParseReference(value string) (ParsedReference, error) {
	fail := func(reason string) (ParsedReference, error) {
		return ParsedReference{}, &ReferenceSyntaxError{Value: value, Reason: reason}
	}
	if strings.TrimSpace(value) == "" {
		return fail("empty reference; " + referenceForms)
	}
	if strings.TrimSpace(value) != value {
		return fail("leading or trailing whitespace; write " + strings.TrimSpace(value))
	}
	for _, r := range value {
		if unicode.IsControl(r) {
			return fail("contains a control character")
		}
	}
	if strings.HasPrefix(value, "/") || strings.HasPrefix(value, "~") {
		return fail("absolute and home-relative paths are not references; " + referenceForms)
	}
	owner, target, qualified := strings.Cut(value, ":")
	if !qualified {
		if ticketIDGrammar.MatchString(value) {
			return ParsedReference{Kind: ReferenceKindTicket, TicketID: value}, nil
		}
		if ticketIDLikeStart.MatchString(value) {
			return fail("not a valid ticket id (YYYY-MM-DD-HHMM-slug using letters, digits, '.', '_' or '-')")
		}
		return fail("unqualified path or unknown ticket id form; " + referenceForms)
	}
	switch strings.ToLower(owner) {
	case "http", "https", "file", "ssh":
		return fail("URLs are not references; " + referenceForms)
	}
	if owner == "" {
		return fail("missing prefix before ':'; " + referenceForms)
	}
	if !repoKeyGrammar.MatchString(owner) {
		return fail("prefix " + fmt.Sprintf("%q", owner) + " is not architect or a repo key (letters, digits, '.', '_', '-')")
	}
	if strings.EqualFold(owner, ArchitectReferencePrefix) && owner != ArchitectReferencePrefix {
		return fail("the architect namespace is spelled " + ArchitectReferencePrefix)
	}
	if err := checkReferenceTarget(target); err != "" {
		return fail(err)
	}
	kind := ReferenceKindRepo
	if owner == ArchitectReferencePrefix {
		kind = ReferenceKindArchitect
	}
	return ParsedReference{Kind: kind, Owner: owner, Target: target}, nil
}

func checkReferenceTarget(target string) string {
	switch {
	case target == "":
		return "empty target after ':'; name a path, or " + ReferenceRootTarget + " for the root"
	case target == ReferenceRootTarget:
		return ""
	case strings.HasPrefix(target, "/"):
		return "the target must be relative to the root (no leading '/')"
	case strings.HasPrefix(target, "~"):
		return "the target must be relative to the root (no '~')"
	case strings.Contains(target, `\`):
		return "use '/' as the separator (no backslash)"
	case strings.HasSuffix(target, "/"):
		return "remove the trailing '/' (directories are named like files)"
	}
	for segment := range strings.SplitSeq(target, "/") {
		switch segment {
		case "":
			return "empty path segment ('//')"
		case ".":
			return "remove './' segments"
		case "..":
			return "'..' segments are not allowed; references stay inside their root"
		}
	}
	return ""
}

// ReferenceStatus is the daemon's verdict on a reference: present, confirmed
// missing, unverified (could not be inspected; Error says why) or invalid
// (not a valid reference for this architect; Error says why and how to fix
// it).
type ReferenceStatus string

const (
	ReferencePresent    ReferenceStatus = "present"
	ReferenceMissing    ReferenceStatus = "missing"
	ReferenceUnverified ReferenceStatus = "unverified"
	ReferenceInvalid    ReferenceStatus = "invalid"
)

// ReferenceLocation is where a reference's target lives: the daemon's host,
// a remote machine (Machine names it) or unknown when it cannot be located.
// Unknown is never local.
type ReferenceLocation string

const (
	ReferenceLocal   ReferenceLocation = "local"
	ReferenceRemote  ReferenceLocation = "remote"
	ReferenceUnknown ReferenceLocation = "unknown"
)

// ReferenceAccess is how a worker reads a reference: the Hiveryn MCP tool to
// call with the reference value.
type ReferenceAccess string

const (
	ReferenceAccessReadTicket  ReferenceAccess = "readTicket"
	ReferenceAccessReadContext ReferenceAccess = "readContext"
	ReferenceAccessListContext ReferenceAccess = "listContext"
)

// PathKind is the filesystem type of a present architect or repo target.
type PathKind string

const (
	PathKindFile      PathKind = "file"
	PathKindDirectory PathKind = "directory"
)
