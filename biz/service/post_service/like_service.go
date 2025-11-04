package post_service


import (
	"context"
	"fmt"

	"zetian-personal-website-hertz/biz/repository/post_repo"
	"zetian-personal-website-hertz/biz/repository/post_fav_repo"
)


// Favorite / Unfavorite
// ----------------------------------------------------

// FavoritePost lets a user favorite a post.
// - Idempotent: multiple calls will keep only one row in post_favorites.
func FavoritePost(ctx context.Context, userID, postID int64) error {
	// 1) ensure post exists
	post, err := post_repo.GetPostByID(ctx, postID)
	if err != nil || post == nil {
		return fmt.Errorf("post not found: %d", postID)
	}

	// 2) favorite
	if err := post_fav_repo.FavoritePost(ctx, userID, postID); err != nil {
		return fmt.Errorf("failed to favorite post: %w", err)
	}
	return nil
}

// UnfavoritePost lets a user remove favorite from a post.
// - Idempotent: removing non-existing favorite is ok.
func UnfavoritePost(ctx context.Context, userID, postID int64) error {
	// 1) optional: ensure post exists
	post, err := post_repo.GetPostByID(ctx, postID)
	if err != nil || post == nil {
		return fmt.Errorf("post not found: %d", postID)
	}

	// 2) unfavorite
	if err := post_fav_repo.UnfavoritePost(ctx, userID, postID); err != nil {
		return fmt.Errorf("failed to unfavorite post: %w", err)
	}
	return nil
}