package config

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestDefaultConfigGameplayDefaults(t *testing.T) {
	cfg := Default()

	if cfg.ShinyRate != 1028 {
		t.Fatalf("ShinyRate = %d, want 1028", cfg.ShinyRate)
	}
	if cfg.ShinyRateBoosted != 10 {
		t.Fatalf("ShinyRateBoosted = %d, want 10", cfg.ShinyRateBoosted)
	}
	if cfg.SpriteSize != "small" || !cfg.ShowName {
		t.Fatalf("sprite defaults = size %q show %v, want small/true", cfg.SpriteSize, cfg.ShowName)
	}
	wantGens := []int{1, 2, 3, 4, 5, 6, 7, 8}
	if !reflect.DeepEqual(cfg.Generations, wantGens) {
		t.Fatalf("Generations = %#v, want %#v", cfg.Generations, wantGens)
	}
	if cfg.XPNormal != 10 || cfg.XPShiny != 150 {
		t.Fatalf("XP defaults = normal %d shiny %d, want 10/150", cfg.XPNormal, cfg.XPShiny)
	}
}

func TestIsExcluded(t *testing.T) {
	cfg := &Config{Excluded: []string{"pikachu", "mewtwo"}}

	if !cfg.IsExcluded("pikachu") {
		t.Fatalf("IsExcluded(pikachu) = false, want true")
	}
	if cfg.IsExcluded("eevee") {
		t.Fatalf("IsExcluded(eevee) = true, want false")
	}
	if cfg.IsExcluded("Pikachu") {
		t.Fatalf("IsExcluded should be exact/canonical-name based")
	}
}

func TestEnsureDefaultBootstrapsConfigFile(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", tmp)

	path, created, err := EnsureDefault()
	if err != nil {
		t.Fatalf("EnsureDefault() error = %v", err)
	}
	if !created {
		t.Fatalf("EnsureDefault() created = false, want true")
	}
	wantPath := filepath.Join(tmp, "pokedex-fetch-go", "config.toml")
	if path != wantPath {
		t.Fatalf("config path = %q, want %q", path, wantPath)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile(config) error = %v", err)
	}
	if !strings.Contains(string(data), "shiny_rate = 1028") {
		t.Fatalf("bootstrapped config missing shiny default: %s", string(data))
	}

	_, created, err = EnsureDefault()
	if err != nil {
		t.Fatalf("second EnsureDefault() error = %v", err)
	}
	if created {
		t.Fatalf("second EnsureDefault() created = true, want false")
	}
}

func TestLoadReadsUserConfigAndKeepsDefaultsForMissingFields(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	if err := os.MkdirAll(DataDir(), 0755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	custom := `shiny_rate = 7
sprite_size = "large"
generations = [1, 3]
excluded = ["pikachu"]
`
	if err := os.WriteFile(ConfigPath(), []byte(custom), 0644); err != nil {
		t.Fatalf("WriteFile(config) error = %v", err)
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.ShinyRate != 7 || cfg.SpriteSize != "large" {
		t.Fatalf("loaded overrides = shiny %d size %q", cfg.ShinyRate, cfg.SpriteSize)
	}
	if !reflect.DeepEqual(cfg.Generations, []int{1, 3}) {
		t.Fatalf("loaded generations = %#v", cfg.Generations)
	}
	if !cfg.IsExcluded("pikachu") {
		t.Fatalf("loaded excluded list does not contain pikachu")
	}
	if cfg.ShinyRateBoosted != 10 || cfg.XPNormal != 10 || cfg.XPShiny != 150 || !cfg.ShowName {
		t.Fatalf("missing fields did not keep defaults: %#v", cfg)
	}
}
