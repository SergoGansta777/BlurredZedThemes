package main

import (
	"fmt"
	"strconv"
	"strings"

	"zed-themes/scripts/themeutil"
)

func solidColorValue(style map[string]any, key string, surfaceBg string) string {
	color, ok := style[key].(string)
	if !ok || color == "" {
		return ""
	}
	return flattenColor(color, surfaceBg)
}

func flattenColor(color string, surfaceBg string) string {
	fg, fgAlpha, ok := parseHexColor(color)
	if !ok {
		return themeutil.StripAlpha(color)
	}
	bg, _, ok := parseHexColor(surfaceBg)
	if !ok {
		return themeutil.StripAlpha(color)
	}
	if fgAlpha == 0xFF {
		return formatHexColor(fg)
	}
	if fgAlpha == 0x00 {
		return formatHexColor(bg)
	}
	return formatHexColor([3]uint8{
		blendChannel(fg[0], bg[0], fgAlpha),
		blendChannel(fg[1], bg[1], fgAlpha),
		blendChannel(fg[2], bg[2], fgAlpha),
	})
}

func parseHexColor(color string) ([3]uint8, uint8, bool) {
	h := strings.TrimPrefix(color, "#")
	if len(h) != 6 && len(h) != 8 {
		return [3]uint8{}, 0, false
	}
	r, err := strconv.ParseUint(h[0:2], 16, 8)
	if err != nil {
		return [3]uint8{}, 0, false
	}
	g, err := strconv.ParseUint(h[2:4], 16, 8)
	if err != nil {
		return [3]uint8{}, 0, false
	}
	b, err := strconv.ParseUint(h[4:6], 16, 8)
	if err != nil {
		return [3]uint8{}, 0, false
	}
	alpha := uint64(0xFF)
	if len(h) == 8 {
		alpha, err = strconv.ParseUint(h[6:8], 16, 8)
		if err != nil {
			return [3]uint8{}, 0, false
		}
	}
	return [3]uint8{uint8(r), uint8(g), uint8(b)}, uint8(alpha), true
}

func formatHexColor(rgb [3]uint8) string {
	return fmt.Sprintf("#%02X%02X%02X", rgb[0], rgb[1], rgb[2])
}

func blendChannel(fg uint8, bg uint8, alpha uint8) uint8 {
	return uint8((uint16(fg)*uint16(alpha) + uint16(bg)*(0xFF-uint16(alpha)) + 127) / 0xFF)
}

func mixOpaqueColors(fg string, bg string, alpha uint8) string {
	fgRGB, _, ok := parseHexColor(themeutil.StripAlpha(fg))
	if !ok {
		return themeutil.StripAlpha(fg)
	}
	bgRGB, _, ok := parseHexColor(themeutil.StripAlpha(bg))
	if !ok {
		return themeutil.StripAlpha(fg)
	}
	return formatHexColor([3]uint8{
		blendChannel(fgRGB[0], bgRGB[0], alpha),
		blendChannel(fgRGB[1], bgRGB[1], alpha),
		blendChannel(fgRGB[2], bgRGB[2], alpha),
	})
}
