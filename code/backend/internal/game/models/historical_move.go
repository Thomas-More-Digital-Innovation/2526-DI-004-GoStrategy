package models

// HistoricalMove represents a move in the game history with its outcome
type HistoricalMove struct {
	MoveIndex int            `json:"moveIndex"`
	PlayerID  int            `json:"playerId"`
	FromX     int            `json:"fromX"`
	FromY     int            `json:"fromY"`
	ToX       int            `json:"toX"`
	ToY       int            `json:"toY"`
	Attacker  *PieceData     `json:"attacker,omitempty"`
	Defender  *PieceData     `json:"defender,omitempty"`
	Result    MoveResultType `json:"result"`
}

// PieceData represents a piece's identity in the history
type PieceData struct {
	Type    string `json:"type"`
	Rank    string `json:"rank"`
	OwnerID int    `json:"ownerId"`
}

// MoveResultType defines the outcome of a move
type MoveResultType string

// Move result types
const (
	MoveResultTypeMove    MoveResultType = "move"    // Normal move to empty cell
	MoveResultTypeWin     MoveResultType = "win"     // Attacker won combat
	MoveResultTypeLoss    MoveResultType = "loss"    // Attacker lost combat
	MoveResultTypeTie     MoveResultType = "tie"     // Both pieces died
	MoveResultTypeCapture MoveResultType = "capture" // Flag captured (game over)
)

// GameHistory represents the full history of a game
type GameHistory struct {
	GameID       string           `json:"gameId"`
	InitialState [][]PieceData    `json:"initialState"`
	Moves        []HistoricalMove `json:"moves"`
	WinnerID     *int             `json:"winnerId"`
}
