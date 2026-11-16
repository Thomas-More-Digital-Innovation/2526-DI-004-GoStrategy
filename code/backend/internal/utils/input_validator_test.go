package utils_test

import (
	"digital-innovation/gostrategy/internal/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

type validatorTest struct {
	input       string
	expectError bool
}

func TestValidateGeneric(t *testing.T) {
	tests := []validatorTest{
		{"user123", false},
		{"user_name", false},
		{"user!name", false},
		{"<script>", true},
		{"   ", false},
		{"Name with ;", true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			err := utils.ValidateGeneric(tt.input, "Description")
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
