package util

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"html/template"
	"log"
	"mime/multipart"
	"net/smtp"
	"os"
	"time"

	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/uploader"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/its-me-debk007/Akatsuki_backend/database"
	"github.com/its-me-debk007/Akatsuki_backend/model"
)

const (
	SMTP_HOST       = "smtp.gmail.com"
	SMTP_PORT       = "587"
	CharCollections = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890-._"
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

func ParseToken(tokenString string, typeShouldBeAccess bool) (string, error) {
	secretKey := os.Getenv("SECRET_KEY")

	registeredClaims := jwt.RegisteredClaims{}

	_, err := jwt.ParseWithClaims(tokenString, &registeredClaims, func(t *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	var errorMsg string
	db := database.DB.First(&model.User{}, "username = ?", registeredClaims.Issuer)

	switch {
	case err != nil:
		errorMsg = "invalid token"

	case typeShouldBeAccess && registeredClaims.Subject != "ACCESS":
		errorMsg = "invalid token (required type is access token)"

	case !typeShouldBeAccess && registeredClaims.Subject == "ACCESS":
		errorMsg = "invalid token (required type is refresh token)"

	case db.Error != nil:
		errorMsg = "user not signed up"

	case time.Since(registeredClaims.ExpiresAt.Time) >= 0:
		errorMsg = "token expired"
	}

	if errorMsg != "" {
		return "", errors.New(errorMsg)
	}

	return registeredClaims.Issuer, nil
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

func SendEmail(receiverEmail string, link string) {

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
		"link": link,
	})

	msg := []byte(subject + mime + buffer.String())

	if err = smtp.SendMail(SMTP_HOST+":"+SMTP_PORT, auth, senderEmail, []string{receiverEmail}, msg); err != nil {
		log.Fatalln("SEND EMAIL ERROR", err.Error())
	}

	log.Printf("LINK SENT TO %s", receiverEmail)
}

func UploadMedia(file multipart.File, id time.Time) (string, error) {
	cldCloudName, cldApiKey, cldApiSecret := os.Getenv("CLOUDINARY_CLOUD_NAME"), os.Getenv("CLOUDINARY_API_KEY"), os.Getenv("CLOUDINARY_API_SECRET")

	cld, _ := cloudinary.NewFromParams(
		cldCloudName,
		cldApiKey,
		cldApiSecret,
	)
	ctx := context.Background()

	resp, err := cld.Upload.Upload(
		ctx,
		file,
		uploader.UploadParams{
			PublicID: fmt.Sprintf("docs/sdk/go/akatsuki_post_%v", id),
		})

	if err != nil {
		return "", fmt.Errorf("CLOUDINARY_ERROR:- %s", err.Error())
	}

	return resp.SecureURL, nil
}
