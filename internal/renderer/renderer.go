package renderer

import (
	"image/color"
	"os"
	"path/filepath"

	"github.com/fogleman/gg"
	"github.com/rizkyandriawan/eddie/internal/config"
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

// Renderer renders terminal output to PNG
type Renderer struct {
	theme      config.Theme
	fontPath   string
	charWidth  float64
	charHeight float64
	fontSize   float64
}

// NewRenderer creates a new renderer
func NewRenderer(theme config.Theme) *Renderer {
	r := &Renderer{
		theme: theme,
	}

	// Use larger font size for better resolution
	r.fontSize = theme.FontSize
	if r.fontSize < 16 {
		r.fontSize = 16 // minimum for readability
	}

	// Try to find a monospace font
	r.fontPath = findFont()

	// Character dimensions based on font size
	r.charWidth = r.fontSize * 0.6
	r.charHeight = r.fontSize * 1.2

	return r
}

// findFont looks for a suitable monospace font
func findFont() string {
	// Common monospace font locations
	paths := []string{
		// Linux
		"/usr/share/fonts/truetype/dejavu/DejaVuSansMono.ttf",
		"/usr/share/fonts/truetype/dejavu/DejaVuSansMono-Bold.ttf",
		"/usr/share/fonts/TTF/DejaVuSansMono.ttf",
		"/usr/share/fonts/truetype/liberation/LiberationMono-Regular.ttf",
		"/usr/share/fonts/truetype/ubuntu/UbuntuMono-R.ttf",
		"/usr/share/fonts/truetype/jetbrains-mono/JetBrainsMono-Regular.ttf",
		"/usr/share/fonts/truetype/firacode/FiraCode-Regular.ttf",
		// macOS
		"/System/Library/Fonts/Menlo.ttc",
		"/System/Library/Fonts/Monaco.ttf",
		"/Library/Fonts/SF-Mono-Regular.otf",
		// Windows
		"C:\\Windows\\Fonts\\consola.ttf",
		"C:\\Windows\\Fonts\\cour.ttf",
	}

	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}

	return ""
}

// RenderBuffer renders a ScreenBuffer to a PNG file with colors
func (r *Renderer) RenderBuffer(buffer *ScreenBuffer, outputPath string) error {
	// Calculate image dimensions with small margin
	padding := float64(r.theme.Padding)
	if padding < 20 {
		padding = 20
	}
	margin := 16.0 // Small fixed margin around the terminal

	imgWidth := int(float64(buffer.Width)*r.charWidth + padding*2 + margin*2)
	imgHeight := int(float64(buffer.Height)*r.charHeight + padding*2 + margin*2)

	// Create drawing context
	dc := gg.NewContext(imgWidth, imgHeight)

	// Fill background with same color everywhere
	bgColor := parseHexColor(r.theme.Background)
	dc.SetColor(bgColor)
	dc.Clear()

	// Load font
	if r.fontPath != "" {
		if err := dc.LoadFontFace(r.fontPath, r.fontSize); err != nil {
			// Fall back to basic font if loading fails
			r.fontPath = ""
		}
	}

	// Render each cell (offset by margin + padding)
	offsetX := margin + padding
	offsetY := margin + padding

	for row, line := range buffer.Lines {
		y := offsetY + float64(row)*r.charHeight + r.charHeight*0.85 // baseline offset

		for col, cell := range line.Cells {
			x := offsetX + float64(col)*r.charWidth

			// Draw background if different from default
			if cell.BG != bgColor && cell.BG.A > 0 {
				dc.SetColor(cell.BG)
				dc.DrawRectangle(x, offsetY+float64(row)*r.charHeight, r.charWidth, r.charHeight)
				dc.Fill()
			}

			// Draw character
			if cell.Char != ' ' && cell.Char != 0 {
				dc.SetColor(cell.FG)
				dc.DrawString(string(cell.Char), x, y)
			}
		}
	}

	// Ensure output directory exists
	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	return dc.SavePNG(outputPath)
}

// Render renders plain text terminal output to a PNG file (legacy method)
func (r *Renderer) Render(output string, width, height int, outputPath string) error {
	// Convert plain text to ScreenBuffer (all default colors)
	fgColor := parseHexColor(r.theme.Foreground)
	bgColor := parseHexColor(r.theme.Background)

	buffer := &ScreenBuffer{
		Width:  width,
		Height: height,
		Lines:  make([]ScreenLine, height),
	}

	lines := splitLines(output)
	for row := 0; row < height; row++ {
		line := ScreenLine{
			Cells: make([]ScreenCell, width),
		}

		var lineText string
		if row < len(lines) {
			lineText = lines[row]
		}

		runes := []rune(lineText)
		for col := 0; col < width; col++ {
			ch := ' '
			if col < len(runes) {
				ch = runes[col]
			}

			line.Cells[col] = ScreenCell{
				Char: ch,
				FG:   fgColor,
				BG:   bgColor,
			}
		}

		buffer.Lines[row] = line
	}

	return r.RenderBuffer(buffer, outputPath)
}

func splitLines(s string) []string {
	var lines []string
	var current []rune

	for _, r := range s {
		if r == '\n' {
			lines = append(lines, string(current))
			current = nil
		} else if r != '\r' {
			current = append(current, r)
		}
	}

	if len(current) > 0 {
		lines = append(lines, string(current))
	}

	return lines
}

// parseHexColor parses a hex color string like "#1a1a1a"
func parseHexColor(hex string) color.RGBA {
	if len(hex) == 0 {
		return color.RGBA{0, 0, 0, 255}
	}

	if hex[0] == '#' {
		hex = hex[1:]
	}

	if len(hex) != 6 {
		return color.RGBA{0, 0, 0, 255}
	}

	r := hexToByte(hex[0:2])
	g := hexToByte(hex[2:4])
	b := hexToByte(hex[4:6])

	return color.RGBA{r, g, b, 255}
}

func hexToByte(s string) uint8 {
	var val uint8
	for _, c := range s {
		val *= 16
		switch {
		case c >= '0' && c <= '9':
			val += uint8(c - '0')
		case c >= 'a' && c <= 'f':
			val += uint8(c - 'a' + 10)
		case c >= 'A' && c <= 'F':
			val += uint8(c - 'A' + 10)
		}
	}
	return val
}

// darkenColor makes a color darker by the given factor (0-1)
func darkenColor(c color.RGBA, factor float64) color.RGBA {
	return color.RGBA{
		R: uint8(float64(c.R) * (1 - factor)),
		G: uint8(float64(c.G) * (1 - factor)),
		B: uint8(float64(c.B) * (1 - factor)),
		A: c.A,
	}
}

// drawRoundedRect draws a rounded rectangle
func drawRoundedRect(dc *gg.Context, x, y, w, h, r float64) {
	dc.NewSubPath()
	dc.MoveTo(x+r, y)
	dc.LineTo(x+w-r, y)
	dc.DrawArc(x+w-r, y+r, r, -gg.Radians(90), 0)
	dc.LineTo(x+w, y+h-r)
	dc.DrawArc(x+w-r, y+h-r, r, 0, gg.Radians(90))
	dc.LineTo(x+r, y+h)
	dc.DrawArc(x+r, y+h-r, r, gg.Radians(90), gg.Radians(180))
	dc.LineTo(x, y+r)
	dc.DrawArc(x+r, y+r, r, gg.Radians(180), gg.Radians(270))
	dc.ClosePath()
}
