package db

import (
	"context"
	"digital-innovation/gostrategy/internal/models"
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// SaveHistoryParams represents the parameters for saving game history
type SaveHistoryParams struct {
	GameID       string
	Player1ID    *int
	Player2ID    *int
	WinnerID     *int
	GameType     string
	InitialState any
	StartTime    time.Time
	FinishedAt   time.Time
	Moves        []models.HistoricalMove
}

// SaveGameHistory persists the game metadata and initial state
func SaveGameHistory(ctx context.Context, params SaveHistoryParams) error {
	stateJSON, err := json.Marshal(params.InitialState)
	if err != nil {
		return fmt.Errorf("failed to marshal initial state: %w", err)
	}

	game := models.Game{
		ID:            params.GameID,
		Player1UserID: params.Player1ID,
		Player2UserID: params.Player2ID,
		WinnerID:      params.WinnerID,
		GameType:      params.GameType,
		InitialState:  string(stateJSON),
		CreatedAt:     params.StartTime,
		FinishedAt:    &params.FinishedAt,
	}

	gameMoves := make([]models.GameMove, 0, len(params.Moves))
	for _, m := range params.Moves {
		gm := models.GameMove{
			GameID:    params.GameID,
			MoveIndex: m.MoveIndex,
			PlayerID:  m.PlayerID,
			FromX:     m.FromX,
			FromY:     m.FromY,
			ToX:       m.ToX,
			ToY:       m.ToY,
			Result:    m.Result,
		}
		if m.Attacker != nil {
			if b, err := json.Marshal(m.Attacker); err == nil {
				s := string(b)
				gm.AttackerData = &s
			}
		}
		if m.Defender != nil {
			if b, err := json.Marshal(m.Defender); err == nil {
				s := string(b)
				gm.DefenderData = &s
			}
		}
		gameMoves = append(gameMoves, gm)
	}

	return DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&game).Error; err != nil {
			return fmt.Errorf("failed to save game: %w", err)
		}
		if len(gameMoves) > 0 {
			if err := tx.Create(&gameMoves).Error; err != nil {
				return fmt.Errorf("failed to save moves: %w", err)
			}
		}
		return nil
	})
}

// GetGameHistory retrieves the full move history of a completed game
func GetGameHistory(ctx context.Context, gameID string) (*models.GameHistory, error) {
	var game models.Game
	err := DB.WithContext(ctx).Preload("Moves", func(db *gorm.DB) *gorm.DB {
		return db.Order("move_index ASC")
	}).Where("id = ?", gameID).First(&game).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get game history: %w", err)
	}

	var history models.GameHistory
	history.GameID = game.ID
	history.WinnerID = game.WinnerID

	if err := json.Unmarshal([]byte(game.InitialState), &history.InitialState); err != nil {
		return nil, fmt.Errorf("failed to unmarshal initial state: %w", err)
	}

	if len(game.Moves) == 0 {
		return nil, fmt.Errorf("replay moves for game %s have expired or do not exist", gameID)
	}

	for _, gm := range game.Moves {
		m := models.HistoricalMove{
			MoveIndex: gm.MoveIndex,
			PlayerID:  gm.PlayerID,
			FromX:     gm.FromX,
			FromY:     gm.FromY,
			ToX:       gm.ToX,
			ToY:       gm.ToY,
			Result:    gm.Result,
		}

		if gm.AttackerData != nil {
			var attacker models.PieceData
			if err := json.Unmarshal([]byte(*gm.AttackerData), &attacker); err == nil {
				m.Attacker = &attacker
			}
		}
		if gm.DefenderData != nil {
			var defender models.PieceData
			if err := json.Unmarshal([]byte(*gm.DefenderData), &defender); err == nil {
				m.Defender = &defender
			}
		}

		history.Moves = append(history.Moves, m)
	}

	return &history, nil
}

// PruneGameHistory removes replay moves (game_moves) older than their retention period.
// Game metadata in the games table is preserved indefinitely for stats and match history.
func PruneGameHistory(ctx context.Context, userGameTTL, anonymousGameTTL time.Duration) error {
	anonCutoff := time.Now().Add(-anonymousGameTTL)
	userCutoff := time.Now().Add(-userGameTTL)

	return DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		subQuery := tx.Model(&models.Game{}).
			Select("id").
			Where("(game_type = ? AND created_at < ?) OR (player1_user_id IS NULL AND player2_user_id IS NULL AND created_at < ?) OR (created_at < ?)",
				models.AiVsAi, anonCutoff, anonCutoff, userCutoff)

		if err := tx.Where("game_id IN (?)", subQuery).Delete(&models.GameMove{}).Error; err != nil {
			return fmt.Errorf("failed to prune game moves: %w", err)
		}

		return nil
	})
}
