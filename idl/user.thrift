namespace go user

// ===============================
// User public profile view
// ===============================

/**
 * Public user profile returned to the client.
 * This combines:
 *   - User.ID
 *   - User.Username
 *   - User.AvatarUrl
 *   - UserStats (followers, following, likes received)
 *
 * Extra fields:
 *   - isFollowing: whether the viewer already follows this user
 *   - isMe: whether the viewer is viewing their own profile
 */
struct UserProfile {
    1: i64    id;
    2: string userName;
    3: string avatarUrl;

    4: i64 followersCount;
    5: i64 followingCount;
    6: i64 postLikeReceivedCount;

    7: bool isFollowing; 
    8: bool isMe;
    9: bool followedYou;        
}

struct SimpleUserProfile {
    1: i64    id;
    2: string userName;
    3: string avatarUrl;

    4: bool isFollowing; 
    5: bool isMe;
    6: bool followedYou; 
}

/**
 * List users who FOLLOW the target user.
 *
 * - viewer is taken from JWT (for isFollowing / followedYou / isMe).
 * - targetUserId: whose followers we are listing.
 * - cursor: opaque offset/cursor for pagination (0 or missing = first page).
 * - limit: page size (optional, server may cap it, e.g. 20).
 */
struct ListFollowersReq {
    1: i64 targetUserId (api.query = "user_id");
    2: i64 cursor       (api.query = "cursor");   // optional, 0 for first page
    3: i32 limit        (api.query = "limit");    // optional, default by server
}

struct ListFollowersResp {
    1: bool   isSuccessful;
    2: string errorMessage;
    // Followers of targetUserId
    3: list<SimpleUserProfile> users;
    // Pagination info
    4: i64  nextCursor;   // 0 when no more data
    5: bool hasMore;
}

/**
 * List users that the target user is FOLLOWING.
 *
 * - viewer is taken from JWT.
 * - targetUserId: whose "following" list we are listing.
 */
struct ListFollowingReq {
    1: i64 targetUserId (api.query = "user_id");
    2: i64 cursor       (api.query = "cursor");   // optional, 0 for first page
    3: i32 limit        (api.query = "limit");    // optional, default by server
}

struct ListFollowingResp {
    1: bool   isSuccessful;
    2: string errorMessage;
    // People that targetUserId is following
    3: list<SimpleUserProfile> users;
    // Pagination info
    4: i64  nextCursor;   // 0 when no more data
    5: bool hasMore;
}

// ===============================
// Get User Info API
// ===============================

/**
 * Request user information (public profile + stats).
 * The authenticated user (viewer) is determined by JWT,
 * and does not need to be passed here.
 */
struct GetUserProfileReq {
    1: i64 id (api.query = "id");
}

struct GetUserProfileResp {
    1: bool   isSuccessful;
    2: string errorMessage;
    3: UserProfile user;
}


// ===============================
// Follow / Unfollow APIs
// ===============================

/**
 * Follow another user.
 * viewerID (the follower) is taken from JWT.
 */
struct FollowUserReq {
    1: i64 targetUserId (api.query = "id");
}

struct FollowUserResp {
    1: bool   isSuccessful;
    2: string errorMessage;
}

/**
 * Unfollow another user.
 * viewerID comes from JWT.
 */
struct UnfollowUserReq {
    1: i64 targetUserId (api.query = "id");
}

struct UnfollowUserResp {
    1: bool   isSuccessful;
    2: string errorMessage;
}


struct LoginReq {
    1: string email (api.body="email");
    2: string password (api.body="password");
}

struct LoginResp {
    1: bool isSuccessful;
    2: string errorMessage;
    3: string userName;
    4: string email;
}

struct SignUpReq {
    1: string username (api.body="username");
    2: string email    (api.body="email");
    3: string password (api.body="password");
}

struct SignUpResp {
    1: bool isSuccessful;
    2: string errorMessage;
    3: string userName;
    4: string email;
}



struct LogoutReq {
}

struct LogoutResp {
    1: bool isSuccessful;
    2: string errorMessage;
}


//get user by ID or name
//cannot pass both ID and name, at least one should be passed
struct GetUserReq {
    1: i64 id (api.query = "id");
    2: string name (api.query = "name");
}


struct GetUserResp {
    1: bool   isSuccessful;
    2: string errorMessage;
    3: string userName;
    4: i64    id;
    5: string avatarUrl;
}


struct ResetPasswordReq {
    1: string email (api.body="email");
    2: string newPassword (api.body="new_password");
}

struct ResetPasswordResp {
    1: bool isSuccessful;
    2: string errorMessage;
}



//place holder
//body里，avatar=要上传的文件
struct UpdateAvatarReq {
}

struct UpdateAvatarResp {
    1: bool isSuccessful;
    2: string errorMessage;
    3: string avatarUrl;
}
