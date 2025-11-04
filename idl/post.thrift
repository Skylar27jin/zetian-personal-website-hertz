namespace go post

// RFC3339Nano time strings, like "2025-11-03T00:12:34.123456789Z"
struct Post {
    1: i64 id,
    2: i64 user_id,
    3: i64 school_id,
    4: string school_name,

    5: string title,
    6: string content,

    7: i32 like_count, //aggregated
    8: i32 fav_count, //aggregated
    9: i32 view_count, //in db

    10: string created_at,
    11: string updated_at,

    12: bool is_liked_by_user,
    13: bool is_fav_by_user,
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

