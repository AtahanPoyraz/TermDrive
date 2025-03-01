package validator

import (
	"errors"
	"regexp"
)

// IsValidUsername checks if the given username meets the format requirements.
// The username must be at least 3 characters long and can only contain
// letters (both uppercase and lowercase), numbers, and underscores.
// It cannot contain spaces or special characters.
func IsValidUsername(username string) error {
	if len(username) < 3 {
		return errors.New("username must be at least 3 characters long")
	}

	var usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	if !usernameRegex.MatchString(username) {
		return errors.New("username can only contain letters, numbers, and underscores, and cannot contain spaces")
	}

	return nil
}

// IsValidEmail validates whether the provided email address follows
// the standard email format (e.g., example@domain.com).
// It checks for the presence of '@' and a valid domain structure.
func IsValidEmail(email string) error {
	var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return errors.New("invalid email format")
	}

	return nil
}

// IsValidPassword checks the strength of the provided password.
// A valid password must contain at least one uppercase letter,
// one lowercase letter, one digit, and one special character
func IsValidPassword(password string) bool {
	var (
		upperCase   = regexp.MustCompile(`[A-Z]`)
		lowerCase   = regexp.MustCompile(`[a-z]`)
		digit       = regexp.MustCompile(`[0-9]`)
		punctuation = regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>]`)
	)

	return upperCase.MatchString(password) &&
		lowerCase.MatchString(password) &&
		digit.MatchString(password) &&
		punctuation.MatchString(password)
}
