package themeutil

import "testing"

func TestSemanticColorsAppliesDefaultsAndOverrides(t *testing.T) {
	colors := SemanticColors(map[string]string{
		"foam":   "#55AAFF",
		"gold":   "#DDAA55",
		"love":   "#CC4455",
		"muted":  "#888888",
		"pine":   "#22AA66",
		"rose":   "#CC8855",
		"subtle": "#666666",
	}, map[string]string{
		"info":      "#4A9BD6",
		"vcs_added": "#549159",
	})

	if got := colors["error"]; got != "#CC4455" {
		t.Fatalf("error = %q, want #CC4455", got)
	}
	if got := colors["info"]; got != "#4A9BD6" {
		t.Fatalf("info = %q, want #4A9BD6", got)
	}
	if got := colors["vcs_added"]; got != "#549159" {
		t.Fatalf("vcs_added = %q, want #549159", got)
	}
}
