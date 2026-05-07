package main

import (
	"fmt"
	"strconv"
	"strings"
)

//go:generate go run -tags generate generate_styles.go styles_decl.go

type rgbColor struct {
	R int
	G int
	B int
}

type colorSpec struct {
	hex string
	rgb *rgbColor
}

type tagStyle struct {
	Names []string
	FG    *colorSpec
	BG    colorSpec
}

func RGB(r, g, b int) colorSpec {
	return colorSpec{rgb: &rgbColor{R: r, G: g, B: b}}
}

func Hex(hex string) colorSpec {
	return colorSpec{hex: hex}
}

func Explicit(spec colorSpec) *colorSpec {
	return &spec
}

var tagStyles = []tagStyle{
	{Names: []string{"red"}, FG: Explicit(RGB(255, 255, 255)), BG: RGB(255, 0, 0)},
	{Names: []string{"blue"}, FG: Explicit(RGB(255, 255, 255)), BG: RGB(0, 0, 255)},
	// CSS green is #008000. Keep the historical ntagcolor shade active.
	{Names: []string{"green"}, FG: Explicit(RGB(0, 0, 0)), BG: RGB(0, 255, 0)},
	// CSS orange is #ffa500. Keep the historical ntagcolor shade active.
	{Names: []string{"orange"}, FG: Explicit(RGB(255, 255, 255)), BG: RGB(255, 120, 0)},
	{Names: []string{"yellow"}, FG: Explicit(RGB(0, 0, 0)), BG: RGB(255, 255, 0)},
	// CSS purple is #800080. Keep the historical ntagcolor shade active.
	{Names: []string{"purple"}, FG: Explicit(RGB(255, 255, 255)), BG: RGB(100, 10, 255)},
	// CSS gray/grey is #808080. Keep the historical ntagcolor shade active.
	{Names: []string{"gray", "grey"}, FG: Explicit(RGB(255, 255, 255)), BG: RGB(100, 100, 100)},
	{Names: []string{"black"}, FG: Explicit(RGB(255, 255, 255)), BG: RGB(0, 0, 0)},
	{Names: []string{"white"}, FG: Explicit(RGB(0, 0, 0)), BG: RGB(255, 255, 255)},
	{Names: []string{"aqua"}, FG: Explicit(RGB(0, 0, 0)), BG: RGB(0, 255, 255)},
	{Names: []string{"teal"}, FG: Explicit(RGB(255, 255, 255)), BG: RGB(0, 128, 128)},
	{Names: []string{"emerald"}, FG: Explicit(RGB(0, 0, 0)), BG: RGB(80, 200, 120)},

	{Names: []string{"aliceblue"}, BG: Hex("#f0f8ff")},
	{Names: []string{"antiquewhite"}, BG: Hex("#faebd7")},
	{Names: []string{"aquamarine"}, BG: Hex("#7fffd4")},
	{Names: []string{"azure"}, BG: Hex("#f0ffff")},
	{Names: []string{"beige"}, BG: Hex("#f5f5dc")},
	{Names: []string{"bisque"}, BG: Hex("#ffe4c4")},
	{Names: []string{"blanchedalmond"}, BG: Hex("#ffebcd")},
	{Names: []string{"blueviolet"}, BG: Hex("#8a2be2")},
	{Names: []string{"brown"}, BG: Hex("#a52a2a")},
	{Names: []string{"burlywood"}, BG: Hex("#deb887")},
	{Names: []string{"cadetblue"}, BG: Hex("#5f9ea0")},
	{Names: []string{"chartreuse"}, BG: Hex("#7fff00")},
	{Names: []string{"chocolate"}, BG: Hex("#d2691e")},
	{Names: []string{"coral"}, BG: Hex("#ff7f50")},
	{Names: []string{"cornflowerblue"}, BG: Hex("#6495ed")},
	{Names: []string{"cornsilk"}, BG: Hex("#fff8dc")},
	{Names: []string{"crimson"}, BG: Hex("#dc143c")},
	{Names: []string{"cyan"}, BG: Hex("#00ffff")},
	{Names: []string{"darkblue"}, BG: Hex("#00008b")},
	{Names: []string{"darkcyan"}, BG: Hex("#008b8b")},
	{Names: []string{"darkgoldenrod"}, BG: Hex("#b8860b")},
	{Names: []string{"darkgray", "darkgrey"}, BG: Hex("#a9a9a9")},
	{Names: []string{"darkgreen"}, BG: Hex("#006400")},
	{Names: []string{"darkkhaki"}, BG: Hex("#bdb76b")},
	{Names: []string{"darkmagenta"}, BG: Hex("#8b008b")},
	{Names: []string{"darkolivegreen"}, BG: Hex("#556b2f")},
	{Names: []string{"darkorange"}, BG: Hex("#ff8c00")},
	{Names: []string{"darkorchid"}, BG: Hex("#9932cc")},
	{Names: []string{"darkred"}, BG: Hex("#8b0000")},
	{Names: []string{"darksalmon"}, BG: Hex("#e9967a")},
	{Names: []string{"darkseagreen"}, BG: Hex("#8fbc8f")},
	{Names: []string{"darkslateblue"}, BG: Hex("#483d8b")},
	{Names: []string{"darkslategray", "darkslategrey"}, BG: Hex("#2f4f4f")},
	{Names: []string{"darkturquoise"}, BG: Hex("#00ced1")},
	{Names: []string{"darkviolet"}, BG: Hex("#9400d3")},
	{Names: []string{"deeppink"}, BG: Hex("#ff1493")},
	{Names: []string{"deepskyblue"}, BG: Hex("#00bfff")},
	{Names: []string{"dimgray", "dimgrey"}, BG: Hex("#696969")},
	{Names: []string{"dodgerblue"}, BG: Hex("#1e90ff")},
	{Names: []string{"firebrick"}, BG: Hex("#b22222")},
	{Names: []string{"floralwhite"}, BG: Hex("#fffaf0")},
	{Names: []string{"forestgreen"}, BG: Hex("#228b22")},
	{Names: []string{"fuchsia"}, BG: Hex("#ff00ff")},
	{Names: []string{"gainsboro"}, BG: Hex("#dcdcdc")},
	{Names: []string{"ghostwhite"}, BG: Hex("#f8f8ff")},
	{Names: []string{"gold"}, BG: Hex("#ffd700")},
	{Names: []string{"goldenrod"}, BG: Hex("#daa520")},
	{Names: []string{"greenyellow"}, BG: Hex("#adff2f")},
	{Names: []string{"honeydew"}, BG: Hex("#f0fff0")},
	{Names: []string{"hotpink"}, BG: Hex("#ff69b4")},
	{Names: []string{"indianred"}, BG: Hex("#cd5c5c")},
	{Names: []string{"indigo"}, BG: Hex("#4b0082")},
	{Names: []string{"ivory"}, BG: Hex("#fffff0")},
	{Names: []string{"khaki"}, BG: Hex("#f0e68c")},
	{Names: []string{"lavender"}, BG: Hex("#e6e6fa")},
	{Names: []string{"lavenderblush"}, BG: Hex("#fff0f5")},
	{Names: []string{"lawngreen"}, BG: Hex("#7cfc00")},
	{Names: []string{"lemonchiffon"}, BG: Hex("#fffacd")},
	{Names: []string{"lightblue"}, BG: Hex("#add8e6")},
	{Names: []string{"lightcoral"}, BG: Hex("#f08080")},
	{Names: []string{"lightcyan"}, BG: Hex("#e0ffff")},
	{Names: []string{"lightgoldenrodyellow"}, BG: Hex("#fafad2")},
	{Names: []string{"lightgray", "lightgrey"}, BG: Hex("#d3d3d3")},
	{Names: []string{"lightgreen"}, BG: Hex("#90ee90")},
	{Names: []string{"lightpink"}, BG: Hex("#ffb6c1")},
	{Names: []string{"lightsalmon"}, BG: Hex("#ffa07a")},
	{Names: []string{"lightseagreen"}, BG: Hex("#20b2aa")},
	{Names: []string{"lightskyblue"}, BG: Hex("#87cefa")},
	{Names: []string{"lightslategray", "lightslategrey"}, BG: Hex("#778899")},
	{Names: []string{"lightsteelblue"}, BG: Hex("#b0c4de")},
	{Names: []string{"lightyellow"}, BG: Hex("#ffffe0")},
	{Names: []string{"lime"}, BG: Hex("#00ff00")},
	{Names: []string{"limegreen"}, BG: Hex("#32cd32")},
	{Names: []string{"linen"}, BG: Hex("#faf0e6")},
	{Names: []string{"magenta"}, BG: Hex("#ff00ff")},
	{Names: []string{"maroon"}, BG: Hex("#800000")},
	{Names: []string{"mediumaquamarine"}, BG: Hex("#66cdaa")},
	{Names: []string{"mediumblue"}, BG: Hex("#0000cd")},
	{Names: []string{"mediumorchid"}, BG: Hex("#ba55d3")},
	{Names: []string{"mediumpurple"}, BG: Hex("#9370db")},
	{Names: []string{"mediumseagreen"}, BG: Hex("#3cb371")},
	{Names: []string{"mediumslateblue"}, BG: Hex("#7b68ee")},
	{Names: []string{"mediumspringgreen"}, BG: Hex("#00fa9a")},
	{Names: []string{"mediumturquoise"}, BG: Hex("#48d1cc")},
	{Names: []string{"mediumvioletred"}, BG: Hex("#c71585")},
	{Names: []string{"midnightblue"}, BG: Hex("#191970")},
	{Names: []string{"mintcream"}, BG: Hex("#f5fffa")},
	{Names: []string{"mistyrose"}, BG: Hex("#ffe4e1")},
	{Names: []string{"moccasin"}, BG: Hex("#ffe4b5")},
	{Names: []string{"navajowhite"}, BG: Hex("#ffdead")},
	{Names: []string{"navy"}, BG: Hex("#000080")},
	{Names: []string{"oldlace"}, BG: Hex("#fdf5e6")},
	{Names: []string{"olive"}, BG: Hex("#808000")},
	{Names: []string{"olivedrab"}, BG: Hex("#6b8e23")},
	{Names: []string{"orangered"}, BG: Hex("#ff4500")},
	{Names: []string{"orchid"}, BG: Hex("#da70d6")},
	{Names: []string{"palegoldenrod"}, BG: Hex("#eee8aa")},
	{Names: []string{"palegreen"}, BG: Hex("#98fb98")},
	{Names: []string{"paleturquoise"}, BG: Hex("#afeeee")},
	{Names: []string{"palevioletred"}, BG: Hex("#db7093")},
	{Names: []string{"papayawhip"}, BG: Hex("#ffefd5")},
	{Names: []string{"peachpuff"}, BG: Hex("#ffdab9")},
	{Names: []string{"peru"}, BG: Hex("#cd853f")},
	{Names: []string{"pink"}, BG: Hex("#ffc0cb")},
	{Names: []string{"plum"}, BG: Hex("#dda0dd")},
	{Names: []string{"powderblue"}, BG: Hex("#b0e0e6")},
	{Names: []string{"rebeccapurple"}, BG: Hex("#663399")},
	{Names: []string{"rosybrown"}, BG: Hex("#bc8f8f")},
	{Names: []string{"royalblue"}, BG: Hex("#4169e1")},
	{Names: []string{"saddlebrown"}, BG: Hex("#8b4513")},
	{Names: []string{"salmon"}, BG: Hex("#fa8072")},
	{Names: []string{"sandybrown"}, BG: Hex("#f4a460")},
	{Names: []string{"seagreen"}, BG: Hex("#2e8b57")},
	{Names: []string{"seashell"}, BG: Hex("#fff5ee")},
	{Names: []string{"sienna"}, BG: Hex("#a0522d")},
	{Names: []string{"silver"}, BG: Hex("#c0c0c0")},
	{Names: []string{"skyblue"}, BG: Hex("#87ceeb")},
	{Names: []string{"slateblue"}, BG: Hex("#6a5acd")},
	{Names: []string{"slategray", "slategrey"}, BG: Hex("#708090")},
	{Names: []string{"snow"}, BG: Hex("#fffafa")},
	{Names: []string{"springgreen"}, BG: Hex("#00ff7f")},
	{Names: []string{"steelblue"}, BG: Hex("#4682b4")},
	{Names: []string{"tan"}, BG: Hex("#d2b48c")},
	{Names: []string{"thistle"}, BG: Hex("#d8bfd8")},
	{Names: []string{"tomato"}, BG: Hex("#ff6347")},
	{Names: []string{"turquoise"}, BG: Hex("#40e0d0")},
	{Names: []string{"violet"}, BG: Hex("#ee82ee")},
	{Names: []string{"wheat"}, BG: Hex("#f5deb3")},
	{Names: []string{"whitesmoke"}, BG: Hex("#f5f5f5")},
	{Names: []string{"yellowgreen"}, BG: Hex("#9acd32")},
}

func parseColorSpec(spec colorSpec) (rgbColor, error) {
	if spec.rgb != nil {
		return *spec.rgb, nil
	}
	hex := strings.TrimPrefix(spec.hex, "#")
	if len(hex) != 6 {
		return rgbColor{}, fmt.Errorf("expected #RRGGBB color, got %q", spec.hex)
	}

	parsed, err := strconv.ParseUint(hex, 16, 32)
	if err != nil {
		return rgbColor{}, fmt.Errorf("parse %q: %w", spec.hex, err)
	}

	return rgbColor{
		R: int(parsed >> 16 & 0xff),
		G: int(parsed >> 8 & 0xff),
		B: int(parsed & 0xff),
	}, nil
}

func mustColor(spec colorSpec) rgbColor {
	color, err := parseColorSpec(spec)
	if err != nil {
		panic(err)
	}
	return color
}

func foregroundFor(bg rgbColor, fg *colorSpec) rgbColor {
	if fg != nil {
		return mustColor(*fg)
	}
	if contrastYIQ(bg) >= 128 {
		return rgbColor{R: 0, G: 0, B: 0}
	}
	return rgbColor{R: 255, G: 255, B: 255}
}

func contrastYIQ(color rgbColor) int {
	return (color.R*299 + color.G*587 + color.B*114) / 1000
}
