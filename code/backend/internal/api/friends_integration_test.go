package api

import (
	"digital-innovation/gostrategy/internal/db"
	"digital-innovation/gostrategy/internal/models"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegration_FriendsFlow(t *testing.T) {
	_, ts := setupTestServer(t)
	defer ts.Close()
	alice := registerFriendsClient(t, ts, "flow_alice")
	bob := registerFriendsClient(t, ts, "flow_bob")
	aliceCode := alice.friendCode()

	t.Run("a friend code is stable and formatted", func(t *testing.T) {
		assert.Regexp(t, `^[0-9A-Z]{4}-[0-9A-Z]{4}$`, aliceCode)
		assert.Equal(t, aliceCode, alice.friendCode())
	})

	t.Run("preview resolves any spelling of the code", func(t *testing.T) {
		messy := strings.ToLower(strings.ReplaceAll(aliceCode, "-", " "))
		var res models.FriendPreviewResponse
		status, body := bob.call(http.MethodGet, "/friends/preview?code="+url.QueryEscape(messy), nil, &res)
		require.Equal(t, http.StatusOK, status, body)
		assert.Equal(t, "flow_alice", res.User.Username)
		assert.Equal(t, alice.id, res.User.ID)
	})

	var sent models.Friendship
	t.Run("sending a request creates a pending friendship", func(t *testing.T) {
		require.Equal(t, http.StatusCreated, bob.sendRequest(aliceCode, &sent))
		assert.Equal(t, models.FriendshipPending, sent.Status)
		assert.Equal(t, bob.id, sent.RequesterID)

		assert.Len(t, bob.list().Outgoing, 1)
		incoming := alice.list().Incoming
		require.Len(t, incoming, 1)
		assert.Equal(t, "flow_bob", incoming[0].User.Username)
	})

	t.Run("only the addressee can accept", func(t *testing.T) {
		path := fmt.Sprintf("/friends/requests/%d/accept", sent.ID)
		status, _ := bob.call(http.MethodPost, path, nil, nil)
		assert.Equal(t, http.StatusNotFound, status)

		var accepted models.Friendship
		status, body := alice.call(http.MethodPost, path, nil, &accepted)
		require.Equal(t, http.StatusOK, status, body)
		assert.Equal(t, models.FriendshipAccepted, accepted.Status)
		assert.Len(t, alice.list().Friends, 1)
		assert.Len(t, bob.list().Friends, 1)
	})

	t.Run("either side can unfriend", func(t *testing.T) {
		status, _ := bob.call(http.MethodDelete, fmt.Sprintf("/friends/%d", sent.ID), nil, nil)
		assert.Equal(t, http.StatusNoContent, status)
		assert.Empty(t, alice.list().Friends)
	})
}

func TestIntegration_FriendRequestErrors(t *testing.T) {
	_, ts := setupTestServer(t)
	defer ts.Close()
	alice := registerFriendsClient(t, ts, "err_alice")
	bob := registerFriendsClient(t, ts, "err_bob")
	aliceCode, bobCode := alice.friendCode(), bob.friendCode()

	t.Run("a malformed code is a 400", func(t *testing.T) {
		assert.Equal(t, http.StatusBadRequest, alice.sendRequest("nope", nil))
	})

	t.Run("a missing body is a 400", func(t *testing.T) {
		status, _ := alice.call(http.MethodPost, "/friends/requests", nil, nil)
		assert.Equal(t, http.StatusBadRequest, status)
	})

	t.Run("an unknown code is a 404", func(t *testing.T) {
		assert.Equal(t, http.StatusNotFound, alice.sendRequest("ZZZZ-ZZZZ", nil))
		status, _ := alice.call(http.MethodGet, "/friends/preview?code=ZZZZ-ZZZZ", nil, nil)
		assert.Equal(t, http.StatusNotFound, status)
	})

	t.Run("your own code is a 400", func(t *testing.T) {
		assert.Equal(t, http.StatusBadRequest, alice.sendRequest(aliceCode, nil))
	})

	t.Run("a repeat request is a 409", func(t *testing.T) {
		require.Equal(t, http.StatusCreated, alice.sendRequest(bobCode, nil))
		assert.Equal(t, http.StatusConflict, alice.sendRequest(bobCode, nil))
	})

	t.Run("a crossed request accepts the existing one with a 200", func(t *testing.T) {
		var res models.Friendship
		require.Equal(t, http.StatusOK, bob.sendRequest(aliceCode, &res))
		assert.Equal(t, models.FriendshipAccepted, res.Status)
	})

	t.Run("friends cannot request each other again", func(t *testing.T) {
		assert.Equal(t, http.StatusConflict, alice.sendRequest(bobCode, nil))
	})

	t.Run("ids must be numeric and known", func(t *testing.T) {
		// why: alice has spent most of her action budget above, and the limiter is per user
		carol := registerFriendsClient(t, ts, "err_carol")
		status, _ := carol.call(http.MethodPost, "/friends/requests/abc/accept", nil, nil)
		assert.Equal(t, http.StatusBadRequest, status)
		status, _ = carol.call(http.MethodPost, "/friends/requests/999999/accept", nil, nil)
		assert.Equal(t, http.StatusNotFound, status)
		status, _ = carol.call(http.MethodDelete, "/friends/999999", nil, nil)
		assert.Equal(t, http.StatusNotFound, status)
	})
}

func TestIntegration_FriendLimits(t *testing.T) {
	_, ts := setupTestServer(t)
	defer ts.Close()
	target := registerFriendsClient(t, ts, "limit_target")
	targetCode := target.friendCode()

	seed := func(userID, n int, status models.FriendshipStatus) {
		rows := make([]models.Friendship, n)
		for i := range rows {
			rows[i] = models.Friendship{RequesterID: userID, AddresseeID: 100000 + i, Status: status}
		}
		require.NoError(t, db.DB.CreateInBatches(rows, 100).Error)
	}

	t.Run("too many pending requests is a 429", func(t *testing.T) {
		sender := registerFriendsClient(t, ts, "limit_pending")
		seed(sender.id, db.MaxOutgoingPending, models.FriendshipPending)
		assert.Equal(t, http.StatusTooManyRequests, sender.sendRequest(targetCode, nil))
	})

	t.Run("a full friend list is a 409", func(t *testing.T) {
		sender := registerFriendsClient(t, ts, "limit_friends")
		seed(sender.id, db.MaxFriends, models.FriendshipAccepted)
		assert.Equal(t, http.StatusConflict, sender.sendRequest(targetCode, nil))
	})
}

func TestIntegration_FriendRoutesRequireLogin(t *testing.T) {
	_, ts := setupTestServer(t)
	defer ts.Close()
	anon := &friendsClient{t: t, baseURL: ts.URL}

	routes := []struct{ method, path string }{
		{http.MethodPost, "/users/me/friend-code"},
		{http.MethodPost, "/users/me/friend-code/regenerate"},
		{http.MethodGet, "/friends"},
		{http.MethodGet, "/friends/preview?code=AAAA-AAAA"},
		{http.MethodPost, "/friends/requests"},
		{http.MethodPost, "/friends/requests/1/accept"},
		{http.MethodDelete, "/friends/1"},
	}
	for _, r := range routes {
		status, _ := anon.call(r.method, r.path, nil, nil)
		assert.Equal(t, http.StatusUnauthorized, status, "%s %s", r.method, r.path)
	}
}

func TestIntegration_FriendCodeRegenerate(t *testing.T) {
	_, ts := setupTestServer(t)
	defer ts.Close()
	alice := registerFriendsClient(t, ts, "regen_alice")
	bob := registerFriendsClient(t, ts, "regen_bob")
	oldCode := alice.friendCode()

	var res models.FriendCodeResponse
	status, body := alice.call(http.MethodPost, "/users/me/friend-code/regenerate", nil, &res)
	require.Equal(t, http.StatusOK, status, body)
	assert.NotEqual(t, oldCode, res.Code)

	status, _ = bob.call(http.MethodGet, "/friends/preview?code="+oldCode, nil, nil)
	assert.Equal(t, http.StatusNotFound, status, "the old code must stop working")
	status, _ = bob.call(http.MethodGet, "/friends/preview?code="+res.Code, nil, nil)
	assert.Equal(t, http.StatusOK, status)
}

func TestIntegration_FriendCodeStaysPrivate(t *testing.T) {
	_, ts := setupTestServer(t)
	defer ts.Close()
	alice := registerFriendsClient(t, ts, "private_alice")
	bob := registerFriendsClient(t, ts, "private_bob")
	code := alice.friendCode()
	require.Equal(t, http.StatusCreated, bob.sendRequest(code, nil))

	paths := map[string]*friendsClient{
		fmt.Sprintf("/users/%d", alice.id): bob,
		"/users/me":                        alice,
		"/friends":                         bob,
	}
	for path, client := range paths {
		status, body := client.call(http.MethodGet, path, nil, nil)
		require.Equal(t, http.StatusOK, status, path)
		assert.NotContains(t, body, "friend_code", path)
		assert.NotContains(t, body, strings.ReplaceAll(code, "-", ""), path)
	}
}

func TestIntegration_FriendRoutesAreRateLimited(t *testing.T) {
	routes := map[string]func(c *friendsClient) (int, string){
		"preview": func(c *friendsClient) (int, string) {
			return c.call(http.MethodGet, "/friends/preview?code=nope", nil, nil)
		},
		"requests": func(c *friendsClient) (int, string) {
			return c.call(http.MethodPost, "/friends/requests", models.AddFriendRequest{Code: "nope"}, nil)
		},
		"accept": func(c *friendsClient) (int, string) {
			return c.call(http.MethodPost, "/friends/requests/999999/accept", nil, nil)
		},
		"delete": func(c *friendsClient) (int, string) {
			return c.call(http.MethodDelete, "/friends/999999", nil, nil)
		},
		"regenerate": func(c *friendsClient) (int, string) {
			return c.call(http.MethodPost, "/users/me/friend-code/regenerate", nil, nil)
		},
	}
	for name, hit := range routes {
		t.Run(name, func(t *testing.T) {
			// why: a fresh server per route keeps the per-IP limiter out of the way, so a 429
			// can only come from the per-user action limiter this test is about
			_, ts := setupTestServer(t)
			defer ts.Close()
			client := registerFriendsClient(t, ts, "rate_"+strings.ReplaceAll(name, "-", "_"))

			var limitedBody string
			for range 15 {
				if status, body := hit(client); status == http.StatusTooManyRequests {
					limitedBody = body
				}
			}
			assert.Contains(t, limitedBody, "this user account", "expected the user limiter to trip after the burst of 10")
		})
	}
}
