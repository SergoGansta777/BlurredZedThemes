package main

type roleMapping struct {
	key  string
	role string
}

type alphaBaseKind uint8

const (
	alphaBaseRole alphaBaseKind = iota
	alphaBaseSemantic
	alphaBaseTerminal
)

type alphaRule struct {
	alphaKey  string
	baseKey   string
	baseKind  alphaBaseKind
	styleKeys []string
	force     bool
}

var defaultRoleMappings = []roleMapping{
	{"background", "surface"},
	{"surface.background", "surface"},
	{"elevated_surface.background", "surface"},
	{"panel.overlay_background", "surface"},
	{"editor.background", "surface"},
	{"editor.gutter.background", "surface"},
	{"editor.subheader.background", "surface"},
	{"editor.active_line.background", "overlay"},
	{"editor.highlighted_line.background", "overlay"},
	{"editor.foreground", "text"},
	{"editor.line_number", "muted"},
	{"editor.active_line_number", "foam"},
	{"editor.invisible", "muted"},
	{"editor.indent_guide", "muted"},
	{"editor.indent_guide_active", "subtle"},
	{"editor.wrap_guide", "muted"},
	{"editor.active_wrap_guide", "muted"},
	{"editor.document_highlight.read_background", "foam"},
	{"editor.document_highlight.write_background", "muted"},
	{"editor.document_highlight.bracket_background", "iris"},
	{"editor.debugger_active_line.background", "rose"},
	{"drop_target.background", "text"},
	{"text", "text"},
	{"text.muted", "muted"},
	{"text.placeholder", "muted"},
	{"text.disabled", "subtle"},
	{"text.accent", "foam"},
	{"link_text.hover", "foam"},
	{"icon", "text"},
	{"icon.muted", "muted"},
	{"icon.placeholder", "muted"},
	{"icon.disabled", "subtle"},
	{"icon.accent", "foam"},
	{"border.variant", "foam"},
	{"border.focused", "foam"},
	{"border.selected", "iris"},
	{"border.disabled", "muted"},
	{"tab.active_foreground", "text"},
	{"tab.inactive_foreground", "muted"},
	{"status_bar.background", "surface"},
	{"title_bar.background", "surface"},
	{"title_bar.inactive_background", "surface"},
	{"status_bar.foreground", "text"},
	{"title_bar.foreground", "text"},
	{"element.active", "highlight_med"},
	{"element.selected", "highlight_med"},
	{"element.hover", "highlight_low"},
	{"element.disabled", "surface"},
	{"element.background", "surface"},
	{"ghost_element.active", "highlight_high"},
	{"ghost_element.selected", "highlight_high"},
	{"ghost_element.hover", "highlight_low"},
	{"ghost_element.disabled", "surface"},
	{"ghost_element.background", "surface"},
	{"minimap.thumb.background", "foam"},
	{"minimap.thumb.hover_background", "foam"},
	{"minimap.thumb.active_background", "foam"},
	{"pane.focused_border", "muted"},
	{"pane_group.border", "muted"},
	{"panel.focused_border", "muted"},
	{"panel.indent_guide", "muted"},
	{"panel.indent_guide_active", "subtle"},
	{"panel.indent_guide_hover", "foam"},
	{"scrollbar.thumb.background", "muted"},
	{"scrollbar.thumb.hover_background", "muted"},
	{"scrollbar.track.background", "surface"},
	{"scrollbar.track.border", "text"},
	{"search.match_background", "foam"},
	{"search.active_match_background", "rose"},
}

var defaultConstMappings = map[string]string{
	"border":                  transparentColor,
	"border.transparent":      transparentColor,
	"tab.inactive_background": transparentColor,
	"tab_bar.background":      transparentColor,
}

var defaultTransparentKeys = []string{
	"panel.background",
	"toolbar.background",
	"tab_bar.background",
	"tab.inactive_background",
	"border",
	"border.transparent",
}

func roleValue(p Palette, name string) string {
	return p.Roles[name]
}

func roleOpaque(p Palette, name string) string {
	return stripAlpha(roleValue(p, name))
}

func semanticOf(p Palette, name string) string {
	return semanticColor(p, name)
}

func terminalBaseOf(p Palette, key string) string {
	if p.Terminal == nil {
		return ""
	}
	if v, ok := p.Terminal[key]; ok {
		return stripAlpha(v)
	}
	return ""
}

func alphaRole(alphaKey, role string, keys ...string) alphaRule {
	return alphaRule{alphaKey: alphaKey, baseKey: role, baseKind: alphaBaseRole, styleKeys: keys}
}

func alphaSemantic(alphaKey, semantic string, keys ...string) alphaRule {
	return alphaRule{alphaKey: alphaKey, baseKey: semantic, baseKind: alphaBaseSemantic, styleKeys: keys}
}

func alphaTerminal(alphaKey, terminalKey string) alphaRule {
	return alphaRule{
		alphaKey:  alphaKey,
		baseKey:   terminalKey,
		baseKind:  alphaBaseTerminal,
		styleKeys: []string{terminalKey},
		force:     true,
	}
}

var alphaRules = []alphaRule{
	alphaRole("ui", "surface", "background", "status_bar.background", "title_bar.background"),
	alphaRole("ui_inactive", "surface", "title_bar.inactive_background"),
	alphaRole("surface", "surface", "surface.background"),
	alphaRole("elevated", "surface", "elevated_surface.background"),
	alphaRole("overlay", "surface", "panel.overlay_background"),
	alphaRole("subheader", "surface", "editor.subheader.background"),
	alphaRole("active_line", "overlay", "editor.active_line.background"),
	alphaRole("highlighted_line", "overlay", "editor.highlighted_line.background"),
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
	alphaRole("conflict_marker", "foam", "version_control.conflict_marker.theirs"),
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
