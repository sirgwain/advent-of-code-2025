package advent

import (
	"strconv"
	"strings"
)

// BrailleColorFunc decides how to color a single Braille cell.
//
// cellX, cellY: braille cell coordinates (0-based) in the output grid
// dots: 8 bytes, one per dot (1..8), 0 = off, non-zero = whatever was in grid
//
// Return values:
//
//	fg:  0–255 for 256-color foreground, or -1 for "no fg color"
//	bg:  0–255 for 256-color background, or -1 for "no bg color"
//	ok:  if false, RenderBrailleWithColor will render without any color
type BrailleColorFunc func(cellX, cellY int, dots [8]byte) (fg, bg int, ok bool)

// RenderBraille renders a grid of bytes as a string of Unicode braille cells.
// grid[y][x] == 0 => pixel off, != 0 => pixel on.
func RenderBraille(grid [][]byte) string {
	if len(grid) == 0 {
		return ""
	}

	height := len(grid)
	width := 0
	for _, row := range grid {
		if len(row) > width {
			width = len(row)
		}
	}

	var b strings.Builder

	// Step by rows of 4 pixels
	for y := 0; y < height; y += 4 {
		// Step by columns of 2 pixels
		for x := 0; x < width; x += 2 {
			var mask uint8

			// Helpers to test if a pixel is "on" safely
			on := func(px, py int) bool {
				if py < 0 || py >= height {
					return false
				}
				row := grid[py]
				if px < 0 || px >= len(row) {
					return false
				}
				return row[px] != 0
			}

			// Map pixels to braille bits
			if on(x+0, y+0) {
				mask |= 1 << 0
			} // dot1
			if on(x+0, y+1) {
				mask |= 1 << 1
			} // dot2
			if on(x+0, y+2) {
				mask |= 1 << 2
			} // dot3
			if on(x+1, y+0) {
				mask |= 1 << 3
			} // dot4
			if on(x+1, y+1) {
				mask |= 1 << 4
			} // dot5
			if on(x+1, y+2) {
				mask |= 1 << 5
			} // dot6
			if on(x+0, y+3) {
				mask |= 1 << 6
			} // dot7
			if on(x+1, y+3) {
				mask |= 1 << 7
			} // dot8

			r := rune(0x2800 + int(mask))
			b.WriteRune(r)
		}
		b.WriteByte('\n')
	}

	return b.String()
}

// RenderBrailleWithColor renders a grid of bytes as a string of Unicode braille
// cells and uses colorFn to optionally wrap each cell in ANSI color codes.
//
// grid[y][x] == 0 => pixel off
// grid[y][x] != 0 => pixel "on" with that value (used only by colorFn)
func RenderBrailleWithColor(grid [][]byte, colorFn BrailleColorFunc) string {
	if len(grid) == 0 {
		return ""
	}

	height := len(grid)
	width := 0
	for _, row := range grid {
		if len(row) > width {
			width = len(row)
		}
	}

	const brailleBase = rune(0x2800)

	var b strings.Builder

	// Helper to safely read a byte from the grid
	valueAt := func(px, py int) byte {
		if py < 0 || py >= height {
			return 0
		}
		row := grid[py]
		if px < 0 || px >= len(row) {
			return 0
		}
		return row[px]
	}

	// Step by rows of 4 pixels → one braille row
	for y := 0; y < height; y += 4 {
		// Step by columns of 2 pixels → one braille column
		for x := 0; x < width; x += 2 {
			var mask uint8
			var dots [8]byte

			// Map pixels to dots + mask + per-dot values
			// dot1
			if v := valueAt(x+0, y+0); v != 0 {
				mask |= 1 << 0
				dots[0] = v
			}
			// dot2
			if v := valueAt(x+0, y+1); v != 0 {
				mask |= 1 << 1
				dots[1] = v
			}
			// dot3
			if v := valueAt(x+0, y+2); v != 0 {
				mask |= 1 << 2
				dots[2] = v
			}
			// dot4
			if v := valueAt(x+1, y+0); v != 0 {
				mask |= 1 << 3
				dots[3] = v
			}
			// dot5
			if v := valueAt(x+1, y+1); v != 0 {
				mask |= 1 << 4
				dots[4] = v
			}
			// dot6
			if v := valueAt(x+1, y+2); v != 0 {
				mask |= 1 << 5
				dots[5] = v
			}
			// dot7
			if v := valueAt(x+0, y+3); v != 0 {
				mask |= 1 << 6
				dots[6] = v
			}
			// dot8
			if v := valueAt(x+1, y+3); v != 0 {
				mask |= 1 << 7
				dots[7] = v
			}

			r := brailleBase + rune(mask)

			cellX := x / 2 // braille X coordinate
			cellY := y / 4 // braille Y coordinate
			colored := false

			if colorFn != nil {
				if fg, bg, ok := colorFn(cellX, cellY, dots); ok {
					// Build ANSI prefix
					// Examples:
					//   fg only:   \x1b[38;5;FGm
					//   bg only:   \x1b[48;5;BGm
					//   both:      \x1b[38;5;FG;48;5;BGm
					if fg >= 0 || bg >= 0 {
						b.WriteString("\x1b[")
						first := true
						if fg >= 0 {
							b.WriteString("38;5;")
							b.WriteString(strconv.Itoa(fg))
							first = false
						}
						if bg >= 0 {
							if !first {
								b.WriteByte(';')
							}
							b.WriteString("48;5;")
							b.WriteString(strconv.Itoa(bg))
						}
						b.WriteByte('m')
						b.WriteRune(r)
						b.WriteString("\x1b[0m")
						colored = true
					}
				}
			}

			if !colored {
				// No color or colorFn said "no color"
				b.WriteRune(r)
			}
		}
		b.WriteByte('\n')
	}

	return b.String()
}

// DensityColor colors cells based on how many dots are on.
// 0 dots  -> gray
// 1–2     -> blue
// 3–4     -> green
// 5–6     -> yellow
// 7–8     -> red
func DensityColor(cellX, cellY int, dots [8]byte) (fg, bg int, ok bool) {
	count := 0
	for _, v := range dots {
		if v != 0 {
			count++
		}
	}
	if count == 0 {
		return 240, -1, true // dark gray, no background
	}

	switch {
	case count <= 2:
		return 33, -1, true // blue-ish
	case count <= 4:
		return 82, -1, true // green
	case count <= 6:
		return 226, -1, true // yellow
	default:
		return 196, -1, true // red
	}
}
