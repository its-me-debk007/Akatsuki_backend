package controller

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/its-me-debk007/Akatsuki_backend/database"
	"github.com/its-me-debk007/Akatsuki_backend/model"
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

	c.JSON(http.StatusOK, gin.H{"data": input})
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
	// input.CreatedAt = time.Now()

	var user model.User
	database.DB.First(&user, "username = ?", input.Username)

	if user.Email != "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "email already registered",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "successfully signed up",
	})
}
