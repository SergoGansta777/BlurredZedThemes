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

	family, err := themeutil.ReadJSONFile[themeutil.ThemeFamily](cfg.ThemePath)
	if err != nil {
		return fmt.Errorf("read theme: %w", err)
	}
	theme, err := themeutil.FirstTheme(family)
	if err != nil {
		return err
	}
	style, err := themeutil.FirstThemeStyle(family)
	if err != nil {
		return err
	}
	style = themeutil.NormalizeImportedStyle(style)

	palette := buildPalette(family, theme, style)
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

func buildPalette(family themeutil.ThemeFamily, theme themeutil.ThemeContent, style map[string]any) Palette {
	backgroundAppearance := stringValue(style, "background.appearance")
	blurMode := ""
	if strings.EqualFold(backgroundAppearance, "blurred") && isTransparentColor(stringValue(style, "editor.background")) {
		blurMode = "flat"
	}
	return Palette{
		Meta: Meta{
			Name:                 family.Name,
			Author:               family.Author,
			Appearance:           theme.Appearance,
			ThemeName:            theme.Name,
			BackgroundAppearance: backgroundAppearance,
			BlurMode:             blurMode,
		},
		Roles:    deriveRoles(style),
		Semantic: deriveSemantic(style),
		Accents:  stringSlice(style, "accents"),
		Terminal: deriveTerminal(style),
	}
}

func deriveRoles(style map[string]any) map[string]string {
	role := map[string]string{}
	role["surface"] = themeutil.StripAlpha(stringValue(style, "editor.background"))
	role["base"] = themeutil.StripAlpha(stringValue(style, "background"))
	role["overlay"] = themeutil.StripAlpha(stringValue(style, "editor.active_line.background"))
	role["muted"] = stringValue(style, "text.muted")
	role["subtle"] = stringValue(style, "text.placeholder")
	role["text"] = stringValue(style, "text")

	role["love"] = stringValue(style, "error")
	role["gold"] = firstNonEmpty(stringValue(style, "warning"), stringValue(style, "modified"))
	role["rose"] = firstNonEmpty(stringValue(style, "modified"), stringValue(style, "conflict"))
	role["pine"] = firstNonEmpty(stringValue(style, "info"), stringValue(style, "success"))
	role["foam"] = firstNonEmpty(stringValue(style, "text.accent"), stringValue(style, "link_text.hover"))
	role["iris"] = firstNonEmpty(stringValue(style, "renamed"), stringValue(style, "keyword"))

	role["highlight_low"] = themeutil.StripAlpha(stringValue(style, "element.hover"))
	role["highlight_med"] = themeutil.StripAlpha(stringValue(style, "element.selected"))
	role["highlight_high"] = themeutil.StripAlpha(stringValue(style, "ghost_element.active"))

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

	overrides := map[string]string{}
	for _, rule := range themeutil.AlphaRules {
		baseColor := alphaBaseValue(p, rule)
		if baseColor == "" {
			continue
		}
		for _, styleKey := range rule.StyleKeys {
			refValue, ok := style[styleKey].(string)
			if !ok || refValue == "" {
				continue
			}
			alpha, ok := themeutil.InferAlpha(refValue, baseColor)
			if !ok {
				continue
			}
			if def := defaults[rule.AlphaKey]; def != "" && strings.EqualFold(def, alpha) {
				continue
			}
			overrides[rule.AlphaKey] = strings.ToUpper(alpha)
		}
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

func alphaBaseValue(p Palette, rule themeutil.AlphaRule) string {
	switch rule.BaseKind {
	case themeutil.AlphaBaseRole:
		return themeutil.StripAlpha(p.Roles[rule.BaseKey])
	case themeutil.AlphaBaseSemantic:
		return themeutil.SemanticColor(p.Roles, p.Semantic, rule.BaseKey)
	case themeutil.AlphaBaseTerminal:
		return themeutil.StripAlpha(p.Terminal[rule.BaseKey])
	default:
		return ""
	}
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}
