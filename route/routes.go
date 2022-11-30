package route

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/its-me-debk007/Akatsuki_backend/controller"
	"github.com/its-me-debk007/Akatsuki_backend/model"
	"github.com/its-me-debk007/Akatsuki_backend/util"
)

func SetupRoutes(app *gin.Engine) {
	api := app.Group("/api/v1")
	{
		authGroup := api.Group("/auth")
		{
			authGroup.POST("/login", controller.Login)
			authGroup.POST("/signup", controller.Signup)
			authGroup.POST("/forgot-password", controller.ForgotPassword)
			authGroup.POST("/verify", controller.Verify)
			authGroup.POST("/reset-password", controller.ResetPassword)
		}

		postGroup := api.Group("/posts")
		{
			postGroup.POST("/create", controller.CreatePost, middleWare)
			postGroup.GET("/random", controller.RandomPosts)
			postGroup.POST("/like", controller.LikePost, middleWare)
		}

		storyGroup := api.Group("/stories", middleWare)
		{
			storyGroup.POST("/create", controller.CreateStory)
			storyGroup.GET("/", controller.GetStories)
		}

		userGroup := api.Group("/users", middleWare)
		{
			userGroup.GET("/follow", controller.Follow)
			userGroup.GET("/profile", controller.Profile)
			userGroup.GET("/suggestions", controller.Suggestion)
		}

		api.GET("/search", middleWare, controller.Search)

		app.NoRoute(func(c *gin.Context) {
			c.JSON(http.StatusNotFound, model.Message{"route doesn't exist"})
		})

		app.NoMethod(func(c *gin.Context) {
			c.JSON(http.StatusMethodNotAllowed, model.Message{"method not allowed"})
		})
	}
}

func middleWare(c *gin.Context) {
	header := c.GetHeader("Authorization")
	if header == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, model.Message{"no token provided"})
	}

	token := header[7:]

	username, err := util.ParseToken(token, true)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, model.Message{err.Error()})
		return
	}

	c.Request.Header.Set("username", username)
}
