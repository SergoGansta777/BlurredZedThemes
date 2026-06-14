package main

func applyTerminalDims(style map[string]any) {
	dims := map[string]string{
		"terminal.ansi.dim_black":   "terminal.ansi.black",
		"terminal.ansi.dim_red":     "terminal.ansi.red",
		"terminal.ansi.dim_green":   "terminal.ansi.green",
		"terminal.ansi.dim_yellow":  "terminal.ansi.yellow",
		"terminal.ansi.dim_blue":    "terminal.ansi.blue",
		"terminal.ansi.dim_magenta": "terminal.ansi.magenta",
		"terminal.ansi.dim_cyan":    "terminal.ansi.cyan",
		"terminal.ansi.dim_white":   "terminal.ansi.white",
	}
	for dimKey, baseKey := range dims {
		if v, ok := style[dimKey]; ok {
			if s, ok := v.(string); ok && !isTodoValue(s) {
				continue
			}
		}
		if v, ok := style[baseKey].(string); ok {
			style[dimKey] = v
		}
	}
}

func applyDerivedVim(style map[string]any, p Palette) {
	s := styleMap(style)
	role := func(name string) string { return roleValue(p, name) }
	normal := firstNonEmpty(derivedColor(p, "vim_normal"), semanticOf(p, "info"), role("foam"), role("pine"))
	insert := firstNonEmpty(derivedColor(p, "vim_insert"), semanticOf(p, "modified"), role("rose"), role("gold"))
	visual := firstNonEmpty(derivedColor(p, "vim_visual"), semanticOf(p, "renamed"), role("iris"), role("rose"))
	replace := firstNonEmpty(derivedColor(p, "vim_replace"), semanticOf(p, "deleted"), role("love"), role("rose"))
	foreground := firstNonEmpty(role("base"), role("surface"), role("text"))

	s.SetRole("vim.mode.text", foreground)
	s.SetRole("vim.normal.background", normal)
	s.SetRole("vim.normal.foreground", foreground)
	s.SetRole("vim.helix_normal.background", normal)
	s.SetRole("vim.helix_normal.foreground", foreground)
	s.SetRole("vim.insert.background", insert)
	s.SetRole("vim.insert.foreground", foreground)
	s.SetRole("vim.visual.background", visual)
	s.SetRole("vim.visual.foreground", foreground)
	s.SetRole("vim.helix_select.background", visual)
	s.SetRole("vim.helix_select.foreground", foreground)
	s.SetRole("vim.visual_line.background", visual)
	s.SetRole("vim.visual_line.foreground", foreground)
	s.SetRole("vim.visual_block.background", visual)
	s.SetRole("vim.visual_block.foreground", foreground)
	s.SetRole("vim.replace.background", replace)
	s.SetRole("vim.replace.foreground", foreground)
}

func applyDerivedPlayers(style map[string]any, p Palette, alpha AlphaConfig) {
	if styleMap(style).HasValue("players") {
		return
	}
	if len(p.Accents) == 0 {
		return
	}
	alphaHex, ok := alphaValue(p.Meta.Appearance, alpha, "selection")
	if !ok {
		alphaHex = selectionAlpha
	}
	var players []map[string]string
	for _, c := range p.Accents {
		if c == "" {
			continue
		}
		players = append(players, map[string]string{
			"cursor":     c,
			"background": c,
			"selection":  withAlpha(c, alphaHex),
		})
	}
	if len(players) > 0 {
		style["players"] = players
	}
}

func applyDerivedEditorLineNumbers(style map[string]any) {
	if styleMap(style).HasValue("editor.hover_line_number") {
		return
	}

	active, _ := style["editor.active_line_number"].(string)
	normal, _ := style["editor.line_number"].(string)

	switch {
	case active != "" && normal != "":
		style["editor.hover_line_number"] = mixOpaqueColors(active, normal, hoverBlendAlpha)
	case active != "":
		style["editor.hover_line_number"] = active
	case normal != "":
		style["editor.hover_line_number"] = normal
	}
}

func applyDerivedDiffHunks(style map[string]any, p Palette, alpha AlphaConfig) {
	s := styleMap(style)
	setDiffHunk := func(kind string, semantic string, strongAlphaKey string) {
		base := semanticOf(p, semantic)
		if base == "" {
			return
		}

		strong := alphaHexOrDefault(p.Meta.Appearance, alpha, strongAlphaKey, "59")
		hollow := withAlpha(base, alphaHexOrDefault(p.Meta.Appearance, alpha, "semantic_bg", "26"))

		s.SetRole("editor.diff_hunk."+kind+".background", withAlpha(base, strong))
		s.SetRole("editor.diff_hunk."+kind+".hollow_background", hollow)
		s.SetRole("editor.diff_hunk."+kind+".hollow_border", hollow)
	}

	setDiffHunk("added", "created", "word_added")
	setDiffHunk("deleted", "deleted", "word_deleted")
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}
