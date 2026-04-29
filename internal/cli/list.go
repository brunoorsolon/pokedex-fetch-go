package cli

import (
	"fmt"

	"github.com/brunoorsolon/pokedex-fetch-go/internal/pokemon"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all Pokemon names",
	RunE: func(cmd *cobra.Command, args []string) error {
		genFlag, _ := cmd.Flags().GetString("gen")
		if genFlag != "" {
			gens, err := parseGenerations(genFlag)
			if err != nil {
				return err
			}
			for _, g := range gens {
				gen, ok := pokemon.GetGeneration(g)
				if !ok {
					continue
				}
				for _, name := range gen.Names {
					fmt.Println(name)
				}
			}
			return nil
		}
		for _, p := range pokemon.GetAll() {
			fmt.Println(p.Name)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
