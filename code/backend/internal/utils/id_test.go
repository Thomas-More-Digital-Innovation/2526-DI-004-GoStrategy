package utils

import (
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewGameID(t *testing.T) {
	t.Run("format and length", func(t *testing.T) {
		id := NewGameID()
		assert.Len(t, id, 26)
		assert.False(t, strings.HasPrefix(id, "game-"))
		assert.True(t, IsValidULID(id))
	})

	t.Run("concurrency and uniqueness", func(t *testing.T) {
		const count = 1000
		ids := make([]string, count)
		var wg sync.WaitGroup
		wg.Add(count)

		for i := 0; i < count; i++ {
			go func(idx int) {
				defer wg.Done()
				ids[idx] = NewGameID()
			}(i)
		}
		wg.Wait()

		seen := make(map[string]struct{}, count)
		for _, id := range ids {
			assert.Len(t, id, 26)
			assert.True(t, IsValidULID(id))
			_, exists := seen[id]
			assert.False(t, exists, "duplicate ULID detected: %s", id)
			seen[id] = struct{}{}
		}
	})
}

func TestIsValidULID(t *testing.T) {
	t.Run("valid ULIDs", func(t *testing.T) {
		assert.True(t, IsValidULID("01AN4Z07BY79KA1307SR9X4MV3"))
		assert.True(t, IsValidULID(NewGameID()))
	})

	t.Run("invalid ULIDs", func(t *testing.T) {
		assert.False(t, IsValidULID(""))
		assert.False(t, IsValidULID("game-12345"))
		assert.False(t, IsValidULID("not-a-valid-ulid"))
		assert.False(t, IsValidULID("01AN4Z07BY79KA1307SR9X4MV!"))
	})
}
