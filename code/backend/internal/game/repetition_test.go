package game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTwoSquareViolation(t *testing.T) {
	player1 := NewPlayer(0, "Alice", "red")
	player2 := NewPlayer(1, "Bob", "blue")

	posA := NewPosition(1, 1)
	posB := NewPosition(1, 2)
	posC := NewPosition(2, 1)

	otherFrom := NewPosition(5, 5)
	otherTo := NewPosition(5, 6)

	t.Run("empty history has no violation", func(t *testing.T) {
		move := NewMove(posA, posB, &player1)
		assert.False(t, IsTwoSquareViolation(nil, 0, move))
	})

	t.Run("up to 3 full cycles are permitted", func(t *testing.T) {
		var history []Move
		// 3 full cycles: A->B, B->A x 3
		for cycle := range 3 {
			// A -> B
			moveAB := NewMove(posA, posB, &player1)
			assert.False(t, IsTwoSquareViolation(history, 0, moveAB), "A->B cycle %d should be valid", cycle+1)
			history = append(history, moveAB)

			// Opponent move interleaved
			history = append(history, NewMove(NewPosition(9, 9), NewPosition(9, 8), &player2))

			// B -> A
			moveBA := NewMove(posB, posA, &player1)
			assert.False(t, IsTwoSquareViolation(history, 0, moveBA), "B->A cycle %d should be valid", cycle+1)
			history = append(history, moveBA)

			// Opponent move interleaved
			history = append(history, NewMove(NewPosition(9, 8), NewPosition(9, 9), &player2))
		}

		// Now piece is at A. 3 full cycles have occurred (6 player moves).
		// Attempting A -> B a 4th time violates the Two-Square Rule.
		move4thAB := NewMove(posA, posB, &player1)
		assert.True(t, IsTwoSquareViolation(history, 0, move4thAB), "4th A->B move must be flagged as violation")

		// Moving to a different square C is legal
		moveAC := NewMove(posA, posC, &player1)
		assert.False(t, IsTwoSquareViolation(history, 0, moveAC), "moving to a third square C must be legal")
	})

	t.Run("moving another piece breaks the chain", func(t *testing.T) {
		var history []Move
		// 2 full cycles
		for range 2 {
			history = append(history, NewMove(posA, posB, &player1))
			history = append(history, NewMove(posB, posA, &player1))
		}

		// Player moves a different piece
		history = append(history, NewMove(otherFrom, otherTo, &player1))

		// Cycle count should have reset; doing A->B should now be legal
		moveAB := NewMove(posA, posB, &player1)
		assert.False(t, IsTwoSquareViolation(history, 0, moveAB))
	})

	t.Run("Game FilterTwoSquareMoves removes only violating moves", func(t *testing.T) {
		g := &Game{
			MoveHistory: []Move{},
		}

		for range 3 {
			g.MoveHistory = append(g.MoveHistory, NewMove(posA, posB, &player1))
			g.MoveHistory = append(g.MoveHistory, NewMove(posB, posA, &player1))
		}

		candidates := []Move{
			NewMove(posA, posB, &player1), // violating
			NewMove(posA, posC, &player1), // legal
		}

		filtered := g.FilterTwoSquareMoves(0, candidates)
		assert.Len(t, filtered, 1)
		assert.Equal(t, posC, filtered[0].GetTo())
	})

	t.Run("Scout boundary crossing with varied distances", func(t *testing.T) {
		var history []Move

		// Scout oscillates across boundary between (1,3) and (1,4) with varied endpoints:
		// 1: (1,1) -> (1,5) [crosses 3|4 forward]
		// 2: (1,5) -> (1,2) [crosses 3|4 backward]
		// 3: (1,2) -> (1,6) [crosses 3|4 forward]
		// 4: (1,6) -> (1,3) [crosses 3|4 backward]
		// 5: (1,3) -> (1,5) [crosses 3|4 forward]
		// 6: (1,5) -> (1,2) [crosses 3|4 backward]
		sequence := []struct {
			from Position
			to   Position
		}{
			{NewPosition(1, 1), NewPosition(1, 5)},
			{NewPosition(1, 5), NewPosition(1, 2)},
			{NewPosition(1, 2), NewPosition(1, 6)},
			{NewPosition(1, 6), NewPosition(1, 3)},
			{NewPosition(1, 3), NewPosition(1, 5)},
			{NewPosition(1, 5), NewPosition(1, 2)},
		}

		for _, step := range sequence {
			move := NewMove(step.from, step.to, &player1)
			assert.False(t, IsTwoSquareViolation(history, 0, move))
			history = append(history, move)
		}

		// Turn 7: proposing (1,2) -> (1,5) crosses 3|4 forward again for the 7th time -> violation!
		violatingMove := NewMove(NewPosition(1, 2), NewPosition(1, 5), &player1)
		assert.True(t, IsTwoSquareViolation(history, 0, violatingMove))

		// Turn 7: stopping before boundary 3|4: (1,2) -> (1,3) does not cross 3|4 -> legal!
		legalMove := NewMove(NewPosition(1, 2), NewPosition(1, 3), &player1)
		assert.False(t, IsTwoSquareViolation(history, 0, legalMove))
	})
}
