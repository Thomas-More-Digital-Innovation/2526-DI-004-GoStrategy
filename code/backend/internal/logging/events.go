package logging

import (
	"log/slog"
)

// UserRegistered logs a registration event
func UserRegistered(username string, userID int) {
	logRecord(slog.LevelInfo, 3, "user registered",
		slog.String("component", "auth"),
		slog.String("user", sanitize(username)),
		slog.Int("user_id", userID),
	)
}

// UserLoggedIn logs a login event
func UserLoggedIn(username string, userID int) {
	logRecord(slog.LevelInfo, 3, "user logged in",
		slog.String("component", "auth"),
		slog.String("user", sanitize(username)),
		slog.Int("user_id", userID),
	)
}

// GameStarted logs a game started event
func GameStarted(gameID, gameType string, username string, userID int) {
	logRecord(slog.LevelInfo, 3, "game started",
		slog.String("component", "game"),
		slog.String("game_id", sanitize(gameID)),
		slog.String("game_type", sanitize(gameType)),
		slog.String("user", sanitize(username)),
		slog.Int("user_id", userID),
	)
}

// GameFinished logs a game finished event
func GameFinished(gameID string, winner, loser string, rounds int) {
	logRecord(slog.LevelInfo, 3, "game finished",
		slog.String("component", "game"),
		slog.String("game_id", sanitize(gameID)),
		slog.String("winner", sanitize(winner)),
		slog.String("loser", sanitize(loser)),
		slog.Int("rounds", rounds),
	)
}

// GameAborted logs a game aborted event
func GameAborted(gameID string, reason string, username string, userID int) {
	logRecord(slog.LevelWarn, 3, "game aborted",
		slog.String("component", "game"),
		slog.String("game_id", sanitize(gameID)),
		slog.String("reason", sanitize(reason)),
		slog.String("user", sanitize(username)),
		slog.Int("user_id", userID),
	)
}

// ConnectionError logs a connection error event
func ConnectionError(gameID string, username string, userID int, err error) {
	logRecord(slog.LevelError, 3, "websocket error",
		slog.String("component", "web"),
		slog.String("game_id", sanitize(gameID)),
		slog.String("user", sanitize(username)),
		slog.Int("user_id", userID),
		slog.Any("error", err),
	)
}

// ConnectionClosed logs a connection closed event
func ConnectionClosed(gameID string, username string, userID int) {
	logRecord(slog.LevelInfo, 3, "connection closed",
		slog.String("component", "web"),
		slog.String("game_id", sanitize(gameID)),
		slog.String("user", sanitize(username)),
		slog.Int("user_id", userID),
	)
}

// Error logs a generic error event
func Error(message string, err error) {
	logRecord(slog.LevelError, 3, message, slog.Any("error", err))
}

// ErrorWithUser logs an error that includes a user in the message
func ErrorWithUser(message string, username string, userID int, err error) {
	logRecord(slog.LevelError, 3, message,
		slog.String("user", sanitize(username)),
		slog.Int("user_id", userID),
		slog.Any("error", err),
	)
}

// ErrorWith2Users logs an error that includes two users in the message
func ErrorWith2Users(message string, username1 string, userID1 int, username2 string, userID2 int, err error) {
	logRecord(slog.LevelError, 3, message,
		slog.String("user1", sanitize(username1)),
		slog.Int("user_id1", userID1),
		slog.String("user2", sanitize(username2)),
		slog.Int("user_id2", userID2),
		slog.Any("error", err),
	)
}

// ErrorWithIP logs an error that includes an IP address in the message
func ErrorWithIP(message string, ip string, err error) {
	logRecord(slog.LevelError, 3, message,
		slog.String("ip", sanitize(ip)),
		slog.Any("error", err),
	)
}

// SecurityWarning logs a security warning that includes a user in the message
func SecurityWarning(message string, details string, username string, userID int) {
	logRecord(slog.LevelWarn, 3, message,
		slog.String("component", "security"),
		slog.String("details", sanitize(details)),
		slog.String("user", sanitize(username)),
		slog.Int("user_id", userID),
	)
}

// SecurityWarningWithIP logs a security warning that includes an IP address in the message
func SecurityWarningWithIP(message string, details string, ip string) {
	logRecord(slog.LevelWarn, 3, message,
		slog.String("component", "security"),
		slog.String("details", sanitize(details)),
		slog.String("ip", sanitize(ip)),
	)
}
