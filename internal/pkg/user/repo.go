package user

import (
	"context"

	"github.com/samuael/Project/MaidLink/internal/pkg/model"
)

type IUserRepo interface {
	GetUserByEmail(context context.Context) (*model.User, error)
	CheckEmailExistance(context.Context) error
	GetUserByID(context context.Context) (*model.User, error)
	UpdateUser(context.Context) (*model.User, error)
	ChangePassword(context context.Context) (*model.User, error)
	ChangeImageUrl(context.Context) error
	GetImageUrl(context.Context) (string, error)
	CreateUser(context.Context) (*model.User, error)
	RemoveUser(context.Context) error
	// DeleteProfilePicture(context.Context) error
}
