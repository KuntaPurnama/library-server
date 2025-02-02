package helpers

import (
	"library-backend/config"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func GenerateToken(username string) (string, error) {
	claims := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.Config.JWTSecret))
}

func ValidateToken(tokenString string) (bool, string) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(config.Config.JWTSecret), nil
	})

	if err != nil || !token.Valid {
		return false, ""
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return false, ""
	}

	username, ok := claims["username"].(string)
	if !ok {
		return false, ""
	}

	return true, username
}
