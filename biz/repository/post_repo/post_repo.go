package post_repo

import (
	"context"
	"time"

	"zetian-personal-website-hertz/biz/domain"
	DB "zetian-personal-website-hertz/biz/repository"

	"gorm.io/gorm"
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


func UpdatePost(ctx context.Context, userID, postID int64, title, content string) error {
	return DB.DB.WithContext(ctx).
		Model(&domain.Post{}).
		Where("id = ? AND user_id = ?", postID, userID).
		Updates(map[string]any{
			"title":   title,
			"content": content,
			//updated_at will be auto updated by gorm
		}).Error
}


//take in user_id and post_id to ensure only the owner can delete the post
func DeletePost(ctx context.Context, userID, postID int64) error {
    tx := DB.DB.WithContext(ctx).
        Where("id = ? AND user_id = ?", postID, userID).
        Delete(&domain.Post{})

    if tx.Error != nil {
        return tx.Error
    }
    if tx.RowsAffected == 0 {
        // 自定义一个 not found / no permission 错误
        return gorm.ErrRecordNotFound
    }
    return nil
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
