package user_service

import (
	"context"

	"github.com/habte/Project/MaidLink/internal/pkg/model"
	"github.com/habte/Project/MaidLink/internal/pkg/user"
)

type UserService struct {
	Repo user.IUserRepo
}

func NewUserService(repo user.IUserRepo) user.IUserService {
	return &UserService{
		Repo: repo,
	}
}

func (userser *UserService) GetUserByEmail(context context.Context) *model.User {
	user, er := userser.Repo.GetUserByEmail(context)
	if er != nil {
		return nil
	}
	return user
}

func (userser *UserService) GetUserByID(context context.Context) *model.User {
	user, err := userser.Repo.GetUserByID(context)
	if err != nil {
		return nil
	}
	return user
}

func (userser *UserService) UpdateUser(context context.Context) *model.User {
	if model, er := userser.Repo.UpdateUser(context); er != nil {
		return nil //
	} else {
		return model
	}
}

func (userser *UserService) ChangePassword(context context.Context) *model.User {
	if model, er := userser.Repo.UpdateUser(context); er != nil {
		return nil //
	} else {
		return model
	}
}

func (userser *UserService) ChangeImageUrl(context context.Context) bool {
	er := userser.Repo.ChangeImageUrl(context)
	return er == nil
}

func (userser *UserService) GetImageUrl(context context.Context) string {
	image, er := userser.Repo.GetImageUrl(context)
	if er != nil {
		return ""
	}
	return image
}

func (userser *UserService) DeleteProfilePicture(contex context.Context) bool {
	ctx := context.WithValue(contex, "image_url", "")                                        // user_id
	ctx = context.WithValue(ctx, "user_id", contex.Value("session").(*model.Session).UserID) //
	era := userser.Repo.ChangeImageUrl(ctx)
	return era == nil
}

// CreateUser ...
func (userser *UserService) CreateUser(context context.Context) *model.User {
	user, er := userser.Repo.CreateUser(context)
	if er == nil {
		return user
	}
	return nil
}

func (userser *UserService) RemoveUser(context context.Context) bool {
	er := userser.Repo.RemoveUser(context)
	return er == nil
}
func (userser *UserService) CheckEmailExistance(context context.Context) bool {
	er := userser.Repo.CheckEmailExistance(context)
	return er == nil
}

func (userser *UserService) DeleteAccount(conte context.Context) bool {
	return userser.Repo.DeleteAccount(conte) == nil
}

func (userser *UserService) ChangeUsername(conte context.Context) bool {
	return userser.Repo.ChangeUsername(conte) == nil
}
