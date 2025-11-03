package post_like_repo

import (
	"context"
	"time"

	"zetian-personal-website-hertz/biz/domain"
	DB "zetian-personal-website-hertz/biz/repository"

	"gorm.io/gorm/clause"
)

// LikePost inserts a like relation. Idempotent: multiple calls are safe.
func LikePost(ctx context.Context, userID, postID int64) error {
	like := &domain.PostLike{
		UserID:    userID,
		PostID:    postID,
		CreatedAt: time.Now(),
	}

	// ON CONFLICT(user_id, post_id) DO NOTHING
	return DB.DB.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "user_id"}, {Name: "post_id"}},
			DoNothing: true,
		}).
		Create(like).Error
}

// UnlikePost removes a like relation. Idempotent: deleting non-existing row is ok.
func UnlikePost(ctx context.Context, userID, postID int64) error {
	return DB.DB.WithContext(ctx).
		Where("user_id = ? AND post_id = ?", userID, postID).
		Delete(&domain.PostLike{}).Error
}

// HasUserLiked checks whether the user has liked the post.
func HasUserLiked(ctx context.Context, userID, postID int64) (bool, error) {
	var count int64
	err := DB.DB.WithContext(ctx).
		Model(&domain.PostLike{}).
		Where("user_id = ? AND post_id = ?", userID, postID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// CountLikes returns total like_count for a post.
func CountLikes(ctx context.Context, postID int64) (int, error) {
	var count int64
	err := DB.DB.WithContext(ctx).
		Model(&domain.PostLike{}).
		Where("post_id = ?", postID).
		Count(&count).Error
	return int(count), err
}

// Optionally: count likes for a batch of posts (for feed)
func CountLikesBatch(ctx context.Context, postIDs []int64) (map[int64]int, error) {
	if len(postIDs) == 0 {
		return map[int64]int{}, nil
	}

	type result struct {
		PostID int64
		Count  int64
	}

	var rows []result
	err := DB.DB.WithContext(ctx).
		Model(&domain.PostLike{}).
		Select("post_id, COUNT(*) as count").
		Where("post_id IN ?", postIDs).
		Group("post_id").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	m := make(map[int64]int, len(rows))
	for _, r := range rows {
		m[r.PostID] = int(r.Count)
	}
	return m, nil
}

// GetUserLikedPostIDs returns a set-like map[postID]bool for given user and postIDs.
func GetUserLikedPostIDs(ctx context.Context, userID int64, postIDs []int64) (map[int64]bool, error) {
	result := make(map[int64]bool)

	if len(postIDs) == 0 {
		return result, nil
	}

	var rows []domain.PostLike
	err := DB.DB.WithContext(ctx).
		Where("user_id = ? AND post_id IN ?", userID, postIDs).
		Find(&rows).Error
	if err != nil {
		return nil, err
	}

	for _, r := range rows {
		result[r.PostID] = true
	}
	return result, nil
}
