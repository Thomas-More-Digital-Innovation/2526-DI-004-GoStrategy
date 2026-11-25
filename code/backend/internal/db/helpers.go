// Package db provides database access and helpers
package db

import "strings"

func isPostgresDialect() bool {
	return DB.Name() == "postgres"
}

// isUniqueViolation reports a unique-constraint failure on both Postgres and SQLite.
//
// ponytail: matches the error text because GORM's TranslateError is off; switch to
// errors.Is(err, gorm.ErrDuplicatedKey) if it is ever enabled.
func isUniqueViolation(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "unique") || strings.Contains(msg, "duplicate")
}
