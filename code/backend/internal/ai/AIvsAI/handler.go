package aivsai

import (
	"digital-innovation/gostrategy/internal/game/models"
	"fmt"
)

// RunAIvsAI runs a tournament between two AI types and prints the results
func RunAIvsAI(ai1, ai2 string, matches int, format string, logging bool) {
	RunAIvsAIWithSetups(ai1, ai2, matches, format, logging, nil, nil)
}

// RunAIvsAIWithSetups runs a tournament between two AI types with custom board setups and prints the results
func RunAIvsAIWithSetups(ai1, ai2 string, matches int, format string, logging bool, setup1, setup2 []string) {
	summary := runAIvsAISetup(ai1, ai2, matches, logging, setup1, setup2)

	switch format {
	case "md":
		printMarkdownSummary(summary, matches)
	default:
		printDefaultSummary(summary, matches)
	}
}

func printMarkdownSummary(summary models.AiGameSummary, matches int) {
	// Top-level summary
	fmt.Printf("\n### AI vs AI Tournament Summary (%d games)\n\n", matches)
	fmt.Printf("**Total Matches:** %d  \n**Total Rounds:** %d  \n**Average Rounds (per game):** %.2f  \n**Shortest Game (rounds):** %d\n\n",
		matches, summary.TotalRounds, summary.AverageRounds, summary.LeastRounds)

	totalMatches := float64(matches)
	wonMatches := float64(matches - summary.Draws)

	winCausePct := func(count int) float64 {
		if wonMatches > 0 {
			return float64(count) * 100.0 / wonMatches
		}
		return 0.0
	}

	matchPct := func(count int) float64 {
		if totalMatches > 0 {
			return float64(count) * 100.0 / totalMatches
		}
		return 0.0
	}

	// Overall win causes (aggregate)
	totalFlag := summary.WinCauseFlagCaptured
	totalNoMoves := summary.WinCauseNoMovesWins
	totalMaxTurns := summary.WinCauseMaxTurns

	fmt.Println("#### Overall Win Causes")
	fmt.Println()
	fmt.Println("| Cause | Count | % |")
	fmt.Println("|-------:|------:|---:|")
	fmt.Printf("| Flag captured | %d | %.1f%% |\n", totalFlag, winCausePct(totalFlag))
	fmt.Printf("| No movable pieces | %d | %.1f%% |\n", totalNoMoves, winCausePct(totalNoMoves))
	fmt.Printf("| Max turns | %d | %.1f%% |\n\n", totalMaxTurns, winCausePct(totalMaxTurns))

	// Per-player summary table
	p1 := summary.Player1data
	p2 := summary.Player2data

	fmt.Println("#### Player Results")
	fmt.Println()
	fmt.Println("| Player | Wins | Win % | Flag captures | No-move wins | Max-turn wins |")
	fmt.Println("|:-------|-----:|-----:|--------------:|-------------:|--------------:|")
	fmt.Printf("| %s | %d | %.1f%% | %d | %d | %d |\n",
		p1.Name, p1.Wins, matchPct(p1.Wins), p1.WinCauseFlagCaptured, p1.WinCauseNoMovesWin, p1.WinCauseMaxTurns)
	fmt.Printf("| %s | %d | %.1f%% | %d | %d | %d |\n\n",
		p2.Name, p2.Wins, matchPct(p2.Wins), p2.WinCauseFlagCaptured, p2.WinCauseNoMovesWin, p2.WinCauseMaxTurns)

	// Draws
	fmt.Printf("**Draws:** %d (%.1f%%)\n", summary.Draws, matchPct(summary.Draws))
}

func printDefaultSummary(summary models.AiGameSummary, matches int) {
	// Human-readable plain text summary
	fmt.Println()
	fmt.Println("========================================")
	fmt.Printf("AI vs AI Tournament Summary (%d games)\n", matches)
	fmt.Println("========================================")
	fmt.Printf("Total Matches: %d\n", matches)
	fmt.Printf("Total Rounds: %d\n", summary.TotalRounds)
	fmt.Printf("Average Rounds (per game): %.2f\n", summary.AverageRounds)
	fmt.Printf("Shortest Game (rounds): %d\n", summary.LeastRounds)
	fmt.Println("----------------------------------------")

	totalMatches := float64(matches)
	wonMatches := float64(matches - summary.Draws)

	winCausePct := func(count int) float64 {
		if wonMatches > 0 {
			return float64(count) * 100.0 / wonMatches
		}
		return 0.0
	}

	matchPct := func(count int) float64 {
		if totalMatches > 0 {
			return float64(count) * 100.0 / totalMatches
		}
		return 0.0
	}

	// Overall win causes
	fmt.Println("Overall Win Causes:")
	totalFlag := summary.WinCauseFlagCaptured
	totalNoMoves := summary.WinCauseNoMovesWins
	totalMaxTurns := summary.WinCauseMaxTurns

	fmt.Printf("  Flag captured:     %d (%.1f%%)\n", totalFlag, winCausePct(totalFlag))
	fmt.Printf("  No movable pieces: %d (%.1f%%)\n", totalNoMoves, winCausePct(totalNoMoves))
	fmt.Printf("  Max turns:         %d (%.1f%%)\n", totalMaxTurns, winCausePct(totalMaxTurns))
	fmt.Println("----------------------------------------")

	// Per-player breakdown
	p1 := summary.Player1data
	p2 := summary.Player2data
	fmt.Printf("Player: %s\n", p1.Name)
	fmt.Printf("  Wins: %d (%.1f%%)\n", p1.Wins, matchPct(p1.Wins))
	fmt.Printf("  Win causes:\n")
	fmt.Printf("    Flag captured:     %d\n", p1.WinCauseFlagCaptured)
	fmt.Printf("    No movable pieces: %d\n", p1.WinCauseNoMovesWin)
	fmt.Printf("    Max turns:         %d\n", p1.WinCauseMaxTurns)
	fmt.Println("----------------------------------------")
	fmt.Printf("Player: %s\n", p2.Name)
	fmt.Printf("  Wins: %d (%.1f%%)\n", p2.Wins, matchPct(p2.Wins))
	fmt.Printf("  Win causes:\n")
	fmt.Printf("    Flag captured:     %d\n", p2.WinCauseFlagCaptured)
	fmt.Printf("    No movable pieces: %d\n", p2.WinCauseNoMovesWin)
	fmt.Printf("    Max turns:         %d\n", p2.WinCauseMaxTurns)
	fmt.Println("----------------------------------------")

	fmt.Printf("Draws: %d (%.1f%%)\n", summary.Draws, matchPct(summary.Draws))
}
