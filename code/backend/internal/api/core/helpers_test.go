package core_test

import (
	"digital-innovation/gostrategy/internal/api/core"
	"digital-innovation/gostrategy/internal/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHelpers(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("SendError", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		core.SendError(c, "Test error", http.StatusBadRequest)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("SendJSON", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		core.SendJSON(c, gin.H{"test": "ok"}, http.StatusOK)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("EnsureAuthenticated", func(t *testing.T) {
		t.Run("Unauthorized", func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			user := core.EnsureAuthenticated(c)
			assert.Nil(t, user)
			assert.Equal(t, http.StatusUnauthorized, w.Code)
		})

		t.Run("Authorized", func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			userObj := &models.User{ID: 1, Username: "test"}
			c.Set("user", userObj)
			user := core.EnsureAuthenticated(c)
			assert.Equal(t, userObj, user)
		})
	})

	t.Run("ParseID", func(t *testing.T) {
		t.Run("FromParam", func(t *testing.T) {
			c, _ := gin.CreateTestContext(httptest.NewRecorder())
			c.Params = []gin.Param{{Key: "id", Value: "123"}}
			id, err := core.ParseID(c, "id")
			assert.NoError(t, err)
			assert.Equal(t, 123, id)
		})

		t.Run("FromQuery", func(t *testing.T) {
			c, _ := gin.CreateTestContext(httptest.NewRecorder())
			c.Request, _ = http.NewRequest("GET", "/?id=456", nil)
			id, err := core.ParseID(c, "id")
			assert.NoError(t, err)
			assert.Equal(t, 456, id)
		})
	})

	t.Run("ParsePagination", func(t *testing.T) {
		t.Run("DefaultValues", func(t *testing.T) {
			c, _ := gin.CreateTestContext(httptest.NewRecorder())
			c.Request, _ = http.NewRequest("GET", "/test", nil)
			limit, offset := core.ParsePagination(c, 10, 50)
			assert.Equal(t, 10, limit)
			assert.Equal(t, 0, offset)
		})

		t.Run("CustomValidValues", func(t *testing.T) {
			c, _ := gin.CreateTestContext(httptest.NewRecorder())
			c.Request, _ = http.NewRequest("GET", "/test?limit=25&offset=50", nil)
			limit, offset := core.ParsePagination(c, 10, 50)
			assert.Equal(t, 25, limit)
			assert.Equal(t, 50, offset)
		})

		t.Run("ClampsExceedingMaxLimit", func(t *testing.T) {
			c, _ := gin.CreateTestContext(httptest.NewRecorder())
			c.Request, _ = http.NewRequest("GET", "/test?limit=100&offset=0", nil)
			limit, offset := core.ParsePagination(c, 10, 50)
			assert.Equal(t, 50, limit)
			assert.Equal(t, 0, offset)
		})

		t.Run("RejectsNegativeOffsetAndZeroLimit", func(t *testing.T) {
			c, _ := gin.CreateTestContext(httptest.NewRecorder())
			c.Request, _ = http.NewRequest("GET", "/test?limit=0&offset=-10", nil)
			limit, offset := core.ParsePagination(c, 10, 50)
			assert.Equal(t, 10, limit)
			assert.Equal(t, 0, offset)
		})
	})
}
