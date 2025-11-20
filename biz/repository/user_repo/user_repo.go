package user_repo

import (
	"context"
	"zetian-personal-website-hertz/biz/domain"
	DB "zetian-personal-website-hertz/biz/repository"
)



func CreateUser(ctx context.Context, user *domain.User) error {
	return DB.DB.WithContext(ctx).Create(user).Error
}

func UpdateUser(ctx context.Context, user *domain.User) error {
	return DB.DB.WithContext(ctx).Updates(user).Error
}

func GetUserByID(ctx context.Context, id int64) (*domain.User, error) {
    var user domain.User
    err := DB.DB.WithContext(ctx).Where("id = ?", id).First(&user).Error
    if err != nil {
        return nil, err
    }
    return &user, nil
}


func GetUserByUsername(ctx context.Context, username string) (*domain.User, error) {
    var user domain.User
    err := DB.DB.WithContext(ctx).Where("username = ?", username).First(&user).Error
    if err != nil {
        return nil, err
    }
    return &user, nil
}

func GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	err := DB.DB.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func GetTableName() string {
	return "users"
}