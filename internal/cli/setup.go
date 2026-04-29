package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var setupCmd = &cobra.Command{
	Use:   "setup [bash|zsh|fish]",
	Short: "Print shell integration snippet",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		shell := "bash"
		if len(args) == 1 {
			shell = args[0]
		}
		switch shell {
		case "bash", "zsh":
			fmt.Println(`# pokedex-fetch-go + fastfetch
pokedex-fetch-go catch --raw | fastfetch --logo-type file-raw --logo -`)
		case "fish":
			fmt.Println(`# pokedex-fetch-go + fastfetch
pokedex-fetch-go catch --raw | fastfetch --logo-type file-raw --logo -`)
		default:
			return fmt.Errorf("unsupported shell %q (expected bash, zsh, or fish)", shell)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(setupCmd)
}
