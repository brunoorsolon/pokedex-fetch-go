package state

import "testing"

func TestTrainerAddXPLevelsAndTracksLifetimeXP(t *testing.T) {
	tr := DefaultTrainer()

	tr.AddXP(999)
	if tr.Level != 1 || tr.XP != 999 || tr.TotalXP != 999 {
		t.Fatalf("after 999 XP: level=%d xp=%d total=%d, want level=1 xp=999 total=999", tr.Level, tr.XP, tr.TotalXP)
	}

	tr.AddXP(1)
	if tr.Level != 2 || tr.XP != 0 || tr.TotalXP != 1000 {
		t.Fatalf("after level-up: level=%d xp=%d total=%d, want level=2 xp=0 total=1000", tr.Level, tr.XP, tr.TotalXP)
	}

	tr.AddXP(2500)
	if tr.Level != 3 || tr.XP != 1000 || tr.TotalXP != 3500 {
		t.Fatalf("after chained XP: level=%d xp=%d total=%d, want level=3 xp=1000 total=3500", tr.Level, tr.XP, tr.TotalXP)
	}
}

func TestTrainerUnlockIsIdempotentAndRecalculatesPoints(t *testing.T) {
	tr := DefaultTrainer()

	if !tr.Unlock("first", "First", "Did first thing", "test", 10) {
		t.Fatalf("first Unlock() = false, want true")
	}
	if tr.Unlock("first", "First renamed", "Did first thing again", "test", 99) {
		t.Fatalf("duplicate Unlock() = true, want false")
	}
	if !tr.Unlock("second", "Second", "Did second thing", "test", 15) {
		t.Fatalf("second Unlock() = false, want true")
	}

	if got := len(tr.HallOfFame); got != 2 {
		t.Fatalf("HallOfFame length = %d, want 2", got)
	}
	if got := tr.AchievementPoints; got != 25 {
		t.Fatalf("AchievementPoints = %d, want 25", got)
	}
	if tr.HallOfFame[0].ID != "first" || tr.HallOfFame[0].Points != 10 || tr.HallOfFame[0].Date == "" {
		t.Fatalf("first achievement not persisted correctly: %#v", tr.HallOfFame[0])
	}

	tr.AchievementPoints = 0
	tr.RecalculateAchievementPoints()
	if got := tr.AchievementPoints; got != 25 {
		t.Fatalf("recalculated AchievementPoints = %d, want 25", got)
	}
}

func TestTrainerSaveLoadRoundTrip(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())

	tr := DefaultTrainer()
	tr.Name = "Ash"
	tr.Region = "Kanto"
	tr.Hometown = "Pallet Town"
	tr.AddXP(1000)
	tr.Unlock("pokedex_first_catch", "Just getting started", "Caught your first Pokemon", "pokedex", 5)

	if err := tr.Save(); err != nil {
		t.Fatalf("Save() error = %v", err)
	}
	loaded, err := LoadTrainer()
	if err != nil {
		t.Fatalf("LoadTrainer() error = %v", err)
	}

	if loaded.Name != "Ash" || loaded.Region != "Kanto" || loaded.Hometown != "Pallet Town" {
		t.Fatalf("loaded profile = %#v", loaded)
	}
	if loaded.Level != 2 || loaded.TotalXP != 1000 {
		t.Fatalf("loaded progression level=%d totalXP=%d, want level=2 totalXP=1000", loaded.Level, loaded.TotalXP)
	}
	if len(loaded.HallOfFame) != 1 || loaded.HallOfFame[0].ID != "pokedex_first_catch" || loaded.AchievementPoints != 5 {
		t.Fatalf("loaded achievements = %#v points=%d", loaded.HallOfFame, loaded.AchievementPoints)
	}
}
