package service

import (
	"golang.org/x/crypto/bcrypt"
)

type Hasher struct {
}

func NewHasher() Hasher {
	return Hasher{}
}

func (h *Hasher) Hash(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return "", err
	}

	return string(hashed), nil
}

func (h *Hasher) IsValid(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
