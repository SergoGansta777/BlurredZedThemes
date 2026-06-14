package themeutil

func SemanticColors(roles map[string]string, semantic map[string]string) map[string]string {
	role := func(n string) string { return roles[n] }
	out := map[string]string{
		"error":       role("love"),
		"warning":     role("gold"),
		"info":        role("foam"),
		"success":     role("pine"),
		"conflict":    role("rose"),
		"created":     role("foam"),
		"deleted":     role("love"),
		"modified":    role("gold"),
		"renamed":     role("iris"),
		"hidden":      role("muted"),
		"hint":        role("subtle"),
		"ignored":     role("muted"),
		"unreachable": role("muted"),
		"predictive":  role("muted"),
	}
	for k, v := range semantic {
		out[k] = v
	}
	return out
}

func SemanticColor(roles map[string]string, semantic map[string]string, name string) string {
	return SemanticColors(roles, semantic)[name]
}
