package ai

import (
	"digital-innovation/gostrategy/internal/game"
	"digital-innovation/gostrategy/internal/game/models"
	"math/rand/v2"
)

// DeterminizeBoard returns a board where all unrevealed opponent pieces are randomized
// while respecting movement constraints (pieces that moved cannot be Bombs or Flags).
func DeterminizeBoard(board *game.Board, ourPlayer *game.Player, memory *Memory) *game.Board {
	opponent := GetOpponent(board, ourPlayer.GetID())
	if opponent == nil {
		return board.FastClone()
	}

	alivePieces := opponent.GetAlivePieces()
	if len(alivePieces) == 0 {
		return board.FastClone()
	}

	revealedPieces := make(map[*game.Piece]bool)
	var mobileTypes []models.PieceType
	var immobileTypes []models.PieceType

	for _, piece := range alivePieces {
		pos, exists := opponent.GetPiecePosition(piece)
		isKnown := piece.IsRevealed()
		if !isKnown && exists && memory != nil {
			if entry := memory.Recall(pos); entry != nil && entry.Confidence == 1.0 {
				isKnown = true
			}
		}

		if isKnown {
			revealedPieces[piece] = true
		} else {
			pt := *piece.GetType()
			if pt.IsMovable() {
				mobileTypes = append(mobileTypes, pt)
			} else {
				immobileTypes = append(immobileTypes, pt)
			}
		}
	}

	//nolint:gosec
	rand.Shuffle(len(mobileTypes), func(i, j int) {
		mobileTypes[i], mobileTypes[j] = mobileTypes[j], mobileTypes[i]
	})
	//nolint:gosec
	rand.Shuffle(len(immobileTypes), func(i, j int) {
		immobileTypes[i], immobileTypes[j] = immobileTypes[j], immobileTypes[i]
	})

	determinized := board.FastClone()

	var unconstrainedPositions []game.Position
	var movedPositions []game.Position

	for y := range 10 {
		for x := range 10 {
			pos := game.NewPosition(x, y)
			piece := determinized.GetPieceAt(pos)
			if piece == nil || piece.GetOwner().GetID() == ourPlayer.GetID() {
				continue
			}
			if revealedPieces[piece] {
				continue
			}

			if memory != nil && memory.HasMoved(pos) {
				movedPositions = append(movedPositions, pos)
			} else {
				unconstrainedPositions = append(unconstrainedPositions, pos)
			}
		}
	}

	//nolint:gosec
	rand.Shuffle(len(unconstrainedPositions), func(i, j int) {
		unconstrainedPositions[i], unconstrainedPositions[j] = unconstrainedPositions[j], unconstrainedPositions[i]
	})

	immobileIdx := 0
	assignedPositions := make(map[game.Position]bool)

	for _, pos := range unconstrainedPositions {
		if immobileIdx < len(immobileTypes) {
			sampledType := immobileTypes[immobileIdx]
			immobileIdx++
			newPiece := game.NewPiece(sampledType, opponent)
			determinized.SetPieceAt(pos, newPiece)
			assignedPositions[pos] = true
		}
	}

	mobileIdx := 0
	for _, pos := range unconstrainedPositions {
		if assignedPositions[pos] {
			continue
		}
		if mobileIdx < len(mobileTypes) {
			sampledType := mobileTypes[mobileIdx]
			mobileIdx++
			newPiece := game.NewPiece(sampledType, opponent)
			determinized.SetPieceAt(pos, newPiece)
		}
	}

	for _, pos := range movedPositions {
		if mobileIdx < len(mobileTypes) {
			sampledType := mobileTypes[mobileIdx]
			mobileIdx++
			newPiece := game.NewPiece(sampledType, opponent)
			determinized.SetPieceAt(pos, newPiece)
		}
	}

	return determinized
}
