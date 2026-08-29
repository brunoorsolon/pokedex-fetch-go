---
title: pokedex-fetch-go
type: PRD
status: draft
created: 2026-04-29
---

# pokedex-fetch-go — Product Requirements Document

> **Repository**: https://github.com/brunoorsolon/pokedex-fetch-go
> **License**: MIT
> **Binary name**: `pokedex-fetch-go`

## 1. Executive Summary

**pokedex-fetch-go** is a single-binary CLI tool written in Go that displays random Pokemon sprites in the terminal and integrates with [fastfetch](https://github.com/fastfetch-cli/fastfetch) for system info display. It replaces the Python-based `pokemon-colorscripts` + `poketerm` stack with a zero-dependency, self-contained binary.

Every time the user opens a terminal, a random Pokemon sprite appears alongside their system info (via fastfetch). That encounter is "caught" and tracked in a persistent Pokedex. Rare shiny variants appear at configurable odds. The tool includes a full interactive Pokedex viewer, trainer profile with XP/leveling, achievements system, and generation-based collection tracking — all without requiring Python, pip, or any runtime dependency.

The core value proposition: **everything poketerm does, in a single binary you download and run.**

## 2. Target Users

### Primary: Terminal Enthusiasts / Ricers
- Customize their shell prompt, fetch tool, and terminal appearance
- Comfortable with CLI tools, dotfiles, shell config
- Run fastfetch/neofetch on terminal open
- Want minimal dependencies — every pip install or runtime requirement is friction

### Secondary: Pokemon Fans Who Code
- Enjoy the gamification of terminal usage
- Motivated by collection progress, achievements, streaks
- Share screenshots of their terminal setup

### Pain Points with Current Solutions
- `pokemon-colorscripts` requires Python 3 runtime (~30-50 MB)
- `poketerm` adds pip dependencies (`readchar`), zsh-only shell integration, depends on `pokemon-colorscripts` as external tool, and requires Homebrew
- Installation involves `sudo`, modifying `.zshrc`, and multiple moving parts
- Updating requires `git pull` + running an installer with migration paths
- No support for fastfetch (hardcoded to neofetch/hyfetch, which are deprecated)

## 3. MVP Scope

### In Scope

**Core Display**
- Display a random Pokemon sprite (ANSI art) in the terminal
- Display a specific Pokemon by name or Pokedex number
- Support small and large sprite variants
- Support alternate forms (mega, gmax, etc.)
- Shiny variant with configurable probability (default 1/4096)
- Boosted shiny rate (1/10) for completed generations
- Print Pokemon name with sprite (optional, can be disabled)
- Output compatible with fastfetch `--logo` / `--logo-type file-raw`

**Pokedex / Collection**
- Persistent Pokedex tracking (JSON file)
- Track normal and shiny catches separately with catch count
- Per-generation progress (Gen 1-8, 905 Pokemon)
- Interactive TUI Pokedex viewer with arrow key navigation, pagination
- Detail view for caught Pokemon (types, stats, height, weight, sprite)
- Generation filtering (`--gen 1`, `--gen 2-5`)

**Trainer Profile**
- Trainer name, home region, home town, favourite Pokemon
- XP system (10 XP normal catch, 150 XP shiny catch)
- Level progression (level 1-50+)
- Rank titles (Beginner Trainer through Grand Champion)
- Catch streaks (daily streak tracking, current + longest)
- Daily catch counter

**Achievements**
- Pokedex milestones (1, 50, 100, 905 caught)
- Catch streaks (3, 7, 30, 100, 365 days)
- Shiny milestones (1, 5, 100 shinies)
- Daily catch milestones (5, 10, 50 in a day)
- Duplicate catch milestones (10x, 25x same Pokemon)
- Level milestones (10, 25, 50)
- XP milestones (1k, 10k, 50k, 100k)
- Generation completion (per-gen + all-gens)
- Special collections (all legendaries, all mythicals, all starters, all eeveelutions, all fossils, all pseudo-legendaries, kanto birds, weather trio)
- Hidden achievements (discoverable)
- Hall of Fame display

**Configuration**
- External TOML config file (`~/.config/pokedex-fetch-go/config.toml`)
- Configurable: shiny rate, sprite size, generation pool, favorites, exclusions
- XDG-compliant paths

**Shell Integration**
- Works with bash, zsh, fish (not locked to zsh)
- Simple one-liner to add to shell rc file
- fastfetch integration via process substitution or temp file

### Out of Scope

- Networking / PokeAPI calls (all data embedded or generated at build time)
- Gen 9 (Paldea) — sprites not available in pokemon-colorscripts upstream
- GUI or web interface
- Multiplayer / leaderboards
- Package manager distribution (Homebrew, AUR, etc.) — post-MVP
- Windows support — post-MVP (focus on Linux/macOS)
- Animated sprites
- Custom/user-provided sprites

## 4. User Stories

1. **As a terminal user**, I want to see a random Pokemon sprite alongside my system info when I open a terminal, so that my shell feels personalized and fun.

2. **As a collector**, I want every Pokemon I see to be tracked in a persistent Pokedex, so that I'm motivated to keep opening terminals and filling my collection.

3. **As a Pokemon fan**, I want a rare chance of encountering a shiny variant, so that there's an element of excitement and surprise in each terminal session.

4. **As a completionist**, I want to view my per-generation catch progress in an interactive Pokedex viewer, so that I can track what I'm missing and see how close I am to completing each generation.

5. **As a power user**, I want to configure shiny rates, generation pools, and sprite sizes via a config file, so that I can customize the experience to my preference.

6. **As a minimal-dependency advocate**, I want to install a single binary with no runtime dependencies, so that I don't need Python, pip, or any package manager to use this tool.

7. **As a gamer**, I want XP, levels, ranks, streaks, and achievements, so that using my terminal feels like a progression system.

8. **As a fastfetch user**, I want the tool to integrate seamlessly with fastfetch's logo system, so that I see the Pokemon sprite as my fetch logo without any wrapper hacks.

## 5. Architecture & Patterns

### High-Level Approach

Single Go binary with all sprites and metadata embedded at compile time via `//go:embed`. State persisted to JSON files in XDG config directory. Interactive TUI for Pokedex viewer. Standard CLI for catch/display operations.

### Directory Structure

```
pokedex-fetch-go/
├── cmd/
│   └── pokedex-fetch-go/
│       └── main.go              # Entry point, binary name derived from directory
├── internal/
│   ├── cli/                     # CLI command definitions (cobra)
│   │   ├── root.go
│   │   ├── catch.go             # Default: random catch + display
│   │   ├── show.go              # Show specific Pokemon by name/number
│   │   ├── pokedex.go           # Interactive Pokedex TUI
│   │   ├── trainer.go           # Trainer profile view/edit
│   │   ├── achievements.go      # Achievement display
│   │   └── config.go            # Config management
│   ├── pokemon/                 # Pokemon data + sprite access
│   │   ├── data.go              # Embedded pokemon.json, generation lists
│   │   ├── sprites.go           # Embedded sprite files, lookup
│   │   └── types.go             # Pokemon struct, generation ranges
│   ├── state/                   # Persistent state management
│   │   ├── pokedex.go           # Pokedex read/write (JSON)
│   │   ├── trainer.go           # Trainer profile read/write
│   │   └── paths.go             # XDG path resolution
│   ├── achievement/             # Achievement checking logic
│   │   ├── checker.go           # All achievement check functions
│   │   └── definitions.go       # Achievement + special collection definitions
│   ├── tui/                     # Terminal UI components
│   │   ├── pokedex.go           # Pokedex list view
│   │   ├── detail.go            # Pokemon detail view
│   │   ├── trainer.go           # Trainer profile view
│   │   ├── achievements.go      # Hall of fame view
│   │   └── common.go            # Shared rendering (box drawing, colors, input)
│   └── config/                  # Configuration
│       └── config.go            # TOML parsing, defaults
├── data/                        # Build-time data (embedded into binary)
│   ├── sprites/
│   │   ├── small/
│   │   │   ├── regular/         # 1329 sprite files
│   │   │   └── shiny/           # 1329 sprite files
│   │   └── large/
│   │       ├── regular/
│   │       └── shiny/
│   ├── pokemon.json             # Pokemon metadata (name, forms, types, stats)
│   └── gen/
│       ├── gen1.txt
│       └── ...gen8.txt
├── config.example.toml          # Example config
├── go.mod
├── go.sum
└── README.md
```

### Design Patterns

- **Embed compressed, all variants**: All sprites gzip-compressed before embedding via `embed.FS`. Decompression cost is ~5-11 microseconds per sprite (benchmarked) — 1000x smaller than terminal rendering time. This allows shipping all 4 variants (small/large × regular/shiny) in ~4 MB of sprite data instead of 53 MB raw. Total binary target: ~12 MB.
- **Atomic writes**: State files (pokedex.json, trainer.json) written via write-to-temp + rename to prevent corruption on crash.
- **Separation of concerns**: CLI layer (cobra commands) -> business logic (internal packages) -> state (JSON files). TUI is its own layer that reads state but doesn't modify it directly.
- **No network calls**: All Pokemon data (types, stats, height, weight) baked into `pokemon.json` at build time. No PokeAPI dependency at runtime.

### Data Flow: Terminal Open

```
shell rc → pokedex-fetch-go catch [--gen 1-8] → roll random Pokemon
                                     → roll shiny check
                                     → update pokedex.json (atomic write)
                                     → update trainer.json (XP, streak, achievements)
                                     → print sprite to stdout
                                     → fastfetch reads sprite via --logo
```

### Data Flow: Pokedex Viewer

```
pokedex-fetch-go pokedex [gen] → read pokedex.json + gen lists
                     → render interactive TUI
                     → arrow keys navigate, Enter for detail
                     → detail view shows embedded stats + sprite
                     → q to quit
```

## 6. Technology Stack

| Component | Choice | Justification |
|-----------|--------|---------------|
| Language | Go 1.22+ | Single binary, embed support, fast startup, zero runtime deps |
| CLI framework | `cobra` | Industry standard Go CLI, subcommands, help generation |
| Config parsing | `BurntSushi/toml` | Lightweight, well-maintained TOML parser |
| TUI framework | `charmbracelet/bubbletea` | Full TUI framework — handles resize, input, rendering loop, component model. Industry standard for Go TUIs (lazygit, glow). Avoids brittle manual ANSI rendering for the Pokedex/trainer/achievements screens. |
| Data format | JSON (state), TOML (config) | JSON stdlib for state; TOML for human-editable config |
| Sprite storage | `embed.FS` | Compiled into binary, no runtime file paths |
| Sprite source | pokemon-colorscripts (MIT) | 1329 sprites × 4 variants (small/large × regular/shiny) |

### Dependencies (Go modules)

| Module | Purpose | Size Impact |
|--------|---------|-------------|
| `github.com/spf13/cobra` | CLI subcommands + flags | Minimal |
| `github.com/BurntSushi/toml` | Config file parsing | Minimal |
| `github.com/charmbracelet/bubbletea` | TUI framework for interactive screens | Moderate (~10 transitive deps, all from Charm ecosystem, all compiled into binary) |
| `github.com/charmbracelet/lipgloss` | TUI styling (comes with bubbletea) | Included with bubbletea |

**Total external deps: 3 direct.** ~10 transitive from Charm ecosystem. All compile into the binary — zero runtime deps for the user.

### Build-Time Data Preparation

A one-time script (not shipped) will:
1. Copy sprites from pokemon-colorscripts into `data/sprites/`
2. Enrich `pokemon.json` with types, stats, height, weight from PokeAPI (cached)
3. Generate generation list files from pokemon.json

This script runs during development, not at user install time.

## 7. Security & Configuration

### Security

- **No network access at runtime** — all data is embedded. No DNS, no HTTP, no sockets.
- **No elevated permissions** — no `sudo` required for install or operation.
- **Atomic state writes** — write to temp file, then rename. Prevents corruption.
- **No shell injection** — tool is called directly, not via eval or shell expansion.
- **State integrity** — optional SHA256 hash verification of pokedex file (detect manual edits), matching poketerm's approach.

### Configuration

Config file: `~/.config/pokedex-fetch-go/config.toml`

```toml
# Shiny encounter rate (1 in N). Default: 4096
shiny_rate = 4096

# Boosted shiny rate for completed generations (1 in N). Default: 10
shiny_rate_boosted = 10

# Sprite size: "small" or "large". Default: "small"
sprite_size = "small"

# Show Pokemon name with sprite. Default: true
show_name = true

# Generations to pick from. Default: [1,2,3,4,5,6,7,8]
generations = [1, 2, 3, 4, 5, 6, 7, 8]

# XP awards
xp_normal = 10
xp_shiny = 150

# Pokemon to never show (by name)
excluded = []
```

### State Files

| File | Path | Purpose |
|------|------|---------|
| Pokedex | `~/.config/pokedex-fetch-go/pokedex.json` | Catch records |
| Trainer | `~/.config/pokedex-fetch-go/trainer.json` | Profile, XP, streaks, achievements |
| Config | `~/.config/pokedex-fetch-go/config.toml` | User preferences |

All paths follow XDG Base Directory spec (`$XDG_CONFIG_HOME` with `~/.config` fallback).

## 8. Success Criteria

### Functional Requirements

- [ ] `pokedex-fetch-go` (no args) displays a random Pokemon sprite and registers the catch
- [ ] `pokedex-fetch-go show --name pikachu` displays a specific Pokemon
- [ ] `pokedex-fetch-go show --name pikachu --shiny` forces shiny variant
- [ ] `pokedex-fetch-go catch --gen 1` restricts to Gen 1 Pokemon
- [ ] `pokedex-fetch-go pokedex` opens interactive Pokedex TUI
- [ ] `pokedex-fetch-go pokedex 3` opens Gen 3 Pokedex
- [ ] `pokedex-fetch-go trainer` shows trainer profile
- [ ] `pokedex-fetch-go achievements` shows Hall of Fame
- [ ] `pokedex-fetch-go list` prints all Pokemon names
- [ ] `pokedex-fetch-go config init` creates default config file
- [ ] Shiny encounters occur at configured rate
- [ ] Completed generation boosts shiny rate
- [ ] XP, levels, and ranks progress correctly
- [ ] Achievements unlock at correct thresholds
- [ ] Streaks track daily activity accurately
- [ ] Special collections (legendaries, starters, etc.) tracked
- [ ] `fastfetch --logo-type file-raw --logo <(pokedex-fetch-go catch --raw)` works
- [ ] State survives across sessions (persistent JSON)
- [ ] Config file changes take effect without rebuild
- [ ] Works on bash, zsh, and fish

### Quality Indicators

- Binary size < 15 MB (all 4 sprite variants, gzip-compressed embed)
- Startup to sprite output < 20ms
- Zero runtime dependencies
- No network calls
- Clean `go vet` and `golangci-lint`

### UX Goals

- Single binary install: download + chmod + move to PATH
- Shell integration is one line in rc file
- Interactive Pokedex feels responsive (no input lag)
- Config file is self-documenting with comments

## 9. Implementation Phases

### Phase 1 — Core: Catch & Display
**Goal**: Random Pokemon sprite in terminal, fastfetch integration working.

- [ ] Project scaffolding (go.mod, cobra setup, directory structure)
- [ ] Data preparation script: copy sprites, enrich pokemon.json with stats from PokeAPI
- [ ] Embed sprites and metadata via `embed.FS`
- [ ] `pokedex-fetch-go` default command: random Pokemon, print sprite to stdout
- [ ] `pokedex-fetch-go show --name <name>` / `pokedex-fetch-go show --number <n>`
- [ ] `--shiny`, `--big`, `--no-title` flags
- [ ] `--gen` flag for generation filtering
- [ ] `--raw` flag (sprite only, no name — for fastfetch piping)
- [ ] `pokedex-fetch-go list` command
- [ ] TOML config file loading with defaults
- [ ] Fastfetch integration tested and documented

**Validation**: `pokedex-fetch-go | cat` outputs valid ANSI art. `fastfetch --logo-type file-raw --logo <(pokedex-fetch-go catch --raw)` displays sprite alongside system info.

### Phase 2 — State: Pokedex & Catches
**Goal**: Every encounter is tracked. Collection progress visible.

- [ ] Pokedex JSON state file: read, write, atomic save
- [ ] Catch registration on every `pokedex-fetch-go` / `pokedex-fetch-go catch` invocation
- [ ] Normal + shiny tracking with counts
- [ ] Shiny roll logic (configurable rate, boosted for completed gens)
- [ ] Interactive Pokedex TUI: list view with pagination, arrow keys, generation tabs
- [ ] Detail view: sprite + types + stats + height/weight for caught Pokemon
- [ ] `pokedex-fetch-go pokedex [gen]` command

**Validation**: Open 20+ terminals. Run `pokedex-fetch-go pokedex` and verify catches appear. Verify shiny tracking.

### Phase 3 — Progression: Trainer, XP, Achievements
**Goal**: Full gamification layer.

- [ ] Trainer profile creation (name, region, town, favourites)
- [ ] Trainer JSON state file
- [ ] XP system: earn on catch, level up formula
- [ ] Rank titles based on level
- [ ] Daily catch counter
- [ ] Streak tracking (current, longest)
- [ ] Achievement definitions (all milestones from poketerm)
- [ ] Special collection tracking (legendaries, starters, eeveelutions, etc.)
- [ ] Achievement checker runs on each catch
- [ ] `pokedex-fetch-go trainer` command with TUI display
- [ ] `pokedex-fetch-go achievements` command with Hall of Fame TUI
- [ ] Trainer profile editing

**Validation**: Verify XP increments, level ups, rank changes. Trigger several achievements. Check streak persists across days.

### Phase 4 — Polish & Distribution
**Goal**: Ready for public use.

- [ ] Shell integration snippets for bash, zsh, fish
- [ ] `pokedex-fetch-go setup` command that prints the appropriate rc snippet
- [ ] Example config file with full documentation
- [ ] Man page or `pokedex-fetch-go help` comprehensive output
- [ ] README with screenshots, install instructions
- [ ] Cross-compilation (Linux amd64/arm64, macOS amd64/arm64)
- [ ] Release binaries via GitHub Releases
- [ ] Migration tool for poketerm users (import existing pokedex.txt)
- [ ] Optional: Homebrew formula, AUR PKGBUILD

**Validation**: Fresh install on clean Linux and macOS systems. Import existing poketerm data and verify continuity.

## 10. Risks & Mitigations

### 1. Binary Size from Embedded Sprites
**Risk**: Embedding all 4 sprite variants raw would be ~53 MB.
**Mitigation**: Gzip-compress sprites before embedding. ANSI text compresses at ~92% ratio, reducing 53 MB to ~4 MB. Decompression cost is ~5-11 microseconds per sprite (benchmarked) — invisible compared to terminal rendering. All 4 variants ship in every binary. Target: ~12 MB total.

### 2. Sprite Data Licensing
**Risk**: pokemon-colorscripts sprites are MIT licensed, but Pokemon IP belongs to The Pokemon Company.
**Mitigation**: Same legal posture as pokemon-colorscripts and poketerm — fan project, non-commercial, credit The Pokemon Company. Standard in this ecosystem. Not a blocker.

### 3. Upstream Sprite Availability for New Generations
**Risk**: Gen 9+ sprites not available in pokemon-colorscripts yet.
**Mitigation**: Design data pipeline to be re-runnable. When upstream adds Gen 9, re-run prep script and rebuild. Generation definitions are data, not code.

### 4. Terminal Compatibility (ANSI rendering)
**Risk**: Sprites use ANSI 256-color or truecolor escapes. Some terminals render differently.
**Mitigation**: pokemon-colorscripts sprites already work across most modern terminals. Same escape codes, same rendering. Test on: kitty, alacritty, wezterm, iTerm2, GNOME Terminal.

### 5. State File Corruption
**Risk**: If the process is killed mid-write, JSON could be truncated.
**Mitigation**: Atomic writes (write temp file, fsync, rename). This is a solved problem and costs ~3 lines of Go. Same pattern used by every serious CLI tool.

## Future Considerations (Post-MVP)

- **Import/export**: Import poketerm pokedex.txt, export collection stats
- **Search**: Fuzzy search Pokemon by name in Pokedex TUI
- **Themes**: Configurable TUI colors
- **Windows support**: Should work with Windows Terminal (truecolor)
- **Gen 9+**: When sprites become available
- **Notification on rare events**: Optional desktop notification on shiny catch
- **Statistics**: Catch rate per Pokemon, luck factor, time-based analytics
- **Alternative fetch tools**: Direct integration with macchina, pfetch, etc.

## Appendix

### Sprite Data Dimensions

Source: `pokemon-colorscripts` (MIT license)

| Category | Files | Total Size (uncompressed) |
|----------|-------|--------------------------|
| small/regular | 1329 | 14 MB |
| small/shiny | 1329 | 14 MB |
| large/regular | 1329 | 13 MB |
| large/shiny | 1329 | 13 MB |
| **Total** | **5316** | **53 MB** |

Average sprite size: ~3-5 KB (small), ~3-4 KB (large).

### Pokemon Data Structure (pokemon.json)

```json
{
  "name": "bulbasaur",
  "number": 1,
  "generation": 1,
  "forms": ["regular"],
  "types": ["grass", "poison"],
  "height": 7,
  "weight": 69,
  "stats": {
    "hp": 45,
    "attack": 49,
    "defense": 49,
    "sp_attack": 65,
    "sp_defense": 65,
    "speed": 45
  },
  "is_legendary": false,
  "is_mythical": false
}
```

### Pokedex State Structure (pokedex.json)

```json
{
  "pokemon": {
    "bulbasaur": {
      "normal": { "count": 3, "first_caught": "2026-04-29" },
      "shiny": { "count": 0, "first_caught": null }
    }
  },
  "version": 1
}
```

### Trainer State Structure (trainer.json)

```json
{
  "name": "Ash",
  "region": "Kanto",
  "hometown": "Pallet Town",
  "created": "2026-04-29",
  "favourites": [25, 6, 150, 149, 131, 143],
  "level": 12,
  "xp": 340,
  "total_xp": 6340,
  "streak": {
    "current": 7,
    "longest": 23,
    "last_active_date": "2026-04-29"
  },
  "daily_catch": {
    "date": "2026-04-29",
    "count": 4
  },
  "hall_of_fame": [
    {
      "name": "Just getting started",
      "event": "Caught your first Pokemon",
      "category": "pokedex",
      "date": "2026-04-15"
    }
  ],
  "version": 1
}
```

### Binary Name & Aliasing

The binary is named `pokedex-fetch-go` to be explicit about what it does and avoid clashing with other tools. The README should document a recommended alias for convenience:

```bash
# Optional: add to your shell rc for a shorter command
alias pfg='pokedex-fetch-go'
```

### Shell Integration Examples

**bash / zsh**:
```bash
# ~/.bashrc or ~/.zshrc
pokedex-fetch-go catch --raw | fastfetch --logo-type file-raw --logo -
```

**fish**:
```fish
# ~/.config/fish/config.fish
pokedex-fetch-go catch --raw | fastfetch --logo-type file-raw --logo -
```

### Fastfetch Integration

fastfetch supports reading a logo from stdin or a file:
```bash
# From process substitution
fastfetch --logo-type file-raw --logo <(pokedex-fetch-go catch --raw)

# From pipe (stdin)
pokedex-fetch-go catch --raw | fastfetch --logo-type file-raw --logo -

# From temp file (most compatible)
pokedex-fetch-go catch --raw > /tmp/pokedex-fetch-go-sprite && fastfetch --logo-type file-raw --logo /tmp/pokedex-fetch-go-sprite
```

The `--raw` flag outputs only the ANSI sprite with no name header, suitable for logo use.

### Credits

- Pokemon designs, names, and branding: [The Pokemon Company](https://www.pokemon.com/)
- Sprite source: [pokemon-colorscripts](https://gitlab.com/phoneybadger/pokemon-colorscripts) (MIT)
- Inspired by: [poketerm](https://github.com/poketerm) for the Pokedex/collection concept
