package cli

import (
	"reflect"
	"testing"

	"github.com/brunoorsolon/pokedex-fetch-go/internal/config"
)

func TestParseGenerations(t *testing.T) {
	tests := []struct {
		input string
		want  []int
	}{
		{input: "1", want: []int{1}},
		{input: " 8 ", want: []int{8}},
		{input: "1-3", want: []int{1, 2, 3}},
		{input: "1,3,6", want: []int{1, 3, 6}},
		{input: "1, 3, 6", want: []int{1, 3, 6}},
	}

	for _, tc := range tests {
		got, err := parseGenerations(tc.input)
		if err != nil {
			t.Fatalf("parseGenerations(%q) error = %v", tc.input, err)
		}
		if !reflect.DeepEqual(got, tc.want) {
			t.Fatalf("parseGenerations(%q) = %#v, want %#v", tc.input, got, tc.want)
		}
	}
}

func TestParseGenerationsRejectsInvalidInput(t *testing.T) {
	for _, input := range []string{"", "0", "9", "a", "3-1", "1-9", "1,a", "1,,2"} {
		if got, err := parseGenerations(input); err == nil {
			t.Fatalf("parseGenerations(%q) = %#v, nil error; want error", input, got)
		}
	}
}

func TestPickRandomRespectsExcludedPokemon(t *testing.T) {
	cfg := &config.Config{Excluded: []string{"bulbasaur", "ivysaur", "venusaur"}}
	for i := 0; i < 20; i++ {
		p := pickRandom([]int{1}, cfg)
		if p == nil {
			t.Fatalf("pickRandom returned nil with most gen 1 Pokemon available")
		}
		if cfg.IsExcluded(p.Name) {
			t.Fatalf("pickRandom returned excluded Pokemon %q", p.Name)
		}
	}
}

func TestPickRandomReturnsNilWhenNoCandidates(t *testing.T) {
	cfg := &config.Config{}
	if p := pickRandom([]int{99}, cfg); p != nil {
		t.Fatalf("pickRandom invalid generation = %#v, want nil", p)
	}
}
