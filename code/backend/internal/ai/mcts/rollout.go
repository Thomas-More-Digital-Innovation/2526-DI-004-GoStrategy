package mcts

import (
	"digital-innovation/gostrategy/internal/ai"
	ai_const "digital-innovation/gostrategy/internal/ai/const"
	"digital-innovation/gostrategy/internal/game"
	"digital-innovation/gostrategy/internal/game/models"
	"math"
	"math/rand/v2"
)

// determinizationPool caches a small set of sampled hidden-information boards
// at the root level so individual rollouts avoid re-running full determinization.
type determinizationPool struct {
	worlds []*game.Board
}

func newDeterminizationPool(board *game.Board, player *game.Player, memory *ai.Memory, size int) *determinizationPool {
	if size <= 0 {
		size = 4
	}
	worlds := make([]*game.Board, size)
	for i := range size {
		worlds[i] = ai.DeterminizeBoard(board, player, memory)
	}
	return &determinizationPool{worlds: worlds}
}

func (dp *determinizationPool) sample() *game.Board {
	if len(dp.worlds) == 0 {
		return nil
	}
	//nolint:gosec
	return dp.worlds[rand.IntN(len(dp.worlds))]
}

// runRollout executes a single simulation rollout for candidate move m.
func (aiObj *AI) runRollout(pool *determinizationPool, m game.Move, ourPlayer, opponent *game.Player) float64 {
	baseWorld := pool.sample()
	if baseWorld == nil {
		return 0.5
	}
	tempBoard := baseWorld.FastClone()

	ourFlagPos := aiObj.findFlagPosition(tempBoard, ourPlayer)
	oppFlagPos := aiObj.findFlagPosition(tempBoard, opponent)

	ourIndex := ai.BuildMobileIndex(tempBoard, ourPlayer)
	var oppIndex []game.Position
	if opponent != nil {
		oppIndex = ai.BuildMobileIndex(tempBoard, opponent)
	}

	// Apply the root candidate move first
	captured := aiObj.applySimulatedMoveInPlace(tempBoard, m)
	updateIndex(&ourIndex, m.GetFrom(), m.GetTo(), captured)

	if aiObj.isFlagCapturedAt(tempBoard, oppFlagPos, opponent) {
		return 1.0
	}
	if aiObj.isFlagCapturedAt(tempBoard, ourFlagPos, ourPlayer) {
		return 0.0
	}

	currentPlayer := opponent
	nextPlayer := ourPlayer
	maxRolloutDepth := 10

	for range maxRolloutDepth {
		if aiObj.isFlagCapturedAt(tempBoard, oppFlagPos, opponent) {
			return 1.0
		}
		if aiObj.isFlagCapturedAt(tempBoard, ourFlagPos, ourPlayer) {
			return 0.0
		}

		var currentIndex *[]game.Position
		if currentPlayer.GetID() == ourPlayer.GetID() {
			currentIndex = &ourIndex
		} else {
			currentIndex = &oppIndex
		}

		move, ok := pickFastRolloutMove(tempBoard, *currentIndex, currentPlayer)
		if !ok {
			if currentPlayer.GetID() == ourPlayer.GetID() {
				return 0.0
			}
			return 1.0
		}

		captured := aiObj.applySimulatedMoveInPlace(tempBoard, move)
		updateIndex(currentIndex, move.GetFrom(), move.GetTo(), captured)

		currentPlayer, nextPlayer = nextPlayer, currentPlayer
	}

	eval := ai.EvaluateBoard(tempBoard, ourPlayer, aiObj.GetMemory(), aiObj.params.Weights, aiObj.params.Aggression)
	score := 0.5 + math.Tanh(eval/50.0)*0.5
	if score > 1.0 {
		return 1.0
	}
	if score < 0.0 {
		return 0.0
	}
	return score
}

// pickFastRolloutMove chooses a move with 80% capture preference without full-board allocation.
func pickFastRolloutMove(board *game.Board, piecePositions []game.Position, player *game.Player) (game.Move, bool) {
	if len(piecePositions) == 0 {
		return game.Move{}, false
	}

	//nolint:gosec
	if rand.Float64() < 0.8 {
		if move, ok := pickCaptureRolloutMove(board, piecePositions, player); ok {
			return move, true
		}
	}

	return pickRandomRolloutMove(board, piecePositions, player)
}

func pickCaptureRolloutMove(board *game.Board, piecePositions []game.Position, player *game.Player) (game.Move, bool) {
	captureCount := 0
	var chosenCapture game.Move

	for _, pos := range piecePositions {
		moves, err := board.ListMoves(pos)
		if err != nil {
			continue
		}
		for _, m := range moves {
			if board.GetPieceAt(m.GetTo()) != nil {
				captureCount++
				//nolint:gosec
				if rand.IntN(captureCount) == 0 {
					chosenCapture = game.NewMove(m.GetFrom(), m.GetTo(), player)
				}
			}
		}
	}
	if captureCount > 0 {
		return chosenCapture, true
	}
	return game.Move{}, false
}

func pickRandomRolloutMove(board *game.Board, piecePositions []game.Position, player *game.Player) (game.Move, bool) {
	//nolint:gosec
	perm := rand.Perm(len(piecePositions))
	for _, idx := range perm {
		moves, err := board.ListMoves(piecePositions[idx])
		if err != nil || len(moves) == 0 {
			continue
		}
		//nolint:gosec
		chosen := moves[rand.IntN(len(moves))]
		return game.NewMove(chosen.GetFrom(), chosen.GetTo(), player), true
	}

	return game.Move{}, false
}

func updateIndex(index *[]game.Position, from, to game.Position, captured bool) {
	s := *index
	for i, pos := range s {
		if pos == from {
			if captured {
				s[i] = s[len(s)-1]
				*index = s[:len(s)-1]
			} else {
				s[i] = to
			}
			return
		}
	}
}

func (aiObj *AI) findFlagPosition(board *game.Board, player *game.Player) game.Position {
	if player == nil {
		return game.NewPosition(-1, -1)
	}
	for y := range 10 {
		for x := range 10 {
			pos := game.NewPosition(x, y)
			piece := board.GetPieceAt(pos)
			if piece != nil && piece.GetType().GetName() == ai_const.Flag && piece.GetOwner().GetID() == player.GetID() {
				return pos
			}
		}
	}
	return game.NewPosition(-1, -1)
}

func (aiObj *AI) isFlagCapturedAt(board *game.Board, flagPos game.Position, player *game.Player) bool {
	if flagPos.X == -1 {
		return true
	}
	piece := board.GetPieceAt(flagPos)
	return piece == nil || piece.GetType().GetName() != ai_const.Flag || piece.GetOwner().GetID() != player.GetID()
}

func (aiObj *AI) applySimulatedMoveInPlace(b *game.Board, move game.Move) bool {
	attacker := b.GetPieceAt(move.GetFrom())
	if attacker == nil {
		return false
	}

	target := b.GetPieceAt(move.GetTo())
	if target == nil {
		b.SetPieceAt(move.GetFrom(), nil)
		b.SetPieceAt(move.GetTo(), attacker)
		return false
	}

	attackerRank := attacker.GetRank()
	defenderRank := target.GetRank()

	if defenderRank == models.Flag.GetRank() {
		b.SetPieceAt(move.GetFrom(), nil)
		b.SetPieceAt(move.GetTo(), attacker)
		return false
	}

	if attackerRank == models.Spy.GetRank() && defenderRank == models.Marshal.GetRank() {
		b.SetPieceAt(move.GetFrom(), nil)
		b.SetPieceAt(move.GetTo(), attacker)
		return false
	}

	if defenderRank == models.Bomb.GetRank() {
		if attacker.GetType().GetName() == ai_const.Miner {
			b.SetPieceAt(move.GetFrom(), nil)
			b.SetPieceAt(move.GetTo(), attacker)
			return false
		}
		b.SetPieceAt(move.GetFrom(), nil)
		return true
	}

	switch {
	case attackerRank > defenderRank:
		b.SetPieceAt(move.GetFrom(), nil)
		b.SetPieceAt(move.GetTo(), attacker)
		return false
	case attackerRank < defenderRank:
		b.SetPieceAt(move.GetFrom(), nil)
		return true
	default:
		b.SetPieceAt(move.GetFrom(), nil)
		b.SetPieceAt(move.GetTo(), nil)
		return true
	}
}
