package handlers

import (
	"digital-innovation/gostrategy/internal/api/core"
	"digital-innovation/gostrategy/internal/db"
	"digital-innovation/gostrategy/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// FriendCodeHandler returns the current user's friend code, creating it on first use
// @Summary Get my friend code
// @Description Return the friend code of the logged-in user. The first call creates it.
// @Tags friends
// @Produce json
// @Success 200 {object} models.FriendCodeResponse
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 429 {object} map[string]string "Too many requests"
// @Failure 500 {object} map[string]string "Failed to get friend code"
// @Router /users/me/friend-code [post]
func (h *Handler) FriendCodeHandler(c *gin.Context) {
	user := core.EnsureAuthenticated(c)
	if user == nil {
		return
	}

	code, err := db.GetOrCreateFriendCode(c.Request.Context(), user.ID)
	if err != nil {
		sendFriendError(c, err, "Failed to get friend code")
		return
	}
	core.SendJSON(c, models.FriendCodeResponse{Code: code}, http.StatusOK)
}

// RegenerateFriendCodeHandler replaces the current user's friend code
// @Summary Regenerate my friend code
// @Description Replace the friend code of the logged-in user. Old codes and invite links stop working.
// @Tags friends
// @Produce json
// @Success 200 {object} models.FriendCodeResponse
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 429 {object} map[string]string "Too many requests"
// @Failure 500 {object} map[string]string "Failed to regenerate friend code"
// @Router /users/me/friend-code/regenerate [post]
func (h *Handler) RegenerateFriendCodeHandler(c *gin.Context) {
	user := core.EnsureAuthenticated(c)
	if user == nil {
		return
	}

	code, err := db.RegenerateFriendCode(c.Request.Context(), user.ID)
	if err != nil {
		sendFriendError(c, err, "Failed to regenerate friend code")
		return
	}
	core.SendJSON(c, models.FriendCodeResponse{Code: code}, http.StatusOK)
}
