package auth

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

const (
	TokenExpiry = 24 * time.Hour
)

type Claims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

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

func GenerateJWT(userId uint) (string, error) {
	envKey := os.Getenv("JWT_SECRET_KEY")
	jwtKey := []byte(envKey)

	// expires in one day
	expirationTime := time.Now().Add(time.Hour * 24)

	claims := &Claims{
		UserID: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// middleware that protects with jwt tokens
func JWTMiddleware(next http.Handler) http.Handler {
	var envKey = os.Getenv("JWT_SECRET_KEY")
	var jwtKey = []byte(envKey)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims := &Claims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		// return unauthorized error if jwt validation failed
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				http.Error(w, "Invalid token signature", http.StatusUnauthorized)
				return
			}
			http.Error(w, fmt.Sprintf("Error when parsing token: %s", err), http.StatusBadRequest)
			return
		}

		if !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Add the user ID to the request context
		ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
