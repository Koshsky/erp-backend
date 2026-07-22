package password

import "golang.org/x/crypto/bcrypt"

type Hasher interface {
	Hash(raw string) (string, error)
	Compare(hashed, raw string) error
}

type BcryptHasher struct {
	cost int
}

func NewBcryptHasher() *BcryptHasher {
	return &BcryptHasher{cost: bcrypt.DefaultCost}
}

func (h *BcryptHasher) Hash(raw string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(raw), h.cost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func (h *BcryptHasher) Compare(hashed, raw string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(raw))
}
