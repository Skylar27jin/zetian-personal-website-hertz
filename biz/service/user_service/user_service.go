package user_service

import (
	"context"
	"errors"
	"fmt"
	"zetian-personal-website-hertz/biz/domain"
	"zetian-personal-website-hertz/biz/pkg/crypto"
	repo "zetian-personal-website-hertz/biz/repository"

	"gorm.io/gorm"
)

func SignUp(ctx context.Context,userName, password, email string) error {


	//first check if the email already exists
	user, err := repo.GetUserByEmail(ctx, email)
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


	err = repo.CreateUser(ctx, user)
	if err != nil {
		return err
	}

	return nil


}