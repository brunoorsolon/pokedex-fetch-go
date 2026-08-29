# Feature: Achievement Points System

Read the referenced codebase files before implementing. Pay attention to existing state persistence patterns and package boundaries. In particular: `internal/achievement` may import `internal/state`, but `internal/state` must not import `internal/achievement` or it will create an import cycle.

## Feature Description

Expand the current initial achievement system into a 50-achievement catalog with point rewards totaling exactly 1000 points. Achievements should remain unlocked automatically on catch/update, persist in trainer state, and display earned points.

The existing system already tracks the right inputs: unique caught Pokemon, total catches, shinies, streaks, daily catches, duplicate catches, level, XP, generation completion, and special collections. The implementation should replace the hardcoded `checkUnlock(...)` chain with a data-driven catalog so point values, IDs, hidden flags, and display metadata live in one place.

## User Story

As a terminal Pokemon collector, I want achievements to award weighted points across 50 goals, so that easy milestones feel rewarding while rare/endgame accomplishments matter more.

## Metadata

- **Type**: Enhancement
- **Complexity**: Medium
- **Systems affected**: achievement checker, trainer state JSON, achievement TUI, trainer TUI, tests/docs
- **Dependencies**: standard Go only; existing Bubble Tea/Lip Gloss UI stack

---

## Context References

### Codebase Files — READ BEFORE IMPLEMENTING

- `internal/achievement/checker.go` (lines 13-113) — Current hardcoded achievement unlock chain; replace with catalog-driven checks.
- `internal/achievement/definitions.go` (lines 4-103) — Current `SpecialCollection` definitions; extend with IDs/points or fold into the main catalog.
- `internal/state/trainer.go` (lines 17-45) — `TrainerState` and `Achievement` JSON structs; add point fields here.
- `internal/state/trainer.go` (lines 119-132) — `LoadTrainer` defaults; ensure new point fields initialize safely.
- `internal/state/trainer.go` (lines 201-214) — Current `Unlock(name, event, category)` behavior; update for IDs and points.
- `internal/state/pokedex.go` (lines 114-176) — Existing stat helpers: `TotalShinies`, `TotalCatches`, `HighestCatchCount`, `CaughtSet`.
- `internal/tui/achievements.go` (lines 12-139) — Current Hall of Fame UI only lists unlocked entries; update to show catalog, lock state, and points.
- `internal/tui/trainer.go` (lines 65-73) — Trainer profile summary; add achievement points `earned/1000`.
- `internal/cli/catch.go` (lines 99-108) — Catch flow calls `achievement.CheckAll(...)`; keep this integration point.
- `internal/pokemon/types.go` (lines 4-15) — Pokemon metadata includes `IsLegendary` and `IsMythical`; useful if replacing hand-maintained legendary/mythical lists later.
- `internal/pokemon/data.go` (lines 77-84) — Generation and full Pokemon accessors used by achievement checks.

### New Files to Create

- `internal/achievement/catalog.go` — Achievement definition catalog, point totals, category constants, condition functions.
- `internal/achievement/catalog_test.go` — Verifies exactly 50 achievements and exactly 1000 total points.
- `internal/achievement/checker_test.go` — Verifies unlock behavior, idempotency, point awarding, and master achievement behavior.

### Documentation — READ BEFORE IMPLEMENTING

- [Bubble Tea Model docs](https://pkg.go.dev/github.com/charmbracelet/bubbletea#Model) — Existing TUI follows `Init`, `Update`, `View`.
- [Lip Gloss width docs](https://github.com/charmbracelet/lipgloss/blob/f497c220573977914dd03ad6fc2d7fee52239470/size.go) — Use `lipgloss.Width`, not `len`, when aligning colored rows.
- [Go encoding/json docs](https://pkg.go.dev/encoding/json) — Useful for confirming JSON persistence behavior for new point fields.
- [poketerm achievements reference](https://github.com/chris-wood-mo/poketerm#achievements) — Original inspiration: 42-ish achievements across pokedex, catch, rank, and special categories.

### Patterns to Follow

- Keep state writes atomic via existing `TrainerState.Save()` and `PokedexState.Save()` patterns.
- Keep achievement checking cheap: current `CheckAll` runs after every catch and only does simple integer/set checks.
- Use stable internal IDs for achievements. Names can change later; IDs should not.
- Avoid package cycles: state structs can store `ID`, `Points`, and `AchievementPoints`, but catalog lookup should live in `internal/achievement`.

---

## Achievement Catalog

Implement exactly these 50 achievements. Point total must equal exactly **1000**.

| # | ID | Name | Event / Requirement | Category | Points |
|---:|---|---|---|---|---:|
| 1 | `pokedex_first_catch` | Just getting started | Caught your first Pokemon | pokedex | 5 |
| 2 | `pokedex_25_unique` | Route 1 Regular | Caught 25 Pokemon | pokedex | 8 |
| 3 | `pokedex_50_unique` | You're getting the hang of this! | Caught 50 Pokemon | pokedex | 10 |
| 4 | `pokedex_100_unique` | A true pro | Caught 100 Pokemon | pokedex | 15 |
| 5 | `pokedex_250_unique` | Regional Researcher | Caught 250 Pokemon | pokedex | 25 |
| 6 | `pokedex_all_905` | You sure can open a terminal! | Caught all 905 Pokemon | pokedex | 50 |
| 7 | `catch_1000_total` | Frequent Flyer | Caught 1000 Pokemon total | catch | 12 |
| 8 | `catch_10000_total` | How many terminals? | Caught 10000 Pokemon total | catch | 25 |
| 9 | `streak_3_days` | Warming Up | 3 Day Catch Streak | catch | 5 |
| 10 | `streak_7_days` | Weeklong Catcher | 7 Day Catch Streak | catch | 10 |
| 11 | `streak_30_days` | No longer casual thing | 30 Day Catch Streak | catch | 20 |
| 12 | `streak_100_days` | At this point, the Pokemon fear you | 100 Day Catch Streak | catch | 30 |
| 13 | `streak_365_days` | You have not missed a single day | 365 Day Catch Streak | catch | 35 |
| 14 | `shiny_first` | Oooo a Shiny! | Caught your first Shiny Pokemon | catch | 15 |
| 15 | `shiny_5` | You're really lucky! | Caught 5 Shiny Pokemon | catch | 25 |
| 16 | `shiny_25` | Sparkle Specialist | Caught 25 Shiny Pokemon | catch | 35 |
| 17 | `shiny_100` | Is this even luck anymore | Caught 100 Shiny Pokemon | catch | 45 |
| 18 | `daily_5` | You can almost make a full team! | Caught 5 Pokemon in a day | catch | 5 |
| 19 | `daily_10` | You must like opening terminals | Caught 10 Pokemon in a day | catch | 8 |
| 20 | `daily_25` | Terminal Marathon | Caught 25 Pokemon in a day | catch | 12 |
| 21 | `daily_50` | Time to log off now! | Caught 50 Pokemon in a day | catch | 20 |
| 22 | `duplicate_10` | I hope its not a weedle again! | Caught a Pokemon 10 times | catch | 10 |
| 23 | `duplicate_25` | Not another one! | Caught a Pokemon 25 times | catch | 15 |
| 24 | `duplicate_50` | Familiar Face | Caught a Pokemon 50 times | catch | 20 |
| 25 | `level_10` | Apprentice Trainer | Reach level 10 | rank | 10 |
| 26 | `level_25` | Veteran Trainer | Reach level 25 | rank | 15 |
| 27 | `level_50` | Master Trainer | Reach level 50 | rank | 25 |
| 28 | `xp_1000` | XP Initiate - The grind has begun | Reach a total XP of 1000 | rank | 5 |
| 29 | `xp_10000` | XP Adept - Momentum is undeniable | Reach a total XP of 10000 | rank | 10 |
| 30 | `xp_50000` | XP Elite - You're operating at scale | Reach a total XP of 50000 | rank | 15 |
| 31 | `xp_100000` | XP Grandmaster - The numbers fear you | Reach a total XP of 100000 | rank | 20 |
| 32 | `xp_250000` | XP Beyond Champion | Reach a total XP of 250000 | rank | 20 |
| 33 | `gen_1_complete` | Kanto Master | Completed Generation 1 | pokedex | 20 |
| 34 | `gen_2_complete` | Johto Master | Completed Generation 2 | pokedex | 15 |
| 35 | `gen_3_complete` | Hoenn Master | Completed Generation 3 | pokedex | 20 |
| 36 | `gen_4_complete` | Sinnoh Master | Completed Generation 4 | pokedex | 15 |
| 37 | `gen_5_complete` | Unova Master | Completed Generation 5 | pokedex | 25 |
| 38 | `gen_6_complete` | Kalos Master | Completed Generation 6 | pokedex | 10 |
| 39 | `gen_7_complete` | Alola Master | Completed Generation 7 | pokedex | 15 |
| 40 | `gen_8_complete` | Galar Master | Completed Generation 8 | pokedex | 20 |
| 41 | `all_generations_complete` | Lets collect some shinys now! | Completed all Generations | pokedex | 40 |
| 42 | `special_kanto_birds` | Legendary Trio Master - Kanto Birds | Completed the Legendary Birds | special | 10 |
| 43 | `special_weather_trio` | Weather Dominator | Completed the Weather Trio | special | 15 |
| 44 | `special_eeveelutions` | Eeveelution Enthusiast | Caught All Eeveelutions | special | 15 |
| 45 | `special_pseudo_legendaries` | Pseudo-Legend Slayer | Caught All Pseudo-Legendaries | special | 20 |
| 46 | `special_fossils` | Paleontologist | Caught All Fossil Pokemon | special | 20 |
| 47 | `special_starters` | Starter Supreme | Caught All Starter Pokemon | special | 25 |
| 48 | `special_mythicals` | Myth Hunter | Caught All Mythical Pokemon | special | 35 |
| 49 | `special_legendaries` | Legendary Conqueror | Caught All Legendary Pokemon | special | 50 |
| 50 | `achievement_completionist` | Poketerm - Completed it mate! | Completed all other achievements | special | 50 |

---

## Implementation Plan

### Phase 1: Foundation

1. Add stable achievement IDs and point storage to trainer state.
2. Create a catalog-driven achievement definition model.
3. Add invariants: exactly 50 definitions, exactly 1000 points, unique IDs, unique names.

### Phase 2: Core Implementation

1. Replace hardcoded checks in `CheckAll` with context-based rule evaluation.
2. Build an `achievement.Context` from `PokedexState`, `TrainerState`, and generations.
3. Unlock achievements by ID and award points only once.
4. Ensure master achievement unlocks only after all 49 non-master achievements are unlocked.

### Phase 3: Integration

1. Keep `internal/cli/catch.go` calling `achievement.CheckAll(...)` after catch updates.
2. Update achievements TUI to display total progress and points.
3. Update trainer profile TUI to show achievement points.
4. Update README usage/feature bullets to mention 50 achievements and 1000 points.

### Phase 4: Testing & Validation

1. Add unit tests for catalog totals and uniqueness.
2. Add unit tests for unlock idempotency and point totals.
3. Manual test TUI display and catch flow using temporary `XDG_CONFIG_HOME`.

---

## Step-by-Step Tasks

Execute in order, top to bottom. Each task is atomic.

### UPDATE `internal/state/trainer.go`

- **IMPLEMENT**:
  - Add `AchievementPoints int \`json:"achievement_points"\`` to `TrainerState`.
  - Add fields to `Achievement`:
    - `ID string \`json:"id,omitempty"\``
    - `Points int \`json:"points"\``
  - Change `Unlock` to accept `id, name, event, category string, points int`.
  - Update duplicate detection to use the stable `ID`.
  - On new unlock, append the points to both the entry and `TrainerState.AchievementPoints`.
  - Add a helper like `RecalculateAchievementPoints()` that sums `HallOfFame[].Points` and sets `AchievementPoints`.
- **PATTERN**: Existing `Unlock` appends if not already present (`internal/state/trainer.go:201-214`).
- **GOTCHA**: Do not import `internal/achievement` here.
- **VALIDATE**: `gofmt -w internal/state/trainer.go`

### CREATE `internal/achievement/catalog.go`

- **IMPLEMENT**:
  - Define `AchievementDefinition` with `ID`, `Name`, `Event`, `Category`, `Points`, optional `Hidden`.
  - Define unexported rule wrapper with `Condition func(Context) bool`.
  - Define `Context` containing all derived counters needed by rules:
    - `Unique`, `TotalCatches`, `TotalShinies`, `HighestCatchCount`, `Level`, `TotalXP`, `LongestStreak`, `DailyCount`, `CaughtSet`, generation completion map.
  - Add `Definitions() []AchievementDefinition` for TUI/tests.
  - Add lookup helper `DefinitionByID` for UI/checker use.
  - Encode the exact 50-row catalog above.
- **PATTERN**: Move current milestone names from `internal/achievement/checker.go:27-106` into this catalog.
- **GOTCHA**: Master definition must be excluded when checking whether all other achievements are unlocked.
- **VALIDATE**: `gofmt -w internal/achievement/catalog.go`

### UPDATE `internal/achievement/definitions.go`

- **IMPLEMENT**:
  - Add `ID string` and `Points int` to `SpecialCollection`, or replace special collection point/name handling with catalog rules.
  - Ensure existing special collection names map to these IDs:
    - `special_legendaries`
    - `special_mythicals`
    - `special_kanto_birds`
    - `special_weather_trio`
    - `special_starters`
    - `special_pseudo_legendaries`
    - `special_eeveelutions`
    - `special_fossils`
- **PATTERN**: Existing collection list starts at `internal/achievement/definitions.go:11`.
- **GOTCHA**: Keep canonical Pokemon names lowercase because `PokedexState.CaughtSet()` returns stored lowercase names.
- **VALIDATE**: `gofmt -w internal/achievement/definitions.go`

### REFACTOR `internal/achievement/checker.go`

- **IMPLEMENT**:
  - Build `Context` once at the top from existing helper methods.
  - Loop through all non-master rules and call `tr.Unlock(def.ID, def.Name, def.Event, def.Category, def.Points)` when condition is true.
  - After that loop, unlock `achievement_completionist` only if all other 49 IDs are unlocked.
  - Make `CheckAll` return a small result struct if useful, e.g. `{NewUnlocks []AchievementDefinition}`. Optional; not required by current UI.
- **PATTERN**: Preserve current cheap stat extraction from `internal/achievement/checker.go:14-22`.
- **GOTCHA**: Current `allUnlocked` chaining is fragile; compute completionist state from actual unlocked ID set after rule evaluation.
- **VALIDATE**: `gofmt -w internal/achievement/checker.go`

### UPDATE `internal/tui/achievements.go`

- **IMPLEMENT**:
  - Display catalog progress, not only `trainer.HallOfFame`.
  - Header should show: `Achievements: X/50   Points: Y/1000`.
  - For each row show status, points, category, name, and event:
    - unlocked: `✓ 025 pts [pokedex] Kanto Master — Completed Generation 1 (2026-04-29)`
    - locked: `□ 025 pts [pokedex] Kanto Master — Completed Generation 1`
  - Sort/group by catalog order for predictable browsing.
  - Keep pagination and existing nav keys.
  - Optionally add category filters later; do not block this feature on filters.
  - Before display, load pokedex and call `achievement.CheckAll(...)` so newly satisfied achievements appear even if the user opens achievements before the next catch.
- **PATTERN**: Current TUI pagination and `lipgloss.Width` truncation are in `internal/tui/achievements.go:57-125`.
- **IMPORTS**: Add `internal/achievement` and `internal/pokemon`; possibly `internal/state` is already present.
- **GOTCHA**: Avoid displaying raw ANSI control widths with `len`; use existing `truncateDisplay` pattern.
- **VALIDATE**: `gofmt -w internal/tui/achievements.go`

### UPDATE `internal/tui/trainer.go`

- **IMPLEMENT**:
  - Replace or supplement `Achievements: N` with `Achievements: N/50`.
  - Add `Achievement Points: Y/1000`.
- **PATTERN**: Existing profile lines are in `internal/tui/trainer.go:65-73`.
- **GOTCHA**: Prefer calling `tr.RecalculateAchievementPoints()` defensively before rendering so the displayed total always matches Hall of Fame entries.
- **VALIDATE**: `gofmt -w internal/tui/trainer.go`

### UPDATE `README.md`

- **IMPLEMENT**:
  - Update feature list to mention `50 achievements` and `1000 total achievement points`.
  - Add a short achievement section explaining categories and points.
- **PATTERN**: Existing feature list in `README.md` near top.
- **VALIDATE**: `grep -n "50 achievements\|1000" README.md`

### CREATE `internal/achievement/catalog_test.go`

- **IMPLEMENT**:
  - Assert `len(Definitions()) == 50`.
  - Assert total points equals `1000`.
  - Assert IDs are non-empty and unique.
  - Assert names are non-empty and unique.
  - Assert no achievement has points <= 0.
- **PATTERN**: No existing tests; use Go standard `testing` package.
- **VALIDATE**: `go test ./internal/achievement`

### CREATE `internal/achievement/checker_test.go`

- **IMPLEMENT**:
  - Build minimal `PokedexState`/`TrainerState` fixtures in-memory.
  - Test first catch unlocks `pokedex_first_catch` and awards 5 points once.
  - Test repeated `CheckAll` does not duplicate points.
  - Test master achievement does not unlock until all non-master IDs are unlocked.
- **PATTERN**: Use state structs directly; no filesystem needed.
- **GOTCHA**: Avoid tests that depend on current date except checking non-empty date.
- **VALIDATE**: `go test ./internal/achievement`

---

## Testing Strategy

### Unit Tests

- Catalog invariants:
  - exactly 50 achievements
  - exactly 1000 total points
  - unique IDs/names
  - all categories are one of `pokedex`, `catch`, `rank`, `special`
- Unlock behavior:
  - achievement unlocks when condition becomes true
  - same achievement does not add duplicate points
  - completionist unlock requires all 49 other achievements

### Integration Tests

- Use temp `XDG_CONFIG_HOME` and run catch commands manually or via CLI-level tests if added later.
- Confirm `trainer.json` stores `achievement_points` and per-achievement `id`/`points`.
- Confirm `achievements` TUI shows `X/50` and `Y/1000`.

### Edge Cases

- Completionist should not count itself toward its own condition.
- Generation completion should handle missing/empty generation data defensively.

---

## Validation Commands

### Level 1: Syntax & Style

```bash
gofmt -w internal/state/trainer.go internal/achievement/catalog.go internal/achievement/definitions.go internal/achievement/checker.go internal/tui/achievements.go internal/tui/trainer.go internal/achievement/catalog_test.go internal/achievement/checker_test.go
```

### Level 2: Unit Tests

```bash
go test ./internal/achievement
```

### Level 3: Full Project Tests

```bash
go test ./...
```

### Level 4: Build

```bash
go build -o bin/pokedex-fetch-go ./cmd/pokedex-fetch-go
```

### Level 5: Manual Validation

```bash
XDG_CONFIG_HOME="$(mktemp -d)" ./bin/pokedex-fetch-go catch --gen 1 --raw >/dev/null
```

```bash
XDG_CONFIG_HOME="$(mktemp -d)" ./bin/pokedex-fetch-go achievements
```

```bash
XDG_CONFIG_HOME="$(mktemp -d)" ./bin/pokedex-fetch-go trainer
```

---

## Acceptance Criteria

- [ ] Achievement catalog contains exactly 50 achievements.
- [ ] Achievement points sum to exactly 1000.
- [ ] Unlocking an achievement awards points once only.
- [ ] Trainer state persists total achievement points and per-achievement points.
- [ ] Achievements TUI shows locked/unlocked catalog progress and points.
- [ ] Trainer TUI shows achievement count and points.
- [ ] Completionist achievement unlocks only after the other 49 achievements.
- [ ] Unit tests cover catalog invariants, idempotency, and completionist behavior.
- [ ] Code follows existing package boundaries and TUI style.
