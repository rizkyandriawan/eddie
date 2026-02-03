# Eddie 🖤

> *"We are Eddie."*

CLI screenshot tool for capturing **Claude Code** sessions. Sister tool to [Parker](https://github.com/rizkyandriawan/parker).

## The Story

At the Daily Bugle, there are two photographers: **Peter Parker** and **Eddie Brock**. Same job—capture moments. Different vibes.

| Tool | Domain | Style | Captures |
|------|--------|-------|----------|
| **Parker** 🕷️ | Web Browser | Colorful, friendly | Web UI screenshots |
| **Eddie** 🖤 | Terminal | Dark, raw | CLI/Claude Code sessions |

Parker shoots the web. Eddie captures the void.

## What Eddie Does

Eddie automates the process of documenting Claude Code interactions:

1. Starts Claude Code in a directory
2. Sends prompts/commands
3. Waits for responses
4. Captures terminal screenshots as PNG
5. Generates manifest.json for tooling

## Installation

```bash
# Clone
git clone https://github.com/rizkyandriawan/eddie.git
cd eddie

# Build
go build -o eddie .

# (Optional) Install globally
sudo mv eddie /usr/local/bin/
```

## Quick Start

1. Create a config file:

```yaml
# eddie.yaml
output: ./screenshots
manifest: true

terminal:
  width: 120
  height: 40

sessions:
  - name: hello-claude
    description: "Basic Claude interaction"
    cwd: ~/projects/myapp
    prompts:
      - input: "what files are in this project?"
        wait: 5000
        capture: true
```

2. Run Eddie:

```bash
eddie -c eddie.yaml
```

3. Check output:

```
./screenshots/
├── hello-claude.png
└── manifest.json
```

## Configuration

### Basic Structure

```yaml
output: ./screenshots      # Output directory
manifest: true             # Generate manifest.json

terminal:
  width: 120               # Terminal columns
  height: 40               # Terminal rows

theme:
  name: dark
  background: "#1a1a1a"
  foreground: "#d4d4d4"
  font_size: 14
  padding: 20

sessions:
  - name: session-name
    description: "What this captures"
    cwd: ~/path/to/project
    prompts:
      - input: "your prompt here"
        wait: 5000           # Wait 5 seconds
        capture: true        # Take screenshot
```

### Multi-Turn Conversation

```yaml
sessions:
  - name: refactor-flow
    description: "Multi-step refactoring"
    cwd: ~/projects/myapp
    prompts:
      - input: "analyze main.go"
        wait: 8000
        capture: true
        capture_name: "01-analyze"

      - input: "refactor the error handling"
        wait: 15000
        capture: true
        capture_name: "02-refactor"

      - input: "/diff"
        wait: 3000
        capture: true
        capture_name: "03-diff"
```

### Wait Until Pattern

```yaml
prompts:
  - input: "run tests"
    wait_until: "All tests passed"   # Wait for this text
    timeout: 60000                    # Max 60 seconds
    capture: true
```

### Send Keystrokes

```yaml
prompts:
  - input: "delete all temp files"
    wait: 3000
    capture: true
    capture_name: "approval-prompt"

  - key: "y"                         # Send 'y' to confirm
    wait: 2000
    capture: true
    capture_name: "confirmed"
```

### With Setup Commands

```yaml
sessions:
  - name: feature-review
    cwd: ~/projects/myapp
    setup:
      - "git checkout feature-branch"
      - "npm install"
    prompts:
      - input: "review the changes in this branch"
        wait: 10000
        capture: true
```

## CLI Reference

```
eddie -c <config.yaml> [options]

OPTIONS:
    -c <path>       Path to YAML config file (required)
    -o <path>       Output directory (overrides config)
    --manifest      Generate manifest.json
    --version       Show version
    --help          Show help
```

### Examples

```bash
# Basic usage
eddie -c config.yaml

# Custom output directory
eddie -c config.yaml -o ./docs/cli-screenshots

# With manifest generation
eddie -c config.yaml --manifest
```

## Output

### Directory Structure

```
screenshots/
├── 01-hello.png
├── 02-refactor.png
├── 03-diff.png
└── manifest.json
```

### Manifest Format

```json
{
  "tool": "eddie",
  "version": "1.0.0",
  "target": "claude-code",
  "generated_at": "2024-01-15T10:30:00Z",
  "terminal": {
    "width": 120,
    "height": 40,
    "theme": "dark"
  },
  "sessions": [
    {
      "name": "01-hello",
      "description": "Basic greeting",
      "cwd": "~/projects/myapp",
      "screenshots": [
        {
          "filename": "01-hello.png",
          "prompt": "what files are in this project?",
          "wait_ms": 5000
        }
      ]
    }
  ],
  "summary": {
    "total_sessions": 1,
    "total_screenshots": 1,
    "success": 1,
    "failed": 0
  }
}
```

## Theming

```yaml
theme:
  name: dark                 # Preset name
  background: "#1a1a1a"      # Terminal background
  foreground: "#d4d4d4"      # Default text color
  font_size: 14              # Font size in pixels
  padding: 20                # Image padding

  # Custom ANSI colors (optional)
  colors:
    black: "#000000"
    red: "#ff5555"
    green: "#50fa7b"
    yellow: "#f1fa8c"
    blue: "#bd93f9"
    magenta: "#ff79c6"
    cyan: "#8be9fd"
    white: "#f8f8f2"
```

## Supported Keys

For the `key` field in prompts:

| Key | Description |
|-----|-------------|
| `enter` | Enter/Return |
| `tab` | Tab |
| `escape` | Escape |
| `backspace` | Backspace |
| `up/down/left/right` | Arrow keys |
| `ctrl+c` | Interrupt |
| `ctrl+d` | EOF |
| `y`, `n`, etc. | Single characters |

## Use Cases

### Documenting CLI Tools

```yaml
sessions:
  - name: cli-help
    cwd: ~/projects/mytool
    prompts:
      - input: "explain how to use this CLI"
        wait: 8000
        capture: true
```

### Recording Tutorials

```yaml
sessions:
  - name: tutorial-01
    description: "Getting started tutorial"
    cwd: ~/projects/demo
    prompts:
      - input: "create a new React component called Button"
        wait: 15000
        capture: true
        capture_name: "step-01-create"

      - input: "add click handler with console.log"
        wait: 10000
        capture: true
        capture_name: "step-02-handler"
```

### CI/CD Documentation

```yaml
sessions:
  - name: deploy-flow
    cwd: ~/projects/app
    prompts:
      - input: "show me the deploy script"
        wait: 5000
        capture: true

      - input: "explain what each step does"
        wait: 10000
        capture: true
```

## Parker + Eddie Workflow

Use both tools for complete documentation:

```bash
# Capture web UI screenshots
parker -c web-config.yaml -o ./docs/web --manifest

# Capture CLI/Claude Code sessions
eddie -c cli-config.yaml -o ./docs/cli --manifest

# Both manifests can feed into doc generation
```

## Requirements

- Go 1.22+
- Claude Code CLI installed (`claude` command available)
- A monospace font (DejaVu Sans Mono, Liberation Mono, etc.)

## License

MIT

---

*Parker shoots the web. Eddie captures the void.* 🖤
