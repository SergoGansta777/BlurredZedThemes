package main

import "zed-themes/scripts/themeutil"

type roleMapping struct {
	key  string
	role string
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
	return themeutil.StripAlpha(roleValue(p, name))
}

func semanticOf(p Palette, name string) string {
	return themeutil.SemanticColor(p.Roles, p.Semantic, name)
}

func terminalBaseOf(p Palette, key string) string {
	if p.Terminal == nil {
		return ""
	}
	if v, ok := p.Terminal[key]; ok {
		return themeutil.StripAlpha(v)
	}
	return ""
}
