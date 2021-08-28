package user

import (
	"context"

	"github.com/samuael/Project/MaidLink/internal/pkg/model"
)

type IUserService interface {
	GetUserByEmail(context context.Context) *model.User
	CheckEmailExistance(context.Context) bool
	GetUserByID(context context.Context) *model.User
	UpdateUser(context.Context) *model.User
	ChangePassword(context context.Context) *model.User
	ChangeImageUrl(context.Context) bool
	GetImageUrl(context.Context) string
	DeleteProfilePicture(context.Context) bool
	CreateUser(context.Context) *model.User
	// RemoveUser using 'user_id' as a users ID input...
	RemoveUser(context.Context) bool
}
