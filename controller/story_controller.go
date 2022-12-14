package controller

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/its-me-debk007/Akatsuki_backend/database"
	"github.com/its-me-debk007/Akatsuki_backend/model"
	"github.com/its-me-debk007/Akatsuki_backend/util"
)

func CreateStory(c *gin.Context) {
	username := c.GetHeader("username")

	form, err := c.MultipartForm()

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.Message{err.Error()})
		return
	}

	media := form.File["media"]

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

	story := model.Story{
		Media:          postCollection,
		AuthorUsername: username,
		ExpiresAt:      time.Now().Add(time.Hour * 24),
	}

	if db := database.DB.Save(&story); db.Error != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.Message{db.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, model.Message{"story created successfully"})
}

func GetStories(c *gin.Context) {

}
