package cli

import (
	"fmt"

	"github.com/brunoorsolon/pokedex-fetch-go/internal/config"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration",
}

var configInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Create default config file",
	RunE: func(cmd *cobra.Command, args []string) error {
		path, created, err := config.EnsureDefault()
		if err != nil {
			return err
		}
		if created {
			fmt.Printf("Config created: %s\n", path)
		} else {
			fmt.Printf("Config already exists: %s\n", path)
		}
		return nil
	},
}

var configPathCmd = &cobra.Command{
	Use:   "path",
	Short: "Print config directory path",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(config.DataDir())
	},
}

func init() {
	configCmd.AddCommand(configInitCmd)
	configCmd.AddCommand(configPathCmd)
	rootCmd.AddCommand(configCmd)
}
