package pokemon

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"strings"

	pfgdata "github.com/brunoorsolon/pokedex-fetch-go/data"
)

var (
	allPokemon    []Pokemon
	pokemonByName map[string]*Pokemon
	pokemonByNum  map[int]*Pokemon
	generations   []Generation
)

func init() {
	if err := json.Unmarshal(pfgdata.PokemonJSON, &allPokemon); err != nil {
		panic(fmt.Sprintf("failed to load pokemon data: %v", err))
	}
	pokemonByName = make(map[string]*Pokemon, len(allPokemon))
	pokemonByNum = make(map[int]*Pokemon, len(allPokemon))
	for i := range allPokemon {
		pokemonByName[allPokemon[i].Name] = &allPokemon[i]
		if allPokemon[i].Number > 0 {
			pokemonByNum[allPokemon[i].Number] = &allPokemon[i]
		}
	}
	loadGenerations()
}

func loadGenerations() {
	for g := 1; g <= 8; g++ {
		path := fmt.Sprintf("gen/gen%d.txt", g)
		data, err := pfgdata.GenFS.ReadFile(path)
		if err != nil {
			continue
		}
		var names []string
		scanner := bufio.NewScanner(bytes.NewReader(data))
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line != "" {
				names = append(names, line)
			}
		}
		generations = append(generations, Generation{Number: g, Names: names})
	}
}

// GetByName returns the Pokemon with the given name, if it exists.
func GetByName(name string) (*Pokemon, bool) {
	p, ok := pokemonByName[strings.ToLower(strings.TrimSpace(name))]
	return p, ok
}

// GetByNumber returns the Pokemon with the given National Dex number.
func GetByNumber(num int) (*Pokemon, bool) {
	p, ok := pokemonByNum[num]
	return p, ok
}

// GetGeneration returns the Generation struct for the given generation number (1-8).
func GetGeneration(gen int) (*Generation, bool) {
	for i := range generations {
		if generations[i].Number == gen {
			return &generations[i], true
		}
	}
	return nil, false
}

// GetAllGenerations returns all loaded generations in order.
func GetAllGenerations() []Generation {
	return generations
}

// GetAll returns all Pokemon in load order.
func GetAll() []Pokemon {
	return allPokemon
}

// GetRandomFromGenerations picks a random Pokemon from the specified generation numbers.
// Returns nil if no candidates are found.
func GetRandomFromGenerations(gens []int) *Pokemon {
	var candidates []*Pokemon
	for _, g := range gens {
		gen, ok := GetGeneration(g)
		if !ok {
			continue
		}
		for _, name := range gen.Names {
			if p, ok := pokemonByName[name]; ok {
				candidates = append(candidates, p)
			}
		}
	}
	if len(candidates) == 0 {
		return nil
	}
	return candidates[rand.IntN(len(candidates))]
}
