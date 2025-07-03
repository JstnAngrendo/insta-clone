package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(tokenID string, userID uint, exp time.Time) (string, error) {
	claims := jwt.MapClaims{
		"token_id": tokenID,
		"user_id":  userID,
		"exp":      exp.Unix(),
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString([]byte(os.Getenv("JWT_SECRET")))
}
