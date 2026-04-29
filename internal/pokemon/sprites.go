package pokemon

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"

	pfgdata "github.com/brunoorsolon/pokedex-fetch-go/data"
)

// GetSprite loads, decompresses, and returns the ANSI sprite string for the given Pokemon.
func GetSprite(name string, size SpriteSize, variant SpriteVariant) (string, error) {
	path := fmt.Sprintf("sprites/%s/%s/%s.gz", size, variant, name)
	compressed, err := pfgdata.SpritesFS.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("sprite not found: %s (%s/%s)", name, size, variant)
	}
	r, err := gzip.NewReader(bytes.NewReader(compressed))
	if err != nil {
		return "", fmt.Errorf("failed to decompress sprite: %w", err)
	}
	defer r.Close()
	raw, err := io.ReadAll(r)
	if err != nil {
		return "", fmt.Errorf("failed to read sprite: %w", err)
	}
	return string(raw), nil
}

// HasSprite returns true if a sprite file exists for the given parameters.
func HasSprite(name string, size SpriteSize, variant SpriteVariant) bool {
	path := fmt.Sprintf("sprites/%s/%s/%s.gz", size, variant, name)
	_, err := pfgdata.SpritesFS.ReadFile(path)
	return err == nil
}
