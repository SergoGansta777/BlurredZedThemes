package main

import "testing"

func TestApplyDerivedEditorLineNumbers(t *testing.T) {
	style := map[string]any{
		"editor.line_number":        "#4E5A5F",
		"editor.active_line_number": "#D0D4DA",
	}

	applyDerivedEditorLineNumbers(style)

	got, _ := style["editor.hover_line_number"].(string)
	if got != "#AFB5BB" {
		t.Fatalf("hover line number = %q, want %q", got, "#AFB5BB")
	}
}

func TestApplyDerivedSyntaxBackfillsExistingMaps(t *testing.T) {
	style := map[string]any{
		"syntax": map[string]any{
			"function":    map[string]any{"color": "#112233"},
			"keyword":     map[string]any{"color": "#445566"},
			"tag":         map[string]any{"color": "#778899"},
			"type":        map[string]any{"color": "#AABBCC"},
			"text":        map[string]any{"color": "#FFFFFF"},
			"comment":     map[string]any{"color": "#999999"},
			"string":      map[string]any{"color": "#00FF00"},
			"number":      map[string]any{"color": "#FF0000"},
			"property":    map[string]any{"color": "#ABCDEF"},
			"constant":    map[string]any{"color": "#FEDCBA"},
			"link_uri":    map[string]any{"color": "#135790"},
			"attribute":   map[string]any{"color": "#2468AC"},
			"punctuation": map[string]any{"color": "#555555"},
			"parameter":   map[string]any{"color": "#CCCCCC"},
			"variable":    map[string]any{"color": "#DDDDDD"},
			"preproc":     map[string]any{"color": "#EEEEEE"},
			"diff.plus":   map[string]any{"color": "#00AA00"},
			"diff.minus":  map[string]any{"color": "#AA0000"},
			"emphasis":    map[string]any{"color": "#FF00FF"},
			"title":       map[string]any{"color": "#00FFFF"},
		},
	}

	applyDerivedSyntax(style, Palette{
		Roles: map[string]string{
			"text":   "#FFFFFF",
			"muted":  "#999999",
			"subtle": "#555555",
			"pine":   "#00AA00",
			"gold":   "#FFCC00",
			"foam":   "#66CCFF",
			"love":   "#FF6677",
			"rose":   "#CC9966",
			"iris":   "#AA88FF",
		},
	})

	syntax, ok := style["syntax"].(map[string]any)
	if !ok {
		t.Fatal("syntax map missing")
	}

	if _, ok := syntax["function.call"]; !ok {
		t.Fatal("function.call was not backfilled")
	}
	if _, ok := syntax["selector"]; !ok {
		t.Fatal("selector was not backfilled")
	}
	if _, ok := syntax["punctuation.markup"]; !ok {
		t.Fatal("punctuation.markup was not backfilled")
	}

	functionCall, _ := syntax["function.call"].(map[string]any)
	if got, _ := functionCall["color"].(string); got != "#112233" {
		t.Fatalf("function.call color = %q, want %q", got, "#112233")
	}
}

func TestBuildStyleDerivesWordDiffHighlights(t *testing.T) {
	style := buildStyle(map[string]any{}, Palette{
		Meta: Meta{
			Appearance: "dark",
		},
		Roles: map[string]string{
			"surface": "#101010",
			"text":    "#FFFFFF",
			"muted":   "#888888",
			"subtle":  "#666666",
			"pine":    "#22AA66",
			"love":    "#CC4455",
			"foam":    "#55AAFF",
			"gold":    "#DDAA55",
			"rose":    "#CC8855",
			"iris":    "#AA88FF",
		},
		Semantic: map[string]string{
			"created": "#22AA66",
			"deleted": "#CC4455",
		},
	}, AlphaConfig{
		Dark: map[string]string{
			"word_added":   "59",
			"word_deleted": "59",
		},
	}, false)

	if got, _ := style["version_control.word_added"].(string); got != "#22AA6659" {
		t.Fatalf("version_control.word_added = %q, want %q", got, "#22AA6659")
	}
	if got, _ := style["version_control.word_deleted"].(string); got != "#CC445559" {
		t.Fatalf("version_control.word_deleted = %q, want %q", got, "#CC445559")
	}
}
