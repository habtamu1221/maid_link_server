package user_service

import (
	"context"

	"github.com/samuael/Project/MaidLink/internal/pkg/model"
	"github.com/samuael/Project/MaidLink/internal/pkg/user"
)

type UserService struct {
	Repo user.IUserRepo
}

func NewUserService(repo user.IUserRepo) user.IUserService {
	return &UserService{
		Repo: repo,
	}
}

func (userser *UserService) GetUser(context context.Context) *model.User {
	user, er := userser.Repo.GetUser(context)
	if er != nil {
		return nil
	}
	user.Password = ""
	return user
}
