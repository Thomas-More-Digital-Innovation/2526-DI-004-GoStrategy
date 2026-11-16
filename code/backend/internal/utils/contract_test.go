package utils_test

import (
	"digital-innovation/gostrategy/internal/contracts"
	"digital-innovation/gostrategy/internal/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

type authContract struct {
	Username []struct {
		Input       string `json:"input"`
		FormatValid bool   `json:"formatValid"`
		Valid       bool   `json:"valid"`
	} `json:"username"`
	Password []struct {
		Input       string `json:"input"`
		FormatValid bool   `json:"formatValid"`
		Accepted    bool   `json:"accepted"`
	} `json:"password"`
}

func TestAuthContract(t *testing.T) {
	var c authContract
	contracts.Read(t, "auth.json", &c)

	t.Run("username", func(t *testing.T) {
		for _, tc := range c.Username {
			t.Run(tc.Input, func(t *testing.T) {
				assert.Equal(t, tc.FormatValid, utils.ValidateUsername(tc.Input) == nil, "format")
				assert.Equal(t, tc.Valid, utils.ValidateNewUsername(tc.Input) == nil, "register")
			})
		}
	})

	t.Run("password", func(t *testing.T) {
		for _, tc := range c.Password {
			t.Run(tc.Input, func(t *testing.T) {
				assert.Equal(t, tc.FormatValid, utils.ValidatePassword(tc.Input) == nil, "format")
				assert.Equal(t, tc.Accepted, utils.ValidateNewPassword(tc.Input) == nil, "new password")
			})
		}
	})
}
