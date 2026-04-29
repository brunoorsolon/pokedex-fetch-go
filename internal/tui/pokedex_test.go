package tui

import "testing"

func TestDisplayPokemonName(t *testing.T) {
	tests := map[string]string{
		"":          "",
		"pikachu":   "Pikachu",
		"nidoran-f": "Nidoran♀",
		"nidoran-m": "Nidoran♂",
		"mr-mime":   "Mr. Mime",
		"mime-jr":   "Mime Jr.",
		"ho-oh":     "Ho-Oh",
		"porygon-z": "Porygon-Z",
		"type-null": "Type: Null",
		"tapu-koko": "Tapu Koko",
		"farfetchd": "Farfetch'd",
		"sirfetchd": "Sirfetch'd",
		"flabebe":   "Flabébé",
	}

	for input, want := range tests {
		if got := displayPokemonName(input); got != want {
			t.Fatalf("displayPokemonName(%q) = %q, want %q", input, got, want)
		}
	}
}
