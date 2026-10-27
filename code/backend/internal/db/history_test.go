package db

import (
	"context"
	"digital-innovation/gostrategy/internal/models"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHistoryLogic(t *testing.T) {
	SetupDBTest(t)
	ctx := context.Background()
	user1, _ := CreateUser(ctx, "p1", "Pass1234!", "")
	user2, _ := CreateUser(ctx, "p2", "Pass1234!", "")
	ctx = WithUserID(ctx, user1.ID)

	t.Run("Save Game and Moves", func(t *testing.T) {
		gameID := "game-123"
		initialState := map[string]any{"board": "initial"}
		moves := []models.HistoricalMove{
			{MoveIndex: 1, PlayerID: 1, FromX: 0, FromY: 0, ToX: 0, ToY: 1, Result: "moved"},
			{MoveIndex: 2, PlayerID: 2, FromX: 9, FromY: 9, ToX: 9, ToY: 8, Result: "moved"},
		}

		err := SaveGameHistory(ctx, SaveHistoryParams{
			GameID:       gameID,
			Player1ID:    &user1.ID,
			Player2ID:    &user2.ID,
			WinnerID:     nil,
			GameType:     "ranked",
			InitialState: initialState,
			StartTime:    time.Now(),
			FinishedAt:   time.Now(),
			Moves:        moves,
		})
		assert.NoError(t, err)

		// Verify history
		history, err := GetGameHistory(ctx, gameID)
		assert.NoError(t, err)
		assert.Len(t, history.Moves, 2)
		assert.Equal(t, 1, history.Moves[0].MoveIndex)
	})

	t.Run("PruneGameHistory", func(t *testing.T) {
		now := time.Now()
		sampleMove := models.HistoricalMove{
			MoveIndex: 1,
			PlayerID:  1,
			FromX:     0,
			FromY:     0,
			ToX:       0,
			ToY:       1,
			Result:    "moved",
		}

		testCases := []struct {
			id        string
			gameType  string
			p1        *int
			p2        *int
			age       time.Duration
			hasReplay bool // true if replay moves should remain
		}{
			{"anon-old", models.AiVsAi, nil, nil, 48 * time.Hour, false},
			{"anon-recent", models.AiVsAi, nil, nil, 2 * time.Hour, true},
			{"user-active", "ranked", &user1.ID, &user2.ID, 48 * time.Hour, true},
			{"user-expired", "ranked", &user1.ID, &user2.ID, 40 * 24 * time.Hour, false},
		}

		for _, tc := range testCases {
			gameTime := now.Add(-tc.age)
			err := SaveGameHistory(ctx, SaveHistoryParams{
				GameID:       tc.id,
				Player1ID:    tc.p1,
				Player2ID:    tc.p2,
				GameType:     tc.gameType,
				InitialState: map[string]any{"board": "state"},
				StartTime:    gameTime,
				FinishedAt:   gameTime,
				Moves:        []models.HistoricalMove{sampleMove},
			})
			assert.NoError(t, err)
		}

		err := PruneGameHistory(ctx, 30*24*time.Hour, 24*time.Hour)
		assert.NoError(t, err)

		for _, tc := range testCases {
			// Every game metadata record must be preserved
			var gameCount int64
			DB.Model(&models.Game{}).Where("id = ?", tc.id).Count(&gameCount)
			assert.Equal(t, int64(1), gameCount, "game metadata should never be deleted for %s", tc.id)

			history, err := GetGameHistory(ctx, tc.id)
			if tc.hasReplay {
				assert.NoError(t, err)
				assert.NotNil(t, history)
				assert.Len(t, history.Moves, 1)
			} else {
				assert.Error(t, err, "GetGameHistory should fail when moves are pruned for %s", tc.id)
				assert.Nil(t, history)

				var movesCount int64
				DB.Model(&models.GameMove{}).Where("game_id = ?", tc.id).Count(&movesCount)
				assert.Zero(t, movesCount, "moves should be pruned for %s", tc.id)
			}
		}

		// Verify GetGamesForUserPaged reflects has_replay accurately
		userGames, err := GetGamesForUserPaged(ctx, user1.ID, 10, 0)
		assert.NoError(t, err)
		for _, ug := range userGames {
			if ug.ID == "user-active" {
				assert.True(t, ug.HasReplay, "user-active should have replay")
			}
			if ug.ID == "user-expired" {
				assert.False(t, ug.HasReplay, "user-expired should not have replay")
			}
		}
	})
}
