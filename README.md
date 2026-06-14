<p align="center">
  <img alt="Blurred Zed Themes" src="https://img.shields.io/badge/Blurred%20Zed%20Themes-curated%20hybrid%20collection-111827?style=for-the-badge" />
</p>

<p align="center">
  A curated collection of blurred, hybrid, and flat themes for the Zed editor, tuned for modern UI surfaces,
  clean contrast, and consistent syntax colors.
</p>

<p align="center">
  <img alt="Theme families" src="https://img.shields.io/badge/families-23-4C9AFF?style=flat-square" />
  <img alt="Published variants" src="https://img.shields.io/badge/variants-69-0A84FF?style=flat-square" />
  <img alt="Last commit" src="https://img.shields.io/github/last-commit/SergoGansta777/BlurredZedThemes?style=flat-square" />
  <img alt="Status" src="https://img.shields.io/badge/status-maintained-30D158?style=flat-square" />
</p>

## Overview

These themes are built around Zed’s blurred UI, with optional flat variants for fully opaque window backgrounds. The editor stays sharp, the chrome stays soft where blur is enabled, and the whole layout keeps good contrast without feeling noisy.

- Stable editor backgrounds with transparent UI layers around them.
- Balanced alpha values for panels, overlays, tabs, and status bars.
- Flat variants with one consistent opaque surface background across editor, panels, tabs, and toolbars.
- Consistent syntax mapping across all themes and variants.
- Three variants per theme: Blur, Hybrid, and Flat.

## Install

As local themes:

```bash
mkdir -p ~/.config/zed/themes
cp themes/*.json ~/.config/zed/themes/
```

Then restart Zed (or reload themes) and select a theme in Settings → Theme.

As a Zed dev extension, install this repository directory via `zed: install dev extension`. Zed expects an `extension.toml` manifest at the repository root and theme files under `themes/`.

## Theme gallery

Grouped by theme family. Previews are added as they become available.

| Theme group    | Preview                                                                                                                                                                                             | Source / inspiration                                          |
| -------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------- |
| Evergarden     | Winter:<br><img width="320" alt="Evergarden Winter (Hybrid)" src="https://github.com/user-attachments/assets/a551c81f-73b1-4aec-a0f8-476ff8aefbac" /><br>Spring: TODO<br>Summer: TODO<br>Fall: TODO | https://github.com/everviolet/nvim                            |
| JetBrains      | Dark: TODO<br>Light: TODO                                                                                                                                                                           | https://github.com/artemevsevev/zed-theme-jetbrains           |
| Xcode          | Dark: TODO<br>Light: TODO                                                                                                                                                                           | https://github.com/skarline/zed-xcode-themes                  |
| Kanagawa       | Dragon: TODO<br>Paper: TODO                                                                                                                                                                         | https://github.com/rebelot/kanagawa.nvim                      |
| Cosmos         | <img width="320" alt="Cosmos (Hybrid)" src="https://github.com/user-attachments/assets/195383d5-5f5d-449d-af62-d9a1d0f79ef3" />                                                                     | https://github.com/nauvalazhar/cosmos                         |
| Darkearth      | <img width="320" alt="Darkearth (Hybrid)" src="https://github.com/user-attachments/assets/5ae80649-35a1-44ed-be45-e3abeb62f6ec" />                                                                  | https://github.com/ptdewey/darkearth-nvim                     |
| Ember          | Dark: TODO<br>Soft: TODO<br>Light: TODO                                                                                                                                                             | https://github.com/ember-theme/nvim                           |
| Everforest     | TODO                                                                                                                                                                                                | https://github.com/neanias/everforest-nvim                    |
| Ayu            | TODO                                                                                                                                                                                                | https://github.com/zed-industries/zed/tree/main/assets/themes/ayu |
| Lunar          | <img width="320" alt="Lunar (Hybrid)" src="https://github.com/user-attachments/assets/a0e76368-8ffb-4d9b-ad9d-99bccc3884d3" />                                                                      | https://github.com/comfysage/lunarfrost                       |
| Miasma Fog     | <img width="320" alt="Miasma Fog (Hybrid)" src="https://github.com/user-attachments/assets/c0308e82-e801-418b-9f1b-c2f2692031d0" />                                                                 | https://github.com/xero/miasma.nvim                           |
| Nordic         | <img width="320" alt="Nordic (Hybrid)" src="https://github.com/user-attachments/assets/be112f4e-6176-411a-92bf-d7659a2838d7" />                                                                     | https://github.com/AlexvZyl/nordic.nvim                       |
| Oldworld       | TODO                                                                                                                                                                                                | https://github.com/dgox16/oldworld.nvim                       |
| Rosé Pine Dawn | <img width="320" alt="Rosé Pine Dawn (Hybrid)" src="https://github.com/user-attachments/assets/1113c3bd-892e-48bf-8200-1ed5105dfbf7" />                                                             | https://github.com/rose-pine/zed                              |
| Vesper         | TODO                                                                                                                                                                                                | https://github.com/raunofreiberg/vesper                       |

## Customization

- Global alpha presets live in `palettes/alpha.json`.
- Per-theme overrides live in `palettes/<theme>.json`.
- Use palette `derived` entries for generator-only base colors such as native search, document highlight, Vim, or scrollbar colors before falling back to raw `overrides`.
- The generator backfills modern syntax captures and practical Zed-only keys such as word diff highlights and hover line numbers.
- Regenerate theme files via Taskfile (see below).

## Upstream references

Use these when auditing or modernizing the collection:

- Theme Builder: https://zed.dev/theme-builder
- Theme docs: https://zed.dev/docs/themes
- Theme schema: https://zed.dev/schema/themes/v0.2.0.json
- Bundled Zed themes: https://github.com/zed-industries/zed/tree/main/assets/themes

Notes:

- The published schema is useful, but it is not the whole story. Zed’s bundled themes currently use a few practical keys beyond the schema, so upstream theme files are the best compatibility reference.
- This repo treats the bundled themes and the Theme Builder as the authoritative guide for new groups and real-world key usage.
- Current Zed-native sync targets are `Ayu Mirage`, `JetBrains Dark`, `JetBrains Light`, `Rosé Pine Dawn`, `Xcode Default Dark`, and `Xcode Default Light`.

## Taskfile workflow

All common workflows are wrapped in `Taskfile.yml`:

```bash
task gen-all
task publish
task verify
task validate
task audit
task check
```

Notes:

- Palettes define roles/semantic/derived/accents/terminal, with optional `style` for `syntax` and `players`.
- `alpha` overrides can be added per theme when needed (merged over `palettes/alpha.json`).
- `overrides` are treated as derived data and can be regenerated from a reference theme.
- For Zed-native upstream themes, prefer `go run ./scripts/generate --palette palettes/<theme>.json --compare <reference.json> --write-style-keys syntax,players` so syntax and player colors can be refreshed without copying the whole upstream style surface.
- The generator fills missing fields with `TODO` placeholders and applies safe defaults.
- `task validate` checks published theme family shape, style keys, duplicate theme names, color syntax, players, and syntax highlight entries.
- Published/reference themes live in `themes/`.

Before publishing or installing as an extension, run:

```bash
task check
```

## Recommended settings

These settings match the screenshots and keep the layout clean. Themes are designed primarily for macOS but should work on other platforms that support blur.

```json
{
  "current_line_highlight": "none", // By your preference
  "project_panel": {
    "sticky_scroll": false // Not fully supported yet
  },
  "sticky_scroll": {
    "enabled": true // By your preference
  }
}
```

## Contributing

- Open issues for visual inconsistencies, contrast/accessibility concerns, or missing mappings.
- PRs are welcome for new variants, improved syntax coverage, or closer alignment with upstream palettes.

## License

Licensed under the Apache License, Version 2.0. See `LICENSE`.
