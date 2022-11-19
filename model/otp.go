package model

import (
	"math/big"
	"time"
)

type Otp struct {
	Email     string `gorm:"primarykey"`
	Otp       *big.Int
	CreatedAt time.Time
}
