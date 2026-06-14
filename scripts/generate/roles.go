package main

import (
	"strings"

	"zed-themes/scripts/themeutil"
)

func applyRoleMappings(style map[string]any, p Palette) {
	if len(p.Roles) == 0 {
		return
	}

	s := styleMap(style)
	role := func(name string) string { return roleValue(p, name) }

	for _, mapping := range defaultRoleMappings {
		s.SetRole(mapping.key, role(mapping.role))
	}

	for key, value := range defaultConstMappings {
		s.SetRole(key, value)
	}
	s.SetAny("minimap.thumb.border", nil)

	s.SetAny("scrollbar.thumb.active_background", nil)
	if scrollbarThumb := derivedColor(p, "scrollbar_thumb"); scrollbarThumb != "" {
		s.SetRole("scrollbar.thumb.border", withAlpha(scrollbarThumb, "00"))
	} else {
		s.SetAny("scrollbar.thumb.border", nil)
	}

	semantic := themeutil.SemanticColors(p.Roles, p.Semantic)
	for k, v := range semantic {
		if strings.HasPrefix(k, "vcs_") {
			continue
		}
		if v != "" {
			s.SetRole(k, v)
			s.SetRole(k+".border", v)
		}
	}

	s.SetRole("version_control.added", firstNonEmpty(semantic["vcs_added"], semantic["created"]))
	s.SetRole("version_control.deleted", firstNonEmpty(semantic["vcs_deleted"], semantic["deleted"]))
	s.SetRole("version_control.modified", firstNonEmpty(semantic["vcs_modified"], semantic["modified"]))
	s.SetRole("version_control.renamed", firstNonEmpty(semantic["vcs_renamed"], semantic["renamed"]))
	s.SetRole("version_control.conflict", firstNonEmpty(semantic["vcs_conflict"], semantic["conflict"], semantic["modified"]))
	s.SetRole("version_control.ignored", semantic["ignored"])
	s.SetRole("version_control.conflict_marker.ours", semantic["warning"])
	s.SetRole("version_control.conflict_marker.theirs", semantic["info"])
	s.SetRole("version_control.word_added", semantic["vcs_word_added"])
	s.SetRole("version_control.word_deleted", semantic["vcs_word_deleted"])

	s.SetRole("debugger.accent", semantic["error"])

	if len(p.Accents) == 0 {
		accents := []string{role("foam"), role("iris"), role("pine"), role("rose"), role("gold"), role("love")}
		var out []string
		for _, a := range accents {
			if a != "" {
				out = append(out, a)
			}
		}
		if len(out) > 0 {
			style["accents"] = out
		}
	}
}

func setSemanticBackgrounds(style map[string]any, p Palette, alpha AlphaConfig, editorBg string, opaqueEditorBg string) {
	if editorBg == "" {
		return
	}

	s := styleMap(style)
	type rule struct {
		key        string
		alpha      string
		forceSolid bool
	}
	semanticAlpha, _ := alphaValue(p.Meta.Appearance, alpha, "semantic_bg")
	rules := []rule{
		{"warning", semanticAlpha, true},
		{"info", semanticAlpha, true},
		{"success", semanticAlpha, true},
		{"unreachable", semanticAlpha, true},
		{"conflict", "26", false},
		{"created", "26", false},
		{"deleted", "26", false},
		{"modified", "26", false},
		{"renamed", "26", false},
		{"ignored", "26", false},
	}

	for _, r := range rules {
		bgKey := r.key + ".background"
		if s.HasValue(bgKey) {
			continue
		}
		if r.forceSolid && opaqueEditorBg != "" {
			style[bgKey] = opaqueEditorBg
			continue
		}
		if fg, ok := style[r.key].(string); ok && fg != "" {
			style[bgKey] = withAlpha(fg, r.alpha)
			continue
		}
		style[bgKey] = editorBg
	}

	fallbackBg := editorBg
	if opaqueEditorBg != "" {
		fallbackBg = opaqueEditorBg
	}
	editorFallback := []string{
		"error",
		"hidden",
		"hint",
		"predictive",
	}
	for _, k := range editorFallback {
		bgKey := k + ".background"
		if s.HasValue(bgKey) {
			continue
		}
		style[bgKey] = fallbackBg
	}
}

func semanticColor(p Palette, name string) string {
	return themeutil.SemanticColor(p.Roles, p.Semantic, name)
}

func derivedColor(p Palette, name string) string {
	if p.Derived == nil {
		return ""
	}
	return p.Derived[name]
}

func isStandardizedKey(key string) bool {
	return themeutil.IsStandardizedKey(key)
}
