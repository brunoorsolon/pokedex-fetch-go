package data

import "embed"

//go:embed pokemon.json
var PokemonJSON []byte

//go:embed all:gen
var GenFS embed.FS

//go:embed all:sprites
var SpritesFS embed.FS
