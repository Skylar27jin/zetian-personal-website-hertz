namespace go post

// RFC3339Nano time strings, like "2025-11-03T00:12:34.123456789Z"
struct Post {
    1: i64 id,
    2: i64 user_id,
    3: i64 school_id,
    4: string school_name,

    5: string title,
    6: string content,

    7: optional string location,   // where is the post created
    8: optional list<string> tags, // hashtags
    9: string media_type,         // "text" / "image" / "video": multi-media post type
    10: list<string> media_urls,   // 图片、视频 URL 列表
    11: optional i64 reply_to,      // reply to what post id

    12: string created_at,
    13: string updated_at,

    14: bool is_liked_by_user,
    15: bool is_fav_by_user,

    16: i32 like_count,
    17: i32 fav_count,
    18: i32 view_count,
    19: i32 comment_count,
    20: i32 share_count,
    21: i64 last_comment_at,       // timestamp of last comment
    22: i64 hot_score,             // sort score for hot posts

    23: optional string user_name,  // 新增


}

//get--------------------------------------------------------------
struct GetPostByIDReq {
    1: i64 id;
}

struct GetPostByIDResp {
    1: bool isSuccessful;
    2: string errorMessage;
    3: optional Post post;
}

//create------------------------------------------------------------
struct CreatePostReq {
    1: i64 user_id,
    2: i64 school_id,
    3: string title,
    4: string content,

    5: optional string location,
    6: optional list<string> tags,
    7: optional string media_type,    // 允许前端不传，后端默认 "text"
    8: optional list<string> media_urls,
    9: optional i64 reply_to,
}

struct CreatePostResp {
    1: bool isSuccessful;
    2: string errorMessage;
    3: Post post;
}


//edit--------------------------------------------------------------
struct EditPostReq {
    1: i64 id;
    2: optional string title;
    3: optional string content;
}

struct EditPostResp {
    1: bool isSuccessful;
    2: string errorMessage;
    3: Post post;
}
//delete------------------------------------------------------------
struct DeletePostReq {
    1: i64 id;
}

struct DeletePostResp {
    1: bool isSuccessful;
    2: string errorMessage;
}

//get -----------------------------------------------------
struct GetSchoolRecentPostsReq {
    1: i64 school_id;
    2: string before;
    3: i32 limit;
}

struct GetSchoolRecentPostsResp {
    1: bool isSuccessful;
    2: string errorMessage;
    3: list<Post> posts;
    4: map<i64, Post> quoted_posts; 
}

struct GetPersonalRecentPostsReq {
    1: i64 user_id;
    2: string before;
    3: i32 limit;
}

struct GetPersonalRecentPostsResp {
    1: bool isSuccessful;
    2: string errorMessage;
    3: list<Post> posts;
    4: map<i64, Post> quoted_posts; 
}

struct LikePostReq {
    1: i64 post_id;
}
struct FavPostReq {
    1: i64 post_id;
}
struct UnlikePostReq {
    1: i64 post_id;
}
struct UnfavPostReq {
    1: i64 post_id;
}

struct UserFlagPostResq {
    1: bool isSuccessful;
    2: string errorMessage;
}

