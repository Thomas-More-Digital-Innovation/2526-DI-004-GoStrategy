package game

import (
	"digital-innovation/gostrategy/internal/game/models"
)

// MakeMove makes a move on the game board and resolves any combat that may occur.
// Returns a slice of [attacker, defender] pieces after resolution.
func (g *Game) MakeMove(move *Move, piece *Piece) []*Piece {
	target := g.Board.GetPieceAt(move.GetTo())
	piece, target = g.resolveCombat(move, piece, target)

	histMove := g.buildHistoricalMove(move, piece)
	g.MoveHistory = append(g.MoveHistory, *move)
	g.HistoricalHistory = append(g.HistoricalHistory, histMove)

	g.notifyObservers(move)
	g.NextTurn()
	return []*Piece{piece, target}
}

func (g *Game) resolveCombat(move *Move, piece, target *Piece) (*Piece, *Piece) {
	if target == nil {
		g.LastCombat = nil
		g.Board.MovePiece(move, piece)
		piece.GetOwner().UpdatePiecePosition(piece, move.GetTo())
		return piece, nil
	}

	piece.Reveal()
	target.Reveal()

	g.LastCombat = &CombatResult{
		Occurred:         true,
		AttackerPiece:    piece,
		DefenderPiece:    target,
		AttackerPosition: move.GetFrom(),
		DefenderPosition: move.GetTo(),
	}

	result := piece.Attack(target)
	piece, target = result[0], result[1]

	if piece.IsAlive() {
		g.Board.MovePiece(move, piece)
		piece.GetOwner().UpdatePiecePosition(piece, move.GetTo())
		return piece, target
	}

	g.mustRemovePieceAt(move.GetFrom())
	if target != nil && !target.IsAlive() {
		g.mustRemovePieceAt(move.GetTo())
	}
	return piece, target
}

func (g *Game) mustRemovePieceAt(pos Position) {
	if err := g.Board.RemovePieceAt(pos); err != nil {
		panic(err)
	}
}

func (g *Game) buildHistoricalMove(move *Move, piece *Piece) models.HistoricalMove {
	histMove := models.HistoricalMove{
		MoveIndex: len(g.HistoricalHistory),
		PlayerID:  move.GetPlayer().GetID(),
		FromX:     move.GetFrom().X,
		FromY:     move.GetFrom().Y,
		ToX:       move.GetTo().X,
		ToY:       move.GetTo().Y,
		Result:    models.MoveResultTypeMove,
		Attacker: &models.PieceData{
			Type:    piece.GetType().GetName(),
			Rank:    string(piece.GetType().GetRank()),
			OwnerID: piece.GetOwner().GetID(),
		},
	}

	if g.LastCombat == nil || !g.LastCombat.Occurred {
		return histMove
	}

	attackerType := g.LastCombat.AttackerPiece.GetType()
	defenderType := g.LastCombat.DefenderPiece.GetType()

	histMove.Attacker = &models.PieceData{
		Type:    attackerType.GetName(),
		Rank:    string(attackerType.GetRank()),
		OwnerID: g.LastCombat.AttackerPiece.GetOwner().GetID(),
	}
	histMove.Defender = &models.PieceData{
		Type:    defenderType.GetName(),
		Rank:    string(defenderType.GetRank()),
		OwnerID: g.LastCombat.DefenderPiece.GetOwner().GetID(),
	}
	histMove.Result = g.combatMoveResult()
	return histMove
}

func (g *Game) combatMoveResult() models.MoveResultType {
	if !g.LastCombat.AttackerPiece.IsAlive() {
		if !g.LastCombat.DefenderPiece.IsAlive() {
			return models.MoveResultTypeTie
		}
		return models.MoveResultTypeLoss
	}
	if g.LastCombat.DefenderPiece.GetType().GetName() == "Flag" {
		return models.MoveResultTypeCapture
	}
	return models.MoveResultTypeWin
}

func (g *Game) notifyObservers(move *Move) {
	round := g.GetRound()
	for _, ctrl := range g.PlayerControllers {
		if ctrl.GetPlayer().GetID() == move.GetPlayer().GetID() {
			if recorder, ok := ctrl.(interface {
				RecordOwnMove(Move)
			}); ok {
				recorder.RecordOwnMove(*move)
			}
			continue
		}

		if analyzer, ok := ctrl.(interface {
			AnalyzeMove(Move, *Player, int)
		}); ok {
			analyzer.AnalyzeMove(*move, move.GetPlayer(), round)
		}

		if g.LastCombat != nil && g.LastCombat.Occurred {
			if observer, ok := ctrl.(interface {
				ObserveCombat(Position, Position, *Piece, *Piece, int)
			}); ok {
				observer.ObserveCombat(
					g.LastCombat.AttackerPosition,
					g.LastCombat.DefenderPosition,
					g.LastCombat.AttackerPiece,
					g.LastCombat.DefenderPiece,
					round,
				)
			}
		}
	}
}

// GetInitialBoardState returns the full board state as PieceData (for history)
func (g *Game) GetInitialBoardState() [][]models.PieceData {
	if g.InitialState != nil {
		return g.InitialState
	}
	field := g.Board.GetField()
	state := make([][]models.PieceData, 10)
	for y := range 10 {
		state[y] = make([]models.PieceData, 10)
		for x := range 10 {
			piece := field[y][x]
			if piece != nil {
				state[y][x] = models.PieceData{
					Type:    piece.GetType().GetName(),
					Rank:    string(piece.GetType().GetRank()),
					OwnerID: piece.GetOwner().GetID(),
				}
			} else {
				state[y][x] = models.PieceData{
					OwnerID: -1,
				}
			}
		}
	}

	return state
}

// GetLastCombat returns the last combat result if any
func (g *Game) GetLastCombat() *CombatResult {
	return g.LastCombat
}

// ClearLastCombat clears the last combat result (called after broadcast)
func (g *Game) ClearLastCombat() {
	g.LastCombat = nil
}

// HideCombatPieces hides the pieces involved in the last combat
func (g *Game) HideCombatPieces() {
	if g.LastCombat != nil && g.LastCombat.Occurred {
		if g.LastCombat.AttackerPiece != nil && g.LastCombat.AttackerPiece.IsAlive() {
			g.LastCombat.AttackerPiece.Hide()
		}
		if g.LastCombat.DefenderPiece != nil && g.LastCombat.DefenderPiece.IsAlive() {
			g.LastCombat.DefenderPiece.Hide()
		}
	}
}

// HideAllRevealedPieces hides all revealed pieces on the board
// Called at the start of each new round to reset piece visibility
func (g *Game) HideAllRevealedPieces() {
	field := g.Board.GetField()
	for y := range 10 {
		for x := range 10 {
			piece := field[y][x]
			if piece != nil && piece.IsAlive() && piece.IsRevealed() {
				piece.Hide()
			}
		}
	}
}

// InitializePieces scans board and tracks all pieces for both players (call once at game start)
func (g *Game) InitializePieces() {
	field := g.Board.GetField()
	for y := range 10 {
		for x := range 10 {
			piece := field[y][x]
			if piece != nil {
				pos := NewPosition(x, y)
				piece.GetOwner().AddPiece(piece, pos)
			}
		}
	}
}
