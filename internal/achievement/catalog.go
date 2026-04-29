package achievement

const (
	CategoryPokedex = "pokedex"
	CategoryCatch   = "catch"
	CategoryRank    = "rank"
	CategorySpecial = "special"

	CompletionistID = "achievement_completionist"
	TotalPoints     = 1000
)

// AchievementDefinition describes a stable, user-visible achievement.
type AchievementDefinition struct {
	ID       string
	Name     string
	Event    string
	Category string
	Points   int
	Hidden   bool
}

// Context contains all derived state used by achievement rules.
type Context struct {
	Unique            int
	TotalCatches      int
	TotalShinies      int
	HighestCatchCount int
	Level             int
	TotalXP           int
	LongestStreak     int
	DailyCount        int
	CaughtSet         map[string]bool
	GenComplete       map[int]bool
}

type achievementRule struct {
	AchievementDefinition
	Condition func(Context) bool
}

var achievementRules = []achievementRule{
	{AchievementDefinition{ID: "pokedex_first_catch", Name: "Just getting started", Event: "Caught your first Pokemon", Category: CategoryPokedex, Points: 5}, func(c Context) bool { return c.Unique >= 1 }},
	{AchievementDefinition{ID: "pokedex_25_unique", Name: "Route 1 Regular", Event: "Caught 25 Pokemon", Category: CategoryPokedex, Points: 8}, func(c Context) bool { return c.Unique >= 25 }},
	{AchievementDefinition{ID: "pokedex_50_unique", Name: "You're getting the hang of this!", Event: "Caught 50 Pokemon", Category: CategoryPokedex, Points: 10}, func(c Context) bool { return c.Unique >= 50 }},
	{AchievementDefinition{ID: "pokedex_100_unique", Name: "A true pro", Event: "Caught 100 Pokemon", Category: CategoryPokedex, Points: 15}, func(c Context) bool { return c.Unique >= 100 }},
	{AchievementDefinition{ID: "pokedex_250_unique", Name: "Regional Researcher", Event: "Caught 250 Pokemon", Category: CategoryPokedex, Points: 25}, func(c Context) bool { return c.Unique >= 250 }},
	{AchievementDefinition{ID: "pokedex_all_905", Name: "You sure can open a terminal!", Event: "Caught all 905 Pokemon", Category: CategoryPokedex, Points: 50}, func(c Context) bool { return c.Unique >= 905 }},
	{AchievementDefinition{ID: "catch_1000_total", Name: "Frequent Flyer", Event: "Caught 1000 Pokemon total", Category: CategoryCatch, Points: 12}, func(c Context) bool { return c.TotalCatches >= 1000 }},
	{AchievementDefinition{ID: "catch_10000_total", Name: "How many terminals?", Event: "Caught 10000 Pokemon total", Category: CategoryCatch, Points: 25}, func(c Context) bool { return c.TotalCatches >= 10000 }},
	{AchievementDefinition{ID: "streak_3_days", Name: "Warming Up", Event: "3 Day Catch Streak", Category: CategoryCatch, Points: 5}, func(c Context) bool { return c.LongestStreak >= 3 }},
	{AchievementDefinition{ID: "streak_7_days", Name: "Weeklong Catcher", Event: "7 Day Catch Streak", Category: CategoryCatch, Points: 10}, func(c Context) bool { return c.LongestStreak >= 7 }},
	{AchievementDefinition{ID: "streak_30_days", Name: "No longer casual thing", Event: "30 Day Catch Streak", Category: CategoryCatch, Points: 20}, func(c Context) bool { return c.LongestStreak >= 30 }},
	{AchievementDefinition{ID: "streak_100_days", Name: "At this point, the Pokemon fear you", Event: "100 Day Catch Streak", Category: CategoryCatch, Points: 30}, func(c Context) bool { return c.LongestStreak >= 100 }},
	{AchievementDefinition{ID: "streak_365_days", Name: "You have not missed a single day", Event: "365 Day Catch Streak", Category: CategoryCatch, Points: 35}, func(c Context) bool { return c.LongestStreak >= 365 }},
	{AchievementDefinition{ID: "shiny_first", Name: "Oooo a Shiny!", Event: "Caught your first Shiny Pokemon", Category: CategoryCatch, Points: 15}, func(c Context) bool { return c.TotalShinies >= 1 }},
	{AchievementDefinition{ID: "shiny_5", Name: "You're really lucky!", Event: "Caught 5 Shiny Pokemon", Category: CategoryCatch, Points: 25}, func(c Context) bool { return c.TotalShinies >= 5 }},
	{AchievementDefinition{ID: "shiny_25", Name: "Sparkle Specialist", Event: "Caught 25 Shiny Pokemon", Category: CategoryCatch, Points: 35}, func(c Context) bool { return c.TotalShinies >= 25 }},
	{AchievementDefinition{ID: "shiny_100", Name: "Is this even luck anymore", Event: "Caught 100 Shiny Pokemon", Category: CategoryCatch, Points: 45}, func(c Context) bool { return c.TotalShinies >= 100 }},
	{AchievementDefinition{ID: "daily_5", Name: "You can almost make a full team!", Event: "Caught 5 Pokemon in a day", Category: CategoryCatch, Points: 5}, func(c Context) bool { return c.DailyCount >= 5 }},
	{AchievementDefinition{ID: "daily_10", Name: "You must like opening terminals", Event: "Caught 10 Pokemon in a day", Category: CategoryCatch, Points: 8}, func(c Context) bool { return c.DailyCount >= 10 }},
	{AchievementDefinition{ID: "daily_25", Name: "Terminal Marathon", Event: "Caught 25 Pokemon in a day", Category: CategoryCatch, Points: 12}, func(c Context) bool { return c.DailyCount >= 25 }},
	{AchievementDefinition{ID: "daily_50", Name: "Time to log off now!", Event: "Caught 50 Pokemon in a day", Category: CategoryCatch, Points: 20}, func(c Context) bool { return c.DailyCount >= 50 }},
	{AchievementDefinition{ID: "duplicate_10", Name: "I hope its not a weedle again!", Event: "Caught a Pokemon 10 times", Category: CategoryCatch, Points: 10}, func(c Context) bool { return c.HighestCatchCount >= 10 }},
	{AchievementDefinition{ID: "duplicate_25", Name: "Not another one!", Event: "Caught a Pokemon 25 times", Category: CategoryCatch, Points: 15}, func(c Context) bool { return c.HighestCatchCount >= 25 }},
	{AchievementDefinition{ID: "duplicate_50", Name: "Familiar Face", Event: "Caught a Pokemon 50 times", Category: CategoryCatch, Points: 20}, func(c Context) bool { return c.HighestCatchCount >= 50 }},
	{AchievementDefinition{ID: "level_10", Name: "Apprentice Trainer", Event: "Reach level 10", Category: CategoryRank, Points: 10}, func(c Context) bool { return c.Level >= 10 }},
	{AchievementDefinition{ID: "level_25", Name: "Veteran Trainer", Event: "Reach level 25", Category: CategoryRank, Points: 15}, func(c Context) bool { return c.Level >= 25 }},
	{AchievementDefinition{ID: "level_50", Name: "Master Trainer", Event: "Reach level 50", Category: CategoryRank, Points: 25}, func(c Context) bool { return c.Level >= 50 }},
	{AchievementDefinition{ID: "xp_1000", Name: "XP Initiate - The grind has begun", Event: "Reach a total XP of 1000", Category: CategoryRank, Points: 5}, func(c Context) bool { return c.TotalXP >= 1000 }},
	{AchievementDefinition{ID: "xp_10000", Name: "XP Adept - Momentum is undeniable", Event: "Reach a total XP of 10000", Category: CategoryRank, Points: 10}, func(c Context) bool { return c.TotalXP >= 10000 }},
	{AchievementDefinition{ID: "xp_50000", Name: "XP Elite - You're operating at scale", Event: "Reach a total XP of 50000", Category: CategoryRank, Points: 15}, func(c Context) bool { return c.TotalXP >= 50000 }},
	{AchievementDefinition{ID: "xp_100000", Name: "XP Grandmaster - The numbers fear you", Event: "Reach a total XP of 100000", Category: CategoryRank, Points: 20}, func(c Context) bool { return c.TotalXP >= 100000 }},
	{AchievementDefinition{ID: "xp_250000", Name: "XP Beyond Champion", Event: "Reach a total XP of 250000", Category: CategoryRank, Points: 20}, func(c Context) bool { return c.TotalXP >= 250000 }},
	{AchievementDefinition{ID: "gen_1_complete", Name: "Kanto Master", Event: "Completed Generation 1", Category: CategoryPokedex, Points: 20}, func(c Context) bool { return c.GenComplete[1] }},
	{AchievementDefinition{ID: "gen_2_complete", Name: "Johto Master", Event: "Completed Generation 2", Category: CategoryPokedex, Points: 15}, func(c Context) bool { return c.GenComplete[2] }},
	{AchievementDefinition{ID: "gen_3_complete", Name: "Hoenn Master", Event: "Completed Generation 3", Category: CategoryPokedex, Points: 20}, func(c Context) bool { return c.GenComplete[3] }},
	{AchievementDefinition{ID: "gen_4_complete", Name: "Sinnoh Master", Event: "Completed Generation 4", Category: CategoryPokedex, Points: 15}, func(c Context) bool { return c.GenComplete[4] }},
	{AchievementDefinition{ID: "gen_5_complete", Name: "Unova Master", Event: "Completed Generation 5", Category: CategoryPokedex, Points: 25}, func(c Context) bool { return c.GenComplete[5] }},
	{AchievementDefinition{ID: "gen_6_complete", Name: "Kalos Master", Event: "Completed Generation 6", Category: CategoryPokedex, Points: 10}, func(c Context) bool { return c.GenComplete[6] }},
	{AchievementDefinition{ID: "gen_7_complete", Name: "Alola Master", Event: "Completed Generation 7", Category: CategoryPokedex, Points: 15}, func(c Context) bool { return c.GenComplete[7] }},
	{AchievementDefinition{ID: "gen_8_complete", Name: "Galar Master", Event: "Completed Generation 8", Category: CategoryPokedex, Points: 20}, func(c Context) bool { return c.GenComplete[8] }},
	{AchievementDefinition{ID: "all_generations_complete", Name: "Lets collect some shinys now!", Event: "Completed all Generations", Category: CategoryPokedex, Points: 40}, allGenerationsComplete},
	{AchievementDefinition{ID: "special_kanto_birds", Name: "Legendary Trio Master - Kanto Birds", Event: "Completed the Legendary Birds", Category: CategorySpecial, Points: 10}, collectionComplete("special_kanto_birds")},
	{AchievementDefinition{ID: "special_weather_trio", Name: "Weather Dominator", Event: "Completed the Weather Trio", Category: CategorySpecial, Points: 15}, collectionComplete("special_weather_trio")},
	{AchievementDefinition{ID: "special_eeveelutions", Name: "Eeveelution Enthusiast", Event: "Caught All Eeveelutions", Category: CategorySpecial, Points: 15}, collectionComplete("special_eeveelutions")},
	{AchievementDefinition{ID: "special_pseudo_legendaries", Name: "Pseudo-Legend Slayer", Event: "Caught All Pseudo-Legendaries", Category: CategorySpecial, Points: 20}, collectionComplete("special_pseudo_legendaries")},
	{AchievementDefinition{ID: "special_fossils", Name: "Paleontologist", Event: "Caught All Fossil Pokemon", Category: CategorySpecial, Points: 20}, collectionComplete("special_fossils")},
	{AchievementDefinition{ID: "special_starters", Name: "Starter Supreme", Event: "Caught All Starter Pokemon", Category: CategorySpecial, Points: 25}, collectionComplete("special_starters")},
	{AchievementDefinition{ID: "special_mythicals", Name: "Myth Hunter", Event: "Caught All Mythical Pokemon", Category: CategorySpecial, Points: 35}, collectionComplete("special_mythicals")},
	{AchievementDefinition{ID: "special_legendaries", Name: "Legendary Conqueror", Event: "Caught All Legendary Pokemon", Category: CategorySpecial, Points: 50}, collectionComplete("special_legendaries")},
	{AchievementDefinition{ID: CompletionistID, Name: "Poketerm - Completed it mate!", Event: "Completed all other achievements", Category: CategorySpecial, Points: 50}, nil},
}

// Definitions returns the ordered achievement catalog.
func Definitions() []AchievementDefinition {
	defs := make([]AchievementDefinition, len(achievementRules))
	for i, rule := range achievementRules {
		defs[i] = rule.AchievementDefinition
	}
	return defs
}

// DefinitionByID returns a catalog definition by stable ID.
func DefinitionByID(id string) (AchievementDefinition, bool) {
	for _, rule := range achievementRules {
		if rule.ID == id {
			return rule.AchievementDefinition, true
		}
	}
	return AchievementDefinition{}, false
}

func nonCompletionistRules() []achievementRule {
	rules := make([]achievementRule, 0, len(achievementRules)-1)
	for _, rule := range achievementRules {
		if rule.ID != CompletionistID {
			rules = append(rules, rule)
		}
	}
	return rules
}

func completionistDefinition() AchievementDefinition {
	def, _ := DefinitionByID(CompletionistID)
	return def
}

func allGenerationsComplete(c Context) bool {
	for gen := 1; gen <= 8; gen++ {
		if !c.GenComplete[gen] {
			return false
		}
	}
	return true
}

func collectionComplete(id string) func(Context) bool {
	return func(c Context) bool {
		for _, col := range specialCollections {
			if col.ID != id {
				continue
			}
			for name := range col.Required {
				if !c.CaughtSet[name] {
					return false
				}
			}
			return len(col.Required) > 0
		}
		return false
	}
}
