package game

import (
	"digital-innovation/gostrategy/internal/contracts"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetupContract(t *testing.T) {
	var c struct {
		PieceCounts map[string]int `json:"pieceCounts"`
		Cases       []struct {
			Rows             []string `json:"rows"`
			SetupData        string   `json:"setupData"`
			FlippedRows      []string `json:"flippedRows"`
			FlippedSetupData string   `json:"flippedSetupData"`
		} `json:"cases"`
	}
	contracts.Read(t, "setup.json", &c)

	t.Run("piece counts", func(t *testing.T) {
		require.Len(t, c.PieceCounts, len(pieceTypes))
		for rank, count := range c.PieceCounts {
			id, ok := GetPieceIDFromRank(rank[0])
			require.True(t, ok, rank)
			assert.Equal(t, count, GetPieceTypeFromID(id).GetCount(), rank)
		}
	})

	t.Run("full setup is valid", func(t *testing.T) {
		assert.NoError(t, ValidateSetup(c.Cases[0].Rows))
		assert.NoError(t, ValidateSetupRows(c.Cases[0].Rows))
	})

	for i, tc := range c.Cases {
		t.Run("case "+strconv.Itoa(i), func(t *testing.T) {
			assert.Equal(t, tc.SetupData, RowsToSetupData(tc.Rows))
			rows, err := SetupDataToRows(tc.SetupData)
			require.NoError(t, err)
			assert.Equal(t, tc.Rows, rows)
			assert.Equal(t, tc.FlippedRows, FlipSetupRows(tc.Rows))
			assert.Equal(t, tc.FlippedSetupData, RowsToSetupData(tc.FlippedRows))
		})
	}
}

func TestBoardContract(t *testing.T) {
	var c struct {
		Size      int               `json:"size"`
		Lakes     [][2]int          `json:"lakes"`
		SetupRows map[string][2]int `json:"setupRows"`
	}
	contracts.Read(t, "board.json", &c)

	t.Run("lakes", func(t *testing.T) {
		lakes := map[[2]int]bool{}
		for _, l := range c.Lakes {
			lakes[l] = true
		}
		board := NewBoard()
		for y := range c.Size {
			for x := range c.Size {
				assert.Equal(t, lakes[[2]int{x, y}], board.IsLake(NewPosition(x, y)), "x=%d y=%d", x, y)
			}
		}
	})

	t.Run("setup rows", func(t *testing.T) {
		gs := &Session{}
		for player, bounds := range c.SetupRows {
			id, err := strconv.Atoi(player)
			require.NoError(t, err)
			for y := range c.Size {
				inRange := y >= bounds[0] && y <= bounds[1]
				assert.Equal(t, inRange, gs.positionToIndex(NewPosition(0, y), id) >= 0, "player=%d y=%d", id, y)
			}
		}
	})
}
