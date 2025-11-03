package post_service

import (
	"context"
	"fmt"

	"zetian-personal-website-hertz/biz/repository/post_repo"
	"zetian-personal-website-hertz/biz/repository/post_like_repo"
)



// ----------------------------------------------------
// Like / Unlike
// ----------------------------------------------------

// LikePost lets a user like a post.
// - Idempotent: calling multiple times is safe (repo 使用 ON CONFLICT DO NOTHING).
// - It will return error if the post does not exist.
func LikePost(ctx context.Context, userID, postID int64) error {
	// 1) ensure post exists
	post, err := post_repo.GetPostByID(ctx, postID)
	if err != nil || post == nil {
		return fmt.Errorf("post not found: %d", postID)
	}

	// 2) like
	if err := post_like_repo.LikePost(ctx, userID, postID); err != nil {
		return fmt.Errorf("failed to like post: %w", err)
	}
	return nil
}

// UnlikePost lets a user remove like from a post.
// - Idempotent: deleting non-existing like is still treated as success.
func UnlikePost(ctx context.Context, userID, postID int64) error {
	// 1) optional: ensure post exists (可以不查，但这里保持风格一致)
	post, err := post_repo.GetPostByID(ctx, postID)
	if err != nil || post == nil {
		return fmt.Errorf("post not found: %d", postID)
	}

	// 2) unlike
	if err := post_like_repo.UnlikePost(ctx, userID, postID); err != nil {
		return fmt.Errorf("failed to unlike post: %w", err)
	}
	return nil
}