package state

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// TrainerState holds the trainer's persistent profile and progression data.
type TrainerState struct {
	Name              string         `json:"name"`
	Region            string         `json:"region"`
	Hometown          string         `json:"hometown"`
	Created           string         `json:"created"`
	Favourites        []int          `json:"favourites"`
	Level             int            `json:"level"`
	XP                int            `json:"xp"`
	TotalXP           int            `json:"total_xp"`
	Streak            StreakData     `json:"streak"`
	DailyCatch        DailyCatchData `json:"daily_catch"`
	HallOfFame        []Achievement  `json:"hall_of_fame"`
	AchievementPoints int            `json:"achievement_points"`
	Version           int            `json:"version"`
}

// StreakData tracks daily catch streak information.
type StreakData struct {
	Current        int    `json:"current"`
	Longest        int    `json:"longest"`
	LastActiveDate string `json:"last_active_date"`
}

// DailyCatchData tracks the daily catch count.
type DailyCatchData struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}

// Achievement represents an unlocked achievement in the Hall of Fame.
type Achievement struct {
	ID       string `json:"id,omitempty"`
	Name     string `json:"name"`
	Event    string `json:"event"`
	Category string `json:"category"`
	Points   int    `json:"points"`
	Date     string `json:"date"`
}

type rankTier struct {
	min, max int
	title    string
	color    string
}

var rankTiers = []rankTier{
	{0, 5, "Beginner Trainer", "gray"},
	{6, 10, "Rookie Trainer", "green"},
	{11, 15, "Ace Trainer", "blue"},
	{16, 20, "Gym Challenger", "brown"},
	{21, 25, "Gym Leader", "black"},
	{26, 30, "Elite Trainer", "pink"},
	{31, 35, "Elite Four", "red"},
	{36, 40, "Champion", "yellow"},
	{41, 45, "Regional Master", "blue"},
	{46, 49, "Pokedex Master", "white"},
	{50, 9999, "Grand Champion", "purple"},
}

var colorMap = map[string]int{
	"black":  242,
	"blue":   117,
	"brown":  180,
	"gray":   252,
	"green":  114,
	"pink":   218,
	"purple": 147,
	"red":    210,
	"white":  255,
	"yellow": 229,
}

// HomeTowns maps region name to a list of town options.
var HomeTowns = map[string][]string{
	"Kanto":  {"Pallet Town", "Viridian City", "Pewter City", "Cerulean City", "Vermilion City"},
	"Johto":  {"New Bark Town", "Cherrygrove City", "Violet City", "Azalea Town", "Goldenrod City"},
	"Hoenn":  {"Littleroot Town", "Oldale Town", "Petalburg City", "Rustboro City", "Mauville City"},
	"Sinnoh": {"Twinleaf Town", "Sandgem Town", "Jubilife City", "Oreburgh City", "Floaroma Town"},
	"Unova":  {"Nuvema Town", "Accumula Town", "Striaton City", "Nacrene City", "Castelia City"},
	"Kalos":  {"Vaniville Town", "Aquacorde Town", "Santalune City", "Lumiose City", "Camphrier Town"},
	"Alola":  {"Iki Town", "Hau'oli City", "Heahea City", "Paniola Town", "Konikoni City"},
	"Galar":  {"Postwick", "Wedgehurst", "Motostoke", "Hammerlocke", "Wyndon"},
}

// GenRegions maps generation number to region name.
var GenRegions = map[int]string{
	1: "Kanto", 2: "Johto", 3: "Hoenn", 4: "Sinnoh",
	5: "Unova", 6: "Kalos", 7: "Alola", 8: "Galar",
}

// DefaultTrainer returns a fresh, unconfigured trainer.
func DefaultTrainer() *TrainerState {
	today := time.Now().Format("2006-01-02")
	return &TrainerState{
		Name:              "None",
		Region:            "None",
		Hometown:          "None",
		Created:           "None",
		Favourites:        []int{},
		Level:             1,
		XP:                0,
		TotalXP:           0,
		Streak:            StreakData{Current: 0, Longest: 0, LastActiveDate: "1996-02-27"},
		DailyCatch:        DailyCatchData{Date: today, Count: 0},
		HallOfFame:        []Achievement{},
		AchievementPoints: 0,
		Version:           1,
	}
}

// LoadTrainer reads trainer.json from disk, or returns a default if absent.
func LoadTrainer() (*TrainerState, error) {
	data, err := os.ReadFile(TrainerPath())
	if os.IsNotExist(err) {
		return DefaultTrainer(), nil
	}
	if err != nil {
		return nil, err
	}
	var t TrainerState
	if err := json.Unmarshal(data, &t); err != nil {
		return nil, err
	}
	if t.HallOfFame == nil {
		t.HallOfFame = []Achievement{}
	}
	return &t, nil
}

// Save writes the trainer state to disk atomically.
func (t *TrainerState) Save() error {
	data, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		return err
	}
	return AtomicWrite(TrainerPath(), data)
}

// AddXP awards XP and triggers level-ups as needed.
// Level cost formula: (level + 1) * 500, capped at 25500 for level >= 50.
func (t *TrainerState) AddXP(amount int) {
	t.XP += amount
	t.TotalXP += amount
	for {
		cost := (t.Level + 1) * 500
		if t.Level >= 50 {
			cost = 25500
		}
		if t.XP >= cost {
			t.XP -= cost
			t.Level++
		} else {
			break
		}
	}
}

// UpdateStreak updates the daily catch streak based on today's date.
func (t *TrainerState) UpdateStreak() {
	today := time.Now().Format("2006-01-02")
	todayTime, _ := time.Parse("2006-01-02", today)
	lastTime, err := time.Parse("2006-01-02", t.Streak.LastActiveDate)
	if err != nil {
		lastTime = time.Time{}
	}

	diff := todayTime.Sub(lastTime).Hours() / 24
	switch {
	case diff == 0:
		// Same day — no streak change
	case diff == 1:
		t.Streak.Current++
	default:
		t.Streak.Current = 1
	}

	if t.Streak.Current > t.Streak.Longest {
		t.Streak.Longest = t.Streak.Current
	}
	t.Streak.LastActiveDate = today
}

// UpdateDailyCatch increments the daily catch counter, resetting on a new day.
func (t *TrainerState) UpdateDailyCatch() {
	today := time.Now().Format("2006-01-02")
	if t.DailyCatch.Date != today {
		t.DailyCatch = DailyCatchData{Date: today, Count: 1}
	} else {
		t.DailyCatch.Count++
	}
}

// Unlock adds an achievement to the Hall of Fame if not already present.
// Returns true if the achievement was newly unlocked.
func (t *TrainerState) Unlock(id, name, event, category string, points int) bool {
	for _, a := range t.HallOfFame {
		if a.ID == id {
			return false
		}
	}
	t.HallOfFame = append(t.HallOfFame, Achievement{
		ID:       id,
		Name:     name,
		Event:    event,
		Category: category,
		Points:   points,
		Date:     time.Now().Format("2006-01-02"),
	})
	t.AchievementPoints += points
	return true
}

// RecalculateAchievementPoints resets the cached achievement score from the
// Hall of Fame entries.
func (t *TrainerState) RecalculateAchievementPoints() {
	total := 0
	for _, a := range t.HallOfFame {
		total += a.Points
	}
	t.AchievementPoints = total
}

// GetRank returns an ANSI-colored rank string for the trainer's current level.
func (t *TrainerState) GetRank() string {
	for _, tier := range rankTiers {
		if t.Level >= tier.min && t.Level <= tier.max {
			colorID := colorMap[tier.color]
			return fmt.Sprintf("\033[48;5;%dm\033[38;5;16m ✦ %s ✦ \033[0m", colorID, tier.title)
		}
	}
	return "Unknown"
}

// GetRankTitle returns just the rank title string (no ANSI).
func (t *TrainerState) GetRankTitle() string {
	for _, tier := range rankTiers {
		if t.Level >= tier.min && t.Level <= tier.max {
			return tier.title
		}
	}
	return "Unknown"
}

// NextLevelXPCost returns the XP cost to reach the next level.
func (t *TrainerState) NextLevelXPCost() int {
	if t.Level >= 50 {
		return 25500
	}
	return (t.Level + 1) * 500
}

// IsConfigured returns true if the trainer has been set up (name is not "None").
func (t *TrainerState) IsConfigured() bool {
	return t.Name != "None" && t.Name != ""
}
