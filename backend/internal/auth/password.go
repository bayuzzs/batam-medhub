package auth

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// dummyBcryptHash is a precomputed bcrypt hash used to make unknown-email lookups
// take similar computational time to known-email lookups, preventing email enumeration.
const dummyBcryptHash = "$2a$10$7EqJtq98hPqEX7fNZaFWoO.8H7yV23l8b7a6Wf2pS1C5J7b3z2f1a"

var (
	ErrPasswordTooShort = errors.New("password must be at least 8 characters")
	ErrPasswordTooLong  = errors.New("password must not exceed 72 bytes")
)

// ValidatePassword enforces the 8-72 UTF-8 byte password length requirement.
func ValidatePassword(password string) error {
	byteLen := len([]byte(password))
	if byteLen < 8 {
		return ErrPasswordTooShort
	}
	if byteLen > 72 {
		return ErrPasswordTooLong
	}
	return nil
}

// HashPassword hashes a validated password using bcrypt with default cost.
func HashPassword(password string) (string, error) {
	if err := ValidatePassword(password); err != nil {
		return "", err
	}
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// CheckPassword compares a bcrypt hashed password with its possible plaintext equivalent.
func CheckPassword(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// DummyCheckPassword runs bcrypt comparison against a dummy hash for constant-time work.
func DummyCheckPassword(password string) {
	_ = bcrypt.CompareHashAndPassword([]byte(dummyBcryptHash), []byte(password))
}
