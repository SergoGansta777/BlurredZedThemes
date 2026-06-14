package main

import (
	"maps"
	"strings"

	"zed-themes/scripts/themeutil"
)

var flatSurfaceBackgroundKeys = []string{
	"background",
	"editor.background",
	"editor.gutter.background",
	"panel.background",
	"status_bar.background",
	"tab.inactive_background",
	"tab_bar.background",
	"title_bar.background",
	"title_bar.inactive_background",
	"toolbar.background",
}

func buildStyle(template map[string]any, p Palette, alpha AlphaConfig, prune bool) map[string]any {
	style := map[string]any{}
	maps.Copy(style, template)

	mergeAny(style, p.Style)
	applyRoleMappings(style, p)

	if p.Meta.BackgroundAppearance != "" {
		style["background.appearance"] = p.Meta.BackgroundAppearance
	}

	if len(p.Accents) > 0 {
		style["accents"] = p.Accents
	}

	mergeStringMap(style, p.Colors)

	mergeStringMap(style, p.Terminal)
	applyTerminalDims(style)
	applyAlphaRules(style, p, alpha)

	applyDerivedVim(style, p)
	applyDerivedPlayers(style, p, alpha)
	applyDerivedSyntax(style, p)
	mergeAny(style, p.Overrides)
	applyDerivedEditorLineNumbers(style)
	baseEditorBg, _ := style["editor.background"].(string)
	applyBlurMode(style, p.Meta)

	editorBg, _ := style["editor.background"].(string)
	islandBg := editorBg
	if strings.EqualFold(p.Meta.BlurMode, blurModeFlat) && baseEditorBg != "" {
		islandBg = baseEditorBg
	}
	opaqueSemanticBg := ""
	if hasOpaqueSemanticBackgrounds(p.Meta) && baseEditorBg != "" {
		opaqueSemanticBg = themeutil.StripAlpha(baseEditorBg)
	}
	if editorBg != "" {
		if _, ok := style["editor.gutter.background"]; !ok {
			style["editor.gutter.background"] = editorBg
		}
		applyEditorIslandElements(style, p, alpha, islandBg)
		setSemanticBackgrounds(style, p, alpha, editorBg, opaqueSemanticBg)
		applyDerivedDiffHunks(style, p, alpha)
	}
	applyFlatMode(style, p.Meta, baseEditorBg)

	if text, ok := style["text"].(string); ok && text != "" {
		s := styleMap(style)
		if s.IsUnset("status_bar.foreground") {
			style["status_bar.foreground"] = text
		}
		if s.IsUnset("title_bar.foreground") {
			style["title_bar.foreground"] = text
		}
	}

	styleMap(style).SetDefaults(defaultTransparentKeys, transparentColor)

	if prune && shouldPruneStyle(p.Style) {
		for k := range style {
			if _, ok := p.Style[k]; !ok {
				delete(style, k)
			}
		}
	}

	return style
}

func applyEditorIslandElements(style map[string]any, p Palette, alpha AlphaConfig, editorBg string) {
	if editorBg == "" {
		return
	}

	setIsland := func(styleKey, alphaKey string) {
		if hasPaletteOverride(p, styleKey) {
			return
		}
		if alphaKey == "" {
			style[styleKey] = editorBg
			return
		}
		if alphaHex, ok := alphaValue(p.Meta.Appearance, alpha, alphaKey); ok {
			style[styleKey] = withAlpha(editorBg, alphaHex)
			return
		}
		style[styleKey] = editorBg
	}

	setIsland("tab.active_background", "tab_active")
}

func hasPaletteOverride(p Palette, key string) bool {
	if p.Overrides == nil {
		return false
	}
	_, ok := p.Overrides[key]
	return ok
}

func mergeStringMap(dst map[string]any, src map[string]string) {
	for k, v := range src {
		dst[k] = v
	}
}

func mergeAny(dst map[string]any, src map[string]any) {
	maps.Copy(dst, src)
}

func applyBlurMode(style map[string]any, meta Meta) {
	if !strings.EqualFold(meta.BlurMode, blurModeFlat) {
		return
	}
	style["editor.background"] = transparentColor
	style["editor.gutter.background"] = transparentColor
}

func applyFlatMode(style map[string]any, meta Meta, baseEditorBg string) {
	if !strings.EqualFold(meta.BackgroundAppearance, backgroundAppearanceOpaque) {
		return
	}
	surfaceBg := themeutil.StripAlpha(baseEditorBg)
	if surfaceBg == "" || strings.EqualFold(surfaceBg, transparentColor) {
		return
	}
	for _, key := range flatSurfaceBackgroundKeys {
		style[key] = surfaceBg
	}
	style["tab.active_background"] = flatTabActiveBackground(style, surfaceBg)
}

func flatTabActiveBackground(style map[string]any, surfaceBg string) string {
	if tabBg := solidColorValue(style, "tab.active_background", surfaceBg); tabBg != "" && !strings.EqualFold(tabBg, surfaceBg) {
		return tabBg
	}
	if selectedBg := solidColorValue(style, "element.selected", surfaceBg); selectedBg != "" && !strings.EqualFold(selectedBg, surfaceBg) {
		return selectedBg
	}
	if hoverBg := solidColorValue(style, "element.hover", surfaceBg); hoverBg != "" && !strings.EqualFold(hoverBg, surfaceBg) {
		return hoverBg
	}
	return surfaceBg
}

func removeTODOs(style map[string]any) {
	for k, v := range style {
		if isTodoValue(v) {
			delete(style, k)
		}
	}
}

func shouldPruneStyle(style map[string]any) bool {
	if len(style) == 0 {
		return false
	}
	if _, ok := style["syntax"]; ok {
		if len(style) == 1 {
			return false
		}
	}
	if _, ok := style["players"]; ok {
		if len(style) == 1 {
			return false
		}
		if len(style) == 2 {
			if _, ok := style["syntax"]; ok {
				return false
			}
		}
	}
	if _, ok := style["background"]; ok {
		return true
	}
	if _, ok := style["editor.background"]; ok {
		return true
	}
	if _, ok := style["text"]; ok {
		return true
	}
	if _, ok := style["terminal.foreground"]; ok {
		return true
	}
	return len(style) > 20
}
