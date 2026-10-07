// Package aivsai provides tools for running AI vs AI matches and gathering statistics.
package aivsai

import (
	"digital-innovation/gostrategy/internal/ai"
	AIhandler "digital-innovation/gostrategy/internal/ai/handler"
	"digital-innovation/gostrategy/internal/game"
	"digital-innovation/gostrategy/internal/game/models"
	"fmt"
	"math/rand/v2"
)

func runAIvsAI(ai1, ai2 string, matches int, logging bool) models.AiGameSummary {
	return runAIvsAIWithOptions(ai1, ai2, matches, logging, nil, nil)
}

func runAIvsAISetup(ai1, ai2 string, matches int, logging bool, setup1, setup2 []string) models.AiGameSummary {
	return runAIvsAIWithOptionsAndSetups(ai1, ai2, matches, logging, nil, nil, setup1, setup2)
}

func runAIvsAIWithOptions(ai1, ai2 string, matches int, logging bool, opts1, opts2 map[string]any) models.AiGameSummary {
	return runAIvsAIWithOptionsAndSetups(ai1, ai2, matches, logging, opts1, opts2, nil, nil)
}

func runAIvsAIWithOptionsAndSetups(ai1, ai2 string, matches int, logging bool, opts1, opts2 map[string]any, setup1, setup2 []string) models.AiGameSummary {
	draws := 0

	flagCaptures := 0
	noMovesWins := 0
	maxTurnsWins := 0
	totalRounds := 0
	leastRounds := 1000

	player1Name := "Alice AI - " + ai1
	player2Name := "Bob AI - " + ai2

	player1Data := models.AiTournamentData{Name: player1Name}
	player2Data := models.AiTournamentData{Name: player2Name}

	for i := range matches {
		playerAlice := game.NewPlayer(0, player1Name, "red")
		playerBob := game.NewPlayer(1, player2Name, "blue")

		controllerAlice, err := AIhandler.CreateAIWithOptions(ai1, &playerAlice, opts1)
		if err != nil {
			panic(err)
		}
		controllerBob, err := AIhandler.CreateAIWithOptions(ai2, &playerBob, opts2)
		if err != nil {
			panic(err)
		}

		var g *game.Game
		if i%2 == 0 {
			p1Pieces := createPlayerPieces(&playerAlice, setup1, false)
			p2Pieces := createPlayerPieces(&playerBob, setup2, true)
			g = game.NewGame(controllerAlice, controllerBob)
			if err := game.SetupGame(g, p1Pieces, p2Pieces); err != nil {
				panic(err)
			}
			if logging {
				fmt.Printf("Game %3d (Alice starts): ", i+1)
			}
		} else {
			p1Pieces := createPlayerPieces(&playerBob, setup2, false)
			p2Pieces := createPlayerPieces(&playerAlice, setup1, true)
			g = game.NewGame(controllerBob, controllerAlice)
			if err := game.SetupGame(g, p1Pieces, p2Pieces); err != nil {
				panic(err)
			}
			if logging {
				fmt.Printf("Game %3d (Bob starts):   ", i+1)
			}
		}

		runner := game.NewRunner(g, 0, 1000)
		winner := runner.RunToCompletion()
		rounds := g.GetRound()
		totalRounds += rounds

		if rounds < leastRounds {
			leastRounds = rounds
		}

		if winner != nil {
			var winnerData *models.AiTournamentData
			winCause := g.GetWinCause()
			if winner.GetName() == player1Name {
				winnerData = &player1Data
				if logging {
					fmt.Printf("%v wins - %s (%d rounds)\n", player1Name, winCause, rounds)
				}
			} else {
				winnerData = &player2Data
				if logging {
					fmt.Printf("%v wins - %s (%d rounds)\n", player2Name, winCause, rounds)
				}
			}

			switch winCause {
			case game.WinCauseFlagCaptured:
				winnerData.WinCauseFlagCaptured++
				flagCaptures++
			case game.WinCauseNoMovablePieces:
				winnerData.WinCauseNoMovesWin++
				noMovesWins++
			default:
				winnerData.WinCauseMaxTurns++
				maxTurnsWins++
			}

			winnerData.Wins++
		} else {
			if logging {
				fmt.Printf("Draw after %d rounds\n", rounds)
			}
			draws++
		}
	}

	avgRounds := float64(totalRounds) / float64(matches)

	return models.AiGameSummary{
		Player1data:          player1Data,
		Player2data:          player2Data,
		Draws:                draws,
		TotalRounds:          totalRounds,
		AverageRounds:        avgRounds,
		LeastRounds:          leastRounds,
		Matches:              matches,
		WinCauseFlagCaptured: flagCaptures,
		WinCauseNoMovesWins:  noMovesWins,
		WinCauseMaxTurns:     maxTurnsWins,
	}
}

// TrainAI executes a hill-climbing optimization pipeline to train AI parameters.
func TrainAI(aiType string, opponentType string, generations int, matchesPerGen int, logging bool) error {
	params, err := ai.Load(aiType, "default")
	if err != nil {
		return err
	}

	for gen := range generations {
		if logging {
			fmt.Printf("--- Generation %d/%d ---\n", gen+1, generations)
		}

		candidate := &ai.Parameters{
			AIType: params.AIType,
			Name:   "candidate",
			//nolint:gosec
			Aggression: params.Aggression + (rand.Float64()-0.5)*0.2,
			Weights:    make(map[string]float64),
			Config:     make(map[string]any),
		}

		if candidate.Aggression < 0.0 {
			candidate.Aggression = 0.0
		} else if candidate.Aggression > 1.0 {
			candidate.Aggression = 1.0
		}

		for k, v := range params.Weights {
			//nolint:gosec
			candidate.Weights[k] = v + (rand.Float64()-0.5)*20.0
		}

		for k, v := range params.Config {
			candidate.Config[k] = v
		}

		err = ai.Save(candidate)
		if err != nil {
			return err
		}

		opts1 := map[string]any{"name": "candidate"}
		opts2 := map[string]any{"name": "default"}

		summary := runAIvsAIWithOptions(aiType, opponentType, matchesPerGen, false, opts1, opts2)

		candidateWins := summary.Player1data.Wins
		opponentWins := summary.Player2data.Wins

		if logging {
			fmt.Printf("Candidate Wins: %d, Opponent Wins: %d, Draws: %d\n", candidateWins, opponentWins, summary.Draws)
		}

		if candidateWins > opponentWins {
			if logging {
				fmt.Println("New parameters promoted to default!")
			}
			params.Aggression = candidate.Aggression
			for k, v := range candidate.Weights {
				params.Weights[k] = v
			}
			err = ai.Save(params)
			if err != nil {
				return err
			}
		} else if logging {
			fmt.Println("Candidate rejected.")
		}
	}

	return nil
}

func createPlayerPieces(player *game.Player, setup []string, flipped bool) []*game.Piece {
	if setup == nil {
		return game.RandomSetup(player)
	}
	targetSetup := setup
	if flipped {
		targetSetup = game.FlipSetupRows(setup)
	}
	pieces, err := game.ParseSetup(player, []byte(game.RowsToSetupData(targetSetup)))
	if err != nil {
		panic(fmt.Sprintf("failed to parse setup for player %s: %v", player.GetName(), err))
	}
	return pieces
}
