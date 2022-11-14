package controller

import (
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
