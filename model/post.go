package model

import (
	"time"

	"github.com/lib/pq"
)

type Post struct {
	Id          uint64         `json:"id"    gorm:"primary_key; auto_increment"`
	CreatedAt   time.Time      `json:"created_at"`
	Media       pq.StringArray `json:"media"    gorm:"type:text[]"`
	Description string         `json:"description"`
	Rating      float32        `json:"rating"` // out of 5
	Likes       string         `json:"likes"`
	// Comments
	LikedByUser     bool   `json:"liked_by_user"`
	CommentedByUser bool   `json:"commented_by_user"`
	AuthorUsername  string `json:"-"`
	Author          User   `json:"author" gorm:"foreign_key:AuthorUsername;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
