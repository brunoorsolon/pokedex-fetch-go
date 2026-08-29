# Feature: pokedex-fetch-go

Read the referenced codebase files before implementing. Pay attention to naming of existing utils, types, and models — import from the right files.

## Feature Description

A single-binary Go CLI tool that displays random Pokemon sprites (ANSI art) in the terminal, tracks catches in a persistent Pokedex, and integrates with fastfetch for system info display. Replaces the Python-based `pokemon-colorscripts` + `poketerm` stack with zero runtime dependencies.

Full feature set: random catch on terminal open, shiny variants, interactive Pokedex TUI (bubbletea), trainer profile with XP/leveling/ranks, achievement system, streak tracking, generation-based collection progress, and TOML config.

## User Story

As a terminal enthusiast who uses fastfetch,
I want a single-binary tool that displays Pokemon sprites, tracks my collection, and gamifies terminal usage,
So that I don't need Python/pip/multiple tools and get a better experience with less friction.

## Metadata

- **Type**: New Capability (greenfield project)
- **Complexity**: High
- **Systems affected**: New standalone project — no existing codebase
- **Dependencies**: cobra, BurntSushi/toml, charmbracelet/bubbletea + lipgloss

---

## Context References

### Source Data — READ BEFORE IMPLEMENTING
- `/workspace/pokemon-colorscripts-main/pokemon.json` — Pokemon names + alternate forms (905 entries). Must be enriched with types, stats, height, weight, generation, legendary/mythical flags at build time.
- `/workspace/pokemon-colorscripts-main/colorscripts/` — 5316 sprite files across 4 directories (small/large × regular/shiny). Each 1329 files. Files are ANSI truecolor (RGB) with Unicode half-block characters (▄▀). Average ~3-5 KB.
- `/workspace/poketerm-main/files/gen_files/gen{1..8}_list.txt` — Canonical generation lists (905 total). These define the generation boundaries.

### Reference Implementations — READ FOR BEHAVIOR PARITY
- `/workspace/pokemon-colorscripts-main/pokemon-colorscripts.py` — Original Python sprite display tool. Shows how sprites are looked up by name, how shiny/forms work, how generation ranges map.
- `/workspace/poketerm-main/files/0.0.6/poketerm` (zsh) — Shell integration: catch flow, shiny roll (1/4096 base, 1/10 for completed gen), XP award, pokedex update + sort, neofetch display.
- `/workspace/poketerm-main/files/0.0.6/pokedex.py` — Interactive pokedex viewer: list view, detail view, arrow key navigation, pagination.
- `/workspace/poketerm-main/files/0.0.6/achievements.py` — Achievement definitions, special collections (legendaries, mythicals, starters, eeveelutions, fossils, pseudo-legendaries, kanto birds, weather trio), checker logic.
- `/workspace/poketerm-main/files/0.0.6/trainer.py` — Trainer profile: creation, XP/leveling, rank tiers, streak tracking, favourite pokemon, rendering.
- `/workspace/poketerm-main/files/0.0.6/trainer.json` — Default trainer state structure.

### New Files to Create

See directory structure below. All files are new (greenfield).

### Key Design Decisions from Conversation
- **Binary name**: `pokedex-fetch-go` (repo: `github.com/brunoorsolon/pokedex-fetch-go`, MIT license)
- **Aliasing**: Document `alias pfg='pokedex-fetch-go'` in README; don't install a short alias binary
- **Sprites**: Gzip-compress before embedding. 92% compression (53 MB → ~4 MB). Decompression ~5-11μs per sprite — invisible. Ship all 4 variants.
- **TUI**: `charmbracelet/bubbletea` for interactive screens (pokedex, trainer, achievements)
- **State**: JSON files in `~/.config/pokedex-fetch-go/` (XDG-compliant). Atomic writes (write-temp + rename).
- **Config**: TOML at `~/.config/pokedex-fetch-go/config.toml`
- **No network**: All Pokemon data baked in at build time. PokeAPI used only in build-time data prep script.
- **No tests/builds in this plan**: Per project rules, do not run `go test`, `go build`, `go vet`, etc. User will validate separately.

---

## Implementation Plan

### Phase 1: Data Preparation & Project Scaffolding

Set up the Go project, prepare the build-time data (enrich pokemon.json, compress sprites, generate generation lists), and establish the embed layer.

### Phase 2: Core — Catch & Display

Implement the default `catch` command (random Pokemon, print sprite), `show` command (specific Pokemon), `list` command, config loading, and `--raw` flag for fastfetch integration.

### Phase 3: State — Pokedex & Catches

Add persistent state (pokedex.json), catch registration on every invocation, shiny roll logic, and the interactive Pokedex TUI with bubbletea.

### Phase 4: Progression — Trainer, XP, Achievements

Add trainer profile, XP/leveling/ranks, streak tracking, achievement system with all milestones and special collections, and the trainer/achievements TUI screens.

### Phase 5: Polish & Distribution

Shell integration snippets, `setup` command, example config, README, cross-compilation, GitHub Releases, poketerm migration tool.

---

## Step-by-Step Tasks

Execute in order, top to bottom. Each task is atomic.

---

### PHASE 1: DATA PREPARATION & PROJECT SCAFFOLDING

---

### 1.1 CREATE Go project scaffolding

Create the project directory structure and initialize Go module.

**Directory structure:**
```
pokedex-fetch-go/
├── cmd/
│   └── pokedex-fetch-go/
│       └── main.go
├── internal/
│   ├── cli/
│   ├── pokemon/
│   ├── state/
│   ├── achievement/
│   ├── tui/
│   └── config/
├── data/
│   ├── sprites/           (populated by prep script)
│   ├── pokemon.json       (populated by prep script)
│   └── gen/               (populated by prep script)
├── scripts/               (build-time only, not shipped)
│   └── prepare-data.sh
├── config.example.toml
├── go.mod
├── LICENSE
└── README.md
```

**IMPLEMENT `cmd/pokedex-fetch-go/main.go`:**
```go
package main

import "github.com/brunoorsolon/pokedex-fetch-go/internal/cli"

func main() {
    cli.Execute()
}
```

**IMPLEMENT `go.mod`:**
```
module github.com/brunoorsolon/pokedex-fetch-go

go 1.22
```

**IMPLEMENT `LICENSE`:** MIT license, copyright Bruno Orsolon.

---

### 1.2 CREATE build-time data preparation script

`scripts/prepare-data.sh` — runs once during development to populate `data/`.

**IMPLEMENT:**
1. **Copy sprites from pokemon-colorscripts:**
   - Source: `/workspace/pokemon-colorscripts-main/colorscripts/`
   - Destination: `data/sprites/`
   - Gzip-compress each file individually: `for f in $(find data/sprites/ -type f ! -name '*.gz'); do gzip -9 "$f" && mv "$f.gz" "$f.gz"; done`
   - Keep the directory structure: `small/regular/`, `small/shiny/`, `large/regular/`, `large/shiny/`
   - Each file gets `.gz` extension: `pikachu` → `pikachu.gz`

2. **Copy generation lists:**
   - Source: `/workspace/poketerm-main/files/gen_files/gen{1..8}_list.txt`
   - Destination: `data/gen/gen1.txt` through `data/gen/gen8.txt`

3. **Enrich pokemon.json:**
   - Source: `/workspace/pokemon-colorscripts-main/pokemon.json` (has name + forms only)
   - Enrich from PokeAPI (cached locally): add number, generation, types, height, weight, base stats, is_legendary, is_mythical
   - Output: `data/pokemon.json` with enriched structure
   - Use a helper script (Python, curl+jq, or Go — whatever runs in dev environment) to call `https://pokeapi.co/api/v2/pokemon/{name}` and `https://pokeapi.co/api/v2/pokemon-species/{name}` for each of the 905 Pokemon
   - Cache API responses to avoid re-fetching
   - **Target schema per entry:**
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
         "hp": 45, "attack": 49, "defense": 49,
         "sp_attack": 65, "sp_defense": 65, "speed": 45
       },
       "is_legendary": false,
       "is_mythical": false
     }
     ```
   - Generation number: derive from which gen list the Pokemon appears in (most reliable), NOT from PokeAPI generation field (which has edge cases for regional forms)

**GOTCHA:** The data prep script is a developer tool, NOT shipped in the binary. It can use Python, curl, jq — whatever is available. The output is what gets embedded.

**GOTCHA:** Some Pokemon names in sprites differ from PokeAPI names (e.g., `nidoran-f` vs `nidoran♀`, `mr-mime` vs `mr. mime`, `farfetchd` vs `farfetch'd`, `type-null` vs `Type: Null`). Use the pokemon-colorscripts names as canonical and map to PokeAPI names when fetching.

**GOTCHA:** Sprite files include alternate forms (e.g., `charizard-mega-x`, `charizard-gmax`). The pokemon.json forms array lists these. The base Pokemon name (e.g., `charizard`) is what gets tracked in the Pokedex — forms are display variants only.

---

### 1.3 CREATE Pokemon data types and embed layer

`internal/pokemon/types.go` — Core data types.

**IMPLEMENT:**
```go
package pokemon

type Pokemon struct {
    Name        string   `json:"name"`
    Number      int      `json:"number"`
    Generation  int      `json:"generation"`
    Forms       []string `json:"forms"`
    Types       []string `json:"types"`
    Height      int      `json:"height"`
    Weight      int      `json:"weight"`
    Stats       Stats    `json:"stats"`
    IsLegendary bool     `json:"is_legendary"`
    IsMythical  bool     `json:"is_mythical"`
}

type Stats struct {
    HP        int `json:"hp"`
    Attack    int `json:"attack"`
    Defense   int `json:"defense"`
    SpAttack  int `json:"sp_attack"`
    SpDefense int `json:"sp_defense"`
    Speed     int `json:"speed"`
}

type SpriteSize string
const (
    SpriteSmall SpriteSize = "small"
    SpriteLarge SpriteSize = "large"
)

type SpriteVariant string
const (
    SpriteRegular SpriteVariant = "regular"
    SpriteShiny   SpriteVariant = "shiny"
)

type Generation struct {
    Number  int
    Names   []string // ordered list of Pokemon names in this generation
}
```

---

### 1.4 CREATE Pokemon data loader (embedded)

`internal/pokemon/data.go` — Loads embedded pokemon.json and generation lists.

**IMPLEMENT:**
```go
package pokemon

import (
    "embed"
    "encoding/json"
    "fmt"
)

//go:embed data/pokemon.json
var pokemonJSON []byte

//go:embed data/gen
var genFS embed.FS
```

Wait — Go embed paths are relative to the source file. Since `data/` is at the project root and `data.go` is in `internal/pokemon/`, the embed directive can't reach `../../data/`. Two approaches:

**Approach A (recommended):** Put the embed directives in the `data/` directory itself as a separate package, or in `cmd/pokedex-fetch-go/main.go`, and pass the FS down.

**Approach B:** Use a top-level `internal/embed.go` that embeds everything and exposes it.

**Use Approach B:**

`internal/embedded/embedded.go`:
```go
package embedded

import "embed"

//go:embed all:sprites
var SpritesFS embed.FS

//go:embed pokemon.json
var PokemonJSON []byte

//go:embed all:gen
var GenFS embed.FS
```

Place `internal/embedded/` adjacent to the data directories, or — simpler — put the embed file at the project root and restructure:

**Final approach:** Place the embed holder at project root as a separate package:

```
data/
├── embed.go              # package data; //go:embed directives
├── sprites/...
├── pokemon.json
└── gen/...
```

`data/embed.go`:
```go
package data

import "embed"

//go:embed pokemon.json
var PokemonJSON []byte

//go:embed all:gen
var GenFS embed.FS

//go:embed all:sprites
var SpritesFS embed.FS
```

Then `internal/pokemon/data.go` imports `github.com/brunoorsolon/pokedex-fetch-go/data` and uses these exported vars.

**GOTCHA:** Using `all:` prefix includes files starting with `.` or `_`. This is fine for sprites but be aware. The sprites directory contains only sprite files — no dotfiles.

**GOTCHA:** The `data` package is part of the module, not truly "internal". This is acceptable — the embed vars are read-only and the package has no behavior.

`internal/pokemon/data.go`:
```go
package pokemon

import (
    "encoding/json"
    "fmt"
    "sort"
    "strconv"
    "strings"

    "github.com/brunoorsolon/pokedex-fetch-go/data"
)

var (
    allPokemon   []Pokemon
    pokemonByName map[string]*Pokemon
    pokemonByNum  map[int]*Pokemon
    generations   []Generation
)

func init() {
    if err := json.Unmarshal(data.PokemonJSON, &allPokemon); err != nil {
        panic(fmt.Sprintf("failed to load pokemon data: %v", err))
    }
    pokemonByName = make(map[string]*Pokemon, len(allPokemon))
    pokemonByNum = make(map[int]*Pokemon, len(allPokemon))
    for i := range allPokemon {
        pokemonByName[allPokemon[i].Name] = &allPokemon[i]
        pokemonByNum[allPokemon[i].Number] = &allPokemon[i]
    }
    loadGenerations()
}

func loadGenerations() {
    // Read gen1.txt through gen8.txt from embedded FS
    // Each file has one Pokemon name per line
    // Build Generation structs
}

func GetByName(name string) (*Pokemon, bool) { ... }
func GetByNumber(num int) (*Pokemon, bool) { ... }
func GetGeneration(gen int) (*Generation, bool) { ... }
func GetAllGenerations() []Generation { ... }
func GetAll() []Pokemon { ... }
func GetRandomFromGenerations(gens []int) *Pokemon { ... }
```

---

### 1.5 CREATE Sprite loader (embedded, gzip-compressed)

`internal/pokemon/sprites.go` — Loads and decompresses sprites.

**IMPLEMENT:**
```go
package pokemon

import (
    "bytes"
    "compress/gzip"
    "fmt"
    "io"

    "github.com/brunoorsolon/pokedex-fetch-go/data"
)

func GetSprite(name string, size SpriteSize, variant SpriteVariant) (string, error) {
    path := fmt.Sprintf("sprites/%s/%s/%s.gz", size, variant, name)
    compressed, err := data.SpritesFS.ReadFile(path)
    if err != nil {
        return "", fmt.Errorf("sprite not found: %s (%s/%s)", name, size, variant)
    }
    r, err := gzip.NewReader(bytes.NewReader(compressed))
    if err != nil {
        return "", fmt.Errorf("failed to decompress sprite: %w", err)
    }
    defer r.Close()
    raw, err := io.ReadAll(r)
    if err != nil {
        return "", fmt.Errorf("failed to read sprite: %w", err)
    }
    return string(raw), nil
}

func HasSprite(name string, size SpriteSize, variant SpriteVariant) bool {
    path := fmt.Sprintf("sprites/%s/%s/%s.gz", size, variant, name)
    _, err := data.SpritesFS.ReadFile(path)
    return err == nil
}
```

---

### PHASE 2: CORE — CATCH & DISPLAY

---

### 2.1 CREATE Config system

`internal/config/config.go` — TOML config with defaults.

**IMPLEMENT:**
```go
package config

import (
    "os"
    "path/filepath"

    "github.com/BurntSushi/toml"
)

type Config struct {
    ShinyRate        int      `toml:"shiny_rate"`
    ShinyRateBoosted int      `toml:"shiny_rate_boosted"`
    SpriteSize       string   `toml:"sprite_size"`
    ShowName         bool     `toml:"show_name"`
    Generations      []int    `toml:"generations"`
    XPNormal         int      `toml:"xp_normal"`
    XPShiny          int      `toml:"xp_shiny"`
    Excluded         []string `toml:"excluded"`
}

func Default() *Config {
    return &Config{
        ShinyRate:        4096,
        ShinyRateBoosted: 10,
        SpriteSize:       "small",
        ShowName:         true,
        Generations:      []int{1, 2, 3, 4, 5, 6, 7, 8},
        XPNormal:         10,
        XPShiny:          150,
        Excluded:         []string{},
    }
}

func DataDir() string {
    if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
        return filepath.Join(xdg, "pokedex-fetch-go")
    }
    home, _ := os.UserHomeDir()
    return filepath.Join(home, ".config", "pokedex-fetch-go")
}

func Load() (*Config, error) {
    cfg := Default()
    path := filepath.Join(DataDir(), "config.toml")
    if _, err := os.Stat(path); os.IsNotExist(err) {
        return cfg, nil // no config file = use defaults
    }
    _, err := toml.DecodeFile(path, cfg)
    return cfg, err
}

func (c *Config) IsExcluded(name string) bool {
    for _, ex := range c.Excluded {
        if ex == name {
            return true
        }
    }
    return false
}
```

**IMPLEMENT `config.example.toml`** at project root:
```toml
# pokedex-fetch-go configuration
# Place this file at ~/.config/pokedex-fetch-go/config.toml

# Shiny encounter rate (1 in N). Default: 4096
# shiny_rate = 4096

# Boosted shiny rate for completed generations (1 in N). Default: 10
# shiny_rate_boosted = 10

# Sprite size: "small" or "large". Default: "small"
# sprite_size = "small"

# Show Pokemon name with sprite. Default: true
# show_name = true

# Generations to pick from. Default: all
# generations = [1, 2, 3, 4, 5, 6, 7, 8]

# XP awards
# xp_normal = 10
# xp_shiny = 150

# Pokemon to never show (by name)
# excluded = []
```

---

### 2.2 CREATE CLI root command

`internal/cli/root.go` — Cobra root command. Default action = catch.

**IMPLEMENT:**
```go
package cli

import (
    "fmt"
    "os"

    "github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
    Use:   "pokedex-fetch-go",
    Short: "Pokemon terminal sprite display and collection tracker",
    Long:  "Display random Pokemon sprites in your terminal, track catches in a persistent Pokedex, and integrate with fastfetch.",
    RunE:  runCatch, // default action is catch
}

func Execute() {
    if err := rootCmd.Execute(); err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
}

func init() {
    // Global flags
    rootCmd.PersistentFlags().String("gen", "", "Generation(s) to pick from (e.g., 1, 1-3, 1,3,6)")
    rootCmd.PersistentFlags().Bool("shiny", false, "Force shiny variant")
    rootCmd.PersistentFlags().Bool("big", false, "Use large sprite")
    rootCmd.PersistentFlags().Bool("no-title", false, "Don't display Pokemon name")
    rootCmd.PersistentFlags().Bool("raw", false, "Output sprite only (for fastfetch piping)")
}
```

---

### 2.3 CREATE `catch` command

`internal/cli/catch.go` — Random Pokemon catch, display sprite, register catch.

**IMPLEMENT:**
```go
package cli

import (
    "fmt"
    "math/rand/v2"

    "github.com/spf13/cobra"
    "github.com/brunoorsolon/pokedex-fetch-go/internal/config"
    "github.com/brunoorsolon/pokedex-fetch-go/internal/pokemon"
)

// Also registered as subcommand for explicit `pokedex-fetch-go catch`
var catchCmd = &cobra.Command{
    Use:   "catch",
    Short: "Catch a random Pokemon (default action)",
    RunE:  runCatch,
}

func runCatch(cmd *cobra.Command, args []string) error {
    cfg, err := config.Load()
    if err != nil {
        return fmt.Errorf("config: %w", err)
    }

    // Parse generation flag or use config default
    gens := cfg.Generations
    if genFlag, _ := cmd.Flags().GetString("gen"); genFlag != "" {
        gens, err = parseGenerations(genFlag)
        if err != nil {
            return err
        }
    }

    // Pick random Pokemon from allowed generations, excluding excluded
    poke := pickRandom(gens, cfg)
    if poke == nil {
        return fmt.Errorf("no Pokemon available for the specified generations")
    }

    // Determine shiny
    forceShiny, _ := cmd.Flags().GetBool("shiny")
    isShiny := forceShiny || rollShiny(poke, cfg)

    // Get sprite
    size := pokemon.SpriteSize(cfg.SpriteSize)
    if big, _ := cmd.Flags().GetBool("big"); big {
        size = pokemon.SpriteLarge
    }
    variant := pokemon.SpriteRegular
    if isShiny {
        variant = pokemon.SpriteShiny
    }

    sprite, err := pokemon.GetSprite(poke.Name, size, variant)
    if err != nil {
        return err
    }

    // Display
    raw, _ := cmd.Flags().GetBool("raw")
    noTitle, _ := cmd.Flags().GetBool("no-title")

    if !raw && !noTitle {
        if isShiny {
            fmt.Printf("\033[33m%s (shiny)\033[0m\n", poke.Name)
        } else {
            fmt.Println(poke.Name)
        }
    }
    fmt.Print(sprite)

    // Register catch (Phase 2 — stub for now, implement in 3.x)
    // state.RegisterCatch(poke.Name, isShiny, cfg)

    return nil
}

func rollShiny(poke *pokemon.Pokemon, cfg *config.Config) bool {
    // Check if Pokemon's generation is complete (Phase 2 — stub)
    // If complete, use boosted rate
    rate := cfg.ShinyRate
    // if genComplete { rate = cfg.ShinyRateBoosted }
    return rand.IntN(rate) == 0
}

func pickRandom(gens []int, cfg *config.Config) *pokemon.Pokemon {
    // Collect all Pokemon from specified generations
    var candidates []pokemon.Pokemon
    for _, g := range gens {
        gen, ok := pokemon.GetGeneration(g)
        if !ok {
            continue
        }
        for _, name := range gen.Names {
            if cfg.IsExcluded(name) {
                continue
            }
            if p, ok := pokemon.GetByName(name); ok {
                candidates = append(candidates, *p)
            }
        }
    }
    if len(candidates) == 0 {
        return nil
    }
    return &candidates[rand.IntN(len(candidates))]
}

// parseGenerations parses "1", "1-3", "1,3,6" into []int
func parseGenerations(s string) ([]int, error) {
    // Handle comma-separated: "1,3,6"
    // Handle range: "1-3" (expands to [1,2,3])
    // Handle single: "1"
    // Validate each is 1-8
    ...
}

func init() {
    rootCmd.AddCommand(catchCmd)
}
```

**PATTERN:** The `runCatch` function is also the default `rootCmd.RunE`, so `pokedex-fetch-go` with no subcommand runs catch. The explicit `catch` subcommand exists for clarity.

---

### 2.4 CREATE `show` command

`internal/cli/show.go` — Show a specific Pokemon by name or number.

**IMPLEMENT:**
```go
package cli

import (
    "fmt"
    "strconv"

    "github.com/spf13/cobra"
    "github.com/brunoorsolon/pokedex-fetch-go/internal/config"
    "github.com/brunoorsolon/pokedex-fetch-go/internal/pokemon"
)

var showCmd = &cobra.Command{
    Use:   "show",
    Short: "Show a specific Pokemon by name or number",
    RunE: func(cmd *cobra.Command, args []string) error {
        cfg, _ := config.Load()
        name, _ := cmd.Flags().GetString("name")
        number, _ := cmd.Flags().GetInt("number")
        form, _ := cmd.Flags().GetString("form")

        var poke *pokemon.Pokemon
        var ok bool
        if name != "" {
            poke, ok = pokemon.GetByName(name)
        } else if number > 0 {
            poke, ok = pokemon.GetByNumber(number)
        } else {
            return fmt.Errorf("specify --name or --number")
        }
        if !ok {
            return fmt.Errorf("Pokemon not found")
        }

        // Handle form: append form to sprite name (e.g., "charizard-mega-x")
        spriteName := poke.Name
        if form != "" {
            spriteName = poke.Name + "-" + form
            // Validate form exists
            if !pokemon.HasSprite(spriteName, pokemon.SpriteSmall, pokemon.SpriteRegular) {
                return fmt.Errorf("form '%s' not found for %s", form, poke.Name)
            }
        }

        isShiny, _ := cmd.Flags().GetBool("shiny")
        size := pokemon.SpriteSize(cfg.SpriteSize)
        if big, _ := cmd.Flags().GetBool("big"); big {
            size = pokemon.SpriteLarge
        }
        variant := pokemon.SpriteRegular
        if isShiny {
            variant = pokemon.SpriteShiny
        }

        sprite, err := pokemon.GetSprite(spriteName, size, variant)
        if err != nil {
            return err
        }

        raw, _ := cmd.Flags().GetBool("raw")
        noTitle, _ := cmd.Flags().GetBool("no-title")
        if !raw && !noTitle {
            if isShiny {
                fmt.Printf("\033[33m%s (shiny)\033[0m\n", spriteName)
            } else {
                fmt.Println(spriteName)
            }
        }
        fmt.Print(sprite)
        return nil
    },
}

func init() {
    showCmd.Flags().StringP("name", "n", "", "Pokemon name")
    showCmd.Flags().IntP("number", "N", 0, "Pokemon number")
    showCmd.Flags().StringP("form", "f", "", "Alternate form (e.g., mega, gmax)")
    rootCmd.AddCommand(showCmd)
}
```

---

### 2.5 CREATE `list` command

`internal/cli/list.go` — List all Pokemon names.

**IMPLEMENT:**
```go
package cli

import (
    "fmt"

    "github.com/spf13/cobra"
    "github.com/brunoorsolon/pokedex-fetch-go/internal/pokemon"
)

var listCmd = &cobra.Command{
    Use:   "list",
    Short: "List all Pokemon names",
    RunE: func(cmd *cobra.Command, args []string) error {
        genFlag, _ := cmd.Flags().GetString("gen")
        if genFlag != "" {
            gens, err := parseGenerations(genFlag)
            if err != nil {
                return err
            }
            for _, g := range gens {
                gen, ok := pokemon.GetGeneration(g)
                if !ok {
                    continue
                }
                for _, name := range gen.Names {
                    fmt.Println(name)
                }
            }
            return nil
        }
        for _, p := range pokemon.GetAll() {
            fmt.Println(p.Name)
        }
        return nil
    },
}

func init() {
    rootCmd.AddCommand(listCmd)
}
```

---

### 2.6 CREATE `config` command

`internal/cli/config.go` — Config init + show.

**IMPLEMENT:**
- `pokedex-fetch-go config init` — creates default config file at `~/.config/pokedex-fetch-go/config.toml`
- `pokedex-fetch-go config path` — prints config directory path
- Copy `config.example.toml` content (hardcoded string) to the config path

---

### PHASE 3: STATE — POKEDEX & CATCHES

---

### 3.1 CREATE XDG path helper

`internal/state/paths.go` — Shared path resolution + directory creation.

**IMPLEMENT:**
```go
package state

import (
    "os"
    "path/filepath"
)

func DataDir() string {
    if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
        return filepath.Join(xdg, "pokedex-fetch-go")
    }
    home, _ := os.UserHomeDir()
    return filepath.Join(home, ".config", "pokedex-fetch-go")
}

func EnsureDataDir() error {
    return os.MkdirAll(DataDir(), 0755)
}

func PokedexPath() string {
    return filepath.Join(DataDir(), "pokedex.json")
}

func TrainerPath() string {
    return filepath.Join(DataDir(), "trainer.json")
}
```

**GOTCHA:** `config.DataDir()` was defined in Phase 2 — consolidate into `state.DataDir()` and have `config` use `state.DataDir()` (or factor out a shared `paths` package to avoid circular deps). Simplest: move `DataDir()` to `state/paths.go` and import from there in both `config` and `state`.

---

### 3.2 CREATE Atomic write helper

`internal/state/write.go` — Write-to-temp + fsync + rename.

**IMPLEMENT:**
```go
package state

import (
    "fmt"
    "os"
    "path/filepath"
)

func AtomicWrite(path string, data []byte) error {
    dir := filepath.Dir(path)
    if err := os.MkdirAll(dir, 0755); err != nil {
        return fmt.Errorf("create dir: %w", err)
    }
    tmp, err := os.CreateTemp(dir, ".tmp-*")
    if err != nil {
        return fmt.Errorf("create temp: %w", err)
    }
    tmpPath := tmp.Name()
    defer func() {
        tmp.Close()
        os.Remove(tmpPath) // cleanup on failure
    }()
    if _, err := tmp.Write(data); err != nil {
        return fmt.Errorf("write: %w", err)
    }
    if err := tmp.Sync(); err != nil {
        return fmt.Errorf("sync: %w", err)
    }
    if err := tmp.Close(); err != nil {
        return fmt.Errorf("close: %w", err)
    }
    if err := os.Rename(tmpPath, path); err != nil {
        return fmt.Errorf("rename: %w", err)
    }
    return nil
}
```

---

### 3.3 CREATE Pokedex state manager

`internal/state/pokedex.go` — Load, save, register catches.

**IMPLEMENT:**

```go
package state

import (
    "encoding/json"
    "os"
    "time"
)

type PokedexState struct {
    Pokemon map[string]*CatchRecord `json:"pokemon"`
    Version int                     `json:"version"`
}

type CatchRecord struct {
    Normal *CatchDetail `json:"normal"`
    Shiny  *CatchDetail `json:"shiny"`
}

type CatchDetail struct {
    Count      int    `json:"count"`
    FirstCaught string `json:"first_caught"`
}

func LoadPokedex() (*PokedexState, error) {
    path := PokedexPath()
    data, err := os.ReadFile(path)
    if os.IsNotExist(err) {
        return &PokedexState{
            Pokemon: make(map[string]*CatchRecord),
            Version: 1,
        }, nil
    }
    if err != nil {
        return nil, err
    }
    var state PokedexState
    if err := json.Unmarshal(data, &state); err != nil {
        return nil, err
    }
    if state.Pokemon == nil {
        state.Pokemon = make(map[string]*CatchRecord)
    }
    return &state, nil
}

func (p *PokedexState) RegisterCatch(name string, isShiny bool) {
    record, ok := p.Pokemon[name]
    if !ok {
        record = &CatchRecord{}
        p.Pokemon[name] = record
    }
    today := time.Now().Format("2006-01-02")
    if isShiny {
        if record.Shiny == nil {
            record.Shiny = &CatchDetail{Count: 0, FirstCaught: today}
        }
        record.Shiny.Count++
    } else {
        if record.Normal == nil {
            record.Normal = &CatchDetail{Count: 0, FirstCaught: today}
        }
        record.Normal.Count++
    }
}

func (p *PokedexState) Save() error {
    data, err := json.MarshalIndent(p, "", "  ")
    if err != nil {
        return err
    }
    return AtomicWrite(PokedexPath(), data)
}

func (p *PokedexState) UniqueCaught() int {
    count := 0
    for _, r := range p.Pokemon {
        if (r.Normal != nil && r.Normal.Count > 0) || (r.Shiny != nil && r.Shiny.Count > 0) {
            count++
        }
    }
    return count
}

func (p *PokedexState) IsCaught(name string) bool { ... }
func (p *PokedexState) IsGenComplete(genNames []string) bool { ... }
func (p *PokedexState) TotalShinies() int { ... }
func (p *PokedexState) HighestCatchCount() int { ... }
func (p *PokedexState) TotalCatches() int { ... }
func (p *PokedexState) GenCaughtCount(genNames []string) int { ... }
```

---

### 3.4 UPDATE `catch` command — wire in state

Update `internal/cli/catch.go` to register catches and check generation completion for boosted shiny rate.

**IMPLEMENT:**
- After displaying sprite, call `pokedexState.RegisterCatch(poke.Name, isShiny)` then `pokedexState.Save()`
- In `rollShiny`, load pokedex state, check if the Pokemon's generation is complete via `pokedexState.IsGenComplete(gen.Names)`, use boosted rate if so
- Update trainer state: XP, streak, daily catch, achievements (calls into Phase 4 — stub until then)

---

### 3.5 CREATE Interactive Pokedex TUI

`internal/tui/pokedex.go` — Bubbletea-based Pokedex list viewer.

**IMPLEMENT using bubbletea:**

Model:
```go
type PokedexModel struct {
    generation   int
    pokemonList  []string          // names in this generation
    pokedexState *state.PokedexState
    cursor       int
    page         int
    pageSize     int               // 15 (matching poketerm)
    width        int
    height       int
}
```

Behavior:
- Arrow up/down: move cursor
- Arrow left/right: change page
- Enter: open detail view (if caught)
- `p`: switch to trainer profile view
- `a`: switch to achievements view
- `q`: quit
- Number keys 1-8: switch generation

Display:
- Box-drawn border (matching poketerm aesthetic)
- Header: "Generation N Pokedex"
- Caught count: "Pokemon Caught: X/Y"
- Table: Pokemon name | Normal (check/cross + count) | Shiny (check/cross + count)
- Footer: page info + controls

**PATTERN from poketerm** (`/workspace/poketerm-main/files/0.0.6/pokedex.py` lines 110-168):
- `PAGE_SIZE = 15`
- `BOX_WIDTH = 37`
- Use check mark (green) for caught, cross (red) for uncaught
- Cursor highlight with ANSI invert

Use `lipgloss` for styling instead of raw ANSI where possible — this is the advantage of bubbletea.

---

### 3.6 CREATE Pokemon detail TUI

`internal/tui/detail.go` — Detail view for a caught Pokemon.

**IMPLEMENT:**

Displays:
- Pokemon number + name
- Sprite (regular or shiny, togglable with left/right if both caught)
- Types
- Height / Weight
- Base stats (HP, Attack, Defense, Sp.Atk, Sp.Def, Speed)

**PATTERN from poketerm** (`/workspace/poketerm-main/files/0.0.6/pokedex.py` lines 171-293):
- Box width 80
- Centered sprite
- Left/right toggles normal/shiny if both owned
- Any key returns to list (or left/right to toggle)

In bubbletea this is a child model that the pokedex model switches to on Enter.

---

### 3.7 CREATE `pokedex` CLI command

`internal/cli/pokedex.go` — Launches TUI.

**IMPLEMENT:**
```go
var pokedexCmd = &cobra.Command{
    Use:   "pokedex [generation]",
    Short: "Interactive Pokedex viewer",
    Args:  cobra.MaximumNArgs(1),
    RunE: func(cmd *cobra.Command, args []string) error {
        gen := 1
        if len(args) > 0 {
            // parse generation number, validate 1-8
        }
        return tui.RunPokedex(gen)
    },
}
```

---

### PHASE 4: PROGRESSION — TRAINER, XP, ACHIEVEMENTS

---

### 4.1 CREATE Trainer state manager

`internal/state/trainer.go` — Load, save, modify trainer profile.

**IMPLEMENT:**

```go
type TrainerState struct {
    Name       string        `json:"name"`
    Region     string        `json:"region"`
    Hometown   string        `json:"hometown"`
    Created    string        `json:"created"`
    Favourites []int         `json:"favourites"` // Pokemon numbers
    Level      int           `json:"level"`
    XP         int           `json:"xp"`
    TotalXP    int           `json:"total_xp"`
    Streak     StreakData    `json:"streak"`
    DailyCatch DailyCatchData `json:"daily_catch"`
    HallOfFame []Achievement `json:"hall_of_fame"`
    Version    int           `json:"version"`
}

type StreakData struct {
    Current        int    `json:"current"`
    Longest        int    `json:"longest"`
    LastActiveDate string `json:"last_active_date"`
}

type DailyCatchData struct {
    Date  string `json:"date"`
    Count int    `json:"count"`
}

type Achievement struct {
    Name     string `json:"name"`
    Event    string `json:"event"`
    Category string `json:"category"`
    Date     string `json:"date"`
}
```

Methods:
- `AddXP(amount int)` — add XP, auto-level-up. Level cost formula: `(level + 1) * 500`, capped at 25500 for level 50+. (Matches `/workspace/poketerm-main/files/0.0.6/trainer.py` lines 149-170)
- `UpdateStreak()` — if same day as last active: no-op. If yesterday: increment current. Else: reset to 1. Update longest. (Matches trainer.py lines 179-206)
- `UpdateDailyCatch()` — increment daily count, reset if new day. (Matches achievements.py lines 273-286)
- `GetRank() string` — return rank title based on level. (Matches trainer.py lines 88-100)
- `Unlock(name, event, category string) bool` — add achievement if not already in hall_of_fame. (Matches achievements.py lines 222-237)
- `Save()` / `Load()` — same atomic write pattern as pokedex.

**Rank tiers** (from trainer.py lines 88-100):
```
0-5:    Beginner Trainer
6-10:   Rookie Trainer
11-15:  Ace Trainer
16-20:  Gym Challenger
21-25:  Gym Leader
26-30:  Elite Trainer
31-35:  Elite Four
36-40:  Champion
41-45:  Regional Master
46-49:  Pokedex Master
50+:    Grand Champion
```

**Regions and towns** (from trainer.py lines 49-82):
- Kanto, Johto, Hoenn, Sinnoh, Unova, Kalos, Alola, Galar — each with 5 town options.

---

### 4.2 CREATE Achievement system

`internal/achievement/definitions.go` — All achievement definitions.

**IMPLEMENT — port from `/workspace/poketerm-main/files/0.0.6/achievements.py`:**

**Pokedex milestones** (lines 306-330):
- 1 caught: "Just getting started"
- 50 caught: "You're getting the hang of this!"
- 100 caught: "A true pro"
- 905 caught: "You sure can open a terminal!"
- 10000 total catches: "How many terminals?"

**Streak achievements** (lines 334-359):
- 3 day: "Warming Up"
- 7 day: "Weeklong Catcher"
- 30 day: "No longer casual thing"
- 100 day: "At this point, the Pokemon fear you"
- 365 day: "You have not missed a single day"

**Shiny milestones** (lines 364-379):
- 1 shiny: "Oooo a Shiny!"
- 5 shinies: "You're really lucky!"
- 100 shinies: "Is this even luck anymore"

**Daily catch milestones** (lines 383-397):
- 5 in a day: "You can almost make a full team!"
- 10 in a day: "You must like opening terminals"
- 50 in a day: "Time to log off now!"

**Highest catch count** (lines 401-410):
- 10 of same: "I hope its not a weedle again!"
- 25 of same: "Not another one!"

**Level milestones** (lines 415-430):
- Level 10: "Apprentice Trainer"
- Level 25: "Veteran Trainer"
- Level 50: "Master Trainer"

**XP milestones** (lines 432-448):
- 1000 XP: "XP Initiate - The grind has begun"
- 10000 XP: "XP Adept - Momentum is undeniable"
- 50000 XP: "XP Elite - You're operating at scale"
- 100000 XP: "XP Grandmaster - The numbers fear you"

**Generation completion** (lines 453-471):
- Each gen: "{Region} Master" — "Completed Generation N"
- All gens: "Lets collect some shinys now!" — "Completed all Generations"

**Special collections** (port SPECIAL_COLLECTIONS from achievements.py lines 36-202):
- All legendaries: "Legendary Conqueror"
- All mythicals: "Myth Hunter"
- Kanto birds: "Legendary Trio Master - Kanto Birds"
- Weather trio: "Weather Dominator"
- All starters: "Starter Supreme"
- All pseudo-legendaries: "Pseudo-Legend Slayer"
- All eeveelutions: "Eeveelution Enthusiast"
- All fossils: "Paleontologist"

**Meta achievement:**
- All achievements unlocked: "Poketerm - Completed it mate!"

---

### 4.3 CREATE Achievement checker

`internal/achievement/checker.go` — Check all achievements against current state.

**IMPLEMENT:**
```go
func CheckAll(pokedex *state.PokedexState, trainer *state.TrainerState, gens []pokemon.Generation) {
    // Run all achievement checks
    // Each check calls trainer.Unlock() if threshold met
    // Runs on every catch — must be fast (just integer comparisons)
}
```

Port the logic from achievements.py `check_pokedex_achievements()` (lines 291-494).

---

### 4.4 UPDATE `catch` command — wire in full progression

Update `internal/cli/catch.go`:

After displaying sprite:
1. Load pokedex state
2. Register catch
3. Save pokedex
4. Load trainer state
5. Add XP (normal or shiny amount from config)
6. Update streak
7. Update daily catch count
8. Run achievement checker
9. Save trainer

This all happens synchronously — it's fast (JSON read + integer math + JSON write).

---

### 4.5 CREATE Trainer TUI

`internal/tui/trainer.go` — Bubbletea trainer profile view.

**IMPLEMENT — port display from trainer.py `render_trainer_profile()` (lines 354-455):**

Display:
- Name, Region, Home Town, Started date
- Level + XP progress bar
- Lifetime XP
- Rank (colored)
- Streak stats (current, longest, last active)
- Daily catch count
- Generation progress bars (per-gen caught/total with visual bar)
- Favourite Pokemon (6, with names)

Controls:
- `1-6`: view favourite Pokemon detail
- `e`: edit profile
- `b`: back

**Trainer creation** (first run or edit):
- Interactive prompts for name (letters only, max 12)
- Choose region from list
- Choose hometown from region's towns
- Choose 6 favourite Pokemon by name or number

For first-run trainer creation, use bubbletea text input + list selection components.

---

### 4.6 CREATE Achievements TUI

`internal/tui/achievements.go` — Hall of Fame display.

**IMPLEMENT — port from achievements.py `display_hall_of_fame()` (lines 523-590):**

Display:
- "HALL OF FAME" header
- Categories: Pokedex Milestones, Catch Achievements, Level Milestones, Special Collections
- Each unlocked achievement: name + event + date
- `b`: back

---

### 4.7 CREATE `trainer` and `achievements` CLI commands

`internal/cli/trainer.go` and `internal/cli/achievements.go`.

**IMPLEMENT:**
```go
// trainer.go
var trainerCmd = &cobra.Command{
    Use:   "trainer",
    Short: "View or edit trainer profile",
    RunE:  func(cmd *cobra.Command, args []string) error {
        return tui.RunTrainer()
    },
}

// achievements.go
var achievementsCmd = &cobra.Command{
    Use:   "achievements",
    Short: "View Hall of Fame",
    RunE:  func(cmd *cobra.Command, args []string) error {
        return tui.RunAchievements()
    },
}
```

---

### PHASE 5: POLISH & DISTRIBUTION

---

### 5.1 CREATE `setup` command

`internal/cli/setup.go` — Prints shell integration snippet.

**IMPLEMENT:**
- Detect current shell from `$SHELL`
- Print appropriate snippet for bash/zsh/fish
- Include the alias suggestion

---

### 5.2 CREATE README.md

**IMPLEMENT:**
- Project description
- Screenshots/demo (placeholder)
- Installation (download binary, chmod, move to PATH)
- Shell integration (bash, zsh, fish + fastfetch)
- Commands reference
- Configuration
- Alias documentation (`alias pfg='pokedex-fetch-go'`)
- Credits (The Pokemon Company, pokemon-colorscripts, poketerm)
- License (MIT)

---

### 5.3 CREATE Makefile / build script

**IMPLEMENT:**
```makefile
BINARY=pokedex-fetch-go

build:
	go build -o $(BINARY) ./cmd/pokedex-fetch-go/

build-all:
	GOOS=linux GOARCH=amd64 go build -o $(BINARY)-linux-amd64 ./cmd/pokedex-fetch-go/
	GOOS=linux GOARCH=arm64 go build -o $(BINARY)-linux-arm64 ./cmd/pokedex-fetch-go/
	GOOS=darwin GOARCH=amd64 go build -o $(BINARY)-darwin-amd64 ./cmd/pokedex-fetch-go/
	GOOS=darwin GOARCH=arm64 go build -o $(BINARY)-darwin-arm64 ./cmd/pokedex-fetch-go/
```

---

### 5.4 CREATE poketerm migration tool (optional)

`internal/cli/migrate.go` — Import existing poketerm pokedex.txt.

**IMPLEMENT:**
- `pokedex-fetch-go migrate --from-poketerm <path-to-pokedex.txt>`
- Parse poketerm's format: `pokemon_name count` or `pokemon_name (shiny) count` per line
- Convert to pokedex.json format
- Merge with existing state (don't overwrite)

---

### 5.5 CREATE `.goreleaser.yml` (optional)

For automated GitHub Releases with goreleaser.

---

## Testing Strategy

> **Note:** Per project rules, do not run `go test`, `go build`, or any automated test/build commands. The user will validate separately.

### Unit Tests (when user decides to run them)
- `internal/pokemon/`: Test data loading, sprite decompression, name/number lookup, generation filtering
- `internal/state/`: Test pokedex register/save/load, trainer XP/level/streak logic, atomic write
- `internal/achievement/`: Test each achievement threshold
- `internal/config/`: Test default values, TOML parsing, excluded check

### Integration Tests
- Full catch flow: pick random → roll shiny → register → save → verify state
- Generation completion detection → boosted shiny rate
- Achievement unlock chain (catch → milestone → unlock)

### Edge Cases
- First run with no state files
- Corrupt JSON state files
- Missing config file (should use defaults)
- Pokemon name with special characters (nidoran-f, mr-mime, type-null, farfetchd, flabebe)
- `--gen` with invalid values
- Concurrent terminal opens (two catches writing state simultaneously — atomic write handles this)

---

## Validation Commands

> Per project rules: do NOT run these. Listed for reference when user chooses to validate.

### Level 1: Syntax & Style
```bash
go vet ./...
gofmt -l .
golangci-lint run
```

### Level 2: Unit Tests
```bash
go test ./...
```

### Level 3: Manual Validation
```bash
# Basic catch
pokedex-fetch-go

# Specific Pokemon
pokedex-fetch-go show --name pikachu
pokedex-fetch-go show --name charizard --form mega-x --shiny --big

# Fastfetch integration
fastfetch --logo-type file-raw --logo <(pokedex-fetch-go catch --raw)

# List
pokedex-fetch-go list --gen 1

# Interactive Pokedex
pokedex-fetch-go pokedex 1

# Trainer
pokedex-fetch-go trainer

# Achievements
pokedex-fetch-go achievements
```

---

## Acceptance Criteria

- [ ] `pokedex-fetch-go` displays a random Pokemon sprite and registers the catch
- [ ] `pokedex-fetch-go show --name pikachu` displays pikachu
- [ ] `pokedex-fetch-go show --name pikachu --shiny --big` displays large shiny pikachu
- [ ] `pokedex-fetch-go catch --gen 1` restricts to Gen 1
- [ ] `pokedex-fetch-go catch --raw` outputs sprite only (no name), suitable for fastfetch
- [ ] `fastfetch --logo-type file-raw --logo <(pokedex-fetch-go catch --raw)` works
- [ ] Catches are persisted in `~/.config/pokedex-fetch-go/pokedex.json`
- [ ] Shiny encounters occur at 1/4096 (or configured rate)
- [ ] Completed generation boosts shiny rate to 1/10 (or configured)
- [ ] `pokedex-fetch-go pokedex` opens interactive TUI with arrow navigation
- [ ] `pokedex-fetch-go pokedex 3` opens Gen 3 view
- [ ] Detail view shows sprite + types + stats for caught Pokemon
- [ ] Trainer profile creation works on first run
- [ ] XP awards correctly (10 normal, 150 shiny)
- [ ] Leveling formula matches poketerm: `(level + 1) * 500`
- [ ] Rank titles display correctly
- [ ] Streaks track correctly across days
- [ ] All 42 achievements can be unlocked
- [ ] Special collections tracked (legendaries, starters, etc.)
- [ ] Config file respected (shiny rate, sprite size, generations, excluded)
- [ ] No config file = works with defaults
- [ ] Binary < 15 MB
- [ ] Startup < 20ms
- [ ] Works on bash, zsh, fish
- [ ] Works on Linux and macOS
