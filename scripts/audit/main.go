package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/tabwriter"

	"zed-themes/scripts/themeutil"
)

type Palette = themeutil.Palette
type AlphaConfig = themeutil.AlphaConfig

type auditRow struct {
	Name         string
	Overrides    int
	Derived      int
	StyleExtras  int
	SyntaxKeys   int
	SyntaxSource string
	Players      int
	PlayerSource string
	AlphaLight   int
	AlphaDark    int
	SyntaxTypos  int
	Issues       string
}

var knownSyntaxTypos = []string{
	"ion.special",
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	paths, err := filepath.Glob("palettes/*.json")
	if err != nil {
		return fmt.Errorf("glob palettes: %w", err)
	}
	if len(paths) == 0 {
		return errors.New("no palettes found")
	}

	var rows []auditRow
	for _, path := range paths {
		if strings.HasSuffix(path, "/alpha.json") || filepath.Base(path) == "alpha.json" {
			continue
		}

		palette, err := themeutil.ReadJSONFile[Palette](path)
		if err != nil {
			return fmt.Errorf("read %s: %w", path, err)
		}

		row := auditRow{
			Name:         strings.TrimSuffix(filepath.Base(path), filepath.Ext(path)),
			Overrides:    len(palette.Overrides),
			Derived:      len(palette.Derived),
			StyleExtras:  countStyleExtras(palette.Style),
			SyntaxKeys:   countNestedMapKeys(palette.Style, "syntax"),
			SyntaxSource: "explicit",
			Players:      countArrayEntries(palette.Style, "players"),
			PlayerSource: "explicit",
			AlphaLight:   len(palette.Alpha.Light),
			AlphaDark:    len(palette.Alpha.Dark),
			SyntaxTypos:  countKnownSyntaxTypos(palette.Style),
		}
		if row.SyntaxKeys == 0 {
			row.SyntaxSource = "derived"
		}
		if row.Players == 0 {
			if len(palette.Accents) > 0 {
				row.Players = len(palette.Accents)
				row.PlayerSource = "derived"
			} else {
				row.PlayerSource = "missing"
			}
		}
		row.Issues = summarizeIssues(row)
		rows = append(rows, row)
	}

	sort.Slice(rows, func(i, j int) bool {
		if rows[i].Overrides != rows[j].Overrides {
			return rows[i].Overrides > rows[j].Overrides
		}
		if rows[i].StyleExtras != rows[j].StyleExtras {
			return rows[i].StyleExtras > rows[j].StyleExtras
		}
		if rows[i].SyntaxKeys != rows[j].SyntaxKeys {
			return rows[i].SyntaxKeys > rows[j].SyntaxKeys
		}
		return rows[i].Name < rows[j].Name
	})

	w := tabwriter.NewWriter(os.Stdout, 0, 4, 2, ' ', 0)
	fmt.Fprintln(w, "palette\toverrides\tderived\tstyle_extras\tsyntax\tsyntax_src\tplayers\tplayer_src\talpha_light\talpha_dark\tissues")
	for _, row := range rows {
		fmt.Fprintf(
			w,
			"%s\t%d\t%d\t%d\t%d\t%s\t%d\t%s\t%d\t%d\t%s\n",
			row.Name,
			row.Overrides,
			row.Derived,
			row.StyleExtras,
			row.SyntaxKeys,
			row.SyntaxSource,
			row.Players,
			row.PlayerSource,
			row.AlphaLight,
			row.AlphaDark,
			row.Issues,
		)
	}
	return w.Flush()
}

func summarizeIssues(row auditRow) string {
	var issues []string
	if row.Overrides >= 10 {
		issues = append(issues, "override-heavy")
	}
	if row.StyleExtras > 0 {
		issues = append(issues, fmt.Sprintf("extra-style:%d", row.StyleExtras))
	}
	if row.SyntaxTypos > 0 {
		issues = append(issues, fmt.Sprintf("syntax-typos:%d", row.SyntaxTypos))
	}
	if row.PlayerSource == "missing" {
		issues = append(issues, "missing-players")
	}
	if len(issues) == 0 {
		return "-"
	}
	return strings.Join(issues, ",")
}

func countStyleExtras(style map[string]any) int {
	if len(style) == 0 {
		return 0
	}
	count := 0
	for key := range style {
		if key == "syntax" || key == "players" {
			continue
		}
		count++
	}
	return count
}

func countNestedMapKeys(style map[string]any, key string) int {
	if len(style) == 0 {
		return 0
	}
	nested, ok := style[key].(map[string]any)
	if !ok {
		return 0
	}
	return len(nested)
}

func countArrayEntries(style map[string]any, key string) int {
	if len(style) == 0 {
		return 0
	}
	values, ok := style[key].([]any)
	if !ok {
		return 0
	}
	return len(values)
}

func countKnownSyntaxTypos(style map[string]any) int {
	if len(style) == 0 {
		return 0
	}
	syntax, ok := style["syntax"].(map[string]any)
	if !ok {
		return 0
	}

	count := 0
	for _, key := range knownSyntaxTypos {
		if _, ok := syntax[key]; ok {
			count++
		}
	}
	return count
}
