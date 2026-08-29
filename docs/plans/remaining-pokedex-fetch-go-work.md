---
title: Remaining pokedex-fetch-go Work
status: ready
created: 2026-05-06
source_plan: [[pokedex-fetch-go]]
---

# Remaining pokedex-fetch-go Work

## Status

The main `pokedex-fetch-go` implementation is mostly complete in `/source-repo`, but it is not fully aligned with `plans/pokedex-fetch-go.md`.

Completed plans moved to `/workspace/done/`:

- `achievement-points-system.md`
- `meaningful-integration-and-render-tests.md`

## Primary Gaps

### 1. Complete Pokemon metadata enrichment

Current `/source-repo/data/pokemon.json` has 905 Pokemon, but metadata is placeholder-only:

- `types` is empty for all entries.
- `height` and `weight` are `0` for all entries.
- all stats are `0`.
- `is_legendary` is always `false`.
- `is_mythical` is always `false`.

#### Tasks

1. Update `/source-repo/scripts/prepare_data.sh` or add a helper script to enrich data from PokeAPI or another reliable cached source.
2. Add a persistent cache directory, for example `/source-repo/data/cache/pokeapi/`, excluded from embed if needed.
3. Add explicit name mappings for known mismatches:
   - `nidoran-f`
   - `nidoran-m`
   - `farfetchd`
   - `sirfetchd`
   - `mr-mime`
   - `mime-jr`
   - `type-null`
   - `ho-oh`
   - `porygon-z`
   - `flabebe`
   - form names that do not map directly to PokeAPI endpoints
4. Populate each entry with:
   - `types`
   - `height`
   - `weight`
   - base stats
   - `is_legendary`
   - `is_mythical`
5. Keep generation assignment based on the local generation list/source used by the app, not PokeAPI species generation.
6. Add data validation checks to detect missing metadata before release.

#### Acceptance criteria

- All 905 Pokemon have at least one type.
- All 905 Pokemon have non-zero height and weight where source data exists.
- All 905 Pokemon have non-zero base stats.
- Legendary and mythical flags are populated correctly.
- Pokedex detail TUI shows real type/stat/height/weight data instead of placeholders.

---

### 2. Fix shiny-rate default mismatch

Plan says default shiny rate is `1/4096`, but current code uses `1028` in `/source-repo/internal/config/config.go` and generated config text.

#### Tasks

1. Decide whether the intended default is `4096` or `1028`.
2. If the plan is authoritative, update:
   - `config.Default()`
   - `DefaultFileContent`
   - `config.example.toml`
   - README docs if they mention the rate
3. Add or update config tests to assert the chosen default.

#### Acceptance criteria

- Code, example config, generated config, README, and plan all agree on the default shiny rate.

---

### 3. Implement real trainer setup / first-run UX

Current first catch silently creates a default trainer with `Name: "None"`. The original plan expects trainer profile creation to work on first run.

#### Tasks

1. Decide first-run behavior:
   - silent default trainer creation, or
   - explicit `trainer --name ...` setup, or
   - interactive setup when stdin is a TTY.
2. If using interactive setup, handle non-TTY contexts safely so `pokedex-fetch-go catch --raw` never blocks in shell/fastfetch usage.
3. Update `setup` or `trainer` docs to explain profile initialization.
4. Consider adding a `pokedex-fetch-go trainer init` or `pokedex-fetch-go setup profile` command.

#### Acceptance criteria

- First-run behavior is intentional and documented.
- `catch --raw` remains non-interactive and fastfetch-safe.
- Trainer profile can be created or updated without editing JSON manually.

---

### 4. Add state/config migration handling

The app now persists versioned JSON state, but load paths do not perform meaningful migrations.

#### Tasks

1. Define current schema versions for:
   - config TOML
   - `pokedex.json`
   - `trainer.json`
2. Add migration logic for old trainer achievement entries without IDs/points.
3. Recalculate `AchievementPoints` on load or migration.
4. Ensure missing newly-added config fields receive defaults.
5. Document migration behavior.

#### Acceptance criteria

- Existing user state from older versions loads safely.
- Old achievements are preserved or intentionally migrated.
- Missing fields do not reset user progress.

---

### 5. Add protection for concurrent catch writes

Atomic writes prevent corrupted files, but two simultaneous catches can still lose updates by read-modify-write clobbering.

#### Tasks

1. Decide whether concurrent catch support is required.
2. If required, add file locking around state load/update/save for `pokedex.json` and `trainer.json`.
3. If not required, document last-write-wins behavior.

#### Acceptance criteria

- Concurrent terminal opens either preserve both catches or the limitation is explicitly documented.

---

### 6. Add poketerm migration support or remove from scope

The original plan mentions a poketerm migration tool, but no migration command was found.

#### Tasks

1. Decide whether migration from poketerm remains in scope.
2. If yes, add a command such as:
   - `pokedex-fetch-go migrate poketerm --pokedex /path/to/pokedex.txt --trainer /path/to/trainer.json`
3. Convert old poketerm catch data into the new JSON state format.
4. Preserve trainer XP, streaks, profile fields, and hall-of-fame data where possible.
5. If not in scope, remove or defer it in docs/plans.

#### Acceptance criteria

- Existing poketerm users have a documented migration path, or migration is explicitly deferred.

---

### 7. Add backup/recovery behavior for state files

No backup/recovery mechanism was found for corrupted `pokedex.json` or `trainer.json`.

#### Tasks

1. Before overwriting state, optionally keep `.bak` copies.
2. Add recovery behavior when JSON unmarshal fails:
   - detect corrupted file
   - try `.bak`
   - report a clear error if recovery fails
3. Document manual recovery steps.

#### Acceptance criteria

- A corrupted state file does not silently destroy user progress.
- Users get a clear recovery path.

---

### 8. Tighten docs and release alignment

Current README is useful but does not cover all planned polish.

#### Tasks

1. Add the short alias recommendation:
   - `alias pfg='pokedex-fetch-go'`
2. Document first-run profile behavior.
3. Document state paths and whether they intentionally use `XDG_CONFIG_HOME` for state.
4. Document default shiny rate after resolving the mismatch.
5. Document migration/recovery status.

#### Acceptance criteria

- README accurately reflects current behavior.
- No documented feature is missing from the binary.

---

## Suggested User-run Validation

> Per project rules, the assistant should not run these commands in this environment.

```bash
go test ./...
```

```bash
go build -o bin/pokedex-fetch-go ./cmd/pokedex-fetch-go
```

```bash
XDG_CONFIG_HOME="$(mktemp -d)" ./bin/pokedex-fetch-go catch --gen 1 --raw >/dev/null
```

```bash
XDG_CONFIG_HOME="$(mktemp -d)" ./bin/pokedex-fetch-go show --name pikachu
```

```bash
XDG_CONFIG_HOME="$(mktemp -d)" ./bin/pokedex-fetch-go trainer
```

```bash
XDG_CONFIG_HOME="$(mktemp -d)" ./bin/pokedex-fetch-go achievements
```

## Final Acceptance Criteria

- [ ] Pokemon metadata is fully enriched, not placeholder-only.
- [ ] Shiny-rate default is consistent across code and docs.
- [ ] First-run trainer setup behavior is intentional and documented.
- [ ] State/config migrations are handled safely.
- [ ] Concurrent catch behavior is handled or documented.
- [ ] Poketerm migration is implemented or explicitly deferred.
- [ ] State backup/recovery exists or is explicitly deferred.
- [ ] README matches actual binary behavior.
