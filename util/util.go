package util

import (
	"bytes"
	"html/template"
	"log"
	"net/smtp"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

const (
	SMTP_HOST = "smtp.gmail.com"
	SMTP_PORT = "587"
)

func GenerateToken(username string, subject string, expirationTime time.Duration) (string, error) {
	registeredClaims := jwt.RegisteredClaims{
		Issuer:  username,
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

func SendEmail(receiverEmail string, otp int) {
	senderEmail := os.Getenv("SENDER_EMAIL")
	senderPassword := os.Getenv("SENDER_PASSWORD")
	subject := "Subject: Verify your account\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

	auth := smtp.PlainAuth("", senderEmail, senderPassword, SMTP_HOST)

	var t *template.Template
	var err error

	t, err = t.ParseFiles("template/template.html")
	if err != nil {
		log.Fatalln("HTML PARSING ERROR", err.Error())
	}

	buffer := new(bytes.Buffer)

	t.Execute(buffer, gin.H{
		"otp": otp,
	})

	msg := []byte(subject + mime + buffer.String())

	if err = smtp.SendMail(SMTP_HOST + ":" + SMTP_PORT, auth, senderEmail, []string{receiverEmail}, msg); err != nil {
		log.Fatalln("SEND EMAIL ERROR", err.Error())
	}

	log.Println("OTP SENT")
}
