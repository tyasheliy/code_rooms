package hasher

import "golang.org/x/crypto/bcrypt"

type BcryptHasher struct{}

func (h *BcryptHasher) Hash(s string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func (h *BcryptHasher) Check(s string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(s))

	return err == nil
}
