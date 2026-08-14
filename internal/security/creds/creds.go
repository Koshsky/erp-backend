// Package creds generates randomized login/password for worker accounts.
package creds

import (
	"crypto/rand"
	"math/big"
	"strings"
)

// passwordAlphabet excludes visually ambiguous characters (0/O, 1/l/I).
const passwordAlphabet = "abcdefghijkmnpqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ23456789"
const usernameAlphabet = "abcdefghijklmnopqrstuvwxyz0123456789"

const (
	passwordLength = 16
	usernameLength = 10
)

// cyrToLat — простая транслитерация кириллицы в латиницу (для login по фамилии).
//
//nolint:gochecknoglobals // статичная таблица
var cyrToLat = map[rune]string{
	'а': "a", 'б': "b", 'в': "v", 'г': "g", 'д': "d", 'е': "e", 'ё': "e",
	'ж': "zh", 'з': "z", 'и': "i", 'й': "y", 'к': "k", 'л': "l", 'м': "m",
	'н': "n", 'о': "o", 'п': "p", 'р': "r", 'с': "s", 'т': "t", 'у': "u",
	'ф': "f", 'х': "kh", 'ц': "ts", 'ч': "ch", 'ш': "sh", 'щ': "shch",
	'ъ': "", 'ы': "y", 'ь': "", 'э': "e", 'ю': "yu", 'я': "ya",
}

// TransliterateSurname возвращает латинскую транслитерацию фамилии (первое
// слово ФИО, порядок «Фамилия Имя Отчество») в нижнем регистре; пустая
// строка, если транслитерировать нечего.
func TransliterateSurname(name string) string {
	lower := strings.ToLower(strings.TrimSpace(name))
	if lower == "" {
		return ""
	}
	words := strings.Fields(lower)
	surname := words[0]

	var b strings.Builder
	for _, r := range surname {
		if s, ok := cyrToLat[r]; ok {
			b.WriteString(s)
		}
	}
	return b.String()
}

// RandomPassword returns a random password without ambiguous characters.
func RandomPassword() (string, error) {
	return randomString(passwordAlphabet, passwordLength)
}

// RandomUsernameSuffix returns a random alphanumeric suffix to build a unique username.
func RandomUsernameSuffix() (string, error) {
	return randomString(usernameAlphabet, usernameLength)
}

func randomString(alphabet string, length int) (string, error) {
	out := make([]byte, length)
	maxIdx := big.NewInt(int64(len(alphabet)))
	for i := range out {
		n, err := rand.Int(rand.Reader, maxIdx)
		if err != nil {
			return "", err
		}
		out[i] = alphabet[n.Int64()]
	}
	return string(out), nil
}
