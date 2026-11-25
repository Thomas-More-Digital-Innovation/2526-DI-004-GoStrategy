package db

import (
	"context"
	"digital-innovation/gostrategy/internal/models"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

const (
	// MaxFriends caps accepted friendships per user.
	MaxFriends = 200
	// MaxOutgoingPending caps the unanswered requests one user can have out at once.
	MaxOutgoingPending = 50
)

// Friendship errors. Handlers map each one to an HTTP status.
var (
	ErrSelfRequest        = errors.New("cannot send a friend request to yourself")
	ErrAlreadyFriends     = errors.New("already friends")
	ErrAlreadyRequested   = errors.New("friend request already sent")
	ErrFriendLimit        = errors.New("friend limit reached")
	ErrPendingLimit       = errors.New("too many pending friend requests")
	ErrFriendshipNotFound = errors.New("friendship not found")
)

// SendFriendRequest asks addresseeID to be friends. If they already asked requesterID, that
// request is accepted instead and the returned row has the accepted status.
//
// All friendship functions take the acting user explicitly and run as that user under RLS.
func SendFriendRequest(ctx context.Context, requesterID, addresseeID int) (*models.Friendship, error) {
	if requesterID == addresseeID {
		return nil, ErrSelfRequest
	}
	if err := ensureOtherCapacity(ctx, addresseeID); err != nil {
		return nil, err
	}
	f, err := sendOnce(ctx, requesterID, addresseeID)
	if isUniqueViolation(err) {
		// why: a concurrent send won the pair index; its row is committed now, so one retry sees it.
		// A failed statement poisons a Postgres transaction, which is why the retry is a new one.
		f, err = sendOnce(ctx, requesterID, addresseeID)
		if isUniqueViolation(err) {
			return nil, ErrAlreadyRequested
		}
	}
	return f, err
}

func sendOnce(ctx context.Context, requesterID, addresseeID int) (*models.Friendship, error) {
	var result models.Friendship
	err := WithRLS(WithUserID(ctx, requesterID), func(tx *gorm.DB) error {
		existing, found, err := findPair(tx, requesterID, addresseeID)
		if err != nil {
			return err
		}
		if found {
			result = existing
			return answerExisting(tx, requesterID, &result)
		}
		if err := ensureCanRequest(tx, requesterID); err != nil {
			return err
		}
		result = models.Friendship{RequesterID: requesterID, AddresseeID: addresseeID, Status: models.FriendshipPending}
		if err := tx.Create(&result).Error; err != nil {
			return fmt.Errorf("create friendship: %w", err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// AcceptFriendRequest lets the addressee accept a pending request.
func AcceptFriendRequest(ctx context.Context, userID, friendshipID int) (*models.Friendship, error) {
	var f models.Friendship
	err := WithRLS(WithUserID(ctx, userID), func(tx *gorm.DB) error {
		err := tx.Where("id = ? AND addressee_id = ? AND status = ?", friendshipID, userID, models.FriendshipPending).Take(&f).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrFriendshipNotFound
		}
		if err != nil {
			return fmt.Errorf("load friendship %d: %w", friendshipID, err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if err := ensureOtherCapacity(ctx, f.RequesterID); err != nil {
		return nil, err
	}
	err = WithRLS(WithUserID(ctx, userID), func(tx *gorm.DB) error {
		return acceptRow(tx, userID, &f)
	})
	if err != nil {
		return nil, err
	}
	return &f, nil
}

// DeleteFriendship removes a row. It serves decline, cancel and unfriend: either party may call it.
func DeleteFriendship(ctx context.Context, userID, friendshipID int) error {
	return WithRLS(WithUserID(ctx, userID), func(tx *gorm.DB) error {
		result := tx.Where("id = ? AND (requester_id = ? OR addressee_id = ?)", friendshipID, userID, userID).Delete(&models.Friendship{})
		if result.Error != nil {
			return fmt.Errorf("delete friendship %d: %w", friendshipID, result.Error)
		}
		if result.RowsAffected == 0 {
			return ErrFriendshipNotFound
		}
		return nil
	})
}

func findPair(tx *gorm.DB, a, b int) (f models.Friendship, found bool, err error) {
	var rows []models.Friendship
	err = tx.Where("(requester_id = ? AND addressee_id = ?) OR (requester_id = ? AND addressee_id = ?)", a, b, b, a).Limit(1).Find(&rows).Error
	if err != nil {
		return f, false, fmt.Errorf("look up friendship: %w", err)
	}
	if len(rows) == 0 {
		return f, false, nil
	}
	return rows[0], true, nil
}

// answerExisting handles a send when the pair already has a row: the other side asking first
// means both want it, so accept.
func answerExisting(tx *gorm.DB, userID int, existing *models.Friendship) error {
	switch {
	case existing.Status == models.FriendshipAccepted:
		return ErrAlreadyFriends
	case existing.RequesterID == userID:
		return ErrAlreadyRequested
	}
	return acceptRow(tx, userID, existing)
}

func acceptRow(tx *gorm.DB, userID int, f *models.Friendship) error {
	if err := ensureFriendCapacity(tx, userID); err != nil {
		return err
	}
	result := tx.Model(f).Where("addressee_id = ? AND status = ?", userID, models.FriendshipPending).Update("status", models.FriendshipAccepted)
	if result.Error != nil {
		return fmt.Errorf("accept friendship %d: %w", f.ID, result.Error)
	}
	if result.RowsAffected == 0 {
		return ErrFriendshipNotFound
	}
	f.Status = models.FriendshipAccepted
	return nil
}

func ensureCanRequest(tx *gorm.DB, userID int) error {
	if err := ensureFriendCapacity(tx, userID); err != nil {
		return err
	}
	pending, err := countFriendships(tx, "status = ? AND requester_id = ?", models.FriendshipPending, userID)
	if err != nil {
		return err
	}
	if pending >= MaxOutgoingPending {
		return ErrPendingLimit
	}
	return nil
}

// ensureOtherCapacity checks the other party's cap while running as that party, because RLS hides
// their friendships from the acting user. It runs before the acting user's transaction so a request
// never holds two connections at once.
//
// shortcut: the check and the write are separate transactions, so two simultaneous requests can
// push someone past MaxFriends by a few; move the count into the write if that ever matters.
func ensureOtherCapacity(ctx context.Context, otherID int) error {
	return WithRLS(WithUserID(ctx, otherID), func(tx *gorm.DB) error {
		return ensureFriendCapacity(tx, otherID)
	})
}

// ensureFriendCapacity checks one user's cap through that user's own RLS view. Rows of soft-deleted
// friends count too, as there is no account deletion flow yet.
func ensureFriendCapacity(tx *gorm.DB, userID int) error {
	friends, err := countFriendships(tx, "status = ? AND (requester_id = ? OR addressee_id = ?)", models.FriendshipAccepted, userID, userID)
	if err != nil {
		return err
	}
	if friends >= MaxFriends {
		return ErrFriendLimit
	}
	return nil
}

func countFriendships(tx *gorm.DB, where string, args ...any) (int64, error) {
	var n int64
	if err := tx.Model(&models.Friendship{}).Where(where, args...).Count(&n).Error; err != nil {
		return 0, fmt.Errorf("count friendships: %w", err)
	}
	return n, nil
}
