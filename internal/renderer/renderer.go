package renderer

import (
	"image/color"
	"os"
	"path/filepath"
	"strings"

	"github.com/fogleman/gg"
	"github.com/rizkyandriawan/eddie/internal/ansi"
	"github.com/rizkyandriawan/eddie/internal/config"
)

// Renderer renders terminal output to PNG
type Renderer struct {
	theme      config.Theme
	fontPath   string
	charWidth  float64
	charHeight float64
}

// NewRenderer creates a new renderer
func NewRenderer(theme config.Theme) *Renderer {
	r := &Renderer{
		theme: theme,
	}

	// Try to find a monospace font
	r.fontPath = findFont()

	// Estimate character dimensions based on font size
	r.charWidth = theme.FontSize * 0.6
	r.charHeight = theme.FontSize * 1.4

	return r
}

// findFont looks for a suitable monospace font
func findFont() string {
	// Common monospace font locations
	paths := []string{
		"/usr/share/fonts/truetype/dejavu/DejaVuSansMono.ttf",
		"/usr/share/fonts/TTF/DejaVuSansMono.ttf",
		"/usr/share/fonts/truetype/liberation/LiberationMono-Regular.ttf",
		"/usr/share/fonts/truetype/ubuntu/UbuntuMono-R.ttf",
		"/System/Library/Fonts/Menlo.ttc",
		"/System/Library/Fonts/Monaco.ttf",
		"C:\\Windows\\Fonts\\consola.ttf",
	}

	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}

	return ""
}

// Render renders terminal output to a PNG file
func (r *Renderer) Render(output string, width, height int, outputPath string) error {
	// Parse terminal dimensions
	cols := width
	rows := height

	// Calculate image dimensions
	padding := float64(r.theme.Padding)
	imgWidth := int(float64(cols)*r.charWidth + padding*2)
	imgHeight := int(float64(rows)*r.charHeight + padding*2)

	// Create drawing context
	dc := gg.NewContext(imgWidth, imgHeight)

	// Fill background
	bgColor := parseHexColor(r.theme.Background)
	dc.SetColor(bgColor)
	dc.Clear()

	// Load font
	if r.fontPath != "" {
		if err := dc.LoadFontFace(r.fontPath, r.theme.FontSize); err != nil {
			// Fall back to basic font if loading fails
			r.fontPath = ""
		}
	}

	// Parse and render the output
	fgColor := parseHexColor(r.theme.Foreground)
	parser := ansi.NewParser(fgColor, bgColor)

	lines := strings.Split(output, "\n")
	y := padding + r.charHeight

	for lineNum, line := range lines {
		if lineNum >= rows {
			break
		}

		// Parse ANSI codes in this line
		segments := parser.Parse(line)

		x := padding
		for _, seg := range segments {
			// Handle special characters
			text := seg.Text
			text = strings.ReplaceAll(text, "\r", "")
			text = strings.ReplaceAll(text, "\t", "    ")

			// Skip control characters
			var cleanText strings.Builder
			for _, ch := range text {
				if ch >= 32 || ch == '\t' {
					cleanText.WriteRune(ch)
				}
			}
			text = cleanText.String()

			if text == "" {
				continue
			}

			// Draw background if different from default
			if seg.Background != bgColor {
				textWidth := float64(len(text)) * r.charWidth
				dc.SetColor(seg.Background)
				dc.DrawRectangle(x, y-r.charHeight+4, textWidth, r.charHeight)
				dc.Fill()
			}

			// Draw text
			dc.SetColor(seg.Foreground)
			dc.DrawString(text, x, y)

			// Draw underline if needed
			if seg.Underline {
				textWidth := float64(len(text)) * r.charWidth
				dc.SetLineWidth(1)
				dc.DrawLine(x, y+2, x+textWidth, y+2)
				dc.Stroke()
			}

			x += float64(len(text)) * r.charWidth
		}

		y += r.charHeight
	}

	// Ensure output directory exists
	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	return dc.SavePNG(outputPath)
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
