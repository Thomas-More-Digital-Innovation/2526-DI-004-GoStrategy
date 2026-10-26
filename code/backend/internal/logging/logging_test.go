package logging

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSanitize(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"normal string", "normal string"},
		{"  trimmed  ", "trimmed"},
		{"string with\nnewline", "string with newline"},
		{"string with\rreturn", "string with return"},
		{"string with\x00null", "string with null"},
		{"\x1b[31mcolor codes\x1b[0m", "[31mcolor codes [0m"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			assert.Equal(t, tt.want, sanitize(tt.input))
		})
	}
}

func TestCleanTag(t *testing.T) {
	assert.Equal(t, "game", CleanTag(TagGame))
	assert.Equal(t, "auth", CleanTag(TagAuth))
	assert.Equal(t, "web", CleanTag(TagWeb))
	assert.Equal(t, "security", CleanTag(TagSecurity))
	assert.Equal(t, "custom", CleanTag("[ CUSTOM ]"))
	assert.Equal(t, "general", CleanTag(""))
}

func TestFormatUser(t *testing.T) {
	tests := []struct {
		username string
		userID   int
		want     string
	}{
		{"Alice", 1, "Alice (ID: 1)"},
		{"Bob\nSmith", 2, "Bob Smith (ID: 2)"},
		{"", 0, GuestUser},
		{"", -1, GuestUser},
	}

	for _, tt := range tests {
		t.Run(tt.username, func(t *testing.T) {
			assert.Equal(t, tt.want, FormatUser(tt.username, tt.userID))
		})
	}
}

func TestStructuredJSONOutput(t *testing.T) {
	var buf bytes.Buffer
	t.Setenv("LOG_FORMAT", "json")
	t.Setenv("LOG_LEVEL", "debug")
	initLogger(&buf)
	defer initLogger(nil)

	UserLoggedIn("testadmin", 42)

	var entry map[string]any
	err := json.Unmarshal(buf.Bytes(), &entry)
	require.NoError(t, err)

	assert.Equal(t, "INFO", entry["level"])
	assert.Equal(t, "user logged in", entry["msg"])
	assert.Equal(t, "auth", entry["component"])
	assert.Equal(t, "testadmin", entry["user"])
	assert.Equal(t, float64(42), entry["user_id"])

	// Check caller source attribution
	source, ok := entry["source"].(map[string]any)
	require.True(t, ok)
	assert.Contains(t, source["function"], "TestStructuredJSONOutput")
	assert.Equal(t, "logging_test.go", source["file"])
}

func TestLogLevelFiltering(t *testing.T) {
	var buf bytes.Buffer
	initLogger(&buf)
	defer initLogger(nil)

	SetLogLevel(slog.LevelInfo)
	Debug(TagGame, "should not appear")
	assert.Empty(t, buf.String())

	SetLogLevel(slog.LevelDebug)
	Debug(TagGame, "should appear %s", "now")
	assert.Contains(t, buf.String(), "should appear now")
	assert.Contains(t, buf.String(), "game")
}

func TestLoggingFunctions(t *testing.T) {
	var buf bytes.Buffer
	initLogger(&buf)
	defer initLogger(nil)
	SetLogLevel(slog.LevelDebug)

	err := fmt.Errorf("test error")

	UserRegistered("testuser", 1)
	UserLoggedIn("testuser", 1)

	GameStarted("game1", "pvp", "testuser", 1)
	GameFinished("game1", "winner", "loser", 10)
	GameAborted("game1", "timeout", "testuser", 1)

	ConnectionError("game1", "testuser", 1, err)
	ConnectionClosed("game1", "testuser", 1)

	Error("generic error", err)
	ErrorWithUser("user error", "testuser", 1, err)
	ErrorWith2Users("2 users error", "user1", 1, "user2", 2, err)
	ErrorWithIP("ip error", "127.0.0.1", err)

	SecurityWarning("warning", "bad origin", "testuser", 1)
	SecurityWarningWithIP("warning ip", "details", "127.0.0.1")

	Info("general info", "key1", "val1", "key2", 123)
	Warn("general warn", "key1", "val1")
	LogRaw("raw message")

	out := buf.String()
	assert.Contains(t, out, "user registered")
	assert.Contains(t, out, "game started")
	assert.Contains(t, out, "generic error")
	assert.Contains(t, out, "bad origin")
	assert.Contains(t, out, "general info")
	assert.Contains(t, out, "raw message")
}

func TestFatal(t *testing.T) {
	defer func() {
		r := recover()
		assert.NotNil(t, r, "Fatal should panic")
	}()
	Fatal("fatal message", fmt.Errorf("fatal error"))
}

func TestFatalf(t *testing.T) {
	defer func() {
		r := recover()
		assert.NotNil(t, r, "Fatalf should panic")
	}()
	Fatalf("fatal message %s", "formatted")
}

func TestDebugWithUser(t *testing.T) {
	var buf bytes.Buffer
	initLogger(&buf)
	defer initLogger(nil)
	SetLogLevel(slog.LevelDebug)

	DebugWithUser(TagWeb, "player1", 10, "state update %d", 42)
	assert.Contains(t, buf.String(), "state update 42")
	assert.Contains(t, buf.String(), "player1")
}

func TestGetLogger(t *testing.T) {
	assert.NotNil(t, GetLogger())
}
