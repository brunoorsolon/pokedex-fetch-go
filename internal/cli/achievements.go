package cli

import (
	"github.com/brunoorsolon/pokedex-fetch-go/internal/tui"
	"github.com/spf13/cobra"
)

var achievementsCmd = &cobra.Command{
	Use:     "achievements",
	Aliases: []string{"hof", "hall-of-fame"},
	Short:   "Open the Hall of Fame achievement viewer",
	RunE: func(cmd *cobra.Command, args []string) error {
		return tui.RunAchievements()
	},
}

func init() {
	rootCmd.AddCommand(achievementsCmd)
}
