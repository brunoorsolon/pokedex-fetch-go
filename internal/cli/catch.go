package cli

import (
	"fmt"
	"math/rand/v2"
	"os"
	"strconv"
	"strings"

	"github.com/brunoorsolon/pokedex-fetch-go/internal/achievement"
	"github.com/brunoorsolon/pokedex-fetch-go/internal/config"
	"github.com/brunoorsolon/pokedex-fetch-go/internal/pokemon"
	"github.com/brunoorsolon/pokedex-fetch-go/internal/state"
	"github.com/spf13/cobra"
)

var catchCmd = &cobra.Command{
	Use:   "catch",
	Short: "Catch a random Pokemon (default action)",
	RunE:  runCatch,
}

func runCatch(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("config: %w", err)
	}

	// Determine generation filter
	gens := cfg.Generations
	if genFlag, _ := cmd.Flags().GetString("gen"); genFlag != "" {
		parsed, err := parseGenerations(genFlag)
		if err != nil {
			return err
		}
		gens = parsed
	}

	// Load pokedex state to check generation completion for boosted shiny rate
	pdx, err := state.LoadPokedex()
	if err != nil {
		pdx, _ = state.NewPokedex()
	}

	// Pick random Pokemon from allowed generations, excluding excluded list
	poke := pickRandom(gens, cfg)
	if poke == nil {
		return fmt.Errorf("no Pokemon available for the specified generations")
	}

	// Determine shiny
	forceShiny, _ := cmd.Flags().GetBool("shiny")
	isShiny := forceShiny || rollShiny(poke, cfg, pdx)

	// Resolve sprite size and variant
	size := pokemon.SpriteSize(cfg.SpriteSize)
	if big, _ := cmd.Flags().GetBool("big"); big {
		size = pokemon.SpriteLarge
	}
	variant := pokemon.SpriteRegular
	if isShiny {
		variant = pokemon.SpriteShiny
	}

	sprite, err := pokemon.GetSprite(poke.Name, size, variant)
	if err != nil {
		// Fallback to small regular on error
		sprite, err = pokemon.GetSprite(poke.Name, pokemon.SpriteSmall, pokemon.SpriteRegular)
		if err != nil {
			return err
		}
		isShiny = false
	}

	// Display sprite (and optionally the name)
	raw, _ := cmd.Flags().GetBool("raw")
	noTitle, _ := cmd.Flags().GetBool("no-title")

	if !raw && !noTitle && cfg.ShowName {
		if isShiny {
			fmt.Printf("\033[33m%s (shiny)\033[0m\n", poke.Name)
		} else {
			fmt.Println(poke.Name)
		}
	}
	fmt.Print(sprite)

	// Register catch
	pdx.RegisterCatch(poke.Name, isShiny)
	if err := pdx.Save(); err != nil {
		fmt.Fprintf(os.Stderr, "warning: failed to save pokedex: %v\n", err)
	}

	// Update trainer (XP, streak, daily catch, achievements)
	tr, err := state.LoadTrainer()
	if err != nil {
		tr = state.DefaultTrainer()
	}
	xpAmount := cfg.XPNormal
	if isShiny {
		xpAmount = cfg.XPShiny
	}
	tr.AddXP(xpAmount)
	tr.UpdateStreak()
	tr.UpdateDailyCatch()
	achievement.CheckAll(pdx, tr, pokemon.GetAllGenerations())
	if err := tr.Save(); err != nil {
		fmt.Fprintf(os.Stderr, "warning: failed to save trainer: %v\n", err)
	}

	return nil
}

func rollShiny(poke *pokemon.Pokemon, cfg *config.Config, pdx *state.PokedexState) bool {
	rate := cfg.ShinyRate
	if rate < 1 {
		rate = config.Default().ShinyRate
	}
	boostedRate := cfg.ShinyRateBoosted
	if boostedRate < 1 {
		boostedRate = config.Default().ShinyRateBoosted
	}
	gen, ok := pokemon.GetGeneration(poke.Generation)
	if ok && pdx.IsGenComplete(gen.Names) {
		return rand.IntN(boostedRate) == 0
	}
	return rand.IntN(rate) == 0
}

func pickRandom(gens []int, cfg *config.Config) *pokemon.Pokemon {
	var candidates []*pokemon.Pokemon
	seen := make(map[string]bool)
	for _, g := range gens {
		gen, ok := pokemon.GetGeneration(g)
		if !ok {
			continue
		}
		for _, name := range gen.Names {
			if seen[name] || cfg.IsExcluded(name) {
				continue
			}
			seen[name] = true
			if p, ok := pokemon.GetByName(name); ok {
				candidates = append(candidates, p)
			}
		}
	}
	if len(candidates) == 0 {
		return nil
	}
	return candidates[rand.IntN(len(candidates))]
}

// parseGenerations parses "1", "1-3", or "1,3,6" into a slice of ints (1-8).
func parseGenerations(s string) ([]int, error) {
	s = strings.TrimSpace(s)
	var result []int

	// Comma-separated list: "1,3,6"
	if strings.Contains(s, ",") {
		for _, part := range strings.Split(s, ",") {
			n, err := strconv.Atoi(strings.TrimSpace(part))
			if err != nil || n < 1 || n > 8 {
				return nil, fmt.Errorf("invalid generation %q: must be 1-8", part)
			}
			result = append(result, n)
		}
		return result, nil
	}

	// Range: "1-3"
	if strings.Contains(s, "-") {
		parts := strings.SplitN(s, "-", 2)
		start, err1 := strconv.Atoi(strings.TrimSpace(parts[0]))
		end, err2 := strconv.Atoi(strings.TrimSpace(parts[1]))
		if err1 != nil || err2 != nil || start < 1 || end > 8 || start > end {
			return nil, fmt.Errorf("invalid generation range %q: must be N-M where 1 <= N <= M <= 8", s)
		}
		for i := start; i <= end; i++ {
			result = append(result, i)
		}
		return result, nil
	}

	// Single: "1"
	n, err := strconv.Atoi(s)
	if err != nil || n < 1 || n > 8 {
		return nil, fmt.Errorf("invalid generation %q: must be 1-8", s)
	}
	return []int{n}, nil
}

func init() {
	rootCmd.AddCommand(catchCmd)
}
