package themeutil

import "maps"

var importedSyntaxKeyRenames = map[string]string{
	"ion.special": "function.special",
}

func NormalizeImportedStyle(style map[string]any) map[string]any {
	if style == nil {
		return nil
	}

	out := maps.Clone(style)
	syntax, ok := style["syntax"].(map[string]any)
	if !ok {
		return out
	}

	normalizedSyntax := maps.Clone(syntax)
	for from, to := range importedSyntaxKeyRenames {
		value, ok := normalizedSyntax[from]
		if !ok {
			continue
		}
		if _, exists := normalizedSyntax[to]; !exists {
			normalizedSyntax[to] = value
		}
		delete(normalizedSyntax, from)
	}

	out["syntax"] = normalizedSyntax
	return out
}
