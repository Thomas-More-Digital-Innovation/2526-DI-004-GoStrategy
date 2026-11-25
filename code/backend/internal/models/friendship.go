package models

import "time"

// FriendshipStatus is the lifecycle state of a friend request.
type FriendshipStatus string

// Friendship statuses. Declining, cancelling and unfriending delete the row.
const (
	FriendshipPending  FriendshipStatus = "pending"
	FriendshipAccepted FriendshipStatus = "accepted"
)

// Friendship is one directed row per relationship: the requester asked, the addressee answers.
//
// why: unlike every other model there is no gorm.DeletedAt. A soft-deleted row would keep the
// pair index occupied and block re-requesting after a decline.
//
// The pair index treats A→B and B→A as the same pair. The CASE expressions carry their own
// parentheses because Postgres requires them around non-function index expressions, and they
// avoid commas because GORM splits index tag options on ",".
type Friendship struct {
	ID          int              `json:"id" gorm:"primaryKey"`
	RequesterID int              `json:"requester_id" gorm:"not null;index;uniqueIndex:idx_friendships_pair,expression:(CASE WHEN requester_id < addressee_id THEN requester_id ELSE addressee_id END)"`
	AddresseeID int              `json:"addressee_id" gorm:"not null;index;uniqueIndex:idx_friendships_pair,expression:(CASE WHEN requester_id < addressee_id THEN addressee_id ELSE requester_id END)"`
	Status      FriendshipStatus `json:"status" gorm:"not null;size:20" swaggertype:"string"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
}

// FriendCodeResponse carries the owner's own friend code, formatted as XXXX-XXXX.
type FriendCodeResponse struct {
	Code string `json:"code"`
}

// FriendEntry is one row of the social tab: the friendship id plus the other user.
type FriendEntry struct {
	ID        int       `json:"id"`
	User      UserDTO   `json:"user"`
	CreatedAt time.Time `json:"created_at"`
}

// FriendListResponse groups a user's friendships by what they can do with each one.
type FriendListResponse struct {
	Friends  []FriendEntry `json:"friends"`
	Incoming []FriendEntry `json:"incoming"`
	Outgoing []FriendEntry `json:"outgoing"`
}

// AddFriendRequest asks to become friends with the owner of a friend code.
type AddFriendRequest struct {
	Code string `json:"code"`
}

// FriendPreviewResponse shows who a friend code belongs to, before anyone sends a request.
type FriendPreviewResponse struct {
	User UserDTO `json:"user"`
}
