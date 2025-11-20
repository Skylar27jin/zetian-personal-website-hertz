package post_service

import (
	"context"
	"fmt"
	"time"

	"zetian-personal-website-hertz/biz/domain"
	"zetian-personal-website-hertz/biz/repository/post_fav_repo"
	"zetian-personal-website-hertz/biz/repository/post_like_repo"
	"zetian-personal-website-hertz/biz/repository/post_repo"
	"zetian-personal-website-hertz/biz/repository/school_repo"
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
// 单个帖子：GetPostWithStats
// ----------------------------------------------------

// GetPostWithStats 返回：基本 Post + school_name + like/fav 数量 + viewer 是否点赞/收藏
func GetPostWithStats(
	ctx context.Context,
	postID int64,
	viewerID int64,
) (*domain.PostWithStats, error) {

	// 1. 基础 Post
	post, err := post_repo.GetPostByID(ctx, postID)
	if err != nil {
		return nil, err
	}

	// 2. 学校信息（填充 SchoolName）
	schoolName := ""
	if post.SchoolID != 0 {
		if s, err := school_repo.GetSchoolByID(ctx, post.SchoolID); err == nil && s != nil {
			// 你可以决定展示 ShortName 还是 Name，这里优先 ShortName
			if s.ShortName != "" {
				schoolName = s.ShortName
			} else {
				schoolName = s.Name
			}
		}
	}

	// 3. 点赞 / 收藏数量
	likeCount, err := post_like_repo.CountLikes(ctx, postID)
	if err != nil {
		return nil, err
	}
	favCount, err := post_fav_repo.CountFavorites(ctx, postID)
	if err != nil {
		return nil, err
	}

	result := &domain.PostWithStats{
		Post:          *post,
		SchoolName:    schoolName,
		LikeCount:     likeCount,
		FavCount:      favCount,
		IsLikedByUser: false,
		IsFavByUser:   false,
	}

	// viewerID == -1 就不查用户行为
	if viewerID == -1 {
		return result, nil
	}

	// 4. 当前 viewer 是否点赞 / 收藏
	liked, err := post_like_repo.HasUserLiked(ctx, viewerID, postID)
	if err != nil {
		return nil, err
	}
	faved, err := post_fav_repo.HasUserFavorited(ctx, viewerID, postID)
	if err != nil {
		return nil, err
	}

	result.IsLikedByUser = liked
	result.IsFavByUser = faved

	return result, nil
}

// ----------------------------------------------------
// School Recent Posts
// ----------------------------------------------------

func GetSchoolRecentPosts(ctx context.Context, schoolID int64, beforeStr string, limit int) ([]domain.Post, error) {
	before, err := parseBeforeTime(beforeStr)
	if err != nil {
		return nil, err
	}
	return post_repo.ListPostsBySchoolIDBefore(ctx, schoolID, before, limit)
}

// GetSchoolRecentPostsWithStats: 带 like/fav 数量，用户行为全 false，带 school_name
func GetSchoolRecentPostsWithStats(
	ctx context.Context,
	schoolID int64,
	beforeStr string,
	limit int,
) ([]domain.PostWithStats, error) {
	posts, err := GetSchoolRecentPosts(ctx, schoolID, beforeStr, limit)
	if err != nil {
		return nil, err
	}
	// viewerID = -1 -> 不填充用户行为
	return buildPostWithStatsList(ctx, posts, -1)
}

// ----------------------------------------------------
// Personal Recent Posts
// ----------------------------------------------------

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

// GetPersonalRecentPostsWithStats
// - 带 like/fav 数量
// - 带当前 viewer 的点赞 / 收藏标记
// - 带 school_name
func GetPersonalRecentPostsWithStats(
	ctx context.Context,
	userID int64,
	viewerID int64,
	beforeStr string,
	limit int,
) ([]domain.PostWithStats, error) {
	posts, err := GetPersonalRecentPosts(ctx, userID, beforeStr, limit)
	if err != nil || len(posts) == 0 {
		return nil, err
	}
	return buildPostWithStatsList(ctx, posts, viewerID)
}

// ----------------------------------------------------
// 通用工具函数
// ----------------------------------------------------

// parse time string from frontend as time.Time
// a valid time format is: "2025-11-03T16:16:17.251927Z"
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

// buildPostWithStatsList:
// - 输入：[]domain.Post
// - 输出：[]domain.PostWithStats，包含 school_name / like_count / fav_count / user flags
func buildPostWithStatsList(
	ctx context.Context,
	posts []domain.Post,
	viewerID int64,
) ([]domain.PostWithStats, error) {
	if len(posts) == 0 {
		return []domain.PostWithStats{}, nil
	}

	// 1) 收集 postIDs / schoolIDs
	postIDs := make([]int64, 0, len(posts))
	schoolIDs := make([]int64, 0, len(posts))
	schoolIDSet := make(map[int64]struct{})

	for _, p := range posts {
		postIDs = append(postIDs, p.ID)

		if p.SchoolID != 0 {
			if _, ok := schoolIDSet[p.SchoolID]; !ok {
				schoolIDSet[p.SchoolID] = struct{}{}
				schoolIDs = append(schoolIDs, p.SchoolID)
			}
		}
	}

	// 2) 批量统计点赞/收藏数量
	likeMap, err := post_like_repo.CountLikesBatch(ctx, postIDs)
	if err != nil {
		return nil, err
	}
	favMap, err := post_fav_repo.CountFavoritesBatch(ctx, postIDs)
	if err != nil {
		return nil, err
	}

	// 3) 批量查学校
	schoolMap, err := school_repo.GetSchoolsByIDs(ctx, schoolIDs)
	if err != nil {
		return nil, err
	}

	// 4) viewer 行为（可选）
	likedSet := map[int64]bool{}
	favSet := map[int64]bool{}
	if viewerID != -1 {
		likedSet, err = post_like_repo.GetUserLikedPostIDs(ctx, viewerID, postIDs)
		if err != nil {
			return nil, err
		}
		favSet, err = post_fav_repo.GetUserFavoritedPostIDs(ctx, viewerID, postIDs)
		if err != nil {
			return nil, err
		}
	}

	// 5) 组装结果
	res := make([]domain.PostWithStats, 0, len(posts))
	for _, p := range posts {
		// 学校显示名：优先 ShortName，其次 Name
		schoolName := ""
		if s, ok := schoolMap[p.SchoolID]; ok && s != nil {
			if s.ShortName != "" {
				schoolName = s.ShortName
			} else {
				schoolName = s.Name
			}
		}

		res = append(res, domain.PostWithStats{
			Post:          p,
			SchoolName:    schoolName,
			LikeCount:     likeMap[p.ID],
			FavCount:      favMap[p.ID],
			IsLikedByUser: likedSet[p.ID],
			IsFavByUser:   favSet[p.ID],
		})
	}
	return res, nil
}
