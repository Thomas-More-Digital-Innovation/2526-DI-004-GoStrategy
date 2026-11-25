package db

import (
	"context"
	"crypto/rand"
	"digital-innovation/gostrategy/internal/models"
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

const (
	// friendCodeAlphabet is Crockford base32: no I, L, O or U, so codes survive being read aloud.
	friendCodeAlphabet    = "0123456789ABCDEFGHJKMNPQRSTVWXYZ"
	friendCodeLength      = 8
	friendCodeGroupSize   = 4
	maxFriendCodeAttempts = 5
)

var (
	// ErrInvalidFriendCode means the input cannot be a friend code however it is spelled.
	ErrInvalidFriendCode = errors.New("invalid friend code")
	// ErrFriendCodeNotFound means no active user holds this code.
	ErrFriendCodeNotFound = errors.New("friend code not found")
	errFriendCodeAttempts = errors.New("could not find an unused friend code")
)

// friendCodeLookalikes maps the letters Crockford treats as digits.
var friendCodeLookalikes = strings.NewReplacer(" ", "", "-", "", "O", "0", "I", "1", "L", "1")

// normalizeFriendCode canonicalises user input: separators dropped, uppercased, lookalikes mapped.
func normalizeFriendCode(raw string) (string, error) {
	code := friendCodeLookalikes.Replace(strings.ToUpper(raw))
	if len(code) != friendCodeLength {
		return "", ErrInvalidFriendCode
	}
	for _, r := range code {
		if !strings.ContainsRune(friendCodeAlphabet, r) {
			return "", ErrInvalidFriendCode
		}
	}
	return code, nil
}

func generateFriendCode() (string, error) {
	buf := make([]byte, friendCodeLength)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("read random bytes: %w", err)
	}
	for i, b := range buf {
		// why: the alphabet has exactly 32 symbols, so masking 5 bits is unbiased
		buf[i] = friendCodeAlphabet[b&31]
	}
	return string(buf), nil
}

func formatFriendCode(code string) string {
	return code[:friendCodeGroupSize] + "-" + code[friendCodeGroupSize:]
}

// GetOrCreateFriendCode returns the user's code, generating one on first use.
func GetOrCreateFriendCode(ctx context.Context, userID int) (string, error) {
	code, err := loadFriendCode(ctx, userID)
	if err != nil {
		return "", err
	}
	if code != nil {
		return formatFriendCode(*code), nil
	}
	return assignFriendCode(ctx, userID, false, generateFriendCode)
}

// RegenerateFriendCode replaces the user's code, so every old code and link stops working.
func RegenerateFriendCode(ctx context.Context, userID int) (string, error) {
	return assignFriendCode(ctx, userID, true, generateFriendCode)
}

// FindUserByFriendCode resolves any spelling of a code to its active owner.
func FindUserByFriendCode(ctx context.Context, raw string) (*models.User, error) {
	code, err := normalizeFriendCode(raw)
	if err != nil {
		return nil, err
	}
	var user models.User
	err = DB.WithContext(ctx).Where("friend_code = ?", code).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrFriendCodeNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find user by friend code: %w", err)
	}
	return &user, nil
}

func loadFriendCode(ctx context.Context, userID int) (*string, error) {
	var user models.User
	if err := DB.WithContext(ctx).Select("id", "friend_code").First(&user, userID).Error; err != nil {
		return nil, fmt.Errorf("load friend code of user %d: %w", userID, err)
	}
	return user.FriendCode, nil
}

// assignFriendCode stores a fresh code, retrying on collisions. Without replace it only fills an
// empty column, so two concurrent first calls end up agreeing on one code.
func assignFriendCode(ctx context.Context, userID int, replace bool, gen func() (string, error)) (string, error) {
	for range maxFriendCodeAttempts {
		code, err := gen()
		if err != nil {
			return "", err
		}
		query := DB.WithContext(ctx).Model(&models.User{}).Where("id = ?", userID)
		if !replace {
			query = query.Where("friend_code IS NULL")
		}
		result := query.Update("friend_code", code)
		switch {
		case isUniqueViolation(result.Error):
			continue
		case result.Error != nil:
			return "", fmt.Errorf("save friend code: %w", result.Error)
		case result.RowsAffected > 0:
			return formatFriendCode(code), nil
		}
		return existingFriendCode(ctx, userID)
	}
	return "", fmt.Errorf("%w after %d attempts", errFriendCodeAttempts, maxFriendCodeAttempts)
}

// existingFriendCode covers the lost race: another request filled the column first.
func existingFriendCode(ctx context.Context, userID int) (string, error) {
	code, err := loadFriendCode(ctx, userID)
	if err != nil {
		return "", err
	}
	if code == nil {
		return "", fmt.Errorf("user %d has no friend code: %w", userID, gorm.ErrRecordNotFound)
	}
	return formatFriendCode(*code), nil
}
