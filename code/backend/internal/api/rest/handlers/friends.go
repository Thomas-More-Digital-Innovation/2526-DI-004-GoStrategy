package handlers

import (
	"digital-innovation/gostrategy/internal/api/core"
	"digital-innovation/gostrategy/internal/db"
	"digital-innovation/gostrategy/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ListFriendsHandler returns the current user's friends and pending requests
// @Summary List friends
// @Description Friends, incoming requests and outgoing requests of the logged-in user
// @Tags friends
// @Produce json
// @Success 200 {object} models.FriendListResponse
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Failed to list friends"
// @Router /friends [get]
func (h *Handler) ListFriendsHandler(c *gin.Context) {
	user := core.EnsureAuthenticated(c)
	if user == nil {
		return
	}

	list, err := db.ListFriendships(c.Request.Context(), user.ID)
	if err != nil {
		sendFriendError(c, err, "Failed to list friends")
		return
	}
	core.SendJSON(c, list, http.StatusOK)
}

// PreviewFriendHandler shows who a friend code belongs to
// @Summary Preview a friend code
// @Description Resolve a friend code to its owner so the invite page can ask for confirmation
// @Tags friends
// @Produce json
// @Param code query string true "Friend code, any spelling"
// @Success 200 {object} models.FriendPreviewResponse
// @Failure 400 {object} map[string]string "Invalid friend code"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "No player has this friend code"
// @Failure 429 {object} map[string]string "Too many requests"
// @Router /friends/preview [get]
func (h *Handler) PreviewFriendHandler(c *gin.Context) {
	if core.EnsureAuthenticated(c) == nil {
		return
	}

	owner, err := db.FindUserByFriendCode(c.Request.Context(), c.Query("code"))
	if err != nil {
		sendFriendError(c, err, "Failed to preview friend code")
		return
	}
	core.SendJSON(c, models.FriendPreviewResponse{User: owner.DTO()}, http.StatusOK)
}

// SendFriendRequestHandler sends a friend request to the owner of a friend code
// @Summary Send a friend request
// @Description Ask the owner of a friend code to be friends. If they already asked you, the request is accepted instead.
// @Tags friends
// @Accept json
// @Produce json
// @Param request body models.AddFriendRequest true "Friend code"
// @Success 200 {object} models.Friendship "Accepted an existing request from them"
// @Success 201 {object} models.Friendship "Request sent"
// @Failure 400 {object} map[string]string "Invalid body or friend code, or your own code"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "No player has this friend code"
// @Failure 409 {object} map[string]string "Already friends, already requested, or friend limit reached"
// @Failure 429 {object} map[string]string "Too many pending requests or too many calls"
// @Router /friends/requests [post]
func (h *Handler) SendFriendRequestHandler(c *gin.Context) {
	user := core.EnsureAuthenticated(c)
	if user == nil {
		return
	}

	var req models.AddFriendRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		core.SendError(c, "Invalid request body", http.StatusBadRequest)
		return
	}

	owner, err := db.FindUserByFriendCode(c.Request.Context(), req.Code)
	if err != nil {
		sendFriendError(c, err, "Failed to look up friend code")
		return
	}
	friendship, err := db.SendFriendRequest(c.Request.Context(), user.ID, owner.ID)
	if err != nil {
		sendFriendError(c, err, "Failed to send friend request")
		return
	}

	status := http.StatusCreated
	if friendship.Status == models.FriendshipAccepted {
		status = http.StatusOK
	}
	core.SendJSON(c, friendship, status)
}

// AcceptFriendRequestHandler accepts a pending request addressed to the current user
// @Summary Accept a friend request
// @Description Only the addressee can accept a pending request
// @Tags friends
// @Produce json
// @Param id path int true "Friendship ID"
// @Success 200 {object} models.Friendship
// @Failure 400 {object} map[string]string "Invalid ID"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "Friend request not found"
// @Failure 409 {object} map[string]string "Friend limit reached"
// @Failure 429 {object} map[string]string "Too many requests"
// @Router /friends/requests/{id}/accept [post]
func (h *Handler) AcceptFriendRequestHandler(c *gin.Context) {
	user := core.EnsureAuthenticated(c)
	if user == nil {
		return
	}
	id, ok := friendshipID(c)
	if !ok {
		return
	}

	friendship, err := db.AcceptFriendRequest(c.Request.Context(), user.ID, id)
	if err != nil {
		sendFriendError(c, err, "Failed to accept friend request")
		return
	}
	core.SendJSON(c, friendship, http.StatusOK)
}

// DeleteFriendshipHandler declines a request, cancels one, or removes a friend
// @Summary Remove a friendship or request
// @Description Either side can decline, cancel or unfriend. The row is deleted.
// @Tags friends
// @Param id path int true "Friendship ID"
// @Success 204 "Removed"
// @Failure 400 {object} map[string]string "Invalid ID"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "Friend request not found"
// @Failure 429 {object} map[string]string "Too many requests"
// @Router /friends/{id} [delete]
func (h *Handler) DeleteFriendshipHandler(c *gin.Context) {
	user := core.EnsureAuthenticated(c)
	if user == nil {
		return
	}
	id, ok := friendshipID(c)
	if !ok {
		return
	}

	if err := db.DeleteFriendship(c.Request.Context(), user.ID, id); err != nil {
		sendFriendError(c, err, "Failed to remove friendship")
		return
	}
	c.Status(http.StatusNoContent)
}

// friendshipID reads the :id path parameter, answering 400 itself when it is not a positive number.
func friendshipID(c *gin.Context) (int, bool) {
	id, err := core.ParseID(c, "id")
	if err != nil || id <= 0 {
		core.SendError(c, "Invalid friendship ID", http.StatusBadRequest)
		return 0, false
	}
	return id, true
}
