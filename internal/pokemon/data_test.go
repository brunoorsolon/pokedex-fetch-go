package pokemon

import "testing"

func TestPokemonLookupByNameAndNumber(t *testing.T) {
	pikachu, ok := GetByName(" PIKACHU ")
	if !ok {
		t.Fatalf("GetByName(PIKACHU) not found")
	}
	if pikachu.Name != "pikachu" || pikachu.Number != 25 || pikachu.Generation != 1 {
		t.Fatalf("pikachu data = %#v", pikachu)
	}
	if len(pikachu.Forms) == 0 || pikachu.Forms[0] != "regular" {
		t.Fatalf("pikachu forms = %#v, want regular first", pikachu.Forms)
	}

	byNumber, ok := GetByNumber(25)
	if !ok {
		t.Fatalf("GetByNumber(25) not found")
	}
	if byNumber.Name != pikachu.Name {
		t.Fatalf("GetByNumber(25).Name = %q, want %q", byNumber.Name, pikachu.Name)
	}
	if _, ok := GetByName("missingno"); ok {
		t.Fatalf("GetByName(missingno) found unexpected Pokemon")
	}
	if _, ok := GetByNumber(0); ok {
		t.Fatalf("GetByNumber(0) found unexpected Pokemon")
	}
}

func TestPokemonDataAndGenerationsAreLoaded(t *testing.T) {
	all := GetAll()
	if len(all) != 905 {
		t.Fatalf("GetAll() length = %d, want 905", len(all))
	}

	gen1, ok := GetGeneration(1)
	if !ok {
		t.Fatalf("GetGeneration(1) not found")
	}
	if gen1.Number != 1 || len(gen1.Names) == 0 || gen1.Names[0] != "bulbasaur" {
		t.Fatalf("generation 1 data = %#v", gen1)
	}
	gen8, ok := GetGeneration(8)
	if !ok || gen8.Number != 8 || len(gen8.Names) == 0 {
		t.Fatalf("generation 8 data = %#v found=%v", gen8, ok)
	}
	if _, ok := GetGeneration(9); ok {
		t.Fatalf("GetGeneration(9) found unexpected generation")
	}
}

func TestGetRandomFromGenerationsReturnsAllowedPokemon(t *testing.T) {
	for i := 0; i < 20; i++ {
		p := GetRandomFromGenerations([]int{1})
		if p == nil {
			t.Fatalf("GetRandomFromGenerations([1]) returned nil")
		}
		if p.Generation != 1 {
			t.Fatalf("random Pokemon generation = %d, want 1: %#v", p.Generation, p)
		}
	}
	if p := GetRandomFromGenerations([]int{99}); p != nil {
		t.Fatalf("GetRandomFromGenerations([99]) = %#v, want nil", p)
	}
}
