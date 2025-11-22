package post_repo

import (
	"context"
	"time"

	"zetian-personal-website-hertz/biz/domain"
	DB "zetian-personal-website-hertz/biz/repository"

	"gorm.io/gorm"
)

/*
PostBaseRepo
------------
This repo handles ONLY the post_base table.
Stats operations should go to PostStatsRepo.
*/

// -----------------------------------------------------------------------------
// Create
// -----------------------------------------------------------------------------

// CreatePostBase inserts a new PostBase record.
// Note: stats row must be filled separately (via PostStatsRepo.CreateEmpty).
func CreatePostBase(ctx context.Context, base *domain.PostBase) error {
	return DB.DB.WithContext(ctx).Create(base).Error
}

// -----------------------------------------------------------------------------
// Query
// -----------------------------------------------------------------------------

// GetPostBaseByID returns the base post row.
// No stats included.
func GetPostBaseByID(ctx context.Context, id int64) (*domain.PostBase, error) {
	var base domain.PostBase
	err := DB.DB.WithContext(ctx).First(&base, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &base, nil
}

// -----------------------------------------------------------------------------
// Update
// -----------------------------------------------------------------------------

// UpdatePostBase ensures only the owner can update title/content.
// UpdatedAt will auto-update via gorm hook.
func UpdatePostBase(ctx context.Context, userID, postID int64, title, content string) error {
	return DB.DB.WithContext(ctx).
		Model(&domain.PostBase{}).
		Where("id = ? AND user_id = ?", postID, userID).
		Updates(map[string]any{
			"title":   title,
			"content": content,
		}).Error
}

// -----------------------------------------------------------------------------
// Delete
// -----------------------------------------------------------------------------

// DeletePostBase deletes a post only if owner matches.
// Note: post_stats row will be auto-deleted (ON DELETE CASCADE).
func DeletePostBase(ctx context.Context, userID, postID int64) error {
	tx := DB.DB.WithContext(ctx).
		Where("id = ? AND user_id = ?", postID, userID).
		Delete(&domain.PostBase{})

	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// -----------------------------------------------------------------------------
// List (for index / personal page / school page)
// -----------------------------------------------------------------------------

// ListRecentPosts fetches most recent posts (base only).
func ListRecentPosts(ctx context.Context, limit int) ([]domain.PostBase, error) {
	var posts []domain.PostBase
	err := DB.DB.WithContext(ctx).
		Order("created_at DESC").
		Limit(limit).
		Find(&posts).Error
	return posts, err
}

// ListPostsBySchoolIDBefore paginates by created_at for school feed.
func ListPostsBySchoolIDBefore(ctx context.Context, schoolID int64, before time.Time, limit int) ([]domain.PostBase, error) {
	var posts []domain.PostBase
	err := DB.DB.WithContext(ctx).
		Where("school_id = ? AND created_at < ?", schoolID, before).
		Order("created_at DESC").
		Limit(limit).
		Find(&posts).Error
	return posts, err
}

// ListPostsByUserIDBefore paginates user’s own posts.
func ListPostsByUserIDBefore(ctx context.Context, userID int64, before time.Time, limit int) ([]domain.PostBase, error) {
	var posts []domain.PostBase
	err := DB.DB.WithContext(ctx).
		Where("user_id = ? AND created_at < ?", userID, before).
		Order("created_at DESC").
		Limit(limit).
		Find(&posts).Error
	return posts, err
}
