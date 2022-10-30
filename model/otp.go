package model

import "time"

type Otp struct {
	Email     string `gorm:"primarykey"`
	Otp       int
	CreatedAt time.Time
}
