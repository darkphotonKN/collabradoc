package auth

import (
	"golang.org/x/crypto/bcrypt"
	"time"
)

const (
	SecretKey   = "your_secret_key" // Change this to a secure key
	TokenExpiry = 24 * time.Hour
)

// -- PASSWORDS --

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

// -- JWT --
func GenerateJWT() (string, error) {

}
