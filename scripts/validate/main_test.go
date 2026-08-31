package main

import (
	"os"
	"path/filepath"
	"testing"

	"zed-themes/scripts/themeutil"
)

func TestIsHexColor(t *testing.T) {
	valid := []string{"#000000", "#ABCDEF", "#abcdef80"}
	for _, value := range valid {
		if !isHexColor(value) {
			t.Fatalf("isHexColor(%q) = false, want true", value)
		}
	}

	invalid := []string{"000000", "#12345", "#123456789", "#GGGGGG"}
	for _, value := range invalid {
		if isHexColor(value) {
			t.Fatalf("isHexColor(%q) = true, want false", value)
		}
	}
}

func TestValidateStyleRejectsUnknownStyleKeys(t *testing.T) {
	issues := validateStyle("theme", map[string]any{
		"unknown.key": "#000000",
	}, map[string]struct{}{
		"background": {},
	})

	if len(issues) != 1 {
		t.Fatalf("issues = %v, want one unknown-key issue", issues)
	}
}

func TestValidateSyntaxFields(t *testing.T) {
	issues := validateSyntax("theme.style.syntax", map[string]any{
		"keyword": map[string]any{
			"color":       "#FF00AA",
			"font_style":  "italic",
			"font_weight": float64(700),
		},
	})
	if len(issues) != 0 {
		t.Fatalf("issues = %v, want none", issues)
	}

	issues = validateSyntax("theme.style.syntax", map[string]any{
		"keyword": map[string]any{
			"font_style":  "bold",
			"font_weight": float64(750),
		},
	})
	if len(issues) != 2 {
		t.Fatalf("issues = %v, want two syntax issues", issues)
	}
}

func TestValidatePaletteRejectsUnknownDerivedKeys(t *testing.T) {
	allowed := derivedKeySet(map[string]struct{}{
		"editor.background": {},
	})

	issues := validatePalette("palette.json", themeutil.Palette{
		Derived: map[string]string{
			"editor.background": "#000000",
			"unknown_key":       "#FFFFFF",
		},
	}, allowed)

	if len(issues) != 1 {
		t.Fatalf("issues = %v, want one unknown-derived issue", issues)
	}
}

func TestDerivedKeySetAllowsAlphaAndVimKeys(t *testing.T) {
	allowed := derivedKeySet(map[string]struct{}{})

	for _, key := range []string{"search_match", "vim_replace"} {
		if _, ok := allowed[key]; !ok {
			t.Fatalf("derivedKeySet missing %q", key)
		}
	}
}

func TestValidateExtensionManifest(t *testing.T) {
	path := filepath.Join(t.TempDir(), "extension.toml")
	if err := os.WriteFile(path, []byte(`
authors = ["Test <test@example.com>"]
description = "Test"
id = "test"
name = "Test"
repository = "https://example.com/test"
schema_version = 1
version = "0.1.0"
`), 0o644); err != nil {
		t.Fatal(err)
	}

	if issues := validateExtensionManifest(path); len(issues) != 0 {
		t.Fatalf("issues = %v, want none", issues)
	}

	if err := os.WriteFile(path, []byte(`id = "test"`), 0o644); err != nil {
		t.Fatal(err)
	}
	if issues := validateExtensionManifest(path); len(issues) == 0 {
		t.Fatal("expected missing-field issues")
	}
}
