package cli

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/brunoorsolon/pokedex-fetch-go/internal/config"
	"github.com/brunoorsolon/pokedex-fetch-go/internal/state"
	"github.com/spf13/cobra"
)

func TestRunCatchWritesConfigPokedexAndTrainerState(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", tmp)

	cmd := newCatchFlowTestCommand("--gen", "1", "--shiny", "--raw")
	out, err := captureStdout(t, cmd.Execute)
	if err != nil {
		t.Fatalf("runCatch() error = %v", err)
	}
	if out == "" {
		t.Fatalf("runCatch() produced no raw sprite output")
	}

	dataDir := filepath.Join(tmp, "pokedex-fetch-go")
	for _, path := range []string{
		filepath.Join(dataDir, "config.toml"),
		filepath.Join(dataDir, "pokedex.json"),
		filepath.Join(dataDir, "trainer.json"),
	} {
		if _, err := os.Stat(path); err != nil {
			t.Fatalf("expected state file %s: %v", path, err)
		}
	}

	pdx, err := state.LoadPokedex()
	if err != nil {
		t.Fatalf("LoadPokedex() error = %v", err)
	}
	if got := pdx.UniqueCaught(); got != 1 {
		t.Fatalf("UniqueCaught() = %d, want 1", got)
	}
	if got := pdx.TotalCatches(); got != 1 {
		t.Fatalf("TotalCatches() = %d, want 1", got)
	}
	if got := pdx.TotalShinies(); got != 1 {
		t.Fatalf("TotalShinies() = %d, want 1", got)
	}

	tr, err := state.LoadTrainer()
	if err != nil {
		t.Fatalf("LoadTrainer() error = %v", err)
	}
	if tr.TotalXP != 150 {
		t.Fatalf("TotalXP = %d, want 150 for forced shiny catch", tr.TotalXP)
	}
	if tr.DailyCatch.Count != 1 {
		t.Fatalf("DailyCatch.Count = %d, want 1", tr.DailyCatch.Count)
	}
	if tr.AchievementPoints != 20 {
		t.Fatalf("AchievementPoints = %d, want 20", tr.AchievementPoints)
	}
	if !hasTrainerAchievement(tr, "pokedex_first_catch") || !hasTrainerAchievement(tr, "shiny_first") {
		t.Fatalf("expected first catch and shiny achievements, got %#v", tr.HallOfFame)
	}
}

func TestRunCatchRejectsInvalidGenerationBeforeCatchStateWrite(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", tmp)

	cmd := newCatchFlowTestCommand("--gen", "9", "--raw")
	_, err := captureStdout(t, cmd.Execute)
	if err == nil {
		t.Fatalf("runCatch() error = nil, want invalid generation error")
	}
	if !strings.Contains(err.Error(), "invalid generation") {
		t.Fatalf("runCatch() error = %q, want invalid generation", err.Error())
	}

	if _, err := os.Stat(config.ConfigPath()); err != nil {
		t.Fatalf("config should be bootstrapped before gen validation: %v", err)
	}
	if _, err := os.Stat(state.PokedexPath()); !os.IsNotExist(err) {
		t.Fatalf("pokedex state exists after invalid generation, err=%v", err)
	}
	if _, err := os.Stat(state.TrainerPath()); !os.IsNotExist(err) {
		t.Fatalf("trainer state exists after invalid generation, err=%v", err)
	}
}

func newCatchFlowTestCommand(args ...string) *cobra.Command {
	cmd := &cobra.Command{
		Use:  "catch-test",
		RunE: runCatch,
	}
	cmd.Flags().String("gen", "", `Generation(s) to pick from (e.g. "1", "1-3", "1,3,6")`)
	cmd.Flags().Bool("shiny", false, "Force shiny variant")
	cmd.Flags().Bool("big", false, "Use large sprite")
	cmd.Flags().Bool("no-title", false, "Don't display Pokemon name")
	cmd.Flags().Bool("raw", false, "Output sprite only (for fastfetch piping)")
	cmd.SetArgs(args)
	return cmd
}

func captureStdout(t *testing.T, fn func() error) (string, error) {
	t.Helper()

	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe() error = %v", err)
	}
	os.Stdout = w
	defer func() { os.Stdout = old }()

	out := make(chan string, 1)
	go func() {
		var buf bytes.Buffer
		_, _ = io.Copy(&buf, r)
		out <- buf.String()
	}()

	runErr := fn()
	_ = w.Close()
	output := <-out
	_ = r.Close()
	return output, runErr
}

func hasTrainerAchievement(tr *state.TrainerState, id string) bool {
	for _, a := range tr.HallOfFame {
		if a.ID == id {
			return true
		}
	}
	return false
}
