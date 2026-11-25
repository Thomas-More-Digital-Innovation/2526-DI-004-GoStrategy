package middleware

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMaskSensitiveParams(t *testing.T) {
	cases := []struct{ name, query, want string }{
		{"empty", "", ""},
		{"friend code", "code=K7QX-M2PA", "code=[REDACTED]"},
		{"friend code among others", "limit=5&Code=K7QX-M2PA", "limit=5&Code=[REDACTED]"},
		{"token", "access_token=abc", "access_token=[REDACTED]"},
		{"harmless", "limit=5&offset=10", "limit=5&offset=10"},
		{"key without value", "code", "code"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			assert.Equal(t, c.want, maskSensitiveParams(c.query))
		})
	}
}
