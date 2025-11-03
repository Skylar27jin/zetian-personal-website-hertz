package post_repo

import (
	"context"
	"time"

	"zetian-personal-website-hertz/biz/domain"
	DB "zetian-personal-website-hertz/biz/repository"
)


func CreatePost(ctx context.Context, post *domain.Post) error {
	return DB.DB.WithContext(ctx).Create(post).Error
}


func GetPostByID(ctx context.Context, id int64) (*domain.Post, error) {
	var post domain.Post
	err := DB.DB.WithContext(ctx).Where("id = ?", id).First(&post).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}


func UpdatePost(ctx context.Context, post *domain.Post) error {
	return DB.DB.WithContext(ctx).Save(post).Error
}


func DeletePost(ctx context.Context, id int64) error {
	return DB.DB.WithContext(ctx).Where("id = ?", id).Delete(&domain.Post{}).Error
}


func ListRecentPosts(ctx context.Context, limit int) ([]domain.Post, error) {
	var posts []domain.Post
	err := DB.DB.WithContext(ctx).
		Order("created_at DESC").
		Limit(limit).
		Find(&posts).Error
	return posts, err
}


func ListPostsBySchoolIDBefore(ctx context.Context, schoolID int64, before time.Time, limit int) ([]domain.Post, error) {
	var posts []domain.Post
	err := DB.DB.WithContext(ctx).
		Where("school_id = ? AND created_at < ?", schoolID, before).
		Order("created_at DESC").
		Limit(limit).
		Find(&posts).Error
	return posts, err
}


func ListPostsByUserIDBefore(ctx context.Context, userID int64, before time.Time, limit int) ([]domain.Post, error) {
	var posts []domain.Post
	err := DB.DB.WithContext(ctx).
		Where("user_id = ? AND created_at < ?", userID, before).
		Order("created_at DESC").
		Limit(limit).
		Find(&posts).Error
	return posts, err
}
