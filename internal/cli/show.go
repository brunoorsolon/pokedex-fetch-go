package cli

import (
	"fmt"

	"github.com/brunoorsolon/pokedex-fetch-go/internal/config"
	"github.com/brunoorsolon/pokedex-fetch-go/internal/pokemon"
	"github.com/spf13/cobra"
)

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Show a specific Pokemon by name or number",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, _ := config.Load()

		name, _ := cmd.Flags().GetString("name")
		number, _ := cmd.Flags().GetInt("number")
		form, _ := cmd.Flags().GetString("form")

		var poke *pokemon.Pokemon
		var ok bool
		switch {
		case name != "":
			poke, ok = pokemon.GetByName(name)
		case number > 0:
			poke, ok = pokemon.GetByNumber(number)
		default:
			return fmt.Errorf("specify --name or --number")
		}
		if !ok {
			return fmt.Errorf("Pokemon not found")
		}

		spriteName := poke.Name
		if form != "" {
			spriteName = poke.Name + "-" + form
			if !pokemon.HasSprite(spriteName, pokemon.SpriteSmall, pokemon.SpriteRegular) {
				return fmt.Errorf("form '%s' not found for %s", form, poke.Name)
			}
		}

		isShiny, _ := cmd.Flags().GetBool("shiny")
		size := pokemon.SpriteSize(cfg.SpriteSize)
		if big, _ := cmd.Flags().GetBool("big"); big {
			size = pokemon.SpriteLarge
		}
		variant := pokemon.SpriteRegular
		if isShiny {
			variant = pokemon.SpriteShiny
		}

		sprite, err := pokemon.GetSprite(spriteName, size, variant)
		if err != nil {
			return err
		}

		raw, _ := cmd.Flags().GetBool("raw")
		noTitle, _ := cmd.Flags().GetBool("no-title")
		if !raw && !noTitle {
			if isShiny {
				fmt.Printf("\033[33m%s (shiny)\033[0m\n", spriteName)
			} else {
				fmt.Println(spriteName)
			}
		}
		fmt.Print(sprite)
		return nil
	},
}

func init() {
	showCmd.Flags().StringP("name", "n", "", "Pokemon name")
	showCmd.Flags().IntP("number", "N", 0, "Pokemon number")
	showCmd.Flags().StringP("form", "f", "", "Alternate form (e.g. mega, gmax)")
	rootCmd.AddCommand(showCmd)
}
