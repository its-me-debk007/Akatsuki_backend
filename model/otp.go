package model

import "time"

type Otp struct {
	Email     string    `json:"email" gorm:"primarykey"`
	Otp       int       `json:"otp"`
	CreatedAt time.Time `json:"created_at"`
}