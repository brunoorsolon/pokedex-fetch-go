package achievement

import "testing"

func TestCatalogInvariants(t *testing.T) {
	defs := Definitions()
	if len(defs) != 50 {
		t.Fatalf("Definitions() length = %d, want 50", len(defs))
	}

	validCategories := map[string]bool{
		CategoryPokedex: true,
		CategoryCatch:   true,
		CategoryRank:    true,
		CategorySpecial: true,
	}
	ids := make(map[string]bool, len(defs))
	names := make(map[string]bool, len(defs))
	total := 0
	for _, def := range defs {
		if def.ID == "" {
			t.Fatalf("empty ID in definition: %+v", def)
		}
		if ids[def.ID] {
			t.Fatalf("duplicate ID: %s", def.ID)
		}
		ids[def.ID] = true

		if def.Name == "" {
			t.Fatalf("empty name in definition: %+v", def)
		}
		if names[def.Name] {
			t.Fatalf("duplicate name: %s", def.Name)
		}
		names[def.Name] = true

		if def.Points <= 0 {
			t.Fatalf("non-positive points for %s: %d", def.ID, def.Points)
		}
		if !validCategories[def.Category] {
			t.Fatalf("invalid category for %s: %s", def.ID, def.Category)
		}
		total += def.Points
	}

	if total != TotalPoints {
		t.Fatalf("total points = %d, want %d", total, TotalPoints)
	}
}
