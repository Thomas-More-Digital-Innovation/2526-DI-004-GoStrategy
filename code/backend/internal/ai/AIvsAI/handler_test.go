package aivsai

import (
	"bytes"
	"digital-innovation/gostrategy/internal/game/models"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func captureOutput(fn func()) string {
	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		return ""
	}
	os.Stdout = w

	outChan := make(chan string)
	go func() {
		var buf bytes.Buffer
		_, _ = io.Copy(&buf, r)
		outChan <- buf.String()
	}()

	fn()

	_ = w.Close()
	os.Stdout = old
	return <-outChan
}

func TestRunAIvsAI(t *testing.T) {
	t.Run("default format without logging", func(t *testing.T) {
		output := captureOutput(func() {
			RunAIvsAI(models.Fato, models.Fato, 1, "default", false)
		})
		assert.Contains(t, output, "AI vs AI Tournament Summary (1 games)")
		assert.Contains(t, output, "Total Matches: 1")
		assert.Contains(t, output, "Player: Alice AI - fato")
		assert.Contains(t, output, "Player: Bob AI - fato")
	})

	t.Run("markdown format without logging", func(t *testing.T) {
		output := captureOutput(func() {
			RunAIvsAI(models.Fato, models.Fato, 1, "md", false)
		})
		assert.Contains(t, output, "### AI vs AI Tournament Summary (1 games)")
		assert.Contains(t, output, "| Player | Wins | Win % |")
		assert.Contains(t, output, "Alice AI - fato")
	})

	t.Run("default format with logging", func(t *testing.T) {
		output := captureOutput(func() {
			RunAIvsAI(models.Fato, models.Fato, 1, "default", true)
		})
		assert.Contains(t, output, "Game   1 (Alice starts):")
		assert.Contains(t, output, "AI vs AI Tournament Summary (1 games)")
	})
}

func TestRunAIvsAIWithSetups(t *testing.T) {
	t.Run("with nil setups and md format", func(t *testing.T) {
		output := captureOutput(func() {
			RunAIvsAIWithSetups(models.Fato, models.Fato, 1, "md", false, nil, nil)
		})
		assert.Contains(t, output, "### AI vs AI Tournament Summary (1 games)")
		assert.Contains(t, output, "**Total Matches:** 1")
	})

	t.Run("with default format fallback", func(t *testing.T) {
		output := captureOutput(func() {
			RunAIvsAIWithSetups(models.Fato, models.Fato, 1, "other", false, nil, nil)
		})
		assert.Contains(t, output, "========================================")
		assert.Contains(t, output, "Total Matches: 1")
	})
}

func TestPrintMarkdownSummary(t *testing.T) {
	t.Run("complete statistics with won matches", func(t *testing.T) {
		summary := models.AiGameSummary{
			Matches:              10,
			TotalRounds:          150,
			AverageRounds:        15.0,
			LeastRounds:          8,
			Draws:                2,
			WinCauseFlagCaptured: 4,
			WinCauseNoMovesWins:  3,
			WinCauseMaxTurns:     1,
			Player1data: models.AiTournamentData{
				Name:                 "Alice",
				Wins:                 5,
				WinCauseFlagCaptured: 3,
				WinCauseNoMovesWin:   1,
				WinCauseMaxTurns:     1,
			},
			Player2data: models.AiTournamentData{
				Name:                 "Bob",
				Wins:                 3,
				WinCauseFlagCaptured: 1,
				WinCauseNoMovesWin:   2,
				WinCauseMaxTurns:     0,
			},
		}

		output := captureOutput(func() {
			printMarkdownSummary(summary, 10)
		})

		assert.Contains(t, output, "### AI vs AI Tournament Summary (10 games)")
		assert.Contains(t, output, "**Total Matches:** 10")
		assert.Contains(t, output, "**Total Rounds:** 150")
		assert.Contains(t, output, "| Flag captured | 4 | 50.0% |")
		assert.Contains(t, output, "| No movable pieces | 3 | 37.5% |")
		assert.Contains(t, output, "| Max turns | 1 | 12.5% |")
		assert.Contains(t, output, "| Alice | 5 | 50.0% | 3 | 1 | 1 |")
		assert.Contains(t, output, "| Bob | 3 | 30.0% | 1 | 2 | 0 |")
		assert.Contains(t, output, "**Draws:** 2 (20.0%)")
	})

	t.Run("zero matches edge case", func(t *testing.T) {
		output := captureOutput(func() {
			printMarkdownSummary(models.AiGameSummary{}, 0)
		})
		assert.Contains(t, output, "### AI vs AI Tournament Summary (0 games)")
		assert.Contains(t, output, "| Flag captured | 0 | 0.0% |")
		assert.Contains(t, output, "**Draws:** 0 (0.0%)")
	})

	t.Run("all draws edge case", func(t *testing.T) {
		summary := models.AiGameSummary{Draws: 5}
		output := captureOutput(func() {
			printMarkdownSummary(summary, 5)
		})
		assert.Contains(t, output, "**Draws:** 5 (100.0%)")
		assert.Contains(t, output, "| Flag captured | 0 | 0.0% |")
	})
}

func TestPrintDefaultSummary(t *testing.T) {
	t.Run("standard tournament output", func(t *testing.T) {
		summary := models.AiGameSummary{
			Matches:              10,
			TotalRounds:          150,
			AverageRounds:        15.0,
			LeastRounds:          8,
			Draws:                2,
			WinCauseFlagCaptured: 4,
			WinCauseNoMovesWins:  3,
			WinCauseMaxTurns:     1,
			Player1data: models.AiTournamentData{
				Name:                 "Alice",
				Wins:                 5,
				WinCauseFlagCaptured: 3,
				WinCauseNoMovesWin:   1,
				WinCauseMaxTurns:     1,
			},
			Player2data: models.AiTournamentData{
				Name:                 "Bob",
				Wins:                 3,
				WinCauseFlagCaptured: 1,
				WinCauseNoMovesWin:   2,
				WinCauseMaxTurns:     0,
			},
		}

		output := captureOutput(func() {
			printDefaultSummary(summary, 10)
		})

		assert.Contains(t, output, "AI vs AI Tournament Summary (10 games)")
		assert.Contains(t, output, "Total Matches: 10")
		assert.Contains(t, output, "Total Rounds: 150")
		assert.Contains(t, output, "Flag captured:     4 (50.0%)")
		assert.Contains(t, output, "Player: Alice")
		assert.Contains(t, output, "Wins: 5 (50.0%)")
		assert.Contains(t, output, "Player: Bob")
		assert.Contains(t, output, "Wins: 3 (30.0%)")
		assert.Contains(t, output, "Draws: 2 (20.0%)")
	})

	t.Run("zero matches edge case", func(t *testing.T) {
		output := captureOutput(func() {
			printDefaultSummary(models.AiGameSummary{}, 0)
		})
		assert.Contains(t, output, "Total Matches: 0")
		assert.Contains(t, output, "Flag captured:     0 (0.0%)")
		assert.Contains(t, output, "Draws: 0 (0.0%)")
	})
}
