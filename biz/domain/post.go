package domain

import (
	"time"
	thrift "zetian-personal-website-hertz/biz/model/post"
)

// Post — database row model
type Post struct {
	ID        int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID    int64     `json:"user_id"`
	SchoolID  int64     `json:"school_id"`
	Title     string    `json:"title" gorm:"type:varchar(255)"`
	Content   string    `json:"content" gorm:"type:text"`
	ViewCount int       `json:"view_count"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// PostWithStats — Post + like/fav count + user interaction
type PostWithStats struct {
	Post          Post  `json:"post"`
	SchoolName	  string `json:"school_name"`
	LikeCount     int   `json:"like_count"`
	FavCount      int   `json:"fav_count"`
	IsLikedByUser bool  `json:"is_liked_by_user"`
	IsFavByUser   bool  `json:"is_fav_by_user"`
}
// -----------------------------------------------------------------------------
// Converters
// -----------------------------------------------------------------------------

// thrift.Post -> domain.Post
// 注意：这里仍然只关心最基础的 Post，不填 school_name / like_count 等
func FromThriftPostToDomainPost(tp thrift.Post) Post {
	createdAt, _ := time.Parse(time.RFC3339Nano, tp.CreatedAt)
	updatedAt, _ := time.Parse(time.RFC3339Nano, tp.UpdatedAt)
	return Post{
		ID:        tp.ID,
		UserID:    tp.UserID,
		SchoolID:  tp.SchoolID,
		Title:     tp.Title,
		Content:   tp.Content,
		ViewCount: int(tp.ViewCount),
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}

// domain.Post -> thrift.Post
// 注意：这里只能填基础字段，school_name/like/fav/user flags 用默认值
func FromDomainPostToThriftPost(p Post) thrift.Post {
	return thrift.Post{
		ID:            p.ID,
		UserID:        p.UserID,
		SchoolID:      p.SchoolID,
		SchoolName:    "",    // 默认空，通常用 PostWithStats 的版本返回带名字的
		Title:         p.Title,
		Content:       p.Content,
		ViewCount:     int32(p.ViewCount),
		LikeCount:     0,
		FavCount:      0,
		IsLikedByUser: false,
		IsFavByUser:   false,
		CreatedAt:     p.CreatedAt.Format(time.RFC3339Nano),
		UpdatedAt:     p.UpdatedAt.Format(time.RFC3339Nano),
	}
}

// []thrift.Post -> []domain.Post
func FromThriftPostListToDomainPostList(tps []thrift.Post) []Post {
	list := make([]Post, len(tps))
	for i, tp := range tps {
		list[i] = FromThriftPostToDomainPost(tp)
	}
	return list
}

// []domain.Post -> []thrift.Post
func FromDomainPostListToThriftPostList(posts []Post) []thrift.Post {
	list := make([]thrift.Post, len(posts))
	for i, p := range posts {
		list[i] = FromDomainPostToThriftPost(p)
	}
	return list
}

// thrift.Post -> domain.PostWithStats
func FromThriftPostToDomainPostWithStats(tp thrift.Post) PostWithStats {
	createdAt, _ := time.Parse(time.RFC3339Nano, tp.CreatedAt)
	updatedAt, _ := time.Parse(time.RFC3339Nano, tp.UpdatedAt)
	return PostWithStats{
		Post: Post{
			ID:        tp.ID,
			UserID:    tp.UserID,
			SchoolID:  tp.SchoolID,
			Title:     tp.Title,
			Content:   tp.Content,
			ViewCount: int(tp.ViewCount),
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		},
		SchoolName:    tp.SchoolName,   // ✅ 带上学校名字
		LikeCount:     int(tp.LikeCount),
		FavCount:      int(tp.FavCount),
		IsLikedByUser: tp.IsLikedByUser,
		IsFavByUser:   tp.IsFavByUser,
	}
}

// domain.PostWithStats -> thrift.Post
func FromDomainPostWithStatsToThriftPost(p PostWithStats) thrift.Post {
	return thrift.Post{
		ID:            p.Post.ID,
		UserID:        p.Post.UserID,
		SchoolID:      p.Post.SchoolID,
		SchoolName:    p.SchoolName, // ✅ 写回 thrift
		Title:         p.Post.Title,
		Content:       p.Post.Content,
		ViewCount:     int32(p.Post.ViewCount),
		LikeCount:     int32(p.LikeCount),
		FavCount:      int32(p.FavCount),
		IsLikedByUser: p.IsLikedByUser,
		IsFavByUser:   p.IsFavByUser,
		CreatedAt:     p.Post.CreatedAt.Format(time.RFC3339Nano),
		UpdatedAt:     p.Post.UpdatedAt.Format(time.RFC3339Nano),
	}
}

// []thrift.Post -> []domain.PostWithStats
func FromThriftPostListToPostWithStatsList(tps []thrift.Post) []PostWithStats {
	list := make([]PostWithStats, len(tps))
	for i, tp := range tps {
		list[i] = FromThriftPostToDomainPostWithStats(tp)
	}
	return list
}

// []domain.PostWithStats -> []thrift.Post
func FromPostWithStatsListToThriftPostList(posts []PostWithStats) []thrift.Post {
	list := make([]thrift.Post, len(posts))
	for i, p := range posts {
		list[i] = FromDomainPostWithStatsToThriftPost(p)
	}
	return list
}