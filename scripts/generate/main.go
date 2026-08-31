package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"zed-themes/scripts/themeutil"
)

const (
	todoValue        = "TODO"
	transparentColor = "#00000000"
	selectionAlpha   = "4D"
	wipSuffix        = " (WIP)"
	hoverBlendAlpha  = 0xBF
)

const (
	backgroundAppearanceBlurred = "blurred"
	backgroundAppearanceOpaque  = "opaque"
	blurModeFlat                = "flat"
)

type Palette = themeutil.Palette
type Meta = themeutil.Meta
type AlphaConfig = themeutil.AlphaConfig

type Config struct {
	TemplatePath   string
	PalettePath    string
	AlphaPath      string
	OutPath        string
	PruneStyle     bool
	ComparePath    string
	WriteStyleKeys string
	WriteAlpha     bool
	WIP            bool
	KeepTODOs      bool
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	cfg := parseFlags()
	if cfg.PalettePath == "" {
		return errors.New("missing --palette")
	}

	palette, template, alphaCfg, err := loadInputs(cfg)
	if err != nil {
		return err
	}

	outPath := resolveOutputPath(cfg)
	style, theme := buildThemeOutput(cfg, palette, template, alphaCfg)

	if cfg.ComparePath != "" {
		if err := compareAndMaybeUpdatePalette(cfg, palette, alphaCfg, style); err != nil {
			return err
		}
	}

	if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
		return fmt.Errorf("create output dir: %w", err)
	}

	if err := writeThemeOutput(outPath, "theme", theme); err != nil {
		return err
	}

	for _, variant := range derivedOutputVariants(palette.Meta, outPath) {
		variantPalette := palette
		variantPalette.Meta = variant.Meta
		_, variantTheme := buildThemeOutput(cfg, variantPalette, template, alphaCfg)
		if err := writeThemeOutput(variant.OutPath, variant.Label, variantTheme); err != nil {
			return err
		}
	}
	return nil
}

func parseFlags() Config {
	cfg := Config{}
	flag.StringVar(&cfg.TemplatePath, "template", "templates/base-style.json", "path to base style template")
	flag.StringVar(&cfg.PalettePath, "palette", "", "path to palette json")
	flag.StringVar(&cfg.AlphaPath, "alpha", "palettes/alpha.json", "path to alpha config")
	flag.StringVar(&cfg.OutPath, "out", "", "output theme json path")
	flag.BoolVar(&cfg.PruneStyle, "prune", true, "drop keys not present in palette style when available")
	flag.StringVar(&cfg.ComparePath, "compare", "", "reference theme json to compare generated style against")
	flag.StringVar(&cfg.WriteStyleKeys, "write-style-keys", "", "comma-separated top-level style keys to copy from reference into palette style")
	flag.BoolVar(&cfg.WriteAlpha, "write-alpha", false, "update palette alpha values to match reference")
	flag.BoolVar(&cfg.WIP, "wip", true, "append WIP suffix to names and filenames")
	flag.BoolVar(&cfg.KeepTODOs, "keep-todos", false, "keep TODO values for debugging")
	flag.Parse()
	return cfg
}

func loadInputs(cfg Config) (Palette, map[string]any, AlphaConfig, error) {
	palette, err := themeutil.ReadJSONFile[Palette](cfg.PalettePath)
	if err != nil {
		return Palette{}, nil, AlphaConfig{}, fmt.Errorf("read palette: %w", err)
	}

	template, err := themeutil.ReadJSONFile[map[string]any](cfg.TemplatePath)
	if err != nil {
		return Palette{}, nil, AlphaConfig{}, fmt.Errorf("read template: %w", err)
	}

	alphaCfg, err := themeutil.ReadJSONFile[AlphaConfig](cfg.AlphaPath)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return Palette{}, nil, AlphaConfig{}, fmt.Errorf("read alpha: %w", err)
	}
	themeutil.MergeAlphaConfig(&alphaCfg, palette.Alpha)

	return palette, template, alphaCfg, nil
}

func resolveOutputPath(cfg Config) string {
	if cfg.OutPath != "" {
		return cfg.OutPath
	}
	base := strings.TrimSuffix(filepath.Base(cfg.PalettePath), filepath.Ext(cfg.PalettePath))
	if cfg.WIP {
		return filepath.Join("generated", "themes", base+".wip.json")
	}
	return filepath.Join("generated", "themes", base+".json")
}

func buildTheme(p Palette, style map[string]any, wip bool) map[string]any {
	name := p.Meta.Name
	themeName := p.Meta.ThemeName
	if wip {
		name = withWIPSuffix(name)
		themeName = withWIPSuffix(themeName)
	}
	return map[string]any{
		"$schema": themeutil.ThemeSchemaURL,
		"name":    name,
		"author":  p.Meta.Author,
		"themes": []any{
			map[string]any{
				"appearance": p.Meta.Appearance,
				"name":       themeName,
				"style":      style,
			},
		},
	}
}

func buildThemeOutput(cfg Config, p Palette, template map[string]any, alpha AlphaConfig) (map[string]any, map[string]any) {
	style := buildStyle(template, p, alpha, cfg.PruneStyle)
	if !cfg.KeepTODOs {
		removeTODOs(style)
	}
	return style, buildTheme(p, style, cfg.WIP)
}

func writeThemeOutput(outPath, label string, theme map[string]any) error {
	if outPath == "" {
		return nil
	}
	if err := themeutil.WriteJSONFile(outPath, theme); err != nil {
		return fmt.Errorf("write %s: %w", label, err)
	}
	return nil
}
