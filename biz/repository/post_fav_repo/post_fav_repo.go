package post_fav_repo

import (
	"context"
	"time"

	"zetian-personal-website-hertz/biz/domain"
	DB "zetian-personal-website-hertz/biz/repository"

	"gorm.io/gorm/clause"
)

// FavoritePost inserts a favorite relation. Idempotent.
func FavoritePost(ctx context.Context, userID, postID int64) error {
	fav := &domain.PostFavorite{
		UserID:    userID,
		PostID:    postID,
		CreatedAt: time.Now(),
	}

	return DB.DB.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "user_id"}, {Name: "post_id"}},
			DoNothing: true,
		}).
		Create(fav).Error
}

// UnfavoritePost removes a favorite relation. Idempotent.
func UnfavoritePost(ctx context.Context, userID, postID int64) error {
	return DB.DB.WithContext(ctx).
		Where("user_id = ? AND post_id = ?", userID, postID).
		Delete(&domain.PostFavorite{}).Error
}

// HasUserFavorited checks whether the user has favorited the post.
func HasUserFavorited(ctx context.Context, userID, postID int64) (bool, error) {
	var count int64
	err := DB.DB.WithContext(ctx).
		Model(&domain.PostFavorite{}).
		Where("user_id = ? AND post_id = ?", userID, postID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// CountFavorites returns total fav_count for a post.
func CountFavorites(ctx context.Context, postID int64) (int, error) {
	var count int64
	err := DB.DB.WithContext(ctx).
		Model(&domain.PostFavorite{}).
		Where("post_id = ?", postID).
		Count(&count).Error
	return int(count), err
}

// Optionally: batch count for feed
func CountFavoritesBatch(ctx context.Context, postIDs []int64) (map[int64]int, error) {
	if len(postIDs) == 0 {
		return map[int64]int{}, nil
	}

	type result struct {
		PostID int64
		Count  int64
	}

	var rows []result
	err := DB.DB.WithContext(ctx).
		Model(&domain.PostFavorite{}).
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


// GetUserFavoritedPostIDs returns a set-like map[postID]bool for given user and postIDs.
func GetUserFavoritedPostIDs(ctx context.Context, userID int64, postIDs []int64) (map[int64]bool, error) {
	result := make(map[int64]bool)

	if len(postIDs) == 0 {
		return result, nil
	}

	var rows []domain.PostFavorite
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

// DeleteFavoritesByPostID deletes all favorites for a given post.
func DeleteFavoritesByPostID(ctx context.Context, postID int64) error {
	return DB.DB.WithContext(ctx).
		Where("post_id = ?", postID).
		Delete(&domain.PostFavorite{}).Error
}