package post_stats_repo

import (
	"context"

	"zetian-personal-website-hertz/biz/domain"
	DB "zetian-personal-website-hertz/biz/repository"

	"gorm.io/gorm"
)



// -----------------------------------------------------------------------------
// Create
// -----------------------------------------------------------------------------

// CreateEmptyStats creates a stats row for a post after creating PostBase.
// If already exists, FirstOrCreate will simply load existing row.
func CreateEmptyStats(ctx context.Context, postID int64) error {
	stats := domain.PostStats{
		PostID: postID,
		// DB defaults handle all zero values
	}
	return DB.DB.WithContext(ctx).
		FirstOrCreate(&stats, "post_id = ?", postID).Error
}

// -----------------------------------------------------------------------------
// Query
// -----------------------------------------------------------------------------

// GetStats retrieves a PostStats row by post_id.
func GetStats(ctx context.Context, postID int64) (*domain.PostStats, error) {
    var stats domain.PostStats
    err := DB.DB.WithContext(ctx).
        Where("post_id = ?", postID).
        Take(&stats).Error   // <-- 用 Take，避免 ORDER BY
    if err != nil {
        return nil, err
    }
    return &stats, nil
}

func GetStatsBatch(ctx context.Context, postIDs []int64) (map[int64]*domain.PostStats, error) {
	var statsList []*domain.PostStats
	err := DB.DB.WithContext(ctx).
		Where("post_id IN ?", postIDs).
		Find(&statsList).Error
	if err != nil {
		return nil, err
	}
	statsMap := make(map[int64]*domain.PostStats)
	for _, stats := range statsList {
		if stats == nil {
			continue
		}
		statsMap[int64(stats.PostID)] = stats
	}
	return statsMap, nil
}


// -----------------------------------------------------------------------------
// Update
// -----------------------------------------------------------------------------

// UpdateStats updates all mutable fields in PostStats.
func UpdateStats(ctx context.Context, stats *domain.PostStats) error {
	return DB.DB.WithContext(ctx).
		Model(&domain.PostStats{}).
		Where("post_id = ?", stats.PostID).
		Updates(stats).Error
}

// -----------------------------------------------------------------------------
// Atomic Increment Helpers
// -----------------------------------------------------------------------------

// incrementColumn is a shared helper for atomic increments.
// Uses gorm.Expr("col = col + ?") for thread-safe increments.
func incrementColumn(ctx context.Context, postID int64, column string, delta int32) error {
	return DB.DB.WithContext(ctx).
		Model(&domain.PostStats{}).
		Where("post_id = ?", postID).
		Update(column, gorm.Expr(column+" + ?", delta)).
		Error
}

func IncrementView(ctx context.Context, postID int64, delta int32) error {
	return incrementColumn(ctx, postID, "view_count", delta)
}

func IncrementLike(ctx context.Context, postID int64, delta int32) error {
	return incrementColumn(ctx, postID, "like_count", delta)
}

func IncrementFav(ctx context.Context, postID int64, delta int32) error {
	return incrementColumn(ctx, postID, "fav_count", delta)
}

func IncrementComment(ctx context.Context, postID int64, delta int32) error {
	return incrementColumn(ctx, postID, "comment_count", delta)
}

func IncrementShare(ctx context.Context, postID int64, delta int32) error {
	return incrementColumn(ctx, postID, "share_count", delta)
}

// -----------------------------------------------------------------------------
// Delete
// -----------------------------------------------------------------------------

// DeleteStats deletes a stats row.
// Note: This usually won't be needed if foreign key cascade is enabled.
func DeleteStats(ctx context.Context, postID int64) error {
	return DB.DB.WithContext(ctx).
		Where("post_id = ?", postID).
		Delete(&domain.PostStats{}).Error
}
