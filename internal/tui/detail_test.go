package tui

import (
	"regexp"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestDetailViewShowsPokemonDataAndCaughtStatus(t *testing.T) {
	view := stripANSI(newDetailModel("pikachu", true, false).View())
	for _, want := range []string{
		"#025 Pikachu",
		"Generation:",
		"Caught:",
		"normal caught",
		"shiny  not caught",
		"enter/esc/q back",
	} {
		if !strings.Contains(view, want) {
			t.Fatalf("detail view missing %q:\n%s", want, view)
		}
	}

	shinyOnly := stripANSI(newDetailModel("pikachu", false, true).View())
	for _, want := range []string{"normal not caught", "shiny  caught"} {
		if !strings.Contains(shinyOnly, want) {
			t.Fatalf("shiny detail view missing %q:\n%s", want, shinyOnly)
		}
	}
}

func TestDetailViewShowsNotFoundForUnknownPokemon(t *testing.T) {
	view := newDetailModel("missingno", false, false).View()
	if !strings.Contains(view, "Pokemon not found") {
		t.Fatalf("missing Pokemon view = %q", view)
	}
	if !strings.Contains(view, "press enter to go back") {
		t.Fatalf("missing Pokemon view lacks back hint = %q", view)
	}
}

func TestDetailModelEnterMarksDone(t *testing.T) {
	m := newDetailModel("pikachu", true, false)
	_, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if !m.done {
		t.Fatalf("detail model done = false after enter")
	}
}

var ansiRE = regexp.MustCompile(`\x1b\[[0-9;]*m`)

func stripANSI(s string) string {
	return ansiRE.ReplaceAllString(s, "")
}
