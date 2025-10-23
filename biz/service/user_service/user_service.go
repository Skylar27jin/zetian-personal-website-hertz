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

/*
SignUp registers a new user with given username, password and email, and store the user into db~
validation should be done before calling this function
*/
func SignUp(ctx context.Context,userName, password, email string) error {
	if userName == "" || password == "" || email == "" {
		return fmt.Errorf("username, password and email cannot be empty")
	}

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

func Login(ctx context.Context, email, password string) (*domain.User, error) {
	user, err := repo.GetUserByEmail(ctx, email)
	if err != nil || user == nil {
		return nil, err
	}

	IspasswordMatch := crypto.CheckPassword(password, user.Password)
	if !IspasswordMatch {
		return nil, fmt.Errorf("email or password is incorrect, please try again~")
	}

	return user, nil

}

