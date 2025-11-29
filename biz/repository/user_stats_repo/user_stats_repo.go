package user_stats_repo

import (
    "context"

    "gorm.io/gorm"

    "zetian-personal-website-hertz/biz/domain"
    DB "zetian-personal-website-hertz/biz/repository"
)

// CreateEmptyStats 在注册用户后创建一条初始为 0 的统计记录
func CreateEmptyStats(ctx context.Context, userID int64) error {
    us := &domain.UserStats{
        UserID: userID,
        // 其他字段默认 0
    }
    return DB.DB.WithContext(ctx).Create(us).Error
}

// GetStats 按 user_id 读取统计信息
func GetStats(ctx context.Context, userID int64) (*domain.UserStats, error) {
    var us domain.UserStats
    err := DB.DB.WithContext(ctx).
        Where("user_id = ?", userID).
        First(&us).Error
    if err != nil {
        return nil, err
    }
    return &us, nil
}

// （预留：以后可以加 IncrementFollowers / IncrementFollowing / IncrementPostLikeReceived 等）
func IncrementFollowers(ctx context.Context, userID int64, delta int64) error {
    return DB.DB.WithContext(ctx).
        Model(&domain.UserStats{}).
        Where("user_id = ?", userID).
        Update("followers_count", gorm.Expr("followers_count + ?", delta)).
        Error
}

func IncrementFollowing(ctx context.Context, userID int64, delta int64) error {
    return DB.DB.WithContext(ctx).
        Model(&domain.UserStats{}).
        Where("user_id = ?", userID).
        Update("following_count", gorm.Expr("following_count + ?", delta)).
        Error
}

func IncrementPostLikeReceived(ctx context.Context, userID int64, delta int64) error {
    return DB.DB.WithContext(ctx).
        Model(&domain.UserStats{}).
        Where("user_id = ?", userID).
        Update("post_like_received_count", gorm.Expr("post_like_received_count + ?", delta)).
        Error
}
