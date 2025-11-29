package domain

import "time"

type UserFollowRecord struct {
	ID         uint      `gorm:"primaryKey"`
	FollowerID int64     `gorm:"not null;index:idx_follower_followee,unique"`
	FolloweeID int64     `gorm:"not null;index:idx_follower_followee,unique"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
}

func (UserFollowRecord) TableName() string {
	return "user_follow_records"
}
