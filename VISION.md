# Eddie Vision Document

## The Origin Story: Parker

### Problem

Dokumentasi UI itu pain. Setiap kali ada update, harus:
1. Buka browser
2. Navigate ke page
3. Screenshot manual
4. Crop, rename, organize
5. Update markdown

Multiply by 20 pages. Multiply by setiap release. **Nightmare.**

### Solution: Parker

**Parker** lahir sebagai automated screenshot tool untuk web apps.

```yaml
# parker-config.yaml
urls:
  - url: http://localhost:5173/
    name: 01-homepage
  - url: http://localhost:5173/dashboard
    name: 02-dashboard
```

```bash
parker -c config.yaml -o ./docs/screenshots --manifest
```

Done. 20 screenshots dalam seconds. Consistent naming. Manifest untuk tracking.

### Naming: Peter Parker

Kenapa "Parker"?

Peter Parker adalah fotografer. Dia *captures moments* — literally his job di Daily Bugle. Dan dia punya web powers. **Web. Screenshots. Get it?**

```
Spider-Man → Web → Web Browser → Parker captures the web
```

Parker bekerja dengan baik. Web apps terdokumentasi dengan rapi.

**Tapi ada gap.**

---

## The Gap: CLI Tools

Banyak products punya dua interface:
1. **Web UI** — untuk end users
2. **CLI** — untuk developers, DevOps, power users

Parker handles #1. Tapi CLI documentation? Still manual.

```bash
# Current workflow untuk dokumentasi CLI:
$ myapp --help
# *manually screenshot terminal*
# *crop*
# *save as help.png*
# *repeat for every command*
```

Same pain. Different domain.

---

## Enter Eddie

### The Parallel

Di Daily Bugle, Peter Parker punya rival: **Eddie Brock**.

Same job — photographer. Same goal — capture moments. Tapi *different approach*:

| Aspect | Peter Parker | Eddie Brock |
|--------|--------------|-------------|
| Persona | Friendly, optimistic | Dark, intense |
| Style | Colorful, dynamic | Raw, minimal |
| Domain | Day shift | Night shift |

Dan yang paling perfect: Eddie Brock becomes **Venom**.

Venom's aesthetic:
- **Black and white**
- Raw, unfiltered
- Terminal-like

*Exactly like terminal output.*

### The Tools

| Tool | Domain | Aesthetic | Captures |
|------|--------|-----------|----------|
| **Parker** 🕷️ | Web Browser | Colorful, modern | Web UI screenshots |
| **Eddie** 🖤 | Terminal | Dark, monospace | CLI output screenshots |

```
Parker shoots the web.
Eddie captures the void.
```

---

## Vision: Documentation as Artifacts

### The Philosophy

Screenshots bukan tests. Screenshots adalah **artifacts** — bukti visual bahwa software works dan looks seperti yang diharapkan.

Good documentation punya:
1. **Written explanations** — what and why
2. **Code examples** — how
3. **Visual proof** — screenshots

Parker dan Eddie menghandle #3.

### The Workflow

```
┌─────────────────────────────────────────────────────────┐
│                    YOUR PRODUCT                         │
├──────────────────────┬──────────────────────────────────┤
│      Web UI          │           CLI                    │
├──────────────────────┼──────────────────────────────────┤
│                      │                                  │
│   ┌──────────┐       │       ┌──────────┐              │
│   │  Parker  │       │       │  Eddie   │              │
│   │    🕷️    │       │       │    🖤    │              │
│   └────┬─────┘       │       └────┬─────┘              │
│        │             │            │                    │
│        ▼             │            ▼                    │
│   screenshots/       │       cli-screenshots/          │
│   ├── 01-home.png    │       ├── 01-help.png          │
│   ├── 02-dash.png    │       ├── 02-status.png        │
│   └── manifest.json  │       └── manifest.json        │
│        │             │            │                    │
│        └─────────────┴────────────┘                    │
│                      │                                  │
│                      ▼                                  │
│              📄 DOCUMENTATION                           │
│              (with visual proof)                        │
└─────────────────────────────────────────────────────────┘
```

### Manifest: The Bridge

Both tools generate `manifest.json` — structured metadata about captures.

```json
{
  "tool": "eddie",
  "screenshots": [
    {
      "name": "01-help",
      "command": "myapp --help",
      "description": "Available commands",
      "filename": "01-help.png"
    }
  ]
}
```

Manifests enable:
- **LLM context** — feed to AI for doc generation
- **Diffing** — detect visual changes between versions
- **Automation** — CI/CD integration
- **Linking** — connect screenshots to code/docs

---

## Eddie: Technical Vision

### Core Principles

1. **Single binary** — No runtime dependencies. Download and run.
2. **YAML-driven** — Declarative config, same philosophy as Parker.
3. **Terminal-native** — Proper ANSI color support. Looks like real terminal.
4. **Fast** — Written in Go. Parallel execution.
5. **Minimal** — Does one thing well: capture CLI output as images.

### What Eddie Captures

```yaml
commands:
  # Help text
  - cmd: "myapp --help"
    name: 01-help

  # Status output
  - cmd: "myapp status"
    name: 02-status

  # Colored output
  - cmd: "git diff --color=always"
    name: 03-diff

  # Tables
  - cmd: "docker ps --format 'table {{.Names}}\t{{.Status}}'"
    name: 04-containers

  # TUI snapshots
  - cmd: "htop"
    name: 05-htop
    mode: tui
    delay: 2000

  # Progress/animations
  - cmd: "./deploy.sh"
    name: 06-deploy
    mode: sequence
    interval: 500
    output_format: gif
```

### Output

```
docs/
├── web/                    # From Parker
│   ├── 01-homepage.png
│   ├── 02-dashboard.png
│   └── manifest.json
│
├── cli/                    # From Eddie
│   ├── 01-help.png
│   ├── 02-status.png
│   ├── 03-diff.png
│   └── manifest.json
│
└── README.md               # References both
```

---

## The Taglines

**Parker:** *"Shoots the web."*

**Eddie:** *"We are Eddie."*

(Yes, it's a Venom reference. "We are Venom." Eddie speaks for the terminal.)

---

## Success Criteria

Eddie is successful when:

1. **Zero friction** — `eddie -c config.yaml` just works
2. **Beautiful output** — Screenshots look like actual terminal, not garbage
3. **Consistent with Parker** — Same config style, same manifest format
4. **Fast** — 100 commands in under 10 seconds
5. **Reliable** — Handles failures gracefully, partial success is OK

---

## Roadmap

### Phase 1: MVP
- [x] Vision & spec
- [ ] Go project setup
- [ ] YAML config parser
- [ ] Command runner (instant mode)
- [ ] PNG renderer (ANSI color support)
- [ ] Basic CLI

### Phase 2: Polish
- [ ] Manifest generation
- [ ] HTML gallery
- [ ] Theming (dark/light/custom)
- [ ] Better error handling

### Phase 3: Advanced
- [ ] TUI mode (htop, vim, etc.)
- [ ] Sequence mode (GIF output)
- [ ] Interactive mode (expect-like)
- [ ] CI/CD integration

---

## Why This Matters

Documentation is the difference between:
- "What does this CLI do?" vs "Here's exactly what it does" *[screenshot]*
- "How do I use this?" vs "Run this command, you'll see this" *[screenshot]*
- "Is this working?" vs "Yes, look" *[screenshot]*

Screenshots are proof. Automated screenshots are **sustainable proof**.

Parker handles the web. Eddie handles the terminal.

Together, they document everything.

---

*"We are Eddie."* 🖤
