package controller

import (
	"log"
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

	var user model.User

	database.DB.First(&user)

	if user.Email == "" {
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

	accessToken, err := util.GenerateToken(input.Email, "ACCESS", 24)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadGateway, model.Message{err.Error()})
		return
	}

	refreshToken, err := util.GenerateToken(input.Email, "REFRESH", 1)
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
	input.Password = strings.TrimSpace(input.Password)

	if validation := util.IsValidPassword(input.Password); validation != "ok" {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.Message{validation})
		return
	}

	// var user model.User
	// database.DB.First(&user, "email = ?", input.Email)

	// if user.Email != "" {
	// 	c.AbortWithStatusJSON(http.StatusBadRequest, model.Message{"email already registered"})
	// 	return
	// }

	// if user.Username != "" {
	// 	c.AbortWithStatusJSON(http.StatusBadRequest, model.Message{"username already taken"})
	// 	return
	// }

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), 10)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadGateway, model.Message{err.Error()})
		return
	}

	input.Password = string(hashedPassword)
	input.CreatedAt = time.Now()

	if err := database.DB.Create(&input); err.Error != nil {
		var msg string
		log.Println("dfgdbdyy", err.Error.Error()[62:70])
		if err.Error.Error()[61:69] == "username" {
			msg = "username already taken"
		} else {
			msg = "email already registered"
		}

		c.AbortWithStatusJSON(http.StatusBadRequest, model.Message{msg})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "successfully signed up",
	})
}

func VerifyOtp(c *gin.Context) {

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

	if err := database.DB.Model(&model.User{}).Where("email = ?", input.Email).Update("password", input.Password); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error.Error())
		return
	}
}
