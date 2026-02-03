package runner

import (
	"image/color"

	"github.com/hinshun/vt10x"
)

// ScreenCell represents a single cell with character and colors
type ScreenCell struct {
	Char rune
	FG   color.RGBA
	BG   color.RGBA
}

// ScreenLine represents a line of cells
type ScreenLine struct {
	Cells []ScreenCell
}

// ScreenBuffer represents the entire screen state with colors
type ScreenBuffer struct {
	Lines  []ScreenLine
	Width  int
	Height int
}

// Default terminal colors (matching vt10x)
var defaultColors = []color.RGBA{
	{0, 0, 0, 255},       // 0: Black
	{205, 49, 49, 255},   // 1: Red
	{13, 188, 121, 255},  // 2: Green
	{229, 229, 16, 255},  // 3: Yellow
	{36, 114, 200, 255},  // 4: Blue
	{188, 63, 188, 255},  // 5: Magenta
	{17, 168, 205, 255},  // 6: Cyan
	{229, 229, 229, 255}, // 7: White
	{102, 102, 102, 255}, // 8: Bright Black (Gray)
	{241, 76, 76, 255},   // 9: Bright Red
	{35, 209, 139, 255},  // 10: Bright Green
	{245, 245, 67, 255},  // 11: Bright Yellow
	{59, 142, 234, 255},  // 12: Bright Blue
	{214, 112, 214, 255}, // 13: Bright Magenta
	{41, 184, 219, 255},  // 14: Bright Cyan
	{255, 255, 255, 255}, // 15: Bright White
}

// vt10xColorToRGBA converts vt10x color to RGBA
func vt10xColorToRGBA(c vt10x.Color, defaultColor color.RGBA) color.RGBA {
	// vt10x.Color is an int representing color index or RGB
	colorVal := int(c)

	// Default color (0 usually means default)
	if colorVal == 0 {
		return defaultColor
	}

	// Standard 16 colors (1-16, but stored as 0-15 index)
	if colorVal >= 1 && colorVal <= 16 {
		return defaultColors[colorVal-1]
	}

	// 256 color palette (17-256)
	if colorVal >= 17 && colorVal <= 256 {
		idx := colorVal - 1
		if idx < 16 {
			return defaultColors[idx]
		} else if idx < 232 {
			// 216 colors (6x6x6 cube)
			idx -= 16
			r := uint8((idx / 36) * 51)
			g := uint8(((idx / 6) % 6) * 51)
			b := uint8((idx % 6) * 51)
			return color.RGBA{r, g, b, 255}
		} else {
			// 24 grayscale
			gray := uint8((idx-232)*10 + 8)
			return color.RGBA{gray, gray, gray, 255}
		}
	}

	// RGB color (encoded as high bits)
	if colorVal > 256 {
		// Try to extract RGB - vt10x might encode it differently
		r := uint8((colorVal >> 16) & 0xFF)
		g := uint8((colorVal >> 8) & 0xFF)
		b := uint8(colorVal & 0xFF)
		if r > 0 || g > 0 || b > 0 {
			return color.RGBA{r, g, b, 255}
		}
	}

	return defaultColor
}

// GetScreenBuffer extracts the full screen state with colors from vt10x
func GetScreenBuffer(term vt10x.Terminal, cols, rows int, defaultFG, defaultBG color.RGBA) *ScreenBuffer {
	buffer := &ScreenBuffer{
		Width:  cols,
		Height: rows,
		Lines:  make([]ScreenLine, rows),
	}

	for row := 0; row < rows; row++ {
		line := ScreenLine{
			Cells: make([]ScreenCell, cols),
		}

		for col := 0; col < cols; col++ {
			cell := term.Cell(col, row)

			ch := cell.Char
			if ch == 0 {
				ch = ' '
			}

			fg := vt10xColorToRGBA(cell.FG, defaultFG)
			bg := vt10xColorToRGBA(cell.BG, defaultBG)

			line.Cells[col] = ScreenCell{
				Char: ch,
				FG:   fg,
				BG:   bg,
			}
		}

		buffer.Lines[row] = line
	}

	return buffer
}
