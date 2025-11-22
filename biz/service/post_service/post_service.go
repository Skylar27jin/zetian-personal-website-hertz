package post_service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"zetian-personal-website-hertz/biz/domain"
	"zetian-personal-website-hertz/biz/repository/post_repo/post_fav_repo"
	"zetian-personal-website-hertz/biz/repository/post_repo/post_like_repo"
	"zetian-personal-website-hertz/biz/repository/post_repo/post_stats_repo"
    "zetian-personal-website-hertz/biz/repository/post_repo/post_base_repo"
	"zetian-personal-website-hertz/biz/repository/school_repo"

	"gorm.io/gorm"
)

///////////////////////////////////////////////////////////////////////////////
// Create / Edit / Delete
///////////////////////////////////////////////////////////////////////////////

// CreatePost creates a new post.
//
// It will:
//   1) Insert into posts (PostBase)
//   2) Insert a corresponding row into post_stats (PostStats)
//   3) Optionally resolve SchoolName from cache
//
// Special behavior:
//   - mediaUrls / tags are stored as JSON strings in DB.
//   - ReplyTo and Location can be nil.
//   - Stats row starts with all zeros.
func CreatePost(
	ctx context.Context,
	userID int64,
	schoolID int64,
	title string,
	content string,
	mediaType string,
	mediaUrls []string,
	location *string,
	tags []string,
	replyTo *int64,
) (*domain.Post, error) {

	mediaUrlsJSON, err := json.Marshal(mediaUrls)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal media urls: %w", err)
	}
	tagsJSON, err := json.Marshal(tags)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal tags: %w", err)
	}

	now := time.Now()

	base := &domain.PostBase{
		UserID:    userID,
		SchoolID:  schoolID,
		Title:     title,
		Content:   content,
		MediaType: mediaType,
		MediaUrls: string(mediaUrlsJSON),
		Location:  location,
		Tags:      string(tagsJSON),
		ReplyTo:   replyTo,
		CreatedAt: now,
		UpdatedAt: now,
	}

	// 1) Create base row
	if err := post_base_repo.CreatePostBase(ctx, base); err != nil {
		return nil, fmt.Errorf("failed to create post: %w", err)
	}

	// 2) Create stats row (all zeros, PostID = base.ID)
	if err := post_stats_repo.CreateEmptyStats(ctx, base.ID); err != nil {
		return nil, fmt.Errorf("failed to create post stats: %w", err)
	}

	// 3) Try to resolve school name from cache (best-effort)
	schoolName := ""
	if s, err := school_repo.GetSchoolByIDInCache(base.SchoolID); err == nil && s != nil {
		if s.ShortName != "" {
			schoolName = s.ShortName
		} else {
			schoolName = s.Name
		}
	}

	// Compose final domain.Post
	p := &domain.Post{
		PostBase: *base,
		PostStats: domain.PostStats{
			PostID: base.ID,
			// all counters remain zero
		},
		SchoolName:    schoolName,
		IsLikedByUser: false,
		IsFavByUser:   false,
	}

	return p, nil
}

// EditPost updates title/content of a post.
// Only the owner (userID) is allowed to edit.
// Stats / user flags are untouched.
func EditPost(
	ctx context.Context,
	userID, postID int64,
	title string,
	content string,
) (*domain.Post, error) {

	if err := post_base_repo.UpdatePostBase(ctx, userID, postID, title, content); err != nil {
		return nil, fmt.Errorf("failed to update post: %w", err)
	}

	// Reload base
	base, err := post_base_repo.GetPostBaseByID(ctx, postID)
	if err != nil {
		return nil, fmt.Errorf("failed to load updated post: %w", err)
	}

	// Load stats (best-effort)
	stats, err := post_stats_repo.GetStats(ctx, postID)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("failed to load post stats: %w", err)
		}
		// stats not found => use zeros
		stats = &domain.PostStats{PostID: postID}
	}

	// Resolve school name (best-effort)
	schoolName := ""
	if s, err := school_repo.GetSchoolByIDInCache(base.SchoolID); err == nil && s != nil {
		if s.ShortName != "" {
			schoolName = s.ShortName
		} else {
			schoolName = s.Name
		}
	}

	return &domain.Post{
		PostBase:      *base,
		PostStats:     *stats,
		SchoolName:    schoolName,
		IsLikedByUser: false, // viewer not known here
		IsFavByUser:   false,
	}, nil
}

// DeletePost deletes a post owned by userID and clears likes/favorites.
// Notes:
//   - Likes / favorites are deleted first (in parallel).
//   - Base row deletion uses (user_id, id) to enforce ownership.
//   - post_stats row is expected to be deleted by FK cascade.
func DeletePost(ctx context.Context, userID, postID int64) error {
	var wg sync.WaitGroup
	errChan := make(chan error, 2)

	// 1) delete likes
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := post_like_repo.DeleteLikesByPostID(ctx, postID); err != nil {
			errChan <- fmt.Errorf("delete likes: %w", err)
		}
	}()

	// 2) delete favorites
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := post_fav_repo.DeleteFavoritesByPostID(ctx, postID); err != nil {
			errChan <- fmt.Errorf("delete favorites: %w", err)
		}
	}()

	wg.Wait()
	close(errChan)
	for err := range errChan {
		if err != nil {
			return err
		}
	}

	// 3) delete post base (ownership enforced)
	if err := post_base_repo.DeletePostBase(ctx, userID, postID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// upper layer (handler) can map this to 404 or "no permission"
			return err
		}
		return fmt.Errorf("failed to delete post: %w", err)
	}
	// post_stats row is deleted by FK cascade
	return nil
}

///////////////////////////////////////////////////////////////////////////////
// Single Post with Stats
///////////////////////////////////////////////////////////////////////////////

// GetPostWithStats returns:
//   - post base
//   - post stats (view/like/fav/etc.)
//   - school_name
//   - IsLikedByUser / IsFavByUser for the given viewer
//
// Special behavior:
//   - If viewerID > 0 and viewerID != authorID, view_count will be incremented by 1 (best-effort).
//   - If stats row is missing, a zero-valued stats object is used.
func GetPostWithStats(
	ctx context.Context,
	postID int64,
	viewerID int64,
) (*domain.Post, error) {

	// 1) load base
	base, err := post_base_repo.GetPostBaseByID(ctx, postID)
	if err != nil {
		return nil, fmt.Errorf("failed to load post: %w", err)
	}

	// 2) load stats
	stats, err := post_stats_repo.GetStats(ctx, postID)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("failed to load post stats: %w", err)
		}
		stats = &domain.PostStats{PostID: postID}
	}

	// 3) auto increment view_count when viewer is not the author
	if viewerID > 0 && viewerID != base.UserID {
		if err := post_stats_repo.IncrementView(ctx, postID, 1); err == nil {
			// best-effort: reflect increment in memory
			stats.ViewCount++
		}
	}

	// 4) resolve school name (from cache)
	schoolName := ""
	if s, err := school_repo.GetSchoolByIDInCache(base.SchoolID); err == nil && s != nil {
		if s.ShortName != "" {
			schoolName = s.ShortName
		} else {
			schoolName = s.Name
		}
	}

	// 5) user interaction flags
	liked := false
	faved := false
	if viewerID > 0 {
		if ok, err := post_like_repo.HasUserLiked(ctx, viewerID, postID); err == nil {
			liked = ok
		}
		if ok, err := post_fav_repo.HasUserFavorited(ctx, viewerID, postID); err == nil {
			faved = ok
		}
	}

	return &domain.Post{
		PostBase:      *base,
		PostStats:     *stats,
		SchoolName:    schoolName,
		IsLikedByUser: liked,
		IsFavByUser:   faved,
	}, nil
}

///////////////////////////////////////////////////////////////////////////////
// Recent Posts: Personal / School
///////////////////////////////////////////////////////////////////////////////

// GetSchoolRecentPosts returns only PostBase (without stats/user flags).
func GetSchoolRecentPosts(
	ctx context.Context,
	schoolID int64,
	beforeStr string,
	limit int,
) ([]domain.PostBase, error) {

	before, err := parseBeforeTime(beforeStr)
	if err != nil {
		return nil, err
	}
	return post_base_repo.ListPostsBySchoolIDBefore(ctx, schoolID, before, limit)
}

// GetSchoolRecentPostsWithStats returns []Post with stats / school name.
// viewerID can be -1 if you don't need IsLikedByUser / IsFavByUser.
func GetSchoolRecentPostsWithStats(
	ctx context.Context,
	schoolID int64,
	viewerID int64,
	beforeStr string,
	limit int,
) ([]domain.Post, error) {

	bases, err := GetSchoolRecentPosts(ctx, schoolID, beforeStr, limit)
	if err != nil {
		return nil, err
	}
	if len(bases) == 0 {
		return []domain.Post{}, nil
	}
	return buildPostListWithStats(ctx, bases, viewerID)
}

// GetPersonalRecentPosts returns only PostBase for a user.
func GetPersonalRecentPosts(
	ctx context.Context,
	userID int64,
	beforeStr string,
	limit int,
) ([]domain.PostBase, error) {

	before, err := parseBeforeTime(beforeStr)
	if err != nil {
		return nil, fmt.Errorf("invalid time format for 'before': %v", err)
	}

	return post_base_repo.ListPostsByUserIDBefore(ctx, userID, before, limit)
}

// GetPersonalRecentPostsWithStats returns:
//   - base fields
//   - stats from post_stats
//   - school_name
//   - viewer's like/fav flags
func GetPersonalRecentPostsWithStats(
	ctx context.Context,
	userID int64,
	viewerID int64,
	beforeStr string,
	limit int,
) ([]domain.Post, error) {

	bases, err := GetPersonalRecentPosts(ctx, userID, beforeStr, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list posts: %v", err)
	}
	if len(bases) == 0 {
		return []domain.Post{}, nil
	}

	return buildPostListWithStats(ctx, bases, viewerID)
}

///////////////////////////////////////////////////////////////////////////////
// Like / Unlike / Favorite / Unfavorite
///////////////////////////////////////////////////////////////////////////////

// LikePost lets a user like a post.
//
// Behavior:
//   - Ensures post exists.
//   - If already liked, it's a no-op (idempotent).
//   - If not liked, insert like row and increment like_count by 1 (best-effort).
//   - There is a small race condition window, but an offline reconciliation
//     job can correct like_count if necessary.
func LikePost(ctx context.Context, userID, postID int64) error {
	// Ensure post exists
	if _, err := post_base_repo.GetPostBaseByID(ctx, postID); err != nil {
		return fmt.Errorf("post not found: %w", err)
	}

	// Check if already liked
	liked, err := post_like_repo.HasUserLiked(ctx, userID, postID)
	if err != nil {
		return fmt.Errorf("failed to check like state: %w", err)
	}
	if liked {
		return nil
	}

	// Insert like
	if err := post_like_repo.LikePost(ctx, userID, postID); err != nil {
		return fmt.Errorf("failed to like post: %w", err)
	}

	// Increment like_count (best-effort)
	_ = post_stats_repo.IncrementLike(ctx, postID, 1)

	return nil
}

// UnlikePost lets a user remove like from a post.
// - If not liked, it's treated as success.
func UnlikePost(ctx context.Context, userID, postID int64) error {
	// Ensure post exists
	if _, err := post_base_repo.GetPostBaseByID(ctx, postID); err != nil {
		return fmt.Errorf("post not found: %w", err)
	}

	// Check if liked
	liked, err := post_like_repo.HasUserLiked(ctx, userID, postID)
	if err != nil {
		return fmt.Errorf("failed to check like state: %w", err)
	}
	if !liked {
		return nil
	}

	// Remove like
	if err := post_like_repo.UnlikePost(ctx, userID, postID); err != nil {
		return fmt.Errorf("failed to unlike post: %w", err)
	}

	// Decrement like_count (best-effort)
	_ = post_stats_repo.IncrementLike(ctx, postID, -1)

	return nil
}

// FavoritePost lets a user favorite a post.
// - Idempotent: multiple calls will keep only one row in post_favorites.
func FavoritePost(ctx context.Context, userID, postID int64) error {
	// Ensure post exists
	if _, err := post_base_repo.GetPostBaseByID(ctx, postID); err != nil {
		return fmt.Errorf("post not found: %w", err)
	}

	// Check if already favorited
	faved, err := post_fav_repo.HasUserFavorited(ctx, userID, postID)
	if err != nil {
		return fmt.Errorf("failed to check favorite state: %w", err)
	}
	if faved {
		return nil
	}

	// Insert favorite
	if err := post_fav_repo.FavoritePost(ctx, userID, postID); err != nil {
		return fmt.Errorf("failed to favorite post: %w", err)
	}

	// Increment fav_count (best-effort)
	_ = post_stats_repo.IncrementFav(ctx, postID, 1)

	return nil
}

// UnfavoritePost lets a user remove favorite from a post.
// - Idempotent: removing non-existing favorite is ok.
func UnfavoritePost(ctx context.Context, userID, postID int64) error {
	// Ensure post exists
	if _, err := post_base_repo.GetPostBaseByID(ctx, postID); err != nil {
		return fmt.Errorf("post not found: %w", err)
	}

	// Check if favorited
	faved, err := post_fav_repo.HasUserFavorited(ctx, userID, postID)
	if err != nil {
		return fmt.Errorf("failed to check favorite state: %w", err)
	}
	if !faved {
		return nil
	}

	// Remove favorite
	if err := post_fav_repo.UnfavoritePost(ctx, userID, postID); err != nil {
		return fmt.Errorf("failed to unfavorite post: %w", err)
	}

	// Decrement fav_count (best-effort)
	_ = post_stats_repo.IncrementFav(ctx, postID, -1)

	return nil
}

///////////////////////////////////////////////////////////////////////////////
// Helpers
///////////////////////////////////////////////////////////////////////////////

// parseBeforeTime parses frontend "before" time string into time.Time.
//   - empty string => now
//   - try RFC3339Nano, fallback to RFC3339.
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

// buildPostListWithStats:
//   - Input: []PostBase
//   - Output: []Post with:
//       * base fields
//       * stats from post_stats
//       * school_name
//       * IsLikedByUser / IsFavByUser for given viewer
//
// Implementation details:
//   - Stats are loaded by looping GetStats per post (N is usually small, e.g. 10/20).
//   - Schools are loaded in a single call from in-memory cache.
//   - Viewer interactions are loaded in batch via GetUserLikedPostIDs / GetUserFavoritedPostIDs.
func buildPostListWithStats(
	ctx context.Context,
	bases []domain.PostBase,
	viewerID int64,
) ([]domain.Post, error) {
	if len(bases) == 0 {
		return []domain.Post{}, nil
	}

	// Collect postIDs / schoolIDs (deduplicated)
	postIDs := make([]int64, 0, len(bases))
	schoolIDs := make([]int64, 0, len(bases))
	schoolIDSet := make(map[int64]struct{})

	for _, b := range bases {
		postIDs = append(postIDs, b.ID)
		if b.SchoolID != 0 {
			if _, ok := schoolIDSet[b.SchoolID]; !ok {
				schoolIDSet[b.SchoolID] = struct{}{}
				schoolIDs = append(schoolIDs, b.SchoolID)
			}
		}
	}

	var (
		statsMap  map[int64]*domain.PostStats
		schoolMap map[int64]*domain.School
		likedSet  map[int64]bool
		favSet    map[int64]bool
	)

	var wg sync.WaitGroup
	errChan := make(chan error, 4)

	// 1) load stats for each post (simple loop; N typically small)
	wg.Add(1)
	go func() {
		defer wg.Done()
		m := make(map[int64]*domain.PostStats, len(postIDs))
		for _, pid := range postIDs {
			s, err := post_stats_repo.GetStats(ctx, pid)
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					m[pid] = &domain.PostStats{PostID: pid}
					continue
				}
				errChan <- fmt.Errorf("get stats for post %d: %w", pid, err)
				return
			}
			m[pid] = s
		}
		statsMap = m
	}()

	// 2) load schools from cache
	wg.Add(1)
	go func() {
		defer wg.Done()
		m, err := school_repo.GetSchoolsByIDsInCache(schoolIDs)
		if err != nil {
			errChan <- fmt.Errorf("get schools from cache: %w", err)
			return
		}
		schoolMap = m
	}()

	// 3) viewer liked set
	wg.Add(1)
	go func() {
		defer wg.Done()
		if viewerID <= 0 {
			// anonymous or no viewer; skip
			likedSet = map[int64]bool{}
			return
		}
		m, err := post_like_repo.GetUserLikedPostIDs(ctx, viewerID, postIDs)
		if err != nil {
			errChan <- fmt.Errorf("get user liked posts: %w", err)
			return
		}
		likedSet = m
	}()

	// 4) viewer favorited set
	wg.Add(1)
	go func() {
		defer wg.Done()
		if viewerID <= 0 {
			favSet = map[int64]bool{}
			return
		}
		m, err := post_fav_repo.GetUserFavoritedPostIDs(ctx, viewerID, postIDs)
		if err != nil {
			errChan <- fmt.Errorf("get user favorited posts: %w", err)
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

	// Compose final result
	res := make([]domain.Post, 0, len(bases))
	for _, b := range bases {
		// stats
		s := statsMap[b.ID]
		if s == nil {
			s = &domain.PostStats{PostID: b.ID}
		}

		// school_name
		schoolName := ""
		if sch, ok := schoolMap[b.SchoolID]; ok && sch != nil {
			if sch.ShortName != "" {
				schoolName = sch.ShortName
			} else {
				schoolName = sch.Name
			}
		}

		res = append(res, domain.Post{
			PostBase:      b,
			PostStats:     *s,
			SchoolName:    schoolName,
			IsLikedByUser: likedSet[b.ID],
			IsFavByUser:   favSet[b.ID],
		})
	}
	return res, nil
}
