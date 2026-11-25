package db

import (
	"context"
	"digital-innovation/gostrategy/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func usernames(entries []models.FriendEntry) []string {
	names := make([]string, len(entries))
	for i, e := range entries {
		names[i] = e.User.Username
	}
	return names
}

func TestListFriendships(t *testing.T) {
	SetupDBTest(t)
	ctx := context.Background()

	me := newTestUser(t, "me")
	friend := newTestUser(t, "friend")
	incoming := newTestUser(t, "incoming")
	outgoing := newTestUser(t, "outgoing")
	stranger := newTestUser(t, "stranger")
	ghost := newTestUser(t, "ghost_friend")

	f, err := SendFriendRequest(ctx, me, friend)
	require.NoError(t, err)
	_, err = AcceptFriendRequest(ctx, friend, f.ID)
	require.NoError(t, err)
	_, err = SendFriendRequest(ctx, incoming, me)
	require.NoError(t, err)
	_, err = SendFriendRequest(ctx, me, outgoing)
	require.NoError(t, err)
	_, err = SendFriendRequest(ctx, stranger, friend)
	require.NoError(t, err)
	g, err := SendFriendRequest(ctx, ghost, me)
	require.NoError(t, err)
	_, err = AcceptFriendRequest(ctx, me, g.ID)
	require.NoError(t, err)
	require.NoError(t, DB.Delete(&models.User{}, ghost).Error)

	list, err := ListFriendships(ctx, me)
	require.NoError(t, err)

	t.Run("buckets by status and direction", func(t *testing.T) {
		assert.Equal(t, []string{"friend"}, usernames(list.Friends))
		assert.Equal(t, []string{"incoming"}, usernames(list.Incoming))
		assert.Equal(t, []string{"outgoing"}, usernames(list.Outgoing))
	})

	t.Run("entries carry the friendship id and the other user", func(t *testing.T) {
		assert.Equal(t, f.ID, list.Friends[0].ID)
		assert.Equal(t, friend, list.Friends[0].User.ID)
		assert.False(t, list.Friends[0].CreatedAt.IsZero())
	})

	t.Run("hides friends who deleted their account", func(t *testing.T) {
		assert.NotContains(t, usernames(list.Friends), "ghost_friend")
	})

	t.Run("returns empty slices, not nil, for a user with no friends", func(t *testing.T) {
		empty, err := ListFriendships(ctx, newTestUser(t, "lonely"))
		require.NoError(t, err)
		assert.NotNil(t, empty.Friends)
		assert.NotNil(t, empty.Incoming)
		assert.NotNil(t, empty.Outgoing)
	})
}
