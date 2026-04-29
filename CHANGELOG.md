# Changelog

## [1.0.0] — 2026-04-29

`pokedex-fetch-go` is a single-binary Go CLI for showing Pokemon ANSI sprites in the terminal, catching random Pokemon, tracking collection progress, and integrating with `fastfetch`.

### Highlights

- Random Pokemon encounters directly in the terminal.
- Persistent Pokedex and trainer progression under `~/.config/pokedex-fetch-go`.
- Interactive terminal UIs for Pokedex, trainer profile, and achievements.
- `fastfetch`-compatible raw sprite output.
- 50 achievements worth 1000 total points.
- First-run config bootstrap with editable TOML defaults.

### Added

#### Core CLI

- `pokedex-fetch-go` / `pokedex-fetch-go catch` to catch a random Pokemon.
- `show` command to display a specific Pokemon by name or National Dex number.
- `list` command to list Pokemon by generation.
- `pokedex` command for the interactive collection browser.
- `trainer` command for trainer profile/progression.
- `achievements`, `hof`, and `hall-of-fame` commands for achievement progress.
- `config init` and `config path` commands.
- `setup bash` helper for shell/fastfetch integration.

#### Pokemon Display

- Embedded ANSI sprites from `pokemon-colorscripts`.
- Regular and shiny variants.
- Small and large sprite support.
- Raw output mode for piping into tools like `fastfetch`.
- Optional Pokemon title display.

#### Pokedex TUI

- Generation-based interactive Pokedex view.
- Clean aligned rows with National Dex numbers.
- Zero-padded Dex numbers, e.g. `#025 Pikachu`.
- Normal and shiny catch counts per Pokemon.
- Proper display capitalization and special-case names like:
  - `Mr. Mime`
  - `Ho-Oh`
  - `Type: Null`
  - `Nidoran♀` / `Nidoran♂`
- `g` go-to-number mode for jumping to a National Dex number.
- Detail view on `enter`, including small sprite plus Pokemon metadata/caught status.

#### Trainer Progression

- Trainer profile with name, region, and hometown.
- XP and level progression.
- Total lifetime XP tracking.
- Daily catch counter.
- Catch streak tracking.
- Trainer rank/title based on level.

#### Achievements

- 50 achievement catalog with stable internal IDs.
- 1000 total achievement points.
- Categories:
  - Pokedex progress
  - Catch milestones
  - Streak milestones
  - Shiny milestones
  - Daily catch milestones
  - Trainer rank / XP milestones
  - Generation completion
  - Special collections
- Completionist achievement for unlocking all other achievements.
- Hall of Fame TUI now shows locked/unlocked achievements, points, categories, and progress.
- Trainer profile shows achievement count and total achievement points.

#### Configuration

- First run automatically creates:

```text
~/.config/pokedex-fetch-go/config.toml
```

- Configurable options include:
  - shiny rate
  - boosted shiny rate after generation completion
  - sprite size
  - title display
  - enabled generations
  - XP rewards
  - excluded Pokemon
- Default shiny rate is `1/1028`.

#### Persistence

- Pokedex state saved as JSON.
- Trainer state saved as JSON.
- Atomic state writes to reduce risk of partial/corrupt files.
- Build artifacts and local runtime files ignored via `.gitignore`.

#### Tests

- Added baseline tests across all packages.
- Coverage includes:
  - achievement catalog invariants
  - achievement unlock idempotency
  - completionist achievement behavior
  - config defaults and bootstrap
  - Pokedex state calculations and save/load
  - trainer XP/level progression and save/load
  - embedded Pokemon data lookup
  - sprite loading
  - CLI generation parsing and catch flow
  - TUI display helpers and detail rendering
  - embedded data smoke tests

### Notes

- This release ships as GitHub release binaries only. No distro packages are provided yet.
- Pokemon names, designs, and branding belong to The Pokemon Company.
- Sprites are sourced from `pokemon-colorscripts`.
- The project is inspired by `poketerm`, rewritten in Go with a different architecture.