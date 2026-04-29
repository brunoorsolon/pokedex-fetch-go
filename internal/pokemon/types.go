package pokemon

// Pokemon represents a single Pokemon entry from the embedded data.
type Pokemon struct {
	Name        string   `json:"name"`
	Number      int      `json:"number"`
	Generation  int      `json:"generation"`
	Forms       []string `json:"forms"`
	Types       []string `json:"types"`
	Height      int      `json:"height"`
	Weight      int      `json:"weight"`
	Stats       Stats    `json:"stats"`
	IsLegendary bool     `json:"is_legendary"`
	IsMythical  bool     `json:"is_mythical"`
}

// Stats holds base stat values for a Pokemon.
type Stats struct {
	HP        int `json:"hp"`
	Attack    int `json:"attack"`
	Defense   int `json:"defense"`
	SpAttack  int `json:"sp_attack"`
	SpDefense int `json:"sp_defense"`
	Speed     int `json:"speed"`
}

// SpriteSize controls which sprite resolution to use.
type SpriteSize string

const (
	SpriteSmall SpriteSize = "small"
	SpriteLarge SpriteSize = "large"
)

// SpriteVariant controls whether to use the regular or shiny sprite.
type SpriteVariant string

const (
	SpriteRegular SpriteVariant = "regular"
	SpriteShiny   SpriteVariant = "shiny"
)

// Generation groups Pokemon by their generation number.
type Generation struct {
	Number int
	Names  []string // ordered list of Pokemon names in this generation
}
