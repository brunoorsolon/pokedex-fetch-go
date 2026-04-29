package state

import (
	"os"
	"path/filepath"
)

// DataDir returns the XDG-compliant config directory for pokedex-fetch-go.
func DataDir() string {
	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		return filepath.Join(xdg, "pokedex-fetch-go")
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "pokedex-fetch-go")
}

// EnsureDataDir creates the data directory if it does not exist.
func EnsureDataDir() error {
	return os.MkdirAll(DataDir(), 0755)
}

// PokedexPath returns the full path to pokedex.json.
func PokedexPath() string {
	return filepath.Join(DataDir(), "pokedex.json")
}

// TrainerPath returns the full path to trainer.json.
func TrainerPath() string {
	return filepath.Join(DataDir(), "trainer.json")
}
