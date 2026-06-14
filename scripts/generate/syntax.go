package main

import "maps"

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

	applySyntaxAliases(syntax)
	style["syntax"] = syntax
}

func applySyntaxAliases(syntax map[string]any) {
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
