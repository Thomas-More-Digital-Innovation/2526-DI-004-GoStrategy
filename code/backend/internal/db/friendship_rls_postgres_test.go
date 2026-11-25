package db

import (
	"context"
	"digital-innovation/gostrategy/internal/models"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// removeFriendTestData deletes the users' friendships and accounts. FORCE RLS hides friendship
// rows from a plain connection, so it impersonates each user. Call it while DB is still set.
func removeFriendTestData(ctx context.Context, ids ...int) {
	for _, id := range ids {
		_ = WithRLS(WithUserID(ctx, id), func(tx *gorm.DB) error {
			return tx.Where("requester_id = ? OR addressee_id = ?", id, id).Delete(&models.Friendship{}).Error
		})
	}
	DB.Unscoped().Delete(&models.User{}, ids)
	DB.Exec("DELETE FROM user_stats WHERE user_id IN ?", ids)
}

func TestFriendshipRowLevelSecurity(t *testing.T) {
	testDB := setupPostgresDB(t)
	oldDB := DB
	DB = testDB
	defer func() { DB = oldDB }()

	ctx := context.Background()
	a, errA := CreateUser(ctx, "rls_friend_a", "Pass1234!", "")
	b, errB := CreateUser(ctx, "rls_friend_b", "Pass1234!", "")
	c, errC := CreateUser(ctx, "rls_friend_c", "Pass1234!", "")
	require.NoError(t, errA)
	require.NoError(t, errB)
	require.NoError(t, errC)
	defer removeFriendTestData(ctx, a.ID, b.ID, c.ID)

	f, err := SendFriendRequest(ctx, a.ID, b.ID)
	require.NoError(t, err)

	asUser := func(id int, fn func(tx *gorm.DB) error) error {
		return WithRLS(WithUserID(ctx, id), fn)
	}

	t.Run("a third user cannot see the row", func(t *testing.T) {
		err := asUser(c.ID, func(tx *gorm.DB) error {
			return tx.First(&models.Friendship{}, f.ID).Error
		})
		assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
	})

	t.Run("a third user cannot update or delete it", func(t *testing.T) {
		var updated, deleted int64
		require.NoError(t, asUser(c.ID, func(tx *gorm.DB) error {
			updated = tx.Model(&models.Friendship{}).Where("id = ?", f.ID).Update("status", models.FriendshipAccepted).RowsAffected
			deleted = tx.Where("id = ?", f.ID).Delete(&models.Friendship{}).RowsAffected
			return nil
		}))
		assert.Zero(t, updated)
		assert.Zero(t, deleted)
	})

	t.Run("the requester cannot accept their own request", func(t *testing.T) {
		var updated int64
		require.NoError(t, asUser(a.ID, func(tx *gorm.DB) error {
			updated = tx.Model(&models.Friendship{}).Where("id = ?", f.ID).Update("status", models.FriendshipAccepted).RowsAffected
			return nil
		}))
		assert.Zero(t, updated, "the update policy admits only the addressee")
	})

	t.Run("you cannot insert a request as someone else", func(t *testing.T) {
		err := asUser(a.ID, func(tx *gorm.DB) error {
			return tx.Create(&models.Friendship{RequesterID: b.ID, AddresseeID: c.ID, Status: models.FriendshipPending}).Error
		})
		assert.Error(t, err)
	})

	t.Run("both parties see the row", func(t *testing.T) {
		for _, id := range []int{a.ID, b.ID} {
			err := asUser(id, func(tx *gorm.DB) error {
				return tx.First(&models.Friendship{}, f.ID).Error
			})
			assert.NoError(t, err, "user %d", id)
		}
	})

	t.Run("the addressee can accept", func(t *testing.T) {
		accepted, err := AcceptFriendRequest(ctx, b.ID, f.ID)
		require.NoError(t, err)
		assert.Equal(t, models.FriendshipAccepted, accepted.Status)
	})

	t.Run("either side can delete", func(t *testing.T) {
		assert.NoError(t, DeleteFriendship(ctx, a.ID, f.ID))
	})
}

func TestSendFriendRequestConcurrentPostgres(t *testing.T) {
	testDB := setupPostgresDB(t)
	oldDB := DB
	DB = testDB
	defer func() { DB = oldDB }()

	ctx := context.Background()
	a, errA := CreateUser(ctx, "race_pg_a", "Pass1234!", "")
	b, errB := CreateUser(ctx, "race_pg_b", "Pass1234!", "")
	require.NoError(t, errA)
	require.NoError(t, errB)
	defer removeFriendTestData(ctx, a.ID, b.ID)

	var wg sync.WaitGroup
	start := make(chan struct{})
	for _, pair := range [][2]int{{a.ID, b.ID}, {b.ID, a.ID}} {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			_, _ = SendFriendRequest(ctx, pair[0], pair[1])
		}()
	}
	close(start)
	wg.Wait()

	rows := 0
	require.NoError(t, WithRLS(WithUserID(ctx, a.ID), func(tx *gorm.DB) error {
		var found []models.Friendship
		if err := tx.Find(&found).Error; err != nil {
			return err
		}
		rows = len(found)
		return nil
	}))
	assert.Equal(t, 1, rows, "crossed sends must settle on exactly one row")
}

func TestListFriendshipsPostgres(t *testing.T) {
	testDB := setupPostgresDB(t)
	oldDB := DB
	DB = testDB
	defer func() { DB = oldDB }()

	ctx := context.Background()
	a, errA := CreateUser(ctx, "list_pg_a", "Pass1234!", "")
	b, errB := CreateUser(ctx, "list_pg_b", "Pass1234!", "")
	c, errC := CreateUser(ctx, "list_pg_c", "Pass1234!", "")
	require.NoError(t, errA)
	require.NoError(t, errB)
	require.NoError(t, errC)
	defer removeFriendTestData(ctx, a.ID, b.ID, c.ID)

	_, err := SendFriendRequest(ctx, a.ID, b.ID)
	require.NoError(t, err)
	_, err = SendFriendRequest(ctx, c.ID, a.ID)
	require.NoError(t, err)

	list, err := ListFriendships(ctx, a.ID)
	require.NoError(t, err)
	assert.Equal(t, []string{"list_pg_b"}, usernames(list.Outgoing))
	assert.Equal(t, []string{"list_pg_c"}, usernames(list.Incoming))
	assert.Empty(t, list.Friends)
}

func TestFriendCapOfOtherPartyPostgres(t *testing.T) {
	testDB := setupPostgresDB(t)
	oldDB := DB
	DB = testDB
	defer func() { DB = oldDB }()

	ctx := context.Background()
	full, errF := CreateUser(ctx, "cap_pg_full", "Pass1234!", "")
	other, errO := CreateUser(ctx, "cap_pg_other", "Pass1234!", "")
	require.NoError(t, errF)
	require.NoError(t, errO)
	defer removeFriendTestData(ctx, full.ID, other.ID)

	rows := make([]models.Friendship, MaxFriends)
	for i := range rows {
		rows[i] = models.Friendship{RequesterID: full.ID, AddresseeID: 100000 + i, Status: models.FriendshipAccepted}
	}
	require.NoError(t, WithRLS(WithUserID(ctx, full.ID), func(tx *gorm.DB) error {
		return tx.CreateInBatches(rows, 100).Error
	}))

	_, err := SendFriendRequest(ctx, other.ID, full.ID)
	assert.ErrorIs(t, err, ErrFriendLimit, "the other party's rows are invisible to the sender, so their cap needs their own view")
}
