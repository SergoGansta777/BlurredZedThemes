<p align="center">
  <img alt="Blurred Zed Themes" src="https://img.shields.io/badge/Blurred%20Zed%20Themes-curated%20hybrid%20collection-111827?style=for-the-badge" />
</p>

<p align="center">
  A curated collection of blurred hybrid themes for the Zed editor, tuned for modern UI surfaces,
  clean contrast, and consistent syntax colors.
</p>

<p align="center">
  <img alt="Themes" src="https://img.shields.io/badge/themes-17-4C9AFF?style=flat-square" />
  <img alt="Last commit" src="https://img.shields.io/github/last-commit/SergoGansta777/BlurredZedThemes?style=flat-square" />
  <img alt="Status" src="https://img.shields.io/badge/status-maintained-30D158?style=flat-square" />
</p>

## Overview

These themes are built around Zed’s blurred UI. The editor stays sharp, the chrome stays soft, and the whole layout keeps good contrast without feeling noisy.

- Stable editor backgrounds with transparent UI layers around them.
- Balanced alpha values for panels, overlays, tabs, and status bars.
- Consistent syntax mapping across all themes and variants.
- Two variants per theme: Blur and Hybrid.

## Install

```bash
mkdir -p ~/.config/zed/themes
cp themes/*.json ~/.config/zed/themes/
```

Then restart Zed (or reload themes) and select a theme in Settings → Theme.

## Theme gallery

Grouped by theme family. Previews are added as they become available.

| Theme group    | Preview                                                                                                                                                                                             | Source / inspiration                                          |
| -------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------- |
| Evergarden     | Winter:<br><img width="320" alt="Evergarden Winter (Hybrid)" src="https://github.com/user-attachments/assets/a551c81f-73b1-4aec-a0f8-476ff8aefbac" /><br>Spring: TODO<br>Summer: TODO<br>Fall: TODO | https://github.com/everviolet/nvim                            |
| JetBrains      | Dark: TODO<br>Light: TODO                                                                                                                                                                           | https://github.com/zed-industries/zed/tree/main/assets/themes |
| Xcode          | Dark: TODO<br>Light: TODO                                                                                                                                                                           | Apple Xcode default colors + local `xcode theme` reference    |
| Kanagawa       | Dragon: TODO<br>Paper: TODO                                                                                                                                                                         | https://github.com/rebelot/kanagawa.nvim                      |
| Cosmos         | <img width="320" alt="Cosmos (Hybrid)" src="https://github.com/user-attachments/assets/195383d5-5f5d-449d-af62-d9a1d0f79ef3" />                                                                     | https://github.com/nauvalazhar/cosmos                         |
| Darkearth      | <img width="320" alt="Darkearth (Hybrid)" src="https://github.com/user-attachments/assets/5ae80649-35a1-44ed-be45-e3abeb62f6ec" />                                                                  | https://github.com/ptdewey/darkearth-nvim                     |
| Everforest     | TODO                                                                                                                                                                                                | https://github.com/sainnhe/everforest                         |
| Ayu            | TODO                                                                                                                                                                                                | https://github.com/ayu-theme/ayu-colors                       |
| Lunar          | <img width="320" alt="Lunar (Hybrid)" src="https://github.com/user-attachments/assets/a0e76368-8ffb-4d9b-ad9d-99bccc3884d3" />                                                                      | https://github.com/LunarVim/Colorschemes                      |
| Miasma Fog     | <img width="320" alt="Miasma Fog (Hybrid)" src="https://github.com/user-attachments/assets/c0308e82-e801-418b-9f1b-c2f2692031d0" />                                                                 | https://github.com/xero/miasma.nvim                           |
| Nordic         | <img width="320" alt="Nordic (Hybrid)" src="https://github.com/user-attachments/assets/be112f4e-6176-411a-92bf-d7659a2838d7" />                                                                     | https://github.com/AlexvZyl/nordic.nvim                       |
| Oldworld       | TODO                                                                                                                                                                                                | https://github.com/nyoom-engineering/oldworld.nvim            |
| Rosé Pine Dawn | <img width="320" alt="Rosé Pine Dawn (Hybrid)" src="https://github.com/user-attachments/assets/1113c3bd-892e-48bf-8200-1ed5105dfbf7" />                                                             | https://github.com/rose-pine/zed                              |
| Vesper         | TODO                                                                                                                                                                                                | https://github.com/raunofreiberg/vesper                       |

## Customization

- Global alpha presets live in `palettes/alpha.json`.
- Per-theme overrides live in `palettes/<theme>.json`.
- Regenerate theme files via Taskfile (see below).

## Taskfile workflow

All common workflows are wrapped in `Taskfile.yml`:

```bash
task gen-all
task publish
```

Notes:

- Palettes define roles/semantic/accents/terminal, with optional `style` for `syntax` and `players`.
- `alpha` overrides can be added per theme when needed (merged over `palettes/alpha.json`).
- `overrides` are treated as derived data and can be regenerated from a reference theme.
- The generator fills missing fields with `TODO` placeholders and applies safe defaults.
- Published/reference themes live in `themes/`.

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
