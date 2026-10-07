package db

import (
	"context"
	"digital-innovation/gostrategy/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestBoardSetupLogic(t *testing.T) {
	SetupDBTest(t)
	ctx := context.Background()
	user, _ := CreateUser(ctx, "setupuser", "Pass1234!", "")
	ctx = WithUserID(ctx, user.ID)

	t.Run("CRUD BoardSetup", func(t *testing.T) {
		setup, err := CreateBoardSetup(ctx, user.ID, "My Setup", "Desc", "DATA", true)
		assert.NoError(t, err)

		// Read
		retrieved, err := GetBoardSetup(ctx, setup.ID, user.ID)
		assert.NoError(t, err)
		assert.Equal(t, "My Setup", retrieved.Name)

		// Update
		err = UpdateBoardSetup(ctx, setup.ID, user.ID, "Updated Name", "", "", false)
		assert.NoError(t, err)

		// Verify update
		DB.First(&setup, setup.ID)
		assert.Equal(t, "Updated Name", setup.Name)

		// Delete
		err = DeleteBoardSetup(ctx, setup.ID, user.ID)
		assert.NoError(t, err)

		// Verify soft delete
		err = DB.First(&models.BoardSetup{}, setup.ID).Error
		assert.Equal(t, gorm.ErrRecordNotFound, err)

		// Verify it still exists in DB (soft deleted)
		var deleted models.BoardSetup
		err = DB.Unscoped().First(&deleted, setup.ID).Error
		assert.NoError(t, err)
	})

	t.Run("Count BoardSetups", func(t *testing.T) {
		u, _ := CreateUser(ctx, "countuser", "Pass1234!", "")
		uCtx := WithUserID(ctx, u.ID)

		// Initial count is 0
		count, err := CountUserBoardSetups(uCtx, u.ID)
		assert.NoError(t, err)
		assert.Equal(t, int64(0), count)

		// Count increments as setups are created
		s1, err := CreateBoardSetup(uCtx, u.ID, "Setup 1", "Desc", "DATA1", false)
		assert.NoError(t, err)
		_, err = CreateBoardSetup(uCtx, u.ID, "Setup 2", "Desc", "DATA2", false)
		assert.NoError(t, err)

		count, err = CountUserBoardSetups(uCtx, u.ID)
		assert.NoError(t, err)
		assert.Equal(t, int64(2), count)

		// User isolation: other user's setups do not affect this count
		otherUser, _ := CreateUser(ctx, "otheruser", "Pass1234!", "")
		otherCtx := WithUserID(ctx, otherUser.ID)
		_, err = CreateBoardSetup(otherCtx, otherUser.ID, "Other Setup", "Desc", "DATA3", false)
		assert.NoError(t, err)

		count, err = CountUserBoardSetups(uCtx, u.ID)
		assert.NoError(t, err)
		assert.Equal(t, int64(2), count)

		// Soft deletion excludes deleted setups from count
		err = DeleteBoardSetup(uCtx, s1.ID, u.ID)
		assert.NoError(t, err)

		count, err = CountUserBoardSetups(uCtx, u.ID)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), count)

		// Context error triggers error return branch
		cancelCtx, cancel := context.WithCancel(uCtx)
		cancel()
		_, err = CountUserBoardSetups(cancelCtx, u.ID)
		assert.Error(t, err)
	})
}
