package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"zed-themes/scripts/themeutil"
)

const (
	todoValue        = "TODO"
	transparentColor = "#00000000"
	selectionAlpha   = "4D"
	wipSuffix        = " (WIP)"
	hoverBlendAlpha  = 0xBF
)

const (
	backgroundAppearanceBlurred = "blurred"
	backgroundAppearanceOpaque  = "opaque"
	blurModeFlat                = "flat"
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

type Palette = themeutil.Palette
type Meta = themeutil.Meta
type AlphaConfig = themeutil.AlphaConfig

type Config struct {
	TemplatePath     string
	PalettePath      string
	AlphaPath        string
	OutPath          string
	PruneStyle       bool
	ComparePath      string
	WriteOverrides   bool
	WriteStyleKeys   string
	WriteAlpha       bool
	PruneAlpha       bool
	PruneOverrides   bool
	RewriteOverrides bool
	WIP              bool
	KeepTODOs        bool
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	cfg := parseFlags()
	if cfg.PalettePath == "" {
		return errors.New("missing --palette")
	}

	palette, template, alphaCfg, err := loadInputs(cfg)
	if err != nil {
		return err
	}

	outPath := resolveOutputPath(cfg)
	style, theme := buildThemeOutput(cfg, palette, template, alphaCfg)

	if cfg.ComparePath != "" {
		if err := compareAndMaybeUpdatePalette(cfg, palette, template, alphaCfg, style); err != nil {
			return err
		}
	}

	if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
		return fmt.Errorf("create output dir: %w", err)
	}

	if err := writeThemeOutput(outPath, "theme", theme); err != nil {
		return err
	}

	for _, variant := range derivedOutputVariants(palette.Meta, outPath) {
		variantPalette := palette
		variantPalette.Meta = variant.Meta
		_, variantTheme := buildThemeOutput(cfg, variantPalette, template, alphaCfg)
		if err := writeThemeOutput(variant.OutPath, variant.Label, variantTheme); err != nil {
			return err
		}
	}
	return nil
}

func parseFlags() Config {
	cfg := Config{}
	flag.StringVar(&cfg.TemplatePath, "template", "templates/base-style.json", "path to base style template")
	flag.StringVar(&cfg.PalettePath, "palette", "", "path to palette json")
	flag.StringVar(&cfg.AlphaPath, "alpha", "palettes/alpha.json", "path to alpha config")
	flag.StringVar(&cfg.OutPath, "out", "", "output theme json path")
	flag.BoolVar(&cfg.PruneStyle, "prune", true, "drop keys not present in palette style when available")
	flag.StringVar(&cfg.ComparePath, "compare", "", "reference theme json to compare generated style against")
	flag.BoolVar(&cfg.WriteOverrides, "write-overrides", false, "update palette overrides to match reference")
	flag.StringVar(&cfg.WriteStyleKeys, "write-style-keys", "", "comma-separated top-level style keys to copy from reference into palette style")
	flag.BoolVar(&cfg.WriteAlpha, "write-alpha", false, "update palette alpha overrides to match reference")
	flag.BoolVar(&cfg.PruneAlpha, "prune-alpha-overrides", false, "remove alpha-derived overrides after writing alpha")
	flag.BoolVar(&cfg.PruneOverrides, "prune-overrides", false, "remove overrides that are already derived without changing the generated style")
	flag.BoolVar(&cfg.RewriteOverrides, "rewrite-overrides", false, "replace overrides with only reference diffs (excluding standardized keys)")
	flag.BoolVar(&cfg.WIP, "wip", true, "append WIP suffix to names and filenames")
	flag.BoolVar(&cfg.KeepTODOs, "keep-todos", false, "keep TODO values for debugging")
	flag.Parse()
	return cfg
}

func loadInputs(cfg Config) (Palette, map[string]any, AlphaConfig, error) {
	palette, err := themeutil.ReadJSONFile[Palette](cfg.PalettePath)
	if err != nil {
		return Palette{}, nil, AlphaConfig{}, fmt.Errorf("read palette: %w", err)
	}

	template, err := themeutil.ReadJSONFile[map[string]any](cfg.TemplatePath)
	if err != nil {
		return Palette{}, nil, AlphaConfig{}, fmt.Errorf("read template: %w", err)
	}

	alphaCfg, err := themeutil.ReadJSONFile[AlphaConfig](cfg.AlphaPath)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return Palette{}, nil, AlphaConfig{}, fmt.Errorf("read alpha: %w", err)
	}
	themeutil.MergeAlphaConfig(&alphaCfg, palette.Alpha)

	return palette, template, alphaCfg, nil
}

func resolveOutputPath(cfg Config) string {
	if cfg.OutPath != "" {
		return cfg.OutPath
	}
	base := strings.TrimSuffix(filepath.Base(cfg.PalettePath), filepath.Ext(cfg.PalettePath))
	if cfg.WIP {
		return filepath.Join("generated", "themes", base+".wip.json")
	}
	return filepath.Join("generated", "themes", base+".json")
}

func buildTheme(p Palette, style map[string]any, wip bool) map[string]any {
	name := p.Meta.Name
	themeName := p.Meta.ThemeName
	if wip {
		name = withWIPSuffix(name)
		themeName = withWIPSuffix(themeName)
	}
	return map[string]any{
		"$schema": themeutil.ThemeSchemaURL,
		"name":    name,
		"author":  p.Meta.Author,
		"themes": []any{
			map[string]any{
				"appearance": p.Meta.Appearance,
				"name":       themeName,
				"style":      style,
			},
		},
	}
}

func buildThemeOutput(cfg Config, p Palette, template map[string]any, alpha AlphaConfig) (map[string]any, map[string]any) {
	style := buildStyle(template, p, alpha, cfg.PruneStyle)
	if !cfg.KeepTODOs {
		removeTODOs(style)
	}
	return style, buildTheme(p, style, cfg.WIP)
}

func writeThemeOutput(outPath, label string, theme map[string]any) error {
	if outPath == "" {
		return nil
	}
	if err := themeutil.WriteJSONFile(outPath, theme); err != nil {
		return fmt.Errorf("write %s: %w", label, err)
	}
	return nil
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

func alphaValue(appearance string, cfg AlphaConfig, key string) (string, bool) {
	if strings.EqualFold(appearance, "light") {
		if v := cfg.Light[key]; v != "" {
			return v, true
		}
		return "", false
	}
	if v := cfg.Dark[key]; v != "" {
		return v, true
	}
	return "", false
}

func alphaBaseValue(p Palette, rule alphaRule) string {
	base, _ := alphaBaseValueForRule(p, rule)
	return base
}

func alphaBaseValueForRule(p Palette, rule alphaRule) (string, bool) {
	if value := derivedColor(p, rule.alphaKey); value != "" {
		return themeutil.StripAlpha(value), true
	}

	switch rule.baseKind {
	case alphaBaseRole:
		return roleOpaque(p, rule.baseKey), false
	case alphaBaseSemantic:
		return semanticOf(p, rule.baseKey), false
	case alphaBaseTerminal:
		return terminalBaseOf(p, rule.baseKey), false
	default:
		return "", false
	}
}

func applyAlphaRules(style map[string]any, p Palette, alpha AlphaConfig) {
	appearance := p.Meta.Appearance
	for _, rule := range alphaRules {
		alphaHex, ok := alphaValue(appearance, alpha, rule.alphaKey)
		if !ok {
			continue
		}
		base, fromDerived := alphaBaseValueForRule(p, rule)
		if base == "" {
			continue
		}
		for _, styleKey := range rule.styleKeys {
			if current, ok := style[styleKey].(string); ok && current != "" && !isTodoValue(current) {
				if !rule.force && !fromDerived && !strings.EqualFold(themeutil.StripAlpha(current), themeutil.StripAlpha(base)) {
					continue
				}
			}
			style[styleKey] = withAlpha(base, alphaHex)
		}
	}
}

func compareAndMaybeUpdatePalette(cfg Config, palette Palette, template map[string]any, alphaCfg AlphaConfig, generated map[string]any) error {
	reference, err := readThemeStyle(cfg.ComparePath)
	if err != nil {
		return fmt.Errorf("read reference theme: %w", err)
	}

	missing, extra, diffs := diffStyle(reference, generated)
	fmt.Printf("compare %s\n", cfg.ComparePath)
	fmt.Printf("  missing in generated: %d\n", len(missing))
	fmt.Printf("  extra in generated: %d\n", len(extra))
	fmt.Printf("  value diffs: %d\n", len(diffs))

	if !cfg.WriteOverrides && cfg.WriteStyleKeys == "" && !cfg.WriteAlpha && !cfg.PruneAlpha && !cfg.PruneOverrides {
		return nil
	}

	updated := palette
	if cfg.WriteStyleKeys != "" {
		if updated.Style == nil {
			updated.Style = map[string]any{}
		}
		for _, key := range parseCommaList(cfg.WriteStyleKeys) {
			if value, ok := reference[key]; ok {
				updated.Style[key] = value
				continue
			}
			delete(updated.Style, key)
		}
	}
	if cfg.WriteOverrides {
		if updated.Overrides == nil {
			updated.Overrides = map[string]any{}
		}
		if cfg.RewriteOverrides {
			updated.Overrides = map[string]any{}
		}
		if updated.Style == nil {
			updated.Style = map[string]any{}
		}
		for _, key := range missing {
			if key == "syntax" || key == "players" {
				updated.Style[key] = reference[key]
				continue
			}
			if isStandardizedKey(key) {
				continue
			}
			updated.Overrides[key] = reference[key]
		}
		for _, key := range diffs {
			if key == "syntax" || key == "players" {
				updated.Style[key] = reference[key]
				continue
			}
			if isStandardizedKey(key) {
				continue
			}
			updated.Overrides[key] = reference[key]
		}
	}

	if cfg.WriteAlpha {
		if updated.Alpha.Light == nil {
			updated.Alpha.Light = map[string]string{}
		}
		if updated.Alpha.Dark == nil {
			updated.Alpha.Dark = map[string]string{}
		}
		applyAlphaOverrides(&updated, alphaCfg, reference)
	}

	if cfg.PruneAlpha {
		pruneAlphaOverrides(&updated, template, alphaCfg, reference)
	}
	if cfg.PruneOverrides {
		pruneRedundantOverrides(&updated, template, alphaCfg, reference, cfg.PruneStyle)
	}

	return themeutil.WriteJSONFile(cfg.PalettePath, updated)
}

func parseCommaList(value string) []string {
	if value == "" {
		return nil
	}
	parts := strings.Split(value, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		out = append(out, part)
	}
	return out
}

func readThemeStyle(path string) (map[string]any, error) {
	theme, err := themeutil.ReadJSONFile[map[string]any](path)
	if err != nil {
		return nil, err
	}
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
	return themeutil.NormalizeImportedStyle(style), nil
}

func diffStyle(reference, generated map[string]any) ([]string, []string, []string) {
	var missing, extra, diffs []string

	for k := range reference {
		if _, ok := generated[k]; !ok {
			missing = append(missing, k)
		}
	}
	for k := range generated {
		if _, ok := reference[k]; !ok {
			extra = append(extra, k)
		}
	}
	for k, v := range reference {
		if gv, ok := generated[k]; ok && !valuesEqual(v, gv) {
			diffs = append(diffs, k)
		}
	}
	return missing, extra, diffs
}

func valuesEqual(a, b any) bool {
	if as, ok := a.(string); ok {
		if bs, ok := b.(string); ok && hexColorsEqual(as, bs) {
			return true
		}
	}
	ab, err := json.Marshal(a)
	if err != nil {
		return false
	}
	bb, err := json.Marshal(b)
	if err != nil {
		return false
	}
	return string(ab) == string(bb)
}

func hexColorsEqual(a, b string) bool {
	if !strings.HasPrefix(a, "#") || !strings.HasPrefix(b, "#") {
		return false
	}
	ah := strings.TrimPrefix(a, "#")
	bh := strings.TrimPrefix(b, "#")
	if len(ah) != len(bh) {
		return false
	}
	if len(ah) != 6 && len(ah) != 8 {
		return false
	}
	return strings.EqualFold(ah, bh)
}

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

func applyDerivedSyntax(style map[string]any, p Palette) {
	if len(p.Roles) == 0 {
		return
	}
	role := func(name string) string { return roleValue(p, name) }

	syntax := map[string]any{
		"embedded":        map[string]any{"color": role("text")},
		"enum":            map[string]any{"color": role("foam")},
		"hint":            map[string]any{"color": firstNonEmpty(semanticOf(p, "hint"), role("muted"))},
		"label":           map[string]any{"color": role("foam")},
		"parameter":       map[string]any{"color": role("text")},
		"predictive":      map[string]any{"color": firstNonEmpty(semanticOf(p, "predictive"), role("muted")), "font_style": "italic"},
		"preproc":         map[string]any{"color": role("pine")},
		"primary":         map[string]any{"color": role("text")},
		"text":            map[string]any{"color": role("text")},
		"comment":         map[string]any{"color": role("muted"), "font_style": "italic"},
		"punctuation":     map[string]any{"color": role("subtle")},
		"operator":        map[string]any{"color": role("subtle")},
		"keyword":         map[string]any{"color": role("pine")},
		"string":          map[string]any{"color": role("gold")},
		"number":          map[string]any{"color": role("foam")},
		"boolean":         map[string]any{"color": role("love")},
		"function":        map[string]any{"color": role("rose")},
		"type":            map[string]any{"color": role("foam")},
		"constant":        map[string]any{"color": role("foam")},
		"variable":        map[string]any{"color": role("text")},
		"property":        map[string]any{"color": role("text")},
		"tag":             map[string]any{"color": role("iris")},
		"attribute":       map[string]any{"color": role("rose")},
		"namespace":       map[string]any{"color": role("iris"), "font_style": "italic"},
		"module":          map[string]any{"color": role("iris"), "font_style": "italic"},
		"string.escape":   map[string]any{"color": role("love")},
		"string.regex":    map[string]any{"color": role("gold")},
		"string.special":  map[string]any{"color": role("pine")},
		"link_text":       map[string]any{"color": role("foam")},
		"link_uri":        map[string]any{"color": role("pine"), "font_style": "italic"},
		"emphasis":        map[string]any{"color": role("iris"), "font_style": "italic"},
		"emphasis.strong": map[string]any{"color": role("iris"), "font_weight": 700},
		"title":           map[string]any{"color": role("text"), "font_weight": 800},
		"diff.plus":       map[string]any{"color": firstNonEmpty(semanticOf(p, "created"), role("pine"))},
		"diff.minus":      map[string]any{"color": firstNonEmpty(semanticOf(p, "deleted"), role("love"))},
		"diff.delta":      map[string]any{"color": firstNonEmpty(semanticOf(p, "modified"), role("gold"))},
		"selector":        map[string]any{"color": firstNonEmpty(role("gold"), role("foam"), role("text"))},
		"selector.pseudo": map[string]any{"color": firstNonEmpty(role("foam"), role("pine"), role("text"))},
	}

	if existing, ok := style["syntax"].(map[string]any); ok {
		for k, v := range existing {
			syntax[k] = v
		}
	}

	setSyntaxAlias(syntax, "comment.doc", "comment")
	setSyntaxAlias(syntax, "comment.documentation", "comment.doc", "comment")
	setSyntaxAlias(syntax, "comment.error", "comment")
	setSyntaxAlias(syntax, "comment.hint", "comment")
	setSyntaxAlias(syntax, "comment.info", "comment")
	setSyntaxAlias(syntax, "comment.note", "comment")
	setSyntaxAlias(syntax, "comment.todo", "comment")
	setSyntaxAlias(syntax, "comment.warn", "comment")
	setSyntaxAlias(syntax, "comment.warning", "comment.warn", "comment")
	setSyntaxAlias(syntax, "constant.builtin", "constant")
	setSyntaxAlias(syntax, "constant.character", "constant")
	setSyntaxAlias(syntax, "constant.macro", "constant", "preproc")
	setSyntaxAlias(syntax, "field", "property", "variable")
	setSyntaxAlias(syntax, "function.builtin", "function")
	setSyntaxAlias(syntax, "function.call", "function")
	setSyntaxAlias(syntax, "function.definition", "function")
	setSyntaxAlias(syntax, "function.decorator", "attribute", "function")
	setSyntaxAlias(syntax, "function.macro", "function", "preproc")
	setSyntaxAlias(syntax, "function.method", "function.call", "function")
	setSyntaxAlias(syntax, "function.method.builtin", "function.method", "function.builtin", "function")
	setSyntaxAlias(syntax, "function.method.call", "function.method", "function.call", "function")
	setSyntaxAlias(syntax, "function.special", "function")
	setSyntaxAlias(syntax, "function.special.definition", "function.special", "function.definition", "function")
	setSyntaxAlias(syntax, "keyword.conditional", "keyword")
	setSyntaxAlias(syntax, "keyword.conditional.ternary", "keyword.conditional", "keyword")
	setSyntaxAlias(syntax, "keyword.coroutine", "keyword")
	setSyntaxAlias(syntax, "keyword.debug", "keyword")
	setSyntaxAlias(syntax, "keyword.directive", "keyword", "preproc")
	setSyntaxAlias(syntax, "keyword.directive.define", "keyword.directive", "keyword", "preproc")
	setSyntaxAlias(syntax, "keyword.exception", "keyword")
	setSyntaxAlias(syntax, "keyword.export", "keyword")
	setSyntaxAlias(syntax, "keyword.function", "keyword")
	setSyntaxAlias(syntax, "keyword.import", "keyword")
	setSyntaxAlias(syntax, "keyword.modifier", "keyword")
	setSyntaxAlias(syntax, "keyword.operator", "keyword", "operator")
	setSyntaxAlias(syntax, "keyword.repeat", "keyword")
	setSyntaxAlias(syntax, "keyword.return", "keyword")
	setSyntaxAlias(syntax, "keyword.type", "keyword", "type")
	setSyntaxAlias(syntax, "namespace.crateRoot", "namespace", "module", "type")
	setSyntaxAlias(syntax, "number.float", "number")
	setSyntaxAlias(syntax, "operator.controlFlow", "keyword", "operator")
	setSyntaxAlias(syntax, "punctuation.bracket", "punctuation")
	setSyntaxAlias(syntax, "punctuation.delimiter", "punctuation")
	setSyntaxAlias(syntax, "punctuation.list_marker", "punctuation.special", "punctuation")
	setSyntaxAlias(syntax, "punctuation.markup", "punctuation.list_marker", "punctuation.special", "punctuation")
	setSyntaxAlias(syntax, "punctuation.special", "punctuation")
	setSyntaxAlias(syntax, "punctuation.special.symbol", "punctuation.special", "punctuation")
	setSyntaxAlias(syntax, "string.doc", "string", "comment.doc")
	setSyntaxAlias(syntax, "string.documentation", "string.doc", "string")
	setSyntaxAlias(syntax, "string.regexp", "string.regex")
	setSyntaxAlias(syntax, "string.special.path", "string.special", "string")
	setSyntaxAlias(syntax, "string.special.symbol", "string.special", "string")
	setSyntaxAlias(syntax, "string.special.url", "link_uri", "string.special", "string")
	setSyntaxAlias(syntax, "tag.attribute", "attribute", "property")
	setSyntaxAlias(syntax, "tag.delimiter", "punctuation.delimiter", "punctuation")
	setSyntaxAlias(syntax, "tag.doctype", "tag", "keyword")
	setSyntaxAlias(syntax, "text.literal", "string", "text")
	setSyntaxAlias(syntax, "type.builtin", "type")
	setSyntaxAlias(syntax, "type.class.definition", "type.definition", "type")
	setSyntaxAlias(syntax, "type.definition", "type")
	setSyntaxAlias(syntax, "type.enum", "type", "variant")
	setSyntaxAlias(syntax, "type.enum.member", "variant", "property", "constant")
	setSyntaxAlias(syntax, "type.interface", "type")
	setSyntaxAlias(syntax, "type.parameter", "type", "parameter")
	setSyntaxAlias(syntax, "type.super", "type")
	setSyntaxAlias(syntax, "variable.builtin", "variable", "constant")
	setSyntaxAlias(syntax, "variable.member", "property", "variable")
	setSyntaxAlias(syntax, "variable.parameter", "parameter", "variable")
	setSyntaxAlias(syntax, "variable.special", "variable", "constant")
	setSyntaxAlias(syntax, "variant", "type", "constructor")
	setSyntaxAlias(syntax, "selector", "tag", "type", "string")
	setSyntaxAlias(syntax, "selector.pseudo", "selector", "keyword", "attribute")

	style["syntax"] = syntax
}

func setSyntaxAlias(syntax map[string]any, key string, candidates ...string) {
	if _, ok := syntaxHighlightStyle(syntax, key); ok {
		return
	}
	for _, candidate := range candidates {
		styleDef, ok := syntaxHighlightStyle(syntax, candidate)
		if !ok {
			continue
		}
		syntax[key] = maps.Clone(styleDef)
		return
	}
}

func syntaxHighlightStyle(syntax map[string]any, key string) (map[string]any, bool) {
	value, ok := syntax[key]
	if !ok {
		return nil, false
	}
	styleDef, ok := value.(map[string]any)
	if !ok || len(styleDef) == 0 {
		return nil, false
	}
	return styleDef, true
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
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

func alphaHexOrDefault(appearance string, cfg AlphaConfig, key string, fallback string) string {
	if key != "" {
		if value, ok := alphaValue(appearance, cfg, key); ok {
			return value
		}
	}
	return fallback
}

func applyAlphaOverrides(palette *Palette, base AlphaConfig, reference map[string]any) {
	appearance := strings.ToLower(palette.Meta.Appearance)
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

	overrides := map[string]string{}
	for _, rule := range alphaRules {
		baseColor := alphaBaseValue(*palette, rule)
		if baseColor == "" {
			continue
		}
		for _, styleKey := range rule.styleKeys {
			refValue, ok := reference[styleKey].(string)
			if !ok || refValue == "" {
				continue
			}
			alpha, ok := themeutil.InferAlpha(refValue, baseColor)
			if !ok {
				continue
			}
			if def := defaults[rule.alphaKey]; def != "" && strings.EqualFold(def, alpha) {
				continue
			}
			overrides[rule.alphaKey] = strings.ToUpper(alpha)
		}
	}

	selectionAlpha := themeutil.InferSelectionAlpha(reference)
	if selectionAlpha != "" && !strings.EqualFold(selectionAlpha, defaults["selection"]) {
		overrides["selection"] = strings.ToUpper(selectionAlpha)
	}

	if appearance == "light" {
		for k, v := range overrides {
			palette.Alpha.Light[k] = v
		}
	} else {
		for k, v := range overrides {
			palette.Alpha.Dark[k] = v
		}
	}
}

func pruneAlphaOverrides(palette *Palette, template map[string]any, alphaCfg AlphaConfig, reference map[string]any) {
	if palette.Overrides == nil {
		return
	}
	mergedAlpha := alphaCfg
	themeutil.MergeAlphaConfig(&mergedAlpha, palette.Alpha)
	candidate := *palette
	candidate.Overrides = maps.Clone(palette.Overrides)

	alphaKeys := alphaDerivedKeys()
	for _, key := range alphaKeys {
		delete(candidate.Overrides, key)
	}
	style := buildStyle(template, candidate, mergedAlpha, false)
	for _, key := range alphaKeys {
		refValue, ok := reference[key]
		if !ok {
			continue
		}
		if genValue, ok := style[key]; ok && valuesEqual(refValue, genValue) {
			delete(palette.Overrides, key)
		}
	}
}

func pruneRedundantOverrides(palette *Palette, template map[string]any, alphaCfg AlphaConfig, reference map[string]any, pruneStyle bool) {
	if len(palette.Overrides) == 0 {
		return
	}

	mergedAlpha := alphaCfg
	themeutil.MergeAlphaConfig(&mergedAlpha, palette.Alpha)

	for {
		changed := false
		for _, key := range sortedMapKeys(palette.Overrides) {
			trial := *palette
			trial.Overrides = maps.Clone(palette.Overrides)
			delete(trial.Overrides, key)
			if len(trial.Overrides) == 0 {
				trial.Overrides = nil
			}

			style := buildStyle(template, trial, mergedAlpha, pruneStyle)
			removeTODOs(style)
			missing, extra, diffs := diffStyle(reference, style)
			if len(missing) != 0 || len(extra) != 0 || len(diffs) != 0 {
				continue
			}
			delete(palette.Overrides, key)
			changed = true
		}
		if !changed {
			break
		}
	}
	if len(palette.Overrides) == 0 {
		palette.Overrides = nil
	}
}

func sortedMapKeys(m map[string]any) []string {
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func alphaDerivedKeys() []string {
	seen := map[string]struct{}{}
	keys := make([]string, 0, len(alphaRules))
	for _, rule := range alphaRules {
		for _, styleKey := range rule.styleKeys {
			if styleKey == "" {
				continue
			}
			if _, ok := seen[styleKey]; ok {
				continue
			}
			seen[styleKey] = struct{}{}
			keys = append(keys, styleKey)
		}
	}
	return keys
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
