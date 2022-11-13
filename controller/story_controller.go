package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/its-me-debk007/Akatsuki_backend/model"
)

func CreateStory(c *gin.Context) {


	c.JSON(http.StatusOK, model.Message{"post created successfully"})
}