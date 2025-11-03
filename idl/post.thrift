namespace go post

// RFC3339Nano time strings, like "2025-11-03T00:12:34.123456789Z"
struct Post {
    1:  i64    id,
    2:  i64    user_id,
    3:  i64    school_id,
    4:  string title,
    5:  string content,
    6:  i32    like_count,        // aggregate
    7:  i32    fav_count,         // aggregate
    8:  i32    view_count,
    9:  string created_at,        // RFC3339Nano
    10: string updated_at,        // RFC3339Nano

    // New: per-viewer flags (viewerID==-1 => server should return false)
    11: bool   is_liked_by_user,
    12: bool   is_fav_by_user,

    // (Optional future fields — keep numbers reserved if needed)
    // 13: optional list<string> hashtags,
    // 14: optional string user_name,
    // 15: optional string school_name,
}

//get--------------------------------------------------------------
struct GetPostByIDReq {
    1: i64 id;
}

struct GetPostByIDResp {
    1: bool isSuccessful;
    2: string errorMessage;
    3: Post post;
}

//create------------------------------------------------------------
struct CreatePostReq {
    1: i64 user_id,
    2: i64 school_id,
    3: string title,
    4: string content
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
    3: optional string conten;
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
    4: string oldestTime; //min(posts.created_at), so that frontend is able to eaisly search for the next group of posts
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
}

struct LikePostReq {
    1: i64 user_id;
    2: i64 post_id;
}
struct FavPostReq {
    1: i64 user_id;
    2: i64 post_id;
}
struct UnlikePostReq {
    1: i64 user_id;
    2: i64 post_id;
}
struct UnfavPostReq {
    1: i64 user_id;
    2: i64 post_id;
}


struct UserFlagPostResq {
    1: bool isSuccessful;
    2: string errorMessage;
}

