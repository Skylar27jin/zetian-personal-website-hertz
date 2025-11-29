package user_follow_repo

import (
	"context"

	"zetian-personal-website-hertz/biz/domain"
	DB "zetian-personal-website-hertz/biz/repository"
)

// Follow 插入一条关注关系（不做幂等判断）
// 幂等逻辑在 service 里面用 IsFollowing 控制
func Follow(ctx context.Context, followerID, followeeID int64) error {
	fr := &domain.UserFollowRecord{
		FollowerID: followerID,
		FolloweeID: followeeID,
	}
	return DB.DB.WithContext(ctx).Create(fr).Error
}

// Unfollow 删除一条关注关系，删除 0 行也视为成功
func Unfollow(ctx context.Context, followerID, followeeID int64) error {
	return DB.DB.WithContext(ctx).
		Where("follower_id = ? AND followee_id = ?", followerID, followeeID).
		Delete(&domain.UserFollowRecord{}).Error
}

// IsFollowing 判断 follower 是否已经关注了 followee
func IsFollowing(ctx context.Context, followerID, followeeID int64) (bool, error) {
	var count int64
	if err := DB.DB.WithContext(ctx).
		Model(&domain.UserFollowRecord{}).
		Where("follower_id = ? AND followee_id = ?", followerID, followeeID).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
