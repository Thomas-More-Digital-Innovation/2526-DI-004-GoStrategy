package models

// GameListResponse is a page of games plus the total count for pagination.
type GameListResponse struct {
	Games []Game `json:"games"`
	Total int64  `json:"total"`
}

// CreateGameResponse is returned after a game session is created.
type CreateGameResponse struct {
	GameID   string `json:"gameId"`
	GameType string `json:"gameType"`
	WsURL    string `json:"wsUrl"`
}

// ReconnectableGameResponse tells the client whether it can rejoin a running game.
type ReconnectableGameResponse struct {
	HasGame   bool   `json:"hasGame"`
	GameID    string `json:"gameId,omitempty"`
	GameType  string `json:"gameType,omitempty"`
	SeatIndex *int   `json:"seatIndex,omitempty"`
}

// MessageResponse is a plain acknowledgement body.
type MessageResponse struct {
	Message string `json:"message"`
}
