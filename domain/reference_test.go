package domain

import (
	"errors"
	"strings"
	"testing"
)

func TestParseReferenceAccepts(t *testing.T) {
	cases := []struct {
		in   string
		want ParsedReference
	}{
		{"2026-10-27-0534-ship-it", ParsedReference{Kind: ReferenceKindTicket, TicketID: "2026-10-27-0534-ship-it"}},
		{"architect:OVERVIEW.md", ParsedReference{Kind: ReferenceKindArchitect, Owner: "architect", Target: "OVERVIEW.md"}},
		{"architect:investigations/2026/notes.md", ParsedReference{Kind: ReferenceKindArchitect, Owner: "architect", Target: "investigations/2026/notes.md"}},
		{"architect:.", ParsedReference{Kind: ReferenceKindArchitect, Owner: "architect", Target: "."}},
		{"daemon:docs", ParsedReference{Kind: ReferenceKindRepo, Owner: "daemon", Target: "docs"}},
		{"ims-sip.e2e_x:a:b.md", ParsedReference{Kind: ReferenceKindRepo, Owner: "ims-sip.e2e_x", Target: "a:b.md"}},
		{"daemon:.github/workflows/ci.yml", ParsedReference{Kind: ReferenceKindRepo, Owner: "daemon", Target: ".github/workflows/ci.yml"}},
	}
	for _, c := range cases {
		got, err := ParseReference(c.in)
		if err != nil {
			t.Fatalf("ParseReference(%q): %v", c.in, err)
		}
		if got != c.want {
			t.Fatalf("ParseReference(%q) = %+v, want %+v", c.in, got, c.want)
		}
		if got.String() != c.in {
			t.Fatalf("String() = %q, want %q", got.String(), c.in)
		}
	}
}

func TestParseReferenceRejects(t *testing.T) {
	cases := map[string]string{
		"":                         "empty reference",
		"  ":                       "empty reference",
		" architect:a":             "whitespace",
		"/Users/k/architects/x.md": "absolute",
		"~/notes.md":               "absolute",
		"docs/notes.md":            "unqualified",
		"notes.md":                 "unqualified",
		"2026-10-27-05-bad":        "unqualified",
		"2026-10-27-0534-":         "not a valid ticket id",
		"https://example.com/x":    "URLs",
		":x":                       "missing prefix",
		"architect:":               "empty target",
		"architect:/etc/passwd":    "relative to the root",
		"daemon:~/x":               "no '~'",
		"daemon:a/../../etc":       "'..'",
		"daemon:..":                "'..'",
		"daemon:./a":               "'./'",
		"daemon:a//b":              "empty path segment",
		"daemon:a/":                "trailing",
		`daemon:a\b`:               "backslash",
		"Architect:x":              "spelled architect",
		"bad key:x":                "not architect or a repo key",
		"daemon:a\nb":              "control character",
		"repo/sub:x":               "not architect or a repo key",
	}
	for in, want := range cases {
		_, err := ParseReference(in)
		var syntax *ReferenceSyntaxError
		if !errors.As(err, &syntax) {
			t.Fatalf("ParseReference(%q) err = %v, want syntax error", in, err)
		}
		if !strings.Contains(err.Error(), want) {
			t.Fatalf("ParseReference(%q) = %q, want it to mention %q", in, err, want)
		}
	}
}

func TestValidRepoKey(t *testing.T) {
	for _, ok := range []string{"daemon", "ims-sip-e2e", "a.b_c", "A1"} {
		if !ValidRepoKey(ok) {
			t.Fatalf("ValidRepoKey(%q) = false", ok)
		}
	}
	for _, bad := range []string{"", "-x", "a b", "a:b", "a/b"} {
		if ValidRepoKey(bad) {
			t.Fatalf("ValidRepoKey(%q) = true", bad)
		}
	}
}
