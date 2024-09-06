package password

import (
	"golang.org/x/crypto/bcrypt"
)

func Hash(plaintextPassword string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func CompareWithHash(plaintextPassword, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plaintextPassword))
}
