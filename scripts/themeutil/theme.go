package themeutil

import "errors"

type ThemeFamily struct {
	Schema string         `json:"$schema"`
	Author string         `json:"author"`
	Name   string         `json:"name"`
	Themes []ThemeContent `json:"themes"`
}

type ThemeContent struct {
	Appearance string         `json:"appearance"`
	Name       string         `json:"name"`
	Style      map[string]any `json:"style"`
}

func FirstTheme(family ThemeFamily) (ThemeContent, error) {
	if len(family.Themes) == 0 {
		return ThemeContent{}, errors.New("invalid theme: missing themes array")
	}
	return family.Themes[0], nil
}

func FirstThemeStyle(family ThemeFamily) (map[string]any, error) {
	theme, err := FirstTheme(family)
	if err != nil {
		return nil, err
	}
	if theme.Style == nil {
		return nil, errors.New("invalid theme: missing style map")
	}
	return theme.Style, nil
}
