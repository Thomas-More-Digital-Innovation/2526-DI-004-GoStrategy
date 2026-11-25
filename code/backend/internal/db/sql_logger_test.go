package db

import (
	"bytes"
	"digital-innovation/gostrategy/internal/models"
	"log"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestSQLLoggerKeepsValuesOutOfTheLog(t *testing.T) {
	var out bytes.Buffer
	db, err := gorm.Open(sqlite.Open("file:sqllogger?mode=memory&cache=shared"), &gorm.Config{Logger: newSQLLogger(log.New(&out, "", 0))})
	require.NoError(t, err)
	t.Cleanup(func() {
		if sqlDB, err := db.DB(); err == nil {
			_ = sqlDB.Close()
		}
	})
	require.NoError(t, db.AutoMigrate(&models.User{}))

	friendCode, secret := "K7QXM2PA", "value-must-not-leak"
	require.NoError(t, db.Create(&models.User{Username: "logtest", PasswordHash: secret, FriendCode: &friendCode}).Error)

	assert.Contains(t, out.String(), "INSERT INTO", "statements are still logged")
	assert.NotContains(t, out.String(), secret)
	assert.NotContains(t, out.String(), friendCode)
}
