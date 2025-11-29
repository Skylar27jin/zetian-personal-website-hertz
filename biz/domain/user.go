package domain

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"size:64;not null"`
	Password string `gorm:"not null"`
    Email    string `gorm:"uniqueIndex;size:255"`
	AvatarUrl string `gorm:"type:text"`
}
//note : gorm note only effect autoMigrate, it is not used to validate input

type UserStats struct {
    UserID                int64     `gorm:"primaryKey"`              // 对应 users.id
    FollowersCount        int64     `gorm:"not null;default:0"`      // 粉丝数
    FollowingCount        int64     `gorm:"not null;default:0"`      // 关注数
    PostLikeReceivedCount int64     `gorm:"not null;default:0"`      // TA 发的帖子被点赞总数

	CreatedAt time.Time `gorm:"autoCreateTime"`
    UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

type UserWithStats struct {
    User  *User
    Stats *UserStats
}


type UserProfile struct {
    Id int64
    UserName string
    AvatarUrl string

    FollowersCount int64
    FollowingCount int64
    PostLikeReceivedCount int64

    IsFollowing bool
    FollowedYou bool
    IsMe bool  
}

type SimpleUserProfile struct {
    Id int64
    UserName string
    AvatarUrl string

    IsFollowing bool
    FollowedYou bool
    IsMe bool  
}