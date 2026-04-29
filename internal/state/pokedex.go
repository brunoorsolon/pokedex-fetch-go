package state

import (
	"encoding/json"
	"os"
	"time"
)

// PokedexState tracks all caught Pokemon.
type PokedexState struct {
	Pokemon map[string]*CatchRecord `json:"pokemon"`
	Version int                     `json:"version"`
}

// CatchRecord holds normal and shiny catch details for one Pokemon species.
type CatchRecord struct {
	Normal *CatchDetail `json:"normal,omitempty"`
	Shiny  *CatchDetail `json:"shiny,omitempty"`
}

// CatchDetail tracks how many times and when a Pokemon was first caught.
type CatchDetail struct {
	Count       int    `json:"count"`
	FirstCaught string `json:"first_caught"`
}

// NewPokedex returns an empty PokedexState.
func NewPokedex() (*PokedexState, error) {
	return &PokedexState{
		Pokemon: make(map[string]*CatchRecord),
		Version: 1,
	}, nil
}

// LoadPokedex reads pokedex.json from disk, or returns a fresh state if absent.
func LoadPokedex() (*PokedexState, error) {
	data, err := os.ReadFile(PokedexPath())
	if os.IsNotExist(err) {
		return NewPokedex()
	}
	if err != nil {
		return nil, err
	}
	var s PokedexState
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, err
	}
	if s.Pokemon == nil {
		s.Pokemon = make(map[string]*CatchRecord)
	}
	return &s, nil
}

// RegisterCatch records a catch for the given Pokemon name.
func (p *PokedexState) RegisterCatch(name string, isShiny bool) {
	record, ok := p.Pokemon[name]
	if !ok {
		record = &CatchRecord{}
		p.Pokemon[name] = record
	}
	today := time.Now().Format("2006-01-02")
	if isShiny {
		if record.Shiny == nil {
			record.Shiny = &CatchDetail{Count: 0, FirstCaught: today}
		}
		record.Shiny.Count++
	} else {
		if record.Normal == nil {
			record.Normal = &CatchDetail{Count: 0, FirstCaught: today}
		}
		record.Normal.Count++
	}
}

// Save writes the pokedex state to disk atomically.
func (p *PokedexState) Save() error {
	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return err
	}
	return AtomicWrite(PokedexPath(), data)
}

// IsCaught returns true if the Pokemon has been caught at least once (normal or shiny).
func (p *PokedexState) IsCaught(name string) bool {
	r, ok := p.Pokemon[name]
	if !ok {
		return false
	}
	return (r.Normal != nil && r.Normal.Count > 0) || (r.Shiny != nil && r.Shiny.Count > 0)
}

// IsGenComplete returns true if every Pokemon name in genNames has been caught.
func (p *PokedexState) IsGenComplete(genNames []string) bool {
	for _, name := range genNames {
		if !p.IsCaught(name) {
			return false
		}
	}
	return len(genNames) > 0
}

// UniqueCaught returns the number of distinct Pokemon species caught.
func (p *PokedexState) UniqueCaught() int {
	count := 0
	for _, r := range p.Pokemon {
		if (r.Normal != nil && r.Normal.Count > 0) || (r.Shiny != nil && r.Shiny.Count > 0) {
			count++
		}
	}
	return count
}

// TotalShinies returns the total number of shiny catches across all species.
func (p *PokedexState) TotalShinies() int {
	total := 0
	for _, r := range p.Pokemon {
		if r.Shiny != nil {
			total += r.Shiny.Count
		}
	}
	return total
}

// TotalCatches returns the sum of all normal and shiny catch counts.
func (p *PokedexState) TotalCatches() int {
	total := 0
	for _, r := range p.Pokemon {
		if r.Normal != nil {
			total += r.Normal.Count
		}
		if r.Shiny != nil {
			total += r.Shiny.Count
		}
	}
	return total
}

// HighestCatchCount returns the highest catch count for any single species.
func (p *PokedexState) HighestCatchCount() int {
	max := 0
	for _, r := range p.Pokemon {
		n := 0
		if r.Normal != nil {
			n += r.Normal.Count
		}
		if r.Shiny != nil {
			n += r.Shiny.Count
		}
		if n > max {
			max = n
		}
	}
	return max
}

// GenCaughtCount returns how many Pokemon from genNames have been caught.
func (p *PokedexState) GenCaughtCount(genNames []string) int {
	count := 0
	for _, name := range genNames {
		if p.IsCaught(name) {
			count++
		}
	}
	return count
}

// CaughtSet returns the set of all caught Pokemon names.
func (p *PokedexState) CaughtSet() map[string]bool {
	s := make(map[string]bool, len(p.Pokemon))
	for name, r := range p.Pokemon {
		if (r.Normal != nil && r.Normal.Count > 0) || (r.Shiny != nil && r.Shiny.Count > 0) {
			s[name] = true
		}
	}
	return s
}
