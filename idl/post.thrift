namespace go post

struct Post {
    1: i64 id,
    2: i64 user_id,
    3: i64 school_id,
    4: string title,
    5: string content,
    6: i32 like_count,
    7: i32 fav_count,
    8: i32 view_count,
    9: string created_at,
    10: string updated_at
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

struct GetAllPersonalPostsReq {
    1: i64 user_id;
    2: string before;
    3: i32 limit;
}

struct GetAllPersonalPostsResp {
    1: bool isSuccessful;
    2: string errorMessage;
    3: list<Post> posts;
}
