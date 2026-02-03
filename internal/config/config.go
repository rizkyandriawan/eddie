package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Output   string    `yaml:"output"`
	Manifest bool      `yaml:"manifest"`
	Terminal Terminal  `yaml:"terminal"`
	Theme    Theme     `yaml:"theme"`
	Sessions []Session `yaml:"sessions"`
}

type Terminal struct {
	Width  int `yaml:"width"`
	Height int `yaml:"height"`
}

type Theme struct {
	Name       string            `yaml:"name"`
	Background string            `yaml:"background"`
	Foreground string            `yaml:"foreground"`
	Font       string            `yaml:"font"`
	FontSize   float64           `yaml:"font_size"`
	Padding    int               `yaml:"padding"`
	Colors     map[string]string `yaml:"colors"`
}

type Session struct {
	Name        string   `yaml:"name"`
	Description string   `yaml:"description"`
	Cwd         string   `yaml:"cwd"`
	Setup       []string `yaml:"setup"`
	Prompts     []Prompt `yaml:"prompts"`
}

type Prompt struct {
	Input       string `yaml:"input"`
	Key         string `yaml:"key"`
	Wait        int    `yaml:"wait"`
	WaitUntil   string `yaml:"wait_until"`
	Timeout     int    `yaml:"timeout"`
	Capture     bool   `yaml:"capture"`
	CaptureName string `yaml:"capture_name"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		Output: "./screenshots",
		Terminal: Terminal{
			Width:  120,
			Height: 40,
		},
		Theme: Theme{
			Name:       "dark",
			Background: "#1a1a1a",
			Foreground: "#d4d4d4",
			Font:       "monospace",
			FontSize:   14,
			Padding:    20,
		},
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	// Expand ~ in paths
	cfg.Output = expandPath(cfg.Output)
	for i := range cfg.Sessions {
		cfg.Sessions[i].Cwd = expandPath(cfg.Sessions[i].Cwd)
	}

	return cfg, nil
}

func expandPath(path string) string {
	if len(path) > 0 && path[0] == '~' {
		home, _ := os.UserHomeDir()
		return filepath.Join(home, path[1:])
	}
	return path
}
