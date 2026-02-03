# Eddie 🖤

## The Story

### Naming: Eddie Brock

Eddie Brock adalah fotografer di Daily Bugle — tempat yang sama dengan Peter Parker. Mereka rival, tapi punya job yang sama: **capture moments**.

Bedanya:
- **Parker** = friendly, colorful, web-based
- **Eddie** = dark, raw, terminal-based

Dan yang paling perfect: Eddie jadi **Venom** — karakter yang identik dengan warna **hitam putih**. Persis kayak terminal output.

### The Parallel

| Aspect | Parker 🕷️ | Eddie 🖤 |
|--------|-----------|----------|
| Domain | Web (browser) | Terminal (shell) |
| Engine | Playwright | termshot / native |
| Output | Colorful screenshots | B&W terminal captures |
| Config | YAML-based | YAML-based (same style) |
| Vibe | Friendly neighborhood | Dark, minimal, raw |
| Tagline | "Shoots the web" | "We are Eddie" |

### Why This Matters

Parker dan Eddie adalah **sister tools**:
- Sama-sama capture screenshots untuk dokumentasi
- Sama-sama YAML config
- Sama-sama output manifest untuk LLM
- Beda domain: web vs terminal

Kalau lo punya web app + CLI tool, lo pakai **dua-duanya**.

---

## What Eddie Does

Eddie adalah CLI screenshot tool untuk **terminal commands**.

### Core Flow

1. User defines commands di YAML
2. Eddie runs each command
3. Captures the terminal output as PNG
4. Generates manifest.json (sama kayak Parker)
5. Optional: HTML gallery

### Example Usage

```bash
# Basic
eddie -c commands.yaml

# With output dir
eddie -c commands.yaml -o ./docs/cli-screenshots

# With manifest
eddie -c commands.yaml --manifest

# Full output
eddie -c commands.yaml --manifest --html
```

---

## Configuration Format

### Basic

```yaml
commands:
  - "ls -la"
  - "docker ps"
  - "git status"
```

### With Metadata

```yaml
commands:
  - cmd: "myapp --help"
    name: 01-help
    description: "Shows all available commands"

  - cmd: "myapp status"
    name: 02-status
    description: "Current application status"

  - cmd: "myapp list --format=table"
    name: 03-list
    description: "List all items in table format"
```

### With Options

```yaml
commands:
  - cmd: "htop"
    name: system-monitor
    description: "System resource monitor"
    interactive: true    # untuk TUI apps
    wait: 2000           # wait 2s before capture

  - cmd: "cat /etc/passwd | head -10"
    name: passwd-sample
    description: "First 10 lines of passwd"

  - cmd: "tree -L 2"
    name: project-structure
    description: "Project directory structure"
    cwd: "/home/user/myproject"  # run dari directory tertentu
```

### Environment Variables

```yaml
env:
  NODE_ENV: production
  DEBUG: "app:*"

commands:
  - cmd: "node app.js --version"
    name: version
```

### Theming (Stretch Goal)

```yaml
theme:
  background: "#1e1e1e"    # terminal background
  foreground: "#d4d4d4"    # text color
  font: "JetBrains Mono"
  font_size: 14
  padding: 20

commands:
  - "echo hello"
```

---

## CLI Interface

```bash
eddie -c <config.yaml> [options]
```

### Flags

| Flag | Description | Default |
|------|-------------|---------|
| `-c, --config` | YAML config file | (required) |
| `-o, --output` | Output directory | `./screenshots` |
| `--manifest` | Generate manifest.json | `false` |
| `--html` | Generate HTML gallery | `false` |
| `--width` | Terminal width (columns) | `120` |
| `--height` | Terminal height (rows) | `40` |
| `--theme` | Color theme | `dark` |
| `--font` | Font family | `monospace` |
| `--font-size` | Font size in px | `14` |

### Examples

```bash
# Basic capture
eddie -c commands.yaml

# Custom dimensions
eddie -c commands.yaml --width 80 --height 24

# Light theme
eddie -c commands.yaml --theme light

# Full documentation output
eddie -c commands.yaml -o ./docs/cli --manifest --html
```

---

## Output

### Directory Structure

```
output/
├── 01-help.png
├── 02-status.png
├── 03-list.png
├── manifest.json
└── gallery.html
```

### Manifest Format

```json
{
  "tool": "eddie",
  "version": "1.0.0",
  "generated_at": "2024-01-15T10:30:00Z",
  "terminal": {
    "width": 120,
    "height": 40,
    "theme": "dark"
  },
  "screenshots": [
    {
      "command": "myapp --help",
      "name": "01-help",
      "filename": "01-help.png",
      "description": "Shows all available commands",
      "exit_code": 0,
      "duration_ms": 150,
      "output_lines": 25,
      "hash": "a1b2c3d4"
    },
    {
      "command": "myapp status",
      "name": "02-status",
      "filename": "02-status.png",
      "description": "Current application status",
      "exit_code": 0,
      "duration_ms": 89,
      "output_lines": 12,
      "hash": "e5f6g7h8"
    }
  ],
  "summary": {
    "total": 3,
    "success": 3,
    "failed": 0
  }
}
```

### HTML Gallery

Same style as Parker:
- Grid of thumbnails
- Click to fullscreen
- Shows command, description, exit code
- Filter by status (success/failed)

---

## Technical Approach

### Option 1: Using `termshot` (Recommended)

[termshot](https://github.com/homeport/termshot) is a Go tool that:
- Runs a command
- Captures output
- Renders as PNG with terminal styling

```bash
termshot -f output.png -- ls -la
```

**Pros:**
- Mature, well-tested
- Good terminal rendering
- Supports colors (ANSI)

**Cons:**
- External dependency (Go binary)
- Limited customization

### Option 2: Using `rich` + `PIL` (Pure Python)

1. Run command, capture output
2. Use `rich` to render with syntax highlighting
3. Use `PIL` to render to image

```python
from rich.console import Console
from rich.terminal_theme import MONOKAI
from PIL import Image

console = Console(record=True, width=120)
# ... capture output ...
console.save_svg("output.svg")
# Convert SVG to PNG
```

**Pros:**
- Pure Python
- Full control over rendering
- Rich syntax highlighting

**Cons:**
- More complex
- SVG → PNG conversion needed

### Option 3: Using `asciinema` + `svg-term`

1. Record command with asciinema
2. Convert to SVG with svg-term
3. Convert to PNG

**Pros:**
- Can capture interactive/animated content

**Cons:**
- Multiple steps
- Heavier dependencies

### Recommendation

**Start with Option 2 (rich + PIL)** karena:
- Pure Python = consistent dengan Parker
- Full control = bisa customize theme
- `rich` sudah handle ANSI colors dengan baik

Fallback ke termshot kalau rendering terlalu complex.

---

## Implementation Plan

### Phase 1: MVP 🎯

1. **Config parser** - Parse YAML config
2. **Command runner** - Execute commands, capture stdout/stderr
3. **Renderer** - Convert output to PNG using rich
4. **CLI** - Basic argparse interface

Deliverables:
- `eddie.py` - Main script
- Basic YAML config support
- PNG output

### Phase 2: Polish ✨

1. **Manifest generation** - JSON output with metadata
2. **HTML gallery** - Interactive viewer
3. **Theming** - Light/dark themes, custom colors
4. **Error handling** - Graceful failures, exit codes

### Phase 3: Advanced 🚀

1. **Interactive apps** - Support for TUI (htop, vim, etc.)
2. **Animated captures** - GIF output for dynamic content
3. **Diff detection** - Compare outputs between runs
4. **CI integration** - GitHub Actions, etc.

---

## Exit Codes

| Code | Meaning |
|------|---------|
| `0` | All commands captured successfully |
| `1` | Some commands failed (partial success) |
| `2` | Critical error (config invalid, etc.) |

---

## File Structure

```
~/product/eddie/
├── eddie.py           # Main script
├── requirements.txt   # Dependencies
├── README.md          # Documentation
├── CONFIG.md          # Config reference
├── examples/
│   ├── basic.yaml
│   ├── with-metadata.yaml
│   └── advanced.yaml
└── tests/
    └── test_eddie.py
```

---

## Dependencies

```txt
pyyaml>=6.0
rich>=13.0
Pillow>=10.0
```

Optional:
```txt
cairosvg>=2.7      # for SVG to PNG conversion
```

---

## Relationship to Parker

Eddie dan Parker adalah **sister tools**. Mereka share:

1. **Config philosophy** - YAML-based, declarative
2. **Output format** - manifest.json compatible
3. **CLI style** - Similar flags and UX
4. **Documentation focus** - Screenshots as artifacts, not tests

Idealnya, user bisa pakai keduanya:

```bash
# Capture web UI
parker -c web-urls.yaml -o ./docs/web --manifest

# Capture CLI
eddie -c cli-commands.yaml -o ./docs/cli --manifest

# Both manifests feed into doc generation
```

---

## Current Status

📍 **Status: SPEC COMPLETE, READY FOR IMPLEMENTATION**

### Done ✅
- [x] Naming & story
- [x] Feature spec
- [x] Config format design
- [x] CLI interface design
- [x] Output format design
- [x] Technical approach decided

### Next 🔜
- [ ] Create `eddie.py` MVP
- [ ] Test with basic commands
- [ ] Add manifest generation
- [ ] Add HTML gallery
- [ ] Write README.md

---

## Notes for Implementation

1. **Keep it simple** - Start with basic command capture, iterate
2. **Match Parker's UX** - Similar flags, similar output structure
3. **Handle failures gracefully** - Commands might fail, that's OK
4. **ANSI colors** - Terminal output often has colors, preserve them
5. **Cross-platform** - Should work on Linux, macOS, (Windows later)

---

*"We are Eddie."* 🖤
