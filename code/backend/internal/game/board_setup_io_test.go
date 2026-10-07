package game

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func validTestRows() []string {
	return []string{
		"22FM22B5BB",
		"9718773644",
		"22342258B5",
		"6663533B4B",
	}
}

func TestNormalizeSetupRows(t *testing.T) {
	input := []string{
		"22fm22b5bb",
		"97s8773644",
	}
	normalized := NormalizeSetupRows(input)
	assert.Equal(t, "220M22B5BB", normalized[0])
	assert.Equal(t, "9718773644", normalized[1])
}

func TestValidateSetupRows(t *testing.T) {
	t.Run("valid setup with F alias", func(t *testing.T) {
		norm := NormalizeSetupRows(validTestRows())
		err := ValidateSetupRows(norm)
		assert.NoError(t, err)
	})

	t.Run("wrong row count", func(t *testing.T) {
		err := ValidateSetupRows([]string{"220M22B5BB"})
		assert.Error(t, err)
	})

	t.Run("wrong row length", func(t *testing.T) {
		rows := validTestRows()
		rows[0] = "220M"
		err := ValidateSetupRows(rows)
		assert.Error(t, err)
	})

	t.Run("invalid character", func(t *testing.T) {
		rows := validTestRows()
		rows[0] = "22XM22B5BB"
		err := ValidateSetupRows(NormalizeSetupRows(rows))
		assert.Error(t, err)
	})

	t.Run("incorrect piece count", func(t *testing.T) {
		rows := []string{
			"2222222222",
			"2222222222",
			"2222222222",
			"2222222222",
		}
		err := ValidateSetupRows(rows)
		assert.Error(t, err)
	})
}

func TestRowsAndSetupDataConversion(t *testing.T) {
	rows := NormalizeSetupRows(validTestRows())
	data := RowsToSetupData(rows)
	assert.Len(t, data, 40)

	recovered, err := SetupDataToRows(data)
	require.NoError(t, err)
	assert.Equal(t, rows, recovered)

	_, err = SetupDataToRows("too_short")
	assert.Error(t, err)
}

func TestFlipSetupRows(t *testing.T) {
	rows := []string{
		"0123456789",
		"abcdefghij",
		"klmnopqrst",
		"uvwxyzABCD",
	}
	flipped := FlipSetupRows(rows)
	require.Len(t, flipped, 4)
	assert.Equal(t, "DCBAzyxwvu", flipped[0])
	assert.Equal(t, "tsrqponmlk", flipped[1])
	assert.Equal(t, "jihgfedcba", flipped[2])
	assert.Equal(t, "9876543210", flipped[3])
}

func TestLoadBoardSetupFromFile(t *testing.T) {
	tmpDir := t.TempDir()

	t.Run("map format (like board-setups.json)", func(t *testing.T) {
		content := `{
			"honey-pot": [
				"22FM22B5BB",
				"9718773644",
				"22342258B5",
				"6663533B4B"
			]
		}`
		filePath := filepath.Join(tmpDir, "map_setup.json")
		require.NoError(t, os.WriteFile(filePath, []byte(content), 0o600))

		rows, err := LoadBoardSetupFromFile(filePath)
		require.NoError(t, err)
		assert.Equal(t, "220M22B5BB", rows[0])
	})

	t.Run("single exported object format", func(t *testing.T) {
		content := `{
			"name": "Single Setup",
			"description": "Test setup",
			"rows": [
				"220M22B5BB",
				"9718773644",
				"22342258B5",
				"6663533B4B"
			]
		}`
		filePath := filepath.Join(tmpDir, "single_setup.json")
		require.NoError(t, os.WriteFile(filePath, []byte(content), 0o600))

		rows, err := LoadBoardSetupFromFile(filePath)
		require.NoError(t, err)
		assert.Equal(t, "220M22B5BB", rows[0])
	})

	t.Run("array of rows format", func(t *testing.T) {
		content := `[
			"220M22B5BB",
			"9718773644",
			"22342258B5",
			"6663533B4B"
		]`
		filePath := filepath.Join(tmpDir, "array_setup.json")
		require.NoError(t, os.WriteFile(filePath, []byte(content), 0o600))

		rows, err := LoadBoardSetupFromFile(filePath)
		require.NoError(t, err)
		assert.Equal(t, "220M22B5BB", rows[0])
	})

	t.Run("plain text format", func(t *testing.T) {
		content := "220M22B5BB971877364422342258B56663533B4B"
		filePath := filepath.Join(tmpDir, "text_setup.txt")
		require.NoError(t, os.WriteFile(filePath, []byte(content), 0o600))

		rows, err := LoadBoardSetupFromFile(filePath)
		require.NoError(t, err)
		assert.Equal(t, "220M22B5BB", rows[0])
	})

	t.Run("nonexistent file", func(t *testing.T) {
		_, err := LoadBoardSetupFromFile(filepath.Join(tmpDir, "missing.json"))
		assert.Error(t, err)
	})
}
