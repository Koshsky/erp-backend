package creds

import (
	"strings"
	"unicode"

	"github.com/Koshsky/erp-backend/pkg/errors"
)

// minPasswordLength — минимальная длина пользовательского пароля.
const minPasswordLength = 8

// passwordPolicyMessage — требования к паролю (зеркалит frontend-правила
// usePasswordValidation: длина, регистры, цифра, спецсимвол).
//
//nolint:gosec // сообщение политики, а не секрет
const passwordPolicyMessage = "пароль не соответствует требованиям: не менее 8 символов, строчные и заглавные буквы, цифра, спецсимвол"

// ValidatePassword проверяет пользовательский пароль по политике сложности
// (AD-09). Применяется только к пользовательским паролям (смена пароля);
// авто-генерируемые пароли админом проходят отдельно.
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
	// Пробелы/пустые строки политике не подходят: спецсимволов и букв нет.
	return hasLower && hasUpper && hasDigit && hasSpecial && strings.TrimSpace(password) != ""
}
