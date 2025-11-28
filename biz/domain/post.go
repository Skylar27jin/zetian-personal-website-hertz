package domain

import (
	"encoding/json"
	"time"
	thrift "zetian-personal-website-hertz/biz/model/post"
)

// PostBase — database row model
type PostBase struct {
	ID        int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID    int64     `json:"user_id"`
	SchoolID  int64     `json:"school_id"`

	CategoryID int64    `json:"category_id"`

	Title     string    `json:"title" gorm:"type:varchar(255)"`
	Content   string    `json:"content" gorm:"type:text"`

	MediaType string    `json:"media_type" gorm:"type:varchar(50)"`
	MediaUrls string    `json:"media_urls" gorm:"type:text"` // 存 JSON 字符串，[]string

	Location *string    `json:"location" gorm:"type:varchar(255)"`
	Tags     string     `json:"tags" gorm:"type:text"` // 存 JSON 字符串，[]string

	ReplyTo *int64      `json:"reply_to" gorm:"default:null"`

	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`

}

//Post's stats
type PostStats struct {
	PostID        int64 `json:"post_id" gorm:"primaryKey"`
	ViewCount     int32 `json:"view_count"`
	LikeCount     int32 `json:"like_count"`
	FavCount      int32 `json:"fav_count"`
	CommentCount  int32 `json:"comment_count"`
	ShareCount    int32 `json:"share_count"`
	LastCommentAt int64 `json:"last_comment_at"`
	HotScore      int64 `json:"hot_score"`
}

type Post struct {
	PostBase
	PostStats

	SchoolName    string `json:"school_name"`
	CategoryName  string `json:"category_name"`
	UserName      string `json:"user_name" gorm:"-"`
	UserAvatarUrl string `json:"user_avatar_url" gorm:"-"`
	IsLikedByUser bool   `json:"is_liked_by_user"`
	IsFavByUser   bool   `json:"is_fav_by_user"`
}


/*
Converter overview
------------------
We follow a clean separation:

- PostBase:  original DB row (post metadata + content)
- PostStats: aggregated counters (likes, views, etc.)
- Post:      combined structure returned to clients

This file converts between:
    thrift.Post ↔ PostBase + PostStats + Post

Special behaviors are explicitly documented below.
*/

// -----------------------------------------------------------------------------
// Thrift → Domain (PostBase)
// -----------------------------------------------------------------------------

// ToDomainPostBase converts thrift.Post → PostBase.
// Special behavior:
//   - Tags and MediaUrls are marshaled into JSON strings and stored as TEXT.
//   - Time strings (RFC3339Nano) are parsed without strict error handling.
//     If parsing fails, zero values of time.Time are used.
func ToDomainPostBase(tp thrift.Post) PostBase {
	createdAt, _ := time.Parse(time.RFC3339Nano, tp.CreatedAt)
	updatedAt, _ := time.Parse(time.RFC3339Nano, tp.UpdatedAt)

	// tags: []string → string(JSON)
	tagsJson, _ := json.Marshal(tp.Tags)

	// media_urls: []string → string(JSON)
	mediaUrlsJson, _ := json.Marshal(tp.MediaUrls)

	return PostBase{
		ID:        tp.ID,
		UserID:    tp.UserID,
		SchoolID:  tp.SchoolID,

		CategoryID: tp.CategoryID,

		Title:     tp.Title,
		Content:   tp.Content,
		MediaType: tp.MediaType,
		MediaUrls: string(mediaUrlsJson),

		// pointer field: nil means no location
		Location: tp.Location,

		// tags stored as JSON string
		Tags: string(tagsJson),

		// pointer: nil means "not a reply"
		ReplyTo: tp.ReplyTo,

		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}

// -----------------------------------------------------------------------------
// Thrift → Domain (PostStats)
// -----------------------------------------------------------------------------

// ToDomainPostStats converts thrift.Post → PostStats.
// Note: all counters fit within int32 / int64 as defined in Thrift.
// No precision loss occurs here.
func ToDomainPostStats(tp thrift.Post) PostStats {
	return PostStats{
		PostID:        tp.ID,
		ViewCount:     tp.ViewCount,
		LikeCount:     tp.LikeCount,
		FavCount:      tp.FavCount,
		CommentCount:  tp.CommentCount,
		ShareCount:    tp.ShareCount,
		LastCommentAt: tp.LastCommentAt,
		HotScore:      tp.HotScore,
	}
}

// -----------------------------------------------------------------------------
// Thrift → Domain (Combined Post)
// -----------------------------------------------------------------------------

// ToDomainPost combines PostBase + PostStats + user interaction flags.
func ToDomainPost(tp thrift.Post) Post {
	return Post{
		PostBase:      ToDomainPostBase(tp),
		PostStats:     ToDomainPostStats(tp),
		SchoolName:    tp.SchoolName,
		IsLikedByUser: tp.IsLikedByUser,
		IsFavByUser:   tp.IsFavByUser,
		UserName:      *tp.UserName,
		UserAvatarUrl: *tp.UserAvatarURL,
	}
}

// -----------------------------------------------------------------------------
// Domain → Thrift (Base + Stats)
// -----------------------------------------------------------------------------

// CombineToThriftPost rebuilds a full thrift.Post.
// Special behavior:
//   - Tags and MediaUrls stored as JSON strings must be unmarshaled back to []string.
//   - Time fields are formatted using RFC3339Nano.
//   - Pointer fields (Location, ReplyTo) are passed through directly.
//   - Missing stats or flags should be passed explicitly by callers.
func CombineToThriftPost(
	base PostBase,
	stats PostStats,
	schoolName string,
	liked bool,
	faved bool,
	userName string,
	UserAvatarUrl string,
	CategoryMame string,
) thrift.Post {

	var tags []string
	var mediaUrls []string

	// JSON string → []string
	_ = json.Unmarshal([]byte(base.Tags), &tags)
	_ = json.Unmarshal([]byte(base.MediaUrls), &mediaUrls)

	return thrift.Post{
		// Base fields
		ID:         base.ID,
		UserID:     base.UserID,

		CategoryID: base.CategoryID,
		CategoryName: CategoryMame,

		SchoolID:   base.SchoolID,
		SchoolName: schoolName,
		Title:      base.Title,
		Content:    base.Content,
		MediaType:  base.MediaType,
		MediaUrls:  mediaUrls,
		Location:   base.Location,
		Tags:       tags,
		ReplyTo:    base.ReplyTo,

		// Time (RFC3339Nano)
		CreatedAt: base.CreatedAt.Format(time.RFC3339Nano),
		UpdatedAt: base.UpdatedAt.Format(time.RFC3339Nano),

		// Aggregation fields
		ViewCount:     stats.ViewCount,
		LikeCount:     stats.LikeCount,
		FavCount:      stats.FavCount,
		CommentCount:  stats.CommentCount,
		ShareCount:    stats.ShareCount,
		LastCommentAt: stats.LastCommentAt,
		HotScore:      stats.HotScore,

		// User interaction flags (not stored in DB)
		IsLikedByUser: liked,
		IsFavByUser:   faved,
		UserName:      &userName,
		UserAvatarURL: &UserAvatarUrl,

	}
}

// -----------------------------------------------------------------------------
// Domain → Thrift (Full Post)
// -----------------------------------------------------------------------------

// DomainPostToThrift is a convenience wrapper for a complete Post struct.
func DomainPostToThrift(p Post) thrift.Post {
	return CombineToThriftPost(
		p.PostBase,
		p.PostStats,
		p.SchoolName,
		p.IsLikedByUser,
		p.IsFavByUser,
		p.UserName,
		p.UserAvatarUrl,
		p.CategoryName,
	)
}

// -----------------------------------------------------------------------------
// List Converters
// -----------------------------------------------------------------------------

// ToDomainPostList converts []thrift.Post → []Post.
func ToDomainPostList(tps []thrift.Post) []Post {
	list := make([]Post, len(tps))
	for i, tp := range tps {
		list[i] = ToDomainPost(tp)
	}
	return list
}

// DomainPostListToThrift converts []Post → []thrift.Post.
func DomainPostListToThrift(posts []Post) []thrift.Post {
	list := make([]thrift.Post, len(posts))
	for i, p := range posts {
		list[i] = DomainPostToThrift(p)
	}
	return list
}



// DomainPostListToThriftPointers converts []Post → []*thrift.Post.
//
// 实现思路：
//   1. 先用 DomainPostListToThrift 得到 []thrift.Post
//   2. 再对这个切片逐个取地址，构造 []*thrift.Post
func DomainPostListToThriftPointers(posts []Post) []*thrift.Post {
	// 先得到值类型的列表
	thriftList := DomainPostListToThrift(posts)

	// 再构造指针切片
	ptrs := make([]*thrift.Post, len(thriftList))
	for i := range thriftList {
		ptrs[i] = &thriftList[i]
	}
	return ptrs
}


// DomainPostMapToThriftPointerMap converts map[int64]Post → map[int64]*thrift.Post.
//
// key 保持不变（postID），value 变成对应的 *thrift.Post。
func DomainPostMapToThriftPointerMap(m map[int64]Post) map[int64]*thrift.Post {
	if len(m) == 0 {
		return map[int64]*thrift.Post{}
	}

	res := make(map[int64]*thrift.Post, len(m))
	for id, p := range m {
		tp := DomainPostToThrift(p) // 先得到一个值类型
		cp := tp                    // 拷贝一份，避免取循环变量地址的坑
		res[id] = &cp
	}
	return res
}


