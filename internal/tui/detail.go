package tui

import (
	"fmt"
	"strings"

	"github.com/brunoorsolon/pokedex-fetch-go/internal/pokemon"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type detailModel struct {
	name      string
	hasNormal bool
	hasShiny  bool
	done      bool
}

func newDetailModel(name string, hasNormal, hasShiny bool) *detailModel {
	return &detailModel{name: name, hasNormal: hasNormal, hasShiny: hasShiny}
}

func (m *detailModel) Init() tea.Cmd { return nil }

func (m *detailModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "backspace", "q", "enter":
			m.done = true
		}
	}
	return m, nil
}

func (m *detailModel) View() string {
	p, ok := pokemon.GetByName(m.name)
	if !ok {
		return "Pokemon not found\n\npress enter to go back"
	}

	borderStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
	headerStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("11")).Bold(true)
	labelStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("14")).Bold(true)
	caughtStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	missStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("9"))

	variant := pokemon.SpriteRegular
	if m.hasShiny && !m.hasNormal {
		variant = pokemon.SpriteShiny
	}
	sprite, err := pokemon.GetSprite(p.Name, pokemon.SpriteSmall, variant)
	if err != nil {
		sprite = "(sprite unavailable)"
	}
	spriteLines, spriteW := splitSpriteLines(sprite)

	normal := missStyle.Render("not caught")
	if m.hasNormal {
		normal = caughtStyle.Render("caught")
	}
	shiny := missStyle.Render("not caught")
	if m.hasShiny {
		shiny = caughtStyle.Render("caught")
	}
	types := strings.Join(p.Types, ", ")
	if types == "" {
		types = "unknown"
	}
	detailLines := []string{
		fmt.Sprintf("%s %d", labelStyle.Render("Generation:"), p.Generation),
		fmt.Sprintf("%s %s", labelStyle.Render("Types:"), types),
		fmt.Sprintf("%s %d  %s %d", labelStyle.Render("Height:"), p.Height, labelStyle.Render("Weight:"), p.Weight),
		fmt.Sprintf("%s HP %d / Atk %d / Def %d", labelStyle.Render("Stats:"), p.Stats.HP, p.Stats.Attack, p.Stats.Defense),
		fmt.Sprintf("       SpA %d / SpD %d / Spe %d", p.Stats.SpAttack, p.Stats.SpDefense, p.Stats.Speed),
		fmt.Sprintf("%s normal %s", labelStyle.Render("Caught:"), normal),
		fmt.Sprintf("        shiny  %s", shiny),
	}

	panelW := 42
	for _, detail := range detailLines {
		if w := lipgloss.Width(detail); w > panelW {
			panelW = w
		}
	}
	contentW := spriteW + 4 + panelW
	boxW := contentW + 3
	if boxW < 64 {
		boxW = 64
	}
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

	var sb strings.Builder
	sb.WriteString(top + "\n")
	title := fmt.Sprintf("#%03d %s", p.Number, displayPokemonName(p.Name))
	sb.WriteString(line(headerStyle.Render(title)) + "\n")
	sb.WriteString(mid + "\n")
	rows := len(spriteLines)
	if len(detailLines) > rows {
		rows = len(detailLines)
	}
	for i := 0; i < rows; i++ {
		left := ""
		if i < len(spriteLines) {
			left = spriteLines[i]
		}
		right := ""
		if i < len(detailLines) {
			right = detailLines[i]
		}
		row := left + strings.Repeat(" ", spriteW-lipgloss.Width(left)+4) + right
		sb.WriteString(line(row) + "\n")
	}
	sb.WriteString(mid + "\n")
	sb.WriteString(line("enter/esc/q back"))
	sb.WriteString("\n" + bot + "\n")
	return sb.String()
}

func splitSpriteLines(sprite string) ([]string, int) {
	sprite = strings.TrimRight(sprite, "\n")
	if sprite == "" {
		return []string{}, 0
	}
	lines := strings.Split(sprite, "\n")
	maxW := 0
	for _, line := range lines {
		if w := lipgloss.Width(line); w > maxW {
			maxW = w
		}
	}
	return lines, maxW
}
