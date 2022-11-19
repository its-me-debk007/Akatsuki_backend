package route

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/its-me-debk007/Akatsuki_backend/controller"
	"github.com/its-me-debk007/Akatsuki_backend/model"
)

func SetupRoutes(app *gin.Engine) {
	api := app.Group("/api/v1")
	{
		authGroup := api.Group("/auth")
		{
			authGroup.POST("/login", controller.Login)
			authGroup.POST("/signup", controller.Signup)
			authGroup.POST("/send-otp", controller.SendOtp)
			authGroup.POST("/verify", controller.VerifyOtp)
			authGroup.POST("/reset", controller.ResetPassword)
		}

		postGroup := api.Group("/posts")
		{
			postGroup.POST("/create", controller.CreatePost)
			postGroup.GET("/random", controller.RandomPosts)
			postGroup.POST("/like", controller.LikePost)
		}

		storyGroup := api.Group("/stories")
		{
			storyGroup.POST("/create", controller.CreateStory)
			storyGroup.GET("/", controller.GetStories)
		}

		userGroup := api.Group("/users")
		{
			userGroup.GET("/follow", controller.Follow)
			userGroup.GET("/profile", controller.Profile)
			userGroup.GET("/suggestions", controller.Suggestion)
		}

		api.GET("/search", controller.Search)

		app.NoRoute(func(c *gin.Context) {
			c.JSON(http.StatusNotFound, model.Message{"route doesn't exist"})
		})

		app.NoMethod(func(c *gin.Context) {
			c.JSON(http.StatusMethodNotAllowed, model.Message{"method not allowed"})
		})
	}
}
