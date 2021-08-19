package user

import (
	"context"

	"github.com/samuael/Project/MaidLink/internal/pkg/model"
)

type IUserService interface {
	GetUserByEmail(context context.Context) *model.User
	GetUserByID(context context.Context) *model.User
	UpdateUser(context.Context) *model.User
	ChangePassword(context context.Context) *model.User
	ChangeImageUrl(context.Context) bool
	GetImageUrl(context.Context) string
	DeleteProfilePicture(context.Context) bool
}
