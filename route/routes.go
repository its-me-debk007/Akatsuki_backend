package route

import (
	"github.com/gin-gonic/gin"
	"github.com/its-me-debk007/Akatsuki_backend/controller"
)

func SetupRoutes(app *gin.Engine) {
	app.POST("/login", controller.Login)
	app.POST("/signup", controller.Signup)
	app.POST("/send_otp", controller.SendOtp)
	app.POST("/verify", controller.VerifyOtp)
	app.POST("/reset", controller.ResetPassword)

	app.POST("/post/create", controller.CreatePost)
	app.GET("/post/random", controller.RandomPosts)

	app.POST("/story/create", controller.CreateStory)
	app.GET("/story", controller.GetStories)

	app.GET("/user/follow", controller.Follow)

	app.GET("/api/search", controller.Search)
	app.GET("/user/profile", controller.Profile)
}
