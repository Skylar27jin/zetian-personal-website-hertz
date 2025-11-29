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


// BatchIsFollowing
// viewerID 是否关注了 targets 中每个人（viewer -> target）
func BatchIsFollowing(ctx context.Context, viewerID int64, targetIDs []int64) (map[int64]bool, error) {
	result := make(map[int64]bool, len(targetIDs))

	// 空列表直接返回空 map
	if len(targetIDs) == 0 {
		return result, nil
	}

	var followedIDs []int64
	if err := DB.DB.WithContext(ctx).
		Model(&domain.UserFollowRecord{}).
		Where("follower_id = ? AND followee_id IN ?", viewerID, targetIDs).
		Pluck("followee_id", &followedIDs).Error; err != nil {
		return nil, err
	}

	for _, id := range followedIDs {
		result[id] = true
	}
	// 没出现在 map 里的，默认就是 false
	return result, nil
}

// BatchIsFollowedBy
// targets 中每个人是否关注了 viewerID（target -> viewer）
func BatchIsFollowedBy(ctx context.Context, viewerID int64, targetIDs []int64) (map[int64]bool, error) {
	result := make(map[int64]bool, len(targetIDs))

	if len(targetIDs) == 0 {
		return result, nil
	}

	var followerIDs []int64
	if err := DB.DB.WithContext(ctx).
		Model(&domain.UserFollowRecord{}).
		Where("followee_id = ? AND follower_id IN ?", viewerID, targetIDs).
		Pluck("follower_id", &followerIDs).Error; err != nil {
		return nil, err
	}

	for _, id := range followerIDs {
		result[id] = true
	}
	return result, nil
}


func ListFollowers(
	ctx context.Context,
	targetUserID int64,
	cursor int64,
	limit int,
) (records []domain.UserFollowRecord, nextCursor int64, hasMore bool, err error) {
	if limit <= 0 {
		limit = 20
	}

	q := DB.DB.WithContext(ctx).
		Where("followee_id = ?", targetUserID)

	if cursor > 0 {
		q = q.Where("id > ?", cursor)
	}

	// 多取 1 条判断 hasMore
	var list []domain.UserFollowRecord
	if err = q.Order("id ASC").Limit(limit + 1).Find(&list).Error; err != nil {
		return
	}

	if len(list) == 0 {
		return nil, 0, false, nil
	}

	if len(list) > limit {
		records = list[:limit]
		hasMore = true
		nextCursor = int64(list[limit-1].ID)
	} else {
		records = list
		hasMore = false
		nextCursor = 0
	}
	return
}

func ListFollowees(
	ctx context.Context,
	targetUserID int64,
	cursor int64,
	limit int,
) (records []domain.UserFollowRecord, nextCursor int64, hasMore bool, err error) {
	if limit <= 0 {
		limit = 20
	}

	q := DB.DB.WithContext(ctx).
		Where("follower_id = ?", targetUserID)

	if cursor > 0 {
		q = q.Where("id > ?", cursor)
	}

	var list []domain.UserFollowRecord
	if err = q.Order("id ASC").Limit(limit + 1).Find(&list).Error; err != nil {
		return
	}

	if len(list) == 0 {
		return nil, 0, false, nil
	}

	if len(list) > limit {
		records = list[:limit]
		hasMore = true
		nextCursor = int64(list[limit-1].ID)
	} else {
		records = list
		hasMore = false
		nextCursor = 0
	}
	return
}