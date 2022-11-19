package controller

import (
	"crypto/rand"
	"math/big"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/its-me-debk007/Akatsuki_backend/database"
	"github.com/its-me-debk007/Akatsuki_backend/model"
	"github.com/its-me-debk007/Akatsuki_backend/util"
	"golang.org/x/crypto/bcrypt"
)

func Login(c *gin.Context) {
	input := new(struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	})

	if err := c.ShouldBindJSON(input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	input.Email = strings.TrimSpace(input.Email)
	input.Email = strings.ToLower(input.Email)
	input.Password = strings.TrimSpace(input.Password)

	var filteredPassword string
	for _, ch := range input.Password {
		if ch != '_' {
			filteredPassword += string(ch)
		}
	}
	
	input.Password = filteredPassword

	var user model.User

	if db := database.DB.First(&user, "email = ?", input.Email); db.Error != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.Message{"no account found with given credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.Message{"invalid password"})
		return
	}

	if !user.IsVerified {
		c.AbortWithStatusJSON(http.StatusUnauthorized, model.Message{"user not verified"})
		return
	}

	accessToken, err := util.GenerateToken(user.Username, "ACCESS", 24*7)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadGateway, model.Message{err.Error()})
		return
	}

	refreshToken, err := util.GenerateToken(user.Username, "REFRESH", 1)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadGateway, model.Message{err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func Signup(c *gin.Context) {
	input := model.User{}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	input.Email = strings.TrimSpace(input.Email)
	input.Email = strings.ToLower(input.Email)
	input.Username = strings.TrimSpace(input.Username)
	input.Name = strings.TrimSpace(input.Name)
	input.Password = strings.TrimSpace(input.Password)

	if validation := util.IsValidPassword(input.Password); validation != "ok" {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.Message{validation})
		return
	}

	passwordLength := int64(len(input.Password))
	firstRandomPosition, _ := rand.Int(rand.Reader, big.NewInt(passwordLength))
	input.Password = input.Password[:firstRandomPosition.Int64()] + "_" + input.Password[firstRandomPosition.Int64():]

	passwordLength = int64(len(input.Password))
	secondRandomPosition, _ := rand.Int(rand.Reader, big.NewInt(passwordLength))
	input.Password = input.Password[:secondRandomPosition.Int64()] + "_" + input.Password[secondRandomPosition.Int64():]

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), 10)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadGateway, model.Message{err.Error()})
		return
	}

	input.Password = string(hashedPassword)

	if err := database.DB.Create(&input); err.Error != nil {
		var msg string

		if err.Error.Error()[61:66] == "email" {
			msg = "email already registered"
		} else {
			msg = "username already taken"
		}

		c.AbortWithStatusJSON(http.StatusBadRequest, model.Message{msg})
		return
	}

	c.JSON(http.StatusOK, model.Message{"successfully signed up"})
}

func VerifyOtp(c *gin.Context) {
	input := new(struct {
		Email string `json:"email"    binding:"required,email"`
		Otp   int    `json:"otp"    binding:"required"`
	})

	if err := c.ShouldBindJSON(input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.Message{err.Error()})
		return
	}

	input.Email = strings.TrimSpace(input.Email)
	input.Email = strings.ToLower(input.Email)

	otpStruct := model.Otp{}

	database.DB.First(&otpStruct, "email = ?", input.Email)

	if otpStruct.Email == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.Message{"otp not generated yet"})
		return
	}

	if timeDiff := time.Now().Sub(otpStruct.CreatedAt); timeDiff > (time.Minute * 5) {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.Message{"otp expired"})
		return
	}

	if otpStruct.Otp != input.Otp {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.Message{"otp incorrect"})
		return
	}

	database.DB.Model(&model.User{}).Where("email = ?", input.Email).Update("is_verified", true)

	c.JSON(http.StatusOK, model.Message{"otp verified successfully"})
}

func ResetPassword(c *gin.Context) {
	input := new(struct {
		Email    string `json:"email"    binding:"required,email"`
		Password string `json:"new_password"    binding:"required"`
	})

	if err := c.ShouldBindJSON(input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	input.Email = strings.TrimSpace(input.Email)
	input.Email = strings.ToLower(input.Email)
	input.Password = strings.TrimSpace(input.Password)

	if validation := util.IsValidPassword(input.Password); validation != "ok" {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.Message{validation})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), 10)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadGateway, model.Message{err.Error()})
		return
	}

	input.Password = string(hashedPassword)

	database.DB.Model(&model.User{}).Where("email = ?", input.Email).Update("password", input.Password)

	c.JSON(http.StatusOK, model.Message{"successfully changed password"})
}

func SendOtp(c *gin.Context) {
	input := new(struct {
		Email string `json:"email"    binding:"required,email"`
	})

	if err := c.ShouldBindJSON(input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.Message{err.Error()})
		return
	}

	input.Email = strings.TrimSpace(input.Email)
	input.Email = strings.ToLower(input.Email)

	var user model.User

	database.DB.First(&user, "email = ?", input.Email)

	if user.Email == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.Message{"user not registered"})
		return
	}

	otp, _ := rand.Int(rand.Reader, big.NewInt(900000))
	otp.Add(otp, big.NewInt(100000))

	go util.SendEmail(input.Email, otp)

	otpStruct := model.Otp{
		Email:     input.Email,
		Otp:       otp,
		CreatedAt: time.Now(),
	}

	if db := database.DB.Save(&otpStruct); db.Error != nil {
		c.AbortWithStatusJSON(http.StatusBadGateway, model.Message{db.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, model.Message{"otp sent successfully"})
}
