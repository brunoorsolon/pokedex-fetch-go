package tui

import (
	"fmt"
	"strings"

	"github.com/brunoorsolon/pokedex-fetch-go/internal/achievement"
	"github.com/brunoorsolon/pokedex-fetch-go/internal/pokemon"
	"github.com/brunoorsolon/pokedex-fetch-go/internal/state"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type trainerModel struct {
	trainer *state.TrainerState
	pokedex *state.PokedexState
}

func (m trainerModel) Init() tea.Cmd { return nil }

func (m trainerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok {
		switch key.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m trainerModel) View() string {
	borderStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("12"))
	headerStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("11")).Bold(true)
	labelStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("14")).Bold(true)

	boxW := 64
	inner := boxW - 2
	top := borderStyle.Render("┌" + strings.Repeat("─", inner) + "┐")
	mid := borderStyle.Render("├" + strings.Repeat("─", inner) + "┤")
	bot := borderStyle.Render("└" + strings.Repeat("─", inner) + "┘")
	line := func(content string) string {
		vis := lipgloss.Width(content)
		pad := inner - vis - 1
		if pad < 0 {
			pad = 0
		}
		return borderStyle.Render("│") + " " + content + strings.Repeat(" ", pad) + borderStyle.Render("│")
	}

	unique := 0
	total := 0
	if m.pokedex != nil {
		unique = m.pokedex.UniqueCaught()
		total = m.pokedex.TotalCatches()
	}
	allPokemon := len(pokemon.GetAll())
	if allPokemon == 0 {
		allPokemon = 905
	}
	m.trainer.RecalculateAchievementPoints()
	next := m.trainer.NextLevelXPCost()

	var sb strings.Builder
	sb.WriteString(top + "\n")
	sb.WriteString(line(headerStyle.Render("Trainer Profile")) + "\n")
	sb.WriteString(mid + "\n")
	sb.WriteString(line(fmt.Sprintf("%s %s", labelStyle.Render("Name:"), m.trainer.Name)) + "\n")
	sb.WriteString(line(fmt.Sprintf("%s %s / %s", labelStyle.Render("Home:"), m.trainer.Region, m.trainer.Hometown)) + "\n")
	sb.WriteString(line(fmt.Sprintf("%s %s", labelStyle.Render("Rank:"), m.trainer.GetRankTitle())) + "\n")
	sb.WriteString(line(fmt.Sprintf("%s %d", labelStyle.Render("Level:"), m.trainer.Level)) + "\n")
	sb.WriteString(line(fmt.Sprintf("%s %d/%d (%d total)", labelStyle.Render("XP:"), m.trainer.XP, next, m.trainer.TotalXP)) + "\n")
	sb.WriteString(line(fmt.Sprintf("%s %d/%d unique, %d total", labelStyle.Render("Pokedex:"), unique, allPokemon, total)) + "\n")
	sb.WriteString(line(fmt.Sprintf("%s current %d, longest %d", labelStyle.Render("Streak:"), m.trainer.Streak.Current, m.trainer.Streak.Longest)) + "\n")
	sb.WriteString(line(fmt.Sprintf("%s %d on %s", labelStyle.Render("Daily catches:"), m.trainer.DailyCatch.Count, m.trainer.DailyCatch.Date)) + "\n")
	sb.WriteString(line(fmt.Sprintf("%s %d/%d", labelStyle.Render("Achievements:"), len(m.trainer.HallOfFame), len(achievement.Definitions()))) + "\n")
	sb.WriteString(line(fmt.Sprintf("%s %d/%d", labelStyle.Render("Achievement Points:"), m.trainer.AchievementPoints, achievement.TotalPoints)) + "\n")
	sb.WriteString(mid + "\n")
	sb.WriteString(line("q quit"))
	sb.WriteString("\n" + bot + "\n")
	return sb.String()
}

// RunTrainer launches the trainer profile TUI.
func RunTrainer() error {
	tr, err := state.LoadTrainer()
	if err != nil {
		tr = state.DefaultTrainer()
	}
	pdx, err := state.LoadPokedex()
	if err != nil {
		pdx, _ = state.NewPokedex()
	}
	_, err = tea.NewProgram(trainerModel{trainer: tr, pokedex: pdx}, tea.WithAltScreen()).Run()
	return err
}
