package creds

import (
	"strings"
	"unicode"

	"github.com/Koshsky/erp-backend/pkg/errors"
)

// minPasswordLength — minimum length of a user password.
const minPasswordLength = 8

// passwordPolicyMessage — password requirements (mirrors the frontend
// usePasswordValidation rules: length, letter cases, digit, special character).
//
//nolint:gosec // policy message, not a secret
const passwordPolicyMessage = "пароль не соответствует требованиям: не менее 8 символов, строчные и заглавные буквы, цифра, спецсимвол"

// ValidatePassword checks a user password against the complexity policy
// (AD-09). Applied only to user passwords (password change);
// admin auto-generated passwords are handled separately.
func ValidatePassword(password string) error {
	if !validPassword(password) {
		return errors.BadRequest(passwordPolicyMessage)
	}
	return nil
}

func validPassword(password string) bool {
	if len([]rune(password)) < minPasswordLength {
		return false
	}
	var hasLower, hasUpper, hasDigit, hasSpecial bool
	for _, r := range password {
		switch {
		case unicode.IsLower(r):
			hasLower = true
		case unicode.IsUpper(r):
			hasUpper = true
		case unicode.IsNumber(r):
			hasDigit = true
		case !unicode.IsLetter(r) && !unicode.IsNumber(r):
			hasSpecial = true
		}
	}
	// Whitespace/empty strings do not satisfy the policy: no special characters or letters.
	return hasLower && hasUpper && hasDigit && hasSpecial && strings.TrimSpace(password) != ""
}
