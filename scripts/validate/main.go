package main

import (
	"errors"
	"flag"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"zed-themes/scripts/themeutil"
)

type config struct {
	ThemesDir     string
	ExtensionPath string
	TemplatePath  string
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	cfg := parseFlags()

	template, err := themeutil.ReadJSONFile[map[string]any](cfg.TemplatePath)
	if err != nil {
		return fmt.Errorf("read template: %w", err)
	}
	allowedStyleKeys := keySet(template)

	paths, err := filepath.Glob(filepath.Join(cfg.ThemesDir, "*.json"))
	if err != nil {
		return fmt.Errorf("glob themes: %w", err)
	}
	if len(paths) == 0 {
		return fmt.Errorf("no themes found in %s", cfg.ThemesDir)
	}
	sort.Strings(paths)

	seenThemeNames := map[string]string{}
	issues := validateExtensionManifest(cfg.ExtensionPath)
	themeCount := 0
	for _, path := range paths {
		family, err := themeutil.ReadJSONFile[themeutil.ThemeFamily](path)
		if err != nil {
			issues = append(issues, fmt.Sprintf("%s: read: %v", path, err))
			continue
		}

		issues = append(issues, validateThemeFamily(path, family, allowedStyleKeys)...)
		for _, theme := range family.Themes {
			if theme.Name == "" {
				continue
			}
			if prev, ok := seenThemeNames[theme.Name]; ok {
				issues = append(issues, fmt.Sprintf("%s: duplicate theme name %q already defined in %s", path, theme.Name, prev))
				continue
			}
			seenThemeNames[theme.Name] = path
		}
		themeCount += len(family.Themes)
	}

	if len(issues) > 0 {
		for _, issue := range issues {
			fmt.Fprintln(os.Stderr, issue)
		}
		return errors.New("validation failed")
	}

	fmt.Printf("validated %d themes in %d files\n", themeCount, len(paths))
	return nil
}

func parseFlags() config {
	var cfg config
	flag.StringVar(&cfg.ThemesDir, "themes-dir", "themes", "directory containing published theme JSON files")
	flag.StringVar(&cfg.ExtensionPath, "extension", "extension.toml", "extension manifest path")
	flag.StringVar(&cfg.TemplatePath, "template", "templates/base-style.json", "base style template used as the allowed style-key set")
	flag.Parse()
	return cfg
}

func validateExtensionManifest(path string) []string {
	if path == "" {
		return nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return []string{fmt.Sprintf("%s: read: %v", path, err)}
	}

	fields := map[string]struct{}{}
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") || !strings.Contains(line, "=") {
			continue
		}
		key, _, _ := strings.Cut(line, "=")
		fields[strings.TrimSpace(key)] = struct{}{}
	}

	required := []string{
		"authors",
		"description",
		"id",
		"name",
		"repository",
		"schema_version",
		"version",
	}
	var issues []string
	for _, key := range required {
		if _, ok := fields[key]; !ok {
			issues = append(issues, fmt.Sprintf("%s: missing %s", path, key))
		}
	}
	return issues
}

func validateThemeFamily(path string, family themeutil.ThemeFamily, allowedStyleKeys map[string]struct{}) []string {
	var issues []string
	if family.Schema != "" && family.Schema != themeutil.ThemeSchemaURL {
		issues = append(issues, fmt.Sprintf("%s: $schema = %q, want %q", path, family.Schema, themeutil.ThemeSchemaURL))
	}
	if strings.TrimSpace(family.Author) == "" {
		issues = append(issues, fmt.Sprintf("%s: missing author", path))
	}
	if strings.TrimSpace(family.Name) == "" {
		issues = append(issues, fmt.Sprintf("%s: missing name", path))
	}
	if len(family.Themes) == 0 {
		issues = append(issues, fmt.Sprintf("%s: missing themes", path))
	}

	for i, theme := range family.Themes {
		location := fmt.Sprintf("%s: themes[%d]", path, i)
		if strings.TrimSpace(theme.Name) == "" {
			issues = append(issues, fmt.Sprintf("%s: missing name", location))
		}
		if theme.Appearance != "dark" && theme.Appearance != "light" {
			issues = append(issues, fmt.Sprintf("%s: appearance = %q, want dark or light", location, theme.Appearance))
		}
		if len(theme.Style) == 0 {
			issues = append(issues, fmt.Sprintf("%s: missing style", location))
			continue
		}
		issues = append(issues, validateStyle(location, theme.Style, allowedStyleKeys)...)
	}
	return issues
}

func validateStyle(location string, style map[string]any, allowedStyleKeys map[string]struct{}) []string {
	var issues []string
	for key, value := range style {
		if _, ok := allowedStyleKeys[key]; !ok {
			issues = append(issues, fmt.Sprintf("%s: style.%s is not in %s", location, key, "templates/base-style.json"))
			continue
		}

		switch key {
		case "accents":
			issues = append(issues, validateColorList(location+".style."+key, value)...)
		case "background.appearance":
			issues = append(issues, validateBackgroundAppearance(location+".style."+key, value)...)
		case "players":
			issues = append(issues, validatePlayers(location+".style."+key, value)...)
		case "syntax":
			issues = append(issues, validateSyntax(location+".style."+key, value)...)
		default:
			issues = append(issues, validateOptionalColor(location+".style."+key, value)...)
		}
	}
	return issues
}

func validateColorList(location string, value any) []string {
	items, ok := value.([]any)
	if !ok {
		return []string{fmt.Sprintf("%s: expected array", location)}
	}

	var issues []string
	for i, item := range items {
		issues = append(issues, validateOptionalColor(fmt.Sprintf("%s[%d]", location, i), item)...)
	}
	return issues
}

func validateBackgroundAppearance(location string, value any) []string {
	s, ok := value.(string)
	if !ok {
		return []string{fmt.Sprintf("%s: expected string", location)}
	}
	if s != "blurred" && s != "opaque" {
		return []string{fmt.Sprintf("%s: value = %q, want blurred or opaque", location, s)}
	}
	return nil
}

func validatePlayers(location string, value any) []string {
	players, ok := value.([]any)
	if !ok {
		return []string{fmt.Sprintf("%s: expected array", location)}
	}

	var issues []string
	for i, rawPlayer := range players {
		player, ok := rawPlayer.(map[string]any)
		playerLocation := fmt.Sprintf("%s[%d]", location, i)
		if !ok {
			issues = append(issues, fmt.Sprintf("%s: expected object", playerLocation))
			continue
		}
		for _, key := range []string{"background", "cursor", "selection"} {
			if value, ok := player[key]; ok {
				issues = append(issues, validateOptionalColor(playerLocation+"."+key, value)...)
			}
		}
		for key := range player {
			if key != "background" && key != "cursor" && key != "selection" {
				issues = append(issues, fmt.Sprintf("%s.%s: unknown player field", playerLocation, key))
			}
		}
	}
	return issues
}

func validateSyntax(location string, value any) []string {
	syntax, ok := value.(map[string]any)
	if !ok {
		return []string{fmt.Sprintf("%s: expected object", location)}
	}

	var issues []string
	for capture, rawEntry := range syntax {
		entryLocation := location + "." + capture
		entry, ok := rawEntry.(map[string]any)
		if !ok {
			issues = append(issues, fmt.Sprintf("%s: expected object", entryLocation))
			continue
		}
		for attr, attrValue := range entry {
			switch attr {
			case "background_color", "color":
				issues = append(issues, validateOptionalColor(entryLocation+"."+attr, attrValue)...)
			case "font_style":
				issues = append(issues, validateOptionalFontStyle(entryLocation+"."+attr, attrValue)...)
			case "font_weight":
				issues = append(issues, validateOptionalFontWeight(entryLocation+"."+attr, attrValue)...)
			default:
				issues = append(issues, fmt.Sprintf("%s.%s: unknown syntax field", entryLocation, attr))
			}
		}
	}
	return issues
}

func validateOptionalColor(location string, value any) []string {
	if value == nil {
		return nil
	}
	s, ok := value.(string)
	if !ok {
		return []string{fmt.Sprintf("%s: expected color string or null", location)}
	}
	if !isHexColor(s) {
		return []string{fmt.Sprintf("%s: invalid color %q", location, s)}
	}
	return nil
}

func validateOptionalFontStyle(location string, value any) []string {
	if value == nil {
		return nil
	}
	s, ok := value.(string)
	if !ok {
		return []string{fmt.Sprintf("%s: expected string or null", location)}
	}
	if s != "normal" && s != "italic" && s != "oblique" {
		return []string{fmt.Sprintf("%s: value = %q, want normal, italic, oblique, or null", location, s)}
	}
	return nil
}

func validateOptionalFontWeight(location string, value any) []string {
	if value == nil {
		return nil
	}
	n, ok := value.(float64)
	if !ok {
		return []string{fmt.Sprintf("%s: expected numeric weight or null", location)}
	}
	if math.Trunc(n) != n || int(n)%100 != 0 || n < 100 || n > 900 {
		return []string{fmt.Sprintf("%s: value = %v, want 100..900 in steps of 100 or null", location, n)}
	}
	return nil
}

func isHexColor(value string) bool {
	if !strings.HasPrefix(value, "#") {
		return false
	}
	hex := strings.TrimPrefix(value, "#")
	if len(hex) != 6 && len(hex) != 8 {
		return false
	}
	for _, r := range hex {
		if !isHexDigit(r) {
			return false
		}
	}
	return true
}

func isHexDigit(r rune) bool {
	return (r >= '0' && r <= '9') || (r >= 'a' && r <= 'f') || (r >= 'A' && r <= 'F')
}

func keySet(m map[string]any) map[string]struct{} {
	out := make(map[string]struct{}, len(m))
	for k := range m {
		out[k] = struct{}{}
	}
	return out
}
