package ansi

import (
	"image/color"
	"regexp"
	"strconv"
	"strings"
)

// Segment represents a piece of text with styling
type Segment struct {
	Text       string
	Foreground color.RGBA
	Background color.RGBA
	Bold       bool
	Italic     bool
	Underline  bool
}

// Screen represents the terminal screen state
type Screen struct {
	Rows   [][]Segment
	Width  int
	Height int
}

// Default ANSI colors (standard 16 colors)
var defaultColors = []color.RGBA{
	{0, 0, 0, 255},       // 0: Black
	{205, 49, 49, 255},   // 1: Red
	{13, 188, 121, 255},  // 2: Green
	{229, 229, 16, 255},  // 3: Yellow
	{36, 114, 200, 255},  // 4: Blue
	{188, 63, 188, 255},  // 5: Magenta
	{17, 168, 205, 255},  // 6: Cyan
	{229, 229, 229, 255}, // 7: White
	{102, 102, 102, 255}, // 8: Bright Black
	{241, 76, 76, 255},   // 9: Bright Red
	{35, 209, 139, 255},  // 10: Bright Green
	{245, 245, 67, 255},  // 11: Bright Yellow
	{59, 142, 234, 255},  // 12: Bright Blue
	{214, 112, 214, 255}, // 13: Bright Magenta
	{41, 184, 219, 255},  // 14: Bright Cyan
	{255, 255, 255, 255}, // 15: Bright White
}

// ANSI escape sequence regex
var ansiRegex = regexp.MustCompile(`\x1b\[([0-9;]*)([A-Za-z])`)

// Parser parses ANSI escape sequences
type Parser struct {
	defaultFg color.RGBA
	defaultBg color.RGBA
	fg        color.RGBA
	bg        color.RGBA
	bold      bool
	italic    bool
	underline bool
}

// NewParser creates a new ANSI parser
func NewParser(fg, bg color.RGBA) *Parser {
	return &Parser{
		defaultFg: fg,
		defaultBg: bg,
		fg:        fg,
		bg:        bg,
	}
}

// Parse parses ANSI text into segments
func (p *Parser) Parse(text string) []Segment {
	var segments []Segment
	var currentText strings.Builder

	i := 0
	for i < len(text) {
		if text[i] == '\x1b' && i+1 < len(text) && text[i+1] == '[' {
			// Found escape sequence
			if currentText.Len() > 0 {
				segments = append(segments, Segment{
					Text:       currentText.String(),
					Foreground: p.fg,
					Background: p.bg,
					Bold:       p.bold,
					Italic:     p.italic,
					Underline:  p.underline,
				})
				currentText.Reset()
			}

			// Find end of escape sequence
			end := i + 2
			for end < len(text) && !isLetter(text[end]) {
				end++
			}
			if end < len(text) {
				end++ // include the letter
				p.processEscape(text[i:end])
				i = end
				continue
			}
		}

		currentText.WriteByte(text[i])
		i++
	}

	if currentText.Len() > 0 {
		segments = append(segments, Segment{
			Text:       currentText.String(),
			Foreground: p.fg,
			Background: p.bg,
			Bold:       p.bold,
			Italic:     p.italic,
			Underline:  p.underline,
		})
	}

	return segments
}

func isLetter(b byte) bool {
	return (b >= 'A' && b <= 'Z') || (b >= 'a' && b <= 'z')
}

func (p *Parser) processEscape(seq string) {
	matches := ansiRegex.FindStringSubmatch(seq)
	if len(matches) < 3 {
		return
	}

	params := matches[1]
	cmd := matches[2]

	if cmd != "m" {
		// Only handle SGR (Select Graphic Rendition) for now
		return
	}

	if params == "" || params == "0" {
		p.reset()
		return
	}

	codes := strings.Split(params, ";")
	for i := 0; i < len(codes); i++ {
		code, _ := strconv.Atoi(codes[i])
		switch {
		case code == 0:
			p.reset()
		case code == 1:
			p.bold = true
		case code == 3:
			p.italic = true
		case code == 4:
			p.underline = true
		case code == 22:
			p.bold = false
		case code == 23:
			p.italic = false
		case code == 24:
			p.underline = false
		case code >= 30 && code <= 37:
			p.fg = defaultColors[code-30]
			if p.bold {
				p.fg = defaultColors[code-30+8]
			}
		case code == 38:
			// Extended foreground color
			if i+1 < len(codes) {
				next, _ := strconv.Atoi(codes[i+1])
				if next == 5 && i+2 < len(codes) {
					// 256 color mode
					colorIdx, _ := strconv.Atoi(codes[i+2])
					p.fg = get256Color(colorIdx)
					i += 2
				} else if next == 2 && i+4 < len(codes) {
					// RGB mode
					r, _ := strconv.Atoi(codes[i+2])
					g, _ := strconv.Atoi(codes[i+3])
					b, _ := strconv.Atoi(codes[i+4])
					p.fg = color.RGBA{uint8(r), uint8(g), uint8(b), 255}
					i += 4
				}
			}
		case code == 39:
			p.fg = p.defaultFg
		case code >= 40 && code <= 47:
			p.bg = defaultColors[code-40]
		case code == 48:
			// Extended background color
			if i+1 < len(codes) {
				next, _ := strconv.Atoi(codes[i+1])
				if next == 5 && i+2 < len(codes) {
					colorIdx, _ := strconv.Atoi(codes[i+2])
					p.bg = get256Color(colorIdx)
					i += 2
				} else if next == 2 && i+4 < len(codes) {
					r, _ := strconv.Atoi(codes[i+2])
					g, _ := strconv.Atoi(codes[i+3])
					b, _ := strconv.Atoi(codes[i+4])
					p.bg = color.RGBA{uint8(r), uint8(g), uint8(b), 255}
					i += 4
				}
			}
		case code == 49:
			p.bg = p.defaultBg
		case code >= 90 && code <= 97:
			p.fg = defaultColors[code-90+8]
		case code >= 100 && code <= 107:
			p.bg = defaultColors[code-100+8]
		}
	}
}

func (p *Parser) reset() {
	p.fg = p.defaultFg
	p.bg = p.defaultBg
	p.bold = false
	p.italic = false
	p.underline = false
}

// get256Color returns a color from the 256-color palette
func get256Color(idx int) color.RGBA {
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
		// 24 grayscale colors
		gray := uint8((idx-232)*10 + 8)
		return color.RGBA{gray, gray, gray, 255}
	}
}

// StripANSI removes all ANSI escape sequences from text
func StripANSI(text string) string {
	return ansiRegex.ReplaceAllString(text, "")
}
