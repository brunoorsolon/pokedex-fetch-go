package data

import "testing"

func TestEmbeddedDataSmoke(t *testing.T) {
	if len(PokemonJSON) == 0 {
		t.Fatalf("PokemonJSON is empty")
	}
	if _, err := GenFS.ReadFile("gen/gen1.txt"); err != nil {
		t.Fatalf("embedded gen1 data missing: %v", err)
	}
	if _, err := SpritesFS.ReadFile("sprites/small/regular/pikachu.gz"); err != nil {
		t.Fatalf("embedded pikachu sprite missing: %v", err)
	}
}
