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

func TestBuildStyleDerivesDiffHunks(t *testing.T) {
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
			"semantic_bg":  "1F",
		},
	}, false)

	if got, _ := style["editor.diff_hunk.added.background"].(string); got != "#22AA6659" {
		t.Fatalf("editor.diff_hunk.added.background = %q, want %q", got, "#22AA6659")
	}
	if got, _ := style["editor.diff_hunk.added.hollow_background"].(string); got != "#22AA661F" {
		t.Fatalf("editor.diff_hunk.added.hollow_background = %q, want %q", got, "#22AA661F")
	}
	if got, _ := style["editor.diff_hunk.deleted.hollow_border"].(string); got != "#CC44551F" {
		t.Fatalf("editor.diff_hunk.deleted.hollow_border = %q, want %q", got, "#CC44551F")
	}
}

func TestPruneRedundantOverridesKeepsRequiredOverrides(t *testing.T) {
	template := map[string]any{}
	palette := Palette{
		Meta: Meta{
			Appearance: "dark",
		},
		Roles: map[string]string{
			"surface":        "#101010",
			"text":           "#FFFFFF",
			"muted":          "#888888",
			"subtle":         "#666666",
			"pine":           "#22AA66",
			"love":           "#CC4455",
			"foam":           "#55AAFF",
			"gold":           "#DDAA55",
			"rose":           "#CC8855",
			"iris":           "#AA88FF",
			"highlight_med":  "#222222",
			"highlight_low":  "#181818",
			"highlight_high": "#333333",
		},
		Overrides: map[string]any{
			"version_control.added": "#55AAFF",
			"editor.line_number":    "#777777",
		},
	}

	reference := buildStyle(template, palette, AlphaConfig{}, false)
	removeTODOs(reference)
	pruneRedundantOverrides(&palette, template, AlphaConfig{}, reference, false)

	if _, ok := palette.Overrides["version_control.added"]; ok {
		t.Fatal("version_control.added should be pruned")
	}
	if got, ok := palette.Overrides["editor.line_number"]; !ok || got != "#777777" {
		t.Fatalf("editor.line_number override = %v, present %v; want preserved #777777", got, ok)
	}
}

func TestValuesEqualTreatsHexColorCaseAsEqual(t *testing.T) {
	if !valuesEqual("#1f2024", "#1F2024") {
		t.Fatal("hex colors should compare case-insensitively")
	}
	if valuesEqual("#1f2024", "#1F202440") {
		t.Fatal("hex colors with different alpha should not compare equal")
	}
}

func TestDerivedPaletteOverridesAlphaBasesAndVim(t *testing.T) {
	style := buildStyle(map[string]any{}, Palette{
		Meta: Meta{
			Appearance: "dark",
		},
		Roles: map[string]string{
			"surface":        "#101010",
			"text":           "#FFFFFF",
			"muted":          "#888888",
			"subtle":         "#666666",
			"pine":           "#22AA66",
			"love":           "#CC4455",
			"foam":           "#55AAFF",
			"gold":           "#DDAA55",
			"rose":           "#CC8855",
			"iris":           "#AA88FF",
			"highlight_med":  "#222222",
			"highlight_low":  "#181818",
			"highlight_high": "#333333",
		},
		Derived: map[string]string{
			"doc_highlight_read":    "#445566",
			"invisible":             "#FFFFFF",
			"scrollbar_thumb":       "#343434",
			"scrollbar_thumb_hover": "#343434",
			"search_match":          "#FFFFFF",
			"vim_replace":           "#D86A5C",
			"vim_visual":            "#9A79C4",
		},
	}, AlphaConfig{
		Dark: map[string]string{
			"doc_highlight_read":    "33",
			"invisible":             "06",
			"scrollbar_thumb":       "50",
			"scrollbar_thumb_hover": "FF",
			"search_match":          "26",
		},
	}, false)

	if got, _ := style["editor.document_highlight.read_background"].(string); got != "#44556633" {
		t.Fatalf("editor.document_highlight.read_background = %q, want %q", got, "#44556633")
	}
	if got, _ := style["editor.invisible"].(string); got != "#FFFFFF06" {
		t.Fatalf("editor.invisible = %q, want %q", got, "#FFFFFF06")
	}
	if got, _ := style["scrollbar.thumb.background"].(string); got != "#34343450" {
		t.Fatalf("scrollbar.thumb.background = %q, want %q", got, "#34343450")
	}
	if got, _ := style["scrollbar.thumb.border"].(string); got != "#34343400" {
		t.Fatalf("scrollbar.thumb.border = %q, want %q", got, "#34343400")
	}
	if got, _ := style["scrollbar.thumb.hover_background"].(string); got != "#343434" {
		t.Fatalf("scrollbar.thumb.hover_background = %q, want %q", got, "#343434")
	}
	if got, _ := style["search.match_background"].(string); got != "#FFFFFF26" {
		t.Fatalf("search.match_background = %q, want %q", got, "#FFFFFF26")
	}
	if got, _ := style["vim.visual.background"].(string); got != "#9A79C4" {
		t.Fatalf("vim.visual.background = %q, want %q", got, "#9A79C4")
	}
	if got, _ := style["vim.replace.background"].(string); got != "#D86A5C" {
		t.Fatalf("vim.replace.background = %q, want %q", got, "#D86A5C")
	}
}

func TestSemanticVersionControlAndVimDerivations(t *testing.T) {
	style := buildStyle(map[string]any{}, Palette{
		Meta: Meta{
			Appearance: "dark",
		},
		Roles: map[string]string{
			"surface":        "#101010",
			"text":           "#FFFFFF",
			"muted":          "#888888",
			"subtle":         "#666666",
			"pine":           "#22AA66",
			"love":           "#CC4455",
			"foam":           "#55AAFF",
			"gold":           "#DDAA55",
			"rose":           "#CC8855",
			"iris":           "#AA88FF",
			"highlight_med":  "#222222",
			"highlight_low":  "#181818",
			"highlight_high": "#333333",
		},
		Semantic: map[string]string{
			"conflict": "#B88432",
			"deleted":  "#D45A4A",
			"info":     "#4A9BD6",
			"modified": "#C99B45",
			"renamed":  "#8D6CB5",
		},
	}, AlphaConfig{}, false)

	if got, _ := style["version_control.conflict"].(string); got != "#B88432" {
		t.Fatalf("version_control.conflict = %q, want %q", got, "#B88432")
	}
	if got, _ := style["version_control.conflict_marker.theirs"].(string); got != "#4A9BD6" {
		t.Fatalf("version_control.conflict_marker.theirs = %q, want %q", got, "#4A9BD6")
	}
	if got, _ := style["vim.normal.background"].(string); got != "#4A9BD6" {
		t.Fatalf("vim.normal.background = %q, want %q", got, "#4A9BD6")
	}
	if got, _ := style["vim.insert.background"].(string); got != "#C99B45" {
		t.Fatalf("vim.insert.background = %q, want %q", got, "#C99B45")
	}
	if got, _ := style["vim.visual.background"].(string); got != "#8D6CB5" {
		t.Fatalf("vim.visual.background = %q, want %q", got, "#8D6CB5")
	}
	if got, _ := style["vim.replace.background"].(string); got != "#D45A4A" {
		t.Fatalf("vim.replace.background = %q, want %q", got, "#D45A4A")
	}
}

func TestSemanticVersionControlOverrides(t *testing.T) {
	style := buildStyle(map[string]any{}, Palette{
		Meta: Meta{
			Appearance: "dark",
		},
		Roles: map[string]string{
			"surface":        "#101010",
			"text":           "#FFFFFF",
			"muted":          "#888888",
			"subtle":         "#666666",
			"pine":           "#22AA66",
			"love":           "#CC4455",
			"foam":           "#55AAFF",
			"gold":           "#DDAA55",
			"rose":           "#CC8855",
			"iris":           "#AA88FF",
			"highlight_med":  "#222222",
			"highlight_low":  "#181818",
			"highlight_high": "#333333",
		},
		Semantic: map[string]string{
			"created":        "#22AA66",
			"deleted":        "#CC4455",
			"modified":       "#DDAA55",
			"renamed":        "#AA88FF",
			"vcs_added":      "#549159",
			"vcs_deleted":    "#A75749",
			"vcs_modified":   "#375FAD",
			"vcs_renamed":    "#6183BB",
			"vcs_word_added": "#294436",
		},
	}, AlphaConfig{
		Dark: map[string]string{
			"word_added": "59",
		},
	}, false)

	if got, _ := style["version_control.added"].(string); got != "#549159" {
		t.Fatalf("version_control.added = %q, want %q", got, "#549159")
	}
	if got, _ := style["version_control.deleted"].(string); got != "#A75749" {
		t.Fatalf("version_control.deleted = %q, want %q", got, "#A75749")
	}
	if got, _ := style["version_control.modified"].(string); got != "#375FAD" {
		t.Fatalf("version_control.modified = %q, want %q", got, "#375FAD")
	}
	if got, _ := style["version_control.renamed"].(string); got != "#6183BB" {
		t.Fatalf("version_control.renamed = %q, want %q", got, "#6183BB")
	}
	if got, _ := style["version_control.word_added"].(string); got != "#294436" {
		t.Fatalf("version_control.word_added = %q, want %q", got, "#294436")
	}
}

func TestBuildStyleAppliesTextAlphaRules(t *testing.T) {
	style := buildStyle(map[string]any{}, Palette{
		Meta: Meta{
			Appearance: "dark",
		},
		Roles: map[string]string{
			"surface":        "#101010",
			"text":           "#FFFFFF",
			"muted":          "#FFFFFF99",
			"subtle":         "#FFFFFF66",
			"pine":           "#22AA66",
			"love":           "#CC4455",
			"foam":           "#55AAFF",
			"gold":           "#DDAA55",
			"rose":           "#CC8855",
			"iris":           "#AA88FF",
			"highlight_med":  "#222222",
			"highlight_low":  "#181818",
			"highlight_high": "#333333",
		},
	}, AlphaConfig{
		Dark: map[string]string{
			"active_line_number": "FF",
			"line_number":        "6B",
			"text_disabled":      "4D",
			"text_placeholder":   "66",
		},
	}, false)

	if got, _ := style["editor.active_line_number"].(string); got != "#FFFFFF" {
		t.Fatalf("editor.active_line_number = %q, want %q", got, "#FFFFFF")
	}
	if got, _ := style["editor.line_number"].(string); got != "#FFFFFF6B" {
		t.Fatalf("editor.line_number = %q, want %q", got, "#FFFFFF6B")
	}
	if got, _ := style["text.disabled"].(string); got != "#FFFFFF4D" {
		t.Fatalf("text.disabled = %q, want %q", got, "#FFFFFF4D")
	}
	if got, _ := style["icon.disabled"].(string); got != "#FFFFFF4D" {
		t.Fatalf("icon.disabled = %q, want %q", got, "#FFFFFF4D")
	}
	if got, _ := style["text.placeholder"].(string); got != "#FFFFFF66" {
		t.Fatalf("text.placeholder = %q, want %q", got, "#FFFFFF66")
	}
}
