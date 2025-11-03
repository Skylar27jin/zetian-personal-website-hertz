package post_repo

import (
	"context"
	"time"

	"zetian-personal-website-hertz/biz/domain"
	DB "zetian-personal-website-hertz/biz/repository"
)

// ------------------------------------------------------------
// CreatePost - 创建新帖子
// ------------------------------------------------------------
func CreatePost(ctx context.Context, post *domain.Post) error {
	return DB.DB.WithContext(ctx).Create(post).Error
}

// ------------------------------------------------------------
// GetPostByID - 根据ID获取帖子
// ------------------------------------------------------------
func GetPostByID(ctx context.Context, id int64) (*domain.Post, error) {
	var post domain.Post
	err := DB.DB.WithContext(ctx).Where("id = ?", id).First(&post).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}

// ------------------------------------------------------------
// UpdatePost - 更新帖子（修改标题/内容等）
// ------------------------------------------------------------
func UpdatePost(ctx context.Context, post *domain.Post) error {
	return DB.DB.WithContext(ctx).Save(post).Error
}

// ------------------------------------------------------------
// DeletePost - 删除帖子（物理删除）
// ⚠️ 如需“逻辑删除”，改成 UpdateColumn("deleted", true)
// ------------------------------------------------------------
func DeletePost(ctx context.Context, id int64) error {
	return DB.DB.WithContext(ctx).Where("id = ?", id).Delete(&domain.Post{}).Error
}

// ------------------------------------------------------------
// ListRecentPosts - 获取全站最新帖子（分页）
// ------------------------------------------------------------
func ListRecentPosts(ctx context.Context, limit int) ([]domain.Post, error) {
	var posts []domain.Post
	err := DB.DB.WithContext(ctx).
		Order("created_at DESC").
		Limit(limit).
		Find(&posts).Error
	return posts, err
}

// ------------------------------------------------------------
// ListPostsBySchoolIDBefore - 获取某学校的最新帖子（分页）
// before: 游标时间（ISO 时间解析后）
// limit: 每页数量
// ------------------------------------------------------------
func ListPostsBySchoolIDBefore(ctx context.Context, schoolID int64, before time.Time, limit int) ([]domain.Post, error) {
	var posts []domain.Post
	err := DB.DB.WithContext(ctx).
		Where("school_id = ? AND created_at < ?", schoolID, before).
		Order("created_at DESC").
		Limit(limit).
		Find(&posts).Error
	return posts, err
}

// ------------------------------------------------------------
// ListPostsByUserIDBefore - 获取某用户的所有帖子（分页）
// before: 游标时间
// limit: 每页数量
// ------------------------------------------------------------
func ListPostsByUserIDBefore(ctx context.Context, userID int64, before time.Time, limit int) ([]domain.Post, error) {
	var posts []domain.Post
	err := DB.DB.WithContext(ctx).
		Where("user_id = ? AND created_at < ?", userID, before).
		Order("created_at DESC").
		Limit(limit).
		Find(&posts).Error
	return posts, err
}
