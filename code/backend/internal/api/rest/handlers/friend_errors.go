package handlers

import (
	"digital-innovation/gostrategy/internal/api/core"
	"digital-innovation/gostrategy/internal/db"
	"digital-innovation/gostrategy/internal/logging"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// friendErrorResponses maps each db sentinel error to its HTTP answer. 404s for friendships
// stay vague on purpose: they must not reveal rows that belong to someone else.
var friendErrorResponses = []struct {
	err     error
	status  int
	message string
}{
	{db.ErrInvalidFriendCode, http.StatusBadRequest, "Invalid friend code"},
	{db.ErrSelfRequest, http.StatusBadRequest, "You cannot add yourself as a friend"},
	{db.ErrFriendCodeNotFound, http.StatusNotFound, "No player has this friend code"},
	{db.ErrFriendshipNotFound, http.StatusNotFound, "Friend request not found"},
	{db.ErrAlreadyFriends, http.StatusConflict, "You are already friends"},
	{db.ErrAlreadyRequested, http.StatusConflict, "You already sent this player a friend request"},
	{db.ErrFriendLimit, http.StatusConflict, "Friend limit reached"},
	{db.ErrPendingLimit, http.StatusTooManyRequests, "Too many pending friend requests"},
}

// sendFriendError answers with the mapped status, or a 500 that hides the internal error.
func sendFriendError(c *gin.Context, err error, fallback string) {
	for _, known := range friendErrorResponses {
		if errors.Is(err, known.err) {
			core.SendError(c, known.message, known.status)
			return
		}
	}
	logging.Error(fallback, err)
	core.SendError(c, fallback, http.StatusInternalServerError)
}
