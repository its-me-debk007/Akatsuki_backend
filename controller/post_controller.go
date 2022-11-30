package controller

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/its-me-debk007/Akatsuki_backend/database"
	"github.com/its-me-debk007/Akatsuki_backend/model"
	"github.com/its-me-debk007/Akatsuki_backend/util"
)

func CreatePost(c *gin.Context) {
	username := c.GetHeader("username")

	form, err := c.MultipartForm()

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.Message{err.Error()})
		return
	}

	media := form.File["media"]
	description := form.Value["description"][0]

	var postCollection []string

	for _, fileHeader := range media {

		file, err := fileHeader.Open()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, model.Message{err.Error()})
			return
		}

		url, err := util.UploadMedia(file, time.Now())

		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, model.Message{err.Error()})
			return
		}

		postCollection = append(postCollection, url)
	}

	post := model.Post{
		Description:    description,
		Media:          postCollection,
		AuthorUsername: username,
	}

	if db := database.DB.Save(&post); db.Error != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.Message{db.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, model.Message{"post created succesfully"})
}

func RandomPosts(c *gin.Context) {
	var posts []model.Post

	database.DB.Raw("SELECT * FROM posts ORDER BY RANDOM() LIMIT 5 ;").Scan(&posts)

	for i, post := range posts {

		var author model.User
		database.DB.First(&author, "username = ?", post.AuthorUsername)

		author.Password = ""
		post.Author = author

		posts[i] = post
	}

	c.JSON(http.StatusOK, posts)
}

func LikePost(c *gin.Context) {

}
