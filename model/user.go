package model

import "time"

type User struct {
	CreatedAt  time.Time `json:"created_at"`
	Name       string    `json:"name"    binding:"required"`
	Email      string    `json:"email"    binding:"required,email"    gorm:"unique"`
	Username   string    `json:"username"    binding:"required"    gorm:"primary_key"`
	Password   string    `json:"password"    binding:"required"`
	IsVerified bool      `json:"-"`
	ProfilePic string    `json:"profile_pic"    gorm:"default:https://res.cloudinary.com/debk007cloud/image/upload/v1668334132/low-resolution-splashes-wallpaper-preview_weaxun.jpg"`
	Token      string    `json:"-"`
}

// func (user *User) FillDefaults() {
// 	if user.ProfilePic == "" {
// 		user.ProfilePic = "https://res.cloudinary.com/debk007cloud/image/upload/v1668334132/low-resolution-splashes-wallpaper-preview_weaxun.jpg"
// 	}
// }
