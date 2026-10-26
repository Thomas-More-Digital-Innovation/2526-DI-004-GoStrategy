// Package logging provides standardized logging utilities for the application
package logging

import (
	"fmt"
	"log/slog"
)

// Debug logs a debug message with a component tag
func Debug(tag string, format string, v ...any) {
	msg := fmt.Sprintf(format, v...)
	logRecord(slog.LevelDebug, 3, msg, slog.String("component", CleanTag(tag)))
}

// DebugWithUser logs a debug message with standardized user information
func DebugWithUser(tag string, username string, userID int, format string, v ...any) {
	msg := fmt.Sprintf(format, v...)
	logRecord(slog.LevelDebug, 3, msg,
		slog.String("component", CleanTag(tag)),
		slog.String("user", sanitize(username)),
		slog.Int("user_id", userID),
	)
}
