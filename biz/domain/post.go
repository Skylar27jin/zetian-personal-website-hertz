package domain

import "time"

// Post 表示一条帖子
type Post struct {
	ID        int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID    int64     `json:"user_id"`                        // 发帖人ID
	SchoolID  int64     `json:"school_id"`                      // 外键 -> School
	Title     string    `json:"title" gorm:"type:varchar(255)"` // 帖子标题
	Content   string    `json:"content" gorm:"type:text"`       // 帖子正文
	LikeCount int       `json:"like_count"`
	FavCount  int       `json:"fav_count"`
	ViewCount int       `json:"view_count"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
