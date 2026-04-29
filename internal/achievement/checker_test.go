package achievement

import (
	"testing"

	"github.com/brunoorsolon/pokedex-fetch-go/internal/state"
)

func TestCheckAllUnlocksFirstCatchWithPoints(t *testing.T) {
	pdx := mustPokedex(t)
	pdx.RegisterCatch("bulbasaur", false)
	tr := state.DefaultTrainer()

	result := CheckAll(pdx, tr, nil)

	if len(result.NewUnlocks) != 1 {
		t.Fatalf("new unlocks = %d, want 1", len(result.NewUnlocks))
	}
	if !hasAchievement(tr, "pokedex_first_catch") {
		t.Fatalf("missing pokedex_first_catch achievement")
	}
	if tr.AchievementPoints != 5 {
		t.Fatalf("achievement points = %d, want 5", tr.AchievementPoints)
	}
	if tr.HallOfFame[0].Date == "" {
		t.Fatalf("achievement date was not set")
	}
}

func TestCheckAllIsIdempotent(t *testing.T) {
	pdx := mustPokedex(t)
	pdx.RegisterCatch("bulbasaur", false)
	tr := state.DefaultTrainer()

	CheckAll(pdx, tr, nil)
	CheckAll(pdx, tr, nil)

	if len(tr.HallOfFame) != 1 {
		t.Fatalf("hall of fame length = %d, want 1", len(tr.HallOfFame))
	}
	if tr.AchievementPoints != 5 {
		t.Fatalf("achievement points = %d, want 5", tr.AchievementPoints)
	}
}

func TestCompletionistRequiresAllOtherAchievements(t *testing.T) {
	pdx := mustPokedex(t)
	tr := state.DefaultTrainer()
	missingID := "pokedex_first_catch"

	for _, def := range Definitions() {
		if def.ID == CompletionistID || def.ID == missingID {
			continue
		}
		tr.Unlock(def.ID, def.Name, def.Event, def.Category, def.Points)
	}

	CheckAll(pdx, tr, nil)
	if hasAchievement(tr, CompletionistID) {
		t.Fatalf("completionist unlocked before all other achievements")
	}

	missing, ok := DefinitionByID(missingID)
	if !ok {
		t.Fatalf("missing test definition %s", missingID)
	}
	tr.Unlock(missing.ID, missing.Name, missing.Event, missing.Category, missing.Points)
	CheckAll(pdx, tr, nil)

	if !hasAchievement(tr, CompletionistID) {
		t.Fatalf("completionist did not unlock after all other achievements")
	}
	if tr.AchievementPoints != TotalPoints {
		t.Fatalf("achievement points = %d, want %d", tr.AchievementPoints, TotalPoints)
	}
}

func mustPokedex(t *testing.T) *state.PokedexState {
	t.Helper()
	pdx, err := state.NewPokedex()
	if err != nil {
		t.Fatalf("NewPokedex error: %v", err)
	}
	return pdx
}

func hasAchievement(tr *state.TrainerState, id string) bool {
	for _, a := range tr.HallOfFame {
		if a.ID == id {
			return true
		}
	}
	return false
}
