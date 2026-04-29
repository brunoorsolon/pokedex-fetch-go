package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "pokedex-fetch-go",
	Short: "Pokemon terminal sprite display and collection tracker",
	Long:  "Display random Pokemon sprites in your terminal, track catches in a persistent Pokedex, and integrate with fastfetch.",
	RunE:  runCatch,
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().String("gen", "", `Generation(s) to pick from (e.g. "1", "1-3", "1,3,6")`)
	rootCmd.PersistentFlags().Bool("shiny", false, "Force shiny variant")
	rootCmd.PersistentFlags().Bool("big", false, "Use large sprite")
	rootCmd.PersistentFlags().Bool("no-title", false, "Don't display Pokemon name")
	rootCmd.PersistentFlags().Bool("raw", false, "Output sprite only (for fastfetch piping)")
}
