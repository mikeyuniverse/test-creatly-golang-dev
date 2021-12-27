package hasher

import (
	"crypto/sha1"
	"fmt"
)

type Hasher struct {
	salt string
}

func New(salt string) *Hasher {
	return &Hasher{salt: salt}
}

func (h *Hasher) Hash(password string) (string, error) {
	hash := sha1.New()

	_, err := hash.Write([]byte(password))
	if err != nil {
		return "", err
	}

	_, err = hash.Write([]byte(h.salt))
	if err != nil {
		return "", err
	}

	passwordHash := fmt.Sprintf("%x", hash.Sum([]byte(h.salt)))
	return passwordHash, nil
}
