package mcts

import (
	"digital-innovation/gostrategy/internal/ai"
	ai_const "digital-innovation/gostrategy/internal/ai/const"
	"digital-innovation/gostrategy/internal/game"
	"digital-innovation/gostrategy/internal/game/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMctsAI_Constructors(t *testing.T) {
	player1 := game.NewPlayer(0, "us", "red")

	t.Run("default constructor with and without memory", func(t *testing.T) {
		aiNoMem := NewAI(&player1, false)
		assert.NotNil(t, aiNoMem)
		assert.Equal(t, &player1, aiNoMem.GetPlayer())
		assert.Nil(t, aiNoMem.GetMemory())

		aiWithMem := NewAI(&player1, true)
		assert.NotNil(t, aiWithMem)
		assert.NotNil(t, aiWithMem.GetMemory())
	})

	t.Run("constructor with predefined params", func(t *testing.T) {
		params := &ai.Parameters{
			Weights: map[string]float64{
				ai_const.Marshal: 100.0,
			},
			Aggression: 0.5,
			Config: map[string]interface{}{
				"iterations":           float64(5),
				"exploration_constant": 1.414,
			},
		}

		aiObj := NewAIWithParams(&player1, true, params)
		assert.NotNil(t, aiObj)
		assert.Equal(t, params, aiObj.params)
	})
}

func TestMctsAI_MakeMove(t *testing.T) {
	player1 := game.NewPlayer(0, "us", "red")
	player2 := game.NewPlayer(1, "them", "blue")

	t.Run("returns empty move when no legal moves", func(t *testing.T) {
		aiObj := NewAI(&player1, false)
		board := game.NewBoard()

		move := aiObj.MakeMove(board)
		assert.Equal(t, game.Move{}, move)
	})

	t.Run("returns single legal move immediately", func(t *testing.T) {
		aiObj := NewAI(&player1, false)
		board := game.NewBoard()

		p1 := game.NewPiece(models.Marshal, &player1)
		board.SetPieceAt(game.NewPosition(0, 0), p1)
		player1.AddPiece(p1, game.NewPosition(0, 0))

		pFlag := game.NewPiece(models.Flag, &player1)
		board.SetPieceAt(game.NewPosition(1, 0), pFlag)
		player1.AddPiece(pFlag, game.NewPosition(1, 0))

		move := aiObj.MakeMove(board)
		assert.Equal(t, game.NewPosition(0, 0), move.GetFrom())
		assert.Equal(t, game.NewPosition(0, 1), move.GetTo())
	})

	t.Run("selects candidate move with total_rollouts config", func(t *testing.T) {
		params := &ai.Parameters{
			Weights: map[string]float64{
				ai_const.Marshal: 100.0,
				ai_const.Flag:    10000.0,
			},
			Aggression: 0.5,
			Config: map[string]interface{}{
				"total_rollouts": float64(10),
			},
		}
		aiObj := NewAIWithParams(&player1, false, params)
		board := game.NewBoard()

		p1Marshal := game.NewPiece(models.Marshal, &player1)
		p1Flag := game.NewPiece(models.Flag, &player1)
		p2Flag := game.NewPiece(models.Flag, &player2)

		board.SetPieceAt(game.NewPosition(0, 0), p1Marshal)
		board.SetPieceAt(game.NewPosition(9, 9), p1Flag)
		board.SetPieceAt(game.NewPosition(5, 5), p2Flag)
		player1.AddPiece(p1Marshal, game.NewPosition(0, 0))
		player1.AddPiece(p1Flag, game.NewPosition(9, 9))
		player2.AddPiece(p2Flag, game.NewPosition(5, 5))

		move := aiObj.MakeMove(board)
		assert.NotEqual(t, game.Move{}, move)
	})

	t.Run("selects candidate move with iterations config", func(t *testing.T) {
		params := &ai.Parameters{
			Weights:    map[string]float64{ai_const.Marshal: 100.0},
			Aggression: 0.5,
			Config: map[string]interface{}{
				"iterations":           float64(2),
				"exploration_constant": 2.0,
			},
		}
		aiObj := NewAIWithParams(&player1, false, params)
		board := game.NewBoard()

		p1Marshal := game.NewPiece(models.Marshal, &player1)
		p1Flag := game.NewPiece(models.Flag, &player1)
		p2Flag := game.NewPiece(models.Flag, &player2)

		board.SetPieceAt(game.NewPosition(0, 0), p1Marshal)
		board.SetPieceAt(game.NewPosition(9, 9), p1Flag)
		board.SetPieceAt(game.NewPosition(5, 5), p2Flag)
		player1.AddPiece(p1Marshal, game.NewPosition(0, 0))
		player1.AddPiece(p1Flag, game.NewPosition(9, 9))
		player2.AddPiece(p2Flag, game.NewPosition(5, 5))

		move := aiObj.MakeMove(board)
		assert.NotEqual(t, game.Move{}, move)
	})
}

func TestMctsAI_WarmupRollouts(t *testing.T) {
	player1 := game.NewPlayer(0, "us", "red")
	player2 := game.NewPlayer(1, "them", "blue")
	aiObj := NewAI(&player1, false)

	board := game.NewBoard()
	p1Flag := game.NewPiece(models.Flag, &player1)
	p2Flag := game.NewPiece(models.Flag, &player2)
	board.SetPieceAt(game.NewPosition(0, 0), p1Flag)
	board.SetPieceAt(game.NewPosition(9, 9), p2Flag)
	player1.AddPiece(p1Flag, game.NewPosition(0, 0))
	player2.AddPiece(p2Flag, game.NewPosition(9, 9))

	pool := &determinizationPool{worlds: []*game.Board{board}}
	moves := []game.Move{
		game.NewMove(game.NewPosition(1, 1), game.NewPosition(1, 2), &player1),
		game.NewMove(game.NewPosition(2, 2), game.NewPosition(2, 3), &player1),
	}

	t.Run("sufficient budget runs two warmups per move", func(t *testing.T) {
		visits := make([]int, len(moves))
		scores := make([]float64, len(moves))
		completed := aiObj.warmupRollouts(pool, moves, &player2, 10, visits, scores)

		assert.Equal(t, 4, completed)
		assert.Equal(t, []int{2, 2}, visits)
	})

	t.Run("constrained budget runs one warmup per move", func(t *testing.T) {
		visits := make([]int, len(moves))
		scores := make([]float64, len(moves))
		completed := aiObj.warmupRollouts(pool, moves, &player2, 2, visits, scores)

		assert.Equal(t, 2, completed)
		assert.Equal(t, []int{1, 1}, visits)
	})
}

func TestMctsAI_ExploreUCB1(t *testing.T) {
	player1 := game.NewPlayer(0, "us", "red")
	player2 := game.NewPlayer(1, "them", "blue")
	aiObj := NewAI(&player1, false)

	board := game.NewBoard()
	p1Flag := game.NewPiece(models.Flag, &player1)
	p2Flag := game.NewPiece(models.Flag, &player2)
	board.SetPieceAt(game.NewPosition(0, 0), p1Flag)
	board.SetPieceAt(game.NewPosition(9, 9), p2Flag)
	player1.AddPiece(p1Flag, game.NewPosition(0, 0))
	player2.AddPiece(p2Flag, game.NewPosition(9, 9))

	pool := &determinizationPool{worlds: []*game.Board{board}}
	moves := []game.Move{
		game.NewMove(game.NewPosition(1, 1), game.NewPosition(1, 2), &player1),
		game.NewMove(game.NewPosition(2, 2), game.NewPosition(2, 3), &player1),
	}

	visits := []int{2, 2}
	scores := []float64{1.0, 1.0}
	budget := 10

	aiObj.exploreUCB1(pool, moves, &player2, budget, 1.414, visits, scores, 4)

	totalVisits := visits[0] + visits[1]
	assert.GreaterOrEqual(t, totalVisits, budget)
}

func TestSelectBestCandidate(t *testing.T) {
	player1 := game.NewPlayer(0, "us", "red")
	m1 := game.NewMove(game.NewPosition(0, 0), game.NewPosition(0, 1), &player1)
	m2 := game.NewMove(game.NewPosition(1, 0), game.NewPosition(1, 1), &player1)
	m3 := game.NewMove(game.NewPosition(2, 0), game.NewPosition(2, 1), &player1)

	moves := []game.Move{m1, m2, m3}
	visits := []int{10, 10, 5}
	scores := []float64{3.0, 8.5, 1.0}

	best := selectBestCandidate(moves, visits, scores)
	assert.Equal(t, m2, best)
}

func BenchmarkMctsAI(b *testing.B) {
	player1 := game.NewPlayer(0, "Piet", "red")
	aiObj := NewAI(&player1, false)
	board := game.NewBoard()

	p1 := game.NewPiece(models.Marshal, &player1)
	board.SetPieceAt(game.NewPosition(4, 4), p1)
	player1.AddPiece(p1, game.NewPosition(4, 4))

	for b.Loop() {
		_ = aiObj.MakeMove(board)
	}
}
