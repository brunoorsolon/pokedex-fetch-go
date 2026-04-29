package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// DefaultFileContent is the config.toml content written on first run and by
// `pokedex-fetch-go config init`.
const DefaultFileContent = `# pokedex-fetch-go configuration
# Place this file at ~/.config/pokedex-fetch-go/config.toml

# Shiny encounter rate (1 in N). Default: 1028
shiny_rate = 1028

# Boosted shiny rate for completed generations (1 in N). Default: 10
shiny_rate_boosted = 10

# Sprite size: "small" or "large". Default: "small"
sprite_size = "small"

# Show Pokemon name with sprite. Default: true
show_name = true

# Generations to pick from. Default: all
generations = [1, 2, 3, 4, 5, 6, 7, 8]

# XP awards per catch
xp_normal = 10
xp_shiny = 150

# Pokemon to never show (by name)
excluded = []
`

// Config holds all user-configurable settings.
type Config struct {
	ShinyRate        int      `toml:"shiny_rate"`
	ShinyRateBoosted int      `toml:"shiny_rate_boosted"`
	SpriteSize       string   `toml:"sprite_size"`
	ShowName         bool     `toml:"show_name"`
	Generations      []int    `toml:"generations"`
	XPNormal         int      `toml:"xp_normal"`
	XPShiny          int      `toml:"xp_shiny"`
	Excluded         []string `toml:"excluded"`
}

// Default returns a Config populated with sensible defaults.
func Default() *Config {
	return &Config{
		ShinyRate:        1028,
		ShinyRateBoosted: 10,
		SpriteSize:       "small",
		ShowName:         true,
		Generations:      []int{1, 2, 3, 4, 5, 6, 7, 8},
		XPNormal:         10,
		XPShiny:          150,
		Excluded:         []string{},
	}
}

// DataDir returns the XDG-compliant config directory path.
func DataDir() string {
	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		return filepath.Join(xdg, "pokedex-fetch-go")
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "pokedex-fetch-go")
}

// ConfigPath returns the full path to config.toml.
func ConfigPath() string {
	return filepath.Join(DataDir(), "config.toml")
}

// EnsureDefault writes the default config.toml if it does not already exist.
func EnsureDefault() (string, bool, error) {
	path := ConfigPath()
	if _, err := os.Stat(path); err == nil {
		return path, false, nil
	} else if !os.IsNotExist(err) {
		return path, false, err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return path, false, fmt.Errorf("create config dir: %w", err)
	}
	if err := os.WriteFile(path, []byte(DefaultFileContent), 0644); err != nil {
		return path, false, fmt.Errorf("write config: %w", err)
	}
	return path, true, nil
}

// Load reads config.toml from the config directory, bootstrapping defaults if absent.
func Load() (*Config, error) {
	cfg := Default()
	path, _, err := EnsureDefault()
	if err != nil {
		return cfg, err
	}
	_, err = toml.DecodeFile(path, cfg)
	return cfg, err
}

// IsExcluded returns true if the given Pokemon name is in the exclusion list.
func (c *Config) IsExcluded(name string) bool {
	for _, ex := range c.Excluded {
		if ex == name {
			return true
		}
	}
	return false
}
