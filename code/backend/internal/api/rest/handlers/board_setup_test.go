package handlers_test

import (
	"bytes"
	"context"
	"digital-innovation/gostrategy/internal/db"
	"digital-innovation/gostrategy/internal/models"
	"digital-innovation/gostrategy/internal/testutils"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBoardSetupHandlers(t *testing.T) {
	_, h, _ := testutils.SetupHandlerTest(t)

	user, _ := db.CreateUser(context.Background(), "alice", "StrongPassword1", "")
	var createdSetup models.BoardSetup

	t.Run("Create", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		setupReq := models.BoardSetup{
			Name:        "Test Setup",
			Description: "A test description",
			SetupData:   "some-binary-data",
		}
		jsonBody, _ := json.Marshal(setupReq)
		c.Request, _ = http.NewRequest("POST", "/users/me/board-setups", bytes.NewBuffer(jsonBody))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Set("user", user)
		h.CreateBoardSetupHandler(c)
		assert.Equal(t, http.StatusCreated, w.Code)

		err := json.Unmarshal(w.Body.Bytes(), &createdSetup)
		assert.NoError(t, err)
	})

	t.Run("CreateLimitEnforced", func(t *testing.T) {
		limitUser, _ := db.CreateUser(context.Background(), "limit_user", "StrongPassword1", "")
		ctx := db.WithUserID(context.Background(), limitUser.ID)
		for i := 1; i <= 10; i++ {
			_, err := db.CreateBoardSetup(ctx, limitUser.ID, fmt.Sprintf("Setup %d", i), "desc", "dummy-data", false)
			require.NoError(t, err)
		}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		eleventhReq := models.BoardSetup{
			Name:        "11th Setup",
			Description: "Should be rejected",
			SetupData:   "dummy-data",
		}
		jsonBody, _ := json.Marshal(eleventhReq)
		c.Request, _ = http.NewRequest("POST", "/users/me/board-setups", bytes.NewBuffer(jsonBody))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Set("user", limitUser)
		h.CreateBoardSetupHandler(c)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Maximum number of board setups reached (10)")
	})

	t.Run("List", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/users/me/board-setups", nil)
		c.Set("user", user)
		h.GetUserBoardSetupsHandler(c)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Get", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", fmt.Sprintf("/users/me/board-setups/%d", createdSetup.ID), nil)
		c.Params = []gin.Param{{Key: "id", Value: fmt.Sprintf("%d", createdSetup.ID)}}
		c.Set("user", user)
		h.GetBoardSetupHandler(c)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Update", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		setupReq := models.BoardSetup{
			Name: "Updated Name",
		}
		jsonBodyUpdate, _ := json.Marshal(setupReq)
		c.Request, _ = http.NewRequest("PUT", fmt.Sprintf("/users/me/board-setups/%d", createdSetup.ID), bytes.NewBuffer(jsonBodyUpdate))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Params = []gin.Param{{Key: "id", Value: fmt.Sprintf("%d", createdSetup.ID)}}
		c.Set("user", user)
		h.UpdateBoardSetupHandler(c)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Delete", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("DELETE", fmt.Sprintf("/users/me/board-setups/%d", createdSetup.ID), nil)
		c.Params = []gin.Param{{Key: "id", Value: fmt.Sprintf("%d", createdSetup.ID)}}
		c.Set("user", user)
		h.DeleteBoardSetupHandler(c)
		assert.Equal(t, http.StatusNoContent, w.Code)
	})
}
