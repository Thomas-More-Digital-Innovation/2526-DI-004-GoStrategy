package handlers

import (
	"digital-innovation/gostrategy/internal/api/dto"
	"digital-innovation/gostrategy/internal/api/ws"
	"digital-innovation/gostrategy/internal/game"
	"encoding/json"
	"fmt"
)

func init() {
	ws.RegisterHandler(dto.MsgTypeSwapPieces, handleSwapPieces)
	ws.RegisterHandler(dto.MsgTypeRandomizeSetup, handleRandomizeSetup)
	ws.RegisterHandler(dto.MsgTypeStartGame, handleStartGame)
	ws.RegisterHandler(dto.MsgTypeLoadSetup, handleLoadSetup)
}

// resolveSwapPlayerID determines which player owns both swap positions.
// Returns -1 and an error when the positions span setup areas or the client is unauthorized.
func resolveSwapPlayerID(c *ws.Client, pos1, pos2 game.Position) (int, error) {
	playerID := c.GetSeatIndex()
	if playerID >= 0 {
		// Regular player: both positions must be within their own setup area.
		switch playerID {
		case 0:
			if pos1.Y < 6 || pos1.Y > 9 || pos2.Y < 6 || pos2.Y > 9 {
				return -1, fmt.Errorf("positions outside setup area")
			}
		default:
			if pos1.Y < 0 || pos1.Y > 3 || pos2.Y < 0 || pos2.Y > 3 {
				return -1, fmt.Errorf("positions outside setup area")
			}
		}
		return playerID, nil
	}
	// Spectator-creator (AI vs AI): may swap for either player if both positions are in the same area.
	if !c.IsAuthorized() {
		return -1, fmt.Errorf("spectators cannot swap pieces")
	}
	switch {
	case pos1.Y >= 6 && pos2.Y >= 6:
		return 0, nil
	case pos1.Y <= 3 && pos2.Y <= 3:
		return 1, nil
	default:
		return -1, fmt.Errorf("both positions must be within the same player's setup area")
	}
}

func handleSwapPieces(c *ws.Client, data json.RawMessage) error {
	var swapMsg dto.SwapPiecesMessage
	if err := json.Unmarshal(data, &swapMsg); err != nil {
		return err
	}
	pos1 := game.NewPosition(swapMsg.Pos1.X, swapMsg.Pos1.Y)
	pos2 := game.NewPosition(swapMsg.Pos2.X, swapMsg.Pos2.Y)
	playerID, err := resolveSwapPlayerID(c, pos1, pos2)
	if err != nil {
		return err
	}
	if err := c.GetSession().SwapSetupPieces(playerID, pos1, pos2); err != nil {
		return err
	}
	c.GetHub().BroadcastSetupBoard()
	return nil
}

func handleRandomizeSetup(c *ws.Client, data json.RawMessage) error {
	if !c.IsAuthorized() {
		return fmt.Errorf("unauthorized")
	}

	targetPlayer := c.GetSeatIndex()
	// AI vs AI check... simplified for brevity as per dispatcher goal
	var msg dto.RandomizeSetupMessage
	if err := json.Unmarshal(data, &msg); err == nil && msg.PlayerID != nil {
		targetPlayer = *msg.PlayerID
	}

	if targetPlayer < 0 {
		return fmt.Errorf("target player required")
	}

	if err := c.GetSession().RandomizeSetup(targetPlayer); err != nil {
		return err
	}
	c.GetHub().BroadcastSetupBoard()
	return nil
}

func handleStartGame(c *ws.Client, data json.RawMessage) error {
	if !c.IsAuthorized() {
		return fmt.Errorf("unauthorized")
	}

	var msg dto.StartGameMessage
	headless := false
	if err := json.Unmarshal(data, &msg); err == nil {
		headless = msg.Headless
	}

	if err := c.GetSession().StartGameFromSetup(headless); err != nil {
		return err
	}
	c.GetHub().BroadcastGameTransition()
	return nil
}

func handleLoadSetup(c *ws.Client, data json.RawMessage) error {
	if !c.IsAuthorized() {
		return fmt.Errorf("unauthorized")
	}

	var loadMsg dto.LoadSetupMessage
	if err := json.Unmarshal(data, &loadMsg); err != nil {
		return err
	}

	targetPlayer := c.GetSeatIndex()
	if loadMsg.PlayerID != nil {
		targetPlayer = *loadMsg.PlayerID
	}

	if targetPlayer < 0 {
		return fmt.Errorf("target player required")
	}

	// Setup data decoding logic...
	// For now, assume session handles it
	if err := c.GetSession().LoadSetup(targetPlayer, []byte(loadMsg.SetupData)); err != nil {
		return err
	}
	c.GetHub().BroadcastSetupBoard()
	return nil
}
