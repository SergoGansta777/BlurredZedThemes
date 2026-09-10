<p align="center">
  <img alt="Blurred Zed Themes" src="https://img.shields.io/badge/Blurred%20Zed%20Themes-curated%20hybrid%20collection-111827?style=for-the-badge" />
</p>

<p align="center">
  A curated collection of blurred, hybrid, and flat themes for the Zed editor, tuned for modern UI surfaces,
  clean contrast, and consistent syntax colors.
</p>

<p align="center">
  <img alt="Theme families" src="https://img.shields.io/badge/families-25-4C9AFF?style=flat-square" />
  <img alt="Published variants" src="https://img.shields.io/badge/variants-73-0A84FF?style=flat-square" />
  <img alt="Last commit" src="https://img.shields.io/github/last-commit/SergoGansta777/BlurredZedThemes?style=flat-square" />
  <img alt="Status" src="https://img.shields.io/badge/status-maintained-30D158?style=flat-square" />
</p>

## Overview

These themes are built around Zed’s blurred UI, with optional flat variants for fully opaque window backgrounds. The editor stays sharp, the chrome stays soft where blur is enabled, and the whole layout keeps good contrast without feeling noisy.

- Stable editor backgrounds with transparent UI layers around them.
- Balanced alpha values for panels, overlays, tabs, and status bars.
- Flat variants with one consistent opaque surface background across editor, panels, tabs, and toolbars.
- Consistent syntax mapping across all themes and variants.
- Three variants per theme: Blur, Hybrid, and Flat, except Token Dark (opaque only).

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
| Evergarden     | Winter:<br><img width="320" alt="Evergarden Winter (Hybrid)" src="https://github.com/user-attachments/assets/a551c81f-73b1-4aec-a0f8-476ff8aefbac" /><br>Spring: TODO<br>Summer: TODO<br>Fall: TODO | https://github.com/evergardentheme/nvim                       |
| JetBrains      | Dark: TODO<br>Light: TODO                                                                                                                                                                           | https://github.com/artemevsevev/zed-theme-jetbrains           |
| Xcode          | Dark: TODO<br>Light: TODO                                                                                                                                                                           | https://github.com/skarline/zed-xcode-themes                  |
| Kanagawa       | Dragon: TODO<br>Paper: TODO                                                                                                                                                                         | https://github.com/rebelot/kanagawa.nvim                      |
| Cosmos         | <img width="320" alt="Cosmos (Hybrid)" src="https://github.com/user-attachments/assets/195383d5-5f5d-449d-af62-d9a1d0f79ef3" />                                                                     | https://github.com/nauvalazhar/cosmos                         |
| Darkearth      | <img width="320" alt="Darkearth (Hybrid)" src="https://github.com/user-attachments/assets/5ae80649-35a1-44ed-be45-e3abeb62f6ec" />                                                                  | https://github.com/ptdewey/darkearth-nvim                     |
| Ember          | Dark: TODO<br>Soft: TODO<br>Light: TODO<br>Lighter: TODO                                                                                                                                             | https://github.com/ember-theme/nvim                           |
| Everforest     | TODO                                                                                                                                                                                                | https://github.com/neanias/everforest-nvim                    |
| Ayu            | TODO                                                                                                                                                                                                | https://github.com/zed-industries/zed/tree/main/assets/themes/ayu |
| Lunar          | <img width="320" alt="Lunar (Hybrid)" src="https://github.com/user-attachments/assets/a0e76368-8ffb-4d9b-ad9d-99bccc3884d3" />                                                                      | https://github.com/comfysage/lunarfrost                       |
| Miasma Fog     | <img width="320" alt="Miasma Fog (Hybrid)" src="https://github.com/user-attachments/assets/c0308e82-e801-418b-9f1b-c2f2692031d0" />                                                                 | https://github.com/xero/miasma.nvim                           |
| Nordic         | <img width="320" alt="Nordic (Hybrid)" src="https://github.com/user-attachments/assets/be112f4e-6176-411a-92bf-d7659a2838d7" />                                                                     | https://github.com/AlexvZyl/nordic.nvim                       |
| Oldworld       | TODO                                                                                                                                                                                                | https://github.com/dgox16/oldworld.nvim                       |
| Rosé Pine Dawn | <img width="320" alt="Rosé Pine Dawn (Hybrid)" src="https://github.com/user-attachments/assets/1113c3bd-892e-48bf-8200-1ed5105dfbf7" />                                                             | https://github.com/rose-pine/zed                              |
| Vesper         | TODO                                                                                                                                                                                                | https://github.com/raunofreiberg/vesper                       |
| Token          | Dark: TODO | https://github.com/ThorstenRhau/token |
| Rusty          | Dark: local draft, licensing unresolved | https://github.com/armannikoyan/rusty |

## Customization

### Token Dark

Adapted from [classic Token](https://github.com/ThorstenRhau/token/tree/86e66d9ab7c74d53e7ac56a02b7f8bee56196cda), checked September 10, 2026. Only the original dark palette is included, not Ultra, Meridian, Flint, Temper, or light variants.

`palettes/token-dark.json` generates `themes/token-dark.json`, selectable as **Token Dark**. It preserves the core upstream syntax colors and typography, the 16 ANSI colors, neutral selection, and diagnostic/line-diff backgrounds. Intentional Zed adaptations:

- Uniform opaque editor/chrome surfaces and selected tabs matching the neutral selection.
- Darker, opaque search and word-diff highlights to keep syntax readable without Neovim's foreground inversion.
- Distinct upstream Vim mode hues: neutral normal, green insert, peach visual, and red replace.
- Explicit dim ANSI colors, with colored entries blended 20% toward the editor background and black darkened separately.
- Shared generator defaults for Zed-only captures and UI groups, with palette-local corrections for semantic comments and macros.

To install only Token Dark, copy `themes/token-dark.json` into `~/.config/zed/themes/` and select **Token Dark**. To regenerate it without publishing the rest of the collection:

```bash
go run ./scripts/generate --palette palettes/token-dark.json --wip=false --out themes/token-dark.json
go run scripts/format/main.go themes/token-dark.json
```

For future upstream reviews, compare the pinned revision's `lua/token/palette.lua`, `groups/treesitter.lua`, `groups/syntax.lua`, `typography.lua`, `terminal.lua`, and `lualine.lua`. Preserve the Zed adaptations above rather than replacing the generated theme wholesale. Retain `licenses/token.txt` when redistributing Token-derived data.

### Rusty (Local Draft)

Adapted from [Rusty at f310310](https://github.com/armannikoyan/rusty/tree/f310310991ecd50ca10745f8960cb5b8bd2ef208), checked September 10, 2026. `drafts/rusty/palette.json` generates the single opaque **Rusty** theme in `drafts/rusty/themes/rusty.json`. Published collection counts above exclude this draft. Keeping both files outside `palettes/` and `themes/` excludes them from normal generation, publishing, extension discovery, and wildcard installation.

The port follows the source palette, not the README's optional transparency example. It preserves the charcoal background, cool-gray selection, italic comments, purple statements, orange types/constants, green strings, aqua Tree-sitter/LSP functions, and neutral LSP variables/enum members. Where legacy and LSP mappings differ, the port favors the explicit modern mappings.

Intentional Zed adaptations:

- Selected tabs and menus use the original selection color rather than Neovim's reversed foreground/background.
- The current-line color uses upstream's `#282a2e`, but visibility follows Zed's `current_line_highlight` setting. Set it to `none` to match Rusty's default, or `line` to enable it.
- Zed diagnostic warnings use the original yellow, while errors retain the original red. Upstream's red `WarningMsg` is a Vim message mapping, not an explicit LSP diagnostic palette.
- Search highlights use dark yellow-tinted backgrounds instead of bright yellow with inverted text.
- Diff hunks use a darker neutral background than upstream's `#494e56`; word additions/deletions use restrained green/red backgrounds because Zed retains syntax foregrounds.
- Vim mode indicators retain upstream Lualine's blue, green, purple, and red hues.
- Terminal ANSI colors are mapped from the named palette hues. Upstream claims terminal support but defines no ANSI assignments at this revision. Bright colors retain the original hues (bright black uses the comment color); dim colors are explicitly subdued.
- Zed-only UI and syntax groups use the existing generator, with explicit overrides where its defaults differ from Rusty. No upstream Lua code is copied.

The original red is retained despite measuring approximately 4.46:1 against the editor background, with lower contrast on selections and highlights. This is a fidelity-first port, not a claim of universal accessibility compliance.

```bash
go run ./scripts/generate --palette drafts/rusty/palette.json --wip=false --out drafts/rusty/themes/rusty.json
go run scripts/format/main.go drafts/rusty/palette.json drafts/rusty/themes/rusty.json
go run ./scripts/validate --palettes-dir drafts/rusty --themes-dir drafts/rusty/themes
go run ./scripts/generate --palette drafts/rusty/palette.json --compare drafts/rusty/themes/rusty.json
```

For local preview, copy `drafts/rusty/themes/rusty.json` into `~/.config/zed/themes/` and select **Rusty**. Review `lua/rusty/colors.lua`, `lua/rusty/init.lua`, and `lua/rusty/plugins/lualine.lua` at the pinned revision when updating. Draft checks above are separate from `task check`, which covers the published collection.

**Publication hold:** no license file or license grant was found in the upstream source or README. This draft is not represented as Apache-licensed. It may be kept in local history, but upstream licensing remains unresolved; do not push or redistribute this draft until clarified. Moving it outside the extension directories does not exclude it from a Git push. Attribution alone is not a license grant.

### Shared Settings

- Global alpha presets live in `palettes/alpha.json`.
- Per-theme values live in `palettes/<theme>.json`.
- Use palette `derived` entries for exact generated colors such as native search, document highlights, Vim, scrollbars, or other style keys.
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
- Per-theme `alpha` values can be added when needed and are merged over `palettes/alpha.json`.
- Use `derived` as the single source for exact generated style-key values.
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

Token-derived palette and theme data retain Thorsten Rhau's [BSD 3-Clause license](licenses/token.txt).

The local Rusty draft is excluded from the Apache license statement; upstream licensing remains unresolved (see above).
