package db

import (
	"context"
	"digital-innovation/gostrategy/internal/models"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var formattedFriendCode = regexp.MustCompile(`^[0-9A-HJKMNP-TV-Z]{4}-[0-9A-HJKMNP-TV-Z]{4}$`)

func TestNormalizeFriendCode(t *testing.T) {
	t.Run("accepts and canonicalises", func(t *testing.T) {
		cases := []struct{ raw, want string }{
			{"K7QXM2PA", "K7QXM2PA"},
			{"k7qx-m2pa", "K7QXM2PA"},
			{" K7QX M2PA ", "K7QXM2PA"},
			{"OOIILL00", "00111100"},
			{"oo-ii-ll-00", "00111100"},
		}
		for _, c := range cases {
			got, err := normalizeFriendCode(c.raw)
			require.NoError(t, err, c.raw)
			assert.Equal(t, c.want, got, c.raw)
		}
	})

	t.Run("rejects malformed input", func(t *testing.T) {
		for _, raw := range []string{"", "K7QXM2P", "K7QXM2PAA", "K7QXM2PU", "K7QX!M2P", "----"} {
			_, err := normalizeFriendCode(raw)
			assert.ErrorIs(t, err, ErrInvalidFriendCode, raw)
		}
	})
}

func TestGenerateFriendCode(t *testing.T) {
	t.Run("uses the alphabet and fixed length", func(t *testing.T) {
		for range 200 {
			code, err := generateFriendCode()
			require.NoError(t, err)
			require.Len(t, code, friendCodeLength)
			normalized, err := normalizeFriendCode(code)
			require.NoError(t, err)
			assert.Equal(t, code, normalized)
		}
	})
}

func TestFormatFriendCode(t *testing.T) {
	assert.Equal(t, "K7QX-M2PA", formatFriendCode("K7QXM2PA"))
}

func TestFriendCodeLifecycle(t *testing.T) {
	SetupDBTest(t)
	ctx := context.Background()

	t.Run("get or create is idempotent and formatted", func(t *testing.T) {
		id := newTestUser(t, "codeowner")

		first, err := GetOrCreateFriendCode(ctx, id)
		require.NoError(t, err)
		assert.Regexp(t, formattedFriendCode, first)

		second, err := GetOrCreateFriendCode(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, first, second)
	})

	t.Run("find resolves any spelling of the code", func(t *testing.T) {
		id := newTestUser(t, "findable")
		code, err := GetOrCreateFriendCode(ctx, id)
		require.NoError(t, err)

		user, err := FindUserByFriendCode(ctx, "  "+code[:4]+" "+code[5:]+" ")
		require.NoError(t, err)
		assert.Equal(t, id, user.ID)
	})

	t.Run("regenerate invalidates the old code", func(t *testing.T) {
		id := newTestUser(t, "regen")
		old, err := GetOrCreateFriendCode(ctx, id)
		require.NoError(t, err)

		fresh, err := RegenerateFriendCode(ctx, id)
		require.NoError(t, err)
		assert.NotEqual(t, old, fresh)

		_, err = FindUserByFriendCode(ctx, old)
		assert.ErrorIs(t, err, ErrFriendCodeNotFound)
		user, err := FindUserByFriendCode(ctx, fresh)
		require.NoError(t, err)
		assert.Equal(t, id, user.ID)
	})

	t.Run("find rejects malformed and unknown codes", func(t *testing.T) {
		_, err := FindUserByFriendCode(ctx, "nope")
		assert.ErrorIs(t, err, ErrInvalidFriendCode)
		_, err = FindUserByFriendCode(ctx, "ZZZZ-ZZZZ")
		assert.ErrorIs(t, err, ErrFriendCodeNotFound)
	})

	t.Run("soft-deleted users are not found by code", func(t *testing.T) {
		id := newTestUser(t, "ghost")
		code, err := GetOrCreateFriendCode(ctx, id)
		require.NoError(t, err)
		require.NoError(t, DB.Delete(&models.User{}, id).Error)

		_, err = FindUserByFriendCode(ctx, code)
		assert.ErrorIs(t, err, ErrFriendCodeNotFound)
	})

	t.Run("unknown user errors instead of creating a code", func(t *testing.T) {
		_, err := GetOrCreateFriendCode(ctx, 424242)
		assert.Error(t, err)
	})
}

func TestFriendCodeCollisionRetry(t *testing.T) {
	SetupDBTest(t)
	ctx := context.Background()
	holder := newTestUser(t, "holder")
	require.NoError(t, DB.Table("users").Where("id = ?", holder).Update("friend_code", "AAAAAAAA").Error)

	t.Run("retries after a collision", func(t *testing.T) {
		id := newTestUser(t, "retrier")
		codes := []string{"AAAAAAAA", "BBBBBBBB"}
		calls := 0
		gen := func() (string, error) { calls++; return codes[calls-1], nil }

		got, err := assignFriendCode(ctx, id, false, gen)
		require.NoError(t, err)
		assert.Equal(t, "BBBB-BBBB", got)
		assert.Equal(t, 2, calls)
	})

	t.Run("gives up after the attempt limit", func(t *testing.T) {
		id := newTestUser(t, "unlucky")
		calls := 0
		gen := func() (string, error) { calls++; return "AAAAAAAA", nil }

		_, err := assignFriendCode(ctx, id, false, gen)
		assert.Error(t, err)
		assert.Equal(t, maxFriendCodeAttempts, calls)
	})
}
