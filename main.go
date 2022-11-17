package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/its-me-debk007/Akatsuki_backend/database"
	"github.com/its-me-debk007/Akatsuki_backend/model"
	"github.com/its-me-debk007/Akatsuki_backend/route"
	"github.com/its-me-debk007/Akatsuki_backend/util"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalln("ENV LOADING ERROR", err.Error())
	}

	database.ConnectDatabase()

	app := gin.Default()

	app.Use(gin.Recovery())
	app.Use(middleWare)

	app.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowCredentials: true,
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPatch},
		AllowHeaders:     []string{"Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
	}))

	route.SetupRoutes(app)

	port := os.Getenv("PORT")

	if err := app.Run(":" + port); err != nil {
		log.Fatal("App listen error:-\n" + err.Error())
	}
}

func middleWare(c *gin.Context) {
	url := c.Request.URL

	if !(url.Path[7:12] == "/auth" || url.Path == "/api/v1/posts/random") {
		header := c.GetHeader("Authorization")
		if header == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, model.Message{"no token provided"})
		}

		token := header[7:]

		username, err := util.ParseToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, model.Message{err.Error()})
			return
		}

		c.Request.Header.Set("username", username)
	}
	// c.Next()
}
