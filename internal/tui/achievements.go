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

type achievementsModel struct {
	trainer *state.TrainerState
	cursor  int
	page    int
}

const achievementPageSize = 10

func (m achievementsModel) Init() tea.Cmd { return nil }

func (m achievementsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	defs := achievement.Definitions()
	if key, ok := msg.(tea.KeyMsg); ok {
		switch key.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
				if m.cursor < (m.page-1)*achievementPageSize {
					m.page--
				}
			}
		case "down", "j":
			if m.cursor < len(defs)-1 {
				m.cursor++
				if m.cursor >= m.page*achievementPageSize {
					m.page++
				}
			}
		case "left":
			if m.page > 1 {
				m.page--
				m.cursor = (m.page - 1) * achievementPageSize
			}
		case "right":
			maxPage := maxAchievementPage(len(defs))
			if m.page < maxPage {
				m.page++
				m.cursor = (m.page - 1) * achievementPageSize
			}
		}
	}
	return m, nil
}

func (m achievementsModel) View() string {
	borderStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("13"))
	headerStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("11")).Bold(true)
	mutedStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	unlockedStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	cursorStyle := lipgloss.NewStyle().Reverse(true)

	boxW := 96
	inner := boxW - 2
	top := borderStyle.Render("┌" + strings.Repeat("─", inner) + "┐")
	mid := borderStyle.Render("├" + strings.Repeat("─", inner) + "┤")
	bot := borderStyle.Render("└" + strings.Repeat("─", inner) + "┘")
	line := func(content string) string {
		if lipgloss.Width(content) > inner-1 {
			content = truncateDisplay(content, inner-2)
		}
		vis := lipgloss.Width(content)
		pad := inner - vis - 1
		if pad < 0 {
			pad = 0
		}
		return borderStyle.Render("│") + " " + content + strings.Repeat(" ", pad) + borderStyle.Render("│")
	}

	m.trainer.RecalculateAchievementPoints()
	defs := achievement.Definitions()
	unlocked := achievementsByID(m.trainer.HallOfFame)
	maxPage := maxAchievementPage(len(defs))
	start := (m.page - 1) * achievementPageSize
	end := start + achievementPageSize
	if end > len(defs) {
		end = len(defs)
	}

	var sb strings.Builder
	sb.WriteString(top + "\n")
	sb.WriteString(line(headerStyle.Render("Hall of Fame")) + "\n")
	sb.WriteString(line(fmt.Sprintf("Achievements: %d/%d   Points: %d/%d", len(m.trainer.HallOfFame), len(defs), m.trainer.AchievementPoints, achievement.TotalPoints)) + "\n")
	sb.WriteString(mid + "\n")
	for i, def := range defs[start:end] {
		row := achievementRow(def, unlocked, mutedStyle, unlockedStyle)
		if start+i == m.cursor {
			row = cursorStyle.Render(row)
		}
		sb.WriteString(line(row) + "\n")
	}
	for i := end - start; i < achievementPageSize; i++ {
		sb.WriteString(line("") + "\n")
	}
	sb.WriteString(mid + "\n")
	sb.WriteString(line(fmt.Sprintf("Page %d/%d  ↑↓ nav  ←→ page  q quit", m.page, maxPage)))
	sb.WriteString("\n" + bot + "\n")
	return sb.String()
}

func achievementRow(def achievement.AchievementDefinition, unlocked map[string]state.Achievement, mutedStyle, unlockedStyle lipgloss.Style) string {
	if a, ok := unlocked[def.ID]; ok {
		return unlockedStyle.Render(fmt.Sprintf("✓ %03d pts [%s] %s — %s (%s)", def.Points, def.Category, def.Name, def.Event, a.Date))
	}
	return mutedStyle.Render(fmt.Sprintf("□ %03d pts [%s] %s — %s", def.Points, def.Category, def.Name, def.Event))
}

func achievementsByID(items []state.Achievement) map[string]state.Achievement {
	m := make(map[string]state.Achievement, len(items))
	for _, a := range items {
		m[a.ID] = a
	}
	return m
}

func maxAchievementPage(count int) int {
	if count == 0 {
		return 1
	}
	return (count-1)/achievementPageSize + 1
}

func truncateDisplay(s string, max int) string {
	if max <= 1 {
		return "…"
	}
	if lipgloss.Width(s) <= max {
		return s
	}
	runes := []rune(s)
	for len(runes) > 0 && lipgloss.Width(string(runes))+1 > max {
		runes = runes[:len(runes)-1]
	}
	return string(runes) + "…"
}

// RunAchievements launches the Hall of Fame TUI.
func RunAchievements() error {
	tr, err := state.LoadTrainer()
	if err != nil {
		tr = state.DefaultTrainer()
	}
	pdx, err := state.LoadPokedex()
	if err != nil {
		pdx, _ = state.NewPokedex()
	}
	beforeCount := len(tr.HallOfFame)
	beforePoints := tr.AchievementPoints
	achievement.CheckAll(pdx, tr, pokemon.GetAllGenerations())
	if len(tr.HallOfFame) != beforeCount || tr.AchievementPoints != beforePoints {
		_ = tr.Save()
	}
	_, err = tea.NewProgram(achievementsModel{trainer: tr, page: 1}, tea.WithAltScreen()).Run()
	return err
}
