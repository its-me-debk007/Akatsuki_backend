package model

import (
	"time"

	"github.com/lib/pq"
)

type Story struct {
	Id          uint64         `gorm:"primary_key; auto_increment"`
	Media       pq.StringArray `gorm:"type:text[]"`
	ExpiresAt   time.Time
	AuthorUsername string `json:"-"`
	Author      User   `gorm:"foreign_key:AuthorUsername"`
}
