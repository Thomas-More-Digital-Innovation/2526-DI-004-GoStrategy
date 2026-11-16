package utils

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

var (
	// These regexes match characters that are NOT allowed
	// The ^ character negates the set of characters inside []. So it matches any character that is NOT in the set.
	// Examples:
	// [^a-zA-Z0-9_] matches any character that is not a letter, number, or underscore
	// [^a-zA-Z0-9!@#$%^&*()_+=\-\. ] matches any character that is not a letter, number, special character, or space.
	// [^<>] matches any character that is not < or >
	usernameForbiddenRegex = regexp.MustCompile(`[^a-zA-Z0-9_]`)
	passwordForbiddenRegex = regexp.MustCompile(`[^a-zA-Z0-9!@#$%^&*()_+=\-\. ]`)
	genericForbiddenRegex  = regexp.MustCompile(`[<>"'%;]`)
)

// ValidateUsername returns an error if the username contains invalid characters
func ValidateUsername(username string) error {
	if usernameForbiddenRegex.MatchString(username) {
		return fmt.Errorf("username contains invalid characters (only alphanumeric and underscores allowed)")
	}
	return nil
}

// Username length bounds enforced on registration.
const (
	UsernameMinLen = 3
	UsernameMaxLen = 50
)

// ValidateNewUsername validates characters and length of a username being registered
func ValidateNewUsername(username string) error {
	if err := ValidateUsername(username); err != nil {
		return err
	}
	if len(username) < UsernameMinLen || len(username) > UsernameMaxLen {
		return fmt.Errorf("username must be %d-%d characters", UsernameMinLen, UsernameMaxLen)
	}
	return nil
}

// ValidatePassword returns an error if the password contains invalid characters or edge spaces
func ValidatePassword(password string) error {
	if strings.HasPrefix(password, " ") || strings.HasSuffix(password, " ") {
		return fmt.Errorf("password cannot start or end with a space")
	}
	if passwordForbiddenRegex.MatchString(password) {
		return fmt.Errorf(`password contains invalid characters (allowed: a-z, A-Z, 0-9, spaces, and !@#$%%^&*()_+=-.)`)
	}
	return nil
}

// PasswordMinLen is the minimum length of a newly set password.
const PasswordMinLen = 8

// ValidateNewPassword validates characters and complexity of a password being set
func ValidateNewPassword(password string) error {
	if err := ValidatePassword(password); err != nil {
		return err
	}
	var hasNumber, hasUpper, hasLower bool
	for _, char := range password {
		switch {
		case unicode.IsDigit(char):
			hasNumber = true
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		}
	}
	if len(password) < PasswordMinLen || !hasNumber || !hasUpper || !hasLower {
		return fmt.Errorf("password must be at least %d characters and contain at least one number, one uppercase and one lowercase letter", PasswordMinLen)
	}
	return nil
}

// ValidateGeneric returns an error if the input contains dangerous characters
func ValidateGeneric(input string, fieldName string) error {
	if genericForbiddenRegex.MatchString(input) {
		return fmt.Errorf("%s contains invalid characters (< > \" ' %% ; are not allowed)", fieldName)
	}
	return nil
}
