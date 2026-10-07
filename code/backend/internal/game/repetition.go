package game

import "slices"

// MaxTwoSquareCycles defines the maximum allowed consecutive back-and-forth cycles
// between the same two squares (or crossing the same boundary) for a single piece
// before repetition is disallowed under the Two-Square Rule.
const MaxTwoSquareCycles = 3

type squareBoundary struct {
	p1 Position
	p2 Position
}

func newSquareBoundary(a, b Position) squareBoundary {
	if a.Y < b.Y || (a.Y == b.Y && a.X < b.X) {
		return squareBoundary{p1: a, p2: b}
	}
	return squareBoundary{p1: b, p2: a}
}

func getBoundariesCrossed(from, to Position, buf *[10]squareBoundary) int {
	dx := 0
	if to.X > from.X {
		dx = 1
	} else if to.X < from.X {
		dx = -1
	}

	dy := 0
	if to.Y > from.Y {
		dy = 1
	} else if to.Y < from.Y {
		dy = -1
	}

	if (dx != 0 && dy != 0) || (dx == 0 && dy == 0) {
		return 0
	}

	cur := from
	count := 0
	for cur != to && count < len(buf) {
		next := NewPosition(cur.X+dx, cur.Y+dy)
		buf[count] = newSquareBoundary(cur, next)
		count++
		cur = next
	}
	return count
}

func boundaryDirection(from, to Position, b squareBoundary) int {
	dx := 0
	if to.X > from.X {
		dx = 1
	} else if to.X < from.X {
		dx = -1
	}
	dy := 0
	if to.Y > from.Y {
		dy = 1
	} else if to.Y < from.Y {
		dy = -1
	}

	cur := from
	for cur != to {
		next := NewPosition(cur.X+dx, cur.Y+dy)
		if cur == b.p1 && next == b.p2 {
			return 1
		}
		if cur == b.p2 && next == b.p1 {
			return -1
		}
		cur = next
	}
	return 0
}

func checkEndpointViolation(history []Move, playerID int, from, to Position) bool {
	expectedFrom := to
	expectedTo := from
	cycleSteps := 0

	for _, hist := range slices.Backward(history) {
		if hist.GetPlayer() == nil || hist.GetPlayer().GetID() != playerID {
			continue
		}

		if hist.GetFrom() != expectedFrom || hist.GetTo() != expectedTo {
			break
		}

		cycleSteps++
		expectedFrom, expectedTo = expectedTo, expectedFrom

		if cycleSteps >= MaxTwoSquareCycles*2 {
			return true
		}
	}

	return false
}

// IsTwoSquareViolation checks if proposing move (from -> to) for playerID
// would violate the Two-Square Rule. It verifies:
//   - Exact two-square endpoint oscillation (any piece).
//   - Scout boundary crossing rule: no square boundary may be crossed in alternating
//     directions more than MaxTwoSquareCycles consecutive times by the same piece.
func IsTwoSquareViolation(history []Move, playerID int, move Move) bool {
	from := move.GetFrom()
	to := move.GetTo()

	if checkEndpointViolation(history, playerID, from, to) {
		return true
	}
	var proposedBuf [10]squareBoundary
	numBoundaries := getBoundariesCrossed(from, to, &proposedBuf)
	if numBoundaries <= 1 {
		return false
	}

	for bIdx := range numBoundaries {
		b := proposedBuf[bIdx]
		propDir := boundaryDirection(from, to, b)
		if propDir == 0 {
			continue
		}

		expectedDir := -propDir
		crossings := 0
		lastPos := from

		for _, hist := range slices.Backward(history) {
			if hist.GetPlayer() == nil || hist.GetPlayer().GetID() != playerID {
				continue
			}

			if hist.GetTo() != lastPos {
				break
			}
			lastPos = hist.GetFrom()

			histDir := boundaryDirection(hist.GetFrom(), hist.GetTo(), b)
			if histDir == 0 || histDir != expectedDir {
				break
			}

			crossings++
			expectedDir = -expectedDir

			if crossings >= MaxTwoSquareCycles*2 {
				return true
			}
		}
	}

	return false
}

// IsTwoSquareViolation checks if the proposed move violates the Two-Square Rule in this game.
func (g *Game) IsTwoSquareViolation(playerID int, move Move) bool {
	if g == nil {
		return false
	}
	return IsTwoSquareViolation(g.MoveHistory, playerID, move)
}

// FilterTwoSquareMoves removes moves that violate the Two-Square Rule for playerID.
func (g *Game) FilterTwoSquareMoves(playerID int, moves []Move) []Move {
	if g == nil || len(moves) == 0 {
		return moves
	}

	filtered := make([]Move, 0, len(moves))
	for _, m := range moves {
		if !g.IsTwoSquareViolation(playerID, m) {
			filtered = append(filtered, m)
		}
	}
	return filtered
}
