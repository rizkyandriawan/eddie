package runner

import (
	"fmt"
	"image/color"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/creack/pty"
	"github.com/hinshun/vt10x"
	"github.com/rizkyandriawan/eddie/internal/config"
	"github.com/rizkyandriawan/eddie/internal/renderer"
)

// Runner runs Claude Code sessions and captures output
type Runner struct {
	config   *config.Config
	renderer *renderer.Renderer
}

// NewRunner creates a new runner
func NewRunner(cfg *config.Config) *Runner {
	return &Runner{
		config:   cfg,
		renderer: renderer.NewRenderer(cfg.Theme),
	}
}

// Screenshot represents a captured screenshot
type Screenshot struct {
	Name        string
	Filename    string
	Description string
	Prompt      string
	WaitMs      int
}

// SessionResult holds the results of a session
type SessionResult struct {
	Name        string
	Description string
	Cwd         string
	Screenshots []Screenshot
	Error       error
}

// RunSession runs a single Claude Code session
func (r *Runner) RunSession(session config.Session) (*SessionResult, error) {
	result := &SessionResult{
		Name:        session.Name,
		Description: session.Description,
		Cwd:         session.Cwd,
	}

	// Run setup commands first
	for _, setupCmd := range session.Setup {
		cmd := exec.Command("sh", "-c", setupCmd)
		cmd.Dir = session.Cwd
		if err := cmd.Run(); err != nil {
			return result, fmt.Errorf("setup command failed: %s: %w", setupCmd, err)
		}
	}

	// Create virtual terminal
	cols := r.config.Terminal.Width
	rows := r.config.Terminal.Height
	term := vt10x.New(vt10x.WithSize(cols, rows))

	// Determine command to run
	cmdStr := session.Command
	if cmdStr == "" {
		cmdStr = "claude"
	}

	// Start Claude Code in PTY
	cmd := exec.Command("sh", "-c", cmdStr)
	cmd.Dir = session.Cwd
	cmd.Env = append(os.Environ(), "TERM=xterm-256color")

	ptmx, err := pty.StartWithSize(cmd, &pty.Winsize{
		Rows: uint16(rows),
		Cols: uint16(cols),
	})
	if err != nil {
		return result, fmt.Errorf("failed to start PTY: %w", err)
	}
	defer ptmx.Close()

	// Channel to signal process exit
	done := make(chan struct{})
	var mu sync.Mutex

	// Read output continuously and feed to virtual terminal
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := ptmx.Read(buf)
			if err != nil {
				close(done)
				return
			}
			mu.Lock()
			term.Write(buf[:n])
			mu.Unlock()
		}
	}()

	// Process each prompt
	for i, prompt := range session.Prompts {
		// Determine capture name
		captureName := prompt.CaptureName
		if captureName == "" {
			if len(session.Prompts) == 1 {
				captureName = session.Name
			} else {
				captureName = fmt.Sprintf("%s-%d", session.Name, i+1)
			}
		}

		// Wait first (before any input)
		if prompt.WaitUntil != "" {
			// Wait until specific text appears
			timeout := prompt.Timeout
			if timeout == 0 {
				timeout = 30000
			}
			err := waitUntilScreen(term, &mu, prompt.WaitUntil, time.Duration(timeout)*time.Millisecond)
			if err != nil {
				return result, fmt.Errorf("timeout waiting for '%s': %w", prompt.WaitUntil, err)
			}
		} else if prompt.Wait > 0 {
			// Wait fixed duration
			time.Sleep(time.Duration(prompt.Wait) * time.Millisecond)
		}

		// Send input or keystroke (after wait)
		if prompt.Input != "" {
			// Send the text input (without automatic newline)
			_, err := ptmx.Write([]byte(prompt.Input))
			if err != nil {
				return result, fmt.Errorf("failed to send input: %w", err)
			}
		}
		if prompt.Key != "" {
			// Send keystroke
			keyBytes := keyToBytes(prompt.Key)
			_, err := ptmx.Write(keyBytes)
			if err != nil {
				return result, fmt.Errorf("failed to send key: %w", err)
			}
		}

		// Small delay after input to let terminal update
		if prompt.Input != "" || prompt.Key != "" {
			time.Sleep(200 * time.Millisecond)
		}

		// Capture screenshot
		if prompt.Capture {
			outputPath := fmt.Sprintf("%s/%s.png", r.config.Output, captureName)

			// Get screen buffer with colors from virtual terminal
			mu.Lock()
			defaultFG := parseHexColor(r.config.Theme.Foreground)
			defaultBG := parseHexColor(r.config.Theme.Background)
			screenBuffer := GetScreenBuffer(term, cols, rows, defaultFG, defaultBG)
			mu.Unlock()

			// Convert to renderer's ScreenBuffer type
			renderBuffer := convertToRenderBuffer(screenBuffer)

			err := r.renderer.RenderBuffer(renderBuffer, outputPath)
			if err != nil {
				return result, fmt.Errorf("failed to render screenshot: %w", err)
			}

			result.Screenshots = append(result.Screenshots, Screenshot{
				Name:        captureName,
				Filename:    captureName + ".png",
				Description: session.Description,
				Prompt:      prompt.Input,
				WaitMs:      prompt.Wait,
			})

			fmt.Printf("  Captured: %s\n", outputPath)
		}
	}

	// Send Ctrl+C to exit Claude Code
	ptmx.Write([]byte{3}) // Ctrl+C
	time.Sleep(500 * time.Millisecond)
	ptmx.Write([]byte{3}) // Again to make sure

	// Wait for process to exit
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		cmd.Process.Kill()
	}

	return result, nil
}


// RunAll runs all sessions
func (r *Runner) RunAll() ([]SessionResult, error) {
	var results []SessionResult

	for _, session := range r.config.Sessions {
		fmt.Printf("Running session: %s\n", session.Name)

		result, err := r.RunSession(session)
		if err != nil {
			result.Error = err
			fmt.Printf("  Error: %v\n", err)
		}

		results = append(results, *result)
	}

	return results, nil
}

// waitUntilScreen waits until the text appears on the virtual terminal screen
func waitUntilScreen(term vt10x.Terminal, mu *sync.Mutex, text string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		mu.Lock()
		screen := term.String()
		mu.Unlock()
		if strings.Contains(screen, text) {
			return nil
		}
		time.Sleep(100 * time.Millisecond)
	}
	return fmt.Errorf("timeout")
}

// keyToBytes converts key name to bytes
func keyToBytes(key string) []byte {
	key = strings.ToLower(key)

	switch key {
	case "enter", "return":
		return []byte{13}
	case "tab":
		return []byte{9}
	case "escape", "esc":
		return []byte{27}
	case "backspace":
		return []byte{127}
	case "delete":
		return []byte{27, 91, 51, 126}
	case "up":
		return []byte{27, 91, 65}
	case "down":
		return []byte{27, 91, 66}
	case "right":
		return []byte{27, 91, 67}
	case "left":
		return []byte{27, 91, 68}
	case "home":
		return []byte{27, 91, 72}
	case "end":
		return []byte{27, 91, 70}
	case "ctrl+c":
		return []byte{3}
	case "ctrl+d":
		return []byte{4}
	case "ctrl+z":
		return []byte{26}
	case "ctrl+l":
		return []byte{12}
	case "y":
		return []byte{'y'}
	case "n":
		return []byte{'n'}
	default:
		// Single character
		if len(key) == 1 {
			return []byte(key)
		}
		return []byte(key)
	}
}

// CaptureScreen captures the current PTY output to a file
func CaptureScreen(ptmx *os.File, output io.Writer) error {
	buf := make([]byte, 32*1024)
	n, err := ptmx.Read(buf)
	if err != nil && err != io.EOF {
		return err
	}
	output.Write(buf[:n])
	return nil
}

// convertToRenderBuffer converts runner.ScreenBuffer to renderer.ScreenBuffer
func convertToRenderBuffer(sb *ScreenBuffer) *renderer.ScreenBuffer {
	rb := &renderer.ScreenBuffer{
		Width:  sb.Width,
		Height: sb.Height,
		Lines:  make([]renderer.ScreenLine, len(sb.Lines)),
	}

	for i, line := range sb.Lines {
		rl := renderer.ScreenLine{
			Cells: make([]renderer.ScreenCell, len(line.Cells)),
		}
		for j, cell := range line.Cells {
			rl.Cells[j] = renderer.ScreenCell{
				Char: cell.Char,
				FG:   cell.FG,
				BG:   cell.BG,
			}
		}
		rb.Lines[i] = rl
	}

	return rb
}

// parseHexColor parses a hex color string like "#1a1a1a"
func parseHexColor(hex string) color.RGBA {
	if len(hex) == 0 {
		return color.RGBA{212, 212, 212, 255} // default gray
	}

	if hex[0] == '#' {
		hex = hex[1:]
	}

	if len(hex) != 6 {
		return color.RGBA{212, 212, 212, 255}
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
