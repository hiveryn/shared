package domain

import "testing"

func TestEffectiveArtifactMode(t *testing.T) {
	for _, tc := range []struct{ mode, contract, want string }{
		{"", "", "none"}, {"", "legacy files", "required"},
		{"none", "", "none"}, {"optional", "requested files", "optional"},
		{"required", "promised files", "required"},
	} {
		if got := EffectiveArtifactMode(tc.mode, tc.contract); got != tc.want {
			t.Fatalf("mode %q contract %q: %q, want %q", tc.mode, tc.contract, got, tc.want)
		}
	}
}
