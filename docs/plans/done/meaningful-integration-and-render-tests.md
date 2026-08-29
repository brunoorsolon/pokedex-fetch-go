# Feature: Meaningful CLI, Sprite, and TUI Tests

Read the referenced codebase files before implementing. The goal is to add high-signal tests for real behavior, not placeholder coverage. Keep tests deterministic, isolated with `t.TempDir()`/`t.Setenv()`, and avoid assertions that depend on a specific random Pokemon unless the test controls the candidate pool.

## Feature Description

Add meaningful test coverage for the remaining important workflows:

- CLI catch behavior using a temporary `XDG_CONFIG_HOME`.
- Catch flow persistence: `config.toml`, `pokedex.json`, and `trainer.json` are written and contain expected state.
- Sprite loading for regular and shiny small sprites.
- Pokedex detail TUI rendering for selected Pokemon, including name/number/caught status.

This complements the existing unit tests for config defaults, parsing, state calculations, embedded data, achievements, and display-name formatting.

## User Story

As a maintainer, I want tests that exercise real CLI, persistence, sprite, and TUI behavior so that changes to catch flow or rendering fail quickly without relying only on manual testing.

## Metadata

- **Type**: Enhancement
- **Complexity**: Low/Medium
- **Systems affected**: `internal/cli`, `internal/pokemon`, `internal/tui`, `internal/state`, `internal/config`
- **Dependencies**: Go standard `testing`, `os`, `io`, `strings`; existing Cobra CLI and Bubble Tea/Lip Gloss code

---

## Context References

### Codebase Files — READ BEFORE IMPLEMENTING

- `internal/cli/catch.go` (line 23) — `runCatch` is the catch-flow entry point to test directly.
- `internal/cli/catch.go` (lines 24-35) — Loads config and parses `--gen`; test invalid generation behavior and temp config bootstrap.
- `internal/cli/catch.go` (lines 37-48) — Loads pokedex and chooses a random Pokemon; assertions must not depend on exact species unless candidate pool is controlled.
- `internal/cli/catch.go` (lines 50-78) — Shiny/sprite/output behavior; use `--shiny --raw` for deterministic shiny catch and quiet-enough output.
- `internal/cli/catch.go` (lines 80-108) — Persists pokedex/trainer state and checks achievements; assert files and loaded state.
- `internal/cli/catch_test.go` (lines 10-61) — Existing package-level test style for unexported CLI helpers.
- `internal/config/config.go` (lines 66-81) — `DataDir`/`ConfigPath` respect `XDG_CONFIG_HOME`; use these in tests.
- `internal/config/config_test.go` (lines 46-75) — Existing `t.TempDir()` + `t.Setenv("XDG_CONFIG_HOME", ...)` pattern.
- `internal/state/pokedex.go` (line 36) — `LoadPokedex` reads persisted catch state.
- `internal/state/trainer.go` (line 122) — `LoadTrainer` reads persisted XP, daily count, and achievements.
- `internal/pokemon/sprites.go` (line 13) — `GetSprite` loads and decompresses embedded ANSI sprites.
- `internal/pokemon/sprites.go` (line 32) — `HasSprite` checks embedded sprite existence.
- `data/embed_test.go` (line 5) — Current data smoke test only verifies an embedded sprite file exists; add behavior tests in `internal/pokemon`.
- `internal/tui/detail.go` (line 19) — `newDetailModel` constructs detail models without launching Bubble Tea.
- `internal/tui/detail.go` (line 36) — `detailModel.View` renders selected Pokemon details and sprite.
- `internal/tui/detail.go` (line 131) — `splitSpriteLines` is pure and useful for targeted layout helper tests if needed.
- `internal/tui/pokedex_test.go` (line 5) — Existing TUI test style for pure helpers in package `tui`.

### New Files to Create

- `internal/cli/catch_flow_test.go` — Integration-style unit tests for `runCatch` using temp config/state directories.
- `internal/pokemon/sprites_test.go` — Tests sprite loading/decompression for regular and shiny small sprites.
- `internal/tui/detail_test.go` — Tests detail view rendering, missing Pokemon fallback, and key handling.

### Documentation — READ BEFORE IMPLEMENTING

- [Go testing package](https://pkg.go.dev/testing) — `T.TempDir`, `T.Setenv`, table-driven tests, and cleanup patterns.
- [Cobra command testing examples](https://github.com/spf13/cobra/blob/main/command_test.go) — Use `SetArgs`, output buffers, or direct command construction where appropriate.
- [Bubble Tea testing guide](https://www.mintlify.com/charmbracelet/bubbletea/guides/testing) — View output can be tested by constructing the model and asserting rendered content.
- [Lip Gloss width/style behavior](https://github.com/charmbracelet/lipgloss/blob/f497c220573977914dd03ad6fc2d7fee52239470/size.go) — Rendered strings include ANSI styling; substring assertions should account for ANSI when needed.

### Patterns to Follow

- Keep tests in the same package (`package cli`, `package tui`) when testing unexported functions/models.
- Use `t.Setenv("XDG_CONFIG_HOME", t.TempDir())` to isolate config and state files.
- Do not use `t.Parallel()` in tests that temporarily replace `os.Stdout` or mutate global Cobra command state.
- Prefer direct `runCatch(cmd, nil)` tests with a fresh test command over invoking `Execute()`, because `Execute()` can call `os.Exit` on error.
- For random catch flow, assert invariants (`TotalCatches == 1`, `TotalShinies == 1` with `--shiny`) instead of exact Pokemon name.
- For styled TUI output, use `strings.Contains` on stable visible text. If ANSI styling breaks a check, add a small test-only `stripANSI` helper.

---

## Implementation Plan

### Phase 1: CLI Catch Test Harness

Create a small helper in `internal/cli/catch_flow_test.go` that builds a fresh `*cobra.Command` with the flags read by `runCatch`:

- `gen` string
- `shiny` bool
- `big` bool
- `no-title` bool
- `raw` bool

Add a helper to capture `os.Stdout` because `runCatch` uses `fmt.Print` directly instead of Cobra's command output writer.

### Phase 2: Catch Flow Persistence Tests

Add tests that:

- Set `XDG_CONFIG_HOME` to a temp dir.
- Run `runCatch` with `--gen 1 --shiny --raw`.
- Assert `config.toml`, `pokedex.json`, and `trainer.json` exist under `<temp>/pokedex-fetch-go/`.
- Load state via `state.LoadPokedex()` and `state.LoadTrainer()`.
- Assert:
  - one total catch
  - one unique catch
  - one shiny catch
  - trainer `TotalXP == 150`
  - trainer daily catch count is `1`
  - first-catch and shiny achievements exist, or at least `AchievementPoints == 20`

Also add a negative test for invalid generation input:

- Run with `--gen 9`.
- Assert error is returned.
- Assert catch state files were not created.
- Config file may be created because `runCatch` loads config before parsing `--gen`; document/assert that behavior explicitly if checking files.

### Phase 3: Sprite Tests

Add sprite behavior tests in `internal/pokemon/sprites_test.go`:

- `HasSprite("pikachu", SpriteSmall, SpriteRegular) == true`
- `HasSprite("pikachu", SpriteSmall, SpriteShiny) == true`
- `GetSprite("pikachu", SpriteSmall, SpriteRegular)` returns non-empty string.
- `GetSprite("pikachu", SpriteSmall, SpriteShiny)` returns non-empty string.
- Regular and shiny sprite strings differ.
- Missing sprite returns an error and `HasSprite` is false.

Avoid exact full-sprite snapshots because ANSI art can change and snapshot failures would be noisy.

### Phase 4: TUI Detail Rendering Tests

Add tests in `internal/tui/detail_test.go`:

- Construct `newDetailModel("pikachu", true, false)` and call `View()`.
- Assert visible output contains:
  - `#025 Pikachu`
  - `Generation:`
  - `Caught:`
  - `normal`
  - `caught`
  - `shiny`
  - `not caught`
  - footer `enter/esc/q back`
- Construct `newDetailModel("pikachu", false, true)` and assert shiny caught state renders.
- Construct `newDetailModel("missingno", false, false)` and assert `Pokemon not found`.
- Send `tea.KeyMsg{Type: tea.KeyEnter}` to `Update` and assert `done == true`.

If ANSI styling prevents reliable substring checks, add a local helper:

```go
var ansiRE = regexp.MustCompile(`\x1b\[[0-9;]*m`)
func stripANSI(s string) string { return ansiRE.ReplaceAllString(s, "") }
```

### Phase 5: Validation

Run formatting and tests locally. In this assistant environment, do not run tests/builds if project instructions prohibit it; the implementer/user should run them locally.

---

## Step-by-Step Tasks

Execute in order, top to bottom. Each task is atomic.

### CREATE `internal/cli/catch_flow_test.go`

- **IMPLEMENT**:
  - Add imports: `io`, `os`, `path/filepath`, `strings` if needed, `testing`, `github.com/spf13/cobra`, plus internal `config`, `state`.
  - Add `newCatchFlowTestCommand()` returning a fresh Cobra command with flags consumed by `runCatch`.
  - Add `captureStdout(t, fn)` helper using `os.Pipe()` and restoring `os.Stdout` with `t.Cleanup` or `defer`.
  - Add `TestRunCatchWritesConfigPokedexAndTrainerState`.
  - Add `TestRunCatchRejectsInvalidGenerationBeforeCatchStateWrite`.
- **PATTERN**: Existing CLI tests in `internal/cli/catch_test.go:10-61` directly test unexported helpers in `package cli`.
- **IMPORTS**:
  - `github.com/spf13/cobra`
  - `github.com/brunoorsolon/pokedex-fetch-go/internal/config`
  - `github.com/brunoorsolon/pokedex-fetch-go/internal/state`
- **GOTCHA**:
  - Do not call `Execute()`; it can exit the process on error.
  - Do not assert exact caught Pokemon due randomness.
  - `runCatch` prints to `os.Stdout`, not `cmd.OutOrStdout()`, so capture stdout with `os.Pipe()`.
  - Config bootstrap happens before generation parsing, so invalid generation may still create `config.toml`.
- **VALIDATE**: `gofmt -w internal/cli/catch_flow_test.go`

### CREATE `internal/pokemon/sprites_test.go`

- **IMPLEMENT**:
  - Add `TestGetSpriteLoadsSmallRegularAndShiny`.
  - Add `TestGetSpriteMissingPokemonReturnsError`.
  - Assert non-empty sprite strings and regular != shiny.
  - Optionally assert sprite contains ANSI escape prefix `\x1b[` if stable in current assets.
- **PATTERN**: Existing data smoke test in `data/embed_test.go:5` verifies embedded files; this test should verify public sprite API behavior.
- **GOTCHA**:
  - Do not snapshot full sprite art.
  - Test `SpriteSmall`; large sprites are intentionally not part of this baseline.
- **VALIDATE**: `gofmt -w internal/pokemon/sprites_test.go`

### CREATE `internal/tui/detail_test.go`

- **IMPLEMENT**:
  - Add `TestDetailViewShowsPokemonDataAndCaughtStatus`.
  - Add `TestDetailViewShowsNotFoundForUnknownPokemon`.
  - Add `TestDetailModelEnterMarksDone`.
  - Use `newDetailModel(...)` directly; do not launch a Bubble Tea program.
  - Use `tea.KeyMsg{Type: tea.KeyEnter}` for update behavior.
  - Add `stripANSI` helper only if simple substring checks fail.
- **PATTERN**: Existing TUI helper test in `internal/tui/pokedex_test.go:5` uses direct function assertions.
- **IMPORTS**:
  - `strings`
  - `testing`
  - `tea "github.com/charmbracelet/bubbletea"` if testing `Update`.
- **GOTCHA**:
  - Labels and caught status are styled by Lip Gloss. Visible text usually remains in the string, but ANSI may require stripping.
  - Current embedded Pokemon data may have empty `Types`, `Height`, and `Weight`; assert labels/status, not exact type/height values.
- **VALIDATE**: `gofmt -w internal/tui/detail_test.go`

### UPDATE `internal/cli/catch_test.go` if needed

- **IMPLEMENT**:
  - Do not merge catch-flow tests here unless the file stays readable.
  - If helper names collide, keep flow helpers in `catch_flow_test.go` with specific names.
- **PATTERN**: Existing parse/pick tests are already focused and should remain unchanged.
- **GOTCHA**: Avoid making the existing random tests brittle.
- **VALIDATE**: `gofmt -w internal/cli/catch_test.go`

---

## Testing Strategy

### Unit Tests

- `internal/pokemon`:
  - sprite API loads/decompresses embedded small regular and shiny sprites
  - missing sprite returns error
- `internal/tui`:
  - detail view content for normal caught / shiny missing
  - detail view fallback for missing Pokemon
  - key handling marks detail view done

### Integration-Style Unit Tests

- `internal/cli`:
  - catch flow runs in isolated temp config directory
  - default config bootstraps on first run
  - catch flow writes pokedex/trainer state
  - forced shiny produces deterministic shiny state and XP
  - invalid generation returns an error without writing catch state

### Edge Cases

- Random Pokemon should not make assertions flaky.
- `--shiny` must force shiny regardless of shiny rate.
- `--raw` should still register catch and write state.
- Invalid generation still may bootstrap config before returning error.
- ANSI styling in TUI should not make tests brittle.

---

## Validation Commands

> Project instructions say the assistant should not run tests/builds in this environment. The implementer/user should run these locally.

### Level 1: Syntax & Style

```bash
gofmt -w internal/cli/catch_flow_test.go internal/pokemon/sprites_test.go internal/tui/detail_test.go
```

### Level 2: Focused Tests

```bash
go test ./internal/cli ./internal/pokemon ./internal/tui
```

### Level 3: Full Test Suite

```bash
go test ./...
```

### Level 4: Manual Validation

```bash
XDG_CONFIG_HOME="$(mktemp -d)" ./bin/pokedex-fetch-go catch --gen 1 --shiny --raw >/dev/null
```

Then inspect generated state if desired:

```bash
find "$XDG_CONFIG_HOME/pokedex-fetch-go" -maxdepth 1 -type f -print
```

---

## Acceptance Criteria

- [ ] CLI catch flow test writes `config.toml`, `pokedex.json`, and `trainer.json` in a temp `XDG_CONFIG_HOME`.
- [ ] Forced shiny catch test verifies shiny catch count and shiny XP behavior.
- [ ] Invalid generation test verifies error behavior without catch state writes.
- [ ] Sprite tests verify small regular and shiny sprite loading through public API.
- [ ] Detail TUI tests verify visible name/number/caught status and missing Pokemon fallback.
- [ ] Tests avoid brittle full ANSI sprite snapshots.
- [ ] `go test ./...` passes locally.
