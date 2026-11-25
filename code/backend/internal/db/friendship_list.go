package db

import (
	"context"
	"digital-innovation/gostrategy/internal/models"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type friendRow struct {
	ID             int
	RequesterID    int
	Status         models.FriendshipStatus
	CreatedAt      time.Time
	UserID         int
	Username       string
	ProfilePicture string
	UserCreatedAt  time.Time
	UserUpdatedAt  time.Time
}

// ListFriendships returns the user's friends plus the requests waiting on them and on others.
func ListFriendships(ctx context.Context, userID int) (*models.FriendListResponse, error) {
	var rows []friendRow
	err := WithRLS(WithUserID(ctx, userID), func(tx *gorm.DB) error {
		// why: Table() plus a join gets no automatic soft-delete filter, so friends who deleted
		// their account are dropped by the explicit u.deleted_at condition.
		return tx.Table("friendships AS f").
			Select("f.id, f.requester_id, f.status, f.created_at, "+
				"u.id AS user_id, u.username, COALESCE(u.profile_picture, '') AS profile_picture, "+
				"u.created_at AS user_created_at, u.updated_at AS user_updated_at").
			Joins("JOIN users u ON u.id = CASE WHEN f.requester_id = ? THEN f.addressee_id ELSE f.requester_id END", userID).
			Where("(f.requester_id = ? OR f.addressee_id = ?) AND u.deleted_at IS NULL", userID, userID).
			Order("f.created_at DESC, f.id DESC").
			Scan(&rows).Error
	})
	if err != nil {
		return nil, fmt.Errorf("list friendships of user %d: %w", userID, err)
	}

	list := &models.FriendListResponse{
		Friends:  []models.FriendEntry{},
		Incoming: []models.FriendEntry{},
		Outgoing: []models.FriendEntry{},
	}
	for _, r := range rows {
		entry := models.FriendEntry{
			ID: r.ID,
			User: models.UserDTO{
				ID:             r.UserID,
				Username:       r.Username,
				ProfilePicture: r.ProfilePicture,
				CreatedAt:      r.UserCreatedAt,
				UpdatedAt:      r.UserUpdatedAt,
			},
			CreatedAt: r.CreatedAt,
		}
		switch {
		case r.Status == models.FriendshipAccepted:
			list.Friends = append(list.Friends, entry)
		case r.RequesterID == userID:
			list.Outgoing = append(list.Outgoing, entry)
		default:
			list.Incoming = append(list.Incoming, entry)
		}
	}
	return list, nil
}
