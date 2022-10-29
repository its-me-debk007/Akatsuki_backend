package util

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func GenerateToken(username string, subject string, expirationTime time.Duration) (string, error) {
	registeredClaims := jwt.RegisteredClaims{
		Issuer: username,
		Subject: subject,
		ExpiresAt: &jwt.NumericDate{
			Time: time.Now().Add(time.Hour * expirationTime),
		},
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, registeredClaims)

	secretKey := os.Getenv("SECRET_KEY")

	token, err := claims.SignedString([]byte(secretKey))

	if err != nil {
		return token, err
	}

	return token, nil
}

func IsValidPassword(password string) string {
	isDigit, isLowercase, isUppercase, isSpecialChar := 0, 0, 0, 0
	for _, ch := range password {
		switch {
		case ch >= '0' && ch <= '9':
			isDigit = 1

		case ch >= 'a' && ch <= 'z':
			isLowercase = 1

		case ch >= 'A' && ch <= 'Z':
			isUppercase = 1

		case ch == '$' || ch == '!' || ch == '@' || ch == '#' || ch == '%' || ch == '&' || ch == '^' || ch == '*' || ch == '/' || ch == '\\':
			isSpecialChar = 1
		}
	}

	switch {
	case len(password) < 8:
		return "password must be at least 8 characters long"

	case isDigit == 0:
		return "password must contain at-least one numeric digit"

	case isLowercase == 0:
		return "password must contain at-least one lowercase alphabet"

	case isUppercase == 0:
		return "password must contain at-least one uppercase alphabet"

	case isSpecialChar == 0:
		return "password must contain at-least one special character"

	default:
		return "ok"
	}
}

func VerifyToken(tokenString string) (*jwt.RegisteredClaims, error) {
	secretKey := os.Getenv("SECRET_KEY")

	registeredClaims := jwt.RegisteredClaims{}

	_, err := jwt.ParseWithClaims(tokenString, &registeredClaims, func(t *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	return &registeredClaims, nil
}
