package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/rizkyandriawan/eddie/internal/config"
	"github.com/rizkyandriawan/eddie/internal/manifest"
	"github.com/rizkyandriawan/eddie/internal/runner"
)

var version = "1.0.0"

func main() {
	// CLI flags
	configPath := flag.String("c", "", "Path to YAML config file (required)")
	outputDir := flag.String("o", "", "Output directory (overrides config)")
	generateManifest := flag.Bool("manifest", false, "Generate manifest.json")
	showVersion := flag.Bool("version", false, "Show version")
	showHelp := flag.Bool("help", false, "Show help")

	flag.Parse()

	if *showVersion {
		fmt.Printf("eddie v%s\n", version)
		os.Exit(0)
	}

	if *showHelp || *configPath == "" {
		printHelp()
		os.Exit(0)
	}

	// Load config
	cfg, err := config.Load(*configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(2)
	}

	// Override output directory if specified
	if *outputDir != "" {
		cfg.Output = *outputDir
	}

	// Override manifest flag
	if *generateManifest {
		cfg.Manifest = true
	}

	// Create output directory
	if err := os.MkdirAll(cfg.Output, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating output directory: %v\n", err)
		os.Exit(2)
	}

	// Run sessions
	fmt.Println("Eddie - Claude Code Screenshot Tool")
	fmt.Println("====================================")
	fmt.Printf("Output: %s\n\n", cfg.Output)

	r := runner.NewRunner(cfg)
	results, err := r.RunAll()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running sessions: %v\n", err)
		os.Exit(1)
	}

	// Generate manifest
	if cfg.Manifest {
		if err := manifest.Generate(cfg, results, cfg.Output); err != nil {
			fmt.Fprintf(os.Stderr, "Error generating manifest: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("\nManifest: %s/manifest.json\n", cfg.Output)
	}

	// Summary
	fmt.Println("\n====================================")
	totalScreenshots := 0
	successSessions := 0
	for _, result := range results {
		totalScreenshots += len(result.Screenshots)
		if result.Error == nil {
			successSessions++
		}
	}
	fmt.Printf("Sessions: %d/%d successful\n", successSessions, len(results))
	fmt.Printf("Screenshots: %d captured\n", totalScreenshots)

	// Exit code based on results
	if successSessions < len(results) {
		os.Exit(1)
	}
}

func printHelp() {
	fmt.Println(`Eddie - Claude Code Screenshot Tool 🖤

"We are Eddie."

USAGE:
    eddie -c <config.yaml> [options]

OPTIONS:
    -c <path>       Path to YAML config file (required)
    -o <path>       Output directory (overrides config)
    --manifest      Generate manifest.json
    --version       Show version
    --help          Show this help

EXAMPLE CONFIG:

    output: ./screenshots
    manifest: true

    terminal:
      width: 120
      height: 40

    theme:
      name: dark
      background: "#1a1a1a"
      foreground: "#d4d4d4"

    sessions:
      - name: 01-hello
        description: "Basic greeting"
        cwd: ~/projects/myapp
        prompts:
          - input: "what files are in this project?"
            wait: 5000
            capture: true

EXAMPLES:

    # Basic usage
    eddie -c config.yaml

    # Custom output directory
    eddie -c config.yaml -o ./docs/screenshots

    # With manifest
    eddie -c config.yaml --manifest

For more information, visit: https://github.com/rizkyandriawan/eddie
`)
}
