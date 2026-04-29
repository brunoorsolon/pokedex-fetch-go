package achievement

import (
	"github.com/brunoorsolon/pokedex-fetch-go/internal/pokemon"
	"github.com/brunoorsolon/pokedex-fetch-go/internal/state"
)

// CheckResult describes achievements unlocked by one CheckAll run.
type CheckResult struct {
	NewUnlocks []AchievementDefinition
}

// CheckAll runs all achievement checks against the current state and unlocks
// any newly-earned achievements. Called on every catch — all checks are simple
// integer comparisons so it runs in microseconds.
func CheckAll(pdx *state.PokedexState, tr *state.TrainerState, gens []pokemon.Generation) CheckResult {
	result := CheckResult{}
	if pdx == nil || tr == nil {
		return result
	}

	tr.RecalculateAchievementPoints()
	ctx := buildContext(pdx, tr, gens)

	for _, rule := range nonCompletionistRules() {
		if rule.Condition != nil && rule.Condition(ctx) {
			def := rule.AchievementDefinition
			if tr.Unlock(def.ID, def.Name, def.Event, def.Category, def.Points) {
				result.NewUnlocks = append(result.NewUnlocks, def)
			}
		}
	}

	if allNonCompletionistUnlocked(tr) {
		def := completionistDefinition()
		if tr.Unlock(def.ID, def.Name, def.Event, def.Category, def.Points) {
			result.NewUnlocks = append(result.NewUnlocks, def)
		}
	}

	return result
}

func buildContext(pdx *state.PokedexState, tr *state.TrainerState, gens []pokemon.Generation) Context {
	genComplete := make(map[int]bool, len(gens))
	for _, gen := range gens {
		genComplete[gen.Number] = pdx.IsGenComplete(gen.Names)
	}

	return Context{
		Unique:            pdx.UniqueCaught(),
		TotalCatches:      pdx.TotalCatches(),
		TotalShinies:      pdx.TotalShinies(),
		HighestCatchCount: pdx.HighestCatchCount(),
		Level:             tr.Level,
		TotalXP:           tr.TotalXP,
		LongestStreak:     tr.Streak.Longest,
		DailyCount:        tr.DailyCatch.Count,
		CaughtSet:         pdx.CaughtSet(),
		GenComplete:       genComplete,
	}
}

func allNonCompletionistUnlocked(tr *state.TrainerState) bool {
	unlocked := make(map[string]bool, len(tr.HallOfFame))
	for _, a := range tr.HallOfFame {
		unlocked[a.ID] = true
	}
	for _, rule := range nonCompletionistRules() {
		if !unlocked[rule.ID] {
			return false
		}
	}
	return true
}
