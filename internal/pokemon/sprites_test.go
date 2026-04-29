package pokemon

import (
	"strings"
	"testing"
)

func TestGetSpriteLoadsSmallRegularAndShiny(t *testing.T) {
	if !HasSprite("pikachu", SpriteSmall, SpriteRegular) {
		t.Fatalf("HasSprite(pikachu, small, regular) = false, want true")
	}
	if !HasSprite("pikachu", SpriteSmall, SpriteShiny) {
		t.Fatalf("HasSprite(pikachu, small, shiny) = false, want true")
	}

	regular, err := GetSprite("pikachu", SpriteSmall, SpriteRegular)
	if err != nil {
		t.Fatalf("GetSprite regular error = %v", err)
	}
	shiny, err := GetSprite("pikachu", SpriteSmall, SpriteShiny)
	if err != nil {
		t.Fatalf("GetSprite shiny error = %v", err)
	}
	if regular == "" {
		t.Fatalf("regular sprite is empty")
	}
	if shiny == "" {
		t.Fatalf("shiny sprite is empty")
	}
	if regular == shiny {
		t.Fatalf("regular and shiny sprites should differ")
	}
	if !strings.Contains(regular, "\x1b[") || !strings.Contains(shiny, "\x1b[") {
		t.Fatalf("sprites should contain ANSI escape sequences")
	}
}

func TestGetSpriteMissingPokemonReturnsError(t *testing.T) {
	if HasSprite("missingno", SpriteSmall, SpriteRegular) {
		t.Fatalf("HasSprite(missingno) = true, want false")
	}
	if sprite, err := GetSprite("missingno", SpriteSmall, SpriteRegular); err == nil {
		t.Fatalf("GetSprite(missingno) = %q, nil error; want error", sprite)
	}
}
