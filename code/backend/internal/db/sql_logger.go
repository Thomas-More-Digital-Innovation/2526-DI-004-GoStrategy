package db

import (
	"time"

	"gorm.io/gorm/logger"
)

// slowQueryThreshold is GORM's own default for flagging slow statements.
const slowQueryThreshold = 200 * time.Millisecond

// newSQLLogger logs every statement with placeholders instead of values.
//
// why: with values inlined, password hashes and friend codes end up in the server log.
func newSQLLogger(w logger.Writer) logger.Interface {
	return logger.New(w, logger.Config{
		SlowThreshold:        slowQueryThreshold,
		LogLevel:             logger.Info,
		Colorful:             true,
		ParameterizedQueries: true,
	})
}
