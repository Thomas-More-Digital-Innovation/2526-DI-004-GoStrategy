// Package ai provides base interfaces and common functionality for AI implementations
package ai

import (
	"digital-innovation/gostrategy/internal/game"
	"digital-innovation/gostrategy/internal/game/models"
)

// AI is the interface that all AI implementations must satisfy.
// It extends the PlayerController interface
type AI interface {
	game.PlayerController
}

// BaseAI provides common functionality for all AI types
type BaseAI struct {
	player      *game.Player
	memory      *Memory
	recentMoves []game.Move
}

// NewBaseAI creates a new BaseAI instance
func NewBaseAI(player *game.Player, hasMemory bool) *BaseAI {
	var memory *Memory
	if hasMemory {
		memory = NewMemory()
	}
	return &BaseAI{
		player:      player,
		memory:      memory,
		recentMoves: make([]game.Move, 0, 16),
	}
}

// GetPlayer returns the player associated with the AI.
func (ai *BaseAI) GetPlayer() *game.Player {
	return ai.player
}

// GetControllerType returns the type of the AI controller, which is AIController.
func (ai *BaseAI) GetControllerType() game.ControllerType {
	return game.AIController
}

// GetMemory returns the AI's memory system (O(1) position lookup)
func (ai *BaseAI) GetMemory() *Memory {
	return ai.memory
}

// AnalyzeMove is called after opponent moves - override in subclasses for learning
// Default implementation updates memory automatically and detects scouts
func (ai *BaseAI) AnalyzeMove(move game.Move, opponent *game.Player, round int) {
	if ai.memory == nil {
		return
	}

	from := move.GetFrom()
	to := move.GetTo()

	ai.memory.MovePiece(from, to)

	deltaX := from.X - to.X
	if deltaX < 0 {
		deltaX = -deltaX
	}
	deltaY := from.Y - to.Y
	if deltaY < 0 {
		deltaY = -deltaY
	}

	if (deltaX > 1 || deltaY > 1) && opponent != nil {
		scoutPiece := game.NewPiece(models.Scout, opponent)
		ai.memory.Remember(to, scoutPiece, 1.0, round)
	}
}

// ObserveCombat is called when combat occurs - override for learning from reveals
// Default implementation updates memory with revealed pieces
func (ai *BaseAI) ObserveCombat(attackerPos, defenderPos game.Position, attackerPiece, defenderPiece *game.Piece, round int) {
	if ai.memory == nil {
		return
	}

	ai.memory.UpdateFromCombat(attackerPos, defenderPos, attackerPiece, defenderPiece, round)
}

// RecordOwnMove records a move made by this AI to enforce Two-Square repetition limits.
func (ai *BaseAI) RecordOwnMove(move game.Move) {
	ai.recentMoves = append(ai.recentMoves, move)
}

// GetRecentMoves returns moves executed by this AI.
func (ai *BaseAI) GetRecentMoves() []game.Move {
	return ai.recentMoves
}

// FilterTwoSquareMoves removes candidate moves that violate the Two-Square Rule.
func (ai *BaseAI) FilterTwoSquareMoves(moves []game.Move) []game.Move {
	if len(moves) == 0 || len(ai.recentMoves) == 0 {
		return moves
	}

	playerID := ai.player.GetID()
	filtered := make([]game.Move, 0, len(moves))
	for _, m := range moves {
		if !game.IsTwoSquareViolation(ai.recentMoves, playerID, m) {
			filtered = append(filtered, m)
		}
	}
	if len(filtered) == 0 {
		return moves
	}
	return filtered
}
