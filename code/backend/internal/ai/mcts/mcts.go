// Package mcts implements the Monte Carlo Tree Search strategy.
package mcts

import (
	"digital-innovation/gostrategy/internal/ai"
	"digital-innovation/gostrategy/internal/game"
	"digital-innovation/gostrategy/internal/game/models"
	"math"
	"runtime"
	"sync"
)

// AI implements the Monte Carlo Tree Search strategy using UCB1 selection.
type AI struct {
	ai.BaseAI
	params *ai.Parameters
}

// NewAI creates a new MCTS AI instance.
func NewAI(player *game.Player, hasMemory bool) *AI {
	params, _ := ai.Load(models.Mcts, "default")
	return &AI{
		BaseAI: *ai.NewBaseAI(player, hasMemory),
		params: params,
	}
}

// NewAIWithParams creates a new MCTS AI instance with predefined params.
func NewAIWithParams(player *game.Player, hasMemory bool, params *ai.Parameters) *AI {
	return &AI{
		BaseAI: *ai.NewBaseAI(player, hasMemory),
		params: params,
	}
}

// MakeMove implements the player controller interface selecting the best move using MCTS rollouts.
func (aiObj *AI) MakeMove(board *game.Board) game.Move {
	opponent := ai.GetOpponent(board, aiObj.GetPlayer().GetID())
	moves := ai.GetMoves(board, aiObj.GetPlayer())
	moves = aiObj.FilterTwoSquareMoves(moves)
	if len(moves) == 0 {
		return game.Move{}
	}
	if len(moves) == 1 {
		return moves[0]
	}

	budget := 500
	if bVal, ok := aiObj.params.Config["total_rollouts"].(float64); ok && bVal > 0 {
		budget = int(bVal)
	} else if iterVal, ok := aiObj.params.Config["iterations"].(float64); ok && iterVal > 0 {
		budget = min(int(iterVal)*len(moves), 1000)
	}

	explConstant := 1.414
	if cVal, ok := aiObj.params.Config["exploration_constant"].(float64); ok && cVal > 0 {
		explConstant = cVal
	}

	pool := newDeterminizationPool(board, aiObj.GetPlayer(), aiObj.GetMemory(), 4)

	numMoves := len(moves)
	visits := make([]int, numMoves)
	totalScores := make([]float64, numMoves)

	completed := aiObj.warmupRollouts(pool, moves, opponent, budget, visits, totalScores)
	aiObj.exploreUCB1(pool, moves, opponent, budget, explConstant, visits, totalScores, completed)

	return selectBestCandidate(moves, visits, totalScores)
}

func (aiObj *AI) warmupRollouts(
	pool *determinizationPool,
	moves []game.Move,
	opponent *game.Player,
	budget int,
	visits []int,
	totalScores []float64,
) int {
	numMoves := len(moves)
	warmupPerMove := 2
	if numMoves*warmupPerMove > budget {
		warmupPerMove = 1
	}

	completed := 0
	for i := range numMoves {
		for range warmupPerMove {
			score := aiObj.runRollout(pool, moves[i], aiObj.GetPlayer(), opponent)
			visits[i]++
			totalScores[i] += score
			completed++
		}
	}
	return completed
}

func (aiObj *AI) exploreUCB1(
	pool *determinizationPool,
	moves []game.Move,
	opponent *game.Player,
	budget int,
	explConstant float64,
	visits []int,
	totalScores []float64,
	initialCompleted int,
) {
	numMoves := len(moves)
	numWorkers := max(runtime.GOMAXPROCS(0), 1)
	totalCompleted := initialCompleted
	var mu sync.Mutex

	var wg sync.WaitGroup
	for range numWorkers {
		wg.Go(func() {
			for {
				mu.Lock()
				if totalCompleted >= budget {
					mu.Unlock()
					return
				}

				bestIdx := 0
				bestUCB := -1e9
				logN := math.Log(float64(totalCompleted + 1))

				for i := range numMoves {
					v := visits[i]
					mean := totalScores[i] / float64(v)
					ucb := mean + explConstant*math.Sqrt(logN/float64(v))
					if ucb > bestUCB {
						bestUCB = ucb
						bestIdx = i
					}
				}

				// virtual visit discourages worker collision
				visits[bestIdx]++
				m := moves[bestIdx]
				mu.Unlock()

				score := aiObj.runRollout(pool, m, aiObj.GetPlayer(), opponent)

				mu.Lock()
				totalScores[bestIdx] += score
				totalCompleted++
				mu.Unlock()
			}
		})
	}

	wg.Wait()
}

func selectBestCandidate(moves []game.Move, visits []int, totalScores []float64) game.Move {
	bestIdx := 0
	bestMean := -1e9
	for i := range moves {
		mean := totalScores[i] / float64(visits[i])
		if mean > bestMean {
			bestMean = mean
			bestIdx = i
		}
	}

	return moves[bestIdx]
}
