package tui

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/brunoorsolon/pokedex-fetch-go/internal/pokemon"
	"github.com/brunoorsolon/pokedex-fetch-go/internal/state"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const pageSize = 15

// pokedexModel is the bubbletea model for the interactive Pokedex list view.
type pokedexModel struct {
	gen       int
	pkgList   []string // pokemon names for current generation
	pdx       *state.PokedexState
	cursor    int
	page      int
	width     int
	height    int
	gotoMode  bool
	gotoInput string
	gotoError string
	detail    *detailModel // non-nil when in detail view
}

func newPokedexModel(gen int) (*pokedexModel, error) {
	pdx, err := state.LoadPokedex()
	if err != nil {
		pdx, _ = state.NewPokedex()
	}
	g, ok := pokemon.GetGeneration(gen)
	if !ok {
		return nil, fmt.Errorf("generation %d not found", gen)
	}
	return &pokedexModel{
		gen:     gen,
		pkgList: g.Names,
		pdx:     pdx,
		cursor:  0,
		page:    1,
	}, nil
}

func (m *pokedexModel) Init() tea.Cmd { return nil }

func (m *pokedexModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Delegate to detail view if active
	if m.detail != nil {
		newModel, cmd := m.detail.Update(msg)
		d := newModel.(*detailModel)
		if d.done {
			m.detail = nil
			return m, nil
		}
		m.detail = d
		return m, cmd
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.gotoMode {
			switch msg.String() {
			case "ctrl+c":
				return m, tea.Quit
			case "esc":
				m.clearGoto()
			case "enter":
				m.commitGoto()
			case "backspace", "ctrl+h":
				if len(m.gotoInput) > 0 {
					m.gotoInput = m.gotoInput[:len(m.gotoInput)-1]
					m.gotoError = ""
				}
			default:
				s := msg.String()
				if len(s) == 1 && s[0] >= '0' && s[0] <= '9' && len(m.gotoInput) < 4 {
					m.gotoInput += s
					m.gotoError = ""
				}
			}
			return m, nil
		}

		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
				if m.cursor < (m.page-1)*pageSize {
					m.page--
				}
			}
		case "down", "j":
			if m.cursor < len(m.pkgList)-1 {
				m.cursor++
				if m.cursor >= m.page*pageSize {
					m.page++
				}
			}
		case "left":
			if m.page > 1 {
				m.page--
				m.cursor = (m.page - 1) * pageSize
			}
		case "right":
			maxPage := (len(m.pkgList)-1)/pageSize + 1
			if m.page < maxPage {
				m.page++
				m.cursor = (m.page - 1) * pageSize
			}
		case "g":
			m.gotoMode = true
			m.gotoInput = ""
			m.gotoError = ""
		case "enter":
			if m.cursor < len(m.pkgList) {
				name := m.pkgList[m.cursor]
				if m.pdx.IsCaught(name) {
					r := m.pdx.Pokemon[name]
					hasNormal := r != nil && r.Normal != nil && r.Normal.Count > 0
					hasShiny := r != nil && r.Shiny != nil && r.Shiny.Count > 0
					m.detail = newDetailModel(name, hasNormal, hasShiny)
				}
			}
		case "1", "2", "3", "4", "5", "6", "7", "8":
			g := int(msg.String()[0] - '0')
			m.setGeneration(g)
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}
	return m, nil
}

func (m *pokedexModel) clearGoto() {
	m.gotoMode = false
	m.gotoInput = ""
	m.gotoError = ""
}

func (m *pokedexModel) setGeneration(g int) bool {
	gen, ok := pokemon.GetGeneration(g)
	if !ok {
		return false
	}
	m.gen = g
	m.pkgList = gen.Names
	m.cursor = 0
	m.page = 1
	return true
}

func (m *pokedexModel) commitGoto() {
	if m.gotoInput == "" {
		m.gotoError = "Enter a National Dex number"
		return
	}
	num, err := strconv.Atoi(m.gotoInput)
	if err != nil {
		m.gotoError = "Enter a valid National Dex number"
		return
	}
	p, ok := pokemon.GetByNumber(num)
	if !ok {
		m.gotoError = fmt.Sprintf("#%03d not found", num)
		return
	}
	if p.Generation != m.gen && !m.setGeneration(p.Generation) {
		m.gotoError = fmt.Sprintf("Generation %d not found", p.Generation)
		return
	}
	for i, name := range m.pkgList {
		if name == p.Name {
			m.cursor = i
			m.page = (i / pageSize) + 1
			m.clearGoto()
			return
		}
	}
	m.gotoError = fmt.Sprintf("#%03d is not listed in generation %d", num, m.gen)
}

var pokemonDisplayNames = map[string]string{
	"nidoran-f": "Nidoran♀",
	"nidoran-m": "Nidoran♂",
	"mr-mime":   "Mr. Mime",
	"mime-jr":   "Mime Jr.",
	"ho-oh":     "Ho-Oh",
	"porygon-z": "Porygon-Z",
	"type-null": "Type: Null",
	"jangmo-o":  "Jangmo-o",
	"hakamo-o":  "Hakamo-o",
	"kommo-o":   "Kommo-o",
	"tapu-koko": "Tapu Koko",
	"tapu-lele": "Tapu Lele",
	"tapu-bulu": "Tapu Bulu",
	"tapu-fini": "Tapu Fini",
	"mr-rime":   "Mr. Rime",
	"farfetchd": "Farfetch'd",
	"sirfetchd": "Sirfetch'd",
	"flabebe":   "Flabébé",
}

func displayPokemonName(name string) string {
	if display, ok := pokemonDisplayNames[name]; ok {
		return display
	}
	if name == "" {
		return name
	}
	return strings.ToUpper(name[:1]) + name[1:]
}

func (m *pokedexModel) View() string {
	if m.detail != nil {
		return m.detail.View()
	}

	borderStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
	headerStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("11")).Bold(true)
	caughtStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	missStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
	shinyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("11"))
	cursorStyle := lipgloss.NewStyle().Reverse(true)

	boxW := 66
	inner := boxW - 2
	hr := borderStyle.Render("┌" + strings.Repeat("─", inner) + "┐")
	hrMid := borderStyle.Render("├" + strings.Repeat("─", inner) + "┤")
	hrBot := borderStyle.Render("└" + strings.Repeat("─", inner) + "┘")

	bord := func(content string) string {
		vis := lipgloss.Width(content)
		pad := inner - vis - 1
		if pad < 0 {
			pad = 0
		}
		return borderStyle.Render("│") + " " + content + strings.Repeat(" ", pad) + borderStyle.Render("│")
	}

	var sb strings.Builder
	sb.WriteString(hr + "\n")

	title := fmt.Sprintf("Generation %d Pokedex", m.gen)
	pad := (inner - len(title)) / 2
	if pad < 0 {
		pad = 0
	}
	sb.WriteString(borderStyle.Render("│") + strings.Repeat(" ", pad) + headerStyle.Render(title) +
		strings.Repeat(" ", inner-len(title)-pad) + borderStyle.Render("│") + "\n")
	sb.WriteString(hrMid + "\n")

	total := len(m.pkgList)
	unique := 0
	for _, name := range m.pkgList {
		if m.pdx.IsCaught(name) {
			unique++
		}
	}
	caughtStr := fmt.Sprintf("Pokemon Caught: %d/%d", unique, total)
	if unique == total && total > 0 {
		sb.WriteString(bord(caughtStyle.Render(caughtStr)) + "\n")
	} else {
		sb.WriteString(bord(caughtStr) + "\n")
	}

	// Keep the count columns lined up while leaving room for long names.
	maxNameLen := 14
	for _, name := range m.pkgList {
		if w := lipgloss.Width(displayPokemonName(name)); w > maxNameLen {
			maxNameLen = w
		}
	}
	sb.WriteString(bord(fmt.Sprintf("No.   %-*s   Normal  Shiny", maxNameLen, "Pokemon")) + "\n")
	sb.WriteString(hrMid + "\n")

	start := (m.page - 1) * pageSize
	end := start + pageSize
	if end > total {
		end = total
	}
	printed := 0

	for i, name := range m.pkgList[start:end] {
		idx := start + i
		r := m.pdx.Pokemon[name]
		var n, s int
		if r != nil {
			if r.Normal != nil {
				n = r.Normal.Count
			}
			if r.Shiny != nil {
				s = r.Shiny.Count
			}
		}
		var ncell, scell string
		if n > 0 {
			ncell = fmt.Sprintf("%s %-3d", caughtStyle.Render("✔"), n)
		} else {
			ncell = fmt.Sprintf("%s %-3d", missStyle.Render("✖"), n)
		}
		if s > 0 {
			scell = fmt.Sprintf("%s %-3d", shinyStyle.Render("✔"), s)
		} else {
			scell = fmt.Sprintf("%s %-3d", missStyle.Render("✖"), s)
		}
		number := 0
		if p, ok := pokemon.GetByName(name); ok {
			number = p.Number
		}
		displayName := displayPokemonName(name)
		row := fmt.Sprintf("#%03d  %-*s   %s  %s", number, maxNameLen, displayName, ncell, scell)
		if idx == m.cursor {
			row = cursorStyle.Render(row)
		}
		sb.WriteString(bord(row) + "\n")
		printed++
	}
	for i := printed; i < pageSize; i++ {
		sb.WriteString(bord("") + "\n")
	}

	maxPage := (total-1)/pageSize + 1
	if total == 0 {
		maxPage = 1
	}
	footer := fmt.Sprintf("Page %d/%d  ↑↓ nav  ←→ page  g goto  1-8 gen  q quit", m.page, maxPage)
	if m.gotoMode {
		if m.gotoError != "" {
			footer = m.gotoError
		} else {
			footer = fmt.Sprintf("Go to National Dex #: %s  enter jump  esc cancel", m.gotoInput)
		}
	}
	sb.WriteString(hrMid + "\n")
	sb.WriteString(bord(footer) + "\n")
	sb.WriteString(hrBot + "\n")
	return sb.String()
}

// RunPokedex launches the interactive Pokedex TUI starting on the given generation.
func RunPokedex(gen int) error {
	m, err := newPokedexModel(gen)
	if err != nil {
		return err
	}
	p := tea.NewProgram(m, tea.WithAltScreen())
	_, err = p.Run()
	return err
}
