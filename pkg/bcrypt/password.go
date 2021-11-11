package bcrypt

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

var ErrMismatchedHashAndPassword = bcrypt.ErrMismatchedHashAndPassword

// BcryptPassword generates a hash from the user's password (adding the pepper).
func BcryptPassword(password, pepper string) (string, error) {
	pwBytes := []byte(password + pepper)
	hashed, err := bcrypt.GenerateFromPassword(pwBytes, bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("bcryptPassword: %w", err)
	}

	return string(hashed), nil
}

// CompareHashAndPassword compares a bcrypt hashed password with its possible
// plaintext equivalent in constant time.
func CompareHashAndPassword(password, hash string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return ErrMismatchedHashAndPassword
		}
		return err
	}
	return nil
}
