package main

import (
	"path/filepath"
	"strings"
)

type outputVariant struct {
	Label   string
	Meta    Meta
	OutPath string
}

func withWIPSuffix(name string) string {
	if name == "" {
		return name
	}
	if strings.HasSuffix(name, wipSuffix) {
		return name
	}
	return name + wipSuffix
}

func shouldGenerateBlurVariant(meta Meta) bool {
	if !strings.EqualFold(meta.BackgroundAppearance, backgroundAppearanceBlurred) {
		return false
	}
	return !strings.EqualFold(meta.BlurMode, blurModeFlat)
}

func derivedOutputVariants(meta Meta, outPath string) []outputVariant {
	var variants []outputVariant
	if shouldGenerateBlurVariant(meta) {
		variants = append(variants, outputVariant{
			Label:   "blur theme",
			Meta:    blurVariantMeta(meta),
			OutPath: blurOutputPath(outPath),
		})
	}
	if shouldGenerateFlatVariant(meta) {
		variants = append(variants, outputVariant{
			Label:   "flat theme",
			Meta:    flatVariantMeta(meta),
			OutPath: flatOutputPath(outPath),
		})
	}
	return variants
}

func blurVariantMeta(meta Meta) Meta {
	out := meta
	out.BlurMode = blurModeFlat
	out.Name = blurName(meta.Name)
	out.ThemeName = blurThemeName(meta.ThemeName)
	return out
}

func blurName(name string) string {
	return variantName(name, "Blur")
}

func blurThemeName(name string) string {
	return variantThemeName(name, "Blur")
}

func blurOutputPath(outPath string) string {
	return variantOutputPath(outPath, "-blur")
}

func shouldGenerateFlatVariant(meta Meta) bool {
	if !strings.EqualFold(meta.BackgroundAppearance, backgroundAppearanceBlurred) {
		return false
	}
	return !strings.EqualFold(meta.BlurMode, blurModeFlat)
}

func flatVariantMeta(meta Meta) Meta {
	out := meta
	out.BackgroundAppearance = backgroundAppearanceOpaque
	out.BlurMode = ""
	out.Name = flatName(meta.Name)
	out.ThemeName = flatThemeName(meta.ThemeName)
	return out
}

func flatName(name string) string {
	return variantName(name, "Flat")
}

func flatThemeName(name string) string {
	return variantThemeName(name, "Flat")
}

func flatOutputPath(outPath string) string {
	return variantOutputPath(outPath, "-flat")
}

func variantOutputPath(outPath string, variantSuffix string) string {
	if outPath == "" {
		return ""
	}
	dir := filepath.Dir(outPath)
	base := strings.TrimSuffix(filepath.Base(outPath), filepath.Ext(outPath))
	wip := false
	if strings.HasSuffix(base, ".wip") {
		wip = true
		base = strings.TrimSuffix(base, ".wip")
	}
	ext := filepath.Ext(outPath)
	if strings.HasSuffix(base, variantSuffix) {
		return ""
	}
	if before, ok := strings.CutSuffix(base, "-hybrid"); ok {
		base = before + variantSuffix
	} else {
		base = base + variantSuffix
	}
	if wip {
		base = base + ".wip"
	}
	return filepath.Join(dir, base+ext)
}

func variantName(name string, variant string) string {
	if name == "" {
		return name
	}
	target := "(" + variant + ")"
	if strings.Contains(name, target) {
		return name
	}
	for _, existing := range []string{"Hybrid", "Blur", "Flat"} {
		token := "(" + existing + ")"
		if strings.Contains(name, token) {
			return strings.ReplaceAll(name, token, target)
		}
		if strings.Contains(name, existing) {
			return strings.ReplaceAll(name, existing, variant)
		}
	}
	return name + " " + variant
}

func variantThemeName(name string, variant string) string {
	if name == "" {
		return name
	}
	target := "(" + variant + ")"
	if strings.Contains(name, target) {
		return name
	}
	for _, existing := range []string{"Hybrid", "Blur", "Flat"} {
		token := "(" + existing + ")"
		if strings.Contains(name, token) {
			return strings.ReplaceAll(name, token, target)
		}
		if strings.Contains(name, existing) {
			return strings.ReplaceAll(name, existing, variant)
		}
	}
	return name + " " + target
}

func hasOpaqueSemanticBackgrounds(meta Meta) bool {
	return strings.EqualFold(meta.BackgroundAppearance, backgroundAppearanceBlurred) ||
		strings.EqualFold(meta.BackgroundAppearance, backgroundAppearanceOpaque)
}
