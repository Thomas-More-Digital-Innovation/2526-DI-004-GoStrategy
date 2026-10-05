package utils

import (
	"crypto/rand"
	"sync"
	"time"

	"github.com/oklog/ulid/v2"
)

var (
	ulidEntropyLock sync.Mutex
	ulidEntropy     = ulid.Monotonic(rand.Reader, 0)
)

// NewGameID generates a monotonically increasing, collision-resistant 26-character ULID
func NewGameID() string {
	ulidEntropyLock.Lock()
	defer ulidEntropyLock.Unlock()
	return ulid.MustNew(ulid.Timestamp(time.Now()), ulidEntropy).String()
}

// IsValidULID verifies whether a given string is a valid 26-character Crockford Base32 ULID
func IsValidULID(id string) bool {
	_, err := ulid.ParseStrict(id)
	return err == nil
}
