package repository

import (
	"context"
	"zetian-personal-website-hertz/biz/domain"
	"zetian-personal-website-hertz/biz/pkg/db"
)



func CreateUser(ctx context.Context, user *domain.User) error {
	return db.DB.WithContext(ctx).Create(user).Error
}

func UpdateUser(ctx context.Context, user *domain.User) error {
	return db.DB.WithContext(ctx).Updates(user).Error
}
func GetUserByUsername(ctx context.Context, username string) (*domain.User, error) {
    var user domain.User
    err := db.DB.WithContext(ctx).Where("username = ?", username).First(&user).Error
    if err != nil {
        return nil, err
    }
    return &user, nil
}

func GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	err := db.DB.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}