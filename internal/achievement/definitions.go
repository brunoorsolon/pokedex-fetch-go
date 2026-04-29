package achievement

// SpecialCollection defines a named set of Pokemon that must all be caught.
type SpecialCollection struct {
	ID       string
	Name     string
	Event    string
	Category string
	Points   int
	Required map[string]bool
}

var specialCollections = []SpecialCollection{
	{
		ID: "special_legendaries", Name: "Legendary Conqueror", Event: "Caught All Legendary Pokemon", Category: CategorySpecial, Points: 50,
		Required: setOf(
			"articuno", "zapdos", "moltres", "mewtwo",
			"raikou", "entei", "suicune", "lugia", "ho-oh", "celebi",
			"regirock", "regice", "registeel", "latias", "latios",
			"kyogre", "groudon", "rayquaza", "jirachi", "deoxys",
			"uxie", "mesprit", "azelf", "dialga", "palkia", "heatran",
			"regigigas", "giratina", "cresselia", "phione", "manaphy",
			"darkrai", "shaymin", "arceus",
			"victini", "cobalion", "terrakion", "virizion", "tornadus",
			"thundurus", "reshiram", "zekrom", "landorus", "kyurem",
			"keldeo", "meloetta", "genesect",
			"xerneas", "yveltal", "zygarde", "diancie", "hoopa", "volcanion",
			"type-null", "silvally", "tapu-koko", "tapu-lele", "tapu-bulu",
			"tapu-fini", "cosmog", "cosmoem", "solgaleo", "lunala",
			"necrozma", "magearna", "marshadow", "zeraora", "meltan", "melmetal",
			"zacian", "zamazenta", "eternatus", "kubfu", "urshifu",
			"zarude", "regieleki", "regidrago", "glastrier", "spectrier",
			"calyrex", "enamorus",
		),
	},
	{
		ID: "special_mythicals", Name: "Myth Hunter", Event: "Caught All Mythical Pokemon", Category: CategorySpecial, Points: 35,
		Required: setOf(
			"mew", "celebi", "jirachi", "deoxys",
			"phione", "manaphy", "darkrai", "shaymin", "arceus",
			"victini", "keldeo", "meloetta", "genesect",
			"diancie", "hoopa", "volcanion",
			"magearna", "marshadow", "zeraora", "meltan", "melmetal", "zarude",
		),
	},
	{
		ID: "special_kanto_birds", Name: "Legendary Trio Master - Kanto Birds", Event: "Completed the Legendary Birds", Category: CategorySpecial, Points: 10,
		Required: setOf("articuno", "zapdos", "moltres"),
	},
	{
		ID: "special_weather_trio", Name: "Weather Dominator", Event: "Completed the Weather Trio", Category: CategorySpecial, Points: 15,
		Required: setOf("kyogre", "groudon", "rayquaza"),
	},
	{
		ID: "special_starters", Name: "Starter Supreme", Event: "Caught All Starter Pokemon", Category: CategorySpecial, Points: 25,
		Required: setOf(
			"bulbasaur", "charmander", "squirtle",
			"chikorita", "cyndaquil", "totodile",
			"treecko", "torchic", "mudkip",
			"turtwig", "chimchar", "piplup",
			"snivy", "tepig", "oshawott",
			"chespin", "fennekin", "froakie",
			"rowlet", "litten", "popplio",
			"grookey", "scorbunny", "sobble",
		),
	},
	{
		ID: "special_pseudo_legendaries", Name: "Pseudo-Legend Slayer", Event: "Caught All Pseudo-Legendaries", Category: CategorySpecial, Points: 20,
		Required: setOf(
			"dragonite", "tyranitar", "salamence", "metagross",
			"garchomp", "hydreigon", "goodra", "kommo-o", "dragapult",
		),
	},
	{
		ID: "special_eeveelutions", Name: "Eeveelution Enthusiast", Event: "Caught All Eeveelutions", Category: CategorySpecial, Points: 15,
		Required: setOf(
			"eevee", "vaporeon", "jolteon", "flareon",
			"espeon", "umbreon", "leafeon", "glaceon", "sylveon",
		),
	},
	{
		ID: "special_fossils", Name: "Paleontologist", Event: "Caught All Fossil Pokemon", Category: CategorySpecial, Points: 20,
		Required: setOf(
			"omanyte", "kabuto", "aerodactyl",
			"lileep", "anorith",
			"cranidos", "shieldon",
			"tirtouga", "archen",
			"tyrunt", "amaura",
			"dracozolt", "arctozolt", "dracovish", "arctovish",
		),
	},
}

func setOf(names ...string) map[string]bool {
	m := make(map[string]bool, len(names))
	for _, n := range names {
		m[n] = true
	}
	return m
}

// SpecialCollections returns the special collections list.
func SpecialCollections() []SpecialCollection {
	return specialCollections
}
