package domain

import (
	"time"
	thrift "zetian-personal-website-hertz/biz/model/post"
)

// Post: the very basic version, which correspond to a row in Posts repo
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
//a higher encaps of Post, with LikeCount and FavCount from post_likes and post_favorites repo
type PostwLikeFav struct {
	Post      Post `json:"post"`       // 帖子主体（数据库映射）
	LikeCount int  `json:"like_count"` // 点赞数量（统计）
	FavCount  int  `json:"fav_count"`  // 收藏数量（统计）
}

type PostwLikeFavAndUser struct {
	PostwLikeFav PostwLikeFav
	IsLikedByUser bool `json:"is_liked_by_user"`
	IsFavByUser   bool `json:"is_fav_by_user"`
}



// FromThriftPostToDomainPost converts thrift.Post -> domain.Post.
//
// 注意：
// - thrift.Post 里有 like_count / fav_count / is_liked_by_user / is_fav_by_user，
//   但 domain.Post 里没有这些字段，这里会直接丢弃。
// - 时间解析失败时，CreatedAt / UpdatedAt 会被置为 time.Time{} (零值)。
func FromThriftPostToDomainPost(tp thrift.Post) Post {
	createdAt, err1 := time.Parse(time.RFC3339Nano, tp.CreatedAt)
	if err1 != nil {
		createdAt = time.Time{} // 默认：零值时间
	}
	updatedAt, err2 := time.Parse(time.RFC3339Nano, tp.UpdatedAt)
	if err2 != nil {
		updatedAt = time.Time{}
	}

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

// FromDomainPostToThriftPost converts domain.Post -> thrift.Post.
//
// 注意：
// - domain.Post 没有 like_count / fav_count / is_liked_by_user / is_fav_by_user，
//   这里会统一填默认值：
//   like_count      = 0
//   fav_count       = 0
//   is_liked_by_user = false
//   is_fav_by_user   = false
func FromDomainPostToThriftPost(p Post) thrift.Post {
	return thrift.Post{
		ID:             p.ID,
		UserID:         p.UserID,
		SchoolID:       p.SchoolID,
		Title:          p.Title,
		Content:        p.Content,
		ViewCount:      int32(p.ViewCount),
		LikeCount:      0,     // 默认值
		FavCount:       0,     // 默认值
		IsLikedByUser:  false, // 默认值
		IsFavByUser:    false, // 默认值
		CreatedAt:      p.CreatedAt.Format(time.RFC3339Nano),
		UpdatedAt:      p.UpdatedAt.Format(time.RFC3339Nano),
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

// FromThriftPostToDomainPostwLikeFav converts thrift.Post -> domain.PostwLikeFav.
//
// 注意：
// - thrift.Post 里的 is_liked_by_user / is_fav_by_user 在 PostwLikeFav 中没有字段，
//   这里会被丢弃（只保留计数）。
func FromThriftPostToDomainPostwLikeFav(tp thrift.Post) PostwLikeFav {
	createdAt, err1 := time.Parse(time.RFC3339Nano, tp.CreatedAt)
	if err1 != nil {
		createdAt = time.Time{}
	}
	updatedAt, err2 := time.Parse(time.RFC3339Nano, tp.UpdatedAt)
	if err2 != nil {
		updatedAt = time.Time{}
	}

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

// FromDomainPostwLikeFavToThriftPost converts domain.PostwLikeFav -> thrift.Post.
//
// 注意：
// - PostwLikeFav 中没有 is_liked_by_user / is_fav_by_user，
//   这里会统一填默认 false。
func FromDomainPostwLikeFavToThriftPost(p PostwLikeFav) thrift.Post {
	return thrift.Post{
		ID:            p.Post.ID,
		UserID:        p.Post.UserID,
		SchoolID:      p.Post.SchoolID,
		Title:         p.Post.Title,
		Content:       p.Post.Content,
		ViewCount:     int32(p.Post.ViewCount),
		LikeCount:     int32(p.LikeCount),
		FavCount:      int32(p.FavCount),
		IsLikedByUser: false, // 默认值
		IsFavByUser:   false, // 默认值
		CreatedAt:     p.Post.CreatedAt.Format(time.RFC3339Nano),
		UpdatedAt:     p.Post.UpdatedAt.Format(time.RFC3339Nano),
	}
}


// []thrift.Post -> []domain.PostwLikeFav
func FromThriftPostListToDomainPostwLikeFavList(tps []thrift.Post) []PostwLikeFav {
	list := make([]PostwLikeFav, len(tps))
	for i, tp := range tps {
		list[i] = FromThriftPostToDomainPostwLikeFav(tp)
	}
	return list
}

// []domain.PostwLikeFav -> []thrift.Post
func FromDomainPostwLikeFavListToThriftPostList(posts []PostwLikeFav) []thrift.Post {
	list := make([]thrift.Post, len(posts))
	for i, p := range posts {
		list[i] = FromDomainPostwLikeFavToThriftPost(p)
	}
	return list
}


// FromThriftPostToDomainPostwLikeFavAndUser
// converts thrift.Post -> domain.PostwLikeFavAndUser.
//
// 对齐关系：
// - tp.like_count      -> LikeCount
// - tp.fav_count       -> FavCount
// - tp.is_liked_by_user -> IsLikedByUser
// - tp.is_fav_by_user   -> IsFavByUser
// - 时间解析失败时：CreatedAt/UpdatedAt = time.Time{} (零值)
func FromThriftPostToDomainPostwLikeFavAndUser(tp thrift.Post) PostwLikeFavAndUser {
	createdAt, err1 := time.Parse(time.RFC3339Nano, tp.CreatedAt)
	if err1 != nil {
		createdAt = time.Time{} // 默认：零值
	}
	updatedAt, err2 := time.Parse(time.RFC3339Nano, tp.UpdatedAt)
	if err2 != nil {
		updatedAt = time.Time{}
	}

	return PostwLikeFavAndUser{
		PostwLikeFav: PostwLikeFav{
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
		},
		IsLikedByUser: tp.IsLikedByUser,
		IsFavByUser:   tp.IsFavByUser,
	}
}



// FromDomainPostwLikeFavAndUserToThriftPost
// converts domain.PostwLikeFavAndUser -> thrift.Post.
//
// 字段一一对应，无需默认值（全部可填满）。
func FromDomainPostwLikeFavAndUserToThriftPost(p PostwLikeFavAndUser) thrift.Post {
	return thrift.Post{
		ID:             p.PostwLikeFav.Post.ID,
		UserID:         p.PostwLikeFav.Post.UserID,
		SchoolID:       p.PostwLikeFav.Post.SchoolID,
		Title:          p.PostwLikeFav.Post.Title,
		Content:        p.PostwLikeFav.Post.Content,
		ViewCount:      int32(p.PostwLikeFav.Post.ViewCount),
		LikeCount:      int32(p.PostwLikeFav.LikeCount),
		FavCount:       int32(p.PostwLikeFav.FavCount),
		IsLikedByUser:  p.IsLikedByUser,
		IsFavByUser:    p.IsFavByUser,
		CreatedAt:      p.PostwLikeFav.Post.CreatedAt.Format(time.RFC3339Nano),
		UpdatedAt:      p.PostwLikeFav.Post.UpdatedAt.Format(time.RFC3339Nano),
	}
}


// []thrift.Post -> []domain.PostwLikeFavAndUser
func FromThriftPostListToDomainPostwLikeFavAndUserList(tps []thrift.Post) []PostwLikeFavAndUser {
	list := make([]PostwLikeFavAndUser, len(tps))
	for i, tp := range tps {
		list[i] = FromThriftPostToDomainPostwLikeFavAndUser(tp)
	}
	return list
}

// []domain.PostwLikeFavAndUser -> []thrift.Post
func FromDomainPostwLikeFavAndUserListToThriftPostList(posts []PostwLikeFavAndUser) []thrift.Post {
	list := make([]thrift.Post, len(posts))
	for i, p := range posts {
		list[i] = FromDomainPostwLikeFavAndUserToThriftPost(p)
	}
	return list
}

