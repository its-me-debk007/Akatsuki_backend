package route

import (
	"github.com/gin-gonic/gin"
	"github.com/its-me-debk007/Akatsuki_backend/controller"
)

func SetupRoutes(app *gin.Engine) {
	authGroup := app.Group("/auth")

	authGroup.POST("/login", controller.Login)
	authGroup.POST("/signup", controller.Signup)
	authGroup.POST("/send_otp", controller.SendOtp)
	authGroup.POST("/verify", controller.VerifyOtp)
	authGroup.POST("/reset", controller.ResetPassword)

	postGroup := app.Group("/post")

	postGroup.POST("/create", controller.CreatePost)
	postGroup.GET("/random", controller.RandomPosts)
	postGroup.POST("/like", controller.LikePost)

	storyGroup := app.Group("/story")

	storyGroup.POST("/create", controller.CreateStory)
	storyGroup.GET("/", controller.GetStories)

	userGroup := app.Group("/user")

	userGroup.GET("/follow", controller.Follow)
	userGroup.GET("/profile", controller.Profile)
	userGroup.GET("/suggestion", controller.Suggestion)

	app.GET("/search", controller.Search)
}
