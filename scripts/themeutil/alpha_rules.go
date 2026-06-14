package themeutil

type AlphaBaseKind uint8

const (
	AlphaBaseRole AlphaBaseKind = iota
	AlphaBaseSemantic
	AlphaBaseTerminal
)

type AlphaRule struct {
	AlphaKey  string
	BaseKey   string
	BaseKind  AlphaBaseKind
	StyleKeys []string
	Force     bool
}

func alphaRole(alphaKey, role string, keys ...string) AlphaRule {
	return AlphaRule{AlphaKey: alphaKey, BaseKey: role, BaseKind: AlphaBaseRole, StyleKeys: keys}
}

func alphaRoleForce(alphaKey, role string, keys ...string) AlphaRule {
	return AlphaRule{AlphaKey: alphaKey, BaseKey: role, BaseKind: AlphaBaseRole, StyleKeys: keys, Force: true}
}

func alphaSemantic(alphaKey, semantic string, keys ...string) AlphaRule {
	return AlphaRule{AlphaKey: alphaKey, BaseKey: semantic, BaseKind: AlphaBaseSemantic, StyleKeys: keys}
}

func alphaTerminal(alphaKey, terminalKey string) AlphaRule {
	return AlphaRule{
		AlphaKey:  alphaKey,
		BaseKey:   terminalKey,
		BaseKind:  AlphaBaseTerminal,
		StyleKeys: []string{terminalKey},
		Force:     true,
	}
}

var AlphaRules = []AlphaRule{
	alphaRole("ui", "surface", "background", "status_bar.background", "title_bar.background"),
	alphaRole("ui_inactive", "surface", "title_bar.inactive_background"),
	alphaRole("surface", "surface", "surface.background"),
	alphaRole("elevated", "surface", "elevated_surface.background"),
	alphaRole("overlay", "surface", "panel.overlay_background"),
	alphaRole("subheader", "surface", "editor.subheader.background"),
	alphaRole("active_line", "overlay", "editor.active_line.background"),
	alphaRoleForce("active_line_number", "text", "editor.active_line_number"),
	alphaRole("highlighted_line", "overlay", "editor.highlighted_line.background"),
	alphaRole("line_number", "text", "editor.line_number"),
	alphaRole("text_disabled", "text", "text.disabled", "icon.disabled"),
	alphaRole("text_placeholder", "text", "text.placeholder"),
	alphaRole("element_active", "highlight_med", "element.active"),
	alphaRole("element_hover", "highlight_low", "element.hover"),
	alphaRole("element_selected", "highlight_med", "element.selected"),
	alphaRole("element_disabled", "surface", "element.disabled"),
	alphaRole("ghost_active", "highlight_high", "ghost_element.active"),
	alphaRole("ghost_hover", "highlight_low", "ghost_element.hover"),
	alphaRole("ghost_selected", "highlight_high", "ghost_element.selected"),
	alphaRole("ghost_disabled", "surface", "ghost_element.disabled"),
	alphaRole("border_variant", "foam", "border.variant"),
	alphaRole("border_focused", "foam", "border.focused"),
	alphaRole("border_selected", "iris", "border.selected"),
	alphaRole("border_disabled", "muted", "border.disabled"),
	alphaSemantic("conflict_marker", "warning", "version_control.conflict_marker.ours"),
	alphaSemantic("conflict_marker", "info", "version_control.conflict_marker.theirs"),
	alphaRole("panel_focus_border", "muted", "panel.focused_border"),
	alphaRole("panel_indent_guide", "muted", "panel.indent_guide"),
	alphaRole("panel_indent_guide_active", "subtle", "panel.indent_guide_active"),
	alphaRole("pane_focus_border", "muted", "pane.focused_border"),
	alphaRole("pane_group_border", "muted", "pane_group.border"),
	alphaRole("scrollbar_thumb", "muted", "scrollbar.thumb.background"),
	alphaRole("scrollbar_thumb_hover", "muted", "scrollbar.thumb.hover_background"),
	alphaRole("scrollbar_track", "surface", "scrollbar.track.background"),
	alphaRole("scrollbar_track_border", "text", "scrollbar.track.border"),
	alphaRole("search_match", "foam", "search.match_background"),
	alphaRole("search_active", "rose", "search.active_match_background"),
	alphaRole("debugger_line", "rose", "editor.debugger_active_line.background"),
	alphaRole("invisible", "muted", "editor.invisible"),
	alphaRole("indent_guide", "muted", "editor.indent_guide"),
	alphaRole("indent_guide_active", "subtle", "editor.indent_guide_active"),
	alphaRole("wrap_guide", "muted", "editor.wrap_guide"),
	alphaRole("active_wrap_guide", "muted", "editor.active_wrap_guide"),
	alphaRole("doc_highlight_read", "foam", "editor.document_highlight.read_background"),
	alphaRole("doc_highlight_write", "muted", "editor.document_highlight.write_background"),
	alphaRole("doc_highlight_bracket", "iris", "editor.document_highlight.bracket_background"),
	alphaRole("drop_target", "text", "drop_target.background"),
	alphaRole("minimap_bg", "foam", "minimap.thumb.background"),
	alphaRole("minimap_hover", "foam", "minimap.thumb.hover_background"),
	alphaRole("minimap_active", "foam", "minimap.thumb.active_background"),
	alphaSemantic("word_added", "created", "version_control.word_added"),
	alphaSemantic("word_deleted", "deleted", "version_control.word_deleted"),
	alphaTerminal("terminal_background", "terminal.background"),
	alphaTerminal("terminal_ansi_background", "terminal.ansi.background"),
}
