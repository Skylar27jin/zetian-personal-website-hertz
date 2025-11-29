package user_service

import (
	"context"
	"errors"
	"fmt"
	"zetian-personal-website-hertz/biz/domain"
	"zetian-personal-website-hertz/biz/pkg/crypto"
	"zetian-personal-website-hertz/biz/repository/user_repo"
	"zetian-personal-website-hertz/biz/repository/user_stats_repo"
	"zetian-personal-website-hertz/biz/repository/user_follow_repo"

	"gorm.io/gorm"
)

// GetUserProfile returns combined User + UserStats + relationship flags.
func GetUserProfile(ctx context.Context, viewerID, targetUserID int64) (*domain.UserProfile, error) {
    if targetUserID <= 0 {
        return nil, fmt.Errorf("invalid target user id")
    }

    // 1) Load user base info
    user, err := user_repo.GetUserByID(ctx, targetUserID)
    if err != nil {
        return nil, fmt.Errorf("user not found: %w", err)
    }

    // 2) Load stats
    stats, err := user_stats_repo.GetStats(ctx, targetUserID)
    if err != nil {
        return nil, fmt.Errorf("user stats not found: %w", err)
    }

    // 3) Determine whether viewer follows this user
    isFollowing := false
    if viewerID > 0 && viewerID != targetUserID {
        follow, err := user_follow_repo.IsFollowing(ctx, viewerID, targetUserID)
        if err == nil {
            isFollowing = follow
        }
    }

    // 4) Determine isMe
    isMe := viewerID == targetUserID

    // 5) Compose UserProfile
    profile := &domain.UserProfile{
        Id:                    int64(user.ID),
        UserName:              user.Username,
        AvatarUrl:             user.AvatarUrl,
        FollowersCount:        stats.FollowersCount,
        FollowingCount:        stats.FollowingCount,
        PostLikeReceivedCount: stats.PostLikeReceivedCount,
        IsFollowing:           isFollowing,
        IsMe:                  isMe,
    }

    return profile, nil
}


/*
SignUp registers a new user with given username, password and email, and store the user into db~
validation should be done before calling this function
*/
func SignUp(ctx context.Context,userName, password, email string) error {
	if userName == "" || password == "" || email == "" {
		return fmt.Errorf("username, password and email cannot be empty")
	}

	//first check if the email already exists
	user, err := user_repo.GetUserByEmail(ctx, email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err // 数据库出错
	}
	if user != nil {
		return fmt.Errorf("email already exists")
	}

	//encrypt the password
	hashedPassword, err := crypto.HashPassword(password)
	if err != nil {
		return err
	}
	user = &domain.User{
		Username: userName,
		Password: hashedPassword,
		Email:    email,
	}


	err = user_repo.CreateUser(ctx, user)
	if err != nil {
		return err
	}

	if err := user_stats_repo.CreateEmptyStats(ctx, int64(user.ID)); err != nil {
		return fmt.Errorf("create user stats: %w", err)
	}

	return nil
}

func Login(ctx context.Context, email, password string) (*domain.User, error) {
	user, err := user_repo.GetUserByEmail(ctx, email)
	if err != nil || user == nil {
		return nil, fmt.Errorf("email or password is incorrect")
	}

	IspasswordMatch := crypto.CheckPassword(password, user.Password)

	if !IspasswordMatch {
		return nil, fmt.Errorf("email or password is incorrect")
	}

	return user, nil

}


func GetUserByID(ctx context.Context, id int64) (*domain.User, error) {
	user, err := user_repo.GetUserByID(ctx, id)
	if err != nil{
		return nil, fmt.Errorf("db Error: %v", err.Error())
	}
	
	return user, nil
}

func GetUserWithStats(ctx context.Context, id int64) (*domain.UserWithStats, error) {
    // 1) 先查 user
    user, err := user_repo.GetUserByID(ctx, id)
    if err != nil {
        return nil, fmt.Errorf("db Error (user): %v", err.Error())
    }

    // 2) 再查 stats（可能不存在，就给个默认值）
    stats, err := user_stats_repo.GetStats(ctx, id)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            // 如果某些老用户没有 stats，给一个默认的（防止前端崩）
            stats = &domain.UserStats{
                UserID:                id,
                FollowersCount:        0,
                FollowingCount:        0,
                PostLikeReceivedCount: 0,
            }
        } else {
            return nil, fmt.Errorf("db Error (user_stats): %v", err.Error())
        }
    }

    return &domain.UserWithStats{
        User:  user,
        Stats: stats,
    }, nil
}


func GetUserByUsername(ctx context.Context, name string) (*domain.User, error) {
	user, err := user_repo.GetUserByUsername(ctx, name)
	if err != nil{
		return nil, fmt.Errorf("db Error: %v", err.Error())
	}

	return user, nil
}


func ResetPassword(ctx context.Context, email, newPassword string) error {
	user, err := user_repo.GetUserByEmail(ctx, email)
	if err != nil || user == nil {
		return fmt.Errorf("email does not exist")
	}

	//encrypt the new password
	hashedPassword, err := crypto.HashPassword(newPassword)
	if err != nil {
		return err
	}

	user.Password = hashedPassword

	err = user_repo.UpdateUser(ctx, user)
	if err != nil {
		return err
	}

	return nil
}


func UpdateAvatarURL(ctx context.Context, userID int64, avatarURL string) error {
	user, err := user_repo.GetUserByID(ctx, userID)
	if err != nil || user == nil {
		return fmt.Errorf("email does not exist")
	}


	user.AvatarUrl = avatarURL

	err = user_repo.UpdateUser(ctx, user)
	if err != nil {
		return err
	}

	return nil
}



// FollowUser 让 followerID 关注 followeeID
func FollowUser(ctx context.Context, followerID, followeeID int64) error {
	if followerID <= 0 || followeeID <= 0 {
		return fmt.Errorf("invalid user id")
	}
	if followerID == followeeID {
		return fmt.Errorf("cannot follow yourself")
	}

	// 1) 先判断是否已经关注（幂等）
	isFollowing, err := user_follow_repo.IsFollowing(ctx, followerID, followeeID)
	if err != nil {
		return fmt.Errorf("check follow state failed: %w", err)
	}
	if isFollowing {
		// 已经关注了，直接返回成功
		return nil
	}

	// 2) 插入关注关系
	if err := user_follow_repo.Follow(ctx, followerID, followeeID); err != nil {
		return fmt.Errorf("create follow record failed: %w", err)
	}

	// 3) 更新统计：被关注者 followers_count +1，关注者 following_count +1
	// 这里都当作 best-effort，失败不影响主流程
	_ = user_stats_repo.IncrementFollowers(ctx, followeeID, 1)
	_ = user_stats_repo.IncrementFollowing(ctx, followerID, 1)

	return nil
}


// UnfollowUser 让 followerID 取消关注 followeeID
func UnfollowUser(ctx context.Context, followerID, followeeID int64) error {
	if followerID <= 0 || followeeID <= 0 {
		return fmt.Errorf("invalid user id")
	}
	if followerID == followeeID {
		return fmt.Errorf("cannot unfollow yourself")
	}

	// 1) 先判断是否当前有关注关系（避免计数乱减）
	isFollowing, err := user_follow_repo.IsFollowing(ctx, followerID, followeeID)
	if err != nil {
		return fmt.Errorf("check follow state failed: %w", err)
	}
	if !isFollowing {
		// 本来就没关注，当作成功
		return nil
	}

	// 2) 删除关注关系
	if err := user_follow_repo.Unfollow(ctx, followerID, followeeID); err != nil {
		return fmt.Errorf("delete follow record failed: %w", err)
	}

	// 3) 更新统计：被关注者 followers_count -1，关注者 following_count -1
	_ = user_stats_repo.IncrementFollowers(ctx, followeeID, -1)
	_ = user_stats_repo.IncrementFollowing(ctx, followerID, -1)

	return nil
}


