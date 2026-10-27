// Package db provides database access and persistence for game data
package db

import (
	"context"
	"digital-innovation/gostrategy/internal/models"
	"fmt"
)

// GetGamesForUserPaged retrieves a paged list of completed games for a specific user
func GetGamesForUserPaged(ctx context.Context, userID int, limit, offset int) ([]models.Game, error) {
	var games []models.Game
	err := DB.WithContext(ctx).
		Model(&models.Game{}).
		Select("games.*, EXISTS(SELECT 1 FROM game_moves WHERE game_moves.game_id = games.id AND game_moves.deleted_at IS NULL) AS has_replay").
		Where("player1_user_id = ? OR player2_user_id = ?", userID, userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&games).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get paged games: %w", err)
	}
	return games, nil
}

// GetGamesCountForUser returns the total count of games for a specific user
func GetGamesCountForUser(ctx context.Context, userID int) (int64, error) {
	var count int64
	err := DB.WithContext(ctx).Model(&models.Game{}).
		Where("player1_user_id = ? OR player2_user_id = ?", userID, userID).
		Count(&count).Error

	if err != nil {
		return 0, fmt.Errorf("failed to count games: %w", err)
	}
	return count, nil
}
