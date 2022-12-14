package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/its-me-debk007/Akatsuki_backend/database"
	"github.com/its-me-debk007/Akatsuki_backend/model"
)

func Follow(c *gin.Context) {
	username := c.Query("username")
	if db := database.DB.First(&model.User{}, "username = ?", username); db.Error != nil {
		c.JSON(http.StatusBadRequest, model.Message{"username invalid"})
		return
	}
}

func Search(c *gin.Context) {
	query := c.Query("query")
	query = fmt.Sprintf("%%%s%%", query)

	var users []struct {
		Username   string `json:"username"`
		Name       string `json:"name"`
		ProfilePic string `json:"profile_pic"`
	}
	database.DB.Model(&model.User{}).Where("username LIKE ?", query).Or("name LIKE ?", query).Find(&users)

	var posts []model.Post
	database.DB.Model(&model.Post{}).Find(&posts, "description LIKE ?", query)

	for i := range posts {
		database.DB.First(&posts[i].Author, "username = ?", posts[i].AuthorUsername)
	}

	c.JSON(http.StatusOK, gin.H{
		"users": users,
		"posts": posts,
	})
}

func Profile(c *gin.Context) {
	username := c.Query("username")
	if username == "" {
		username = c.GetHeader("username")
	}

	var user model.User
	if db := database.DB.First(&user, "username = ?", username); db.Error != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.Message{"username doesn't exist"})
		return
	}

	var posts []model.Post
	database.DB.Find(&posts, "author_username = ?", username)

	c.JSON(http.StatusOK, gin.H{
		"about": user,
		"posts": posts,
	})
}

func Suggestion(c *gin.Context) {
	var users []model.User
	database.DB.Raw("SELECT * FROM users ORDER BY RANDOM() LIMIT 5").Scan(&users)

	c.JSON(http.StatusOK, users)
}
