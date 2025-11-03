package domain

import (
	"time"
	thrift "zetian-personal-website-hertz/biz/model/post"
)

// Post 表示一条帖子
type Post struct {
	ID        int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID    int64     `json:"user_id"`                        // 发帖人ID
	SchoolID  int64     `json:"school_id"`                      // 外键 -> School
	Title     string    `json:"title" gorm:"type:varchar(255)"` // 帖子标题
	Content   string    `json:"content" gorm:"type:text"`       // 帖子正文
	ViewCount int       `json:"view_count"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PostwLikeFav struct {
	Post          Post  `json:"post"`				// 帖子主体（数据库映射）
	LikeCount     int   `json:"like_count"`       	// 点赞数量（统计）
	FavCount      int   `json:"fav_count"`        	// 收藏数量（统计）
}


// thrift.Post → domain.Post
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

// domain.Post → thrift.Post
func FromDomainPostToThriftPost(p Post) thrift.Post {
	return thrift.Post{
		ID:        p.ID,
		UserID:    p.UserID,
		SchoolID:  p.SchoolID,
		Title:     p.Title,
		Content:   p.Content,
		ViewCount: int32(p.ViewCount),
		CreatedAt: p.CreatedAt.Format(time.RFC3339Nano),
		UpdatedAt: p.UpdatedAt.Format(time.RFC3339Nano),
	}
}

// thrift.Post → domain.PostwLikeFav
func FromThriftPostToDomainPostwLikeFav(tp thrift.Post) PostwLikeFav {
	createdAt, _ := time.Parse(time.RFC3339Nano, tp.CreatedAt)
	updatedAt, _ := time.Parse(time.RFC3339Nano, tp.UpdatedAt)

	return PostwLikeFav{
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
		LikeCount: int(tp.LikeCount),
		FavCount:  int(tp.FavCount),
	}
}

// domain.PostwLikeFav → thrift.Post
func FromDomainPostwLikeFavToThriftPost(p PostwLikeFav) thrift.Post {
	return thrift.Post{
		ID:        p.Post.ID,
		UserID:    p.Post.UserID,
		SchoolID:  p.Post.SchoolID,
		Title:     p.Post.Title,
		Content:   p.Post.Content,
		ViewCount: int32(p.Post.ViewCount),
		LikeCount: int32(p.LikeCount),
		FavCount:  int32(p.FavCount),
		CreatedAt: p.Post.CreatedAt.Format(time.RFC3339Nano),
		UpdatedAt: p.Post.UpdatedAt.Format(time.RFC3339Nano),
	}
}


// []thrift.Post → []domain.Post
func FromThriftPostListToDomainPostList(tps []thrift.Post) []Post {
	list := make([]Post, len(tps))
	for i, tp := range tps {
		list[i] = FromThriftPostToDomainPost(tp)
	}
	return list
}

// []domain.Post → []thrift.Post
func FromDomainPostListToThriftPostList(posts []Post) []thrift.Post {
	list := make([]thrift.Post, len(posts))
	for i, p := range posts {
		list[i] = FromDomainPostToThriftPost(p)
	}
	return list
}

// []thrift.Post → []domain.PostwLikeFav
func FromThriftPostListToDomainPostwLikeFavList(tps []thrift.Post) []PostwLikeFav {
	list := make([]PostwLikeFav, len(tps))
	for i, tp := range tps {
		list[i] = FromThriftPostToDomainPostwLikeFav(tp)
	}
	return list
}

// []domain.PostwLikeFav → []thrift.Post
func FromDomainPostwLikeFavListToThriftPostList(posts []PostwLikeFav) []thrift.Post {
	list := make([]thrift.Post, len(posts))
	for i, p := range posts {
		list[i] = FromDomainPostwLikeFavToThriftPost(p)
	}
	return list
}

