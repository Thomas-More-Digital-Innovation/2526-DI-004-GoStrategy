package game

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ExportedBoardSetup represents a single board setup in standard export format.
type ExportedBoardSetup struct {
	Name        string   `json:"name"`
	Description string   `json:"description,omitempty"`
	SetupData   string   `json:"setup_data"`
	Rows        []string `json:"rows"`
}

// NormalizeSetupString maps common notation aliases ('F'->'0', 'S'->'1') and normalizes chars.
func NormalizeSetupString(raw string) string {
	var sb strings.Builder
	sb.Grow(len(raw))
	for _, r := range raw {
		switch r {
		case 'f', 'F':
			sb.WriteByte('0')
		case 's', 'S':
			sb.WriteByte('1')
		case 'm':
			sb.WriteByte('M')
		case 'b':
			sb.WriteByte('B')
		default:
			sb.WriteRune(r)
		}
	}
	return sb.String()
}

// NormalizeSetupRows normalizes aliases and whitespace for each row.
func NormalizeSetupRows(rows []string) []string {
	normalized := make([]string, len(rows))
	for i, r := range rows {
		normalized[i] = NormalizeSetupString(strings.TrimSpace(r))
	}
	return normalized
}

// ValidateSetupRows checks if rows form a valid 4x10 setup with correct piece counts.
func ValidateSetupRows(rows []string) error {
	if len(rows) != PlayerSetupRows {
		return fmt.Errorf("expected %d rows, got %d", PlayerSetupRows, len(rows))
	}

	typeCounts := make(map[string]int)
	expectedCounts := make(map[string]int)
	for _, pt := range pieceTypes {
		expectedCounts[pt.GetName()] = pt.GetCount()
	}

	for y, row := range rows {
		if len(row) != PlayerSetupCols {
			return fmt.Errorf("row %d must be %d characters, got %d", y, PlayerSetupCols, len(row))
		}
		for x := 0; x < len(row); x++ {
			c := row[x]
			pieceID, ok := GetPieceIDFromRank(c)
			if !ok {
				return fmt.Errorf("invalid piece character '%c' at row %d col %d", c, y, x)
			}
			pieceType := GetPieceTypeFromID(pieceID)
			if pieceType == nil {
				return fmt.Errorf("unknown piece type for character '%c'", c)
			}
			typeCounts[pieceType.GetName()]++
		}
	}

	for typeName, expected := range expectedCounts {
		got := typeCounts[typeName]
		if got != expected {
			return fmt.Errorf("invalid piece count for %s: got %d, expected %d", typeName, got, expected)
		}
	}

	return nil
}

// RowsToSetupData joins 4 rows into a 40-character setup string.
func RowsToSetupData(rows []string) string {
	return strings.Join(rows, "")
}

// SetupDataToRows splits a 40-character setup string into 4 rows of 10.
func SetupDataToRows(setupData string) ([]string, error) {
	if len(setupData) != PlayerSetupCells {
		return nil, fmt.Errorf("setup data must be %d characters, got %d", PlayerSetupCells, len(setupData))
	}
	rows := make([]string, PlayerSetupRows)
	for i := range PlayerSetupRows {
		rows[i] = setupData[i*PlayerSetupCols : (i+1)*PlayerSetupCols]
	}
	return rows, nil
}

// FlipSetupRows rotates a 4x10 setup 180 degrees for opposite player placement.
func FlipSetupRows(rows []string) []string {
	flipped := make([]string, len(rows))
	for i := range rows {
		source := rows[len(rows)-1-i]
		runes := []rune(source)
		for j, k := 0, len(runes)-1; j < k; j, k = j+1, k-1 {
			runes[j], runes[k] = runes[k], runes[j]
		}
		flipped[i] = string(runes)
	}
	return flipped
}

// LoadBoardSetupFromFile loads, parses, and validates a setup from a JSON or text file.
func LoadBoardSetupFromFile(path string) ([]string, error) {
	//nolint:gosec
	data, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return nil, fmt.Errorf("failed to read setup file: %w", err)
	}

	rows, err := parseSetupBytes(data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse setup from %s: %w", path, err)
	}

	normRows := NormalizeSetupRows(rows)
	if err := ValidateSetupRows(normRows); err != nil {
		return nil, fmt.Errorf("invalid setup in %s: %w", path, err)
	}

	return normRows, nil
}

func parseSetupBytes(data []byte) ([]string, error) {
	var single struct {
		Name      string   `json:"name"`
		SetupData string   `json:"setup_data"`
		Rows      []string `json:"rows"`
		Setup     []string `json:"setup"`
	}
	if err := json.Unmarshal(data, &single); err == nil {
		if len(single.Rows) == PlayerSetupRows {
			return single.Rows, nil
		}
		if len(single.Setup) == PlayerSetupRows {
			return single.Setup, nil
		}
		if len(single.SetupData) == PlayerSetupCells {
			return SetupDataToRows(single.SetupData)
		}
	}

	var arr []string
	if err := json.Unmarshal(data, &arr); err == nil && len(arr) == PlayerSetupRows {
		return arr, nil
	}

	var rawMap map[string]any
	if err := json.Unmarshal(data, &rawMap); err == nil && len(rawMap) > 0 {
		for _, val := range rawMap {
			switch v := val.(type) {
			case []any:
				rows := make([]string, 0, len(v))
				for _, item := range v {
					if str, ok := item.(string); ok {
						rows = append(rows, str)
					}
				}
				if len(rows) == PlayerSetupRows {
					return rows, nil
				}
			case []string:
				if len(v) == PlayerSetupRows {
					return v, nil
				}
			case string:
				return ParseBoardSetupSmart(v)
			}
		}
	}

	return ParseBoardSetupSmart(string(data))
}
