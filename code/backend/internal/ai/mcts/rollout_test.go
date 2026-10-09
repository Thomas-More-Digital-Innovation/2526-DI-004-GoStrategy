package mcts

import (
	"digital-innovation/gostrategy/internal/game"
	"digital-innovation/gostrategy/internal/game/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeterminizationPool(t *testing.T) {
	board := game.NewBoard()
	player := game.NewPlayer(0, "us", "red")

	t.Run("default size when non-positive", func(t *testing.T) {
		dp := newDeterminizationPool(board, &player, nil, 0)
		assert.NotNil(t, dp)
		assert.Len(t, dp.worlds, 4)
		assert.NotNil(t, dp.sample())
	})

	t.Run("custom pool size", func(t *testing.T) {
		dp := newDeterminizationPool(board, &player, nil, 2)
		assert.NotNil(t, dp)
		assert.Len(t, dp.worlds, 2)
		assert.NotNil(t, dp.sample())
	})

	t.Run("sample empty pool returns nil", func(t *testing.T) {
		dp := &determinizationPool{worlds: nil}
		assert.Nil(t, dp.sample())
	})
}

func TestRunRollout(t *testing.T) {
	player1 := game.NewPlayer(0, "us", "red")
	player2 := game.NewPlayer(1, "them", "blue")
	aiObj := NewAI(&player1, false)

	t.Run("nil pool sample returns neutral score", func(t *testing.T) {
		score := aiObj.runRollout(&determinizationPool{worlds: nil}, game.Move{}, &player1, &player2)
		assert.Equal(t, 0.5, score)
	})

	t.Run("root move captures opponent flag", func(t *testing.T) {
		board := game.NewBoard()
		p1Marshal := game.NewPiece(models.Marshal, &player1)
		p2Flag := game.NewPiece(models.Flag, &player2)
		p1Flag := game.NewPiece(models.Flag, &player1)

		board.SetPieceAt(game.NewPosition(0, 0), p1Marshal)
		board.SetPieceAt(game.NewPosition(0, 1), p2Flag)
		board.SetPieceAt(game.NewPosition(9, 9), p1Flag)
		player1.AddPiece(p1Marshal, game.NewPosition(0, 0))
		player2.AddPiece(p2Flag, game.NewPosition(0, 1))
		player1.AddPiece(p1Flag, game.NewPosition(9, 9))

		pool := &determinizationPool{worlds: []*game.Board{board}}
		m := game.NewMove(game.NewPosition(0, 0), game.NewPosition(0, 1), &player1)
		score := aiObj.runRollout(pool, m, &player1, &player2)
		assert.Equal(t, 1.0, score)
	})

	t.Run("our flag missing returns zero score", func(t *testing.T) {
		board := game.NewBoard()
		p1Marshal := game.NewPiece(models.Marshal, &player1)
		p2Flag := game.NewPiece(models.Flag, &player2)

		board.SetPieceAt(game.NewPosition(0, 0), p1Marshal)
		board.SetPieceAt(game.NewPosition(5, 5), p2Flag)
		player1.AddPiece(p1Marshal, game.NewPosition(0, 0))
		player2.AddPiece(p2Flag, game.NewPosition(5, 5))

		pool := &determinizationPool{worlds: []*game.Board{board}}
		m := game.NewMove(game.NewPosition(0, 0), game.NewPosition(0, 1), &player1)
		score := aiObj.runRollout(pool, m, &player1, &player2)
		assert.Equal(t, 0.0, score)
	})

	t.Run("opponent has no mobile pieces returns win", func(t *testing.T) {
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

		pool := &determinizationPool{worlds: []*game.Board{board}}
		m := game.NewMove(game.NewPosition(0, 0), game.NewPosition(0, 1), &player1)
		score := aiObj.runRollout(pool, m, &player1, &player2)
		assert.Equal(t, 1.0, score)
	})

	t.Run("rollout reaching depth evaluates bounded score", func(t *testing.T) {
		board := game.NewBoard()
		p1Scout := game.NewPiece(models.Scout, &player1)
		p2Scout := game.NewPiece(models.Scout, &player2)
		p1Flag := game.NewPiece(models.Flag, &player1)
		p2Flag := game.NewPiece(models.Flag, &player2)

		board.SetPieceAt(game.NewPosition(0, 0), p1Scout)
		board.SetPieceAt(game.NewPosition(9, 0), p2Scout)
		board.SetPieceAt(game.NewPosition(0, 9), p1Flag)
		board.SetPieceAt(game.NewPosition(9, 9), p2Flag)
		player1.AddPiece(p1Scout, game.NewPosition(0, 0))
		player2.AddPiece(p2Scout, game.NewPosition(9, 0))
		player1.AddPiece(p1Flag, game.NewPosition(0, 9))
		player2.AddPiece(p2Flag, game.NewPosition(9, 9))

		pool := &determinizationPool{worlds: []*game.Board{board}}
		m := game.NewMove(game.NewPosition(0, 0), game.NewPosition(0, 1), &player1)
		score := aiObj.runRollout(pool, m, &player1, &player2)
		assert.GreaterOrEqual(t, score, 0.0)
		assert.LessOrEqual(t, score, 1.0)
	})
}

func TestPickRolloutMoves(t *testing.T) {
	player1 := game.NewPlayer(0, "us", "red")
	player2 := game.NewPlayer(1, "them", "blue")

	t.Run("pickFastRolloutMove with empty positions", func(t *testing.T) {
		b := game.NewBoard()
		move, ok := pickFastRolloutMove(b, nil, &player1)
		assert.False(t, ok)
		assert.Equal(t, game.Move{}, move)
	})

	t.Run("pickFastRolloutMove with valid capture", func(t *testing.T) {
		b := game.NewBoard()
		p1 := game.NewPiece(models.Marshal, &player1)
		p2 := game.NewPiece(models.General, &player2)
		b.SetPieceAt(game.NewPosition(0, 0), p1)
		b.SetPieceAt(game.NewPosition(0, 1), p2)
		player1.AddPiece(p1, game.NewPosition(0, 0))
		player2.AddPiece(p2, game.NewPosition(0, 1))

		move, ok := pickFastRolloutMove(b, []game.Position{game.NewPosition(0, 0)}, &player1)
		assert.True(t, ok)
		assert.Equal(t, game.NewPosition(0, 0), move.GetFrom())
	})

	t.Run("pickCaptureRolloutMove returns false when no captures exist", func(t *testing.T) {
		b := game.NewBoard()
		p1 := game.NewPiece(models.Marshal, &player1)
		b.SetPieceAt(game.NewPosition(0, 0), p1)
		player1.AddPiece(p1, game.NewPosition(0, 0))

		move, ok := pickCaptureRolloutMove(b, []game.Position{game.NewPosition(0, 0)}, &player1)
		assert.False(t, ok)
		assert.Equal(t, game.Move{}, move)
	})

	t.Run("pickRandomRolloutMove returns valid move", func(t *testing.T) {
		b := game.NewBoard()
		p1 := game.NewPiece(models.Marshal, &player1)
		b.SetPieceAt(game.NewPosition(0, 0), p1)
		player1.AddPiece(p1, game.NewPosition(0, 0))

		move, ok := pickRandomRolloutMove(b, []game.Position{game.NewPosition(0, 0)}, &player1)
		assert.True(t, ok)
		assert.Equal(t, game.NewPosition(0, 0), move.GetFrom())
	})

	t.Run("pickRandomRolloutMove returns false when piece has no moves", func(t *testing.T) {
		b := game.NewBoard()
		p1 := game.NewPiece(models.Flag, &player1)
		b.SetPieceAt(game.NewPosition(0, 0), p1)
		player1.AddPiece(p1, game.NewPosition(0, 0))

		move, ok := pickRandomRolloutMove(b, []game.Position{game.NewPosition(0, 0)}, &player1)
		assert.False(t, ok)
		assert.Equal(t, game.Move{}, move)
	})
}

func TestFlagPositionAndCapture(t *testing.T) {
	player1 := game.NewPlayer(0, "us", "red")
	player2 := game.NewPlayer(1, "them", "blue")
	aiObj := NewAI(&player1, false)

	board := game.NewBoard()
	p1Flag := game.NewPiece(models.Flag, &player1)
	p2Flag := game.NewPiece(models.Flag, &player2)
	board.SetPieceAt(game.NewPosition(0, 0), p1Flag)
	board.SetPieceAt(game.NewPosition(0, 1), p2Flag)
	player1.AddPiece(p1Flag, game.NewPosition(0, 0))
	player2.AddPiece(p2Flag, game.NewPosition(0, 1))

	assert.Equal(t, game.NewPosition(-1, -1), aiObj.findFlagPosition(board, nil))
	assert.Equal(t, game.NewPosition(0, 0), aiObj.findFlagPosition(board, &player1))

	boardNoFlag := game.NewBoard()
	assert.Equal(t, game.NewPosition(-1, -1), aiObj.findFlagPosition(boardNoFlag, &player1))

	assert.True(t, aiObj.isFlagCapturedAt(board, game.NewPosition(-1, -1), &player1))
	assert.False(t, aiObj.isFlagCapturedAt(board, game.NewPosition(0, 0), &player1))
	assert.True(t, aiObj.isFlagCapturedAt(board, game.NewPosition(0, 1), &player1))
	assert.True(t, aiObj.isFlagCapturedAt(boardNoFlag, game.NewPosition(0, 0), &player1))
}

func TestSimulateMoveInPlace(t *testing.T) {
	player1 := game.NewPlayer(0, "us", "red")
	player2 := game.NewPlayer(1, "them", "blue")
	aiObj := NewAI(&player1, false)

	t.Run("nil attacker", func(t *testing.T) {
		b := game.NewBoard()
		captured := aiObj.applySimulatedMoveInPlace(b, game.NewMove(game.NewPosition(0, 0), game.NewPosition(0, 1), &player1))
		assert.False(t, captured)
	})

	t.Run("nil target", func(t *testing.T) {
		b := game.NewBoard()
		p1 := game.NewPiece(models.General, &player1)
		b.SetPieceAt(game.NewPosition(0, 0), p1)
		captured := aiObj.applySimulatedMoveInPlace(b, game.NewMove(game.NewPosition(0, 0), game.NewPosition(0, 1), &player1))
		assert.False(t, captured)
		assert.Equal(t, p1, b.GetPieceAt(game.NewPosition(0, 1)))
		assert.Nil(t, b.GetPieceAt(game.NewPosition(0, 0)))
	})

	t.Run("target is flag", func(t *testing.T) {
		b := game.NewBoard()
		p1 := game.NewPiece(models.General, &player1)
		p2 := game.NewPiece(models.Flag, &player2)
		b.SetPieceAt(game.NewPosition(0, 0), p1)
		b.SetPieceAt(game.NewPosition(0, 1), p2)
		captured := aiObj.applySimulatedMoveInPlace(b, game.NewMove(game.NewPosition(0, 0), game.NewPosition(0, 1), &player1))
		assert.False(t, captured)
		assert.Equal(t, p1, b.GetPieceAt(game.NewPosition(0, 1)))
	})

	t.Run("spy vs marshal", func(t *testing.T) {
		b := game.NewBoard()
		p1 := game.NewPiece(models.Spy, &player1)
		p2 := game.NewPiece(models.Marshal, &player2)
		b.SetPieceAt(game.NewPosition(0, 0), p1)
		b.SetPieceAt(game.NewPosition(0, 1), p2)
		captured := aiObj.applySimulatedMoveInPlace(b, game.NewMove(game.NewPosition(0, 0), game.NewPosition(0, 1), &player1))
		assert.False(t, captured)
		assert.Equal(t, p1, b.GetPieceAt(game.NewPosition(0, 1)))
	})

	t.Run("miner defuses bomb", func(t *testing.T) {
		b := game.NewBoard()
		p1 := game.NewPiece(models.Miner, &player1)
		p2 := game.NewPiece(models.Bomb, &player2)
		b.SetPieceAt(game.NewPosition(0, 0), p1)
		b.SetPieceAt(game.NewPosition(0, 1), p2)
		captured := aiObj.applySimulatedMoveInPlace(b, game.NewMove(game.NewPosition(0, 0), game.NewPosition(0, 1), &player1))
		assert.False(t, captured)
		assert.Equal(t, p1, b.GetPieceAt(game.NewPosition(0, 1)))
	})

	t.Run("non-miner eliminated by bomb", func(t *testing.T) {
		b := game.NewBoard()
		p1 := game.NewPiece(models.General, &player1)
		p2 := game.NewPiece(models.Bomb, &player2)
		b.SetPieceAt(game.NewPosition(0, 0), p1)
		b.SetPieceAt(game.NewPosition(0, 1), p2)
		captured := aiObj.applySimulatedMoveInPlace(b, game.NewMove(game.NewPosition(0, 0), game.NewPosition(0, 1), &player1))
		assert.True(t, captured)
		assert.Nil(t, b.GetPieceAt(game.NewPosition(0, 0)))
		assert.Equal(t, p2, b.GetPieceAt(game.NewPosition(0, 1)))
	})

	t.Run("attacker vs defender ranks", func(t *testing.T) {
		b := game.NewBoard()
		p1 := game.NewPiece(models.General, &player1)
		p2 := game.NewPiece(models.Miner, &player2)
		b.SetPieceAt(game.NewPosition(0, 0), p1)
		b.SetPieceAt(game.NewPosition(0, 1), p2)
		captured := aiObj.applySimulatedMoveInPlace(b, game.NewMove(game.NewPosition(0, 0), game.NewPosition(0, 1), &player1))
		assert.False(t, captured)
		assert.Equal(t, p1, b.GetPieceAt(game.NewPosition(0, 1)))

		b = game.NewBoard()
		p1 = game.NewPiece(models.Miner, &player1)
		p2 = game.NewPiece(models.General, &player2)
		b.SetPieceAt(game.NewPosition(0, 0), p1)
		b.SetPieceAt(game.NewPosition(0, 1), p2)
		captured = aiObj.applySimulatedMoveInPlace(b, game.NewMove(game.NewPosition(0, 0), game.NewPosition(0, 1), &player1))
		assert.True(t, captured)
		assert.Nil(t, b.GetPieceAt(game.NewPosition(0, 0)))
		assert.Equal(t, p2, b.GetPieceAt(game.NewPosition(0, 1)))

		b = game.NewBoard()
		p1 = game.NewPiece(models.Miner, &player1)
		p2 = game.NewPiece(models.Miner, &player2)
		b.SetPieceAt(game.NewPosition(0, 0), p1)
		b.SetPieceAt(game.NewPosition(0, 1), p2)
		captured = aiObj.applySimulatedMoveInPlace(b, game.NewMove(game.NewPosition(0, 0), game.NewPosition(0, 1), &player1))
		assert.True(t, captured)
		assert.Nil(t, b.GetPieceAt(game.NewPosition(0, 0)))
		assert.Nil(t, b.GetPieceAt(game.NewPosition(0, 1)))
	})
}

func TestUpdateIndex(t *testing.T) {
	idx := []game.Position{
		game.NewPosition(0, 0),
		game.NewPosition(1, 1),
	}

	updateIndex(&idx, game.NewPosition(0, 0), game.NewPosition(0, 1), false)
	assert.Len(t, idx, 2)
	assert.Contains(t, idx, game.NewPosition(0, 1))

	updateIndex(&idx, game.NewPosition(0, 1), game.NewPosition(0, 2), true)
	assert.Len(t, idx, 1)
	assert.NotContains(t, idx, game.NewPosition(0, 1))
	assert.NotContains(t, idx, game.NewPosition(0, 2))

	updateIndex(&idx, game.NewPosition(5, 5), game.NewPosition(5, 6), false)
	assert.Len(t, idx, 1)
}
