package game

import (
	"digital-innovation/gostrategy/internal/game/models"
	"errors"
	"fmt"

	// nosemgrep: go.lang.security.audit.crypto.math_random.math-random-used
	"math/rand/v2"
)

// SetupGame initializes the board with pieces for both players and prepares the game
func SetupGame(game *Game, player1Pieces, player2Pieces []*Piece) error {
	// Validate piece counts
	if len(player1Pieces) != 40 || len(player2Pieces) != 40 {
		return errors.New("each player must have exactly 40 pieces")
	}

	// Place player 1 pieces (bottom 4 rows: 6-9)
	if err := placePiecesInRows(game.Board, player1Pieces, 6, 9); err != nil {
		return err
	}

	// Place player 2 pieces (top 4 rows: 0-3)
	if err := placePiecesInRows(game.Board, player2Pieces, 0, 3); err != nil {
		return err
	}

	game.InitializePieces()

	game.Players[0].InitializePieceScore(GetPieceListStrategicValue(player1Pieces))
	game.Players[1].InitializePieceScore(GetPieceListStrategicValue(player2Pieces))

	return nil
}

// placePiecesInRows places pieces sequentially in the specified row range
func placePiecesInRows(board *Board, pieces []*Piece, startRow, endRow int) error {
	pieceIndex := 0
	for y := startRow; y <= endRow; y++ {
		for x := range 10 {
			if pieceIndex >= len(pieces) {
				return nil
			}
			pos := NewPosition(x, y)
			board.SetPieceAt(pos, pieces[pieceIndex])
			pieceIndex++
		}
	}
	return nil
}

// RandomSetup creates a random valid piece placement for a player
func RandomSetup(player *Player) []*Piece {
	pieces := GetPieceList(player)
	//nolint:gosec
	rand.Shuffle(len(pieces), func(i, j int) {
		pieces[i], pieces[j] = pieces[j], pieces[i]
	})
	return pieces
}

// CustomSetup allows manual piece placement with validation
func CustomSetup(_ *Player, piecePositions map[Position]*Piece) ([]*Piece, error) {
	if len(piecePositions) != 40 {
		return nil, errors.New("must place exactly 40 pieces")
	}

	pieces := make([]*Piece, 0, 40)
	for _, piece := range piecePositions {
		pieces = append(pieces, piece)
	}

	return pieces, nil
}

func parseCell(cell byte, isBitpacked bool) (*models.PieceType, error) {
	if isBitpacked {
		if cell&BitOccupied == 0 {
			return nil, fmt.Errorf("cell is empty, all 40 cells must be occupied")
		}
		return GetPieceTypeFromCell(cell), nil
	}
	if cell == '.' || cell == ' ' {
		return nil, fmt.Errorf("cell is empty in rank string, all 40 cells must be occupied")
	}
	pieceID, ok := GetPieceIDFromRank(cell)
	if !ok {
		return nil, fmt.Errorf("invalid piece rank '%c'", cell)
	}
	return GetPieceTypeFromID(pieceID), nil
}

func validatePieceCounts(typeCounts map[string]int) error {
	for _, pt := range pieceTypes {
		expected := pt.GetCount()
		got := typeCounts[pt.GetName()]
		if got != expected {
			return fmt.Errorf("invalid setup: incorrect count for %s (got %d, expected %d)", pt.GetName(), got, expected)
		}
	}
	return nil
}

// ParseSetup converts 40-byte binary setup data or 40-character rank strings into pieces
func ParseSetup(player *Player, data []byte) ([]*Piece, error) {
	if len(data) != 40 {
		return nil, fmt.Errorf("setup data must be 40 bytes, got %d", len(data))
	}

	isBitpacked := false
	for _, b := range data {
		if b&BitOccupied != 0 {
			isBitpacked = true
			break
		}
	}

	typeCounts := make(map[string]int)
	pieces := make([]*Piece, 0, 40)
	for i := range 40 {
		pieceType, err := parseCell(data[i], isBitpacked)
		if err != nil {
			return nil, fmt.Errorf("cell %d: %w", i, err)
		}
		if pieceType == nil {
			return nil, fmt.Errorf("unknown piece type at index %d", i)
		}
		typeCounts[pieceType.GetName()]++
		pieces = append(pieces, NewPiece(*pieceType, player))
	}

	if err := validatePieceCounts(typeCounts); err != nil {
		return nil, err
	}
	return pieces, nil
}

// QuickStart creates a game with random setups for both players
func QuickStart(controller1, controller2 PlayerController) *Game {
	game := NewGame(controller1, controller2)

	player1 := controller1.GetPlayer()
	player2 := controller2.GetPlayer()

	player1Pieces := RandomSetup(player1)
	player2Pieces := RandomSetup(player2)

	if err := SetupGame(game, player1Pieces, player2Pieces); err != nil {
		panic("Failed to setup game: " + err.Error())
	}

	return game
}
