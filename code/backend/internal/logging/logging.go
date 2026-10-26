package logging

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"regexp"
	"runtime"
	"strings"
	"time"
)

var (
	nonPrintableRegex = regexp.MustCompile(`[[:cntrl:]]`)
	ansiRegex         = regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)
	logLevelVar       = new(slog.LevelVar)
	logger            *slog.Logger
)

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[41m"
	colorYellow = "\033[43m"
	colorBlue   = "\033[44m"
	colorPurple = "\033[45m"
	colorGreen  = "\033[42m"
	colorBlack  = "\033[30m"
	colorWhite  = "\033[37m"
	colorBold   = "\033[1m"
)

// Legacy log tags maintained for caller compatibility
const (
	TagDebug    = colorBold + colorBlack + colorGreen + "[ DEBUG ]" + colorReset
	TagAuth     = colorBold + colorWhite + colorBlue + "[  AUTH  ]" + colorReset
	TagGame     = colorBold + colorWhite + colorBlue + "[  GAME  ]" + colorReset
	TagWeb      = colorBold + colorWhite + colorBlue + "[  WEB   ]" + colorReset
	TagSecurity = colorBold + colorBlack + colorYellow + "[SECURITY]" + colorReset
	TagError    = colorBold + colorWhite + colorRed + "[ ERROR  ]" + colorReset
	TagFatal    = colorBold + colorWhite + colorPurple + "[ FATAL  ]" + colorReset
)

// GuestUser is the default username for unauthenticated users
const GuestUser = "Guest"

func init() {
	initLogger(os.Stdout)
}

func initLogger(w io.Writer) {
	level := slog.LevelInfo
	envLevel := strings.ToLower(os.Getenv("LOG_LEVEL"))
	debugEnv := strings.ToLower(os.Getenv("GOSTRATEGY_DEBUG"))

	switch {
	case envLevel == "debug" || debugEnv == "true":
		level = slog.LevelDebug
	case envLevel == "warn":
		level = slog.LevelWarn
	case envLevel == "error":
		level = slog.LevelError
	}
	logLevelVar.Set(level)

	opts := &slog.HandlerOptions{
		Level:     logLevelVar,
		AddSource: true,
		ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
			if a.Key == slog.SourceKey {
				if src, ok := a.Value.Any().(*slog.Source); ok && src != nil {
					idx := strings.LastIndex(src.File, "/")
					if idx != -1 {
						src.File = src.File[idx+1:]
					}
				}
			}
			return a
		},
	}

	var handler slog.Handler
	if strings.ToLower(os.Getenv("LOG_FORMAT")) == "json" || strings.ToLower(os.Getenv("GIN_MODE")) == "release" {
		handler = slog.NewJSONHandler(w, opts)
	} else {
		handler = slog.NewTextHandler(w, opts)
	}

	logger = slog.New(handler)
	slog.SetDefault(logger)
}

// sanitize strips control characters and CRLF to prevent log injection
func sanitize(s string) string {
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.ReplaceAll(s, "\r", " ")
	s = nonPrintableRegex.ReplaceAllString(s, " ")
	return strings.TrimSpace(s)
}

// CleanTag extracts a clean lowercase component string from legacy tags
func CleanTag(tag string) string {
	tag = ansiRegex.ReplaceAllString(tag, "")
	s := sanitize(tag)
	s = strings.Trim(s, "[] ")
	if s == "" {
		return "general"
	}
	return strings.ToLower(s)
}

// logRecord dispatches a structured record at the true caller depth
func logRecord(level slog.Level, depth int, msg string, attrs ...slog.Attr) {
	ctx := context.Background()
	if !logger.Enabled(ctx, level) {
		return
	}
	var pcs [1]uintptr
	runtime.Callers(depth, pcs[:])
	r := slog.NewRecord(time.Now(), level, sanitize(msg), pcs[0])
	r.AddAttrs(attrs...)
	_ = logger.Handler().Handle(ctx, r)
}

// FormatUser returns a standardized user identifier for logs
func FormatUser(username string, userID int) string {
	if username == "" && userID <= 0 {
		return GuestUser
	}
	return fmt.Sprintf("%s (ID: %d)", sanitize(username), userID)
}

// GetLogger returns the active structured slog logger
func GetLogger() *slog.Logger {
	return logger
}

// SetLogLevel sets the active logging level dynamically
func SetLogLevel(level slog.Level) {
	logLevelVar.Set(level)
}

// Info logs a general structured informational message
func Info(msg string, args ...any) {
	logRecord(slog.LevelInfo, 3, msg, argsToAttrs(args)...)
}

// Warn logs a general structured warning message
func Warn(msg string, args ...any) {
	logRecord(slog.LevelWarn, 3, msg, argsToAttrs(args)...)
}

func argsToAttrs(args []any) []slog.Attr {
	attrs := make([]slog.Attr, 0, len(args)/2)
	for i := 0; i < len(args); i += 2 {
		if i+1 < len(args) {
			attrs = append(attrs, slog.Any(fmt.Sprint(args[i]), args[i+1]))
		} else {
			attrs = append(attrs, slog.Any("extra", args[i]))
		}
	}
	return attrs
}

// Fatal logs a message and exits the application
func Fatal(message string, err error) {
	logRecord(slog.LevelError, 3, message, slog.Any("error", err))
	panic(err)
}

// Fatalf logs a formatted message and exits the application
func Fatalf(format string, v ...any) {
	msg := fmt.Sprintf(format, v...)
	logRecord(slog.LevelError, 3, msg)
	panic(msg)
}

// LogRaw logs a message without additional attributes
func LogRaw(s string) {
	logRecord(slog.LevelInfo, 3, s)
}
