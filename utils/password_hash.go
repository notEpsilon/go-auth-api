package utils

import "golang.org/x/crypto/bcrypt"

type Hasher interface {
	Hash([]byte) ([]byte, error)
	Compare([]byte, []byte) bool
}

func NewHasher(name string) Hasher {
	switch name {
	case "bcrypt":
		return BCryptHasher{}
	default:
		return BCryptHasher{}
	}
}

type BCryptHasher struct{}

func (h BCryptHasher) Hash(password []byte) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return hash, nil
}

func (h BCryptHasher) Compare(hashed []byte, plain []byte) bool {
	return bcrypt.CompareHashAndPassword(hashed, plain) == nil
}

var _ Hasher = BCryptHasher{}
