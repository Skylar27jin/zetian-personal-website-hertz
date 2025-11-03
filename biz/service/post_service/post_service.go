package post_service

import (
	"context"
	"fmt"
	"time"

	"zetian-personal-website-hertz/biz/domain"
	post_favorite_repo "zetian-personal-website-hertz/biz/repository/post_fav_repo.go"
	"zetian-personal-website-hertz/biz/repository/post_like_repo"
	"zetian-personal-website-hertz/biz/repository/post_repo"
)

// CreatePost
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

// EditPost
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

// DeletePost
func DeletePost(ctx context.Context, id int64) error {
	return post_repo.DeletePost(ctx, id)
}

// GetPostByID
func GetPostByID(ctx context.Context, id int64) (*domain.Post, error) {
	return post_repo.GetPostByID(ctx, id)
}

// ----------------------------------------------------
// GetPostwLikeFavAndUserFlags
// - 返回：单个帖子 + like/fav 数量 + 当前 viewer 是否点赞/收藏
// - viewerID == -1 时：IsLikedByUser / IsFavByUser 一律为 false
// ----------------------------------------------------
func GetPostwLikeFavAndUserFlags(
	ctx context.Context,
	postID int64,
	viewerID int64,
) (*domain.PostwLikeFavAndUser, error) {

	// 1. 基础帖子
	post, err := post_repo.GetPostByID(ctx, postID)
	if err != nil {
		return nil, err
	}

	// 2. 点赞 / 收藏数量
	likeCount, err := post_like_repo.CountLikes(ctx, postID)
	if err != nil {
		return nil, err
	}
	favCount, err := post_favorite_repo.CountFavorites(ctx, postID)
	if err != nil {
		return nil, err
	}

	// 3. 先构造带计数的结构，用户相关 flag 先默认 false
	result := &domain.PostwLikeFavAndUser{
		PostwLikeFav: domain.PostwLikeFav{
			Post:      *post,
			LikeCount: likeCount,
			FavCount:  favCount,
		},
		IsLikedByUser: false, // 默认
		IsFavByUser:   false, // 默认
	}

	// 4. 如果 viewerID == -1，就不查“是否点赞/收藏”，直接返回默认 false
	if viewerID == -1 {
		return result, nil
	}

	// 5. 查询当前 viewer 是否点赞/收藏
	liked, err := post_like_repo.HasUserLiked(ctx, viewerID, postID)
	if err != nil {
		return nil, err
	}
	faved, err := post_favorite_repo.HasUserFavorited(ctx, viewerID, postID)
	if err != nil {
		return nil, err
	}

	result.IsLikedByUser = liked
	result.IsFavByUser = faved

	return result, nil
}


// GetSchoolRecentPosts
func GetSchoolRecentPosts(ctx context.Context, schoolID int64, beforeStr string, limit int) ([]domain.Post, error) {
	before, err := parseBeforeTime(beforeStr)
	if err != nil {
		return nil, err
	}
	return post_repo.ListPostsBySchoolIDBefore(ctx, schoolID, before, limit)
}

// GetSchoolRecentPostsWithLikeFav returns posts with like/fav counts.
func GetSchoolRecentPostsWithLikeFav(
	ctx context.Context,
	schoolID int64,
	beforeStr string,
	limit int,
) ([]domain.PostwLikeFav, error) {
	posts, err := GetSchoolRecentPosts(ctx, schoolID, beforeStr, limit)
	if err != nil {
		return nil, err
	}
	return buildPostwLikeFavList(ctx, posts)
}




func GetPersonalRecentPosts(ctx context.Context, userID int64, beforeStr string, limit int) ([]domain.Post, error) {
	before, err := parseBeforeTime(beforeStr)
	if err != nil {
		return nil, fmt.Errorf("invalid time format for 'before': %v", err)
	}

	posts, err := post_repo.ListPostsByUserIDBefore(ctx, userID, before, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list posts: %v", err)
	}
	return posts, nil
}


// GetPersonalRecentPostsWithLikeFavAndUserFlags
// Note: viewerID usually == userID here, but可以区分“看的人”和“发帖的人”.
// if viewerID == -1, then liked_by_user and fav_by_user will be false
func GetPersonalRecentPostsWithLikeFavAndUserFlags(
	ctx context.Context,
	userID int64,
	viewerID int64,
	beforeStr string,
	limit int,
) ([]domain.PostwLikeFavAndUser, error) {
	// 1. get posts + like/fav counts
	postsWithCounts, err := GetPersonalRecentPostsWithLikeFav(ctx, userID, beforeStr, limit)
	if err != nil {
		return nil, err
	}
	if len(postsWithCounts) == 0 {
		return []domain.PostwLikeFavAndUser{}, nil
	}

	// 2. collect post IDs
	postIDs := make([]int64, 0, len(postsWithCounts))
	for _, p := range postsWithCounts {
		postIDs = append(postIDs, p.Post.ID)
	}

	// 3. query whether viewer liked/favorited
	likedSet := map[int64]bool{}
	favSet := map[int64]bool{}
	if viewerID != -1{
		likedSet, err = post_like_repo.GetUserLikedPostIDs(ctx, viewerID, postIDs)
		if err != nil {
			return nil, err
		}
		favSet, err = post_favorite_repo.GetUserFavoritedPostIDs(ctx, viewerID, postIDs)
		if err != nil {
			return nil, err
		}
	} else {
		for _, postWithCount := range postsWithCounts {
			likedSet[postWithCount.Post.ID] = false
			favSet[postWithCount.Post.ID] = false
		}
	}

	// 4. build final result
	res := make([]domain.PostwLikeFavAndUser, 0, len(postsWithCounts))
	for _, p := range postsWithCounts {
		res = append(res, domain.PostwLikeFavAndUser{
			PostwLikeFav: p,
			IsLikedByUser: likedSet[p.Post.ID],
			IsFavByUser:   favSet[p.Post.ID],
		})
	}
	return res, nil
}


// GetPersonalRecentPostsWithLikeFav returns user's posts with like/fav counts.
func GetPersonalRecentPostsWithLikeFav(
	ctx context.Context,
	userID int64,
	beforeStr string,
	limit int,
) ([]domain.PostwLikeFav, error) {
	posts, err := GetPersonalRecentPosts(ctx, userID, beforeStr, limit)
	if err != nil {
		return nil, err
	}
	return buildPostwLikeFavList(ctx, posts)
}




//parse time string from frontend as time.Time
//a valid time format is:"2025-11-03T16:16:17.251927Z"
func parseBeforeTime(beforeStr string) (time.Time, error) {
	if beforeStr == "" {
		return time.Now(), nil
	}
	// try RFC3339Nano first
	if t, err := time.Parse(time.RFC3339Nano, beforeStr); err == nil {
		return t, nil
	}
	// fallback to RFC3339
	return time.Parse(time.RFC3339, beforeStr)
}


//take a list of domain.Post, search their like and fav counts, then returns as []domain.PostwLikeFav
func buildPostwLikeFavList(
	ctx context.Context,
	posts []domain.Post,
) ([]domain.PostwLikeFav, error) {
	if len(posts) == 0 {
		return []domain.PostwLikeFav{}, nil
	}

	postIDs := make([]int64, 0, len(posts))
	for _, p := range posts {
		postIDs = append(postIDs, p.ID)
	}

	likeMap, err := post_like_repo.CountLikesBatch(ctx, postIDs)
	if err != nil {
		return nil, err
	}
	favMap, err := post_favorite_repo.CountFavoritesBatch(ctx, postIDs)
	if err != nil {
		return nil, err
	}

	res := make([]domain.PostwLikeFav, 0, len(posts))
	for _, p := range posts {
		res = append(res, domain.PostwLikeFav{
			Post:      p,
			LikeCount: likeMap[p.ID],
			FavCount:  favMap[p.ID],
		})
	}
	return res, nil
}

