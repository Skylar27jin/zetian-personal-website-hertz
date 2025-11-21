package post_service

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"zetian-personal-website-hertz/biz/domain"
	"zetian-personal-website-hertz/biz/repository/post_fav_repo"
	"zetian-personal-website-hertz/biz/repository/post_like_repo"
	"zetian-personal-website-hertz/biz/repository/post_repo"
	"zetian-personal-website-hertz/biz/repository/school_repo"

	"gorm.io/gorm"
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
func EditPost(ctx context.Context, userID, postID int64, title string, content string) (*domain.Post, error) {

	return nil, nil
}


// DeletePost deletes a post owned by userID and clears likes/favorites.
func DeletePost(ctx context.Context, userID, postID int64) error {

	var wg sync.WaitGroup
	wg.Add(2)
	errChan := make(chan error, 2)

	// 1. 先清理 likes
	go func() {
		defer wg.Done()
		if err := post_like_repo.DeleteLikesByPostID(ctx, postID); err != nil {
			errChan <- err
		}
	}()
	// 2. 再清理 favorites
	go func() {
		defer wg.Done()
		if err := post_fav_repo.DeleteFavoritesByPostID(ctx, postID); err != nil {
			errChan <- err
		}
	}()

	wg.Wait()
	close(errChan)

	for err := range errChan{
		if err != nil{
			return err
		}
	}
	// 3. 最后删除 post 本身（带 userID 条件，保证只能删自己的）
	if err := post_repo.DeletePost(ctx, userID, postID); err != nil {
		// 区分没找到 vs 真错误
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		return err
	}

	return nil
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

	post, err := post_repo.GetPostByID(ctx, postID)
	if err != nil {
		return nil, err
	}
	postsWithStats, err := buildPostWithStatsList(ctx, []domain.Post{*post}, viewerID)
	if err != nil {
		return nil, err
	}
	return &postsWithStats[0], nil
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

    var (
        likeMap   map[int64]int
        favMap    map[int64]int
        schoolMap map[int64]*domain.School
        likedSet  map[int64]bool
        favSet    map[int64]bool
    )

    var wg sync.WaitGroup
    wg.Add(5)

    errChan := make(chan error, 5)

    // 2) 批量统计点赞
    go func() {
        defer wg.Done()
        m, err := post_like_repo.CountLikesBatch(ctx, postIDs)
        if err != nil {
            errChan <- err
            return
        }
        likeMap = m
    }()

    // 3) 批量统计收藏
    go func() {
        defer wg.Done()
        m, err := post_fav_repo.CountFavoritesBatch(ctx, postIDs)
        if err != nil {
            errChan <- err
            return
        }
        favMap = m
    }()

    // 4) 批量查学校（从缓存）
    go func() {
        defer wg.Done()
        m, err := school_repo.GetSchoolsByIDsInCache(schoolIDs)
        if err != nil {
            errChan <- err
            return
        }
        schoolMap = m
    }()

    // 5) viewer 点赞信息
    go func() {
        defer wg.Done()
        if viewerID == -1 {
            return
        }
        m, err := post_like_repo.GetUserLikedPostIDs(ctx, viewerID, postIDs)
        if err != nil {
            errChan <- err
            return
        }
        likedSet = m
    }()

    // 6) viewer 收藏信息
    go func() {
        defer wg.Done()
        if viewerID == -1 {
            return
        }
        m, err := post_fav_repo.GetUserFavoritedPostIDs(ctx, viewerID, postIDs)
        if err != nil {
            errChan <- err
            return
        }
        favSet = m
    }()

    wg.Wait()
    close(errChan)
    for err := range errChan {
        if err != nil {
            return nil, err
        }
    }

    // nil map 这里是安全的，直接用 likeMap[p.ID] 也不会 panic
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
