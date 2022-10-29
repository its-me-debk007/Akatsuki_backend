package model

import "time"

type User struct {
	CreatedAt  time.Time `json:"created_at" `
	Name       string    `json:"name"    binding:"required"`
	Email      string    `json:"email"    binding:"required"    gorm:"primarykey"`
	Username   string    `json:"username"    binding:"required"    gorm:"unique"`
	Password   string    `json:"password"    binding:"required"`
	IsVerified bool      `json:"is_verified"`
}
