package post_repo

import (
	"context"
	"time"
	"zetian-personal-website-hertz/biz/domain"
	DB "zetian-personal-website-hertz/biz/repository"
)

// CreatePost
func CreatePost(ctx context.Context, post *domain.Post) error {
	return DB.DB.WithContext(ctx).Create(post).Error
}

// GetPostByID
func GetPostByID(ctx context.Context, id int64) (*domain.Post, error) {
	var post domain.Post
	err := DB.DB.WithContext(ctx).Where("id = ?", id).First(&post).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}

// ListRecentPosts
func ListRecentPosts(ctx context.Context, limit int) ([]domain.Post, error) {
	var posts []domain.Post
	err := DB.DB.WithContext(ctx).
		Order("created_at DESC").
		Limit(limit).
		Find(&posts).Error
	return posts, err
}

//get a school's newest n posts before x, where n is "limit" and x is "before"
func ListPostsBySchoolIDBefore(ctx context.Context, schoolID int64, before time.Time, limit int) ([]domain.Post, error) {
	var posts []domain.Post
	err := DB.DB.WithContext(ctx).
		Where("school_id = ? AND created_at < ?", schoolID, before).
		Order("created_at DESC").
		Limit(limit).
		Find(&posts).Error
	return posts, err
}



