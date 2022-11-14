package model

import (
	"time"

	"github.com/lib/pq"
)

type Post struct {
	Id          uint64 `gorm:"primary_key; auto_increment"`
	CreatedAt   time.Time
	Media       pq.StringArray `gorm:"type:text[]"`
	Description string
	Rating      float32 // out of 5
	Likes       string
	// Comments
	LikedByUser     bool
	CommentedByUser bool
	AuthorEmail     string `json:"-"`
	Author          User   `gorm:"foreign_key:AuthorEmail"`
}
