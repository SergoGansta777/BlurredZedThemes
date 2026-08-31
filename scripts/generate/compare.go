package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"zed-themes/scripts/themeutil"
)

func compareAndMaybeUpdatePalette(cfg Config, palette Palette, alphaCfg AlphaConfig, generated map[string]any) error {
	reference, err := readThemeStyle(cfg.ComparePath)
	if err != nil {
		return fmt.Errorf("read reference theme: %w", err)
	}

	missing, extra, diffs := diffStyle(reference, generated)
	fmt.Printf("compare %s\n", cfg.ComparePath)
	fmt.Printf("  missing in generated: %d\n", len(missing))
	fmt.Printf("  extra in generated: %d\n", len(extra))
	fmt.Printf("  value diffs: %d\n", len(diffs))

	if !shouldUpdatePalette(cfg) {
		return nil
	}

	updated := palette
	applyStyleKeyUpdates(&updated, reference, cfg.WriteStyleKeys)

	if cfg.WriteAlpha {
		if updated.Alpha.Light == nil {
			updated.Alpha.Light = map[string]string{}
		}
		if updated.Alpha.Dark == nil {
			updated.Alpha.Dark = map[string]string{}
		}
		applyAlphaOverrides(&updated, alphaCfg, reference)
	}

	return themeutil.WriteJSONFile(cfg.PalettePath, updated)
}

func shouldUpdatePalette(cfg Config) bool {
	return cfg.WriteStyleKeys != "" || cfg.WriteAlpha
}

func applyStyleKeyUpdates(palette *Palette, reference map[string]any, keysCSV string) {
	if keysCSV == "" {
		return
	}
	if palette.Style == nil {
		palette.Style = map[string]any{}
	}
	for _, key := range parseCommaList(keysCSV) {
		if value, ok := reference[key]; ok {
			palette.Style[key] = value
			continue
		}
		delete(palette.Style, key)
	}
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
	family, err := themeutil.ReadJSONFile[themeutil.ThemeFamily](path)
	if err != nil {
		return nil, err
	}
	style, err := themeutil.FirstThemeStyle(family)
	if err != nil {
		return nil, err
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
	for _, rule := range themeutil.AlphaRules {
		baseColor := alphaBaseValue(*palette, rule)
		if baseColor == "" {
			continue
		}
		for _, styleKey := range rule.StyleKeys {
			refValue, ok := reference[styleKey].(string)
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
