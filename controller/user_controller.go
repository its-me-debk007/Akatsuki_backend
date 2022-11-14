package controller

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/its-me-debk007/Akatsuki_backend/database"
	"github.com/its-me-debk007/Akatsuki_backend/model"
	"github.com/its-me-debk007/Akatsuki_backend/util"
)

func Follow(c *gin.Context) {
	token := c.GetHeader("Authorization")[7:]
	userEmail, err := util.ParseToken(token)
	log.Println(userEmail)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, model.Message{err.Error()})
		return
	}

	username := c.Query("username")
	if db := database.DB.First(&model.User{}, "username = ?", username); db.Error != nil {
		c.JSON(http.StatusBadRequest, model.Message{"username invalid"})
		return
	}
}

func Search(c *gin.Context) {
	token := c.GetHeader("Authorization")[7:]
	_, err := util.ParseToken(token)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, model.Message{err.Error()})
		return
	}

	query := c.Query("query")
	query = fmt.Sprintf("%%%s%%", query)

	var users []struct {
		Username   string `json:"username"`
		Name       string `json:"name"`
		ProfilePic string `json:"profile_pic"`
	}
	database.DB.Model(&model.User{}).Find(&users, "username LIKE ?", query)

	var posts []model.Post
	database.DB.Model(&model.Post{}).Find(&posts, "description LIKE ?", query)

	c.JSON(http.StatusOK, gin.H{
		"users": users,
		"posts": posts,
	})
}

func Profile(c *gin.Context) {
	token := c.GetHeader("Authorization")[7:]
	_, err := util.ParseToken(token)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, model.Message{err.Error()})
		return
	}

	username := c.Query("username")

	var user model.User
	if db := database.DB.First(&user, "username = ?", username); db.Error != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.Message{"username doesn't exist"})
		return
	}

	var posts []model.Post
	database.DB.Find(&posts, "authorusername = ?", username)

	c.JSON(http.StatusOK, gin.H{
		"about": user,
		"posts": posts,
	})
}
