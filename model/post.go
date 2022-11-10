package model

import (
	"time"

	"github.com/lib/pq"
)

type Post struct {
	Id          uint64 `gorm:"primaryKey;autoIncrement"`
	CreatedAt   time.Time
	Media       pq.StringArray
	Description string
	Likes       string
	// Comments
	LikedByUser     bool
	CommentedByUser bool
	// AuthorID        uint64
	Author          User `gorm:"foreignKey:Email"`
}
