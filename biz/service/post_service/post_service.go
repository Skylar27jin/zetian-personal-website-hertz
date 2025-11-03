package post_service

import (
	"context"
	"time"

	"zetian-personal-website-hertz/biz/domain"
	"zetian-personal-website-hertz/biz/repository/post_repo"
)

// ----------------------------------------------------
// CreatePost
// ----------------------------------------------------
func CreatePost(ctx context.Context, userID int64, schoolID int64, title string, content string) (*domain.Post, error) {
	newPost := &domain.Post{
		UserID:    userID,
		SchoolID:  schoolID,
		Title:     title,
		Content:   content,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := post_repo.CreatePost(ctx, newPost); err != nil {
		return nil, err
	}
	return newPost, nil
}

// ----------------------------------------------------
// EditPost
// ----------------------------------------------------
func EditPost(ctx context.Context, id int64, title string, content string) (*domain.Post, error) {
	post, err := post_repo.GetPostByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if title != "" {
		post.Title = title
	}
	if content != "" {
		post.Content = content
	}
	post.UpdatedAt = time.Now()

	if err := post_repo.UpdatePost(ctx, post); err != nil {
		return nil, err
	}
	return post, nil
}

// ----------------------------------------------------
// DeletePost
// ----------------------------------------------------
func DeletePost(ctx context.Context, id int64) error {
	return post_repo.DeletePost(ctx, id)
}

// ----------------------------------------------------
// GetPostByID
// ----------------------------------------------------
func GetPostByID(ctx context.Context, id int64) (*domain.Post, error) {
	return post_repo.GetPostByID(ctx, id)
}

// ----------------------------------------------------
// GetSchoolRecentPosts
// ----------------------------------------------------
func GetSchoolRecentPosts(ctx context.Context, schoolID int64, beforeStr string, limit int) ([]domain.Post, error) {
	var before time.Time
	var err error

	if beforeStr == "" {
		before = time.Now()
	} else {
		before, err = time.Parse(time.RFC3339, beforeStr)
		if err != nil {
			return nil, err
		}
	}

	return post_repo.ListPostsBySchoolIDBefore(ctx, schoolID, before, limit)
}

// ----------------------------------------------------
// GetAllPersonalPosts
// ----------------------------------------------------
func GetAllPersonalPosts(ctx context.Context, userID int64, beforeStr string, limit int) ([]domain.Post, error) {
	var before time.Time
	var err error

	if beforeStr == "" {
		before = time.Now()
	} else {
		before, err = time.Parse(time.RFC3339, beforeStr)
		if err != nil {
			return nil, err
		}
	}

	return post_repo.ListPostsByUserIDBefore(ctx, userID, before, limit)
}
