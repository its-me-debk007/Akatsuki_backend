package database

import (
	"log"
	"os"

	"github.com/its-me-debk007/Akatsuki_backend/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	dbUrl := os.Getenv("DATABASE_URL")

	db, err := gorm.Open(postgres.Open(dbUrl), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	DB = db
	if err := db.AutoMigrate(
		new(model.User),
		new(model.Otp),
		new(model.Post),
	); err != nil {
		log.Fatalln("AUTO_MIGRATION_ERROR")
	}

	// populateDatabase()
}

// func populateDatabase() {
// 	for i := 0; i < 25; i++ {
// 		post := model.Post{
// 			Media: []string{
// 				"https://res.cloudinary.com/debk007cloud/image/upload/v1668334132/low-resolution-splashes-wallpaper-preview_weaxun.jpg",
// 				"https://res.cloudinary.com/debk007cloud/image/upload/v1668334132/low-resolution-splashes-wallpaper-preview_weaxun.jpg",
// 			},
// 			Description: fmt.Sprintf("%d description", i),
// 			AuthorEmail: fmt.Sprintf("abc%d@gmail.com", i),
// 		}
// 		DB.Save(&post)
// 	}
// }
