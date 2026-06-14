package themeutil

import "testing"

func TestNormalizeImportedStyleRenamesKnownSyntaxTypos(t *testing.T) {
	style := map[string]any{
		"syntax": map[string]any{
			"ion.special": map[string]any{
				"color": "#286983",
			},
		},
	}

	normalized := NormalizeImportedStyle(style)
	syntax, ok := normalized["syntax"].(map[string]any)
	if !ok {
		t.Fatal("syntax map missing")
	}
	if _, ok := syntax["ion.special"]; ok {
		t.Fatal("ion.special should be removed")
	}

	entry, ok := syntax["function.special"].(map[string]any)
	if !ok {
		t.Fatal("function.special should be created")
	}
	if got, _ := entry["color"].(string); got != "#286983" {
		t.Fatalf("function.special color = %q, want %q", got, "#286983")
	}
}

func TestWithAlphaReturnsShortOpaqueColors(t *testing.T) {
	if got := WithAlpha("#ffffff", "FF"); got != "#FFFFFF" {
		t.Fatalf("WithAlpha opaque = %q, want %q", got, "#FFFFFF")
	}
}

func TestStripAlpha(t *testing.T) {
	if got := StripAlpha("#AABBCC80"); got != "#AABBCC" {
		t.Fatalf("StripAlpha = %q, want %q", got, "#AABBCC")
	}
	if got := StripAlpha("#AABBCC"); got != "#AABBCC" {
		t.Fatalf("StripAlpha opaque = %q, want %q", got, "#AABBCC")
	}
}
