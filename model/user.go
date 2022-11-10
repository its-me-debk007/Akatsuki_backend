package model

import "time"

type User struct {
	CreatedAt  time.Time
	Name       string `binding:"required"`
	Email      string `binding:"required"    gorm:"primarykey"`
	Username   string `binding:"required"    gorm:"unique"`
	Password   string `binding:"required"`
	IsVerified bool
	ProfilePic string `binding:"required"`
}
