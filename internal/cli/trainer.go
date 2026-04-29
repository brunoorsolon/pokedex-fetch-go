package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/brunoorsolon/pokedex-fetch-go/internal/pokemon"
	"github.com/brunoorsolon/pokedex-fetch-go/internal/state"
	"github.com/brunoorsolon/pokedex-fetch-go/internal/tui"
	"github.com/spf13/cobra"
)

var trainerCmd = &cobra.Command{
	Use:   "trainer",
	Short: "Show or edit trainer profile",
	RunE: func(cmd *cobra.Command, args []string) error {
		name, _ := cmd.Flags().GetString("name")
		region, _ := cmd.Flags().GetString("region")
		hometown, _ := cmd.Flags().GetString("hometown")
		favourites, _ := cmd.Flags().GetString("favourites")

		if name == "" && region == "" && hometown == "" && favourites == "" {
			return tui.RunTrainer()
		}

		tr, err := state.LoadTrainer()
		if err != nil {
			tr = state.DefaultTrainer()
		}
		if name != "" {
			tr.Name = name
		}
		if region != "" {
			if _, ok := state.HomeTowns[region]; !ok && region != "None" {
				return fmt.Errorf("unknown region %q", region)
			}
			tr.Region = region
		}
		if hometown != "" {
			tr.Hometown = hometown
		}
		if favourites != "" {
			parsed, err := parseFavouriteNumbers(favourites)
			if err != nil {
				return err
			}
			tr.Favourites = parsed
		}
		if err := tr.Save(); err != nil {
			return err
		}
		fmt.Println("Trainer profile updated")
		return nil
	},
}

func parseFavouriteNumbers(value string) ([]int, error) {
	parts := strings.Split(value, ",")
	result := make([]int, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		n, err := strconv.Atoi(part)
		if err != nil {
			return nil, fmt.Errorf("invalid favourite Pokemon number %q", part)
		}
		if _, ok := pokemon.GetByNumber(n); !ok {
			return nil, fmt.Errorf("Pokemon number %d not found", n)
		}
		result = append(result, n)
	}
	return result, nil
}

func init() {
	trainerCmd.Flags().String("name", "", "Set trainer name")
	trainerCmd.Flags().String("region", "", "Set home region (Kanto, Johto, Hoenn, Sinnoh, Unova, Kalos, Alola, Galar)")
	trainerCmd.Flags().String("hometown", "", "Set home town")
	trainerCmd.Flags().String("favourites", "", "Set favourite Pokemon numbers as comma-separated National Dex IDs")
	rootCmd.AddCommand(trainerCmd)
}
