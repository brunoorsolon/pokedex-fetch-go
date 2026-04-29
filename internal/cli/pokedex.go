package cli

import (
	"fmt"
	"strconv"

	"github.com/brunoorsolon/pokedex-fetch-go/internal/tui"
	"github.com/spf13/cobra"
)

var pokedexCmd = &cobra.Command{
	Use:   "pokedex [gen]",
	Short: "Open the interactive Pokedex viewer",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		gen := 1
		if len(args) == 1 {
			parsed, err := strconv.Atoi(args[0])
			if err != nil || parsed < 1 || parsed > 8 {
				return fmt.Errorf("invalid generation %q: must be 1-8", args[0])
			}
			gen = parsed
		}
		return tui.RunPokedex(gen)
	},
}

func init() {
	rootCmd.AddCommand(pokedexCmd)
}
