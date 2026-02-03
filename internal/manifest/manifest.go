package manifest

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/rizkyandriawan/eddie/internal/config"
	"github.com/rizkyandriawan/eddie/internal/runner"
)

// Manifest represents the output manifest
type Manifest struct {
	Tool        string           `json:"tool"`
	Version     string           `json:"version"`
	Target      string           `json:"target"`
	GeneratedAt string           `json:"generated_at"`
	Terminal    TerminalInfo     `json:"terminal"`
	Sessions    []SessionManifest `json:"sessions"`
	Summary     Summary          `json:"summary"`
}

// TerminalInfo describes terminal settings
type TerminalInfo struct {
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Theme  string `json:"theme"`
}

// SessionManifest describes a session in the manifest
type SessionManifest struct {
	Name        string               `json:"name"`
	Description string               `json:"description"`
	Cwd         string               `json:"cwd"`
	Screenshots []ScreenshotManifest `json:"screenshots"`
}

// ScreenshotManifest describes a screenshot
type ScreenshotManifest struct {
	Filename string `json:"filename"`
	Prompt   string `json:"prompt,omitempty"`
	WaitMs   int    `json:"wait_ms,omitempty"`
}

// Summary provides aggregate stats
type Summary struct {
	TotalSessions    int `json:"total_sessions"`
	TotalScreenshots int `json:"total_screenshots"`
	Success          int `json:"success"`
	Failed           int `json:"failed"`
}

// Generate creates a manifest from session results
func Generate(cfg *config.Config, results []runner.SessionResult, outputDir string) error {
	manifest := Manifest{
		Tool:        "eddie",
		Version:     "1.0.0",
		Target:      "claude-code",
		GeneratedAt: time.Now().UTC().Format(time.RFC3339),
		Terminal: TerminalInfo{
			Width:  cfg.Terminal.Width,
			Height: cfg.Terminal.Height,
			Theme:  cfg.Theme.Name,
		},
	}

	totalScreenshots := 0
	successSessions := 0
	failedSessions := 0

	for _, result := range results {
		session := SessionManifest{
			Name:        result.Name,
			Description: result.Description,
			Cwd:         result.Cwd,
		}

		for _, ss := range result.Screenshots {
			session.Screenshots = append(session.Screenshots, ScreenshotManifest{
				Filename: ss.Filename,
				Prompt:   ss.Prompt,
				WaitMs:   ss.WaitMs,
			})
			totalScreenshots++
		}

		manifest.Sessions = append(manifest.Sessions, session)

		if result.Error != nil {
			failedSessions++
		} else {
			successSessions++
		}
	}

	manifest.Summary = Summary{
		TotalSessions:    len(results),
		TotalScreenshots: totalScreenshots,
		Success:          successSessions,
		Failed:           failedSessions,
	}

	// Write manifest
	outputPath := filepath.Join(outputDir, "manifest.json")
	data, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(outputPath, data, 0644)
}
