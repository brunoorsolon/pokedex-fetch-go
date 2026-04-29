package state

import "testing"

func TestPokedexRegisterCatchAndDerivedStats(t *testing.T) {
	pdx, err := NewPokedex()
	if err != nil {
		t.Fatalf("NewPokedex() error = %v", err)
	}

	if pdx.IsCaught("bulbasaur") {
		t.Fatalf("empty pokedex reported bulbasaur as caught")
	}

	pdx.RegisterCatch("bulbasaur", false)
	pdx.RegisterCatch("bulbasaur", false)
	pdx.RegisterCatch("pikachu", true)
	pdx.RegisterCatch("pikachu", true)

	if !pdx.IsCaught("bulbasaur") || !pdx.IsCaught("pikachu") {
		t.Fatalf("registered catches were not reported as caught")
	}
	if pdx.IsCaught("charmander") {
		t.Fatalf("uncaught pokemon reported as caught")
	}
	if got := pdx.UniqueCaught(); got != 2 {
		t.Fatalf("UniqueCaught() = %d, want 2", got)
	}
	if got := pdx.TotalCatches(); got != 4 {
		t.Fatalf("TotalCatches() = %d, want 4", got)
	}
	if got := pdx.TotalShinies(); got != 2 {
		t.Fatalf("TotalShinies() = %d, want 2", got)
	}
	if got := pdx.HighestCatchCount(); got != 2 {
		t.Fatalf("HighestCatchCount() = %d, want 2", got)
	}
	if got := pdx.GenCaughtCount([]string{"bulbasaur", "pikachu", "charmander"}); got != 2 {
		t.Fatalf("GenCaughtCount() = %d, want 2", got)
	}
	if !pdx.IsGenComplete([]string{"bulbasaur", "pikachu"}) {
		t.Fatalf("IsGenComplete() = false, want true for caught set")
	}
	if pdx.IsGenComplete([]string{"bulbasaur", "charmander"}) {
		t.Fatalf("IsGenComplete() = true, want false with uncaught pokemon")
	}
	if pdx.IsGenComplete(nil) {
		t.Fatalf("IsGenComplete(nil) = true, want false")
	}

	caught := pdx.CaughtSet()
	if !caught["bulbasaur"] || !caught["pikachu"] || caught["charmander"] {
		t.Fatalf("CaughtSet() = %#v", caught)
	}
}

func TestPokedexSaveLoadRoundTrip(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())

	pdx, err := NewPokedex()
	if err != nil {
		t.Fatalf("NewPokedex() error = %v", err)
	}
	pdx.RegisterCatch("eevee", false)
	pdx.RegisterCatch("eevee", true)

	if err := pdx.Save(); err != nil {
		t.Fatalf("Save() error = %v", err)
	}
	loaded, err := LoadPokedex()
	if err != nil {
		t.Fatalf("LoadPokedex() error = %v", err)
	}

	if got := loaded.UniqueCaught(); got != 1 {
		t.Fatalf("loaded UniqueCaught() = %d, want 1", got)
	}
	if got := loaded.TotalCatches(); got != 2 {
		t.Fatalf("loaded TotalCatches() = %d, want 2", got)
	}
	if got := loaded.TotalShinies(); got != 1 {
		t.Fatalf("loaded TotalShinies() = %d, want 1", got)
	}
}
