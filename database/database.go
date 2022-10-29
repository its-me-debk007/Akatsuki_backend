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
	log.Println(dbUrl)

	db, err := gorm.Open(postgres.Open(dbUrl), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	DB = db
	if err := db.AutoMigrate(new(model.User), new(model.Otp)); err != nil {
		log.Fatalln("AUTO_MIGRATION_ERROR")
	}
}
