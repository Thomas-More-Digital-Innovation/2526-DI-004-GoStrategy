package db

import (
	"context"
	"digital-innovation/gostrategy/internal/models"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// newTestUser inserts a user without hashing a password, which keeps cap tests fast.
func newTestUser(t *testing.T, name string) int {
	t.Helper()
	user := models.User{Username: name, PasswordHash: "x"}
	require.NoError(t, DB.Create(&user).Error)
	return user.ID
}

// seedFriendships gives userID n outgoing rows toward users that need not exist.
func seedFriendships(t *testing.T, userID, n int, status models.FriendshipStatus) {
	t.Helper()
	rows := make([]models.Friendship, n)
	for i := range rows {
		rows[i] = models.Friendship{RequesterID: userID, AddresseeID: 100000 + i, Status: status}
	}
	require.NoError(t, DB.CreateInBatches(rows, 100).Error)
}

func pairRows(t *testing.T, a, b int) []models.Friendship {
	t.Helper()
	var rows []models.Friendship
	require.NoError(t, DB.Where("(requester_id = ? AND addressee_id = ?) OR (requester_id = ? AND addressee_id = ?)", a, b, b, a).Find(&rows).Error)
	return rows
}

func TestSendFriendRequest(t *testing.T) {
	SetupDBTest(t)
	ctx := context.Background()

	t.Run("creates a pending request", func(t *testing.T) {
		a, b := newTestUser(t, "send_a"), newTestUser(t, "send_b")
		f, err := SendFriendRequest(ctx, a, b)
		require.NoError(t, err)
		assert.Equal(t, models.FriendshipPending, f.Status)
		assert.Equal(t, a, f.RequesterID)
		assert.Equal(t, b, f.AddresseeID)
	})

	t.Run("rejects a request to yourself", func(t *testing.T) {
		a := newTestUser(t, "self")
		_, err := SendFriendRequest(ctx, a, a)
		assert.ErrorIs(t, err, ErrSelfRequest)
	})

	t.Run("rejects a duplicate in the same direction", func(t *testing.T) {
		a, b := newTestUser(t, "dup_a"), newTestUser(t, "dup_b")
		_, err := SendFriendRequest(ctx, a, b)
		require.NoError(t, err)
		_, err = SendFriendRequest(ctx, a, b)
		assert.ErrorIs(t, err, ErrAlreadyRequested)
	})

	t.Run("a crossed request accepts the existing one", func(t *testing.T) {
		a, b := newTestUser(t, "cross_a"), newTestUser(t, "cross_b")
		first, err := SendFriendRequest(ctx, a, b)
		require.NoError(t, err)

		second, err := SendFriendRequest(ctx, b, a)
		require.NoError(t, err)
		assert.Equal(t, first.ID, second.ID)
		assert.Equal(t, models.FriendshipAccepted, second.Status)
		assert.Len(t, pairRows(t, a, b), 1)
	})

	t.Run("rejects when already friends, from either side", func(t *testing.T) {
		a, b := newTestUser(t, "fr_a"), newTestUser(t, "fr_b")
		f, err := SendFriendRequest(ctx, a, b)
		require.NoError(t, err)
		_, err = AcceptFriendRequest(ctx, b, f.ID)
		require.NoError(t, err)

		_, err = SendFriendRequest(ctx, a, b)
		assert.ErrorIs(t, err, ErrAlreadyFriends)
		_, err = SendFriendRequest(ctx, b, a)
		assert.ErrorIs(t, err, ErrAlreadyFriends)
	})

	t.Run("caps outgoing pending requests", func(t *testing.T) {
		a, b := newTestUser(t, "cap_pending_a"), newTestUser(t, "cap_pending_b")
		seedFriendships(t, a, MaxOutgoingPending, models.FriendshipPending)
		_, err := SendFriendRequest(ctx, a, b)
		assert.ErrorIs(t, err, ErrPendingLimit)
	})

	t.Run("caps the requester's friends", func(t *testing.T) {
		a, b := newTestUser(t, "cap_friends_a"), newTestUser(t, "cap_friends_b")
		seedFriendships(t, a, MaxFriends, models.FriendshipAccepted)
		_, err := SendFriendRequest(ctx, a, b)
		assert.ErrorIs(t, err, ErrFriendLimit)
	})

	t.Run("a crossed accept respects the acceptor's cap", func(t *testing.T) {
		a, b := newTestUser(t, "cap_cross_a"), newTestUser(t, "cap_cross_b")
		_, err := SendFriendRequest(ctx, a, b)
		require.NoError(t, err)
		seedFriendships(t, b, MaxFriends, models.FriendshipAccepted)

		_, err = SendFriendRequest(ctx, b, a)
		assert.ErrorIs(t, err, ErrFriendLimit)
		assert.Equal(t, models.FriendshipPending, pairRows(t, a, b)[0].Status)
	})

	t.Run("allows a new request after the old one was deleted", func(t *testing.T) {
		a, b := newTestUser(t, "again_a"), newTestUser(t, "again_b")
		f, err := SendFriendRequest(ctx, a, b)
		require.NoError(t, err)
		require.NoError(t, DeleteFriendship(ctx, b, f.ID))

		_, err = SendFriendRequest(ctx, a, b)
		assert.NoError(t, err)
	})
}

func TestFriendshipPairIndex(t *testing.T) {
	SetupDBTest(t)
	a, b := newTestUser(t, "idx_a"), newTestUser(t, "idx_b")
	require.NoError(t, DB.Create(&models.Friendship{RequesterID: a, AddresseeID: b, Status: models.FriendshipPending}).Error)

	err := DB.Create(&models.Friendship{RequesterID: b, AddresseeID: a, Status: models.FriendshipPending}).Error
	assert.Error(t, err, "the reverse direction must hit the unique pair index")
}

func TestSendFriendRequestConcurrent(t *testing.T) {
	SetupDBTest(t)
	ctx := context.Background()
	a, b := newTestUser(t, "race_a"), newTestUser(t, "race_b")

	var wg sync.WaitGroup
	start := make(chan struct{})
	for _, pair := range [][2]int{{a, b}, {b, a}} {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			_, _ = SendFriendRequest(ctx, pair[0], pair[1])
		}()
	}
	close(start)
	wg.Wait()

	assert.LessOrEqual(t, len(pairRows(t, a, b)), 1, "the pair index must never allow two rows")
}

func TestAcceptFriendRequest(t *testing.T) {
	SetupDBTest(t)
	ctx := context.Background()

	pending := func(t *testing.T, name string) (a, b int, id int) {
		t.Helper()
		a, b = newTestUser(t, name+"_a"), newTestUser(t, name+"_b")
		f, err := SendFriendRequest(ctx, a, b)
		require.NoError(t, err)
		return a, b, f.ID
	}

	t.Run("the addressee accepts", func(t *testing.T) {
		_, b, id := pending(t, "acc")
		f, err := AcceptFriendRequest(ctx, b, id)
		require.NoError(t, err)
		assert.Equal(t, models.FriendshipAccepted, f.Status)
	})

	t.Run("the requester cannot accept", func(t *testing.T) {
		a, _, id := pending(t, "acc_req")
		_, err := AcceptFriendRequest(ctx, a, id)
		assert.ErrorIs(t, err, ErrFriendshipNotFound)
	})

	t.Run("a third user cannot accept", func(t *testing.T) {
		_, _, id := pending(t, "acc_third")
		c := newTestUser(t, "acc_third_c")
		_, err := AcceptFriendRequest(ctx, c, id)
		assert.ErrorIs(t, err, ErrFriendshipNotFound)
	})

	t.Run("an accepted request cannot be accepted again", func(t *testing.T) {
		_, b, id := pending(t, "acc_twice")
		_, err := AcceptFriendRequest(ctx, b, id)
		require.NoError(t, err)
		_, err = AcceptFriendRequest(ctx, b, id)
		assert.ErrorIs(t, err, ErrFriendshipNotFound)
	})

	t.Run("an unknown id is not found", func(t *testing.T) {
		_, err := AcceptFriendRequest(ctx, newTestUser(t, "acc_none"), 999999)
		assert.ErrorIs(t, err, ErrFriendshipNotFound)
	})

	t.Run("respects the acceptor's friend cap", func(t *testing.T) {
		a, b, id := pending(t, "acc_cap")
		seedFriendships(t, b, MaxFriends, models.FriendshipAccepted)
		_, err := AcceptFriendRequest(ctx, b, id)
		assert.ErrorIs(t, err, ErrFriendLimit)
		assert.Equal(t, models.FriendshipPending, pairRows(t, a, b)[0].Status)
	})
}

func TestDeleteFriendship(t *testing.T) {
	SetupDBTest(t)
	ctx := context.Background()

	t.Run("the addressee declines, the requester cancels", func(t *testing.T) {
		a, b := newTestUser(t, "del_a"), newTestUser(t, "del_b")
		f, err := SendFriendRequest(ctx, a, b)
		require.NoError(t, err)
		require.NoError(t, DeleteFriendship(ctx, b, f.ID))

		f, err = SendFriendRequest(ctx, a, b)
		require.NoError(t, err)
		require.NoError(t, DeleteFriendship(ctx, a, f.ID))
		assert.Empty(t, pairRows(t, a, b))
	})

	t.Run("either side can unfriend", func(t *testing.T) {
		a, b := newTestUser(t, "unf_a"), newTestUser(t, "unf_b")
		f, err := SendFriendRequest(ctx, a, b)
		require.NoError(t, err)
		_, err = AcceptFriendRequest(ctx, b, f.ID)
		require.NoError(t, err)
		require.NoError(t, DeleteFriendship(ctx, a, f.ID))
		assert.Empty(t, pairRows(t, a, b))
	})

	t.Run("a third user cannot delete", func(t *testing.T) {
		a, b, c := newTestUser(t, "del3_a"), newTestUser(t, "del3_b"), newTestUser(t, "del3_c")
		f, err := SendFriendRequest(ctx, a, b)
		require.NoError(t, err)

		assert.ErrorIs(t, DeleteFriendship(ctx, c, f.ID), ErrFriendshipNotFound)
		assert.Len(t, pairRows(t, a, b), 1)
	})

	t.Run("an unknown id is not found", func(t *testing.T) {
		assert.ErrorIs(t, DeleteFriendship(ctx, newTestUser(t, "del_none"), 999999), ErrFriendshipNotFound)
	})
}
