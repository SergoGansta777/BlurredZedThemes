package themeutil

import "strings"

type AlphaConfig struct {
	Light map[string]string `json:"light"`
	Dark  map[string]string `json:"dark"`
}

func MergeAlphaConfig(dst *AlphaConfig, src AlphaConfig) {
	for k, v := range src.Light {
		if dst.Light == nil {
			dst.Light = map[string]string{}
		}
		if v != "" {
			dst.Light[k] = v
		}
	}
	for k, v := range src.Dark {
		if dst.Dark == nil {
			dst.Dark = map[string]string{}
		}
		if v != "" {
			dst.Dark[k] = v
		}
	}
}

func AlphaFor(appearance string, cfg AlphaConfig, key string) string {
	if strings.EqualFold(appearance, "light") {
		if v := cfg.Light[key]; v != "" {
			return v
		}
	} else {
		if v := cfg.Dark[key]; v != "" {
			return v
		}
	}
	return ""
}

func WithAlpha(hex string, alpha string) string {
	h := strings.TrimPrefix(hex, "#")
	if len(h) == 8 {
		h = h[:6]
	}
	if len(h) == 6 && len(alpha) == 2 {
		if strings.EqualFold(alpha, "FF") {
			return "#" + strings.ToUpper(h)
		}
		return "#" + strings.ToUpper(h+alpha)
	}
	if strings.HasPrefix(hex, "#") {
		return strings.ToUpper(hex)
	}
	return hex
}

func StripAlpha(hex string) string {
	h := strings.TrimPrefix(hex, "#")
	if len(h) == 8 {
		return "#" + h[:6]
	}
	if strings.HasPrefix(hex, "#") {
		return hex
	}
	return hex
}

func InferAlpha(value, base string) (string, bool) {
	v := strings.TrimPrefix(strings.ToUpper(value), "#")
	b := strings.TrimPrefix(strings.ToUpper(base), "#")
	if len(b) != 6 {
		return "", false
	}
	if len(v) == 8 && strings.HasPrefix(v, b) {
		return v[6:8], true
	}
	if len(v) == 6 && v == b {
		return "FF", true
	}
	return "", false
}

func InferSelectionAlpha(style map[string]any) string {
	raw, ok := style["players"].([]any)
	if !ok {
		return ""
	}
	for _, entry := range raw {
		obj, ok := entry.(map[string]any)
		if !ok {
			continue
		}
		bg, _ := obj["background"].(string)
		if bg == "" {
			bg, _ = obj["cursor"].(string)
		}
		sel, _ := obj["selection"].(string)
		if bg == "" || sel == "" {
			continue
		}
		if alpha, ok := InferAlpha(sel, bg); ok {
			return alpha
		}
	}
	return ""
}
