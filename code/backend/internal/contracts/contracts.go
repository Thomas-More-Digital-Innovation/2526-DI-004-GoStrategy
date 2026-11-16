// Package contracts loads the JSON fixtures in code/contracts that Go and frontend tests share.
package contracts

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
)

// Read decodes code/contracts/<name> into v, failing the test on any error.
func Read(t testing.TB, name string, v any) {
	t.Helper()
	// why: resolve from this file, not the caller's working directory, so any package can use it
	_, self, _, _ := runtime.Caller(0)
	path := filepath.Join(filepath.Dir(self), "..", "..", "..", "contracts", name)
	data, err := os.ReadFile(path) //nolint:gosec // test-only, name is a literal fixture file
	require.NoError(t, err)
	require.NoError(t, json.Unmarshal(data, v))
}
