package main

import (
	"strings"

	"zed-themes/scripts/themeutil"
)

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

func alphaHexOrDefault(appearance string, cfg AlphaConfig, key string, fallback string) string {
	if key != "" {
		if value, ok := alphaValue(appearance, cfg, key); ok {
			return value
		}
	}
	return fallback
}
