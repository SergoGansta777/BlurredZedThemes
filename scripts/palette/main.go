package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"zed-themes/scripts/themeutil"
)

type Palette = themeutil.Palette
type Meta = themeutil.Meta
type AlphaConfig = themeutil.AlphaConfig

type Config struct {
	ThemePath string
	OutPath   string
	StyleKeys string
	AlphaPath string
	WithAlpha bool
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	cfg := parseFlags()
	if cfg.ThemePath == "" {
		return errors.New("missing --theme")
	}

	theme, err := themeutil.ReadJSONFile[map[string]any](cfg.ThemePath)
	if err != nil {
		return fmt.Errorf("read theme: %w", err)
	}
	style, err := themeStyle(theme)
	if err != nil {
		return err
	}
	style = themeutil.NormalizeImportedStyle(style)

	palette := buildPalette(theme, style)
	if cfg.StyleKeys != "" {
		palette.Style = pickStyleKeys(style, cfg.StyleKeys)
	}
	if cfg.WithAlpha {
		alphaCfg, err := themeutil.ReadJSONFile[AlphaConfig](cfg.AlphaPath)
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("read alpha: %w", err)
		}
		palette.Alpha = inferAlphaOverrides(palette, alphaCfg, style)
	}

	if cfg.OutPath == "" {
		cfg.OutPath = defaultPalettePath(cfg.ThemePath)
	}

	if err := os.MkdirAll(filepath.Dir(cfg.OutPath), 0o755); err != nil {
		return fmt.Errorf("create output dir: %w", err)
	}

	return themeutil.WriteJSONFile(cfg.OutPath, palette)
}

func parseFlags() Config {
	var cfg Config
	flag.StringVar(&cfg.ThemePath, "theme", "", "path to theme json")
	flag.StringVar(&cfg.OutPath, "out", "", "output palette json path")
	flag.StringVar(&cfg.StyleKeys, "style-keys", "", "comma-separated style keys to copy into palette style")
	flag.StringVar(&cfg.AlphaPath, "alpha", "palettes/alpha.json", "path to alpha config")
	flag.BoolVar(&cfg.WithAlpha, "with-alpha", false, "derive alpha overrides from theme")
	flag.Parse()
	return cfg
}

func defaultPalettePath(themePath string) string {
	base := strings.TrimSuffix(filepath.Base(themePath), filepath.Ext(themePath))
	return filepath.Join("palettes", base+".json")
}

func buildPalette(theme map[string]any, style map[string]any) Palette {
	backgroundAppearance := stringValue(style, "background.appearance")
	blurMode := ""
	if strings.EqualFold(backgroundAppearance, "blurred") && isTransparentColor(stringValue(style, "editor.background")) {
		blurMode = "flat"
	}
	return Palette{
		Meta: Meta{
			Name:                 stringField(theme, "name"),
			Author:               stringField(theme, "author"),
			Appearance:           themeString(theme, "appearance"),
			ThemeName:            themeString(theme, "name"),
			BackgroundAppearance: backgroundAppearance,
			BlurMode:             blurMode,
		},
		Roles:    deriveRoles(style),
		Semantic: deriveSemantic(style),
		Accents:  stringSlice(style, "accents"),
		Terminal: deriveTerminal(style),
	}
}

func themeStyle(theme map[string]any) (map[string]any, error) {
	themes, ok := theme["themes"].([]any)
	if !ok || len(themes) == 0 {
		return nil, errors.New("invalid theme: missing themes array")
	}
	first, ok := themes[0].(map[string]any)
	if !ok {
		return nil, errors.New("invalid theme: themes[0] not object")
	}
	style, ok := first["style"].(map[string]any)
	if !ok {
		return nil, errors.New("invalid theme: missing style map")
	}
	return style, nil
}

func deriveRoles(style map[string]any) map[string]string {
	role := map[string]string{}
	role["surface"] = stripAlpha(stringValue(style, "editor.background"))
	role["base"] = stripAlpha(stringValue(style, "background"))
	role["overlay"] = stripAlpha(stringValue(style, "editor.active_line.background"))
	role["muted"] = stringValue(style, "text.muted")
	role["subtle"] = stringValue(style, "text.placeholder")
	role["text"] = stringValue(style, "text")

	role["love"] = stringValue(style, "error")
	role["gold"] = firstNonEmpty(stringValue(style, "warning"), stringValue(style, "modified"))
	role["rose"] = firstNonEmpty(stringValue(style, "modified"), stringValue(style, "conflict"))
	role["pine"] = firstNonEmpty(stringValue(style, "info"), stringValue(style, "success"))
	role["foam"] = firstNonEmpty(stringValue(style, "text.accent"), stringValue(style, "link_text.hover"))
	role["iris"] = firstNonEmpty(stringValue(style, "renamed"), stringValue(style, "keyword"))

	role["highlight_low"] = stripAlpha(stringValue(style, "element.hover"))
	role["highlight_med"] = stripAlpha(stringValue(style, "element.selected"))
	role["highlight_high"] = stripAlpha(stringValue(style, "ghost_element.active"))

	for k, v := range role {
		if v == "" {
			delete(role, k)
		}
	}
	return role
}

func deriveSemantic(style map[string]any) map[string]string {
	keys := []string{
		"error", "warning", "info", "success", "conflict",
		"created", "deleted", "modified", "renamed",
		"hidden", "hint", "ignored", "unreachable", "predictive",
	}
	out := map[string]string{}
	for _, k := range keys {
		if v := stringValue(style, k); v != "" {
			out[k] = v
		}
	}
	return out
}

func deriveTerminal(style map[string]any) map[string]string {
	out := map[string]string{}
	for k, v := range style {
		if strings.HasPrefix(k, "terminal.") {
			if s, ok := v.(string); ok {
				out[k] = s
			}
		}
	}
	return out
}

func themeString(theme map[string]any, key string) string {
	themes, ok := theme["themes"].([]any)
	if !ok || len(themes) == 0 {
		return ""
	}
	obj, ok := themes[0].(map[string]any)
	if !ok {
		return ""
	}
	if v, ok := obj[key].(string); ok {
		return v
	}
	return ""
}

func stringField(m map[string]any, key string) string {
	if v, ok := m[key].(string); ok {
		return v
	}
	return ""
}

func stringValue(style map[string]any, key string) string {
	if v, ok := style[key].(string); ok {
		return v
	}
	return ""
}

func isTransparentColor(value string) bool {
	if value == "" {
		return false
	}
	v := strings.TrimPrefix(strings.ToUpper(value), "#")
	if len(v) == 8 && v[6:8] == "00" {
		return true
	}
	return v == "00000000"
}

func stringSlice(style map[string]any, key string) []string {
	var out []string
	arr, ok := style[key].([]any)
	if !ok {
		return out
	}
	for _, v := range arr {
		if s, ok := v.(string); ok {
			out = append(out, s)
		}
	}
	return out
}

func pickStyleKeys(style map[string]any, keysCSV string) map[string]any {
	out := map[string]any{}
	for _, raw := range strings.Split(keysCSV, ",") {
		key := strings.TrimSpace(raw)
		if key == "" {
			continue
		}
		if v, ok := style[key]; ok {
			out[key] = v
		}
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

func inferAlphaOverrides(p Palette, base AlphaConfig, style map[string]any) AlphaConfig {
	out := AlphaConfig{Light: map[string]string{}, Dark: map[string]string{}}
	appearance := strings.ToLower(p.Meta.Appearance)
	defaults := map[string]string{}
	if appearance == "light" {
		for k, v := range base.Light {
			defaults[k] = v
		}
	} else {
		for k, v := range base.Dark {
			defaults[k] = v
		}
	}

	role := func(name string) string { return p.Roles[name] }
	semantic := func(name string) string { return themeutil.SemanticColor(p.Roles, p.Semantic, name) }

	specs := []struct {
		alphaKey string
		styleKey string
		base     func() string
	}{
		{"ui", "background", func() string { return role("surface") }},
		{"ui_inactive", "title_bar.inactive_background", func() string { return role("surface") }},
		{"surface", "surface.background", func() string { return role("surface") }},
		{"elevated", "elevated_surface.background", func() string { return role("surface") }},
		{"overlay", "panel.overlay_background", func() string { return role("surface") }},
		{"subheader", "editor.subheader.background", func() string { return role("surface") }},
		{"active_line", "editor.active_line.background", func() string { return role("overlay") }},
		{"active_line_number", "editor.active_line_number", func() string { return role("text") }},
		{"highlighted_line", "editor.highlighted_line.background", func() string { return role("overlay") }},
		{"line_number", "editor.line_number", func() string { return role("text") }},
		{"text_disabled", "text.disabled", func() string { return role("text") }},
		{"text_disabled", "icon.disabled", func() string { return role("text") }},
		{"text_placeholder", "text.placeholder", func() string { return role("text") }},
		{"element_active", "element.active", func() string { return role("highlight_med") }},
		{"element_selected", "element.selected", func() string { return role("highlight_med") }},
		{"element_hover", "element.hover", func() string { return role("highlight_low") }},
		{"element_disabled", "element.disabled", func() string { return role("surface") }},
		{"ghost_active", "ghost_element.active", func() string { return role("highlight_high") }},
		{"ghost_selected", "ghost_element.selected", func() string { return role("highlight_high") }},
		{"ghost_hover", "ghost_element.hover", func() string { return role("highlight_low") }},
		{"ghost_disabled", "ghost_element.disabled", func() string { return role("surface") }},
		{"border_variant", "border.variant", func() string { return role("foam") }},
		{"border_focused", "border.focused", func() string { return role("foam") }},
		{"border_selected", "border.selected", func() string { return role("iris") }},
		{"border_disabled", "border.disabled", func() string { return role("muted") }},
		{"tab_active", "tab.active_background", func() string { return role("surface") }},
		{"conflict_marker", "version_control.conflict_marker.ours", func() string { return semantic("warning") }},
		{"panel_focus_border", "panel.focused_border", func() string { return role("muted") }},
		{"panel_indent_guide", "panel.indent_guide", func() string { return role("muted") }},
		{"panel_indent_guide_active", "panel.indent_guide_active", func() string { return role("subtle") }},
		{"pane_focus_border", "pane.focused_border", func() string { return role("muted") }},
		{"pane_group_border", "pane_group.border", func() string { return role("muted") }},
		{"scrollbar_thumb", "scrollbar.thumb.background", func() string { return role("muted") }},
		{"scrollbar_thumb_hover", "scrollbar.thumb.hover_background", func() string { return role("muted") }},
		{"scrollbar_track", "scrollbar.track.background", func() string { return role("surface") }},
		{"scrollbar_track_border", "scrollbar.track.border", func() string { return role("text") }},
		{"search_match", "search.match_background", func() string { return role("foam") }},
		{"search_active", "search.active_match_background", func() string { return role("rose") }},
		{"debugger_line", "editor.debugger_active_line.background", func() string { return role("rose") }},
		{"indent_guide", "editor.indent_guide", func() string { return role("muted") }},
		{"indent_guide_active", "editor.indent_guide_active", func() string { return role("subtle") }},
		{"wrap_guide", "editor.wrap_guide", func() string { return role("muted") }},
		{"active_wrap_guide", "editor.active_wrap_guide", func() string { return role("muted") }},
		{"doc_highlight_read", "editor.document_highlight.read_background", func() string { return role("foam") }},
		{"doc_highlight_write", "editor.document_highlight.write_background", func() string { return role("muted") }},
		{"doc_highlight_bracket", "editor.document_highlight.bracket_background", func() string { return role("iris") }},
		{"drop_target", "drop_target.background", func() string { return role("text") }},
		{"minimap_bg", "minimap.thumb.background", func() string { return role("foam") }},
		{"minimap_hover", "minimap.thumb.hover_background", func() string { return role("foam") }},
		{"minimap_active", "minimap.thumb.active_background", func() string { return role("foam") }},
	}

	overrides := map[string]string{}
	for _, spec := range specs {
		baseColor := spec.base()
		if baseColor == "" {
			continue
		}
		refValue, ok := style[spec.styleKey].(string)
		if !ok || refValue == "" {
			continue
		}
		alpha, ok := themeutil.InferAlpha(refValue, baseColor)
		if !ok {
			continue
		}
		if def := defaults[spec.alphaKey]; def != "" && strings.EqualFold(def, alpha) {
			continue
		}
		overrides[spec.alphaKey] = strings.ToUpper(alpha)
	}

	selectionAlpha := themeutil.InferSelectionAlpha(style)
	if selectionAlpha != "" && !strings.EqualFold(selectionAlpha, defaults["selection"]) {
		overrides["selection"] = strings.ToUpper(selectionAlpha)
	}

	if appearance == "light" {
		for k, v := range overrides {
			out.Light[k] = v
		}
	} else {
		for k, v := range overrides {
			out.Dark[k] = v
		}
	}

	if len(out.Light) == 0 && len(out.Dark) == 0 {
		return AlphaConfig{}
	}
	return out
}

func stripAlpha(hex string) string {
	h := strings.TrimPrefix(hex, "#")
	if len(h) == 8 {
		return "#" + h[:6]
	}
	if strings.HasPrefix(hex, "#") {
		return hex
	}
	return hex
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}
