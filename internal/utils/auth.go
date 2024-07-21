package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// hashes password into a hashed slice of bytes
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// authenticates user's password by comparing hashes
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
