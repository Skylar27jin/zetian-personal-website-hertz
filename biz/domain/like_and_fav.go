package domain

import "time"

// PostLike represents a like relation between a user and a post.
type PostLike struct {
	UserID    int64     `json:"user_id" gorm:"primaryKey"`
	PostID    int64     `json:"post_id" gorm:"primaryKey"`
	CreatedAt time.Time `json:"created_at"`
}

// PostFavorite represents a favorite relation between a user and a post.
type PostFavorite struct {
	UserID    int64     `json:"user_id" gorm:"primaryKey"`
	PostID    int64     `json:"post_id" gorm:"primaryKey"`
	CreatedAt time.Time `json:"created_at"`
}
