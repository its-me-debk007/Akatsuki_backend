package model

import "time"

type User struct {
	CreatedAt  time.Time
	Name       string `binding:"required"`
	Email      string `binding:"required"    gorm:"unique"`
	Username   string `binding:"required"    gorm:"primary_key"`
	Password   string `json:",omitempty"    binding:"required"`
	IsVerified bool   `json:"-"`
	ProfilePic string `gorm:"default:https://res.cloudinary.com/debk007cloud/image/upload/v1668334132/low-resolution-splashes-wallpaper-preview_weaxun.jpg"`
}

// func (user *User) FillDefaults() {
// 	if user.ProfilePic == "" {
// 		user.ProfilePic = "https://res.cloudinary.com/debk007cloud/image/upload/v1668334132/low-resolution-splashes-wallpaper-preview_weaxun.jpg"
// 	}
// }
