package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"sort"
	"strings"

	"zed-themes/scripts/themeutil"
)

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
			if themeutil.IsStandardizedKey(key) {
				continue
			}
			updated.Overrides[key] = reference[key]
		}
		for _, key := range diffs {
			if key == "syntax" || key == "players" {
				updated.Style[key] = reference[key]
				continue
			}
			if themeutil.IsStandardizedKey(key) {
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
	keys := make([]string, 0, len(themeutil.AlphaRules))
	for _, rule := range themeutil.AlphaRules {
		for _, styleKey := range rule.StyleKeys {
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
