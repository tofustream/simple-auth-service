package security

import (
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// トークンのデフォルト有効期限（秒）
const (
	DefaultAccessTokenExpiry  = 15 * 60          // 15分（900秒）
	DefaultRefreshTokenExpiry = 7 * 24 * 60 * 60 // 7日（604800秒）
)

func GenerateJWT(userID uint, expirySeconds int) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(time.Duration(expirySeconds) * time.Second).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func ValidateJWT(tokenString string) (*jwt.Token, error) {
	secret := os.Getenv("JWT_SECRET")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		log.Println("Failed to validate token:", err)
	}
	return token, err
}
